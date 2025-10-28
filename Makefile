.PHONY: help proto mocks build run test test-coverage lint lint-fix format migrate-up migrate-down migrate-create clean docker-build docker-up docker-down docker-logs dev dev-up dev-down dev-logs dev-restart air install-air act-test act-all act-list act-lint act-build act-clean install-tools

# Variables
BINARY_NAME=order-delivery-service
PROTO_DIR=proto
MIGRATIONS_DIR=migrations
DB_URL=postgres://postgres:postgres@localhost:5432/order_delivery_db?sslmode=disable

# Ensure ~/go/bin is in PATH
export PATH := $(HOME)/go/bin:$(PATH)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

proto: ## Generate Go code from proto files
	@echo "Generating proto files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=. --openapiv2_opt=allow_merge=true,merge_file_name=api \
		--proto_path=. --proto_path=third_party \
		$(PROTO_DIR)/*.proto
	@echo "Proto files generated successfully"

mocks: ## Generate mocks using mockgen
	@echo "Generating mocks..."
	@go generate ./internal/service/...
	@echo "Mocks generated successfully"

build: proto ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: bin/$(BINARY_NAME)"

run: ## Run the application
	@echo "Starting $(BINARY_NAME)..."
	@go run ./cmd/server

test: ## Run unit tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

lint: ## Run golangci-lint
	@echo "Running linter..."
	@$(HOME)/go/bin/golangci-lint run ./...

lint-fix: ## Fix linting issues and format imports
	@echo "Fixing linting issues..."
	@$(HOME)/go/bin/golangci-lint run --fix ./...

format: ## Format code and organize imports
	@echo "Formatting code and organizing imports..."
	@$(HOME)/go/bin/goimports -local  github.com/mohamadchoker/order-delivery-service -w $$(find . -name "*.go" -not -path "./vendor/*" -not -path "./proto/*.pb.go")

migrate-up: ## Run database migrations up (local - requires PostgreSQL running locally)
	@echo "Running migrations locally..."
	@$(HOME)/go/bin/migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down: ## Rollback database migrations (local - requires PostgreSQL running locally)
	@echo "Rolling back migrations locally..."
	@$(HOME)/go/bin/migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-create: ## Create new migration (use NAME=migration_name)
	@echo "Creating migration: $(NAME)"
	@$(HOME)/go/bin/migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

migrate-up-docker: ## Run migrations via Docker (recommended - works with dev containers)
	@echo "Running migrations via Docker..."
	@CONTAINER=$$(docker ps --filter "name=order-delivery-service" --format "{{.Names}}" | head -1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "❌ No running container found. Start with 'make dev-up' or 'make dev-debug'"; \
		exit 1; \
	fi; \
	echo "Using container: $$CONTAINER"; \
	docker exec $$CONTAINER sh -c 'migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/order_delivery_db?sslmode=disable" up'

migrate-down-docker: ## Rollback migrations via Docker
	@echo "Rolling back migrations via Docker..."
	@CONTAINER=$$(docker ps --filter "name=order-delivery-service" --format "{{.Names}}" | head -1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "❌ No running container found. Start with 'make dev-up' or 'make dev-debug'"; \
		exit 1; \
	fi; \
	echo "Using container: $$CONTAINER"; \
	docker exec $$CONTAINER sh -c 'migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/order_delivery_db?sslmode=disable" down 1'

migrate-status-docker: ## Check migration status via Docker
	@echo "Checking migration status via Docker..."
	@CONTAINER=$$(docker ps --filter "name=order-delivery-service" --format "{{.Names}}" | head -1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "❌ No running container found. Start with 'make dev-up' or 'make dev-debug'"; \
		exit 1; \
	fi; \
	echo "Using container: $$CONTAINER"; \
	docker exec $$CONTAINER sh -c 'migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/order_delivery_db?sslmode=disable" version'

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-up: ## Start services with docker-compose
	@echo "Starting Docker services..."
	@docker-compose up -d

docker-down: ## Stop services with docker-compose
	@echo "Stopping Docker services..."
	@docker-compose down

docker-logs: ## View docker-compose logs
	@docker-compose logs -f

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing migrate with PostgreSQL support..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/air-verse/air@latest
	@echo "Tools installed successfully"

install-air: ## Install Air for hot-reloading
	@echo "Installing Air..."
	@go install github.com/air-verse/air@v1.61.5
	@echo "Air installed successfully"

# Development with hot-reloading

dev: dev-up ## Start development environment with hot-reloading

dev-up: ## Start development containers with hot-reloading (no debugger, fmt.Println visible)
	@echo "Starting development environment with hot-reloading..."
	@docker-compose -f docker-compose.dev.yaml up --build

dev-down: ## Stop development containers
	@echo "Stopping development environment..."
	@docker-compose -f docker-compose.dev.yaml down

dev-logs: ## View development container logs
	@docker-compose -f docker-compose.dev.yaml logs -f service

dev-restart: ## Restart development service container
	@echo "Restarting service..."
	@docker-compose -f docker-compose.dev.yaml restart service

dev-debug: ## Start development with Delve debugger (port 2345, fmt.Println may be buffered)
	@echo "Starting development environment with Delve debugger..."
	@echo "Debugger available at localhost:2345"
	@docker-compose -f docker-compose.debug.yaml up --build

dev-debug-down: ## Stop debug containers
	@echo "Stopping debug environment..."
	@docker-compose -f docker-compose.debug.yaml down

dev-debug-logs: ## View debug container logs
	@docker-compose -f docker-compose.debug.yaml logs -f service

air: ## Run Air locally (hot-reload without Docker)
	@echo "Starting Air for hot-reloading..."
	@air -c .air.toml

# CI/CD Commands

ci-test: ## Run full CI pipeline locally (lint + migrate + test + build)
	@echo "Running CI pipeline..."
	@docker-compose -f docker-compose.ci.yaml up --build --abort-on-container-exit --exit-code-from service

ci-test-only: ## Run tests only in CI environment
	@echo "Running tests in CI environment..."
	@docker-compose -f docker-compose.ci.yaml run --rm service sh -c "go test -v -race -coverprofile=coverage.out -covermode=atomic ./..."

ci-lint: ## Run linter in CI environment
	@echo "Running linter in CI environment..."
	@docker-compose -f docker-compose.ci.yaml run --rm service sh -c "golangci-lint run --timeout=5m ./..."

ci-build: ## Build application in CI environment
	@echo "Building in CI environment..."
	@docker-compose -f docker-compose.ci.yaml run --rm service sh -c "go build -o bin/server ./cmd/server"

ci-coverage: ## Generate coverage report in CI environment
	@echo "Generating coverage report..."
	@docker-compose -f docker-compose.ci.yaml run --rm service sh -c "go test -v -race -coverprofile=coverage.out -covermode=atomic ./... && go tool cover -html=coverage.out -o coverage.html"
	@echo "Coverage report generated: coverage.html"

ci-clean: ## Stop and remove CI containers
	@echo "Cleaning CI environment..."
	@docker-compose -f docker-compose.ci.yaml down -v

act-test: ## Run tests locally with act
	@echo "Running tests with act..."
	@act -j test

act-all: ## Run all CI jobs locally with act
	@echo "Running all CI jobs with act..."
	@act

act-list: ## List all available act jobs
	@echo "Available jobs:"
	@act -l

act-lint: ## Run lint locally with act (fast)
	@echo "Running linter with act..."
	@act -j lint

act-build: ## Run build locally with act (fast)
	@echo "Running build with act..."
	@act -j build

act-clean: ## Clean act cache and containers
	@echo "Cleaning act cache..."
	@docker container prune -f
	@docker image prune -f
	@echo "Act cache cleaned"

.DEFAULT_GOAL := help
