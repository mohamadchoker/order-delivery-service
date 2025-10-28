package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds logger configuration options
type Config struct {
	Level            string
	Development      bool
	EnableStacktrace bool
}

// New creates a new logger instance
func New(level string, development bool) (*zap.Logger, error) {
	return NewWithConfig(Config{
		Level:            level,
		Development:      development,
		EnableStacktrace: development, // Default: enable stacktrace only in dev mode
	})
}

// NewWithConfig creates a new logger instance with explicit configuration
func NewWithConfig(cfg Config) (*zap.Logger, error) {
	var zapConfig zap.Config

	if cfg.Development {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Always disable automatic stack traces - we'll control this manually
	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Set log level
	zapLevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zapConfig.Level = zapLevel

	// Build logger with options
	var opts []zap.Option

	// Always add caller info (shows file:line)
	opts = append(opts, zap.AddCaller())

	// Optionally add stack traces
	if cfg.EnableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zapConfig.Build(opts...)
}
