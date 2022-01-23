package server

import (
	"encoding/json"
	"errors"
	"fmt"
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
	var alert map[string]interface{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&alert); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return errors.New("invalid request body")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))

	return nil
}
