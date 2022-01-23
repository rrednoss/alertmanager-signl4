package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleAlertHttpMethods(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		header                 http.Header
		body                   *strings.Reader
		expectedHTTPStatusCode int
	}{
		{
			name:                   "accept POST",
			method:                 http.MethodPost,
			header:                 http.Header{"Content-Type": []string{"application/json"}},
			body:                   strings.NewReader("{\"key\":\"value\"}"),
			expectedHTTPStatusCode: http.StatusOK,
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
			request.Header = tt.header
			response := httptest.NewRecorder()

			handleAlert(response, request)

			got := response.Code
			if got != tt.expectedHTTPStatusCode {
				t.Errorf("got %d, expected %d", got, tt.expectedHTTPStatusCode)
			}
		})
	}
}

func TestHandleAlertPostHeader(t *testing.T) {
	tests := []struct {
		name                   string
		header                 http.Header
		body                   *strings.Reader
		expectedHTTPStatusCode int
	}{
		{
			name:                   "should accept POST with Content-Type header",
			header:                 http.Header{"Content-Type": []string{"application/json"}},
			body:                   strings.NewReader("{\"key\":\"value\"}"),
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "should refuse POST without Content-Type header",
			header:                 nil,
			body:                   strings.NewReader("{\"key\":\"value\"}"),
			expectedHTTPStatusCode: http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/", tt.body)
			request.Header = tt.header
			response := httptest.NewRecorder()

			handleAlert(response, request)

			got := response.Code
			if got != tt.expectedHTTPStatusCode {
				t.Errorf("got %d, expected %d", got, tt.expectedHTTPStatusCode)
			}
		})
	}
}

func TestHandleAlertPostBody(t *testing.T) {
	tests := []struct {
		name                   string
		header                 http.Header
		body                   *strings.Reader
		expectedHTTPStatusCode int
	}{
		{
			name:                   "should accept POST with payload",
			header:                 http.Header{"Content-Type": []string{"application/json"}},
			body:                   strings.NewReader("{\"key\":\"value\"}"),
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "should refuse POST without payload",
			header:                 http.Header{"Content-Type": []string{"application/json"}},
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/", tt.body)
			request.Header = tt.header
			response := httptest.NewRecorder()

			handleAlert(response, request)

			got := response.Code
			if got != tt.expectedHTTPStatusCode {
				t.Errorf("got %d, expected %d", got, tt.expectedHTTPStatusCode)
			}
		})
	}
}
