package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:           "valid path /analysis",
			path:           "/analysis",
			expectedStatus: http.StatusOK,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name:           "invalid path",
			path:           "/invalid",
			expectedStatus: http.StatusNotFound,
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := &MockHandler{
				AnalysisHandlerFunc: tt.handlerFunc,
			}
			server := New(mockHandler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Result().StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Result().StatusCode)
			}
		})
	}
}

/////// Helpers

// MockHandler is a mock implementation of the Handler interface.
type MockHandler struct {
	AnalysisHandlerFunc func(w http.ResponseWriter, r *http.Request)
}

func (m *MockHandler) AnalysisHandler(w http.ResponseWriter, r *http.Request) {
	m.AnalysisHandlerFunc(w, r)
}
