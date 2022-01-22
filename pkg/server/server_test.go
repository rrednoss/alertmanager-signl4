package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestStatusCodesHandleAlert(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		body                   *strings.Reader
		expectedHTTPStatusCode int
	}{
		{
			name:                   "accept POST with body",
			method:                 http.MethodPost,
			body:                   strings.NewReader("{'key':'value'}"),
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "refuse POST with empty body",
			method:                 http.MethodPost,
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusBadRequest,
		},
		{
			name:                   "accept HEAD",
			method:                 http.MethodHead,
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "refuse DELETE",
			method:                 http.MethodDelete,
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:                   "refuse GET",
			method:                 http.MethodGet,
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:                   "refuse PUT",
			method:                 http.MethodPut,
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(tt.method, "/", tt.body)
			response := httptest.NewRecorder()

			handleAlert(response, request)

			got := response.Code

			if got != tt.expectedHTTPStatusCode {
				t.Errorf("got %d, expected %d", got, tt.expectedHTTPStatusCode)
			}
		})
	}
}

func TestPostHandleAlert(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		filename               string
		expectedHTTPStatusCode int
	}{
		{
			name:                   "should read valid JSON",
			method:                 http.MethodPost,
			filename:               "example-alert-valid.json",
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "should refuse invalid JSON",
			method:                 http.MethodPost,
			filename:               "example-alert-invalid.json",
			expectedHTTPStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd, err := os.Getwd()
			if err != nil {
				t.Error("can't get working directory", err)
			}

			json, err := os.Open(fmt.Sprintf("%s/%s", wd, tt.filename))
			if err != nil {
				t.Errorf("missing file: %s", tt.filename)
			}

			request, _ := http.NewRequest(tt.method, "/", json)
			response := httptest.NewRecorder()

			handleAlert(response, request)

		})
	}
}
