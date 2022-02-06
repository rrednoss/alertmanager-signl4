package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rrednoss/alertmanager-signl4/pkg/client"
	"github.com/rrednoss/alertmanager-signl4/pkg/config"
)

var sc client.Signl4Client = client.NewSignl4Client()

func NewServer() *http.Server {
	// initialize config
	config.Signl4 = config.NewSignl4Config()

	s := &http.Server{
		Addr:         ":9095",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      http.HandlerFunc(handleAlert),
	}
	return s
}

func handleAlert(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "HEAD":
		handleHEAD(w, r)
	case "POST":
		handlePOST(w, r)
	default:
		http.Error(w, fmt.Sprintf("the HTTP method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func handleHEAD(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
	if err := handlePOSTHeader(w, r); err != nil {
		return
	}
	if err := handlePOSTBody(w, r); err != nil {
		http.Error(w, fmt.Sprintf("error processing the request, %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func handlePOSTHeader(w http.ResponseWriter, r *http.Request) error {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(w, "invalid Content-Type header", http.StatusUnsupportedMediaType)
		return errors.New("invalid Content-Type header")
	}
	return nil
}

func handlePOSTBody(w http.ResponseWriter, r *http.Request) error {
	alert, err := decodeBody(r.Body)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	tAlert, err := transform(config.Signl4.Template, alert)
	if err != nil {
		return err
	}
	status, err := determineStatus(alert)
	if err != nil {
		return err
	}
	code, err := sc.SendAlert(status, strings.NewReader(tAlert))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(code)
		return err
	}
	w.WriteHeader(code)

	return nil
}

func decodeBody(body io.ReadCloser) (map[string]interface{}, error) {
	var alert map[string]interface{}

	d := json.NewDecoder(body)
	if err := d.Decode(&alert); err != nil {
		return nil, err
	}
	return alert, nil
}

func determineStatus(alert map[string]interface{}) (client.AlertStatus, error) {
	if v, ok := alert[config.Signl4.StatusKey]; ok {
		if v == "Firing" {
			return client.Firing, nil
		} else if v == "Resolved" {
			return client.Resolved, nil
		}
	}
	return client.Unknown, fmt.Errorf("couldn't determine alert status")
}
