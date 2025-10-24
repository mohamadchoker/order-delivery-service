package domain

import (
	"time"

	"github.com/google/uuid"
)

// DeliveryStatus represents the current status of a delivery
type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "PENDING"
	DeliveryStatusAssigned  DeliveryStatus = "ASSIGNED"
	DeliveryStatusPickedUp  DeliveryStatus = "PICKED_UP"
	DeliveryStatusInTransit DeliveryStatus = "IN_TRANSIT"
	DeliveryStatusDelivered DeliveryStatus = "DELIVERED"
	DeliveryStatusFailed    DeliveryStatus = "FAILED"
	DeliveryStatusCancelled DeliveryStatus = "CANCELED"
)

// Address represents a physical address with coordinates
type Address struct {
	Street     string  `json:"street"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

// DeliveryAssignment represents a delivery assignment in the domain
type DeliveryAssignment struct {
	ID                    uuid.UUID      `json:"id"`
	OrderID               string         `json:"order_id"`
	DriverID              *string        `json:"driver_id,omitempty"`
	Status                DeliveryStatus `json:"status"`
	PickupAddress         Address        `json:"pickup_address"`
	DeliveryAddress       Address        `json:"delivery_address"`
	ScheduledPickupTime   time.Time      `json:"scheduled_pickup_time"`
	EstimatedDeliveryTime time.Time      `json:"estimated_delivery_time"`
	ActualPickupTime      *time.Time     `json:"actual_pickup_time,omitempty"`
	ActualDeliveryTime    *time.Time     `json:"actual_delivery_time,omitempty"`
	Notes                 string         `json:"notes"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

// NewDeliveryAssignment creates a new delivery assignment with default values
func NewDeliveryAssignment(
	orderID string,
	pickupAddress Address,
	deliveryAddress Address,
	scheduledPickupTime time.Time,
	estimatedDeliveryTime time.Time,
	notes string,
) *DeliveryAssignment {
	now := time.Now()
	return &DeliveryAssignment{
		ID:                    uuid.New(),
		OrderID:               orderID,
		Status:                DeliveryStatusPending,
		PickupAddress:         pickupAddress,
		DeliveryAddress:       deliveryAddress,
		ScheduledPickupTime:   scheduledPickupTime,
		EstimatedDeliveryTime: estimatedDeliveryTime,
		Notes:                 notes,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// AssignDriver assigns a driver to the delivery
func (d *DeliveryAssignment) AssignDriver(driverID string) error {
	if d.Status != DeliveryStatusPending {
		return ErrInvalidStatusTransition
	}
	d.DriverID = &driverID
	d.Status = DeliveryStatusAssigned
	d.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus updates the delivery status with validation
func (d *DeliveryAssignment) UpdateStatus(status DeliveryStatus) error {
	if !d.isValidStatusTransition(status) {
		return ErrInvalidStatusTransition
	}

	d.Status = status
	d.UpdatedAt = time.Now()

	// Set timestamps based on status
	now := time.Now()
	switch status {
	case DeliveryStatusPickedUp:
		d.ActualPickupTime = &now
	case DeliveryStatusDelivered:
		d.ActualDeliveryTime = &now
	}

	return nil
}

// isValidStatusTransition checks if a status transition is valid
func (d *DeliveryAssignment) isValidStatusTransition(newStatus DeliveryStatus) bool {
	validTransitions := map[DeliveryStatus][]DeliveryStatus{
		DeliveryStatusPending:   {DeliveryStatusAssigned, DeliveryStatusCancelled},
		DeliveryStatusAssigned:  {DeliveryStatusPickedUp, DeliveryStatusCancelled},
		DeliveryStatusPickedUp:  {DeliveryStatusInTransit, DeliveryStatusFailed},
		DeliveryStatusInTransit: {DeliveryStatusDelivered, DeliveryStatusFailed},
		DeliveryStatusDelivered: {},
		DeliveryStatusFailed:    {},
		DeliveryStatusCancelled: {},
	}

	allowed, exists := validTransitions[d.Status]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}

	return false
}

// DeliveryMetrics contains aggregated delivery statistics
type DeliveryMetrics struct {
	TotalDeliveries            int32   `json:"total_deliveries"`
	CompletedDeliveries        int32   `json:"completed_deliveries"`
	FailedDeliveries           int32   `json:"failed_deliveries"`
	CancelledDeliveries        int32   `json:"canceled_deliveries"`
	AverageDeliveryTimeMinutes float64 `json:"average_delivery_time_minutes"`
	OnTimeDeliveryRate         float64 `json:"on_time_delivery_rate"`
}
