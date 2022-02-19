package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/rrednoss/alertmanager-signl4/pkg/client"
	"github.com/rrednoss/alertmanager-signl4/pkg/config"
)

// Prometheus metric that indicates the number of requests received.
// This metric can be used to understand the total load of this application.
var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of http requests.",
	},
	[]string{"method", "path"},
)

// Prometheus metric that indicates the number of individual status codes.
// This metric can be used to understand how many requests were processed correctly.
var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_status_total",
		Help: "Number of http response status.",
	},
	[]string{"status"},
)

// Prometheus metric that indicates how long requests need to be fullfilled.
// This metric can be used to measure the performance of the app and if the Kubernetes resources need to be adjusted.
var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests",
	}, []string{"path"},
)

func NewServer(alertHandler AlertHandler, healthHandler HealthHandler) *http.Server {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	mux := http.NewServeMux()
	mux.Handle("/v1/alert", alertHandler)
	mux.Handle("/healthz", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:         ":9095",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      mux,
	}
	return s
}

// The AlertHandler receives requests from the alert manager, transforms them depending on the configuration and
// sends them to the Signl4 app afterwards.
type AlertHandler struct {
	config config.AppConfig
	client client.Client
}

func NewAlertHandler(appConfig config.AppConfig, client client.Client) AlertHandler {
	h := AlertHandler{
		config: appConfig,
		client: client,
	}
	return h
}

func (h AlertHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(fmt.Sprintf("received new %s request", r.Method))

	totalRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path))

	switch r.Method {
	case "HEAD":
		h.handleHEAD(w, r)
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, fmt.Sprintf("the HTTP method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
		responseStatus.WithLabelValues(fmt.Sprintf("%d", http.StatusMethodNotAllowed))
	}

	timer.ObserveDuration()
}

func (h AlertHandler) handleHEAD(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	responseStatus.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Inc()
}

func (h AlertHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	if err := h.handlePOSTHeader(w, r); err != nil {
		return
	}
	if err := h.handlePOSTBody(w, r); err != nil {
		http.Error(w, fmt.Sprintf("error processing the request, %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h AlertHandler) handlePOSTHeader(w http.ResponseWriter, r *http.Request) error {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		log.Error("invalid alert request Content-Type header")
		http.Error(w, "invalid Content-Type header", http.StatusUnsupportedMediaType)
		return errors.New("invalid Content-Type header")
	}
	return nil
}

func (h AlertHandler) handlePOSTBody(w http.ResponseWriter, r *http.Request) error {
	alert, err := h.decodeBody(r.Body)
	if err != nil {
		return err
	}
	tAlert, err := transform(h.config.Template, alert)
	if err != nil {
		return err
	}
	status, err := h.determineStatus(alert)
	if err != nil {
		return err
	}
	code, err := h.client.SendAlert(status, strings.NewReader(tAlert))
	if err != nil {
		w.WriteHeader(code)
		return err
	}
	w.WriteHeader(code)
	responseStatus.WithLabelValues(fmt.Sprintf("%d", code)).Inc()
	return nil
}

func (h AlertHandler) decodeBody(body io.ReadCloser) (map[string]interface{}, error) {
	var alert map[string]interface{}

	d := json.NewDecoder(body)
	if err := d.Decode(&alert); err != nil {
		log.Error("couldn't decode alert request")
		return nil, err
	}
	log.Info("decoded alert")
	return alert, nil
}

func (h AlertHandler) determineStatus(alert map[string]interface{}) (client.AlertStatus, error) {
	if v, ok := alert[h.config.StatusKey]; ok {
		if v == "Firing" {
			log.Info("determined alert status firing")
			return client.Firing, nil
		} else if v == "Resolved" {
			log.Info("determined alert status resolving")
			return client.Resolved, nil
		}
	}
	log.Error("couldn't determine alert status")
	return client.Unknown, fmt.Errorf("couldn't determine alert status")
}

// The HealthHandler is used for the Kubernetes Liveness and Readiness probes. It checks for the StatusCode 2xx.
// The output "OK" is intended to increase the comfort on the developer and operator side.
type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("received new HealthHandler request")
	w.Write([]byte("OK"))
}
