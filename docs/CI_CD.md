# CI/CD Guide

## Overview

This project includes a complete CI/CD setup optimized for automated testing, building, and deployment.

## CI Components

### Docker Compose CI (`docker-compose.ci.yaml`)

Optimized for automated pipelines with:
- ✅ Fresh PostgreSQL on tmpfs (faster, no state)
- ✅ Production-like logging
- ✅ Automated migrations
- ✅ Complete test pipeline (lint → migrate → test → build)
- ✅ Exit on completion (perfect for CI)

### CI Dockerfile (`Dockerfile.ci`)

Multi-stage optimized build with:
- ✅ Layer caching for Go modules
- ✅ Pre-installed linter and migrate tool
- ✅ Minimal dependencies
- ✅ Fast rebuilds

### GitHub Actions (`.github/workflows/ci.yaml`)

Complete CI pipeline with parallel jobs:
- ✅ Linting
- ✅ Testing with coverage
- ✅ Building
- ✅ Docker image creation
- ✅ Codecov integration

## Local CI Testing

### Full CI Pipeline

Run the complete pipeline locally:

```bash
make ci-test
```

This runs:
1. Lint code
2. Run migrations
3. Run tests with coverage
4. Build application

**Output:** Same as what runs in GitHub Actions!

### Individual CI Commands

```bash
# Lint only
make ci-lint

# Tests only
make ci-test-only

# Build only
make ci-build

# Generate coverage report
make ci-coverage

# Clean up CI containers
make ci-clean
```

## GitHub Actions Setup

### Automatic Triggers

The CI pipeline runs automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

### Manual Trigger

You can also trigger manually from GitHub UI:
- Go to Actions → CI Pipeline → Run workflow

### Pipeline Jobs

#### 1. Lint Job
- Runs `golangci-lint`
- Catches code quality issues
- Fast feedback (~1-2 minutes)

#### 2. Test Job
- Starts PostgreSQL service
- Runs migrations
- Executes all tests with race detection
- Generates coverage report
- Uploads to Codecov (optional)

#### 3. Build Job
- Builds the application binary
- Uploads artifact for download
- Runs only if lint and test pass

#### 4. Docker Job
- Builds Docker image
- Uses layer caching for speed
- Can push to registry (commented out by default)

## Configuration

### Environment Variables (CI)

The `docker-compose.ci.yaml` sets:

```yaml
DB_HOST: postgres
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: postgres
DB_NAME: order_delivery_db
DB_LOG_SQL: false     # Clean logs in CI
LOG_LEVEL: info
LOG_DEV: false        # Production-like
CI: true
```

### Customize Pipeline

Edit `.github/workflows/ci.yaml` to:
- Add deployment steps
- Enable Docker Hub push
- Add integration tests
- Configure notifications

## Best Practices

### Testing Locally Before Push

Always test CI locally before pushing:

```bash
# Run full pipeline
make ci-test

# If it passes locally, it should pass in GitHub Actions
```

### Coverage Requirements

Currently no minimum coverage enforced. To add:

```yaml
# In .github/workflows/ci.yaml
- name: Check coverage
  run: |
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$coverage < 80" | bc -l) )); then
      echo "Coverage $coverage% is below 80%"
      exit 1
    fi
```

### Caching Strategy

The pipeline caches:
- ✅ Go modules (`~/go/pkg/mod`)
- ✅ Docker layers (BuildKit cache)

This speeds up CI by ~50-70%.

### PostgreSQL Service

Uses GitHub Actions services for PostgreSQL:
- Starts before tests
- Health checks ensure it's ready
- Isolated per workflow run
- Automatic cleanup

## Deploying

### Docker Registry Push

To enable Docker Hub push, uncomment in `.github/workflows/ci.yaml`:

```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKER_USERNAME }}
    password: ${{ secrets.DOCKER_PASSWORD }}

- name: Push Docker image
  uses: docker/build-push-action@v5
  with:
    push: true
    tags: |
      your-dockerhub-username/order-delivery-service:latest
      your-dockerhub-username/order-delivery-service:${{ github.sha }}
```

Then add secrets:
1. Go to GitHub repo → Settings → Secrets
2. Add `DOCKER_USERNAME`
3. Add `DOCKER_PASSWORD`

### GitHub Container Registry (GHCR)

Alternative to Docker Hub:

```yaml
- name: Login to GHCR
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Push to GHCR
  uses: docker/build-push-action@v5
  with:
    push: true
    tags: |
      ghcr.io/${{ github.repository }}:latest
      ghcr.io/${{ github.repository }}:${{ github.sha }}
```

### Kubernetes Deployment

Add deployment step after build:

```yaml
- name: Deploy to Kubernetes
  run: |
    kubectl set image deployment/order-delivery-service \
      app=ghcr.io/${{ github.repository }}:${{ github.sha }} \
      --namespace=production
```

## Troubleshooting

### "Linter failed"

Run locally to see issues:
```bash
make ci-lint
```

Fix issues:
```bash
make lint-fix
```

### "Tests failed in CI but pass locally"

Common causes:
1. **Race conditions** - CI uses `-race` flag
2. **Different environment** - CI uses production-like config
3. **Database state** - CI uses fresh DB each time

Debug:
```bash
# Run with same flags as CI
go test -v -race ./...

# Run in CI environment
make ci-test
```

### "Build takes too long"

Optimize:
1. Check if Go modules are cached
2. Ensure Docker layer caching works
3. Consider using `go.sum` hash for cache key

### "Coverage report not uploaded"

Check Codecov token:
```bash
# Add CODECOV_TOKEN to GitHub secrets
# Get it from codecov.io
```

## CI Performance

### Current Pipeline Times

| Job | Duration | Parallelized |
|-----|----------|--------------|
| Lint | ~1-2 min | ✅ Yes |
| Test | ~2-3 min | ✅ Yes |
| Build | ~1-2 min | ✅ Yes (after lint+test) |
| Docker | ~2-3 min | ✅ Yes (after lint+test) |

**Total:** ~3-5 minutes (with parallel execution)

### Optimization Tips

1. **Use caching** - Already enabled for Go modules and Docker
2. **Run jobs in parallel** - Lint and Test run simultaneously
3. **Fail fast** - Linter runs first (fastest)
4. **Skip unnecessary work** - Build only runs if tests pass

## Adding New Tests

1. Write test in appropriate `*_test.go` file
2. Run locally: `go test -v ./internal/...`
3. Ensure it passes with race detection: `go test -race ./...`
4. Push - CI runs automatically
5. Check GitHub Actions for results

## Monitoring CI

### GitHub Actions UI

View pipeline status:
- Repo → Actions tab
- See all workflow runs
- Click run to see job details
- Download artifacts (build binaries)

### Badges

Add to README.md:

```markdown
![CI](https://github.com/your-username/order-delivery-service/workflows/CI%20Pipeline/badge.svg)
```

### Notifications

Configure in repo settings:
- Email notifications
- Slack integration
- Custom webhooks

## CI vs Local Development

| Aspect | Local Dev | CI |
|--------|-----------|-----|
| **Environment** | `make dev-up` | `docker-compose.ci.yaml` |
| **Database** | Persistent | Fresh (tmpfs) |
| **Logging** | Development | Production-like |
| **SQL Logs** | Enabled | Disabled |
| **Speed** | Hot-reload | Full build each time |
| **Purpose** | Fast feedback | Quality gate |

## Future Enhancements

Consider adding:
- [ ] Integration tests
- [ ] Performance benchmarks
- [ ] Security scanning (Snyk, Trivy)
- [ ] Dependency updates (Dependabot)
- [ ] Automated releases
- [ ] Staging deployment
- [ ] E2E tests

## Resources

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Docker BuildKit Caching](https://docs.docker.com/build/cache/)
- [golangci-lint Configuration](https://golangci-lint.run/)
- [Codecov Documentation](https://docs.codecov.com/)
