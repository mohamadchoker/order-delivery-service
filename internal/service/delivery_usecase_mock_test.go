package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/company/order-delivery-service/internal/domain"
	"github.com/company/order-delivery-service/internal/mocks"
	"github.com/company/order-delivery-service/internal/service"
)

// Example test using uber-go/mock generated mocks
func TestCreateDeliveryAssignment_WithMockgen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock using generated mock
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

	// Set expectations using gomock
	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	// Execute
	result, err := uc.CreateDeliveryAssignment(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ORDER-123", result.OrderID)
	assert.Equal(t, domain.DeliveryStatus("PENDING"), result.Status)
}

// Example test with error case
func TestGetDeliveryAssignment_NotFound_WithMockgen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeliveryRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	uc := service.NewDeliveryUseCase(mockRepo, logger)

	ctx := context.Background()
	id := uuid.New()

	// Set expectation for not found error
	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(nil, domain.ErrNotFound).
		Times(1)

	// Execute
	result, err := uc.GetDeliveryAssignment(ctx, id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrNotFound, err)
}
