package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mohamadchoker/order-delivery-service/internal/domain"
)

//go:generate mockgen -destination=../mocks/repository_mock.go -package=mocks  github.com/mohamadchoker/order-delivery-service/internal/service DeliveryRepository

// DeliveryRepository defines the interface for delivery data access.
// This interface belongs to the service layer (Dependency Inversion Principle).
// Concrete implementations (like postgres) depend on this interface, not vice versa.
type DeliveryRepository interface {
	// Create creates a new delivery assignment
	Create(ctx context.Context, assignment *domain.DeliveryAssignment) error

	// GetByID retrieves a delivery assignment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.DeliveryAssignment, error)

	// Update updates an existing delivery assignment
	Update(ctx context.Context, assignment *domain.DeliveryAssignment) error

	// List retrieves delivery assignments with filters and pagination
	List(ctx context.Context, filters ListFilters) ([]*domain.DeliveryAssignment, int64, error)

	// GetMetrics retrieves delivery metrics for a time range
	GetMetrics(ctx context.Context, startTime, endTime time.Time, driverID *string) (*domain.DeliveryMetrics, error)

	// Delete soft-deletes a delivery assignment
	Delete(ctx context.Context, id uuid.UUID) error

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(repo DeliveryRepository) error) error
}

// ListFilters defines filters for listing delivery assignments
type ListFilters struct {
	Page     int
	PageSize int
	Status   *domain.DeliveryStatus
	DriverID *string
}
