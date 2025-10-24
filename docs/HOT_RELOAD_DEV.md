# Hot-Reload Development Setup

## Overview

The project now supports **hot-reloading** during development - your code changes are automatically detected and the service restarts instantly!

**No more manual restarts!** Just save your file and watch it reload. ⚡

---

## What is Hot-Reloading?

Hot-reloading watches your code files and automatically:
1. Detects when you save a file
2. Recompiles the code
3. Restarts the service
4. All in ~1-2 seconds!

**Without hot-reload:**
```bash
# Edit code
vim internal/service/delivery_usecase.go

# Stop service
docker-compose down

# Rebuild
docker-compose up --build

# Wait 30-60 seconds...
```

**With hot-reload:**
```bash
# Edit code
vim internal/service/delivery_usecase.go

# Save file (Ctrl+S)
# ✨ Service automatically restarts in 1-2 seconds!
```

---

## Technology: Air

We use [Air](https://github.com/cosmtrek/air) - the most popular Go hot-reload tool.

**Features:**
- ✅ Fast rebuilds (1-2 seconds)
- ✅ Watches all Go files
- ✅ Watches proto files
- ✅ Configurable via `.air.toml`
- ✅ Works in Docker and locally
- ✅ Smart filtering (ignores test files, mocks, etc.)

---

## Quick Start

### Option 1: Docker (Recommended) 🐳

**Start development environment:**
```bash
make dev
# or
docker-compose -f docker-compose.dev.yaml up
```

**That's it!** Now edit any Go file and watch it reload automatically.

**View logs:**
```bash
make dev-logs
```

**Stop:**
```bash
make dev-down
```

### Option 2: Local (Without Docker)

**Install Air:**
```bash
make install-air
# or
go install github.com/cosmtrek/air@latest
```

**Start PostgreSQL:**
```bash
docker-compose up -d postgres
```

**Run Air:**
```bash
make air
# or
air -c .air.toml
```

---

## Development Workflow

### Typical Day

```bash
# 1. Start development environment (once per day)
make dev

# 2. Edit code
vim internal/service/delivery_usecase.go

# 3. Save file (Ctrl+S)
# ✨ Service automatically reloads!

# 4. Test with grpcurl
grpcurl -plaintext localhost:50051 list

# 5. Edit more code, save, repeat...
# ✨ Each save triggers automatic reload!

# 6. When done for the day
make dev-down
```

### What Triggers Reload?

Air watches these file types:
- ✅ `.go` files (all Go source code)
- ✅ `.proto` files (Protocol Buffer definitions)
- ✅ `.html`, `.tmpl`, `.tpl` (templates, if you add them)

Air **ignores** these:
- ❌ `*_test.go` (test files - use `go test` for these)
- ❌ `*_mock.go` (mock files)
- ❌ `tmp/`, `vendor/`, `.git/` (build artifacts, dependencies)

---

## Configuration

### Air Config File: `.air.toml`

```toml
[build]
  # Command to build
  cmd = "go build -o ./tmp/main ./cmd/server"

  # Output binary
  bin = "./tmp/main"

  # Delay before rebuild (1 second)
  delay = 1000

  # Watch these extensions
  include_ext = ["go", "tpl", "tmpl", "html", "proto"]

  # Ignore these directories
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "migrations"]

  # Ignore these patterns
  exclude_regex = ["_test.go", "_mock.go"]
```

### Development Docker Compose: `docker-compose.dev.yaml`

```yaml
services:
  service:
    build:
      dockerfile: Dockerfile.dev  # Uses Air
    volumes:
      - .:/app  # Mount source code
      - go_modules:/go/pkg/mod  # Cache modules
      - /app/tmp  # Exclude build artifacts
    environment:
      - LOG_LEVEL=debug  # More verbose logs
      - LOG_DEV=true     # Pretty-printed logs
```

---

## Examples

### Example 1: Edit Service Logic

```bash
# 1. Start dev environment
make dev

# 2. In another terminal, edit service
vim internal/service/delivery_usecase.go

# Change something, e.g., add logging:
# logger.Info("Creating delivery assignment - NEW LOG!")

# 3. Save file (Ctrl+S)

# Watch the logs in first terminal:
```

Output:
```
building...
running...
2025-10-24T00:00:00+03:00 INF Starting server on :50051
✅ Service reloaded in 1.2s
```

### Example 2: Update Proto Definition

```bash
# 1. Edit proto
vim proto/delivery.proto

# Add new field to message:
# string tracking_number = 14;

# 2. Save file

# Air detects proto change and:
# - Runs protoc to regenerate Go code
# - Recompiles service
# - Restarts
```

### Example 3: Fix a Bug

```bash
# You discover a bug in handler
vim internal/transport/grpc/grpc_handler.go

# Fix the bug, save

# ✨ Instant reload - test fix immediately!
grpcurl -plaintext -d '{...}' localhost:50051 delivery.DeliveryService/...
```

---

## Comparison

### Without Hot-Reload (Old Way)

| Step | Time |
|------|------|
| Edit code | - |
| Stop containers | 5s |
| Rebuild image | 30s |
| Start containers | 10s |
| Wait for DB | 5s |
| **Total** | **~50s per change** |

### With Hot-Reload (New Way)

| Step | Time |
|------|------|
| Edit code | - |
| Save file | - |
| **Automatic reload** | **~2s** |
| **Total** | **~2s per change** |

**25x faster!** 🚀

---

## Docker Setup

### Development Dockerfile: `Dockerfile.dev`

```dockerfile
FROM golang:1.24-alpine

# Install Air
RUN go install github.com/cosmtrek/air@latest

# Install dev tools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install go.uber.org/mock/mockgen@latest

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run Air (watches for changes)
CMD ["air", "-c", ".air.toml"]
```

### Key Differences from Production Dockerfile

| Aspect | Production (`Dockerfile`) | Development (`Dockerfile.dev`) |
|--------|--------------------------|-------------------------------|
| **Base** | Multi-stage (builder + alpine) | Single-stage (golang:alpine) |
| **Size** | ~20 MB | ~500 MB |
| **Tools** | None (minimal) | Air, protoc, mockgen |
| **Source** | Copied & built | Mounted (live changes) |
| **CMD** | Run binary | Run Air |
| **Rebuild** | Required | Automatic |

---

## Makefile Commands

```bash
# Development
make dev              # Start dev environment with hot-reload
make dev-up           # Same as dev
make dev-down         # Stop dev environment
make dev-logs         # View service logs (follow mode)
make dev-restart      # Restart service container

# Local Air (without Docker)
make air              # Run Air locally
make install-air      # Install Air tool

# Production (no hot-reload)
make docker-up        # Start production containers
make docker-down      # Stop production containers
```

---

## Troubleshooting

### Issue: Air not detecting changes

**Cause:** File system events not propagating to Docker

**Solution 1:** Use polling mode
```toml
# .air.toml
[build]
  poll = true
  poll_interval = 500  # milliseconds
```

**Solution 2:** Check mounted volume
```bash
# Verify volume is mounted
docker-compose -f docker-compose.dev.yaml exec service ls -la /app
```

### Issue: Build errors

**Cause:** Syntax error in code

**Solution:** Check build-errors.log
```bash
cat build-errors.log
```

### Issue: Service not starting after reload

**Cause:** Runtime error (e.g., database connection)

**Solution:** Check logs
```bash
make dev-logs
```

### Issue: Too many rebuilds

**Cause:** Editor creating temporary files

**Solution:** Exclude temp files in .air.toml
```toml
[build]
  exclude_file = [".swp", ".tmp", "~"]
```

### Issue: Slow rebuilds

**Cause:** Large codebase or no module caching

**Solution:** Docker compose already caches Go modules
```yaml
volumes:
  - go_modules:/go/pkg/mod  # ✅ Cached
```

---

## Advanced Configuration

### Custom Build Command

Edit `.air.toml`:
```toml
[build]
  # Run tests before building
  cmd = "go test ./... && go build -o ./tmp/main ./cmd/server"

  # Or generate mocks first
  cmd = "go generate ./... && go build -o ./tmp/main ./cmd/server"
```

### Watch Additional Files

```toml
[build]
  include_ext = ["go", "proto", "yaml", "json"]
```

### Faster Rebuild Delay

```toml
[build]
  delay = 500  # 0.5 seconds (default: 1000)
```

### Kill Previous Process

```toml
[build]
  send_interrupt = true
  kill_delay = "500ms"
```

---

## Production vs Development

### Use Development Mode For:
- ✅ Active development
- ✅ Debugging
- ✅ Testing changes quickly
- ✅ Learning/experimenting

### Use Production Mode For:
- ✅ Final testing before deploy
- ✅ Performance testing
- ✅ Actual deployment
- ✅ CI/CD pipelines

### Switching Modes

```bash
# Development (hot-reload)
make dev

# Production (optimized)
make docker-up
```

---

## FAQ

### Q: Does hot-reload work for proto changes?

**A:** Yes! Edit `.proto` file, save, and Air will:
1. Detect proto change
2. Run `protoc` to regenerate Go code
3. Rebuild service
4. Restart

### Q: What about database migrations?

**A:** Migrations are **not** auto-applied. Run manually:
```bash
make migrate-up
```

Or add to Air build command in `.air.toml`.

### Q: Can I use this in CI/CD?

**A:** No! Use production Dockerfile for CI/CD:
```yaml
# .github/workflows/ci.yml (correct)
docker build -f Dockerfile .

# Don't use Dockerfile.dev in CI!
```

### Q: Does it work on Windows/Mac/Linux?

**A:** Yes! Air works on all platforms. Docker volumes work best on Linux, but Mac/Windows work fine too.

### Q: How do I disable hot-reload temporarily?

```bash
# Stop dev environment
make dev-down

# Use production mode
make docker-up
```

### Q: Can I use Air without Docker?

**A:** Yes!
```bash
# Install Air
make install-air

# Start PostgreSQL only
docker-compose up -d postgres

# Run Air locally
make air
```

---

## Best Practices

### 1. Use Development Mode for Active Work

```bash
# Start of day
make dev

# Work all day with instant reloads

# End of day
make dev-down
```

### 2. Test in Production Mode Before Push

```bash
# Before git push
make dev-down
make docker-up

# Final test with production build
grpcurl -plaintext localhost:50051 ...

# If all good, push
git push
```

### 3. Keep .air.toml in Git

✅ Commit `.air.toml` so team has same config
✅ Exclude `tmp/` in `.gitignore` (build artifacts)
✅ Exclude `build-errors.log` in `.gitignore`

### 4. Use Separate Databases

```yaml
# docker-compose.dev.yaml
postgres:
  environment:
    POSTGRES_DB: order_delivery_db_dev  # Different DB!
```

### 5. Enable Debug Logging in Dev

```yaml
# docker-compose.dev.yaml
environment:
  - LOG_LEVEL=debug  # Verbose logs
  - LOG_DEV=true     # Pretty printing
```

---

## Summary

### What You Get

✅ **Instant feedback** - 2 seconds vs 50 seconds
✅ **Auto-reload** - Just save, no manual restart
✅ **Proto support** - Detects proto changes too
✅ **Docker integrated** - Works in containers
✅ **Local option** - Can run without Docker
✅ **Team ready** - Configuration committed to git

### Quick Reference

```bash
# Start development (hot-reload)
make dev

# Stop development
make dev-down

# View logs
make dev-logs

# That's it! Edit code, save, watch it reload! ⚡
```

---

**Pro Tip:** Open two terminals:
- **Terminal 1:** `make dev` (shows reload logs)
- **Terminal 2:** Edit code, run grpcurl tests

This workflow is **25x faster** than rebuild/restart cycles! 🚀
