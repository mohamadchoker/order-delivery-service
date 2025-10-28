package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/mohamadchoker/order-delivery-service/internal/config"
)

// Connect establishes a connection to PostgreSQL database
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger
	logLevel := gormlogger.Silent
	if cfg.LogSQL {
		logLevel = gormlogger.Info // Shows all SQL queries
	}

	// Create custom logger that outputs to stderr for better visibility in Docker
	gormLogger := gormlogger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags), // Use stderr instead of stdout
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond, // Warn on queries slower than 200ms
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false, // Log "record not found" errors
			Colorful:                  true,  // Colorful output in terminal
			ParameterizedQueries:      false, // Show actual values, not placeholders
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
