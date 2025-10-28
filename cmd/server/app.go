package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/mohamadchoker/order-delivery-service/internal/config"
	"github.com/mohamadchoker/order-delivery-service/internal/repository/postgres"
	"github.com/mohamadchoker/order-delivery-service/internal/service"
	grpchandler "github.com/mohamadchoker/order-delivery-service/internal/transport/grpc"
	"github.com/mohamadchoker/order-delivery-service/pkg/logger"
	dbpkg "github.com/mohamadchoker/order-delivery-service/pkg/postgres"
)

// App represents the application with all its dependencies
type App struct {
	config *config.Config
	logger *zap.Logger
	db     *gorm.DB

	grpcServer    *GRPCServer
	metricsServer *MetricsServer
}

// NewApp creates a new application instance with all dependencies initialized
func NewApp(version, buildDate, gitCommit string) (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log, err := logger.NewWithConfig(logger.Config{
		Level:            cfg.Logger.Level,
		Development:      cfg.Logger.Development,
		EnableStacktrace: cfg.Logger.EnableStacktrace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	log.Info("Starting order delivery service",
		zap.String("version", version),
		zap.String("build_date", buildDate),
		zap.String("git_commit", gitCommit),
		zap.Int("grpc_port", cfg.Server.Port),
	)

	// Connect to database
	db, err := dbpkg.Connect(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Info("Database connection established")

	// Initialize business layer (dependency injection)
	repo := postgres.NewRepository(db)
	useCase := service.NewDeliveryUseCase(repo, log)
	handler := grpchandler.NewHandler(useCase, log)

	// Create gRPC server
	grpcServer, err := NewGRPCServer(GRPCConfig{
		Port:           cfg.Server.Port,
		RequestTimeout: 30 * time.Second,
		Logger:         log,
	}, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC server: %w", err)
	}

	// Create metrics server
	metricsServer := NewMetricsServer(MetricsConfig{
		Port:   9090, // TODO: Add to config
		Logger: log,
	})

	return &App{
		config:        cfg,
		logger:        log,
		db:            db,
		grpcServer:    grpcServer,
		metricsServer: metricsServer,
	}, nil
}

// Run starts all servers and blocks until shutdown signal is received
func (a *App) Run() error {
	// Start metrics server in background
	go func() {
		if err := a.metricsServer.Start(); err != nil {
			a.logger.Error("Metrics server error", zap.Error(err))
		}
	}()

	// Start gRPC server in background
	go func() {
		if err := a.grpcServer.Start(); err != nil {
			a.logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// Graceful shutdown
	return a.Shutdown()
}

// Shutdown gracefully shuts down all servers and closes resources
func (a *App) Shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown gRPC server with timeout
	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-shutdownCtx.Done():
		a.logger.Warn("Shutdown timeout exceeded, forcing stop")
		a.grpcServer.Stop()
	case <-stopped:
		// Graceful stop completed
	}

	// Shutdown metrics server
	metricsCtx, metricsCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer metricsCancel()
	if err := a.metricsServer.Shutdown(metricsCtx); err != nil {
		a.logger.Error("Failed to shutdown metrics server", zap.Error(err))
	}

	// Close database connection
	if err := dbpkg.Close(a.db); err != nil {
		a.logger.Error("Failed to close database connection", zap.Error(err))
		return err
	}

	// Sync logger
	if err := a.logger.Sync(); err != nil {
		// Ignore sync errors on stderr (common on some systems)
		return nil
	}

	return nil
}
