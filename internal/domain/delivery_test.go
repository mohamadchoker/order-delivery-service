package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeliveryAssignment(t *testing.T) {
	orderID := "ORDER-123"
	pickupAddr := Address{City: "New York"}
	deliveryAddr := Address{City: "Boston"}
	scheduledTime := time.Now().Add(1 * time.Hour)
	estimatedTime := time.Now().Add(3 * time.Hour)
	notes := "Handle with care"

	assignment := NewDeliveryAssignment(
		orderID,
		pickupAddr,
		deliveryAddr,
		scheduledTime,
		estimatedTime,
		notes,
	)

	assert.NotNil(t, assignment)
	assert.NotEqual(t, assignment.ID.String(), "")
	assert.Equal(t, orderID, assignment.OrderID)
	assert.Equal(t, DeliveryStatusPending, assignment.Status)
	assert.Nil(t, assignment.DriverID)
	assert.Equal(t, pickupAddr, assignment.PickupAddress)
	assert.Equal(t, deliveryAddr, assignment.DeliveryAddress)
	assert.Equal(t, notes, assignment.Notes)
}

func TestAssignDriver(t *testing.T) {
	tests := []struct {
		name        string
		status      DeliveryStatus
		driverID    string
		expectError bool
	}{
		{
			name:        "successful assignment from pending",
			status:      DeliveryStatusPending,
			driverID:    "DRIVER-123",
			expectError: false,
		},
		{
			name:        "cannot assign from assigned status",
			status:      DeliveryStatusAssigned,
			driverID:    "DRIVER-123",
			expectError: true,
		},
		{
			name:        "cannot assign from delivered status",
			status:      DeliveryStatusDelivered,
			driverID:    "DRIVER-123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignment := &DeliveryAssignment{
				Status: tt.status,
			}

			err := assignment.AssignDriver(tt.driverID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidStatusTransition, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.driverID, *assignment.DriverID)
				assert.Equal(t, DeliveryStatusAssigned, assignment.Status)
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		currentStatus  DeliveryStatus
		newStatus      DeliveryStatus
		expectError    bool
		checkTimestamp bool
	}{
		{
			name:           "pending to assigned",
			currentStatus:  DeliveryStatusPending,
			newStatus:      DeliveryStatusAssigned,
			expectError:    false,
			checkTimestamp: false,
		},
		{
			name:           "assigned to picked up",
			currentStatus:  DeliveryStatusAssigned,
			newStatus:      DeliveryStatusPickedUp,
			expectError:    false,
			checkTimestamp: true,
		},
		{
			name:           "picked up to in transit",
			currentStatus:  DeliveryStatusPickedUp,
			newStatus:      DeliveryStatusInTransit,
			expectError:    false,
			checkTimestamp: false,
		},
		{
			name:           "in transit to delivered",
			currentStatus:  DeliveryStatusInTransit,
			newStatus:      DeliveryStatusDelivered,
			expectError:    false,
			checkTimestamp: true,
		},
		{
			name:          "invalid: pending to delivered",
			currentStatus: DeliveryStatusPending,
			newStatus:     DeliveryStatusDelivered,
			expectError:   true,
		},
		{
			name:          "invalid: delivered to pending",
			currentStatus: DeliveryStatusDelivered,
			newStatus:     DeliveryStatusPending,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignment := &DeliveryAssignment{
				Status: tt.currentStatus,
			}

			err := assignment.UpdateStatus(tt.newStatus)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidStatusTransition, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newStatus, assignment.Status)

				if tt.checkTimestamp {
					if tt.newStatus == DeliveryStatusPickedUp {
						assert.NotNil(t, assignment.ActualPickupTime)
					} else if tt.newStatus == DeliveryStatusDelivered {
						assert.NotNil(t, assignment.ActualDeliveryTime)
					}
				}
			}
		})
	}
}

func TestIsValidStatusTransition(t *testing.T) {
	assignment := &DeliveryAssignment{
		Status: DeliveryStatusPending,
	}

	// Valid transitions from PENDING
	assert.True(t, assignment.isValidStatusTransition(DeliveryStatusAssigned))
	assert.True(t, assignment.isValidStatusTransition(DeliveryStatusCancelled))

	// Invalid transitions from PENDING
	assert.False(t, assignment.isValidStatusTransition(DeliveryStatusPickedUp))
	assert.False(t, assignment.isValidStatusTransition(DeliveryStatusDelivered))

	// Test no transitions from final states
	assignment.Status = DeliveryStatusDelivered
	assert.False(t, assignment.isValidStatusTransition(DeliveryStatusPending))
	assert.False(t, assignment.isValidStatusTransition(DeliveryStatusAssigned))
}
