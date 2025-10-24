package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/company/order-delivery-service/internal/domain"
	"github.com/company/order-delivery-service/internal/repository/postgres/model"
	"github.com/company/order-delivery-service/internal/service"
)

// repository implements service.DeliveryRepository using PostgreSQL
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository
func NewRepository(db *gorm.DB) service.DeliveryRepository {
	return &repository{db: db}
}

// Create creates a new delivery assignment
func (r *repository) Create(ctx context.Context, assignment *domain.DeliveryAssignment) error {
	dbModel := model.FromEntity(assignment)

	if err := r.db.WithContext(ctx).Create(dbModel).Error; err != nil {
		return err
	}

	*assignment = *dbModel.ToEntity()
	return nil
}

// GetByID retrieves a delivery assignment by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DeliveryAssignment, error) {
	var dbModel model.DeliveryAssignment

	if err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return dbModel.ToEntity(), nil
}

// Update updates an existing delivery assignment
func (r *repository) Update(ctx context.Context, assignment *domain.DeliveryAssignment) error {
	dbModel := model.FromEntity(assignment)

	result := r.db.WithContext(ctx).
		Model(&model.DeliveryAssignment{}).
		Where("id = ?", assignment.ID).
		Updates(dbModel)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves delivery assignments with pagination and filters
func (r *repository) List(ctx context.Context, filters service.ListFilters) ([]*domain.DeliveryAssignment, int64, error) {
	var dbModels []model.DeliveryAssignment
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&model.DeliveryAssignment{})

	// Apply filters
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.DriverID != nil {
		query = query.Where("driver_id = ?", *filters.DriverID)
	}

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filters.Page - 1) * filters.PageSize
	if err := query.
		Order("created_at DESC").
		Limit(filters.PageSize).
		Offset(offset).
		Find(&dbModels).Error; err != nil {
		return nil, 0, err
	}

	// Convert to entities
	assignments := make([]*domain.DeliveryAssignment, len(dbModels))
	for i, dbModel := range dbModels {
		assignments[i] = dbModel.ToEntity()
	}

	return assignments, totalCount, nil
}

// GetMetrics retrieves delivery metrics for a time range
func (r *repository) GetMetrics(ctx context.Context, startTime, endTime time.Time, driverID *string) (*domain.DeliveryMetrics, error) {
	var metrics domain.DeliveryMetrics

	query := r.db.WithContext(ctx).Model(&model.DeliveryAssignment{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime)

	if driverID != nil {
		query = query.Where("driver_id = ?", *driverID)
	}

	// Total deliveries
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	metrics.TotalDeliveries = int32(totalCount)
	// Count by status
	type StatusCount struct {
		Status domain.DeliveryStatus
		Count  int32
	}

	var statusCounts []StatusCount
	if err := query.
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case domain.DeliveryStatusDelivered:
			metrics.CompletedDeliveries = sc.Count
		case domain.DeliveryStatusFailed:
			metrics.FailedDeliveries = sc.Count
		case domain.DeliveryStatusCancelled:
			metrics.CancelledDeliveries = sc.Count
		}
	}

	// Average delivery time for completed deliveries
	type AvgTime struct {
		AvgMinutes float64
	}
	var avgTime AvgTime

	if err := r.db.WithContext(ctx).
		Model(&model.DeliveryAssignment{}).
		Where("status = ? AND actual_pickup_time IS NOT NULL AND actual_delivery_time IS NOT NULL", domain.DeliveryStatusDelivered).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Select("AVG(EXTRACT(EPOCH FROM (actual_delivery_time - actual_pickup_time))/60) as avg_minutes").
		Scan(&avgTime).Error; err != nil {
		return nil, err
	}
	metrics.AverageDeliveryTimeMinutes = avgTime.AvgMinutes

	// On-time delivery rate
	type OnTimeCount struct {
		OnTime int32
		Total  int32
	}
	var onTimeCount OnTimeCount

	if err := r.db.WithContext(ctx).
		Model(&model.DeliveryAssignment{}).
		Where("status = ? AND actual_delivery_time IS NOT NULL", domain.DeliveryStatusDelivered).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Select("SUM(CASE WHEN actual_delivery_time <= estimated_delivery_time THEN 1 ELSE 0 END) as on_time, COUNT(*) as total").
		Scan(&onTimeCount).Error; err != nil {
		return nil, err
	}

	if onTimeCount.Total > 0 {
		metrics.OnTimeDeliveryRate = float64(onTimeCount.OnTime) / float64(onTimeCount.Total) * 100
	}

	return &metrics, nil
}

// Delete soft deletes a delivery assignment
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.DeliveryAssignment{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete delivery assignment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// WithTransaction executes a function within a database transaction
func (r *repository) WithTransaction(ctx context.Context, fn func(repo service.DeliveryRepository) error) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create a new repository instance with the transaction
	txRepo := &repository{db: tx}

	// Execute the function
	if err := fn(txRepo); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("failed to rollback transaction after error %v: %w", err, rbErr)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
