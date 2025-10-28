# Codebase Structure Improvements

## Overview
This document outlines the structural improvements made to the order-delivery-service codebase to follow enterprise Go best practices and Clean Architecture principles.

---

## ✅ Completed Improvements

### 1. Repository Interface Dependency Inversion ⭐ **CRITICAL**

**Problem:** Service layer was tightly coupled to the concrete postgres implementation

**Before:**
```go
// internal/service/delivery_usecase.go
import "github.com/mohamadchoker/order-delivery-service/internal/repository/postgres"

type deliveryUseCase struct {
    repo postgres.DeliveryRepository  // ❌ Depends on implementation
}
```

**After:**
```go
// internal/service/repository.go - NEW FILE
package service

type DeliveryRepository interface {
    Create(ctx context.Context, assignment *domain.DeliveryAssignment) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.DeliveryAssignment, error)
    // ... other methods
}

type ListFilters struct {
    Page     int
    PageSize int
    Status   *domain.DeliveryStatus
    DriverID *string
}

// internal/service/delivery_usecase.go
type deliveryUseCase struct {
    repo DeliveryRepository  // ✅ Depends on interface
}
```

**Benefits:**
- ✅ Service layer now owns its dependency interface (Dependency Inversion Principle)
- ✅ Repository implementation depends on service interface, not vice versa
- ✅ Easy to swap PostgreSQL for MongoDB, DynamoDB, etc.
- ✅ Mock testing simplified (no need to import postgres package)
- ✅ True Clean Architecture - domain and service layers have no infrastructure dependencies

**Files Changed:**
- ✅ Created: `internal/service/repository.go`
- ✅ Updated: `internal/service/delivery_usecase.go`
- ✅ Updated: `internal/repository/postgres/postgres_repository.go`
- ✅ Updated: `internal/service/delivery_usecase_test.go`
- ✅ Deleted: `internal/repository/postgres/repository.go` (interface moved to service)

---

## 📋 Recommended Structural Improvements (Not Yet Implemented)

### 2. Configuration Package Structure

**Current:**
```
internal/config/
└── config.go  (330 lines - too large)
```

**Recommended:**
```
internal/config/
├── config.go           # Main config struct and Load()
├── server.go           # ServerConfig type and validation
├── database.go         # DatabaseConfig type and validation
├── logger.go           # LoggerConfig type and validation
└── validator.go        # Config validation logic
```

**Benefits:**
- Separation of concerns
- Easier to test each config component
- Cleaner file structure

---

### 3. Domain Package Organization

**Current:**
```
internal/domain/
├── delivery.go          # Entity, constructor, business logic
├── delivery_test.go     # Tests
└── errors.go            # All errors
```

**Recommended:**
```
internal/domain/
├── delivery/
│   ├── delivery.go           # DeliveryAssignment entity
│   ├── delivery_test.go      # Entity tests
│   ├── repository.go         # Repository interface (optional, alternative to service layer)
│   └── events.go             # Domain events (if using DDD events)
├── address/
│   ├── address.go            # Address value object
│   └── address_test.go
├── status/
│   ├── status.go             # DeliveryStatus enum and transitions
│   └── status_test.go
└── errors/
    ├── errors.go             # Error types
    └── codes.go              # Error codes constants
```

**Benefits:**
- Clear separation by aggregate root
- Value objects are first-class citizens
- Status transition logic isolated
- Easier to add new aggregates (e.g., Driver, Vehicle)

---

### 4. Service Layer Structure (For Multiple Aggregates)

**Current:**
```
internal/service/
├── delivery_usecase.go
├── delivery_usecase_test.go
└── repository.go
```

**Recommended (When Growing):**
```
internal/service/
├── delivery/
│   ├── service.go            # Delivery service implementation
│   ├── service_test.go
│   ├── dto.go                # Input/Output DTOs
│   └── repository.go         # Repository interface
├── driver/                   # Future: driver management
│   └── ...
├── notification/             # Future: notification service
│   └── ...
└── orchestration/            # Complex workflows spanning multiple services
    └── assignment_workflow.go
```

**Benefits:**
- Clear boundaries between business capabilities
- Each service is independently testable
- Easier to extract to microservices later
- Reduces god-object anti-pattern

---

### 5. Create Dedicated DTO Package

**Current:** DTOs mixed in service files
```go
// internal/service/delivery_usecase.go
type CreateDeliveryInput struct { ... }
type ListDeliveryInput struct { ... }
```

**Recommended:**
```
internal/service/
└── dto/
    ├── delivery_dto.go      # All delivery-related DTOs
    ├── common_dto.go        # Shared DTOs (Pagination, etc.)
    └── mapper.go            # Domain <-> DTO conversion
```

**Benefits:**
- Clear contract definition
- Reusable across transport layers (gRPC, REST, GraphQL)
- Easy to version (v1/v2 DTOs)
- Separation from business logic

---

### 6. Add Validation Layer Abstraction

**Current:** Validator is a concrete package
```
pkg/validator/validator.go
```

**Recommended:**
```
pkg/validation/
├── validator.go         # Validator interface
├── fluent/
│   └── fluent_validator.go  # Current implementation
├── struct/
│   └── struct_validator.go  # go-playground/validator wrapper
└── chain/
    └── chain_validator.go   # Chain multiple validators
```

**Benefits:**
- Pluggable validation strategies
- Can combine fluent + struct tag validation
- Easy to add custom validators

---

### 7. Middleware Organization

**Current:**
```
pkg/middleware/
├── request_id.go
└── timeout.go
```

**Recommended:**
```
pkg/middleware/
├── grpc/
│   ├── request_id.go
│   ├── timeout.go
│   ├── auth.go          # Future: authentication
│   ├── ratelimit.go     # Future: rate limiting
│   └── recovery.go      # Future: panic recovery
├── http/                # If adding REST API
│   └── ...
└── middleware.go        # Common middleware interface
```

**Benefits:**
- Protocol-specific middleware isolated
- Easy to add HTTP/GraphQL support
- Clear organization as middleware grows

---

### 8. Observability Package Structure

**Current:**
```
pkg/
├── logger/
│   └── logger.go
└── metrics/
    └── metrics.go
```

**Recommended:**
```
pkg/observability/
├── logging/
│   ├── logger.go              # Logger interface
│   ├── zap/
│   │   └── zap_logger.go      # Zap implementation
│   └── context.go             # Context helpers
├── metrics/
│   ├── metrics.go             # Metrics interface
│   ├── prometheus/
│   │   └── prometheus_metrics.go
│   └── collector.go           # Custom collectors
└── tracing/
    ├── tracer.go              # Tracer interface
    └── otel/
        └── otel_tracer.go     # OpenTelemetry implementation
```

**Benefits:**
- All observability concerns in one place
- Easy to swap implementations (e.g., Zap -> Zerolog)
- Centralized tracing support
- Consistent interface across observability tools

---

### 9. Add Infrastructure Package

**Current:** Infrastructure scattered across pkg/
```
pkg/
├── postgres/
└── ...
```

**Recommended:**
```
internal/infrastructure/
├── database/
│   ├── postgres/
│   │   ├── connection.go
│   │   ├── health.go
│   │   └── migrator.go
│   └── transaction.go         # Transaction manager interface
├── cache/
│   ├── redis/
│   │   └── client.go
│   └── memory/
│       └── cache.go
├── messaging/
│   ├── kafka/
│   └── rabbitmq/
└── storage/
    ├── s3/
    └── gcs/
```

**Benefits:**
- Clear separation of infrastructure concerns
- Easy to mock in tests
- Supports multiple implementations
- Infrastructure decisions isolated from business logic

---

### 10. Create Application Layer

**Current:** Main file handles everything
```
cmd/server/main.go (207 lines)
```

**Recommended:**
```
internal/app/
├── app.go                    # Application struct
├── dependencies.go           # Dependency injection
├── server.go                 # Server lifecycle
└── config.go                 # Config loading

cmd/server/
└── main.go                   # Just bootstrapping (< 50 lines)
```

**Example:**
```go
// internal/app/app.go
package app

type Application struct {
    config  *config.Config
    logger  *zap.Logger
    db      *gorm.DB
    server  *grpc.Server
}

func New(configPath string) (*Application, error) {
    // Load config
    // Setup logger
    // Connect DB
    // Wire dependencies
}

func (a *Application) Start() error {
    // Start servers
    // Start background jobs
}

func (a *Application) Stop() error {
    // Graceful shutdown
}

// cmd/server/main.go
func main() {
    app, err := app.New("config/config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

**Benefits:**
- Main file is minimal and clean
- Application lifecycle clearly defined
- Easy to add CLI commands (migrate, seed, etc.)
- Testable application initialization

---

### 11. Add API Versioning Structure

**Current:**
```
api/grpc/
├── delivery.proto
├── delivery.pb.go
└── delivery_grpc.pb.go
```

**Recommended:**
```
api/
├── v1/
│   ├── delivery.proto
│   ├── driver.proto        # Future
│   └── common.proto        # Shared messages
├── v2/                     # Future version
│   └── ...
└── docs/
    ├── api.html            # Generated docs
    └── openapi.yaml        # If adding REST
```

**Benefits:**
- Clear API versioning strategy
- Backward compatibility support
- Easy to deprecate old versions
- Generated documentation organized

---

### 12. Tests Organization

**Current:** Tests alongside code
```
internal/domain/delivery_test.go
internal/service/delivery_usecase_test.go
```

**Recommended:** Add integration tests structure
```
tests/
├── integration/
│   ├── delivery_integration_test.go
│   ├── repository_integration_test.go
│   └── fixtures/
│       ├── seed.sql
│       └── test_data.go
├── e2e/
│   ├── grpc_e2e_test.go
│   └── scenarios/
│       ├── happy_path_test.go
│       └── error_cases_test.go
├── performance/
│   └── load_test.go
└── testutil/
    ├── database.go          # Test DB helpers
    ├── fixtures.go          # Test data builders
    └── assertions.go        # Custom assertions
```

**Benefits:**
- Clear separation between unit, integration, E2E tests
- Shared test utilities
- Easy to run specific test suites
- Performance testing infrastructure

---

### 13. Add Scripts Organization

**Current:** Single Makefile
```
Makefile (single file with all commands)
```

**Recommended:**
```
scripts/
├── build/
│   ├── build.sh
│   └── docker-build.sh
├── db/
│   ├── migrate.sh
│   ├── seed.sh
│   └── backup.sh
├── dev/
│   ├── setup.sh            # Development environment setup
│   └── start-deps.sh       # Start postgres, redis, etc.
├── ci/
│   ├── lint.sh
│   ├── test.sh
│   └── coverage.sh
└── deploy/
    ├── k8s-deploy.sh
    └── rollback.sh

Makefile                    # Delegates to scripts/
```

**Benefits:**
- Scripts are reusable outside Make
- CI/CD can use scripts directly
- Easier to test scripts
- Better organization of automation

---

### 14. Add Deployment Configurations

**Current:** Only docker-compose.yaml
```
docker-compose.yaml
Dockerfile
```

**Recommended:**
```
deployments/
├── docker/
│   ├── Dockerfile
│   ├── Dockerfile.dev      # Development variant
│   ├── docker-compose.yaml
│   └── docker-compose.prod.yaml
├── kubernetes/
│   ├── base/               # Kustomize base
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── configmap.yaml
│   └── overlays/
│       ├── dev/
│       ├── staging/
│       └── prod/
└── terraform/              # Infrastructure as Code
    ├── main.tf
    ├── variables.tf
    └── modules/
```

**Benefits:**
- Environment-specific configurations
- Production-ready deployment artifacts
- Infrastructure as Code
- Easy to manage multi-environment deployments

---

### 15. Add Development Tools

**Recommended:**
```
tools/
├── tools.go                # Go tools dependencies
├── mockgen/                # Mock generation configs
├── protoc/                 # Proto generation scripts
└── linter/
    └── .golangci.yml      # Linter configuration

.devcontainer/              # VSCode dev container
├── devcontainer.json
└── Dockerfile

.vscode/                    # VSCode settings
├── settings.json
├── launch.json            # Debug configurations
└── tasks.json             # Build tasks
```

**Benefits:**
- Reproducible development environment
- Consistent tooling across team
- Easy onboarding for new developers
- IDE integration

---

## 📊 Structure Comparison

### Before (Current)
```
order-delivery-service/
├── cmd/server/              ✅ Good
├── internal/
│   ├── config/              ✅ Good
│   ├── constants/           ✅ Good
│   ├── domain/              ✅ Good
│   ├── service/             ✅ Improved (added repository.go)
│   ├── repository/postgres/ ✅ Good
│   └── transport/grpc/      ✅ Good
├── pkg/                     ⚠️  Could be better organized
├── api/grpc/                ⚠️  No versioning
├── migrations/              ✅ Good
└── config/                  ⚠️  Contains sensitive data
```

### After (Recommended Future State)
```
order-delivery-service/
├── cmd/
│   ├── server/              # Main application
│   ├── migrate/             # Migration CLI
│   └── seed/                # Data seeding CLI
├── internal/
│   ├── app/                 # Application layer
│   ├── domain/              # Domain layer (aggregates)
│   │   ├── delivery/
│   │   ├── address/
│   │   └── errors/
│   ├── service/             # Service layer
│   │   ├── delivery/
│   │   └── dto/
│   ├── repository/          # Repository implementations
│   │   └── postgres/
│   ├── transport/           # Transport layer
│   │   └── grpc/
│   ├── infrastructure/      # Infrastructure concerns
│   │   ├── database/
│   │   └── cache/
│   ├── config/              # Configuration
│   └── constants/           # Constants
├── pkg/                     # Public reusable packages
│   ├── observability/
│   │   ├── logging/
│   │   ├── metrics/
│   │   └── tracing/
│   ├── validation/
│   └── middleware/
├── api/                     # API definitions
│   ├── v1/
│   └── docs/
├── tests/                   # Test suites
│   ├── integration/
│   ├── e2e/
│   └── testutil/
├── scripts/                 # Automation scripts
│   ├── build/
│   ├── db/
│   └── dev/
├── deployments/             # Deployment configs
│   ├── docker/
│   ├── kubernetes/
│   └── terraform/
├── migrations/              # Database migrations
├── tools/                   # Development tools
├── docs/                    # Documentation
└── config.example.yaml      # Example config (no secrets)
```

---

## 🎯 Migration Priority

### Phase 1: Critical (Do Now)
1. ✅ **Repository interface inversion** - COMPLETED
2. ⏳ Move sensitive data out of config files
3. ⏳ Add .dockerignore and .env support

### Phase 2: High Priority (Next Sprint)
4. ⏳ Create application layer (internal/app/)
5. ⏳ Add integration tests structure
6. ⏳ Organize observability package
7. ⏳ Add Kubernetes manifests

### Phase 3: Medium Priority (This Quarter)
8. ⏳ Organize domain by aggregates (if complexity grows)
9. ⏳ Add DTO package
10. ⏳ Organize middleware by protocol
11. ⏳ Add API versioning

### Phase 4: Nice to Have (When Needed)
12. ⏳ Split config package
13. ⏳ Add infrastructure package
14. ⏳ Create scripts organization
15. ⏳ Add development tools setup

---

## 🏆 Best Practices Achieved

✅ **Dependency Inversion Principle** - Service owns repository interface
✅ **Clean Architecture** - Clear separation of layers
✅ **Interface Segregation** - Small, focused interfaces
✅ **Single Responsibility** - Each package has one reason to change
✅ **Testability** - Easy to mock and test
✅ **Go Project Layout** - Follows golang-standards/project-layout
✅ **Domain-Driven Design** - Rich domain models

---

## 📚 References

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Domain-Driven Design](https://martinfowler.com/tags/domain%20driven%20design.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)

---

**Last Updated:** October 23, 2025
**Status:** Phase 1 Complete ✅
