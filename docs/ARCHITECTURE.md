# Architecture Documentation

## Overview

The Order Delivery Service is built using Clean Architecture principles, ensuring separation of concerns, testability, and maintainability. The service follows Domain-Driven Design (DDD) patterns and implements industry best practices.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    Delivery Layer                       │
│                   (gRPC Handlers)                       │
│  • Protocol conversion (Proto ↔ Entity)                │
│  • Request validation                                   │
│  • Error handling                                       │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                   Use Case Layer                        │
│                 (Business Logic)                        │
│  • Orchestrates business operations                     │
│  • Enforces business rules                             │
│  • Coordinates repository interactions                  │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                       │
│                  (Data Access)                          │
│  • Database operations                                  │
│  • Data persistence                                     │
│  • Query optimization                                   │
└─────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│                   Domain Layer                          │
│                 (Entities & Rules)                      │
│  • Core business entities                              │
│  • Domain logic                                        │
│  • Business rules                                      │
└─────────────────────────────────────────────────────────┘
```

## Project Structure

```
order-delivery-service/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go          # Server initialization
├── internal/                # Private application code
│   ├── config/             # Configuration management
│   │   └── config.go
│   ├── delivery/           # Delivery layer (gRPC handlers)
│   │   └── grpc_handler.go
│   ├── entity/             # Domain entities
│   │   ├── delivery.go
│   │   ├── delivery_test.go
│   │   └── errors.go
│   ├── repository/         # Repository implementations
│   │   ├── model/          # Database models
│   │   │   └── delivery.go
│   │   ├── postgres_repository.go
│   │   └── repository.go   # Repository interfaces
│   └── usecase/            # Business logic
│       ├── delivery_usecase.go
│       └── delivery_usecase_test.go
├── pkg/                    # Public libraries
│   ├── logger/            # Logging utilities
│   │   └── logger.go
│   └── postgres/          # Database utilities
│       └── postgres.go
├── proto/                 # Protocol buffer definitions
│   └── delivery.proto
├── migrations/            # Database migrations
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── config/               # Configuration files
│   ├── config.yaml
│   └── config.example.yaml
├── docs/                # Documentation
│   ├── API.md
│   └── ARCHITECTURE.md
└── tests/              # Integration tests
    └── integration/
```

## Design Patterns

### 1. Clean Architecture

The application follows Clean Architecture with clear boundaries:

- **Entities**: Core business logic and rules
- **Use Cases**: Application-specific business rules
- **Interface Adapters**: Controllers and presenters (gRPC handlers)
- **Frameworks & Drivers**: External frameworks (GORM, gRPC)

**Dependencies flow inward**: External layers depend on internal layers, never the reverse.

### 2. Repository Pattern

Abstracts data access behind an interface:

```go
type DeliveryRepository interface {
    Create(ctx context.Context, assignment *entity.DeliveryAssignment) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.DeliveryAssignment, error)
    Update(ctx context.Context, assignment *entity.DeliveryAssignment) error
    List(ctx context.Context, filters ListFilters) ([]*entity.DeliveryAssignment, int64, error)
    // ...
}
```

Benefits:
- Easy to swap implementations (e.g., different databases)
- Simplified testing with mocks
- Clear separation of concerns

### 3. Dependency Injection

Dependencies are injected through constructors:

```go
func NewDeliveryUseCase(repo repository.DeliveryRepository, logger *zap.Logger) DeliveryUseCase {
    return &deliveryUseCase{
        repo:   repo,
        logger: logger,
    }
}
```

Benefits:
- Testability
- Flexibility
- Clear dependencies

### 4. Domain-Driven Design

Domain entities encapsulate business logic:

```go
func (d *DeliveryAssignment) UpdateStatus(status DeliveryStatus) error {
    if !d.isValidStatusTransition(status) {
        return ErrInvalidStatusTransition
    }
    // Update logic...
}
```

## Data Flow

### Creating a Delivery Assignment

```
1. gRPC Request arrives at Handler
2. Handler validates request
3. Handler converts Proto → Domain Entity
4. Handler calls UseCase.CreateDeliveryAssignment()
5. UseCase validates business rules
6. UseCase creates Entity with domain logic
7. UseCase calls Repository.Create()
8. Repository converts Entity → DB Model
9. Repository persists to PostgreSQL
10. Response flows back: DB Model → Entity → Proto → gRPC Response
```

### Updating Delivery Status

```
1. gRPC Request with ID and new status
2. UseCase retrieves existing entity
3. Entity validates status transition (domain logic)
4. If valid, entity updates its state
5. Repository persists changes
6. Response with updated entity
```

## Database Design

### delivery_assignments Table

```sql
CREATE TABLE delivery_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id VARCHAR(100) NOT NULL,
    driver_id VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    pickup_address JSONB NOT NULL,
    delivery_address JSONB NOT NULL,
    scheduled_pickup_time TIMESTAMP NOT NULL,
    estimated_delivery_time TIMESTAMP NOT NULL,
    actual_pickup_time TIMESTAMP,
    actual_delivery_time TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

**Indexes:**
- `order_id` - Fast lookups by order
- `driver_id` - Filter deliveries by driver
- `status` - Filter by delivery status
- `scheduled_pickup_time` - Time-based queries
- `created_at` - Chronological ordering
- `deleted_at` - Soft delete support

**JSONB Columns:**
- `pickup_address` and `delivery_address` use JSONB for flexible address storage
- Allows querying nested fields efficiently
- Schema-less for address variations

## Concurrency & Safety

### Context Propagation

All operations accept `context.Context` for:
- Request cancellation
- Timeout management
- Request-scoped values
- Distributed tracing

### Database Transactions

Repository methods can be wrapped in transactions for atomic operations.

### Race Condition Handling

- Database-level constraints prevent invalid states
- Optimistic locking through `updated_at` checks
- Status transitions validated at entity level

## Observability

### Logging

Structured logging with Zap:
```go
logger.Info("Creating delivery assignment",
    zap.String("order_id", input.OrderID),
)
```

**Log Levels:**
- DEBUG: Detailed diagnostic information
- INFO: General operational events
- WARN: Warning conditions
- ERROR: Error conditions

### Metrics

Ready for integration with:
- Prometheus (metrics collection)
- Grafana (visualization)
- OpenTelemetry (distributed tracing)

### Health Checks

gRPC health check protocol implemented:
```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

## Error Handling

### Domain Errors

Defined in entity layer:
```go
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
)
```

### gRPC Status Codes

Mapped appropriately:
- `NOT_FOUND` → Entity not found
- `INVALID_ARGUMENT` → Invalid input
- `FAILED_PRECONDITION` → Invalid state transition
- `INTERNAL` → Unexpected errors

## Testing Strategy

### Unit Tests

Test individual components in isolation:
- Entity business logic
- Use case orchestration
- Handler request/response conversion

**Coverage Target:** >80%

### Integration Tests

Test complete flows with real database:
- End-to-end gRPC calls
- Database operations
- Migration testing

### Mocking

Use `testify/mock` for repository mocks in use case tests.

## Scalability Considerations

### Horizontal Scaling

- Stateless service design
- Multiple instances behind load balancer
- Connection pooling configured

### Database Optimization

- Proper indexing strategy
- Connection pool management
- Query optimization
- JSONB for flexible data

### Caching Strategy

Ready for integration with:
- Redis for frequently accessed data
- Application-level caching
- Query result caching

## Security

### Input Validation

- Request validation at handler level
- Domain validation in entities
- SQL injection prevention (prepared statements)

### Connection Security

- TLS support for gRPC (configurable)
- Database SSL connections (configurable)
- Secrets management via environment variables

## Deployment

### Container Orchestration

Ready for:
- Kubernetes
- Docker Swarm
- ECS/EKS

### Configuration Management

- Environment-based configuration
- Configuration file support
- Environment variable overrides

## Future Enhancements

Potential additions:
- Event sourcing for audit trail
- CQRS for read/write optimization
- Message queue integration (RabbitMQ, Kafka)
- Real-time updates via WebSocket
- GraphQL API layer
- Multi-tenancy support
