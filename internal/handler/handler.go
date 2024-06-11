// Package handler provides HTTP handlers for processing and analyzing social media posts
// data. It utilizes the aggregator package to aggregate data over a specified duration
// and calculate statistical results such as total posts, minimum timestamp, maximum timestamp,
// and average value for a given dimension.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"upfcc/internal/aggregator"
	"upfcc/internal/sseclient"
	"upfcc/internal/types"
)

// Aggregator is an interface that defines the methods required
// for aggregating social media posts data.
type Aggregator interface {
	AggregateData(duration time.Duration, dimension types.Dimension, resultChan chan aggregator.AnalysisResult)
}

// SSEClientInterface defines the interface for an SSE client that reads a stream of posts.
type SSEClientInterface interface {
	ReadStream(duration time.Duration) chan sseclient.Post
}

// Handler is responsible for handling HTTP requests and using the aggregator to process data.
type Handler struct {
	aggregator Aggregator
}

// New creates a new Handler with the provided SSE client.
//
// Parameters:
//   - sseClient: An instance of SSEClientInterface to read the stream of posts.
//
// Returns:
//   - A pointer to the newly created Handler.
func New(sseClient SSEClientInterface, aggregator Aggregator) *Handler {
	return &Handler{
		aggregator: aggregator,
	}
}

// AnalysisHandler handles HTTP requests for analyzing social media posts data.
// It reads the 'duration' and 'dimension' query parameters from the URL, validates them,
// and uses the aggregator to process the data. The results are then returned as a JSON response.
//
// Parameters:
//   - w: An http.ResponseWriter to write the HTTP response.
//   - r: An http.Request representing the HTTP request.
func (h *Handler) AnalysisHandler(w http.ResponseWriter, r *http.Request) {
	duration, err := h.parseDuration(w, r)
	if err != nil {
		return
	}

	dimension, err := h.parseDimension(w, r)
	if err != nil {
		return
	}

	resultChan := make(chan aggregator.AnalysisResult)
	go h.aggregator.AggregateData(duration, dimension, resultChan)

	result := <-resultChan
	h.writeJSONResponse(w, result)
}

// parseDuration reads and parses the 'duration' query parameter from the URL.
// If the parameter is missing or invalid, it writes an HTTP error response.
//
// Parameters:
//   - w: An http.ResponseWriter to write the HTTP response.
//   - r: An http.Request representing the HTTP request.
//
// Returns:
//   - A time.Duration value if parsing is successful.
//   - An error if parsing fails.
func (h *Handler) parseDuration(w http.ResponseWriter, r *http.Request) (time.Duration, error) {
	durationStr := r.URL.Query().Get("duration")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		http.Error(w, "Invalid duration: "+err.Error(), http.StatusBadRequest)
		return 0, err
	}
	return duration, nil
}

// parseDimension reads and validates the 'dimension' query parameter from the URL.
// If the parameter is missing or invalid, it writes an HTTP error response.
//
// Parameters:
//   - w: An http.ResponseWriter to write the HTTP response.
//   - r: An http.Request representing the HTTP request.
//
// Returns:
//   - A types.Dimension value if the dimension is valid.
//   - An error if the dimension is invalid.
func (h *Handler) parseDimension(w http.ResponseWriter, r *http.Request) (types.Dimension, error) {
	dimensionStr := r.URL.Query().Get("dimension")
	dimension := types.Dimension(dimensionStr)
	if !types.IsValidDimension(dimension) {
		http.Error(w, "Invalid dimension: "+dimensionStr, http.StatusBadRequest)
		return "", errors.New("invalid dimension")
	}
	return dimension, nil
}

// writeJSONResponse writes the given result as a JSON response.
//
// Parameters:
//   - w: An http.ResponseWriter to write the HTTP response.
//   - result: The result to write as a JSON response.
func (h *Handler) writeJSONResponse(w http.ResponseWriter, result aggregator.AnalysisResult) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
