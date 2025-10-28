package main

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/company/order-delivery-service/pkg/metrics"
	"github.com/company/order-delivery-service/pkg/middleware"
	pb "github.com/company/order-delivery-service/proto"
)

// GRPCServer wraps the gRPC server with its dependencies
type GRPCServer struct {
	server       *grpc.Server
	listener     net.Listener
	healthServer *health.Server
	logger       *zap.Logger
}

// GRPCConfig holds configuration for the gRPC server
type GRPCConfig struct {
	Port           int
	RequestTimeout time.Duration
	Logger         *zap.Logger
}

// NewGRPCServer creates and configures a new gRPC server
func NewGRPCServer(cfg GRPCConfig, handler pb.DeliveryServiceServer) (*GRPCServer, error) {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	// Create gRPC server with middleware chain
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDUnaryInterceptor(),
			middleware.TimeoutUnaryInterceptor(cfg.RequestTimeout),
			metrics.MetricsUnaryInterceptor(),
			middleware.LoggingUnaryInterceptor(cfg.Logger),
		),
	)

	// Register business service
	pb.RegisterDeliveryServiceServer(grpcServer, handler)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	return &GRPCServer{
		server:       grpcServer,
		listener:     lis,
		healthServer: healthServer,
		logger:       cfg.Logger,
	}, nil
}

// Start starts the gRPC server (blocking)
func (s *GRPCServer) Start() error {
	s.logger.Info("gRPC server listening", zap.String("address", s.listener.Addr().String()))
	return s.server.Serve(s.listener)
}

// GracefulStop gracefully stops the gRPC server
func (s *GRPCServer) GracefulStop() {
	s.logger.Info("Stopping gRPC server gracefully...")
	s.server.GracefulStop()
	s.logger.Info("gRPC server stopped gracefully")
}

// Stop immediately stops the gRPC server
func (s *GRPCServer) Stop() {
	s.logger.Warn("Force stopping gRPC server...")
	s.server.Stop()
}
