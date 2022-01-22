package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func NewServer() *http.Server {
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
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleHEAD(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 2048)
	if n, err := r.Body.Read(buf); n == 0 && err == io.EOF {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		n := &notification{}
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&n); err != nil && err != io.EOF {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte("Success"))
	}
}

type notification struct {
	Status string  `json:"status"`
	Alerts []alert `json:"alerts"`
}

type alert struct {
	Status string `json:"status"`
	Labels labels `json:"labels"`
}

type labels struct {
	Alertname string `json:"alertname"`
	Severity  string `json:"severity"`
}
