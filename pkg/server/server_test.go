package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodsHandleAlert(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		expectedHTTPStatusCode int
	}{
		{
			name:                   "accept POST",
			method:                 http.MethodPost,
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "accept HEAD",
			method:                 http.MethodHead,
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "refuse DELETE",
			method:                 http.MethodDelete,
			expectedHTTPStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:                   "refuse GET",
			method:                 http.MethodGet,
			expectedHTTPStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:                   "refuse PUT",
			method:                 http.MethodPut,
			expectedHTTPStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(tt.method, "/", nil)
			response := httptest.NewRecorder()

			handleAlert(response, request)

			got := response.Code

			if got != tt.expectedHTTPStatusCode {
				t.Errorf("got %d, expected %d", got, tt.expectedHTTPStatusCode)
			}
		})
	}
}
