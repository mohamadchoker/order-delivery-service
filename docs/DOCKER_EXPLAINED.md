# Docker & Docker Compose Explained

This guide explains all the Docker concepts used in this project.

---

## Table of Contents
1. [Dockerfile Stages (Multi-stage Builds)](#dockerfile-stages)
2. [Docker Compose Services](#docker-compose-services)
3. [Why We Have Two Dockerfiles](#why-two-dockerfiles)
4. [Why We Have Two Docker Compose Files](#why-two-docker-compose-files)
5. [Volumes Explained](#volumes-explained)
6. [Environment Variables](#environment-variables)

---

## Dockerfile Stages

### What is `AS builder`?

In `Dockerfile`, you'll see:
```dockerfile
FROM golang:1.25-alpine AS builder
```

**What does `AS builder` mean?**
- This is a **named stage** in a multi-stage build
- `builder` is just a name we chose (you could call it anything: `build-stage`, `compile`, etc.)
- It's used to reference this stage later

### Multi-Stage Build Explained

**Dockerfile (Production)**:
```dockerfile
# Stage 1: Build the application
FROM golang:1.25-alpine AS builder    ← Named stage "builder"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/order-delivery-service ./cmd/server

# Stage 2: Create minimal runtime image
FROM alpine:latest                     ← New stage (unnamed)
WORKDIR /app
COPY --from=builder /bin/order-delivery-service .  ← Copy from "builder" stage
CMD ["./order-delivery-service"]
```

**Why use multi-stage builds?**

| Without Multi-stage | With Multi-stage |
|---------------------|------------------|
| **Size**: ~500 MB (includes Go compiler, tools) | **Size**: ~20 MB (only binary + Alpine) |
| Includes: Go compiler, git, make, source code | Includes: Only compiled binary |
| Security: More attack surface | Security: Minimal attack surface |
| Use case: Development | Use case: Production |

**Stages Breakdown**:

1. **Stage 1 (builder)**:
   - Uses `golang:1.25-alpine` (has Go compiler)
   - Downloads dependencies
   - Compiles the code
   - Produces binary: `/bin/order-delivery-service`

2. **Stage 2 (final)**:
   - Uses `alpine:latest` (tiny Linux, ~5 MB)
   - Only copies the **compiled binary** from stage 1
   - Doesn't include Go compiler, source code, or build tools
   - Result: Tiny, secure production image

**The `AS` keyword**:
- `AS builder` → Names the stage so you can reference it later
- `COPY --from=builder` → Copy files from the "builder" stage

### Dockerfile.dev (No Multi-stage)

```dockerfile
FROM golang:1.24-alpine    ← Single stage, no "AS"
# Install development tools
RUN go install github.com/air-verse/air@v1.61.5
WORKDIR /app
COPY . .
CMD ["/entrypoint-dev.sh"]
```

**Why no multi-stage?**
- Development needs the Go compiler (for hot-reload)
- We want source code mounted (not compiled)
- Air recompiles on file changes
- Size doesn't matter in development

---

## Docker Compose Services

### What is a "Service"?

In Docker Compose, a **service** is a **container** you want to run.

**docker-compose.yaml**:
```yaml
services:
  postgres:       ← Service 1: PostgreSQL database
    image: postgres:14-alpine
    ...

  service:        ← Service 2: Our Go application
    build:
      dockerfile: Dockerfile
    ...
```

**Why is our app called `service`?**
- It's just a name! We could call it anything:
  - `app`
  - `api`
  - `backend`
  - `order-delivery-service`
  - `grpc-server`

**Common naming conventions**:

| Name | When to Use | Example |
|------|-------------|---------|
| `service` | Generic, single service | This project |
| `api` | REST/GraphQL API | `api`, `rest-api` |
| `app` | Main application | `app`, `backend` |
| `{service-name}` | Descriptive | `order-service`, `user-service` |
| `{service-name}-api` | Multiple services | `order-api`, `delivery-api` |

**We use `service` because**:
- Simple, clear
- Single application in this project
- Easy to reference: `docker-compose logs service`

### Service vs Container

| Concept | Meaning | Example |
|---------|---------|---------|
| **Service** | Definition in docker-compose.yaml | `service:` block |
| **Container** | Running instance of a service | `order-delivery-service-dev` |

**Example**:
```yaml
services:
  service:                                  ← Service name (in compose file)
    container_name: order-delivery-service  ← Container name (in Docker)
```

Commands:
```bash
# Using service name (docker-compose commands)
docker-compose logs service
docker-compose restart service

# Using container name (docker commands)
docker logs order-delivery-service
docker restart order-delivery-service
```

---

## Why Two Dockerfiles?

### Dockerfile (Production)

**Purpose**: Optimized for production deployment

**Characteristics**:
- ✅ Multi-stage build (small final image: ~20 MB)
- ✅ No development tools
- ✅ No source code in final image
- ✅ Secure (minimal attack surface)
- ✅ Fast startup
- ❌ No hot-reload (need to rebuild to see changes)

**When to use**:
- Production deployment
- CI/CD pipelines
- Performance testing
- Docker Hub releases

**Build**:
```bash
docker build -t order-delivery-service -f Dockerfile .
```

### Dockerfile.dev (Development)

**Purpose**: Optimized for development workflow

**Characteristics**:
- ✅ Hot-reload with Air (see changes instantly)
- ✅ Development tools installed (protoc, mockgen, migrate)
- ✅ Source code mounted as volume
- ✅ Debugging capabilities
- ❌ Large image (~500 MB)
- ❌ Not optimized for production

**When to use**:
- Local development
- Debugging
- Testing new features
- Learning/experimentation

**Build**:
```bash
docker build -t order-delivery-service:dev -f Dockerfile.dev .
```

### Comparison

| Aspect | Dockerfile (Prod) | Dockerfile.dev |
|--------|-------------------|----------------|
| **Size** | ~20 MB | ~500 MB |
| **Stages** | Multi-stage | Single-stage |
| **Go version** | 1.25 | 1.24 |
| **Hot-reload** | ❌ No | ✅ Yes (Air) |
| **Source code** | Not included | Mounted as volume |
| **Tools** | None | protoc, mockgen, migrate, Air |
| **Build time** | ~2 min (first), ~30s (cached) | ~3 min (first), ~10s (cached) |
| **Startup** | Instant | ~2s (Air) |
| **Use case** | Production | Development |

---

## Why Two Docker Compose Files?

### docker-compose.yaml (Production)

**Purpose**: Production-like local environment

```yaml
services:
  service:
    build:
      dockerfile: Dockerfile      ← Uses production Dockerfile
    restart: unless-stopped       ← Auto-restart on failure
    volumes: []                   ← No volumes (self-contained)
```

**Characteristics**:
- Uses optimized production Dockerfile
- Auto-restarts on failure
- No source code mounting
- Mirrors production setup
- Good for final testing before deploy

**Use case**:
```bash
# Test production build locally
docker-compose up

# Verify everything works like production
make docker-up
```

### docker-compose.dev.yaml (Development)

**Purpose**: Development environment with hot-reload

```yaml
services:
  service:
    build:
      dockerfile: Dockerfile.dev  ← Uses development Dockerfile
    restart: "no"                 ← Don't restart (we want to see crashes)
    volumes:
      - .:/app                    ← Mount source code
      - go_modules:/go/pkg/mod    ← Cache Go modules
      - /app/tmp                  ← Exclude build artifacts
    environment:
      - LOG_LEVEL=debug           ← Verbose logging
      - LOG_DEV=true              ← Pretty-printed logs
```

**Characteristics**:
- Uses development Dockerfile with Air
- Source code mounted (edit files, see changes)
- No auto-restart (want to see errors)
- Debug logging enabled
- Go modules cached (faster rebuilds)

**Use case**:
```bash
# Start development environment
docker-compose -f docker-compose.dev.yaml up

# Or use the shortcut
make dev
```

### Comparison

| Aspect | docker-compose.yaml | docker-compose.dev.yaml |
|--------|---------------------|-------------------------|
| **Dockerfile** | Dockerfile (prod) | Dockerfile.dev |
| **Hot-reload** | ❌ No | ✅ Yes |
| **Volumes** | None | Source code mounted |
| **Restart policy** | `unless-stopped` | `no` |
| **Logging** | INFO level | DEBUG level |
| **Container name** | `order-delivery-service` | `order-delivery-service-dev` |
| **Volume name** | `postgres_data` | `postgres_data_dev` |
| **Use case** | Production testing | Active development |

---

## Volumes Explained

### What are Volumes?

Volumes are persistent storage for containers.

### Types of Volumes

#### 1. Named Volumes (Managed by Docker)

**docker-compose.yaml**:
```yaml
services:
  postgres:
    volumes:
      - postgres_data:/var/lib/postgresql/data  ← Named volume

volumes:
  postgres_data:  ← Define the volume
```

**What this does**:
- Docker creates and manages the volume
- Data persists even when container is deleted
- Located in Docker's storage area (e.g., `/var/lib/docker/volumes/`)
- Shared between container restarts

**Use case**: Database data, uploaded files, logs

#### 2. Bind Mounts (Host → Container)

**docker-compose.dev.yaml**:
```yaml
services:
  service:
    volumes:
      - .:/app  ← Bind mount: current directory → /app in container
```

**What this does**:
- Mounts a directory from your host machine into the container
- `.` = current directory (where docker-compose.yaml is)
- `/app` = path inside the container
- Changes on host are immediately visible in container
- Changes in container are immediately visible on host

**Use case**: Development (edit files locally, see changes in container)

#### 3. Anonymous Volumes (Temporary)

**docker-compose.dev.yaml**:
```yaml
services:
  service:
    volumes:
      - /app/tmp  ← Anonymous volume: exclude this directory
```

**What this does**:
- Creates a volume NOT linked to host
- `/app/tmp` is excluded from bind mount
- Air's build artifacts stay in container (not synced to host)
- Faster (no syncing binary files to host)

**Use case**: Exclude directories from bind mounts

### Volume Examples in This Project

**docker-compose.dev.yaml**:
```yaml
volumes:
  # Bind mount: host source code → container /app
  - .:/app

  # Named volume: cache Go modules (faster rebuilds)
  - go_modules:/go/pkg/mod

  # Anonymous volume: exclude Air build artifacts
  - /app/tmp
```

**Why these volumes?**

| Volume | Type | Purpose |
|--------|------|---------|
| `.:/app` | Bind mount | Edit code locally, see changes instantly |
| `go_modules:/go/pkg/mod` | Named | Cache downloaded Go modules (10x faster) |
| `/app/tmp` | Anonymous | Keep Air's binaries in container (faster) |

**What happens**:
1. You edit `internal/service/delivery_usecase.go` on your Mac
2. Change is immediately visible in container at `/app/internal/service/delivery_usecase.go`
3. Air detects change and rebuilds
4. Air's temporary binary is saved to `/app/tmp` (inside container only)
5. Go modules stay cached in `go_modules` volume (don't re-download)

---

## Environment Variables

### Why Environment Variables?

Instead of hardcoding values in code:
```go
// ❌ Bad
db := "postgres://postgres:postgres@localhost:5432/mydb"
```

Use environment variables:
```go
// ✅ Good
dbHost := os.Getenv("DB_HOST")
dbPort := os.Getenv("DB_PORT")
```

### Setting Environment Variables

**docker-compose.yaml**:
```yaml
services:
  service:
    environment:
      - DB_HOST=postgres      ← Set environment variable
      - DB_PORT=5432
      - LOG_LEVEL=info
```

**Inside container**:
```bash
# These environment variables are available:
echo $DB_HOST     # → postgres
echo $DB_PORT     # → 5432
echo $LOG_LEVEL   # → info
```

**In Go code** (`cmd/server/main.go`):
```go
dbHost := os.Getenv("DB_HOST")        // → "postgres"
dbPort := os.Getenv("DB_PORT")        // → "5432"
logLevel := os.Getenv("LOG_LEVEL")    // → "info"
```

### Different Values for Dev vs Prod

**docker-compose.yaml (Production)**:
```yaml
environment:
  - LOG_LEVEL=info       ← Less verbose
  - LOG_DEV=false        ← JSON logs
```

**docker-compose.dev.yaml (Development)**:
```yaml
environment:
  - LOG_LEVEL=debug      ← More verbose
  - LOG_DEV=true         ← Pretty-printed logs
```

---

## Quick Reference

### Commands

```bash
# Production
docker-compose up                          # Start production setup
docker-compose down                        # Stop production setup
docker-compose logs service                # View production logs

# Development
docker-compose -f docker-compose.dev.yaml up     # Start dev setup
docker-compose -f docker-compose.dev.yaml down   # Stop dev setup
docker-compose -f docker-compose.dev.yaml logs service  # View dev logs

# Shortcuts (via Makefile)
make docker-up        # Production
make docker-down      # Stop production
make dev              # Development
make dev-down         # Stop development
make dev-logs         # View dev logs
```

### File Overview

```
project/
├── Dockerfile                  # Production (multi-stage, optimized)
├── Dockerfile.dev              # Development (Air, tools)
├── docker-compose.yaml         # Production setup
├── docker-compose.dev.yaml     # Development setup (hot-reload)
└── entrypoint-dev.sh           # Development startup script
```

### When to Use What

| Task | Command |
|------|---------|
| **Daily development** | `make dev` |
| **Test prod build** | `make docker-up` |
| **View dev logs** | `make dev-logs` |
| **Restart dev service** | `docker-compose -f docker-compose.dev.yaml restart service` |
| **Clean everything** | `make dev-down && docker system prune -a` |

---

## Common Questions

### Q: Why is the service called "service"?
**A:** It's just a name we chose. You can rename it to `app`, `api`, `backend`, etc. in docker-compose.yaml

### Q: What is `AS builder`?
**A:** It names a stage in a multi-stage Dockerfile so you can reference it later with `COPY --from=builder`

### Q: Why two Dockerfiles?
**A:**
- `Dockerfile` = Production (small, optimized)
- `Dockerfile.dev` = Development (hot-reload, tools)

### Q: Why mount source code in dev?
**A:** So you can edit files on your Mac and see changes instantly in the container (hot-reload)

### Q: Why cache Go modules?
**A:** Downloading dependencies is slow (~30s). Caching makes it instant on subsequent builds.

### Q: Why exclude `/app/tmp`?
**A:** Air's binary files are large. No need to sync them to your host machine. Keep them in container only.

### Q: Can I use one Dockerfile for both?
**A:** Yes, but not recommended. Dev and prod have different needs. Separate files = clearer separation.

---

This explains all the Docker concepts in this project! 🐳
