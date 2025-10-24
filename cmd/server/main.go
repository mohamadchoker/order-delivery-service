package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/company/order-delivery-service/internal/config"
	"github.com/company/order-delivery-service/internal/repository/postgres"
	"github.com/company/order-delivery-service/internal/service"
	grpchandler "github.com/company/order-delivery-service/internal/transport/grpc"
	"github.com/company/order-delivery-service/pkg/logger"
	"github.com/company/order-delivery-service/pkg/metrics"
	"github.com/company/order-delivery-service/pkg/middleware"
	dbpkg "github.com/company/order-delivery-service/pkg/postgres"
	pb "github.com/company/order-delivery-service/proto"
)

// Version information - set via ldflags during build
var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Logger.Level, cfg.Logger.Development)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting order delivery service",
		zap.String("version", version),
		zap.String("build_date", buildDate),
		zap.String("git_commit", gitCommit),
		zap.Int("port", cfg.Server.Port),
	)

	// Connect to database
	db, err := dbpkg.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err := dbpkg.Close(db); err != nil {
			log.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	log.Info("Database connection established")

	// Initialize dependencies
	repo := postgres.NewRepository(db)
	uc := service.NewDeliveryUseCase(repo, log)
	handler := grpchandler.NewHandler(uc, log)

	// Create gRPC server with middleware chain
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDUnaryInterceptor(),
			middleware.TimeoutUnaryInterceptor(30*time.Second),
			metrics.MetricsUnaryInterceptor(),
			loggingInterceptor(log),
		),
	)

	// Register services
	pb.RegisterDeliveryServiceServer(grpcServer, handler)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Start metrics HTTP server
	metricsPort := 9090
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Info("Metrics server listening", zap.Int("port", metricsPort))
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Metrics server error", zap.Error(err))
		}
	}()

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	// Start server in goroutine
	go func() {
		log.Info("gRPC server listening", zap.String("address", lis.Addr().String()))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		log.Warn("Shutdown timeout exceeded, forcing stop")
		grpcServer.Stop()
	case <-stopped:
		log.Info("Server stopped gracefully")
	}

	// Shutdown metrics server
	metricsCtx, metricsCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer metricsCancel()
	if err := metricsServer.Shutdown(metricsCtx); err != nil {
		log.Error("Failed to shutdown metrics server", zap.Error(err))
	}
}

// loggingInterceptor logs gRPC requests with request ID
func loggingInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		requestID := middleware.GetRequestID(ctx)

		// Call handler
		resp, err := handler(ctx, req)

		// Log request
		duration := time.Since(start)
		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("request_id", requestID),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			log.Error("gRPC request failed", fields...)
		} else {
			log.Info("gRPC request completed", fields...)
		}

		return resp, err
	}
}
