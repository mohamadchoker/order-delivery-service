# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

An enterprise-grade microservice for managing order delivery assignments built with Go, gRPC, PostgreSQL, and REST following Clean Architecture principles. The service tracks delivery lifecycles from creation through completion with validated state transitions.

**✨ Recently Refactored (October 2025)**: The codebase has been fully restructured to follow enterprise Go best practices with:
- Proper package naming conventions (`domain`, `service`, `transport`)
- Clean directory structure following Go project layout standards
- Comprehensive validation, custom error types, transaction support
- Request tracing, timeout handling, and Prometheus metrics
- **Dependency Inversion**: Repository interface now owned by service layer
- Security improvements: Credentials externalized, .env support
- Version tracking and build information
- **gRPC-Gateway**: REST/HTTP API alongside gRPC endpoints

## Project Structure

```
order-delivery-service/
├── proto/                      # Protocol Buffer definitions
│   ├── delivery.proto          # Service definitions with HTTP annotations
│   ├── delivery.pb.go          # Generated protobuf code
│   ├── delivery_grpc.pb.go     # Generated gRPC server/client
│   └── delivery.pb.gw.go       # Generated gRPC-Gateway reverse proxy
├── third_party/                # Third-party proto dependencies
│   └── google/api/             # Google API annotations for HTTP
├── cmd/
│   └── server/                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── constants/              # Application-wide constants
│   ├── domain/                 # Domain layer (business entities & rules)
│   ├── repository/             # Repository layer (data access)
│   │   └── postgres/           # PostgreSQL implementation
│   │       └── model/          # Database models
│   ├── service/                # Service layer (business logic/use cases)
│   └── transport/              # Transport layer (protocol handlers)
│       └── grpc/               # gRPC handlers
├── pkg/                        # Public reusable packages
│   ├── logger/                 # Structured logging
│   ├── metrics/                # Prometheus metrics
│   ├── middleware/             # Request ID, timeout middleware
│   ├── postgres/               # Database connection utilities
│   └── validator/              # Input validation
├── migrations/                 # Database migrations
├── config/                     # Configuration files
├── docs/                       # Documentation
├── api.swagger.json            # Generated OpenAPI/Swagger spec
└── scripts/                    # Utility scripts
```

## Development Commands

### Environment Setup

**IMPORTANT**: Ensure `~/go/bin` is in your PATH for protoc generators and tools:
```bash
export PATH=$PATH:~/go/bin
# Add to ~/.zshrc or ~/.bashrc to make permanent:
echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc
```

**IDE Setup**: If using GoLand/IntelliJ and seeing `Cannot resolve import 'google/api/annotations.proto'` errors:
- See [IDE Setup Guide](docs/IDE_SETUP.md) for configuration instructions
- The proto files compile correctly with `make proto` - IDE errors are cosmetic

### Building and Running
```bash
# Generate proto files (required after proto changes)
make proto

# Build the application
make build

# Run the service (default ports: gRPC 50051, HTTP 8080)
make run

# Install required development tools
make install-tools
```

**Default Ports**:
- gRPC: `50051` (configurable via `PORT` env var)
- HTTP/REST: `8080` (configurable via `HTTP_PORT` env var)
- Metrics: `9090` (configurable via `METRICS_PORT` env var)
```

### Testing
```bash
# Run all unit tests with race detection
make test

# Run tests with coverage report (generates coverage.html)
make test-coverage

# Run integration tests (requires PostgreSQL running)
make test-integration

# Run tests for a specific package
go test -v ./internal/domain/
go test -v ./internal/service/
```

### Linting
```bash
# Run linter
make lint

# Auto-fix linting issues
make lint-fix
```

### Database Migrations
```bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_new_field
```

### Docker
```bash
# Build Docker image
make docker-build

# Start all services (PostgreSQL + app)
docker-compose up

# Stop all services
docker-compose down
```

### Testing the APIs

#### gRPC API (port 50051)
```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Health check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# Create delivery assignment
grpcurl -plaintext -d '{
  "order_id": "ORDER-123",
  "pickup_address": {"street": "123 Main St", "city": "NYC", "state": "NY", "postal_code": "10001", "country": "USA"},
  "delivery_address": {"street": "456 Oak Ave", "city": "Boston", "state": "MA", "postal_code": "02101", "country": "USA"},
  "scheduled_pickup_time": "2024-01-15T10:00:00Z",
  "estimated_delivery_time": "2024-01-15T14:00:00Z"
}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment
```

#### REST/HTTP API (port 8080)

The service automatically exposes all gRPC endpoints as REST APIs via grpc-gateway.

**Create delivery assignment:**
```bash
curl -X POST http://localhost:8080/v1/deliveries \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORDER-123",
    "pickup_address": {
      "street": "123 Main St",
      "city": "NYC",
      "state": "NY",
      "postal_code": "10001",
      "country": "USA"
    },
    "delivery_address": {
      "street": "456 Oak Ave",
      "city": "Boston",
      "state": "MA",
      "postal_code": "02101",
      "country": "USA"
    },
    "scheduled_pickup_time": "2024-01-15T10:00:00Z",
    "estimated_delivery_time": "2024-01-15T14:00:00Z"
  }'
```

**Get delivery assignment:**
```bash
curl http://localhost:8080/v1/deliveries/{id}
```

**Update delivery status:**
```bash
curl -X PATCH http://localhost:8080/v1/deliveries/{id}/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PICKED_UP",
    "notes": "Package picked up"
  }'
```

**List delivery assignments:**
```bash
curl "http://localhost:8080/v1/deliveries?page=1&page_size=20&status=PENDING"
```

**Assign driver:**
```bash
curl -X POST http://localhost:8080/v1/deliveries/{id}/assign-driver \
  -H "Content-Type: application/json" \
  -d '{
    "driver_id": "DRIVER-456"
  }'
```

**Get delivery metrics:**
```bash
curl "http://localhost:8080/v1/deliveries/metrics?driver_id=DRIVER-456"
```

**OpenAPI/Swagger Documentation:**
The service generates an OpenAPI specification at `api.swagger.json`. You can view it using any OpenAPI viewer or import it into tools like Postman.
```

## Architecture

### Clean Architecture Layers

The codebase follows strict Clean Architecture with dependencies flowing inward:

1. **Domain Layer** (`internal/domain/`): Core business entities and domain logic
   - `DeliveryAssignment` entity with business rules
   - Status validation and state transitions
   - Domain errors

2. **Service Layer** (`internal/service/`): Application business logic
   - Business logic orchestration
   - Use case implementation
   - Coordinates repository interactions
   - Enforces business rules

3. **Repository Layer** (`internal/repository/postgres/`): Data access abstraction
   - `DeliveryRepository` interface defines contract
   - `PostgresRepository` implements persistence
   - `model/` package contains database models separate from domain entities

4. **Transport Layer** (`internal/transport/grpc/`): Protocol handlers
   - gRPC request handlers
   - Protocol buffer conversion (Proto ↔ Domain)
   - Request validation
   - Error handling and gRPC status mapping

5. **API Layer** (`proto/`): Protocol definitions
   - Protocol buffer definitions with HTTP annotations
   - Generated gRPC code
   - Generated gRPC-Gateway reverse proxy code

### gRPC-Gateway Integration

The service uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to automatically expose REST/HTTP endpoints alongside gRPC:

**How it works**:
1. Proto files include `google.api.http` annotations defining REST mappings
2. `make proto` generates both gRPC and HTTP gateway code
3. At runtime, the gateway acts as a reverse proxy:
   - HTTP requests → Gateway translates to gRPC → gRPC handlers
   - Responses flow back: gRPC → Gateway translates to JSON → HTTP response
4. Same business logic serves both protocols - no code duplication

**Benefits**:
- Single source of truth (proto definitions)
- Automatic OpenAPI/Swagger spec generation
- Type-safe REST API with validation
- Easy integration for clients preferring REST over gRPC
- Both protocols can coexist on different ports

**Example mapping** (`proto/delivery.proto`):
```protobuf
rpc CreateDeliveryAssignment(CreateDeliveryAssignmentRequest) returns (DeliveryAssignment) {
  option (google.api.http) = {
    post: "/v1/deliveries"
    body: "*"
  };
}
```

### Key Architectural Patterns

**Dependency Injection**: All dependencies injected via constructors. The main.go file wires everything together:
```
DB → Repository → Service → Handler → gRPC Server
                                    ↘ HTTP Gateway (reverse proxy to gRPC)
```

**Repository Pattern**: Data access abstracted behind interfaces. Repository works with domain entities, not database models. Internal conversion happens in repository layer.

**Domain-Driven Design**: Business logic lives in domain entities. For example, `DeliveryAssignment.UpdateStatus()` validates state transitions using domain rules.

### Data Flow Example

**Creating a delivery via gRPC**:
1. gRPC request arrives at `Handler.CreateDeliveryAssignment()`
2. Handler converts Proto → Domain Entity
3. Handler calls `Service.CreateDeliveryAssignment()`
4. Service creates entity with `domain.NewDeliveryAssignment()`
5. Service calls `Repository.Create()`
6. Repository converts Entity → DB Model and persists
7. Response flows back: DB Model → Entity → Proto → gRPC Response

**Creating a delivery via REST/HTTP**:
1. HTTP POST to `/v1/deliveries` arrives at Gateway
2. Gateway translates JSON → Proto and forwards to gRPC Handler
3. Steps 2-6 are identical (same gRPC handler and business logic)
4. Response flows back: Proto → Gateway translates to JSON → HTTP Response

## Delivery Status Transitions

Valid status transitions enforced at entity level (`internal/domain/delivery.go:107`):

```
PENDING → ASSIGNED → PICKED_UP → IN_TRANSIT → DELIVERED
   ↓         ↓          ↓            ↓
CANCELLED  CANCELLED  FAILED      FAILED
```

Terminal states (no further transitions): DELIVERED, FAILED, CANCELLED

## Database Schema

The `delivery_assignments` table uses JSONB for flexible address storage. Key indexes on:
- `order_id` - Fast lookups by order
- `driver_id` - Filter deliveries by driver
- `status` - Filter by delivery status
- `scheduled_pickup_time` - Time-based queries

Soft deletes supported via `deleted_at` timestamp.

## Code Organization Principles

**Separation of Concerns**: Database models (`repository/postgres/model/`) are separate from domain entities (`domain/`). This allows domain to remain database-agnostic.

**Interface-Based Design**: Repository and Service are interfaces, making testing and swapping implementations easy.

**⭐ Dependency Inversion Principle** (Key Architectural Improvement):
- Repository interface is defined in the **service layer** (`internal/service/repository.go`)
- Service layer owns the abstraction it depends on
- PostgreSQL implementation (`internal/repository/postgres/`) implements the service interface
- This follows Clean Architecture: inner layers (service) don't depend on outer layers (infrastructure)

**Why This Matters**:
```go
// Service layer defines what it needs
// internal/service/repository.go
type DeliveryRepository interface {
    Create(ctx context.Context, assignment *domain.DeliveryAssignment) error
    // ... other methods
}

// Repository implementation depends on service interface
// internal/repository/postgres/postgres_repository.go
func NewRepository(db *gorm.DB) service.DeliveryRepository {
    return &postgresRepository{db: db}
}
```

Benefits:
- Easy to swap PostgreSQL for MongoDB, DynamoDB, etc. without touching service layer
- Service tests don't need to import repository packages
- True dependency inversion - infrastructure depends on domain, not vice versa

**Context Propagation**: All operations accept `context.Context` for cancellation, timeouts, and request-scoped values.

**Error Handling**: Domain errors defined in `internal/domain/errors.go`. These are mapped to appropriate gRPC status codes in handlers.

## Testing Strategy

**Unit Tests**: Test components in isolation with mocks
- Domain business logic (e.g., status transition validation)
- Service orchestration (mock repository)
- Handler conversion logic

**Integration Tests**: Test with real database
- Use build tag `//go:build integration`
- Require PostgreSQL running (via docker-compose)

**Coverage Target**: >80% for unit tests

**Mock Usage**: Use `testify/mock` for repository mocks. See `internal/service/delivery_usecase_test.go` for examples.

## Configuration

### Configuration Management

Configuration loaded from `config/config.yaml` with environment variable overrides via Viper.

**⚠️ IMPORTANT - Security**:
- `config/config.yaml` is gitignored and should NOT be committed
- Use `config/config.example.yaml` as a template
- Sensitive values (passwords, secrets) should be set via environment variables
- See `.env.example` for required environment variables

**Environment Variable Priority**:
1. Environment variables (highest priority)
2. config/config.yaml file
3. Default values (lowest priority)

**Key Configuration**:
- gRPC server port (default: 50051)
- HTTP gateway port (default: 8080)
- Metrics port (default: 9090)
- Database connection (host, port, user, password, dbname)
- Logger level and mode
- Shutdown timeout
- Connection pool settings

**Environment Variables**:
```bash
# Server
PORT=50051                    # gRPC server port
HTTP_PORT=8080                # HTTP gateway port
METRICS_PORT=9090             # Prometheus metrics port
SHUTDOWN_TIMEOUT=30s          # Graceful shutdown timeout

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=order_delivery_db
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
DB_LOG_SQL=false

# Logger
LOG_LEVEL=info                # debug, info, warn, error
LOG_DEV=false                 # Enable development mode
LOG_STACKTRACE=false          # Enable stack traces
```

**Setup Steps**:
1. Copy `config/config.example.yaml` to `config/config.yaml`
2. Copy `.env.example` to `.env` (optional, for local development)
3. Update values or use environment variables for sensitive data

## Important Implementation Notes

**Proto Generation**: Always run `make proto` after modifying `proto/delivery.proto`. Generated files include:
- `delivery.pb.go` - Protobuf message definitions
- `delivery_grpc.pb.go` - gRPC server and client code
- `delivery.pb.gw.go` - HTTP gateway reverse proxy code
- `api.swagger.json` - OpenAPI/Swagger specification

**Database Connection**: Uses GORM with PostgreSQL driver. Connection pooling configured in `pkg/postgres/`.

**Logging**: Structured logging with Zap. Log levels: DEBUG, INFO, WARN, ERROR. Use appropriate level.

**Graceful Shutdown**: Server handles SIGINT/SIGTERM with configurable timeout. In-flight requests complete before shutdown.

**Health Checks**: Standard gRPC health protocol implemented. Check via `grpc.health.v1.Health/Check`.

## Enterprise Best Practices Implemented

### 1. Comprehensive Input Validation
**Location**: `pkg/validator/`

Dedicated validation package with fluent API for:
- Required field validation
- String length constraints
- Time validation (future dates, ranges, sequencing)
- Address validation with postal code regex
- UUID validation
- Enum validation
- Range validation

**Usage in services**:
```go
v := validator.New()
v.ValidateRequired("order_id", input.OrderID)
v.ValidateAddress("pickup_address", input.PickupAddress)
if v.HasErrors() {
    return nil, v.Errors()
}
```

### 2. Custom Error Types with Context
**Location**: `internal/domain/errors.go`

Rich error types that implement `Unwrap()` and `Is()` for proper error chain handling:
- `DomainError`: Operation-specific errors with codes
- `ValidationError`: Field-level validation errors
- `NotFoundError`: Resource not found with context
- `ConflictError`: State conflict errors

**Error wrapping example**:
```go
if err != nil {
    return fmt.Errorf("failed to create delivery: %w", err)
}
```

### 3. Constants Package
**Location**: `internal/constants/`

Centralized constants for:
- Pagination defaults (page size: 20, max: 100)
- ID length constraints
- Time constraints (min schedule advance: 30 min)
- Database connection pool settings
- Context timeouts (default: 30s, DB query: 10s)
- Request ID header name
- Metrics namespace

**No more magic numbers!**

### 4. Transaction Support
**Location**: `internal/repository/postgres/postgres_repository.go:206`

Repository supports atomic operations:
```go
err := repo.WithTransaction(ctx, func(txRepo postgres.DeliveryRepository) error {
    if err := txRepo.Create(ctx, delivery1); err != nil {
        return err // automatic rollback
    }
    if err := txRepo.Create(ctx, delivery2); err != nil {
        return err // automatic rollback
    }
    return nil // automatic commit
})
```

### 5. Request ID Tracing
**Location**: `pkg/middleware/request_id.go`

Every request gets a unique ID for end-to-end tracing:
- Extracts from `X-Request-ID` header or generates UUID
- Adds to context for use in logging
- Includes in response headers
- Accessible via `middleware.GetRequestID(ctx)`

### 6. Context Timeout Handling
**Location**: `pkg/middleware/timeout.go`

Automatic timeout enforcement (default 30s):
- Prevents runaway requests
- Configurable per-request
- Returns `DeadlineExceeded` error
- Cleans up resources properly

### 7. Prometheus Metrics
**Location**: `pkg/metrics/`

Production-ready observability with:
- **Request metrics**: Total count, duration histogram, active gauge
- **Business metrics**: Delivery assignments by status/operation
- **Database metrics**: Query count, duration by operation
- **HTTP endpoint**: `:9090/metrics`

**Available metrics**:
```
order_delivery_service_grpc_requests_total{method,code}
order_delivery_service_grpc_request_duration_seconds{method}
order_delivery_service_grpc_requests_active{method}
order_delivery_service_delivery_assignments_total{status,operation}
order_delivery_service_database_queries_total{operation,status}
order_delivery_service_database_query_duration_seconds{operation}
```

### 8. Middleware Chain
**Location**: `cmd/server/main.go:69`

Layered middleware for cross-cutting concerns:
```go
grpc.ChainUnaryInterceptor(
    middleware.RequestIDUnaryInterceptor(),     // Request tracing
    middleware.TimeoutUnaryInterceptor(30*time.Second),  // Timeout enforcement
    metrics.MetricsUnaryInterceptor(),          // Prometheus metrics
    loggingInterceptor(log),                    // Structured logging with request ID
)
```

## Module Path

Go module: `github.com/company/order-delivery-service`

When importing internal packages:
```go
import (
    "github.com/company/order-delivery-service/internal/constants"
    "github.com/company/order-delivery-service/internal/domain"
    "github.com/company/order-delivery-service/internal/repository/postgres"
    "github.com/company/order-delivery-service/internal/service"
    "github.com/company/order-delivery-service/internal/transport/grpc"
    "github.com/company/order-delivery-service/pkg/validator"
    "github.com/company/order-delivery-service/pkg/middleware"
    "github.com/company/order-delivery-service/pkg/metrics"
    pb "github.com/company/order-delivery-service/api/grpc"
)
```

## Common Development Patterns

**Adding a New Field to DeliveryAssignment**:
1. Update `domain.DeliveryAssignment` struct
2. Update `repository/postgres/model.DeliveryAssignment` struct
3. Update `proto/delivery.proto` message
4. Run `make proto` to regenerate proto code (updates gRPC and HTTP gateway automatically)
5. Update conversion logic in transport handlers
6. Add validation in `pkg/validator` if needed
7. Create and apply database migration
8. Update tests

**Adding a New gRPC/REST Endpoint**:
1. Define RPC in `proto/delivery.proto` with HTTP annotation:
   ```protobuf
   rpc GetDeliveryStats(GetDeliveryStatsRequest) returns (DeliveryStats) {
     option (google.api.http) = {
       get: "/v1/deliveries/stats"
     };
   }
   ```
2. Run `make proto` (automatically generates both gRPC and REST endpoints)
3. Implement method in `internal/transport/grpc/delivery_handler.go`
4. Add input validation using `pkg/validator`
5. Add corresponding method to `Service` interface and implementation
6. Add repository method if needed
7. Add metrics recording if appropriate
8. Write tests
9. The REST endpoint is automatically available at the specified path

**Using Validation in Services**:
```go
import "github.com/company/order-delivery-service/pkg/validator"

func (s *service) Create(ctx context.Context, input Input) error {
    v := validator.New()
    v.ValidateRequired("order_id", input.OrderID)
    v.ValidateAddress("pickup_address", input.PickupAddress)
    v.ValidateTimeFuture("scheduled_time", input.ScheduledTime)

    if v.HasErrors() {
        return v.Errors() // Returns ValidationErrors type
    }
    // ... proceed with business logic
}
```

**Using Transactions**:
```go
err := repo.WithTransaction(ctx, func(txRepo postgres.DeliveryRepository) error {
    // All operations here are atomic
    if err := txRepo.Create(ctx, entity1); err != nil {
        return fmt.Errorf("failed to create entity1: %w", err)
    }
    if err := txRepo.Update(ctx, entity2); err != nil {
        return fmt.Errorf("failed to update entity2: %w", err)
    }
    return nil // Commit happens here
})
// Rollback happens automatically on any error return
```

**Recording Custom Metrics**:
```go
import "github.com/company/order-delivery-service/pkg/metrics"

// In your service
metrics.RecordDeliveryOperation("create", string(assignment.Status))

// In your repository
start := time.Now()
err := r.db.WithContext(ctx).Create(model).Error
metrics.RecordDatabaseQuery("create_delivery", time.Since(start), err)
```

**Using Constants**:
```go
import "github.com/company/order-delivery-service/internal/constants"

// Pagination
if input.PageSize < constants.MinPageSize {
    input.PageSize = constants.DefaultPageSize
}

// Timeouts
ctx, cancel := context.WithTimeout(ctx, constants.DatabaseQueryTimeout)
defer cancel()

// Error context
return domain.NewDomainError(
    constants.OpCreate,
    "VALIDATION_FAILED",
    "invalid input",
    err,
)
```

**Creating Database Migration**:
```bash
make migrate-create NAME=descriptive_name
# Edit generated files in migrations/
make migrate-up
```

## Monitoring and Observability

**Metrics Endpoint**: `http://localhost:9090/metrics`

**View Prometheus metrics**:
```bash
curl http://localhost:9090/metrics
```

**Key metrics to monitor**:
- Request rate: `rate(order_delivery_service_grpc_requests_total[5m])`
- Error rate: `rate(order_delivery_service_grpc_requests_total{code!="OK"}[5m])`
- Request duration P95: `histogram_quantile(0.95, order_delivery_service_grpc_request_duration_seconds_bucket)`
- Database query duration: `order_delivery_service_database_query_duration_seconds`
- Active requests: `order_delivery_service_grpc_requests_active`

**Request tracing**:
All log lines include `request_id` field for correlation:
```bash
# Find all logs for a specific request
grep "request_id=550e8400-e29b-41d4-a716-446655440000" logs.txt
```
