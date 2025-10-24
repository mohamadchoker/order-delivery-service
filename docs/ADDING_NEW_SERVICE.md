# Adding a New gRPC Service

This guide shows you how to add a new gRPC service to the project following Clean Architecture principles.

## Example: Adding a "Driver Service"

Let's say you want to add a service to manage drivers (create driver, get driver info, list drivers, etc.)

---

## Step-by-Step Guide

### 1. Create Proto Definition

**File**: `proto/driver.proto`

```protobuf
syntax = "proto3";

package driver;

option go_package = "github.com/company/order-delivery-service/proto";

import "google/protobuf/timestamp.proto";

// DriverService manages driver operations
service DriverService {
  // CreateDriver creates a new driver
  rpc CreateDriver(CreateDriverRequest) returns (Driver);

  // GetDriver retrieves a driver by ID
  rpc GetDriver(GetDriverRequest) returns (Driver);

  // ListDrivers lists all drivers with pagination
  rpc ListDrivers(ListDriversRequest) returns (ListDriversResponse);

  // UpdateDriverStatus updates driver availability
  rpc UpdateDriverStatus(UpdateDriverStatusRequest) returns (Driver);
}

// Driver status enum
enum DriverStatus {
  AVAILABLE = 0;
  BUSY = 1;
  OFFLINE = 2;
}

// Driver message
message Driver {
  string id = 1;
  string name = 2;
  string phone = 3;
  string email = 4;
  DriverStatus status = 5;
  string vehicle_type = 6;
  string license_plate = 7;
  google.protobuf.Timestamp created_at = 8;
  google.protobuf.Timestamp updated_at = 9;
}

// Request messages
message CreateDriverRequest {
  string name = 1;
  string phone = 2;
  string email = 3;
  string vehicle_type = 4;
  string license_plate = 5;
}

message GetDriverRequest {
  string id = 1;
}

message ListDriversRequest {
  int32 page = 1;
  int32 page_size = 2;
  DriverStatus status = 3;
}

message ListDriversResponse {
  repeated Driver drivers = 1;
  int32 total_count = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message UpdateDriverStatusRequest {
  string id = 1;
  DriverStatus status = 2;
}
```

**Generate proto code**:
```bash
make proto
```

---

### 2. Create Domain Entity

**File**: `internal/domain/driver.go`

```go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// DriverStatus represents driver availability status
type DriverStatus string

const (
	DriverStatusAvailable DriverStatus = "AVAILABLE"
	DriverStatusBusy      DriverStatus = "BUSY"
	DriverStatusOffline   DriverStatus = "OFFLINE"
)

// Driver represents a driver in the domain
type Driver struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Phone        string       `json:"phone"`
	Email        string       `json:"email"`
	Status       DriverStatus `json:"status"`
	VehicleType  string       `json:"vehicle_type"`
	LicensePlate string       `json:"license_plate"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// NewDriver creates a new driver with default values
func NewDriver(name, phone, email, vehicleType, licensePlate string) *Driver {
	now := time.Now()
	return &Driver{
		ID:           uuid.New(),
		Name:         name,
		Phone:        phone,
		Email:        email,
		Status:       DriverStatusAvailable,
		VehicleType:  vehicleType,
		LicensePlate: licensePlate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UpdateStatus updates the driver status
func (d *Driver) UpdateStatus(status DriverStatus) {
	d.Status = status
	d.UpdatedAt = time.Now()
}
```

---

### 3. Create Repository Interface

**File**: `internal/repository/postgres/driver_repository.go`

```go
package postgres

import (
	"context"

	"github.com/company/order-delivery-service/internal/domain"
	"github.com/google/uuid"
)

// DriverRepository defines the interface for driver data access
type DriverRepository interface {
	Create(ctx context.Context, driver *domain.Driver) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	List(ctx context.Context, page, pageSize int, status *domain.DriverStatus) ([]*domain.Driver, int64, error)
	Update(ctx context.Context, driver *domain.Driver) error
	Delete(ctx context.Context, id uuid.UUID) error
}
```

**Implementation**: `internal/repository/postgres/driver_repository_impl.go`

```go
package postgres

import (
	"context"
	"errors"

	"github.com/company/order-delivery-service/internal/domain"
	"github.com/company/order-delivery-service/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type driverRepository struct {
	db *gorm.DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) Create(ctx context.Context, driver *domain.Driver) error {
	dbModel := model.DriverFromDomain(driver)
	if err := r.db.WithContext(ctx).Create(dbModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *driverRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	var dbModel model.Driver
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return dbModel.ToDomain(), nil
}

func (r *driverRepository) List(ctx context.Context, page, pageSize int, status *domain.DriverStatus) ([]*domain.Driver, int64, error) {
	var drivers []model.Driver
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&model.Driver{})

	if status != nil {
		query = query.Where("status = ?", string(*status))
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&drivers).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain
	result := make([]*domain.Driver, len(drivers))
	for i, d := range drivers {
		result[i] = d.ToDomain()
	}

	return result, totalCount, nil
}

func (r *driverRepository) Update(ctx context.Context, driver *domain.Driver) error {
	dbModel := model.DriverFromDomain(driver)
	if err := r.db.WithContext(ctx).Save(dbModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *driverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&model.Driver{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
```

---

### 4. Create Database Model

**File**: `internal/repository/postgres/model/driver.go`

```go
package model

import (
	"time"

	"github.com/company/order-delivery-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Driver is the database model for drivers
type Driver struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key"`
	Name         string         `gorm:"type:varchar(255);not null"`
	Phone        string         `gorm:"type:varchar(50);not null"`
	Email        string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status       string         `gorm:"type:varchar(50);not null;index"`
	VehicleType  string         `gorm:"type:varchar(100)"`
	LicensePlate string         `gorm:"type:varchar(50)"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (Driver) TableName() string {
	return "drivers"
}

// ToDomain converts database model to domain entity
func (d *Driver) ToDomain() *domain.Driver {
	return &domain.Driver{
		ID:           d.ID,
		Name:         d.Name,
		Phone:        d.Phone,
		Email:        d.Email,
		Status:       domain.DriverStatus(d.Status),
		VehicleType:  d.VehicleType,
		LicensePlate: d.LicensePlate,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

// DriverFromDomain converts domain entity to database model
func DriverFromDomain(d *domain.Driver) *Driver {
	return &Driver{
		ID:           d.ID,
		Name:         d.Name,
		Phone:        d.Phone,
		Email:        d.Email,
		Status:       string(d.Status),
		VehicleType:  d.VehicleType,
		LicensePlate: d.LicensePlate,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
```

---

### 5. Create Service Layer (Use Case)

**File**: `internal/service/driver_usecase.go`

```go
package service

import (
	"context"

	"github.com/company/order-delivery-service/internal/domain"
	"github.com/company/order-delivery-service/internal/repository/postgres"
	"github.com/company/order-delivery-service/pkg/validator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DriverUseCase defines the interface for driver business logic
type DriverUseCase interface {
	CreateDriver(ctx context.Context, input CreateDriverInput) (*domain.Driver, error)
	GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	ListDrivers(ctx context.Context, input ListDriverInput) ([]*domain.Driver, int64, error)
	UpdateDriverStatus(ctx context.Context, id uuid.UUID, status domain.DriverStatus) (*domain.Driver, error)
}

type driverUseCase struct {
	repo   postgres.DriverRepository
	logger *zap.Logger
}

// NewDriverUseCase creates a new driver use case
func NewDriverUseCase(repo postgres.DriverRepository, logger *zap.Logger) DriverUseCase {
	return &driverUseCase{
		repo:   repo,
		logger: logger,
	}
}

// CreateDriverInput contains input for creating a driver
type CreateDriverInput struct {
	Name         string
	Phone        string
	Email        string
	VehicleType  string
	LicensePlate string
}

func (uc *driverUseCase) CreateDriver(ctx context.Context, input CreateDriverInput) (*domain.Driver, error) {
	uc.logger.Info("Creating driver", zap.String("name", input.Name))

	// Validate input
	v := validator.New()
	v.ValidateRequired("name", input.Name)
	v.ValidateRequired("phone", input.Phone)
	v.ValidateRequired("email", input.Email)
	if v.HasErrors() {
		return nil, v.Errors()
	}

	// Create domain entity
	driver := domain.NewDriver(
		input.Name,
		input.Phone,
		input.Email,
		input.VehicleType,
		input.LicensePlate,
	)

	// Persist
	if err := uc.repo.Create(ctx, driver); err != nil {
		uc.logger.Error("Failed to create driver", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Driver created successfully", zap.String("id", driver.ID.String()))
	return driver, nil
}

func (uc *driverUseCase) GetDriver(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	uc.logger.Debug("Getting driver", zap.String("id", id.String()))

	driver, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get driver", zap.Error(err))
		return nil, err
	}

	return driver, nil
}

// ListDriverInput contains input for listing drivers
type ListDriverInput struct {
	Page     int
	PageSize int
	Status   *domain.DriverStatus
}

func (uc *driverUseCase) ListDrivers(ctx context.Context, input ListDriverInput) ([]*domain.Driver, int64, error) {
	uc.logger.Debug("Listing drivers", zap.Int("page", input.Page))

	drivers, total, err := uc.repo.List(ctx, input.Page, input.PageSize, input.Status)
	if err != nil {
		uc.logger.Error("Failed to list drivers", zap.Error(err))
		return nil, 0, err
	}

	return drivers, total, nil
}

func (uc *driverUseCase) UpdateDriverStatus(ctx context.Context, id uuid.UUID, status domain.DriverStatus) (*domain.Driver, error) {
	uc.logger.Info("Updating driver status", zap.String("id", id.String()), zap.String("status", string(status)))

	// Get existing driver
	driver, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status
	driver.UpdateStatus(status)

	// Persist
	if err := uc.repo.Update(ctx, driver); err != nil {
		uc.logger.Error("Failed to update driver status", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Driver status updated successfully", zap.String("id", id.String()))
	return driver, nil
}
```

---

### 6. Create gRPC Transport Layer

**File**: `internal/transport/grpc/driver_handler.go`

```go
package grpc

import (
	"context"

	"github.com/company/order-delivery-service/internal/service"
	pb "github.com/company/order-delivery-service/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DriverHandler implements the gRPC DriverService
type DriverHandler struct {
	pb.UnimplementedDriverServiceServer
	useCase service.DriverUseCase
	logger  *zap.Logger
}

// NewDriverHandler creates a new gRPC driver handler
func NewDriverHandler(useCase service.DriverUseCase, logger *zap.Logger) *DriverHandler {
	return &DriverHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (h *DriverHandler) CreateDriver(ctx context.Context, req *pb.CreateDriverRequest) (*pb.Driver, error) {
	h.logger.Info("Received CreateDriver request", zap.String("name", req.Name))

	// Validate
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// Convert to service input
	input := service.CreateDriverInput{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		VehicleType:  req.VehicleType,
		LicensePlate: req.LicensePlate,
	}

	// Call use case
	driver, err := h.useCase.CreateDriver(ctx, input)
	if err != nil {
		return nil, handleError(err)
	}

	return driverToProto(driver), nil
}

func (h *DriverHandler) GetDriver(ctx context.Context, req *pb.GetDriverRequest) (*pb.Driver, error) {
	h.logger.Debug("Received GetDriver request", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	driver, err := h.useCase.GetDriver(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return driverToProto(driver), nil
}

func (h *DriverHandler) ListDrivers(ctx context.Context, req *pb.ListDriversRequest) (*pb.ListDriversResponse, error) {
	h.logger.Debug("Received ListDrivers request", zap.Int32("page", req.Page))

	input := service.ListDriverInput{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}

	if req.Status != pb.DriverStatus_AVAILABLE { // Assuming AVAILABLE is default
		domainStatus := protoStatusToDriverDomain(req.Status)
		input.Status = &domainStatus
	}

	drivers, total, err := h.useCase.ListDrivers(ctx, input)
	if err != nil {
		return nil, handleError(err)
	}

	protoDrivers := make([]*pb.Driver, len(drivers))
	for i, driver := range drivers {
		protoDrivers[i] = driverToProto(driver)
	}

	return &pb.ListDriversResponse{
		Drivers:    protoDrivers,
		TotalCount: int32(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (h *DriverHandler) UpdateDriverStatus(ctx context.Context, req *pb.UpdateDriverStatusRequest) (*pb.Driver, error) {
	h.logger.Info("Received UpdateDriverStatus request", zap.String("id", req.Id))

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	domainStatus := protoStatusToDriverDomain(req.Status)

	driver, err := h.useCase.UpdateDriverStatus(ctx, id, domainStatus)
	if err != nil {
		return nil, handleError(err)
	}

	return driverToProto(driver), nil
}
```

**File**: `internal/transport/grpc/driver_converter.go`

```go
package grpc

import (
	"github.com/company/order-delivery-service/internal/domain"
	pb "github.com/company/order-delivery-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Driver conversions

func protoStatusToDriverDomain(s pb.DriverStatus) domain.DriverStatus {
	switch s {
	case pb.DriverStatus_AVAILABLE:
		return domain.DriverStatusAvailable
	case pb.DriverStatus_BUSY:
		return domain.DriverStatusBusy
	case pb.DriverStatus_OFFLINE:
		return domain.DriverStatusOffline
	default:
		return domain.DriverStatusAvailable
	}
}

func driverStatusToProto(s domain.DriverStatus) pb.DriverStatus {
	switch s {
	case domain.DriverStatusAvailable:
		return pb.DriverStatus_AVAILABLE
	case domain.DriverStatusBusy:
		return pb.DriverStatus_BUSY
	case domain.DriverStatusOffline:
		return pb.DriverStatus_OFFLINE
	default:
		return pb.DriverStatus_AVAILABLE
	}
}

func driverToProto(d *domain.Driver) *pb.Driver {
	return &pb.Driver{
		Id:           d.ID.String(),
		Name:         d.Name,
		Phone:        d.Phone,
		Email:        d.Email,
		Status:       driverStatusToProto(d.Status),
		VehicleType:  d.VehicleType,
		LicensePlate: d.LicensePlate,
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
	}
}
```

---

### 7. Create Database Migration

```bash
make migrate-create NAME=create_drivers_table
```

**File**: `migrations/000002_create_drivers_table.up.sql`

```sql
CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    vehicle_type VARCHAR(100),
    license_plate VARCHAR(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_drivers_status ON drivers(status);
CREATE INDEX idx_drivers_email ON drivers(email);
CREATE INDEX idx_drivers_deleted_at ON drivers(deleted_at);
```

**File**: `migrations/000002_create_drivers_table.down.sql`

```sql
DROP TABLE IF EXISTS drivers;
```

---

### 8. Wire Everything in main.go

**File**: `cmd/server/main.go`

```go
func main() {
	// ... existing code ...

	// Initialize repositories
	deliveryRepo := postgres.NewDeliveryRepository(db)
	driverRepo := postgres.NewDriverRepository(db)  // NEW

	// Initialize use cases
	deliveryUseCase := service.NewDeliveryUseCase(deliveryRepo, logger)
	driverUseCase := service.NewDriverUseCase(driverRepo, logger)  // NEW

	// Initialize handlers
	deliveryHandler := grpc.NewHandler(deliveryUseCase, logger)
	driverHandler := grpc.NewDriverHandler(driverUseCase, logger)  // NEW

	// Register services
	pb.RegisterDeliveryServiceServer(grpcServer, deliveryHandler)
	pb.RegisterDriverServiceServer(grpcServer, driverHandler)  // NEW

	// ... rest of code ...
}
```

---

## Naming Conventions

### Proto Files
- **File name**: `{service_name}.proto` (lowercase, snake_case)
  - ✅ `driver.proto`
  - ✅ `order.proto`
  - ❌ `driverService.proto`

### Package Names
- **Proto package**: Lowercase, singular
  - ✅ `package driver;`
  - ✅ `package order;`
  - ❌ `package drivers;`

### Service Names
- **gRPC service**: PascalCase, ends with "Service"
  - ✅ `service DriverService`
  - ✅ `service OrderService`
  - ❌ `service Driver`

### Messages
- **Message names**: PascalCase, descriptive
  - ✅ `message Driver`
  - ✅ `message CreateDriverRequest`
  - ✅ `message ListDriversResponse`

### Enums
- **Enum names**: PascalCase
  - ✅ `enum DriverStatus`
  - ✅ `enum OrderType`

- **Enum values**: UPPERCASE_SNAKE_CASE
  - ✅ `AVAILABLE`
  - ✅ `IN_TRANSIT`
  - ❌ `Available`

### Go Files

#### Domain Layer
- **File**: `internal/domain/{entity}.go` (singular)
  - ✅ `driver.go`
  - ✅ `order.go`
  - ❌ `drivers.go`

#### Repository Layer
- **Interface file**: `internal/repository/postgres/{entity}_repository.go`
  - ✅ `driver_repository.go`

- **Implementation file**: `internal/repository/postgres/{entity}_repository_impl.go`
  - ✅ `driver_repository_impl.go`

- **Model file**: `internal/repository/postgres/model/{entity}.go`
  - ✅ `model/driver.go`

#### Service Layer
- **File**: `internal/service/{entity}_usecase.go`
  - ✅ `driver_usecase.go`
  - ✅ `order_usecase.go`

#### Transport Layer
- **Handler file**: `internal/transport/grpc/{entity}_handler.go`
  - ✅ `driver_handler.go`

- **Converter file**: `internal/transport/grpc/{entity}_converter.go`
  - ✅ `driver_converter.go`

### Interface Names
- **Repository**: `{Entity}Repository`
  - ✅ `DriverRepository`
  - ✅ `OrderRepository`

- **Use Case**: `{Entity}UseCase`
  - ✅ `DriverUseCase`
  - ✅ `OrderUseCase`

### Struct Names
- **Handler**: `{Entity}Handler` (private in package)
  - ✅ `type DriverHandler struct`

- **Use Case impl**: `{entity}UseCase` (private)
  - ✅ `type driverUseCase struct`

- **Repository impl**: `{entity}Repository` (private)
  - ✅ `type driverRepository struct`

---

## Quick Checklist

When adding a new service:

- [ ] Create proto file (`proto/{service}.proto`)
- [ ] Run `make proto` to generate code
- [ ] Create domain entity (`internal/domain/{entity}.go`)
- [ ] Create repository interface (`internal/repository/postgres/{entity}_repository.go`)
- [ ] Create repository implementation (`internal/repository/postgres/{entity}_repository_impl.go`)
- [ ] Create database model (`internal/repository/postgres/model/{entity}.go`)
- [ ] Create use case interface and implementation (`internal/service/{entity}_usecase.go`)
- [ ] Create gRPC handler (`internal/transport/grpc/{entity}_handler.go`)
- [ ] Create converter (`internal/transport/grpc/{entity}_converter.go`)
- [ ] Create database migration (`migrations/...`)
- [ ] Wire everything in `cmd/server/main.go`
- [ ] Write tests for each layer
- [ ] Update documentation

---

## File Structure Summary

```
order-delivery-service/
├── proto/
│   ├── delivery.proto          ← Existing service
│   └── driver.proto            ← New service
├── internal/
│   ├── domain/
│   │   ├── delivery.go
│   │   └── driver.go           ← New domain entity
│   ├── repository/postgres/
│   │   ├── delivery_repository.go
│   │   ├── driver_repository.go      ← New repo interface
│   │   ├── driver_repository_impl.go ← New repo implementation
│   │   └── model/
│   │       ├── delivery.go
│   │       └── driver.go       ← New DB model
│   ├── service/
│   │   ├── delivery_usecase.go
│   │   └── driver_usecase.go   ← New use case
│   └── transport/grpc/
│       ├── grpc_handler.go     ← Delivery handler
│       ├── converter.go        ← Delivery converter
│       ├── driver_handler.go   ← New driver handler
│       └── driver_converter.go ← New driver converter
├── migrations/
│   ├── 000001_create_delivery_assignments_table.up.sql
│   └── 000002_create_drivers_table.up.sql ← New migration
└── cmd/server/
    └── main.go                 ← Wire everything here
```

---

## Testing Your New Service

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Should show:
# delivery.DeliveryService
# driver.DriverService
# grpc.health.v1.Health

# Test create driver
grpcurl -plaintext -d '{
  "name": "John Doe",
  "phone": "+1234567890",
  "email": "john@example.com",
  "vehicle_type": "Van",
  "license_plate": "ABC-123"
}' localhost:50051 driver.DriverService/CreateDriver

# Test list drivers
grpcurl -plaintext localhost:50051 driver.DriverService/ListDrivers
```

---

## Best Practices

1. **Keep layers separate**: Domain → Service → Repository → Transport
2. **One service per proto file**: Don't mix multiple services in one proto
3. **Use converters**: Keep conversion logic in separate converter files
4. **Consistent naming**: Follow the conventions above
5. **Test each layer**: Write unit tests for domain, service, and integration tests
6. **Migrations**: Always create up and down migrations
7. **Validation**: Validate in both handler and service layers
8. **Errors**: Use domain errors, map to gRPC codes in converter

---

This structure ensures your codebase stays clean, maintainable, and scalable!
