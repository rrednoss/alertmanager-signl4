package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rrednoss/alertmanager-signl4/pkg/client"
	"github.com/rrednoss/alertmanager-signl4/pkg/config"
)

type FakeClient struct{}

func (fc FakeClient) SendAlert(status client.AlertStatus, body io.Reader) (int, error) {
	return http.StatusOK, nil
}

func TestAlertHandlerServeHTTP(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		header                 http.Header
		body                   *strings.Reader
		expectedHTTPStatusCode int
	}{
		{
			name:                   "accept valid POST",
			method:                 http.MethodPost,
			header:                 http.Header{"Content-Type": []string{"application/json"}},
			body:                   strings.NewReader("{\"status\":\"Firing\"}"),
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "refuse POST without Content-Type header",
			method:                 http.MethodPost,
			header:                 nil,
			body:                   strings.NewReader("{\"status\":\"Firing\"}"),
			expectedHTTPStatusCode: http.StatusUnsupportedMediaType,
		},
		{
			name:                   "refuse POST without payload",
			method:                 http.MethodPost,
			header:                 http.Header{"Content-Type": []string{"application/json"}},
			body:                   strings.NewReader(""),
			expectedHTTPStatusCode: http.StatusInternalServerError,
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

	alertHandler := NewAlertHandler(config.AppConfig{StatusKey: "status"}, FakeClient{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(tt.method, "/", tt.body)
			request.Header = tt.header
			response := httptest.NewRecorder()

			alertHandler.ServeHTTP(response, request)

			got := response.Code
			if got != tt.expectedHTTPStatusCode {
				t.Errorf("got %d, expected %d", got, tt.expectedHTTPStatusCode)
			}
		})
	}
}
