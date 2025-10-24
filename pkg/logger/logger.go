package logger

import (
	"go.uber.org/zap"
)

// New creates a new logger instance
func New(level string, development bool) (*zap.Logger, error) {
	var config zap.Config

	if development {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// Set log level
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	config.Level = zapLevel

	return config.Build()
}
