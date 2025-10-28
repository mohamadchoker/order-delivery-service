#!/bin/bash

echo "Testing HTTP logging middleware..."
echo ""
echo "Starting demo server on port 9999..."

# Start a simple HTTP server that demonstrates the logging
go run - << 'GOEOF' &
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDHeader = "X-Request-ID"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

func HTTPLoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set(RequestIDHeader, requestID)
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK, written: false}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)
			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.Duration("duration", duration),
				zap.String("request_id", requestID),
			}
			if rw.statusCode >= 500 {
				logger.Error("HTTP request failed", fields...)
			} else if rw.statusCode >= 400 {
				logger.Warn("HTTP request client error", fields...)
			} else {
				logger.Info("HTTP request completed", fields...)
			}
		})
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not found"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"success"}`))
		}
	})
	server := &http.Server{Addr: ":9999", Handler: HTTPLoggingMiddleware(logger)(handler)}
	fmt.Println("Demo server listening on :9999")
	log.Fatal(server.ListenAndServe())
}
GOEOF

SERVER_PID=$!
sleep 2

echo ""
echo "Making test requests..."
echo ""

echo "1. GET request to /v1/deliveries (200 OK):"
curl -s http://localhost:9999/v1/deliveries
echo ""
echo ""

echo "2. GET request to /error (404 Not Found):"
curl -s http://localhost:9999/error
echo ""
echo ""

echo "3. POST request with custom Request-ID:"
curl -s -H "X-Request-ID: my-custom-id-123" -X POST http://localhost:9999/v1/deliveries
echo ""
echo ""

sleep 1
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo "Demo complete! Check the logs above."
