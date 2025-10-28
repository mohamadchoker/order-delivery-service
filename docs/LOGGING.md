# Logging Configuration Guide

Complete guide to configuring structured logging with Zap in the Order Delivery Service.

---

## 🎯 Quick Reference

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error`, `fatal` |
| `LOG_DEV` | `false` | Development mode (human-readable, colored output) |
| `LOG_STACKTRACE` | `false` | Enable stack traces in error logs |

### Common Configurations

**Production (Clean, JSON logs):**
```bash
LOG_LEVEL=info
LOG_DEV=false
LOG_STACKTRACE=false
```

**Development (Readable, colored):**
```bash
LOG_LEVEL=debug
LOG_DEV=true
LOG_STACKTRACE=false
```

**Debugging (With stack traces):**
```bash
LOG_LEVEL=debug
LOG_DEV=true
LOG_STACKTRACE=true
```

---

## 📊 Log Levels

Logs are filtered by level. Only messages at or above the configured level will appear.

| Level | When to use | Example |
|-------|-------------|---------|
| `debug` | Development, troubleshooting | Request details, variable values |
| `info` | Normal operations | Service started, request completed |
| `warn` | Potential issues | Deprecated feature used, retry attempted |
| `error` | Failures | Database connection failed, validation error |
| `fatal` | Critical failures | Service cannot start, config missing |

**Example:**
```bash
# With LOG_LEVEL=info, you'll see: info, warn, error, fatal
# You WON'T see: debug
```

---

## 🎨 Development vs Production Mode

### Production Mode (`LOG_DEV=false`)

**Output:** JSON format (machine-readable, for log aggregation)

```json
{
  "level": "error",
  "ts": 1698345678.123456,
  "caller": "grpc/grpc_handler.go:125",
  "msg": "Failed to create delivery assignment",
  "error": "invalid input: order_id is required",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "/delivery.DeliveryService/CreateDeliveryAssignment"
}
```

**Best for:**
- Production environments
- Log aggregation tools (ELK, Splunk, Datadog)
- Automated parsing and analysis

---

### Development Mode (`LOG_DEV=true`)

**Output:** Human-readable, colored console output

```
2024-10-24T15:30:45.123Z  ERROR  grpc/grpc_handler.go:125  Failed to create delivery assignment
        error: invalid input: order_id is required
        request_id: 550e8400-e29b-41d4-a716-446655440000
        method: /delivery.DeliveryService/CreateDeliveryAssignment
```

**Best for:**
- Local development
- Reading logs in terminal
- Quick debugging

---

## 🔍 Stack Traces

Stack traces show the call chain when errors occur. They help identify where errors originated.

### Stack Trace Behavior

| `LOG_STACKTRACE` | `LOG_DEV` | Stack traces shown? |
|-----------------|-----------|---------------------|
| `false` | `false` | ❌ No |
| `false` | `true` | ❌ No |
| `true` | `false` | ✅ Yes (on errors) |
| `true` | `true` | ✅ Yes (on errors) |

### Example Stack Trace

**Without stack traces (`LOG_STACKTRACE=false`):**
```
2024-10-24T15:30:45.123Z  ERROR  grpc/grpc_handler.go:125  Failed to create delivery
        error: database connection failed
```

**With stack traces (`LOG_STACKTRACE=true`):**
```
2024-10-24T15:30:45.123Z  ERROR  grpc/grpc_handler.go:125  Failed to create delivery
        error: database connection failed
 github.com/mohamadchoker/order-delivery-service/internal/transport/grpc.(*Handler).CreateDeliveryAssignment
        /app/internal/transport/grpc/grpc_handler.go:125
 github.com/mohamadchoker/order-delivery-service/internal/service.(*deliveryUseCase).CreateDeliveryAssignment
        /app/internal/service/delivery_usecase.go:45
 github.com/mohamadchoker/order-delivery-service/internal/repository/postgres.(*postgresRepository).Create
        /app/internal/repository/postgres/repository.go:67
```

### When to Enable Stack Traces

**Enable (`true`):**
- 🐛 Debugging complex issues
- 🔬 Understanding error propagation
- 📊 Analyzing call chains

**Disable (`false`):**
- 🚀 Production (cleaner logs, better performance)
- 📝 Normal development (less noise)
- 🧪 Testing (easier to read test output)

---

## 🛠️ Configuration Examples

### Docker Compose

**docker-compose.dev.yaml** (Clean development logs):
```yaml
environment:
  - LOG_LEVEL=debug
  - LOG_DEV=true
  - LOG_STACKTRACE=false  # Clean logs, no stack traces
```

**docker-compose.debug.yaml** (Full debugging):
```yaml
environment:
  - LOG_LEVEL=debug
  - LOG_DEV=true
  - LOG_STACKTRACE=true  # Enable when investigating errors
```

**docker-compose.yaml** (Production):
```yaml
environment:
  - LOG_LEVEL=info
  - LOG_DEV=false
  - LOG_STACKTRACE=false
```

---

### Local Development (.env)

For local development outside Docker:

**.env:**
```bash
# Normal development
LOG_LEVEL=debug
LOG_DEV=true
LOG_STACKTRACE=false

# When debugging specific errors, change to:
# LOG_STACKTRACE=true
```

---

## 📋 Log Fields

All logs include these fields automatically:

| Field | Description | Example |
|-------|-------------|---------|
| `level` | Log severity | `info`, `error` |
| `ts` | Timestamp | `2024-10-24T15:30:45.123Z` |
| `caller` | Source location | `grpc/grpc_handler.go:125` |
| `msg` | Log message | `Failed to create delivery` |
| `request_id` | Request trace ID | `550e8400-...` (from X-Request-ID header) |

Additional fields depend on the context (method, error, user_id, etc.).

---

## 🎯 Usage in Code

### Logging with Context

```go
import (
    "go.uber.org/zap"
    "github.com/mohamadchoker/order-delivery-service/pkg/middleware"
)

func (h *Handler) CreateDelivery(ctx context.Context, req *pb.CreateRequest) (*pb.Delivery, error) {
    // Get request ID from context (automatically added by middleware)
    requestID := middleware.GetRequestID(ctx)

    h.logger.Info("Creating delivery assignment",
        zap.String("request_id", requestID),
        zap.String("order_id", req.OrderId),
    )

    // ... business logic ...

    if err != nil {
        h.logger.Error("Failed to create delivery",
            zap.String("request_id", requestID),
            zap.String("order_id", req.OrderId),
            zap.Error(err),
        )
        return nil, err
    }

    return delivery, nil
}
```

### Log Levels in Code

```go
// DEBUG - Detailed information for troubleshooting
logger.Debug("Processing delivery assignment",
    zap.String("delivery_id", id),
    zap.Any("input", input),
)

// INFO - Normal operational messages
logger.Info("Delivery created successfully",
    zap.String("delivery_id", delivery.ID),
    zap.String("order_id", delivery.OrderID),
)

// WARN - Warning but not an error
logger.Warn("Driver not available, retrying",
    zap.String("driver_id", driverID),
    zap.Int("retry_count", retries),
)

// ERROR - Error occurred but service continues
logger.Error("Failed to update delivery status",
    zap.String("delivery_id", id),
    zap.Error(err),
)

// FATAL - Critical error, service cannot continue
logger.Fatal("Failed to connect to database",
    zap.Error(err),
)
```

---

## 🔧 Troubleshooting

### Logs are too noisy

**Problem:** Too many debug logs in production.

**Solution:** Set `LOG_LEVEL=info` or `LOG_LEVEL=warn`

```bash
# In docker-compose.yaml or .env
LOG_LEVEL=info
```

---

### Can't read logs in terminal

**Problem:** JSON logs are hard to read during development.

**Solution:** Enable development mode

```bash
LOG_DEV=true
```

---

### Need to debug error propagation

**Problem:** Can't see where error originated.

**Solution:** Enable stack traces temporarily

```bash
# In docker-compose.dev.yaml
LOG_STACKTRACE=true
```

Then restart:
```bash
make dev-down
make dev-up
```

**Remember to disable after debugging:**
```bash
LOG_STACKTRACE=false
```

---

### Logs missing caller info (file:line)

**Problem:** Can't see which file/line logged the message.

**Solution:** Caller info is always included. Check your log viewer or format.

In production JSON:
```json
{
  "caller": "grpc/grpc_handler.go:125",
  ...
}
```

In development:
```
grpc/grpc_handler.go:125  Failed to create delivery
```

---

## 📊 Monitoring and Analysis

### Finding Logs by Request ID

All logs for a single request share the same `request_id`:

**Production (JSON logs):**
```bash
# Using jq
cat logs.json | jq 'select(.request_id == "550e8400-...")'

# Using grep
grep "550e8400-..." logs.json
```

**Development (console logs):**
```bash
grep "request_id: 550e8400-..." logs.txt
```

---

### Common Log Queries

**Find all errors:**
```bash
# JSON logs
cat logs.json | jq 'select(.level == "error")'

# Console logs
grep "ERROR" logs.txt
```

**Find slow requests (using Prometheus metrics):**
```bash
# Check metrics endpoint
curl http://localhost:9090/metrics | grep grpc_request_duration
```

**Track delivery lifecycle:**
```bash
# Find all logs for a specific order
grep "order_id.*ORDER-123" logs.txt
```

---

## 🎓 Best Practices

1. **Use appropriate log levels**
   - Don't log everything as ERROR
   - Use DEBUG for development-only information

2. **Include context**
   - Always include `request_id` for tracing
   - Add relevant IDs (order_id, driver_id, etc.)

3. **Disable stack traces in production**
   - Cleaner logs
   - Better performance
   - Less noise

4. **Use structured logging**
   - Use `zap.String()`, `zap.Int()`, etc.
   - Don't use string formatting in log messages

   **Good:**
   ```go
   logger.Info("Created delivery", zap.String("id", id))
   ```

   **Bad:**
   ```go
   logger.Info(fmt.Sprintf("Created delivery: %s", id))
   ```

5. **Keep log messages concise**
   - Log message should be a brief summary
   - Details go in structured fields

---

## 📚 Related Documentation

- [Configuration Guide](../config/config.example.yaml)
- [Middleware Documentation](../pkg/middleware/)
- [Metrics Documentation](./MONITORING.md)
- [Request ID Tracing](../pkg/middleware/request_id.go)

---

## 🔗 External Resources

- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Structured Logging Best Practices](https://www.datadoghq.com/blog/structured-logging/)
- [12-Factor App Logs](https://12factor.net/logs)
