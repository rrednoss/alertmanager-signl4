package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/rrednoss/alertmanager-signl4/pkg/config"
)

// Prometheus metric that indicates the number of transformed alerts send to the Signl4 webhook.
// This metric can be used to understand the flow of alerts through this application.
var alertsSend = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "alerts_send_total",
		Help: "Number of alerts send to Signl4 webhook.",
	},
	[]string{},
)

type Client interface {
	SendAlert(status AlertStatus, body io.Reader) (int, error)
}

type AlertStatus int

const (
	Firing AlertStatus = iota
	Resolved
	Unknown
)

type Signl4Client struct {
	Client     *http.Client
	FiringURL  string
	ResolveURL string
}

func NewSignl4Client(config config.AppConfig) Signl4Client {
	prometheus.Register(alertsSend)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return Signl4Client{
		Client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: tr,
		},
		FiringURL:  buildFiringURL(config.TeamSecret),
		ResolveURL: buildResolveURL(config.TeamSecret, config.GroupKey, config.StatusKey),
	}
}

func buildFiringURL(teamSecret string) string {
	return fmt.Sprintf("https://connect.signl4.com/webhook/%s", teamSecret)
}

func buildResolveURL(teamSecret string, groupKey string, statusKey string) string {
	return fmt.Sprintf("https://connect.signl4.com/webhook/%s?ExtIDParam=%s&ExtStatusParam=%s&ResolvedStatus=resolved", teamSecret, groupKey, statusKey)
}

func (sc Signl4Client) SendAlert(status AlertStatus, body io.Reader) (int, error) {
	req, err := http.NewRequest(http.MethodPost, sc.getUrl(status), body)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := sc.Client.Do(req)
	if err != nil {
		log.Error("couldn't send transformed alert")
		return 0, err
	}
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		log.Info("sent transformed alert")
		alertsSend.WithLabelValues().Inc()
		return http.StatusOK, nil
	}
	return http.StatusInternalServerError, nil
}

func (sc Signl4Client) getUrl(status AlertStatus) string {
	if status == Firing {
		return sc.FiringURL
	} else if status == Resolved {
		return sc.ResolveURL
	}
	return ""
}
