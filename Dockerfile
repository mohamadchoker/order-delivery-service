# ==============================================================================
# PRODUCTION BUILD - Multi-stage for minimal final image (~20 MB)
# ==============================================================================

# Stage 1: Build
FROM golang:1.25-alpine AS builder

ARG VERSION=dev
ARG BUILD_DATE=unknown
ARG GIT_COMMIT=unknown

# Install build dependencies
RUN apk add --no-cache git make protobuf-dev

WORKDIR /app

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static;// binary
# CGO_ENABLED=0: Static binary (works on any Linux)
# -ldflags: Inject version info into binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.gitCommit=${GIT_COMMIT}" \
    -o /bin/order-delivery-service \
    ./cmd/server

# Stage 2: Runtime (final image)
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /bin/order-delivery-service .

# Expose ports
EXPOSE 50051 9090

# Run the service
CMD ["./order-delivery-service"]
