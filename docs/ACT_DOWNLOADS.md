# Understanding Act Downloads

## Why Does Act Download So Much?

Act replicates the GitHub Actions environment locally using Docker, which requires downloading various components. Here's what happens and how to optimize it.

---

## What Gets Downloaded

### 1. Docker Images (Biggest Download)

Act uses Docker images that emulate GitHub's runners. You have three options:

| Image Size | Name | Download Size | Use Case |
|------------|------|---------------|----------|
| **Micro** | `node:16-buster-slim` | ~150 MB | Quick tests, simple jobs |
| **Medium** ⭐ | `catthehacker/ubuntu:act-latest` | ~500 MB | **Recommended** - Good balance |
| **Large** | `catthehacker/ubuntu:full-latest` | ~17 GB | Full compatibility |

**Default configuration uses Medium images** (~500 MB).

### 2. GitHub Actions (Per Workflow)

Each action in your workflow gets downloaded:

```yaml
# From .github/workflows/ci.yml
actions/checkout@v4              # ~50 MB
actions/setup-go@v5              # ~100 MB
actions/cache@v4                 # ~30 MB
actions/upload-artifact@v4       # ~40 MB
docker/setup-buildx-action@v3    # ~80 MB
docker/build-push-action@v5      # ~100 MB
golangci/golangci-lint-action@v6 # ~200 MB
codecov/codecov-action@v4        # ~50 MB

Total: ~650 MB
```

### 3. Go Dependencies

During job execution:

```bash
# Go modules (from go.mod)
go mod download                  # ~200-300 MB

# Go tools
mockgen                          # ~50 MB
golangci-lint                    # ~100 MB
protoc-gen-go                    # ~30 MB
```

### 4. PostgreSQL Image

For tests:

```bash
postgres:14-alpine               # ~80 MB
```

---

## Total Download Breakdown

### First Run (Everything)

```
Docker Image (medium):    ~500 MB
GitHub Actions:           ~650 MB
Go modules:               ~250 MB
Go tools:                 ~180 MB
PostgreSQL:                ~80 MB
─────────────────────────────────
TOTAL FIRST RUN:         ~1.66 GB
```

**Time: 5-10 minutes** (depends on internet speed)

### Subsequent Runs (Cached)

```
Everything is cached!      ~0 MB
Only new/updated packages  ~5-20 MB
─────────────────────────────────
TOTAL SUBSEQUENT RUNS:    ~10 MB
```

**Time: 30-60 seconds** ⚡

---

## Download Timeline

### First Time Using Act

```bash
$ act -j test

# Step 1: Pull Docker image (500 MB) - 2-3 minutes
[Test/test] 🚀  Start image=catthehacker/ubuntu:act-latest
[Test/test]   🐳  docker pull catthehacker/ubuntu:act-latest
Pulling from catthehacker/ubuntu... ████████████ 100%

# Step 2: Download actions (~650 MB) - 2-3 minutes
[Test/test]   ✅  Success - Checkout code
[Test/test]   🐳  docker pull docker://v1/actions/checkout@v4
[Test/test]   🐳  docker pull docker://v1/actions/setup-go@v5
... (downloading more actions)

# Step 3: Install Go modules (~250 MB) - 1-2 minutes
[Test/test] ⭐  Run Install dependencies
go: downloading github.com/...
go: downloading google.golang.org/...

# Step 4: Run tests - 1-2 minutes
[Test/test] ⭐  Run Run tests
Running tests...

TOTAL TIME: 5-10 minutes
```

### Second Time (Everything Cached)

```bash
$ act -j test

# Docker image: Already cached ✅
[Test/test] 🚀  Start image=catthehacker/ubuntu:act-latest (cached)

# Actions: Already cached ✅
[Test/test]   ✅  Success - Checkout code (cached)
[Test/test]   ✅  Success - Set up Go (cached)

# Go modules: Already cached ✅
[Test/test] ⭐  Run Install dependencies (using cache)

# Run tests - 30 seconds
[Test/test] ⭐  Run Run tests
Running tests...

TOTAL TIME: 30-60 seconds ⚡
```

---

## How to Optimize

### Option 1: Use Smaller Image (Fastest Downloads)

**Trade-off**: Some actions may not work

```bash
# Edit .actrc
-P ubuntu-latest=node:16-buster-slim  # Only 150 MB!

# Or use the pre-configured small profile
act --rc-file .actrc.small -j test
```

**Saves**: ~350 MB vs medium image

### Option 2: Run Specific Jobs Only

Don't run all jobs if you only need one:

```bash
# Only run linting (fastest, minimal downloads)
make act-lint

# Only run build (fast)
make act-build

# Only run tests
make act-test

# Don't run all jobs unless needed
# make act-all  ← Downloads everything
```

### Option 3: Reuse Containers

Add `--reuse` flag to keep containers between runs:

```bash
# Edit .actrc
--reuse

# Now containers stay alive between runs
act -j test  # First run
act -j test  # Second run uses same container (very fast!)
```

**Warning**: Reused containers may have stale state.

### Option 4: Use GitHub Actions Cache

The workflow is already optimized with caching:

```yaml
# Go modules are cached
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/go/pkg/mod
      ~/go/bin
```

After first run, Go modules are cached and won't re-download.

### Option 5: Pre-pull Images

Pull images ahead of time when on fast internet:

```bash
# Pull the medium image (recommended)
docker pull catthehacker/ubuntu:act-latest

# Pull PostgreSQL
docker pull postgres:14-alpine

# Now act won't need to download these
act -j test
```

---

## Configuration Files

### Default (`.actrc`) - Balanced

```bash
# Good balance: 500 MB image
-P ubuntu-latest=catthehacker/ubuntu:act-latest
--bind
--verbose
```

**Use for**: Most development work
**First download**: ~1.7 GB
**Subsequent runs**: ~10 MB

### Small (`.actrc.small`) - Minimal

```bash
# Minimal: 150 MB image
-P ubuntu-latest=node:16-buster-slim
--bind
```

**Use for**: Quick tests, simple jobs
**First download**: ~1.4 GB (saves 350 MB)
**Subsequent runs**: ~10 MB
**Caveat**: Some actions may not work

### Optimized (`.actrc.optimized`) - Recommended

```bash
# Best of both worlds
-P ubuntu-latest=catthehacker/ubuntu:act-latest
--reuse
--bind
--container-architecture linux/amd64
```

**Use for**: Regular development
**First download**: ~1.7 GB
**Subsequent runs**: ~5 MB (containers reused)

---

## Usage Examples

### Use Default Configuration

```bash
# Uses .actrc (medium image)
make act-test
```

### Use Small Image (Less Downloads)

```bash
# Uses .actrc.small (micro image)
act --rc-file .actrc.small -j test
```

### Use Optimized (Fastest After First Run)

```bash
# Uses .actrc.optimized (reuses containers)
act --rc-file .actrc.optimized -j test
```

---

## What's Cached and Where

### Docker Images

```bash
# Location
~/.docker/

# View cached images
docker images | grep act

# Output
catthehacker/ubuntu    act-latest    abc123    500 MB
postgres               14-alpine     def456     80 MB
```

### Go Modules

```bash
# Location (in container)
~/go/pkg/mod/

# Persisted via GitHub Actions cache
```

### Action Runners

```bash
# Location
~/.cache/act/

# View size
du -sh ~/.cache/act
# Output: ~650 MB
```

---

## Cleaning Up

### Remove Act Cache

```bash
# Clean Docker containers and images
make act-clean

# Or manually
docker container prune -f
docker image prune -f

# Remove act cache
rm -rf ~/.cache/act
```

### Remove Specific Image

```bash
# Remove medium image
docker rmi catthehacker/ubuntu:act-latest

# Remove PostgreSQL
docker rmi postgres:14-alpine

# Next run will re-download
```

### Check Disk Usage

```bash
# Docker disk usage
docker system df

# Output
TYPE            TOTAL     ACTIVE    SIZE
Images          5         2         1.2GB
Containers      3         1         500MB
Local Volumes   1         1         100MB
```

---

## FAQ

### Q: Why so big compared to just running `go test`?

**A**: Act creates a full GitHub Actions environment with:
- Ubuntu container
- GitHub Actions runtime
- All actions you use
- PostgreSQL service container

This is the price for 100% CI/CD parity.

### Q: Can I use my local Go cache?

**A**: Yes! Use bind mounts:

```bash
# Mount local Go cache (already in .actrc)
--bind

# Now act uses your local ~/go/pkg/mod
```

### Q: Do I need to download every time?

**A**: No! Only first time. Everything is cached after that.

```bash
First run:  ~1.7 GB download, 5-10 minutes
Second run: ~10 MB updates, 30-60 seconds
```

### Q: Is there a way to use GitHub Actions cache locally?

**A**: Partially. The workflow uses `actions/cache@v4` which caches Go modules in the container. This persists between act runs if you use `--reuse`.

### Q: Can I pre-download everything?

**A**: Yes!

```bash
# Pull images
docker pull catthehacker/ubuntu:act-latest
docker pull postgres:14-alpine

# Run dry-run (downloads actions)
act -n

# Now everything is cached
```

---

## Comparison: Local vs Act

| Aspect | Local `go test` | Act `act -j test` |
|--------|----------------|-------------------|
| **First run** | 0 MB | ~1.7 GB |
| **Subsequent** | 0 MB | ~10 MB |
| **Speed** | 5 seconds | 30 seconds |
| **Environment** | Your machine | GitHub Actions replica |
| **CI parity** | No | Yes ✅ |
| **Value** | Quick tests | CI verification |

---

## Recommendations

### For Daily Development

```bash
# Use local testing for speed
go test ./...              # 5 seconds

# Use act before pushing
make act-test              # 30 seconds (cached)

# Use GitHub Actions for final validation
git push                   # 2-5 minutes
```

### For First Time Setup

```bash
# Pre-download when on fast internet
docker pull catthehacker/ubuntu:act-latest
docker pull postgres:14-alpine

# Then use act regularly
make act-test
```

### For Limited Bandwidth

```bash
# Use smallest image
act --rc-file .actrc.small -j test

# Or just use local testing
go test ./...

# Push to GitHub and let cloud CI run
git push
```

---

## Size Optimization Summary

| Method | Saves | Trade-off |
|--------|-------|-----------|
| Use micro image | 350 MB | Some actions may fail |
| Run specific jobs only | 500 MB | Don't test everything |
| Pre-pull images | 0 MB (but faster) | Initial time investment |
| Clean cache regularly | 1.7 GB | Re-download next time |
| Use `--reuse` | Faster runs | May have stale state |

---

## Conclusion

**Yes, act downloads a lot (~1.7 GB first time)**, but:

✅ **Only once** - Everything is cached after first run
✅ **Subsequent runs** - Only ~10 MB updates
✅ **Worth it** - Catches CI failures before pushing
✅ **Optimizable** - Use smaller images or specific jobs
✅ **Time saved** - No waiting for GitHub Actions to fail

**Recommended workflow:**

```bash
# Quick local tests (fast)
go test ./...

# Before pushing (catches CI issues)
make act-test

# After all local tests pass
git push
```

This approach minimizes downloads while maximizing confidence! 🚀
