package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mohamadchoker/order-delivery-service/pkg/middleware"
	pb "github.com/mohamadchoker/order-delivery-service/proto"
)

// HTTPServer wraps the HTTP/REST gateway server
type HTTPServer struct {
	server *http.Server
	logger *zap.Logger
}

// HTTPConfig holds configuration for the HTTP gateway server
type HTTPConfig struct {
	Port     int
	GRPCPort int
	Logger   *zap.Logger
}

// NewHTTPServer creates and configures a new HTTP gateway server
func NewHTTPServer(ctx context.Context, cfg HTTPConfig) (*HTTPServer, error) {
	// Create gRPC-Gateway mux
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register gateway handlers
	grpcAddress := fmt.Sprintf("localhost:%d", cfg.GRPCPort)
	if err := pb.RegisterDeliveryServiceHandlerFromEndpoint(ctx, gwMux, grpcAddress, opts); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	// Wrap with HTTP logging middleware
	httpHandler := middleware.HTTPLoggingMiddleware(cfg.Logger)(gwMux)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: httpHandler,
	}

	return &HTTPServer{
		server: httpServer,
		logger: cfg.Logger,
	}, nil
}

// Start starts the HTTP gateway server (blocking)
func (s *HTTPServer) Start() error {
	s.logger.Info("HTTP gateway listening", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP gateway...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP gateway: %w", err)
	}
	s.logger.Info("HTTP gateway stopped gracefully")
	return nil
}
