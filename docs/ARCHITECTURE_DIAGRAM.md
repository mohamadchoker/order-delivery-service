# Architecture Diagram

## Clean Architecture Layers (Current Implementation)

```
┌─────────────────────────────────────────────────────────────┐
│                     Transport Layer                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/transport/grpc/grpc_handler.go            │   │
│  │                                                       │   │
│  │   - CreateDeliveryAssignment()                      │   │
│  │   - GetDeliveryAssignment()                         │   │
│  │   - UpdateDeliveryStatus()                          │   │
│  │   - Proto <-> Domain conversion                     │   │
│  │   - gRPC error mapping                              │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓ depends on ↓
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/service/delivery_usecase.go               │   │
│  │                                                       │   │
│  │   - CreateDeliveryAssignment()                      │   │
│  │   - Business logic orchestration                    │   │
│  │   - Input validation                                │   │
│  │   - Transaction coordination                        │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/service/repository.go ⭐ NEW!             │   │
│  │                                                       │   │
│  │   interface DeliveryRepository {                    │   │
│  │     Create(), GetByID(), Update(),                  │   │
│  │     List(), GetMetrics(), Delete()                  │   │
│  │     WithTransaction()                               │   │
│  │   }                                                  │   │
│  │                                                       │   │
│  │   type ListFilters struct { ... }                   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓ depends on ↓
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/domain/delivery.go                        │   │
│  │                                                       │   │
│  │   type DeliveryAssignment struct {                  │   │
│  │     ID, OrderID, Status, ...                        │   │
│  │   }                                                  │   │
│  │                                                       │   │
│  │   - UpdateStatus() - Business rules                 │   │
│  │   - AssignDriver() - Business rules                 │   │
│  │   - State transition validation                     │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/domain/errors.go                          │   │
│  │                                                       │   │
│  │   - DomainError, ValidationError                    │   │
│  │   - NotFoundError, ConflictError                    │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↑ implements ↑
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                           │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/repository/postgres/postgres_repository.go│   │
│  │                                                       │   │
│  │   implements service.DeliveryRepository ⭐          │   │
│  │                                                       │   │
│  │   - GORM database operations                        │   │
│  │   - Entity <-> DB model conversion                  │   │
│  │   - Transaction management                          │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │   internal/repository/postgres/model/delivery.go     │   │
│  │                                                       │   │
│  │   type DeliveryAssignment struct {                  │   │
│  │     // GORM model with JSON tags                    │   │
│  │   }                                                  │   │
│  │                                                       │   │
│  │   - ToEntity() -> domain.DeliveryAssignment         │   │
│  │   - FromEntity() <- domain.DeliveryAssignment       │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Dependency Flow

```
┌──────────────────────────────────────────────────────────────┐
│                   DEPENDENCY RULE                            │
│                                                              │
│  Dependencies point INWARD (toward domain)                   │
│  Outer layers depend on inner layers, NEVER the reverse      │
│                                                              │
│  Transport  →  Service  →  Domain                            │
│                   ↓                                          │
│             Repository Interface (owned by Service)           │
│                   ↑                                          │
│          Repository Implementation                           │
└──────────────────────────────────────────────────────────────┘
```

## Key Architecture Improvement ⭐

### Before: Incorrect Dependency
```
Service Layer (internal/service/)
    ↓ imports
Repository Package (internal/repository/postgres/)
    ↓ defines
DeliveryRepository interface

❌ Problem: Service depends on Infrastructure!
```

### After: Dependency Inversion ✅
```
Service Layer (internal/service/)
    ↓ defines
DeliveryRepository interface (service/repository.go)
    ↑ implements
Repository Package (internal/repository/postgres/)

✅ Service owns the interface, repository implements it!
```

## Request Flow Example

### Creating a Delivery Assignment

```
1. gRPC Client
      ↓ gRPC Request
┌─────────────────────────────────────┐
│  Transport Layer                    │
│  grpc_handler.CreateDeliveryAssignment()
│    - Validate proto message         │
│    - Convert proto → domain.Address │
│    - Convert proto → service.Input  │
└─────────────────────────────────────┘
      ↓ Call usecase.CreateDeliveryAssignment()
┌─────────────────────────────────────┐
│  Service Layer                      │
│  delivery_usecase.CreateDeliveryAssignment()
│    - Validate business rules        │
│    - Create domain entity           │
│    - Call repo.Create()             │
└─────────────────────────────────────┘
      ↓ Call Create(entity)
┌─────────────────────────────────────┐
│  Domain Layer                       │
│  domain.NewDeliveryAssignment()     │
│    - Apply business rules           │
│    - Validate state                 │
│    - Return entity                  │
└─────────────────────────────────────┘
      ↓ Entity passed to repository
┌─────────────────────────────────────┐
│  Repository Layer                   │
│  postgres_repository.Create()       │
│    - Convert entity → DB model      │
│    - Execute SQL INSERT             │
│    - Convert DB model → entity      │
└─────────────────────────────────────┘
      ↓ Return entity
      ↑ Back through layers
      ↓ Convert to proto
2. gRPC Client ← Response
```

## Middleware Chain

```
gRPC Request
    ↓
┌─────────────────────────────────────┐
│  Request ID Middleware              │
│  - Generate/extract request ID      │
│  - Add to context                   │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│  Timeout Middleware                 │
│  - Enforce 30s timeout              │
│  - Cancel context if exceeded       │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│  Metrics Middleware                 │
│  - Record request count             │
│  - Record duration                  │
│  - Record active requests           │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│  Logging Middleware                 │
│  - Log request start                │
│  - Log request completion           │
│  - Include request ID               │
└─────────────────────────────────────┘
    ↓
Handler (Business Logic)
```

## Data Flow: Domain Entity vs DB Model

```
┌──────────────────────────────────────┐
│  Domain Layer                        │
│                                      │
│  type DeliveryAssignment struct {   │
│    ID        uuid.UUID              │
│    OrderID   string                 │
│    Status    DeliveryStatus         │
│    Address   Address (struct)       │
│  }                                   │
│                                      │
│  ✅ Pure Go types                    │
│  ✅ No database tags                 │
│  ✅ Business logic methods           │
└──────────────────────────────────────┘
            ↕
    FromEntity() / ToEntity()
            ↕
┌──────────────────────────────────────┐
│  Repository Model Layer              │
│                                      │
│  type DeliveryAssignment struct {   │
│    ID        uuid.UUID   `gorm:...` │
│    OrderID   string      `gorm:...` │
│    Status    string      `gorm:...` │
│    Address   []byte      `gorm:...` │
│  }                                   │
│                                      │
│  ✅ GORM tags                         │
│  ✅ JSON serialization                │
│  ✅ Database-specific types           │
└──────────────────────────────────────┘
```

## Package Dependencies

```
cmd/server/
    ↓ imports
internal/app/ (future)
    ↓ imports
internal/transport/grpc/
    ↓ imports
internal/service/
    ↓ imports
internal/domain/
    ←─────────────┐
                  │ imports
internal/repository/postgres/
    ↑             │
    └─────────────┘
    implements service.DeliveryRepository
```

## Cross-Cutting Concerns

```
┌─────────────────────────────────────────────────────────────┐
│                    pkg/ (Reusable Packages)                  │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   logger/    │  │  metrics/    │  │ middleware/  │      │
│  │   Zap        │  │  Prometheus  │  │  Request ID  │      │
│  │   logging    │  │  counters    │  │  Timeout     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  validator/  │  │  postgres/   │  │  constants/  │      │
│  │  Input       │  │  Connection  │  │  App-wide    │      │
│  │  validation  │  │  pool        │  │  constants   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                              │
│  Used by all layers without creating dependencies           │
└─────────────────────────────────────────────────────────────┘
```

## Testing Strategy

```
┌─────────────────────────────────────────────────────────────┐
│  Unit Tests (No external dependencies)                       │
│                                                              │
│  Domain Layer Tests                                          │
│    ✅ Pure business logic                                     │
│    ✅ State transitions                                       │
│    ✅ Validation rules                                        │
│                                                              │
│  Service Layer Tests (with mocks)                            │
│    ✅ Mock DeliveryRepository interface                       │
│    ✅ Test orchestration logic                                │
│    ✅ Test error handling                                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  Integration Tests (with real database)                      │
│                                                              │
│  Repository Layer Tests                                      │
│    ✅ Real PostgreSQL                                         │
│    ✅ Test CRUD operations                                    │
│    ✅ Test transactions                                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  E2E Tests (full stack)                                      │
│                                                              │
│  gRPC Handler Tests                                          │
│    ✅ Real gRPC server                                        │
│    ✅ Real database                                           │
│    ✅ Test entire request flow                                │
└─────────────────────────────────────────────────────────────┘
```

## Benefits of Current Architecture

### ✅ Testability
- Service layer can be tested with mock repository
- Domain layer has zero dependencies
- Each layer independently testable

### ✅ Maintainability
- Clear separation of concerns
- Easy to locate code
- Changes isolated to specific layers

### ✅ Flexibility
- Easy to swap PostgreSQL for another DB
- Easy to add REST/GraphQL alongside gRPC
- Easy to add caching layer

### ✅ SOLID Principles
- **S**ingle Responsibility - Each package has one job
- **O**pen/Closed - Extend via interfaces, not modification
- **L**iskov Substitution - Any repo implementation works
- **I**nterface Segregation - Small, focused interfaces
- **D**ependency Inversion - ⭐ **Service owns repo interface**

---

**Last Updated:** October 23, 2025
