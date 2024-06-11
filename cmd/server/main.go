package main

import (
	"upfcc/internal/aggregator"
	"upfcc/internal/handler"
	"upfcc/internal/server"
	"upfcc/internal/sseclient"

	"log"
	"net/http"
)

func main() {
	sseClient := sseclient.New("https://stream.upfluence.co/stream")
	aggregator := aggregator.New(sseClient)
	handler := handler.New(sseClient, aggregator)

	srv := server.New(handler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
