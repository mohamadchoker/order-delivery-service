# Getting Started Guide

Welcome to the Order Delivery Service! This guide will help you get up and running quickly.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

1. **Go 1.21 or later**
   ```bash
   go version  # Should show 1.21 or higher
   ```
   Download: https://golang.org/dl/

2. **PostgreSQL 14 or later** (or Docker)
   ```bash
   psql --version
   ```
   Download: https://www.postgresql.org/download/
   
   OR use Docker:
   ```bash
   docker --version
   ```

3. **Protocol Buffer Compiler (protoc)**
   ```bash
   protoc --version
   ```
   
   Install:
   - **macOS**: `brew install protobuf`
   - **Linux**: `apt-get install protobuf-compiler`
   - **Windows**: Download from https://github.com/protocolbuffers/protobuf/releases

4. **grpcurl** (for testing, optional)
   ```bash
   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
   ```

## 🚀 Quick Start (5 minutes)

### Option 1: Using Docker (Easiest)

```bash
# 1. Navigate to project directory
cd order-delivery-service

# 2. Start everything with Docker Compose
docker-compose up

# 3. Service is now running on localhost:50051
```

That's it! Skip to "Testing Your Setup" below.

### Option 2: Local Setup

```bash
# 1. Navigate to project directory
cd order-delivery-service

# 2. Run the setup script (installs tools, dependencies, etc.)
chmod +x scripts/setup.sh
./scripts/setup.sh

# 3. Start PostgreSQL (if not already running)
# Using Docker:
docker-compose up -d postgres

# OR if you have PostgreSQL installed locally:
# Create the database
createdb order_delivery_db

# 4. Run database migrations
make migrate-up

# 5. Generate protocol buffer files
make proto

# 6. Start the service
make run
```

The service will start on `localhost:50051`.

## ✅ Testing Your Setup

Once the service is running, test it with these commands:

### 1. Check if service is running

```bash
grpcurl -plaintext localhost:50051 list
```

You should see:
```
delivery.DeliveryService
grpc.health.v1.Health
grpc.reflection.v1alpha.ServerReflection
```

### 2. Check health

```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

Expected response:
```json
{
  "status": "SERVING"
}
```

### 3. Create a delivery assignment

```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER-12345",
  "pickup_address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "postal_code": "10001",
    "country": "USA",
    "latitude": 40.7128,
    "longitude": -74.0060
  },
  "delivery_address": {
    "street": "456 Oak Ave",
    "city": "Boston",
    "state": "MA",
    "postal_code": "02101",
    "country": "USA",
    "latitude": 42.3601,
    "longitude": -71.0589
  },
  "scheduled_pickup_time": "2024-12-31T10:00:00Z",
  "estimated_delivery_time": "2024-12-31T14:00:00Z",
  "notes": "First test delivery"
}' localhost:50051 delivery.DeliveryService/CreateDeliveryAssignment
```

If successful, you'll see a JSON response with the created delivery assignment including an ID!

### 4. List deliveries

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 delivery.DeliveryService/ListDeliveryAssignments
```

## 📚 Next Steps

Now that you have the service running, explore these features:

### 1. Update Delivery Status

```bash
# Replace <DELIVERY_ID> with actual ID from creation response
grpcurl -plaintext -d '{
  "id": "<DELIVERY_ID>",
  "status": "DELIVERY_STATUS_ASSIGNED",
  "notes": "Driver John assigned"
}' localhost:50051 delivery.DeliveryService/UpdateDeliveryStatus
```

### 2. Assign a Driver

```bash
grpcurl -plaintext -d '{
  "id": "<DELIVERY_ID>",
  "driver_id": "DRIVER-789"
}' localhost:50051 delivery.DeliveryService/AssignDriver
```

### 3. Get Delivery Details

```bash
grpcurl -plaintext -d '{
  "id": "<DELIVERY_ID>"
}' localhost:50051 delivery.DeliveryService/GetDeliveryAssignment
```

### 4. Get Metrics

```bash
grpcurl -plaintext -d '{
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-12-31T23:59:59Z"
}' localhost:50051 delivery.DeliveryService/GetDeliveryMetrics
```

## 🧪 Running Tests

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run linter
make lint
```

## 📖 Learn More

- **API Documentation**: See `docs/API.md` for complete API reference
- **Architecture**: See `docs/ARCHITECTURE.md` for design details
- **Testing Guide**: See `docs/TESTING.md` for testing best practices
- **Quick Reference**: See `QUICK_REFERENCE.md` for common commands

## 🛠️ Development Workflow

### Making Changes

1. **Create a new feature**
   ```bash
   # Make your changes to the code
   
   # Generate proto files if you modified proto
   make proto
   
   # Run tests
   make test
   
   # Run linter
   make lint
   
   # Build
   make build
   ```

2. **Create database migration**
   ```bash
   make migrate-create NAME=add_new_column
   # Edit the generated migration files in migrations/
   make migrate-up
   ```

3. **Test your changes**
   ```bash
   # Unit tests
   go test ./internal/...
   
   # Integration tests (requires DB)
   make test-integration
   ```

## 🐛 Troubleshooting

### Service won't start

**Problem**: Port 50051 is already in use

**Solution**:
```bash
# Find and kill the process
lsof -i :50051
kill -9 <PID>
```

### Database connection error

**Problem**: Can't connect to PostgreSQL

**Solution**:
```bash
# Check if PostgreSQL is running
docker-compose ps

# If not running, start it
docker-compose up -d postgres

# Check logs
docker-compose logs postgres
```

### Proto generation fails

**Problem**: `protoc` command not found

**Solution**:
```bash
# Install protoc
# macOS:
brew install protobuf

# Linux:
sudo apt-get install protobuf-compiler

# Then install Go plugins
make install-tools
```

### Tests failing

**Problem**: Database tests fail

**Solution**:
```bash
# Ensure test database is running
docker-compose up -d postgres

# Run migrations
make migrate-up

# Try tests again
make test
```

## 🎯 Common Tasks

### Restart the service
```bash
# Stop: Ctrl+C
# Start: make run
```

### Reset database
```bash
make migrate-down
make migrate-up
```

### View logs
```bash
# Application logs are printed to stdout
# With Docker:
docker-compose logs -f service
```

### Clean build artifacts
```bash
make clean
```

## 📞 Getting Help

- Check the documentation in `/docs`
- Review the `QUICK_REFERENCE.md` for common commands
- Look at the code examples in test files
- Check the `PROJECT_SUMMARY.md` for an overview

## 🎉 Success!

You now have a fully functional Order Delivery Service running! 

**What you can do next:**
- Explore the API with different requests
- Read the architecture documentation
- Modify the code and add new features
- Write tests for your changes
- Deploy to production (see deployment guides)

Happy coding! 🚀
