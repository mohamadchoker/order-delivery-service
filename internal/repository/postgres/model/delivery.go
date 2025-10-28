package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mohamadchoker/order-delivery-service/internal/domain"
)

// Address is a custom type for storing address as JSONB in PostgreSQL
type Address domain.Address

// Scan implements the sql.Scanner interface for Address
func (a *Address) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// Value implements the driver.Valuer interface for Address
func (a *Address) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// DeliveryAssignment is the GORM model for delivery_assignments table
type DeliveryAssignment struct {
	ID                    uuid.UUID             `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID               string                `gorm:"type:varchar(100);not null;index"`
	DriverID              *string               `gorm:"type:varchar(100);index"`
	Status                domain.DeliveryStatus `gorm:"type:varchar(50);not null;index"`
	PickupAddress         Address               `gorm:"type:jsonb;not null"`
	DeliveryAddress       Address               `gorm:"type:jsonb;not null"`
	ScheduledPickupTime   time.Time             `gorm:"not null;index"`
	EstimatedDeliveryTime time.Time             `gorm:"not null"`
	ActualPickupTime      *time.Time
	ActualDeliveryTime    *time.Time
	Notes                 string         `gorm:"type:text"`
	CreatedAt             time.Time      `gorm:"not null;index"`
	UpdatedAt             time.Time      `gorm:"not null"`
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for DeliveryAssignment
func (DeliveryAssignment) TableName() string {
	return "delivery_assignments"
}

// ToEntity converts the GORM model to domain entity
func (d *DeliveryAssignment) ToEntity() *domain.DeliveryAssignment {
	return &domain.DeliveryAssignment{
		ID:                    d.ID,
		OrderID:               d.OrderID,
		DriverID:              d.DriverID,
		Status:                d.Status,
		PickupAddress:         domain.Address(d.PickupAddress),
		DeliveryAddress:       domain.Address(d.DeliveryAddress),
		ScheduledPickupTime:   d.ScheduledPickupTime,
		EstimatedDeliveryTime: d.EstimatedDeliveryTime,
		ActualPickupTime:      d.ActualPickupTime,
		ActualDeliveryTime:    d.ActualDeliveryTime,
		Notes:                 d.Notes,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
	}
}

// FromEntity converts domain entity to GORM model
func FromEntity(e *domain.DeliveryAssignment) *DeliveryAssignment {
	return &DeliveryAssignment{
		ID:                    e.ID,
		OrderID:               e.OrderID,
		DriverID:              e.DriverID,
		Status:                e.Status,
		PickupAddress:         Address(e.PickupAddress),
		DeliveryAddress:       Address(e.DeliveryAddress),
		ScheduledPickupTime:   e.ScheduledPickupTime,
		EstimatedDeliveryTime: e.EstimatedDeliveryTime,
		ActualPickupTime:      e.ActualPickupTime,
		ActualDeliveryTime:    e.ActualDeliveryTime,
		Notes:                 e.Notes,
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
	}
}
