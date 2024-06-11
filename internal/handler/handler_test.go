package handler

import (
	"upfcc/internal/aggregator"
	"upfcc/internal/sseclient"
	"upfcc/internal/types"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAnalysisHandler(t *testing.T) {
	tests := []struct {
		name       string
		duration   string
		dimension  string
		wantStatus int
	}{
		{
			name:       "InvalidDuration",
			duration:   "invalid",
			dimension:  string(types.Likes),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "MissingDuration",
			duration:   "",
			dimension:  string(types.Likes),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidDimension",
			duration:   "5s",
			dimension:  "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ValidDimension",
			duration:   "5s",
			dimension:  string(types.Likes),
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockSSEClient{posts: []sseclient.Post{}}
			MockAggregator := &MockAggregator{result: aggregator.AnalysisResult{}}
			handler := New(mockClient, MockAggregator)

			req := httptest.NewRequest("GET", "/analysis?duration="+tt.duration+"&dimension="+tt.dimension, nil)
			rr := httptest.NewRecorder()

			handler.AnalysisHandler(rr, req)

			if gotStatus := rr.Code; gotStatus != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", gotStatus, tt.wantStatus)
			}

		})
	}
}

func TestWriteJSONResponseError(t *testing.T) {
	mockAggregator := &MockAggregator{
		result: aggregator.AnalysisResult{},
	}
	handler := New(nil, mockAggregator)
	rr := httptest.NewRecorder()
	failingWriter := &FailingResponseWriter{rr}

	handler.writeJSONResponse(failingWriter, mockAggregator.result)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

//// helpers

// MockSSEClient simulates an SSE client for testing purposes.
type MockSSEClient struct {
	posts []sseclient.Post
}

// ReadStream simulates reading a stream of posts for the specified duration.
func (m *MockSSEClient) ReadStream(duration time.Duration) chan sseclient.Post {
	postChan := make(chan sseclient.Post)
	go func() {
		defer close(postChan)
		for _, post := range m.posts {
			postChan <- post
		}
	}()
	return postChan
}

// MockAggregator simulates an Aggregator for testing purposes.
type MockAggregator struct {
	result aggregator.AnalysisResult
}

func (m *MockAggregator) AggregateData(duration time.Duration, dimension types.Dimension, resultChan chan aggregator.AnalysisResult) {
	resultChan <- m.result
}

// FailingResponseWriter simulates a ResponseWriter that always returns an error on Write.
type FailingResponseWriter struct {
	http.ResponseWriter
}

func (f *FailingResponseWriter) Write(b []byte) (int, error) {
	return 0, errors.New("simulated write error")
}
