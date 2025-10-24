package constants

import "time"

const (
	// Pagination constants
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
	MinPageSize     = 1

	// Order ID constraints
	OrderIDMinLength = 1
	OrderIDMaxLength = 100

	// Driver ID constraints
	DriverIDMinLength = 1
	DriverIDMaxLength = 100

	// Time constraints
	MinScheduleAdvance  = 30 * time.Minute    // Minimum time before scheduled pickup
	MaxScheduleAdvance  = 30 * 24 * time.Hour // Maximum time for scheduling (30 days)
	MinDeliveryDuration = 15 * time.Minute    // Minimum time between pickup and delivery

	// Database
	DefaultMaxOpenConns    = 25
	DefaultMaxIdleConns    = 5
	DefaultConnMaxLifetime = 5 * time.Minute

	// Context timeouts
	DefaultContextTimeout   = 30 * time.Second
	DatabaseQueryTimeout    = 10 * time.Second
	LongRunningQueryTimeout = 60 * time.Second

	// Request ID
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"

	// Metrics
	MetricsNamespace = "order_delivery"
	MetricsSubsystem = "service"
)

// Resource names for logging and errors
const (
	ResourceDeliveryAssignment = "delivery_assignment"
	ResourceDriver             = "driver"
	ResourceOrder              = "order"
)

// Operation names for domain errors
const (
	OpCreate       = "create"
	OpGet          = "get"
	OpUpdate       = "update"
	OpDelete       = "delete"
	OpList         = "list"
	OpAssignDriver = "assign_driver"
	OpUpdateStatus = "update_status"
	OpGetMetrics   = "get_metrics"
)
