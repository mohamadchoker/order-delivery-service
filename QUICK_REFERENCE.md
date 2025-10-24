# Quick Reference Guide

## Common Commands

### Development

```bash
# Setup project
./scripts/setup.sh

# Generate proto files
make proto

# Run the service
make run

# Build binary
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Run linter
make lint
```

### Database

```bash
# Start PostgreSQL (Docker)
docker-compose up -d postgres

# Create database
createdb order_delivery_db

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create NAME=add_new_column
```

### Docker

```bash
# Build image
make docker-build

# Start all services
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs
```

### Testing with grpcurl

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Describe service
grpcurl -plaintext localhost:50051 describe delivery.DeliveryService

# Health check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# Create delivery
grpcurl -plaintext -d '{"order_id":"ORDER-123","pickup_address":{"city":"NYC"},"delivery_address":{"city":"Boston"},"scheduled_pickup_time":"2024-01-15T10:00:00Z","estimated_delivery_time":"2024-01-15T14:00:00Z"}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment

# Get delivery
grpcurl -plaintext -d '{"id":"<UUID>"}' localhost:50051 delivery.DeliveryService/GetDeliveryAssignment

# Update status
grpcurl -plaintext -d '{"id":"<UUID>","status":"DELIVERY_STATUS_ASSIGNED"}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus

# List deliveries
grpcurl -plaintext -d '{"page":1,"page_size":10}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

## Project Structure at a Glance

```
order-delivery-service/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── delivery/               # gRPC handlers
│   ├── entity/                 # Domain models
│   ├── repository/             # Data access
│   └── usecase/                # Business logic
├── pkg/                        # Public libraries
├── proto/                      # gRPC definitions
├── migrations/                 # DB migrations
└── config/config.yaml         # Configuration
```

## Key Files

- `cmd/server/main.go` - Application entry point
- `internal/delivery/grpc_handler.go` - gRPC API handlers
- `internal/usecase/delivery_usecase.go` - Business logic
- `internal/repository/postgres_repository.go` - Database layer
- `internal/entity/delivery.go` - Domain model
- `proto/delivery.proto` - API definition
- `Makefile` - Build commands
- `docker-compose.yaml` - Docker setup

## Environment Variables

```bash
# Database
DELIVERY_DATABASE_HOST=localhost
DELIVERY_DATABASE_PORT=5432
DELIVERY_DATABASE_USER=postgres
DELIVERY_DATABASE_PASSWORD=postgres
DELIVERY_DATABASE_DBNAME=order_delivery_db
DELIVERY_DATABASE_SSLMODE=disable

# Server
DELIVERY_SERVER_PORT=50051
DELIVERY_SERVER_SHUTDOWN_TIMEOUT=30s

# Logger
DELIVERY_LOGGER_LEVEL=info
DELIVERY_LOGGER_DEVELOPMENT=false
```

## Status Codes

- `OK` - Success
- `INVALID_ARGUMENT` - Bad request
- `NOT_FOUND` - Resource not found
- `FAILED_PRECONDITION` - Invalid state transition
- `INTERNAL` - Server error

## Delivery Status Flow

```
PENDING → ASSIGNED → PICKED_UP → IN_TRANSIT → DELIVERED
   ↓         ↓          ↓            ↓
CANCELLED  CANCELLED  FAILED      FAILED
```

## Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests
go test -tags=integration ./tests/integration/...

# With race detection
go test -race ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Troubleshooting

### Port already in use
```bash
# Find process
lsof -i :50051
# Kill process
kill -9 <PID>
```

### Database connection failed
```bash
# Check if PostgreSQL is running
docker-compose ps
# Start PostgreSQL
docker-compose up -d postgres
```

### Proto generation failed
```bash
# Install protoc
brew install protobuf  # macOS
apt-get install protobuf-compiler  # Linux

# Install Go plugins
make install-tools
```

## Performance Tips

1. Increase database connection pool:
   - `max_open_conns: 50`
   - `max_idle_conns: 10`

2. Enable gRPC connection pooling on client side

3. Add caching layer for frequently accessed data

4. Use database read replicas for read-heavy workloads

5. Monitor with Prometheus + Grafana

## Security Checklist

- [ ] TLS enabled for gRPC
- [ ] Database uses SSL
- [ ] Secrets in environment variables
- [ ] Input validation enabled
- [ ] Rate limiting configured
- [ ] Authentication middleware added
- [ ] Authorization checks in place
- [ ] Audit logging enabled

## Deployment Checklist

- [ ] Configuration for environment
- [ ] Database migrations applied
- [ ] Environment variables set
- [ ] Health checks working
- [ ] Logging configured
- [ ] Monitoring set up
- [ ] Backups configured
- [ ] Load balancer configured
- [ ] Auto-scaling enabled
- [ ] CI/CD pipeline tested

## Useful Links

- gRPC docs: https://grpc.io/docs/
- GORM docs: https://gorm.io/docs/
- Protocol Buffers: https://protobuf.dev/
- Go best practices: https://go.dev/doc/effective_go
