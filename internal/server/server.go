// Package server provides an HTTP server implementation that routes
// incoming requests to appropriate handlers based on the URL path.
package server

import (
	"net/http"
)

// Handler is an interface that defines the methods required to handle
// different types of HTTP requests.
type Handler interface {
	AnalysisHandler(w http.ResponseWriter, r *http.Request)
}

// Server represents an HTTP server with a specific handler for processing requests.
type Server struct {
	handler Handler
}

// New creates a new instance of the Server with the given handler.
// The handler is used to process incoming HTTP requests.
func New(handler Handler) *Server {
	return &Server{handler: handler}
}

// ServeHTTP routes incoming HTTP requests to the appropriate handler function
// based on the request URL path. If the URL path matches "/analysis", it invokes
// the AnalysisHandler function of the provided handler. For any other paths, it
// returns a 404 Not Found response.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/analysis":
		s.handler.AnalysisHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
