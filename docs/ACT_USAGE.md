# Act - Local GitHub Actions Testing Guide

## Overview

[Act](https://github.com/nektos/act) allows you to run GitHub Actions locally using Docker. This enables you to test your workflows before pushing to GitHub, significantly speeding up development and reducing failed CI runs.

## Installation

### macOS (Homebrew)
```bash
brew install act
```

### Linux
```bash
curl -s https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

### Manual Installation
Download the latest release from [GitHub releases](https://github.com/nektos/act/releases).

## Prerequisites

- **Docker**: Act runs workflows in Docker containers
- **GitHub Actions workflows**: `.github/workflows/*.yml` files

```bash
# Verify Docker is running
docker ps

# Verify act is installed
act --version
```

## Configuration

### Project Configuration Files

#### `.actrc` (Main Configuration)
Located at the project root, this file configures act's behavior:

```bash
# Use medium-sized Docker images
-P ubuntu-latest=catthehacker/ubuntu:act-latest

# Bind local directory
--bind

# Enable verbose output
--verbose
```

#### `.act/.secrets` (Local Secrets)
Store secrets for local testing (NOT committed to git):

```bash
GITHUB_TOKEN=your_token_here
CODECOV_TOKEN=your_codecov_token_here
```

### Image Sizes

Act supports different Docker image sizes:

| Size | Image | Description | Use Case |
|------|-------|-------------|----------|
| Micro | `node:16-buster-slim` | ~150MB | Quick tests, basic jobs |
| Medium | `catthehacker/ubuntu:act-latest` | ~500MB | **Recommended** - Good balance |
| Large | `catthehacker/ubuntu:full-latest` | ~17GB | Full GitHub Actions compatibility |

**We use Medium images** for the best balance of speed and compatibility.

## Basic Usage

### List Available Workflows

```bash
# List all workflows and their jobs
act -l
```

Output:
```
Stage  Job ID     Job name       Workflow name  Workflow file  Events
0      generate   Generate Mocks  CI             ci.yml         push,pull_request
0      lint       Lint           CI             ci.yml         push,pull_request
0      build      Build          CI             ci.yml         push,pull_request
0      docker     Docker Build   CI             ci.yml         push,pull_request
1      test       Test           CI             ci.yml         push,pull_request
```

### Run All Jobs

```bash
# Run all jobs (default event: push)
act

# Or explicitly
act push
```

### Run Specific Jobs

```bash
# Run only the test job
act -j test

# Run only the build job
act -j build

# Run only the lint job
act -j lint

# Run only the docker job
act -j docker

# Run only the generate job
act -j generate
```

### Run Specific Workflows

```bash
# Run specific workflow file
act -W .github/workflows/ci.yml
```

## Advanced Usage

### Dry Run (List Steps Without Executing)

```bash
# See what would run without executing
act -n

# Dry run for specific job
act -j test -n
```

### Run with Secrets

```bash
# Use secrets from file
act --secret-file .act/.secrets

# Pass individual secret
act --secret GITHUB_TOKEN=your_token_here

# Use environment variable
act --env GITHUB_TOKEN=$GITHUB_TOKEN
```

### Run Specific Events

```bash
# Trigger push event
act push

# Trigger pull_request event
act pull_request

# Trigger custom event
act workflow_dispatch
```

### Debugging

```bash
# Verbose output
act -v

# Very verbose output
act -v -v

# Use different shell for debugging
act --use-gitignore=false
```

### Container Options

```bash
# Run in specific container
act -P ubuntu-latest=ubuntu:22.04

# Use custom Docker image
act -P ubuntu-latest=myimage:latest

# Run with privileged mode (for Docker-in-Docker)
act --privileged
```

## Common Workflows for This Project

### Quick Test Before Push

```bash
# Run tests locally before pushing
act -j test

# If tests pass, then push
git push
```

### Full CI Check

```bash
# Run all CI jobs locally
act

# Or run each job individually
act -j generate
act -j test
act -j lint
act -j build
act -j docker
```

### Debug Failed Workflow

```bash
# Run with verbose output
act -j test -v

# Or reuse container for faster iteration
act -j test --reuse
```

### Test Specific Branches

```bash
# Test as if running on main branch
act -e .act/push-main.json

# Test as if running on PR
act pull_request
```

## Project-Specific Commands

### Run Tests Locally (Recommended Workflow)

```bash
# 1. Generate mocks
act -j generate

# 2. Run tests (includes PostgreSQL service)
act -j test

# 3. Run linter
act -j lint

# 4. Build binary
act -j build

# 5. Build Docker image
act -j docker
```

### Quick Validation

```bash
# Run all jobs in parallel (if your machine can handle it)
act --parallel

# Or run sequentially (safer)
act
```

## Troubleshooting

### Issue: Docker Permission Denied

```bash
# Add your user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Or use sudo
sudo act
```

### Issue: Service Containers Not Working

```bash
# Use larger image size
act -P ubuntu-latest=catthehacker/ubuntu:full-latest
```

### Issue: PostgreSQL Connection Failed

The workflow includes PostgreSQL as a service. Act handles this automatically, but if you encounter issues:

```bash
# Check if PostgreSQL service is running in container
act -j test -v

# Look for PostgreSQL startup logs
# Wait for "database system is ready to accept connections"
```

### Issue: Mock Generation Fails

```bash
# Ensure Go is properly set up in container
act -j generate -v

# Check PATH includes ~/go/bin
```

### Issue: Secrets Not Loading

```bash
# Verify secrets file exists
cat .act/.secrets

# Use explicit secret file
act --secret-file .act/.secrets

# Or pass secrets directly
act --secret DB_PASSWORD=postgres
```

### Issue: Out of Disk Space

```bash
# Clean up Docker resources
docker system prune -a

# Remove act cache
rm -rf ~/.cache/act
```

## Performance Tips

### 1. Reuse Containers

```bash
# Reuse containers between runs (faster iteration)
act -j test --reuse
```

**Warning**: Reused containers may have stale state. Clean run for final validation.

### 2. Use Smaller Images for Quick Tests

```bash
# For simple builds, use micro image
act -j build -P ubuntu-latest=node:16-buster-slim
```

### 3. Cache Dependencies

Act automatically caches dependencies like Go modules. The cache persists between runs.

### 4. Run Specific Jobs

```bash
# Don't run all jobs if you only need tests
act -j test
```

### 5. Use Dry Run First

```bash
# Check what will run before executing
act -n
```

## Environment Variables

Act supports environment variables just like GitHub Actions:

```bash
# Set environment variable
act --env DB_HOST=localhost

# Set multiple variables
act --env DB_HOST=localhost --env DB_PORT=5432

# Use .env file
act --env-file .env.test
```

## Comparison: Act vs GitHub Actions

| Feature | Act (Local) | GitHub Actions (Cloud) |
|---------|-------------|------------------------|
| **Speed** | Fast (local execution) | Slower (remote execution) |
| **Cost** | Free (uses local resources) | Free tier, then paid |
| **Debugging** | Easy (direct access) | Limited (logs only) |
| **Environment** | Docker containers | GitHub-hosted runners |
| **Services** | Supported | Fully supported |
| **Secrets** | Local file | GitHub Secrets |
| **Artifacts** | Local filesystem | Cloud storage |
| **Limitations** | Some actions may not work | Full compatibility |

## Best Practices

### 1. Test Locally Before Pushing

```bash
# Always run locally first
act -j test

# If pass, push to GitHub
git push
```

### 2. Use .actrc for Team Consistency

Commit `.actrc` to ensure everyone uses the same configuration.

### 3. Keep Secrets Secure

```bash
# Add .act/.secrets to .gitignore
echo ".act/.secrets" >> .gitignore

# Never commit secrets
```

### 4. Use Dry Run for Workflow Changes

```bash
# After modifying workflow
act -n

# Verify steps look correct before running
act
```

### 5. Clean Up Regularly

```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune
```

## CI/CD Workflow Overview

Our CI workflow (`.github/workflows/ci.yml`) includes:

1. **Generate** - Generate mocks using mockgen
2. **Test** - Run unit tests with PostgreSQL
3. **Lint** - Run golangci-lint
4. **Build** - Build Go binary
5. **Docker** - Build Docker image

### Job Dependencies

```
generate ──> test
```

The test job depends on generate completing successfully.

## Example: Full Local CI Run

```bash
# Step 1: Check what will run
act -l

# Step 2: Dry run
act -n

# Step 3: Run all jobs
act

# Step 4: Check results
# If all jobs pass ✅, safe to push!
git push
```

## Example: Debugging Failed Test

```bash
# Run with verbose output
act -j test -v

# Check PostgreSQL startup
# Look for: "database system is ready to accept connections"

# Check test output
# Look for specific test failures

# Fix code, then rerun
act -j test --reuse
```

## Integration with Makefile

Add act targets to Makefile for convenience:

```makefile
.PHONY: act-test act-all act-list

act-test:  ## Run tests locally with act
	act -j test

act-all:   ## Run all CI jobs locally
	act

act-list:  ## List all available jobs
	act -l
```

Usage:
```bash
make act-test  # Run tests
make act-all   # Run full CI
make act-list  # List jobs
```

## Resources

- **Act GitHub**: https://github.com/nektos/act
- **Act Documentation**: https://nektosact.com
- **GitHub Actions Docs**: https://docs.github.com/en/actions
- **Docker Hub - Act Images**: https://hub.docker.com/r/catthehacker/ubuntu

## Summary

Act provides a powerful way to test GitHub Actions locally:

✅ **Fast feedback** - Test workflows in seconds, not minutes
✅ **Cost effective** - No CI minutes used
✅ **Easy debugging** - Direct access to containers
✅ **Consistent environment** - Same as GitHub Actions
✅ **Offline capable** - Work without internet

**Recommended workflow:**
1. Write code
2. Test with `act -j test`
3. If pass, push to GitHub
4. GitHub Actions runs as final validation

This approach reduces failed CI runs and speeds up development! 🚀
