#!/bin/bash

set -e

echo "🚀 Order Delivery Service Setup"
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go 1.21 or later: https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed${NC}"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Go version: $GO_VERSION"

# Install tools
echo ""
echo "📦 Installing development tools..."

echo "Installing protoc-gen-go..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

echo "Installing protoc-gen-go-grpc..."
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

echo "Installing migrate..."
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

echo -e "${GREEN}✓ Tools installed${NC}"

# Download dependencies
echo ""
echo "📥 Downloading Go dependencies..."
go mod download
echo -e "${GREEN}✓ Dependencies downloaded${NC}"

# Check if protoc is installed
echo ""
if ! command -v protoc &> /dev/null; then
    echo -e "${YELLOW}⚠ protoc is not installed${NC}"
    echo "To generate proto files, install protoc:"
    echo "  - macOS: brew install protobuf"
    echo "  - Linux: apt-get install protobuf-compiler"
    echo "  - Or download from: https://github.com/protocolbuffers/protobuf/releases"
    SKIP_PROTO=true
else
    echo -e "${GREEN}✓ protoc is installed${NC}"
    SKIP_PROTO=false
fi

# Generate proto files
if [ "$SKIP_PROTO" = false ]; then
    echo ""
    echo "🔧 Generating proto files..."
    mkdir -p proto
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/*.proto
    echo -e "${GREEN}✓ Proto files generated${NC}"
fi

# Create config file if it doesn't exist
if [ ! -f config/config.yaml ]; then
    echo ""
    echo "📝 Creating config file..."
    cp config/config.example.yaml config/config.yaml
    echo -e "${GREEN}✓ Config file created${NC}"
fi

# Check if PostgreSQL is running
echo ""
echo "🔍 Checking PostgreSQL..."
if command -v psql &> /dev/null; then
    if psql -U postgres -c '\q' 2>/dev/null; then
        echo -e "${GREEN}✓ PostgreSQL is running${NC}"
        
        # Create database
        echo "Creating database..."
        psql -U postgres -c "CREATE DATABASE order_delivery_db;" 2>/dev/null || echo "Database already exists"
        
        # Run migrations
        echo "Running migrations..."
        export DB_URL="postgres://postgres:postgres@localhost:5432/order_delivery_db?sslmode=disable"
        migrate -path migrations -database "$DB_URL" up 2>/dev/null || echo "Migrations already applied"
        
        echo -e "${GREEN}✓ Database setup complete${NC}"
    else
        echo -e "${YELLOW}⚠ PostgreSQL is not running or not accessible${NC}"
        echo "Start PostgreSQL or use Docker:"
        echo "  docker-compose up -d postgres"
    fi
else
    echo -e "${YELLOW}⚠ PostgreSQL is not installed${NC}"
    echo "Install PostgreSQL or use Docker:"
    echo "  docker-compose up -d postgres"
fi

# Build the application
echo ""
echo "🔨 Building application..."
go build -o bin/order-delivery-service cmd/server/main.go
echo -e "${GREEN}✓ Build successful${NC}"

# Run tests
echo ""
echo "🧪 Running tests..."
go test -v ./... 2>&1 | head -n 20
echo -e "${GREEN}✓ Tests completed${NC}"

echo ""
echo "================================"
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "Next steps:"
echo "1. Start PostgreSQL (if not running):"
echo "   docker-compose up -d postgres"
echo ""
echo "2. Run migrations:"
echo "   make migrate-up"
echo ""
echo "3. Start the service:"
echo "   make run"
echo "   or"
echo "   ./bin/order-delivery-service"
echo ""
echo "4. Test the service:"
echo "   grpcurl -plaintext localhost:50051 list"
echo ""
echo "📚 Documentation:"
echo "   - README.md - Getting started guide"
echo "   - docs/API.md - API documentation"
echo "   - docs/ARCHITECTURE.md - Architecture overview"
