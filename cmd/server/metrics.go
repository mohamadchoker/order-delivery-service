package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// MetricsServer wraps the Prometheus metrics HTTP server
type MetricsServer struct {
	server *http.Server
	logger *zap.Logger
}

// MetricsConfig holds configuration for the metrics server
type MetricsConfig struct {
	Port   int
	Logger *zap.Logger
}

// NewMetricsServer creates and configures a new metrics server
func NewMetricsServer(cfg MetricsConfig) *MetricsServer {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: promhttp.Handler(),
	}

	return &MetricsServer{
		server: server,
		logger: cfg.Logger,
	}
}

// Start starts the metrics server (blocking)
func (s *MetricsServer) Start() error {
	s.logger.Info("Metrics server listening", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("metrics server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the metrics server
func (s *MetricsServer) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down metrics server...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics server: %w", err)
	}
	s.logger.Info("Metrics server stopped gracefully")
	return nil
}
