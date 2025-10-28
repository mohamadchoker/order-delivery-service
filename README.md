# Order Delivery Service

[![CI](https://github.com/company/order-delivery-service/workflows/CI/badge.svg)](https://github.com/company/order-delivery-service/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/company/order-delivery-service)](https://goreportcard.com/report/github.com/company/order-delivery-service)
[![codecov](https://codecov.io/gh/company/order-delivery-service/branch/main/graph/badge.svg)](https://codecov.io/gh/company/order-delivery-service)
[![License](https://img.shields.io/badge/License-Proprietary-blue.svg)](LICENSE)

An enterprise-grade microservice for managing order delivery assignments built with Go, gRPC, REST, and PostgreSQL, following Clean Architecture and Domain-Driven Design principles. Provides both gRPC and HTTP/REST APIs via grpc-gateway.

---

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Technology Stack](#-technology-stack)
- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Testing](#-testing)
- [Local CI Testing with Act](#-local-ci-testing-with-act)
- [Deployment](#-deployment)
- [API Documentation](#-api-documentation)
- [Configuration](#-configuration)
- [Database](#-database)
- [Monitoring](#-monitoring)
- [Contributing](#-contributing)
- [Documentation](#-documentation)

---

## ✨ Features

### Core Functionality
- ✅ **Delivery Management** - Create, update, and track delivery assignments
- ✅ **Driver Assignment** - Assign and reassign drivers to deliveries
- ✅ **Status Tracking** - Validated state transitions with business rules
- ✅ **Metrics & Analytics** - Real-time delivery metrics and performance tracking
- ✅ **Flexible Filtering** - List deliveries with pagination and filters (status, driver, date range)

### Enterprise Features
- ✅ **Clean Architecture** - Clear separation of concerns (Domain → Service → Repository → Transport)
- ✅ **Domain-Driven Design** - Business logic encapsulated in domain entities
- ✅ **Dual Protocol Support** - gRPC and REST/HTTP APIs from single implementation (grpc-gateway)
- ✅ **OpenAPI/Swagger** - Auto-generated API documentation
- ✅ **Auto-Generated Mocks** - Type-safe mocks using uber-go/mock
- ✅ **Comprehensive Validation** - Input validation at all layers
- ✅ **Request Tracing** - X-Request-ID tracking for end-to-end correlation
- ✅ **Timeout Handling** - Automatic timeout enforcement with configurable limits
- ✅ **Prometheus Metrics** - Production-ready observability
- ✅ **Transaction Support** - ACID-compliant database operations
- ✅ **Graceful Shutdown** - Clean service termination with in-flight request completion
- ✅ **Health Checks** - gRPC health protocol implementation
- ✅ **12-Factor App** - Environment-based configuration, no config files

### Code Quality
- ✅ **High Test Coverage** - >80% coverage with unit and integration tests
- ✅ **CI/CD Ready** - GitHub Actions with local testing via `act`
- ✅ **Linting** - golangci-lint with strict rules
- ✅ **Type Safety** - Comprehensive error handling with custom error types
- ✅ **Database Migrations** - Version-controlled schema changes
- ✅ **Docker Ready** - Multi-stage builds with optimized images
- ✅ **Production Logging** - Structured logging with Zap

---

## 🏗️ Architecture

### Clean Architecture Layers

```
                        gRPC Client          REST/HTTP Client
                             │                      │
                             ▼                      ▼
┌─────────────────────────────────────────────────────────┐
│                   Transport Layer                        │
│     gRPC Handlers              HTTP Gateway (Proxy)      │
│          Proto ↔ Domain Conversion                       │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                   Service Layer                          │
│              (Business Logic / Use Cases)                │
│         Orchestrates Repository Operations               │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                 Repository Layer                         │
│              (Data Access Interface)                     │
│         Entity ↔ DB Model Conversion                     │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                   Domain Layer                           │
│         (Business Entities & Domain Logic)               │
│          Status Validation, Business Rules               │
└─────────────────────────────────────────────────────────┘
```

### Key Architectural Patterns

- **Dependency Inversion** - Service layer owns repository interface
- **Repository Pattern** - Data access abstracted behind interfaces
- **Domain-Driven Design** - Business logic in domain entities
- **Dependency Injection** - All dependencies injected via constructors
- **Interface Segregation** - Small, focused interfaces

### Data Flow Example

Creating a delivery assignment:

```
1. gRPC Request → Handler.CreateDeliveryAssignment()
2. Handler validates and converts Proto → Domain Entity
3. Handler calls Service.CreateDeliveryAssignment()
4. Service validates business rules
5. Service creates entity: domain.NewDeliveryAssignment()
6. Service calls Repository.Create()
7. Repository converts Entity → DB Model
8. Repository persists to PostgreSQL
9. Response: DB Model → Entity → Proto → gRPC Response
```

---

## 🔧 Technology Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.24+ | Core application language |
| **API** | gRPC + REST | - | Dual protocol support |
| **Gateway** | grpc-gateway | 2.x | gRPC to REST/HTTP proxy |
| **Protocol** | Protocol Buffers | 3 | API definition and serialization |
| **OpenAPI** | Swagger | 2.0 | Auto-generated API docs |
| **Database** | PostgreSQL | 14+ | Primary data store |
| **ORM** | GORM | 1.25+ | Database abstraction |
| **Logging** | Uber Zap | 1.27+ | Structured logging |
| **Config** | stdlib | - | Simple environment variables |
| **Testing** | testify + gomock | - | Assertions and mocking |
| **Mocks** | uber-go/mock | - | Auto-generated type-safe mocks |
| **Migrations** | golang-migrate | 4.17+ | Schema versioning |
| **Metrics** | Prometheus | - | Application monitoring |
| **Validation** | Custom | - | Fluent validation API |
| **Container** | Docker | 20+ | Containerization |
| **Orchestration** | Docker Compose | - | Local development |
| **CI/CD** | GitHub Actions | - | Continuous integration |
| **Local CI** | act | - | Test GitHub Actions locally |
| **Linter** | golangci-lint | 1.55+ | Code quality |

---

## 🚀 Quick Start

### Prerequisites

```bash
# Required
- Go 1.24 or higher
- PostgreSQL 14 or higher
- Protocol Buffer Compiler (protoc)

# Optional
- Docker & Docker Compose (recommended for local development)
- act (for local GitHub Actions testing)
- grpcurl (for API testing)
- curl (for testing REST API)
```

> **IDE Setup**: If using GoLand/IntelliJ and seeing proto import errors, see [IDE Setup Guide](docs/IDE_SETUP.md) to configure proto paths.

### Installation

#### 1. Install Go Tools

```bash
# Install development tools
make install-tools

# This installs:
# - protoc-gen-go (Protocol Buffer compiler for Go)
# - protoc-gen-go-grpc (gRPC plugin for protoc)
# - protoc-gen-grpc-gateway (gRPC to REST gateway generator)
# - protoc-gen-openapiv2 (OpenAPI/Swagger generator)
# - mockgen (Mock generator)
# - golang-migrate (Database migrations)
# - golangci-lint (Linter)
# - air (Hot reload for development)

# Ensure ~/go/bin is in your PATH
export PATH=$PATH:~/go/bin
# Add to ~/.zshrc or ~/.bashrc:
echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc
```

#### 2. Clone and Setup

```bash
# Clone repository
git clone https://github.com/company/order-delivery-service.git
cd order-delivery-service

# Install dependencies
go mod download

# Generate protocol buffers
make proto

# Generate mocks
make mocks
```

#### 3. Setup Database

**Option A: Using Docker (Recommended)**

```bash
# Start PostgreSQL
docker-compose up -d postgres

# Run migrations
make migrate-up
```

**Option B: Local PostgreSQL**

```bash
# Create database
createdb order_delivery_db

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=order_delivery_db

# Run migrations
make migrate-up
```

#### 4. Run the Service

**Option A: Development Mode with Hot-Reload** ⚡ (Recommended for development)

```bash
make dev
# Service starts with hot-reload enabled
# Any code change automatically restarts the service in 1-2 seconds!
# Perfect for rapid development
```

**Option B: Using Make**

```bash
make run
# Service starts on :50051 (gRPC)
# Metrics on :9090/metrics
```

**Option C: Using Docker Compose**

```bash
docker-compose up --build
```

**Option D: Directly**

```bash
# Set required environment variables
export DB_HOST=localhost
export DB_PASSWORD=postgres

# Run
go run cmd/server/main.go
```

#### 5. Test the Service

```bash
# Health check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# Expected response:
{
  "status": "SERVING"
}

# Create a delivery assignment (see API Documentation section for more examples)
grpcurl -plaintext -d '{
  "order_id": "ORDER-123",
  "pickup_address": {
    "street": "123 Main St",
    "city": "New York",
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
}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment
```

---

## 📁 Project Structure

```
order-delivery-service/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD workflow
├── .act/                       # Act configuration for local CI testing
│   └── .secrets               # Local secrets (not committed)
├── cmd/
│   └── server/
│       └── main.go            # Application entry point
├── internal/
│   ├── config/                # Configuration management
│   │   └── config.go         # Simple env-based config (no Viper!)
│   ├── constants/             # Application-wide constants
│   │   └── constants.go      # Pagination, timeouts, validation rules
│   ├── domain/                # Domain layer (core business logic)
│   │   ├── delivery.go       # DeliveryAssignment entity with business rules
│   │   └── errors.go         # Custom error types (DomainError, ValidationError, etc.)
│   ├── mocks/                 # Auto-generated mocks (DO NOT EDIT)
│   │   ├── repository_mock.go # Mock for DeliveryRepository
│   │   └── usecase_mock.go    # Mock for DeliveryUseCase
│   ├── repository/            # Repository layer (data access)
│   │   └── postgres/
│   │       ├── repository.go  # PostgreSQL implementation
│   │       └── model/
│   │           └── delivery.go # Database models (separate from domain)
│   ├── service/               # Service layer (business logic)
│   │   ├── delivery_usecase.go       # Use case implementation
│   │   ├── delivery_usecase_test.go  # Tests using auto-generated mocks
│   │   └── repository.go             # Repository interface (with go:generate)
│   └── transport/             # Transport layer (protocol handlers)
│       └── grpc/
│           └── handler.go     # gRPC request handlers
├── pkg/                       # Public reusable packages
│   ├── logger/                # Structured logging utilities
│   ├── metrics/               # Prometheus metrics
│   ├── middleware/            # Request ID, timeout middleware
│   ├── postgres/              # Database connection utilities
│   └── validator/             # Input validation (fluent API)
├── proto/                     # Protocol buffer definitions
│   ├── delivery.proto         # gRPC service definition
│   ├── delivery.pb.go         # Generated Go code
│   └── delivery_grpc.pb.go    # Generated gRPC code
├── migrations/                # Database migrations
│   ├── 000001_create_delivery_assignments.up.sql
│   └── 000001_create_delivery_assignments.down.sql
├── docs/                      # Documentation
│   ├── ACT_USAGE.md          # Local CI testing guide
│   ├── FINAL_IMPROVEMENTS.md # Refactoring summary
│   └── REFACTORING_COMPLETE.md # Test migration guide
├── scripts/                   # Utility scripts
├── .actrc                     # Act configuration
├── .env.example              # Environment variables template
├── .gitignore
├── docker-compose.yaml       # Docker Compose configuration
├── Dockerfile                # Multi-stage Docker build
├── go.mod                    # Go module dependencies
├── go.sum
├── Makefile                  # Development commands
├── CLAUDE.md                 # Project guidance for Claude Code
└── README.md                 # This file
```

### Key Files

| File | Purpose |
|------|---------|
| `internal/domain/delivery.go` | Core business entity with status validation |
| `internal/service/delivery_usecase.go` | Business logic orchestration |
| `internal/repository/postgres/repository.go` | Data access implementation |
| `internal/transport/grpc/handler.go` | gRPC request handlers |
| `pkg/validator/validator.go` | Input validation fluent API |
| `pkg/middleware/request_id.go` | Request tracing |
| `pkg/metrics/metrics.go` | Prometheus metrics |
| `proto/delivery.proto` | API definition |
| `.github/workflows/ci.yml` | CI/CD pipeline |
| `Makefile` | Development commands |

---

## 🛠️ Development

### Available Make Commands

```bash
make help              # Show all available commands

# Code Generation
make proto             # Generate Protocol Buffer code
make mocks             # Generate mocks using mockgen

# Build & Run
make build             # Build the application binary
make run               # Run the application locally
make clean             # Clean build artifacts

# Development with Hot-Reload ⚡
make dev               # Start dev environment with hot-reload (recommended!)
make dev-up            # Same as dev
make dev-down          # Stop dev environment
make dev-logs          # View dev logs (follow mode)
make dev-restart       # Restart dev service
make air               # Run Air locally (without Docker)
make install-air       # Install Air tool

# Testing
make test              # Run all tests with race detection
make test-coverage     # Run tests with coverage report (generates coverage.html)
make test-integration  # Run integration tests (requires PostgreSQL)

# Code Quality
make lint              # Run golangci-lint
make lint-fix          # Auto-fix linting issues

# Database
make migrate-up        # Apply all pending migrations
make migrate-down      # Rollback last migration
make migrate-create NAME=your_migration_name  # Create new migration

# Docker
make docker-build      # Build Docker image
make docker-up         # Start all services with Docker Compose
make docker-down       # Stop all Docker services
make docker-logs       # View docker logs

# Local CI Testing (with act)
make act-test          # Run tests locally using act
make act-all           # Run all CI jobs locally
make act-list          # List all available jobs
```

### Development Workflow

#### 1. Making Code Changes

```bash
# 1. Create a feature branch
git checkout -b feature/my-feature

# 2. Make your changes

# 3. Generate mocks if interfaces changed
make mocks

# 4. Run tests
make test

# 5. Run linter
make lint

# 6. Test locally with act (optional but recommended)
make act-test

# 7. Commit and push
git add .
git commit -m "feat: add new feature"
git push origin feature/my-feature
```

#### 2. Adding a New gRPC Endpoint

```bash
# 1. Update proto/delivery.proto
vim proto/delivery.proto

# 2. Generate proto code
make proto

# 3. Implement handler in internal/transport/grpc/handler.go

# 4. Add service method in internal/service/delivery_usecase.go

# 5. Update interface and regenerate mocks
make mocks

# 6. Write tests

# 7. Test
make test
```

#### 3. Adding a New Database Field

```bash
# 1. Create migration
make migrate-create NAME=add_new_field

# 2. Edit generated migration files in migrations/

# 3. Apply migration
make migrate-up

# 4. Update domain entity (internal/domain/delivery.go)

# 5. Update DB model (internal/repository/postgres/model/delivery.go)

# 6. Update proto (proto/delivery.proto) if needed

# 7. Regenerate proto
make proto

# 8. Update tests
```

#### 4. Working with Mocks

Mocks are **auto-generated** - never edit them manually!

```bash
# Generate all mocks
make mocks

# Mocks are generated from go:generate directives:
# internal/service/repository.go        → internal/mocks/repository_mock.go
# internal/service/delivery_usecase.go  → internal/mocks/usecase_mock.go

# Use in tests:
import "github.com/company/order-delivery-service/internal/mocks"

ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockRepo := mocks.NewMockDeliveryRepository(ctrl)
mockRepo.EXPECT().
    Create(gomock.Any(), gomock.Any()).
    Return(nil).
    Times(1)
```

---

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
# Opens coverage.html in browser

# Run specific package tests
go test -v ./internal/service/
go test -v ./internal/domain/

# Run with race detection (default in make test)
go test -race ./...

# Run integration tests (requires PostgreSQL)
make test-integration
```

### Test Coverage

Current coverage: **>80%**

```bash
# Generate coverage report
make test-coverage

# View coverage in terminal
go test -cover ./...

# Detailed coverage by function
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Writing Tests

#### Unit Tests with Auto-Generated Mocks

```go
package service_test

import (
    "testing"
    "github.com/company/order-delivery-service/internal/mocks"
    "github.com/company/order-delivery-service/internal/service"
    "go.uber.org/mock/gomock"
    "github.com/stretchr/testify/require"
)

func TestCreateDelivery(t *testing.T) {
    // Setup
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockDeliveryRepository(ctrl)
    logger, _ := zap.NewDevelopment()
    uc := service.NewDeliveryUseCase(mockRepo, logger)

    // Set expectations
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil).
        Times(1)

    // Execute
    result, err := uc.CreateDeliveryAssignment(ctx, input)

    // Assert
    require.NoError(t, err)
    require.NotNil(t, result)
}
```

#### Integration Tests

```go
//go:build integration

package integration_test

import (
    "testing"
    // Uses real database
)

func TestCreateDelivery_Integration(t *testing.T) {
    // Setup test database
    // Run actual operations
    // Verify in database
}
```

### Test Organization

- **Unit tests**: `*_test.go` in same package
- **Integration tests**: `tests/integration/` with `//go:build integration` tag
- **Test package suffix**: Use `package service_test` for black-box testing
- **Mocks**: Auto-generated in `internal/mocks/`

---

## 🎭 Local CI Testing with Act

### What is Act?

[Act](https://github.com/nektos/act) lets you run GitHub Actions workflows locally using Docker. This means:

✅ **Fast feedback** - Test CI before pushing
✅ **Save time** - No waiting for GitHub runners
✅ **Debug easily** - Direct access to containers
✅ **Cost effective** - No CI minutes used

### Installation

```bash
# macOS
brew install act

# Linux
curl -s https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Verify
act --version
```

### Basic Usage

```bash
# List all available jobs
act -l

# Run all CI jobs
act

# Run specific job
act -j test
act -j lint
act -j build
act -j docker

# Dry run (see what would execute)
act -n

# Verbose output for debugging
act -v
```

### Recommended Workflow

```bash
# 1. Make changes to code

# 2. Test locally with act
act -j test

# 3. If tests pass, run full CI
act

# 4. If all pass, push to GitHub
git push
```

### Common Commands

```bash
# Quick test before push
make act-test

# Run full CI locally
make act-all

# List available jobs
make act-list

# Debug failed job
act -j test -v
```

### Configuration

**`.actrc`** (project configuration):
```bash
-P ubuntu-latest=catthehacker/ubuntu:act-latest
--bind
--verbose
```

**`.act/.secrets`** (local secrets):
```bash
GITHUB_TOKEN=your_token_here
```

For detailed usage, see [docs/ACT_USAGE.md](docs/ACT_USAGE.md).

---

## 🚢 Deployment

### Docker

#### Build Image

```bash
# Build with version info
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  -t order-delivery-service:latest \
  .

# Or use make
make docker-build
```

#### Run Container

```bash
# Using Docker Compose (recommended)
docker-compose up -d

# Or Docker directly
docker run -d \
  -p 50051:50051 \
  -p 9090:9090 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=postgres \
  order-delivery-service:latest
```

### Environment Variables

All configuration is done via environment variables (12-factor app):

```bash
# Server
PORT=50051                    # gRPC port (default: 50051)
METRICS_PORT=9090            # Metrics port (default: 9090)
SHUTDOWN_TIMEOUT=30s         # Graceful shutdown timeout

# Database
DB_HOST=localhost            # Database host (required)
DB_PORT=5432                # Database port (default: 5432)
DB_USER=postgres            # Database user (default: postgres)
DB_PASSWORD=postgres        # Database password (required)
DB_NAME=order_delivery_db   # Database name
DB_SSLMODE=disable          # SSL mode (disable, require, verify-full)
DB_MAX_OPEN_CONNS=25        # Max open connections
DB_MAX_IDLE_CONNS=5         # Max idle connections
DB_CONN_MAX_LIFETIME=5m     # Connection max lifetime

# Logging
LOG_LEVEL=info              # Log level (debug, info, warn, error)
LOG_DEV=false               # Development mode (pretty printing)
```

### Docker Compose

```yaml
# docker-compose.yaml
version: '3.8'

services:
  service:
    build: .
    ports:
      - "50051:50051"
      - "9090:9090"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=postgres
      - LOG_LEVEL=info
    depends_on:
      - postgres

  postgres:
    image: postgres:14-alpine
    environment:
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=order_delivery_db
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### Kubernetes

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-delivery-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-delivery-service
  template:
    metadata:
      labels:
        app: order-delivery-service
    spec:
      containers:
      - name: service
        image: order-delivery-service:latest
        ports:
        - containerPort: 50051
          name: grpc
        - containerPort: 9090
          name: metrics
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DB_HOST
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
        livenessProbe:
          exec:
            command: ["/bin/grpc_health_probe", "-addr=:50051"]
          initialDelaySeconds: 10
        readinessProbe:
          exec:
            command: ["/bin/grpc_health_probe", "-addr=:50051"]
          initialDelaySeconds: 5
```

---

## 📡 API Documentation

### gRPC Service

> **💡 Enum Values - Both Formats Work!**
>
> You can use **either** short or long enum names in requests:
> - ✅ Short: `"status": "PENDING"` (recommended - cleaner)
> - ✅ Long: `"status": "DELIVERY_STATUS_PENDING"` (also works)
>
> **Available status values:**
> | Short Name | Long Name | Value |
> |------------|-----------|-------|
> | `PENDING` | `DELIVERY_STATUS_PENDING` | 1 |
> | `ASSIGNED` | `DELIVERY_STATUS_ASSIGNED` | 2 |
> | `PICKED_UP` | `DELIVERY_STATUS_PICKED_UP` | 3 |
> | `IN_TRANSIT` | `DELIVERY_STATUS_IN_TRANSIT` | 4 |
> | `DELIVERED` | `DELIVERY_STATUS_DELIVERED` | 5 |
> | `FAILED` | `DELIVERY_STATUS_FAILED` | 6 |
> | `CANCELLED` | `DELIVERY_STATUS_CANCELLED` | 7 |
>
> **Note:** Responses always show the long format (gRPC standard).
>
> For complete API examples, see [`docs/API_EXAMPLES.md`](docs/API_EXAMPLES.md)

#### List All Services

```bash
grpcurl -plaintext localhost:50051 list
```

#### Health Check

```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# Response:
{
  "status": "SERVING"
}
```

#### Create Delivery Assignment

```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER-12345",
  "pickup_address": {
    "street": "123 Main St",
    "city": "New York",
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
  "estimated_delivery_time": "2024-01-15T14:00:00Z",
  "notes": "Handle with care"
}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment
```

#### Get Delivery Assignment

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 delivery.DeliveryService/GetDeliveryAssignment
```

#### Assign Driver

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/AssignDriver
```

#### Update Status

**You can use short names** (recommended):

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PICKED_UP",
  "notes": "Package collected"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

#### List Deliveries

**You can use short names** (recommended):

```bash
# List all deliveries
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments

# Filter by status (short name - clean!)
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "status": "PENDING"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments

# Filter by driver
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20,
  "driver_id": "DRIVER-123"
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

#### Get Metrics

```bash
grpcurl -plaintext -d '{
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-01-31T23:59:59Z"
}' localhost:50051 delivery.DeliveryService/GetDeliveryMetrics
```

### Status Flow

Valid status transitions:

```
PENDING → ASSIGNED → PICKED_UP → IN_TRANSIT → DELIVERED
   ↓         ↓          ↓            ↓
CANCELLED  CANCELLED  FAILED      FAILED
```

Terminal states (no further transitions): **DELIVERED**, **FAILED**, **CANCELLED**

---

## ⚙️ Configuration

### Configuration Philosophy

This service follows **12-factor app principles**:

- ✅ All configuration via environment variables
- ✅ No config files in production
- ✅ Simple stdlib-based loading (no Viper!)
- ✅ Sensible defaults
- ✅ Validation on startup

### Environment Variables Reference

See [Environment Variables](#environment-variables) section above for complete list.

### Configuration Loading

```go
// internal/config/config.go
cfg, err := config.Load()  // Reads from environment

// Helper functions:
func getEnv(key, default string) string
func getEnvAsInt(key string, default int) int
func getEnvAsBool(key string, default bool) bool
func getEnvAsDuration(key string, default time.Duration) time.Duration
```

### Local Development

Use `.env.example` as a template:

```bash
cp .env.example .env
# Edit .env with your values
# Then: export $(cat .env | xargs)
```

---

## 🗄️ Database

### Schema

```sql
CREATE TABLE delivery_assignments (
    id UUID PRIMARY KEY,
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
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    INDEX idx_order_id (order_id),
    INDEX idx_driver_id (driver_id),
    INDEX idx_status (status),
    INDEX idx_scheduled_pickup_time (scheduled_pickup_time),
    INDEX idx_deleted_at (deleted_at)
);
```

### Migrations

```bash
# Create new migration
make migrate-create NAME=add_tracking_number

# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" version
```

### Connection Pooling

Configured via environment variables:

```bash
DB_MAX_OPEN_CONNS=25     # Max open connections (default: 25)
DB_MAX_IDLE_CONNS=5      # Max idle connections (default: 5)
DB_CONN_MAX_LIFETIME=5m  # Connection max lifetime (default: 5m)
```

---

## 📊 Monitoring

### Prometheus Metrics

Metrics exposed at `http://localhost:9090/metrics`

#### Available Metrics

**gRPC Metrics:**
```
# Total requests
order_delivery_service_grpc_requests_total{method="CreateDeliveryAssignment",code="OK"}

# Request duration (histogram)
order_delivery_service_grpc_request_duration_seconds{method="CreateDeliveryAssignment"}

# Active requests (gauge)
order_delivery_service_grpc_requests_active{method="CreateDeliveryAssignment"}
```

**Business Metrics:**
```
# Delivery operations
order_delivery_service_delivery_assignments_total{status="PENDING",operation="create"}
```

**Database Metrics:**
```
# Query count
order_delivery_service_database_queries_total{operation="create_delivery",status="success"}

# Query duration
order_delivery_service_database_query_duration_seconds{operation="create_delivery"}
```

#### Querying Metrics

```bash
# View all metrics
curl http://localhost:9090/metrics

# Grep specific metric
curl -s http://localhost:9090/metrics | grep grpc_requests_total
```

#### Prometheus Queries

```promql
# Request rate (requests per second)
rate(order_delivery_service_grpc_requests_total[5m])

# Error rate
rate(order_delivery_service_grpc_requests_total{code!="OK"}[5m])

# P95 latency
histogram_quantile(0.95, order_delivery_service_grpc_request_duration_seconds_bucket)

# Active deliveries by status
order_delivery_service_delivery_assignments_total
```

### Logging

Structured logging with Zap:

```json
{
  "level": "info",
  "ts": 1704067200.123,
  "caller": "service/delivery_usecase.go:58",
  "msg": "Creating delivery assignment",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "order_id": "ORDER-123"
}
```

Log levels: **debug**, **info**, **warn**, **error**

```bash
# Set log level
export LOG_LEVEL=debug

# Enable development mode (pretty printing)
export LOG_DEV=true
```

### Request Tracing

Every request gets a unique ID for end-to-end tracing:

```bash
# Send request with X-Request-ID header
grpcurl -H "X-Request-ID: my-trace-id" ...

# Or let system generate UUID
# All logs will include request_id field

# Grep logs by request ID
grep "request_id=550e8400-e29b-41d4-a716-446655440000" logs.txt
```

---

## 🤝 Contributing

### Development Process

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes
4. **Generate** mocks if needed (`make mocks`)
5. **Test** locally (`make test`)
6. **Lint** your code (`make lint`)
7. **Test** with act (`make act-test`)
8. **Commit** your changes (`git commit -m 'feat: add amazing feature'`)
9. **Push** to branch (`git push origin feature/amazing-feature`)
10. **Open** a Pull Request

### Coding Standards

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `golangci-lint` (`make lint`)
- Write tests for new features
- Maintain >80% code coverage
- Add comments for exported functions
- Use semantic commit messages

### Commit Message Format

```
<type>: <description>

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Examples:
```
feat: add delivery cancellation endpoint
fix: handle nil pointer in status update
docs: update API examples in README
test: add integration tests for metrics
```

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [README.md](README.md) | This file - Getting started and overview |
| [CLAUDE.md](CLAUDE.md) | Project guidance for Claude Code |
| [docs/GOLAND_QUICKSTART.md](docs/GOLAND_QUICKSTART.md) | **🎯 GoLand/IntelliJ quick setup (5 min)** |
| [docs/EDITOR_SETUP.md](docs/EDITOR_SETUP.md) | **Editor setup (VS Code, GoLand, Vim)** |
| [docs/LOGGING.md](docs/LOGGING.md) | **Logging configuration & troubleshooting** |
| [docs/CI_CD.md](docs/CI_CD.md) | **CI/CD pipeline & GitHub Actions** |
| [docs/MIGRATIONS.md](docs/MIGRATIONS.md) | **Database migrations guide** |
| [docs/ADDING_NEW_SERVICE.md](docs/ADDING_NEW_SERVICE.md) | Step-by-step guide to adding new gRPC services |
| [docs/NAMING_CONVENTIONS.md](docs/NAMING_CONVENTIONS.md) | Naming conventions quick reference |
| [docs/DOCKER_EXPLAINED.md](docs/DOCKER_EXPLAINED.md) | Docker & Docker Compose concepts explained |
| [docs/HOT_RELOAD_DEV.md](docs/HOT_RELOAD_DEV.md) | Hot-reload development setup with Air |
| [docs/ACT_USAGE.md](docs/ACT_USAGE.md) | Complete guide to local CI testing with act |
| [docs/API_EXAMPLES.md](docs/API_EXAMPLES.md) | gRPC API usage examples |
| [docs/FINAL_IMPROVEMENTS.md](docs/FINAL_IMPROVEMENTS.md) | Summary of major refactoring (proto location, config, mocks, etc.) |
| [docs/REFACTORING_COMPLETE.md](docs/REFACTORING_COMPLETE.md) | Test migration to auto-generated mocks |

### Additional Resources

- **gRPC Documentation**: https://grpc.io/docs/languages/go/
- **Protocol Buffers**: https://developers.google.com/protocol-buffers
- **GORM**: https://gorm.io/docs/
- **uber-go/mock**: https://github.com/uber-go/mock
- **Act**: https://github.com/nektos/act
- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---

## 📝 License

Copyright © 2024. All rights reserved.

---

## 🎯 Project Status

### Latest Updates (2025)

✅ **Complete refactoring to enterprise Go best practices**
- Moved protos from `api/grpc/` to `proto/`
- Simplified config (removed Viper, using simple env vars)
- Renamed `postgres_repository.go` to `repository.go`
- Added auto-generated mocks with uber-go/mock
- Migrated all tests to use auto-generated mocks
- Updated Docker to Go 1.24
- Added act configuration for local CI testing
- Updated GitHub Actions workflow

✅ **All tests passing**
✅ **Build successful**
✅ **Docker working**
✅ **Production ready**

### Roadmap

Future enhancements:
- [ ] Event sourcing for audit trail
- [ ] Real-time updates via Server-Sent Events
- [ ] GraphQL API layer
- [ ] Distributed tracing with OpenTelemetry
- [ ] Circuit breaker pattern
- [ ] Rate limiting
- [ ] Multi-tenancy support

---

## 💬 Support

For issues and questions:

1. **Check documentation** in `/docs` folder
2. **Search existing issues** on GitHub
3. **Create a new issue** with:
   - Clear description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)

---

## 🙏 Acknowledgments

Built with:
- [Go](https://golang.org/) - The Go Programming Language
- [gRPC](https://grpc.io/) - High-performance RPC framework
- [GORM](https://gorm.io/) - ORM library for Go
- [Zap](https://github.com/uber-go/zap) - Blazing fast structured logging
- [uber-go/mock](https://github.com/uber-go/mock) - Auto-generated mocks
- [act](https://github.com/nektos/act) - Local GitHub Actions testing
- [PostgreSQL](https://www.postgresql.org/) - The world's most advanced open source database

---

**Built with ❤️ using Go and enterprise best practices** 🚀
