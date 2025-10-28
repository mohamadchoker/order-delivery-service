package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	" github.com/mohamadchoker/order-delivery-service/internal/domain"
)

//go:generate mockgen -destination=../mocks/usecase_mock.go -package=mocks  github.com/mohamadchoker/order-delivery-service/internal/service DeliveryUseCase

// DeliveryUseCase defines the business logic interface
type DeliveryUseCase interface {
	CreateDeliveryAssignment(ctx context.Context, input CreateDeliveryInput) (*domain.DeliveryAssignment, error)
	GetDeliveryAssignment(ctx context.Context, id uuid.UUID) (*domain.DeliveryAssignment, error)
	UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus, notes string) (*domain.DeliveryAssignment, error)
	ListDeliveryAssignments(ctx context.Context, input ListDeliveryInput) ([]*domain.DeliveryAssignment, int64, error)
	AssignDriver(ctx context.Context, id uuid.UUID, driverID string) (*domain.DeliveryAssignment, error)
	GetDeliveryMetrics(ctx context.Context, startTime, endTime time.Time, driverID *string) (*domain.DeliveryMetrics, error)
	DeleteDeliveryAssignment(ctx context.Context, id uuid.UUID) error
}

// CreateDeliveryInput contains input for creating a delivery assignment
type CreateDeliveryInput struct {
	OrderID               string
	PickupAddress         domain.Address
	DeliveryAddress       domain.Address
	ScheduledPickupTime   time.Time
	EstimatedDeliveryTime time.Time
	Notes                 string
}

// ListDeliveryInput contains input for listing delivery assignments
type ListDeliveryInput struct {
	Page     int
	PageSize int
	Status   *domain.DeliveryStatus
	DriverID *string
}

// deliveryUseCase implements DeliveryUseCase
type deliveryUseCase struct {
	repo   DeliveryRepository
	logger *zap.Logger
}

// NewDeliveryUseCase creates a new delivery use case
func NewDeliveryUseCase(repo DeliveryRepository, logger *zap.Logger) DeliveryUseCase {
	return &deliveryUseCase{
		repo:   repo,
		logger: logger,
	}
}

// CreateDeliveryAssignment creates a new delivery assignment
func (u *deliveryUseCase) CreateDeliveryAssignment(ctx context.Context, input CreateDeliveryInput) (*domain.DeliveryAssignment, error) {
	// Validate input
	if input.OrderID == "" {
		return nil, domain.ErrInvalidInput
	}

	if input.ScheduledPickupTime.IsZero() || input.EstimatedDeliveryTime.IsZero() {
		return nil, domain.ErrInvalidInput
	}

	// Create entity
	assignment := domain.NewDeliveryAssignment(
		input.OrderID,
		input.PickupAddress,
		input.DeliveryAddress,
		input.ScheduledPickupTime,
		input.EstimatedDeliveryTime,
		input.Notes,
	)

	// Save to repository
	if err := u.repo.Create(ctx, assignment); err != nil {
		u.logger.Error("Failed to create delivery assignment",
			zap.Error(err),
			zap.String("order_id", input.OrderID),
		)
		return nil, err
	}

	return assignment, nil
}

// GetDeliveryAssignment retrieves a delivery assignment by ID
func (u *deliveryUseCase) GetDeliveryAssignment(ctx context.Context, id uuid.UUID) (*domain.DeliveryAssignment, error) {
	assignment, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Failed to get delivery assignment",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, err
	}

	return assignment, nil
}

// UpdateDeliveryStatus updates the status of a delivery assignment
func (u *deliveryUseCase) UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status domain.DeliveryStatus, notes string) (*domain.DeliveryAssignment, error) {
	// Get existing assignment
	assignment, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status using domain logic
	if err := assignment.UpdateStatus(status); err != nil {
		u.logger.Error("Failed to update status",
			zap.Error(err),
			zap.String("id", id.String()),
			zap.String("current_status", string(assignment.Status)),
			zap.String("new_status", string(status)),
		)
		return nil, err
	}

	// Update notes if provided
	if notes != "" {
		assignment.Notes = notes
	}

	// Save changes
	if err := u.repo.Update(ctx, assignment); err != nil {
		u.logger.Error("Failed to update delivery assignment",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, err
	}

	return assignment, nil
}

// ListDeliveryAssignments retrieves delivery assignments with pagination
func (u *deliveryUseCase) ListDeliveryAssignments(ctx context.Context, input ListDeliveryInput) ([]*domain.DeliveryAssignment, int64, error) {
	// Set defaults
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	filters := ListFilters(input)

	assignments, totalCount, err := u.repo.List(ctx, filters)
	if err != nil {
		u.logger.Error("Failed to list delivery assignments", zap.Error(err))
		return nil, 0, err
	}

	return assignments, totalCount, nil
}

// AssignDriver assigns a driver to a delivery assignment
func (u *deliveryUseCase) AssignDriver(ctx context.Context, id uuid.UUID, driverID string) (*domain.DeliveryAssignment, error) {
	// Validate driver ID
	if driverID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Get existing assignment
	assignment, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Assign driver using domain logic
	if err := assignment.AssignDriver(driverID); err != nil {
		u.logger.Error("Failed to assign driver",
			zap.Error(err),
			zap.String("id", id.String()),
			zap.String("driver_id", driverID),
		)
		return nil, err
	}

	// Save changes
	if err := u.repo.Update(ctx, assignment); err != nil {
		u.logger.Error("Failed to update delivery assignment",
			zap.Error(err),
			zap.String("id", id.String()),
		)
		return nil, err
	}

	return assignment, nil
}

// GetDeliveryMetrics retrieves delivery metrics
func (u *deliveryUseCase) GetDeliveryMetrics(ctx context.Context, startTime, endTime time.Time, driverID *string) (*domain.DeliveryMetrics, error) {
	// Validate time range
	if startTime.After(endTime) {
		return nil, domain.ErrInvalidInput
	}

	metrics, err := u.repo.GetMetrics(ctx, startTime, endTime, driverID)
	if err != nil {
		u.logger.Error("Failed to get delivery metrics", zap.Error(err))
		return nil, err
	}

	return metrics, nil
}

func (u *deliveryUseCase) DeleteDeliveryAssignment(ctx context.Context, id uuid.UUID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		u.logger.Error("Failed to delete delivery assignment")
		return err
	}

	return nil
}
