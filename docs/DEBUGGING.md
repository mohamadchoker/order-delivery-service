# Debugging Guide

## Development Modes

The project supports two development modes with SQL query logging enabled by default:

### 1. Normal Development Mode (Recommended for most work)
```bash
make dev-up
# or
make dev
```

**Features:**
- ✅ Hot-reloading with Air
- ✅ `fmt.Println` and `fmt.Printf` output visible in logs
- ✅ Fast rebuild times
- ✅ All stdout/stderr visible immediately
- ❌ No debugger attached

**Use this when:**
- You're actively developing features
- You need to see `fmt.Println` debug output
- You want fast feedback loops
- You don't need to set breakpoints

**View logs:**
```bash
make dev-logs
# or
docker logs -f order-delivery-service-dev
```

### 2. Debug Mode with Delve (For breakpoint debugging)
```bash
make dev-debug
```

**Features:**
- ✅ Hot-reloading with Air
- ✅ Delve debugger attached on port 2345
- ✅ Can set breakpoints and inspect variables
- ⚠️ `fmt.Println` output may be buffered/delayed by Delve
- ⚠️ Slightly slower startup

**Use this when:**
- You need to set breakpoints
- You want to step through code
- You need to inspect variable values
- You're debugging complex logic

**View logs:**
```bash
make dev-debug-logs
# or
docker logs -f order-delivery-service-debug
```

**Connect debugger:**
- Host: `localhost`
- Port: `2345`
- API Version: 2

## Debugging Tips

### fmt.Println Not Showing?

If you're using `make dev-debug` and can't see `fmt.Println` output:

1. **Switch to normal dev mode:**
   ```bash
   make dev-debug-down
   make dev-up
   ```

2. **Or use stderr instead of stdout:**
   ```go
   fmt.Fprintf(os.Stderr, "Debug message: %v\n", value)
   os.Stderr.Sync()  // Force immediate flush
   ```

3. **Or use the logger (recommended):**
   ```go
   // Add logger to your repository/service
   logger.Debug("Debug message", zap.Any("value", value))
   ```

### VS Code Debug Configuration

Add to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Delve",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app",
      "port": 2345,
      "host": "localhost"
    }
  ]
}
```

### GoLand/IntelliJ Debug Configuration

1. Run > Edit Configurations
2. Add New Configuration > Go Remote
3. Set Host: `localhost`, Port: `2345`
4. Click Debug

## Common Commands

```bash
# Start normal dev mode
make dev-up

# Stop normal dev mode
make dev-down

# View normal dev logs
make dev-logs

# Start debug mode
make dev-debug

# Stop debug mode
make dev-debug-down

# View debug logs
make dev-debug-logs

# Restart service (keeps DB running)
make dev-restart
```

## SQL Query Logging

SQL query logging is **enabled by default** in both dev and debug modes.

### What You'll See

Every database query is logged to stderr with:
- **File and line number** where the query was executed
- **Execution time** in milliseconds
- **Number of rows** returned
- **Actual SQL** with real values (no `?` placeholders)
- **Color coding** for readability

### Example Output

```
2025/10/24 11:53:53 /app/internal/repository/postgres/repository.go:124
[3.134ms] [rows:1] SELECT count(*) FROM "delivery_assignments"
WHERE (created_at BETWEEN '2024-01-01 00:00:00' AND '2025-12-31 23:59:59')
AND "delivery_assignments"."deleted_at" IS NULL
```

### Controlling SQL Logging

**Disable SQL logging:**
```bash
# Edit docker-compose.dev.yaml or docker-compose.debug.yaml
environment:
  - DB_LOG_SQL=false  # Change from true to false
```

**Enable for local runs (outside Docker):**
```bash
DB_LOG_SQL=true make run
```

### Slow Query Detection

Queries taking longer than **200ms** are highlighted in yellow as slow queries. This helps identify performance bottlenecks.

## Troubleshooting

### Port Already in Use
```bash
# Kill all dev containers
make dev-down
make dev-debug-down

# Or manually
docker ps | grep order-delivery
docker stop <container-id>
```

### Debugger Not Connecting
1. Ensure you're using `make dev-debug` (not `make dev-up`)
2. Check port 2345 is not blocked: `lsof -i :2345`
3. Wait for service to fully start (check logs)
4. Try restarting: `make dev-debug-down && make dev-debug`

### Hot-Reload Not Working
1. Check file was actually changed
2. Ensure file is not in `exclude_dir` in Air config
3. Check container logs for build errors
4. Try manual restart: `make dev-restart`
