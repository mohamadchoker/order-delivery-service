# Go PATH Setup

## Fixed: October 24, 2025

Your `~/.zshrc` has been updated to permanently include `~/go/bin` in your PATH.

## What Was Changed

**Before:**
```bash
echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc  # ❌ Wrong - writes to file
export PATH=$PATH:~/go/bin  # Repeated 8 times
```

**After:**
```bash
# Add Go bin to PATH
export PATH=$PATH:$HOME/go/bin  # ✅ Clean single line
```

## Verify It Works

Open a **new terminal** and run:

```bash
# Check PATH includes go/bin
echo $PATH | grep go/bin

# Verify migrate is accessible
which migrate
# Output: /Users/mohamadchoker/go/bin/migrate

migrate -version
# Output: dev

# Test migrations
make migrate-up
# Should work without PATH errors
```

## Backup

Your original `.zshrc` was backed up to:
```
~/.zshrc.backup.20251024_150429
```

To restore (if needed):
```bash
cp ~/.zshrc.backup.20251024_150429 ~/.zshrc
source ~/.zshrc
```

## Tools Now Accessible

With `~/go/bin` in PATH, all Go-installed tools work globally:

| Tool | Purpose | Install Command |
|------|---------|----------------|
| `migrate` | Database migrations | `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| `protoc-gen-go` | Protocol Buffers codegen | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` |
| `protoc-gen-go-grpc` | gRPC codegen | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| `golangci-lint` | Linter | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| `mockgen` | Mock generation | `go install go.uber.org/mock/mockgen@latest` |
| `air` | Hot-reload | `go install github.com/air-verse/air@latest` |

## For Other Developers

If other developers get "command not found" errors, they should:

1. Check if `~/go/bin` is in PATH:
   ```bash
   echo $PATH | grep go/bin
   ```

2. If not, add to `~/.zshrc` (Zsh) or `~/.bashrc` (Bash):
   ```bash
   echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
   source ~/.zshrc
   ```

3. Or run once:
   ```bash
   make install-tools  # Installs all tools
   ```

## Troubleshooting

### "migrate: command not found"

**In current terminal:**
```bash
source ~/.zshrc
```

**Or open a new terminal** (recommended).

### Tools not found after adding to .zshrc

Check if the line was added correctly:
```bash
tail -5 ~/.zshrc
# Should show: export PATH=$PATH:$HOME/go/bin
```

Make sure there's no `echo` command before it!

### Verify go/bin location

```bash
go env GOPATH
# Output: /Users/mohamadchoker/go

ls -la ~/go/bin/
# Should list: migrate, protoc-gen-go, air, etc.
```
