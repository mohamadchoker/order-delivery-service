# Act Quick Reference Card

## 🚀 Common Commands

```bash
# List all jobs
make act-list

# Run specific jobs (recommended)
make act-test      # Run tests only
make act-lint      # Run linter only (fastest)
make act-build     # Run build only

# Run all jobs
make act-all

# Clean up
make act-clean
```

---

## ⚡ Fast Workflow

### Daily Development

```bash
# 1. Make code changes
vim internal/service/delivery_usecase.go

# 2. Quick local test (5 seconds)
go test ./internal/service/

# 3. Before pushing: Test with act (30 seconds)
make act-test

# 4. If pass: Push confidently
git push
```

---

## 📊 Download Sizes

| First Time | Subsequent Runs |
|------------|-----------------|
| ~1.7 GB (5-10 min) | ~10 MB (30 sec) |

**Don't worry!** Downloads only happen once. Everything is cached.

---

## 🎯 What to Run

### Before Every Push

```bash
make act-test    # Catches test failures
```

### Before Pull Request

```bash
make act-all     # Run full CI
```

### Quick Lint Check

```bash
make act-lint    # Fastest (no database needed)
```

### Quick Build Check

```bash
make act-build   # Verify build works
```

---

## 🐛 Troubleshooting

### Act is slow

```bash
# Clean up and rerun
make act-clean
make act-test
```

### Out of disk space

```bash
# Clean Docker
docker system prune -a
```

### Tests failing locally but not in act

```bash
# Run with verbose output
act -j test -v
```

---

## 💡 Pro Tips

1. **Run specific jobs** - Don't run all jobs every time
2. **Use local tests first** - `go test` is faster
3. **Act before pushing** - Catches CI failures
4. **Cache is your friend** - First run is slow, then fast
5. **Clean regularly** - `make act-clean` weekly

---

## 📚 Full Documentation

- **ACT_USAGE.md** - Comprehensive guide
- **ACT_DOWNLOADS.md** - Download optimization
- **ACT_CI_SETUP.md** - Setup details
- **README.md** - Main project docs

---

## 🎮 Advanced Usage

```bash
# Run with different config
act --rc-file .actrc.small -j test

# Dry run (see what would run)
act -n

# Verbose output
act -j test -v

# Reuse containers (faster)
act -j test --reuse

# Run specific event
act push
act pull_request
```

---

## ✅ Verification

```bash
# Check act is installed
act --version

# Check available jobs
make act-list

# Should show:
# - generate
# - test
# - lint
# - build
# - docker
```

---

## 🆘 Help

```bash
# See all make commands
make help

# See all act options
act --help

# Get help on specific job
act -j test --help
```

---

**Remember**: Act replicates GitHub Actions locally. First run downloads everything (~1.7 GB), but then it's cached and super fast! 🚀
