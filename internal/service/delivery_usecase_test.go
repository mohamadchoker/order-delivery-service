package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/mohamadchoker/order-delivery-service/internal/domain"
	"github.com/mohamadchoker/order-delivery-service/internal/mocks"
	"github.com/mohamadchoker/order-delivery-service/internal/service"
)

func TestCreateDeliveryAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	now := time.Now()

	input := service.CreateDeliveryInput{
		OrderID:               "ORDER-123",
		PickupAddress:         domain.Address{City: "New York"},
		DeliveryAddress:       domain.Address{City: "Boston"},
		ScheduledPickupTime:   now.Add(1 * time.Hour),
		EstimatedDeliveryTime: now.Add(3 * time.Hour),
		Notes:                 "Test delivery",
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	result, err := uc.CreateDeliveryAssignment(ctx, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ORDER-123", result.OrderID)
	assert.Equal(t, domain.DeliveryStatus("PENDING"), result.Status)
}

func TestCreateDeliveryAssignment_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name  string
		input service.CreateDeliveryInput
	}{
		{
			name: "empty order ID",
			input: service.CreateDeliveryInput{
				OrderID:               "",
				PickupAddress:         domain.Address{City: "New York"},
				DeliveryAddress:       domain.Address{City: "Boston"},
				ScheduledPickupTime:   now.Add(1 * time.Hour),
				EstimatedDeliveryTime: now.Add(3 * time.Hour),
			},
		},
		{
			name: "zero scheduled pickup time",
			input: service.CreateDeliveryInput{
				OrderID:               "ORDER-123",
				PickupAddress:         domain.Address{City: "New York"},
				DeliveryAddress:       domain.Address{City: "Boston"},
				ScheduledPickupTime:   time.Time{},
				EstimatedDeliveryTime: now.Add(3 * time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.CreateDeliveryAssignment(ctx, tt.input)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestGetDeliveryAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()

	expectedAssignment := &domain.DeliveryAssignment{
		ID:      id,
		OrderID: "ORDER-123",
		Status:  domain.DeliveryStatus("PENDING"),
	}

	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(expectedAssignment, nil).
		Times(1)

	result, err := uc.GetDeliveryAssignment(ctx, id)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, "ORDER-123", result.OrderID)
}

func TestGetDeliveryAssignment_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()

	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(nil, domain.ErrNotFound).
		Times(1)

	result, err := uc.GetDeliveryAssignment(ctx, id)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestUpdateDeliveryStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()

	existingAssignment := &domain.DeliveryAssignment{
		ID:      id,
		OrderID: "ORDER-123",
		Status:  domain.DeliveryStatus("PENDING"),
	}

	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(existingAssignment, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	result, err := uc.UpdateDeliveryStatus(ctx, id, domain.DeliveryStatus("ASSIGNED"), "")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, domain.DeliveryStatus("ASSIGNED"), result.Status)
}

func TestUpdateDeliveryStatus_InvalidTransition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()

	existingAssignment := &domain.DeliveryAssignment{
		ID:      id,
		OrderID: "ORDER-123",
		Status:  domain.DeliveryStatus("PENDING"),
	}

	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(existingAssignment, nil).
		Times(1)

	// Should not call Update because validation fails
	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Times(0)

	result, err := uc.UpdateDeliveryStatus(ctx, id, domain.DeliveryStatus("DELIVERED"), "")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAssignDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()
	driverID := "DRIVER-123"

	existingAssignment := &domain.DeliveryAssignment{
		ID:      id,
		OrderID: "ORDER-123",
		Status:  domain.DeliveryStatus("PENDING"),
	}

	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(existingAssignment, nil).
		Times(1)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	result, err := uc.AssignDriver(ctx, id, driverID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.DriverID)
	assert.Equal(t, driverID, *result.DriverID)
	assert.Equal(t, domain.DeliveryStatus("ASSIGNED"), result.Status)
}

func TestAssignDriver_EmptyDriverID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()

	result, err := uc.AssignDriver(ctx, id, "")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListDeliveryAssignments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()

	input := service.ListDeliveryInput{
		Page:     1,
		PageSize: 20,
	}

	expectedAssignments := []*domain.DeliveryAssignment{
		{
			ID:      uuid.New(),
			OrderID: "ORDER-123",
			Status:  domain.DeliveryStatus("PENDING"),
		},
	}

	mockRepo.EXPECT().
		List(ctx, gomock.Any()).
		Return(expectedAssignments, int64(1), nil).
		Times(1)

	result, totalCount, err := uc.ListDeliveryAssignments(ctx, input)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), totalCount)
}

func TestGetDeliveryMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	expectedMetrics := &domain.DeliveryMetrics{
		TotalDeliveries:     10,
		CompletedDeliveries: 8,
		FailedDeliveries:    1,
		CancelledDeliveries: 1,
	}

	mockRepo.EXPECT().
		GetMetrics(ctx, startTime, endTime, nil).
		Return(expectedMetrics, nil).
		Times(1)

	result, err := uc.GetDeliveryMetrics(ctx, startTime, endTime, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int32(10), result.TotalDeliveries)
	assert.Equal(t, int32(8), result.CompletedDeliveries)
}

func TestGetDeliveryMetrics_InvalidTimeRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	startTime := time.Now()
	endTime := time.Now().Add(-24 * time.Hour) // End before start

	result, err := uc.GetDeliveryMetrics(ctx, startTime, endTime, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidInput, err)
}
