package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendAlert(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
				http.Error(w, "invalid Content-Type header", http.StatusUnsupportedMediaType)
			}

			var body map[string]interface{}
			d := json.NewDecoder(r.Body)
			if err := d.Decode(&body); err != nil {
				http.Error(w, "couldn't decode json payload", http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("successfully received request"))
		}))

	sc := Signl4Client{
		Client:     server.Client(),
		FiringURL:  "https://go.dev/",
		ResolveURL: "https://go.dev/",
	}

	tests := []struct {
		name                   string
		status                 AlertStatus
		body                   io.Reader
		expectedHTTPStatusCode int
	}{
		{
			name:                   "should successfully fire an alert",
			status:                 Firing,
			body:                   strings.NewReader("{\"key\":\"value\"}"),
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "should successfully resolve an alert",
			status:                 Resolved,
			body:                   strings.NewReader("{\"key\":\"value\"}"),
			expectedHTTPStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, err := sc.SendAlert(tt.status, tt.body)
			if err != nil {
				t.Errorf("unexpected error during alert sending %v", err)
			}
			if statusCode != tt.expectedHTTPStatusCode {
				t.Errorf("expected HTTP status code %d, got %d", tt.expectedHTTPStatusCode, statusCode)
			}
		})
	}
}
