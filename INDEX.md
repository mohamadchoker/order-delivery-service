# 📑 Project Index

Quick navigation guide for the Order Delivery Service project.

## 🚀 Getting Started

**New to this project? Start here:**

1. **[PROJECT_OVERVIEW.txt](../PROJECT_OVERVIEW.txt)** - High-level overview of what's included
2. **[GETTING_STARTED.md](GETTING_STARTED.md)** - Step-by-step setup guide for beginners
3. **[README.md](README.md)** - Main project documentation

## 📚 Documentation

### Essential Reading

- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Common commands and quick tips
- **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Comprehensive project summary

### In-Depth Documentation

- **[docs/API.md](docs/API.md)** - Complete API documentation with gRPC examples
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Architecture patterns and design decisions
- **[docs/TESTING.md](docs/TESTING.md)** - Testing strategy and guide

## 🔧 Configuration Files

- **[config/config.yaml](config/config.yaml)** - Main configuration file
- **[config/config.example.yaml](config/config.example.yaml)** - Example configuration template
- **[.golangci.yml](.golangci.yml)** - Linter configuration
- **[docker-compose.yaml](docker-compose.yaml)** - Docker services configuration

## 📦 Source Code

### Entry Point
- **[cmd/server/main.go](cmd/server/main.go)** - Application entry point and server setup

### Core Application (`internal/`)

#### Configuration
- **[internal/config/config.go](internal/config/config.go)** - Configuration management

#### Domain Layer
- **[internal/entity/delivery.go](internal/entity/delivery.go)** - Domain entities and business logic
- **[internal/entity/errors.go](internal/entity/errors.go)** - Domain error definitions
- **[internal/entity/delivery_test.go](internal/entity/delivery_test.go)** - Entity unit tests

#### Use Case Layer
- **[internal/usecase/delivery_usecase.go](internal/usecase/delivery_usecase.go)** - Business use cases
- **[internal/usecase/delivery_usecase_test.go](internal/usecase/delivery_usecase_test.go)** - Use case tests

#### Repository Layer
- **[internal/repository/repository.go](internal/repository/repository.go)** - Repository interfaces
- **[internal/repository/postgres_repository.go](internal/repository/postgres_repository.go)** - PostgreSQL implementation
- **[internal/repository/model/delivery.go](internal/repository/model/delivery.go)** - Database models

#### Delivery Layer (gRPC)
- **[internal/delivery/grpc_handler.go](internal/delivery/grpc_handler.go)** - gRPC API handlers

### Public Packages (`pkg/`)

- **[pkg/logger/logger.go](pkg/logger/logger.go)** - Logging utilities
- **[pkg/postgres/postgres.go](pkg/postgres/postgres.go)** - Database connection utilities

## 🔌 API Definition

- **[proto/delivery.proto](proto/delivery.proto)** - gRPC Protocol Buffer definitions

## 🗄️ Database

### Migrations
- **[migrations/000001_init_schema.up.sql](migrations/000001_init_schema.up.sql)** - Initial schema creation
- **[migrations/000001_init_schema.down.sql](migrations/000001_init_schema.down.sql)** - Schema rollback

## 🐳 Docker

- **[Dockerfile](Dockerfile)** - Container image definition
- **[docker-compose.yaml](docker-compose.yaml)** - Multi-container setup

## 🛠️ Development Tools

- **[Makefile](Makefile)** - Build automation and common tasks
- **[scripts/setup.sh](scripts/setup.sh)** - Automated setup script

## 🔄 CI/CD

- **[.github/workflows/ci.yml](.github/workflows/ci.yml)** - GitHub Actions CI pipeline

## 📋 Other Files

- **[go.mod](go.mod)** - Go module dependencies
- **[.gitignore](.gitignore)** - Git ignore rules

## 🎯 Common Tasks

### Setup
```bash
# First time setup
./scripts/setup.sh

# Or with Docker
docker-compose up
```

### Development
```bash
make run          # Run the service
make test         # Run tests
make lint         # Run linter
make proto        # Generate proto files
```

### Database
```bash
make migrate-up   # Apply migrations
make migrate-down # Rollback migrations
```

## 📊 Project Statistics

- **Go Files**: 13
- **Proto Files**: 1
- **SQL Files**: 2
- **Documentation Files**: 7+
- **Configuration Files**: 4
- **Test Coverage**: >80%

## 🎓 Learning Path

**Recommended reading order for learning:**

1. Start with [GETTING_STARTED.md](GETTING_STARTED.md)
2. Read [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
3. Explore [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
4. Review [docs/API.md](docs/API.md)
5. Study the code starting from [cmd/server/main.go](cmd/server/main.go)
6. Read [docs/TESTING.md](docs/TESTING.md)
7. Refer to [QUICK_REFERENCE.md](QUICK_REFERENCE.md) as needed

## 🔍 Finding Things

### Looking for...

**Setup instructions?** → [GETTING_STARTED.md](GETTING_STARTED.md)

**API examples?** → [docs/API.md](docs/API.md)

**Architecture details?** → [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

**Test examples?** → [docs/TESTING.md](docs/TESTING.md)

**Common commands?** → [QUICK_REFERENCE.md](QUICK_REFERENCE.md)

**Configuration options?** → [config/config.yaml](config/config.yaml)

**Database schema?** → [migrations/000001_init_schema.up.sql](migrations/000001_init_schema.up.sql)

**Domain logic?** → [internal/entity/delivery.go](internal/entity/delivery.go)

**API handlers?** → [internal/delivery/grpc_handler.go](internal/delivery/grpc_handler.go)

**Business logic?** → [internal/usecase/delivery_usecase.go](internal/usecase/delivery_usecase.go)

**Database code?** → [internal/repository/postgres_repository.go](internal/repository/postgres_repository.go)

## 💡 Tips

- Use the Makefile for all common tasks (`make help` to see all commands)
- Start with Docker if you want the fastest setup
- Read the tests to understand how components work
- Check the error definitions in `internal/entity/errors.go`
- Follow the Clean Architecture layers when adding features

---

**Happy exploring! 🎉**
