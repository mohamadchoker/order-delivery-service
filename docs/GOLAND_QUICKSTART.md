# GoLand Quick Start Guide

Quick reference for setting up GoLand with automatic import sorting and formatting for this project.

## 🎯 Quick Setup (5 minutes)

### 1. Import Sorting Configuration

**Settings** (`Cmd+,` on Mac) → **Editor** → **Code Style** → **Go** → **Imports**

Configure these settings:

```
✅ Group stdlib imports
✅ Group current project imports
Sorting type: goimports
Add import with project local prefix: github.com/company/order-delivery-service
```

### 2. Format on Save

**Settings** → **Tools** → **Actions on Save**

```
✅ Reformat code
✅ Optimize imports
✅ Run code cleanup (optional)
```

### 3. Use goimports (NOT gofmt)

**Settings** → **Go** → **Gofmt**

```
✅ On code reformat
Tool to run on reformat: goimports
```

**Done!** Now when you save (`Cmd+S`), imports will automatically sort into 3 groups.

---

## 📋 Import Grouping Result

After configuration, your imports will look like this:

```go
package mypackage

import (
    // Group 1: Standard library
    "context"
    "fmt"
    "time"

    // Group 2: Third-party packages
    "github.com/google/uuid"
    "go.uber.org/zap"
    "google.golang.org/grpc"

    // Group 3: Local project packages
    "github.com/company/order-delivery-service/internal/domain"
    "github.com/company/order-delivery-service/pkg/logger"
)
```

---

## ⌨️ Keyboard Shortcuts

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| Format Code | `Cmd+Option+L` | `Ctrl+Alt+L` |
| Optimize Imports | `Ctrl+Option+O` | `Ctrl+Alt+O` |
| Save File | `Cmd+S` | `Ctrl+S` |
| Settings | `Cmd+,` | `Ctrl+Alt+S` |
| Run Tests | `Ctrl+Shift+R` | `Ctrl+Shift+F10` |
| Run golangci-lint | Right-click → External Tools → golangci-lint |

---

## 🔧 External Tools Setup

### golangci-lint Integration

**Settings** → **Tools** → **External Tools** → **+**

```
Name: golangci-lint
Program: /Users/mohamadchoker/go/bin/golangci-lint
Arguments: run $FileDir$
Working directory: $ProjectFileDir$
```

**Usage**: Right-click on any directory → **External Tools** → **golangci-lint**

### goimports Tool

**Settings** → **Tools** → **External Tools** → **+**

```
Name: goimports
Program: /Users/mohamadchoker/go/bin/goimports
Arguments: -local github.com/company/order-delivery-service -w $FilePath$
Working directory: $ProjectFileDir$
```

**Usage**: Right-click on file → **External Tools** → **goimports**

---

## 🐛 Troubleshooting

### Imports not sorting automatically?

1. **Verify goimports is selected** (NOT gofmt):
   - Settings → Go → Gofmt → "Tool to run on reformat" = **goimports**

2. **Check grouping is enabled**:
   - Settings → Editor → Code Style → Go → Imports
   - Both checkboxes should be ✅
   - Local prefix should be: `github.com/company/order-delivery-service`

3. **Test manually**:
   - Open any `.go` file
   - Press `Cmd+Option+L` (format code)
   - Check if imports reorganize

4. **Invalidate caches**:
   - File → Invalidate Caches → Invalidate and Restart

### goimports not found?

GoLand looks for goimports in `$GOPATH/bin`. Verify it's installed:

```bash
ls -la ~/go/bin/goimports
```

If not found, install it:

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

Then restart GoLand.

---

## 📝 Makefile Commands

Run these from GoLand's terminal (`Option+F12`):

```bash
make lint          # Run linter (zero errors!)
make lint-fix      # Auto-fix linting issues
make format        # Format all Go files
make test          # Run tests
make dev-up        # Start development environment
```

---

## 🎨 Code Style Settings

**Settings** → **Editor** → **Code Style** → **Go**

Recommended settings (already in project):

```
Tabs and Indents:
  ✅ Use tab character
  Tab size: 4
  Indent: 4

Wrapping and Braces:
  Hard wrap at: 120

Imports:
  ✅ Group stdlib imports
  ✅ Group current project imports
  Sorting type: goimports
  Local prefix: github.com/company/order-delivery-service
```

---

## 🚀 Pro Tips

1. **Auto-save**: GoLand auto-saves on focus loss. Just switch windows!

2. **Format on paste**: Settings → Editor → General → Smart Keys → Go
   - ✅ "Reformat on paste"

3. **Run tests quickly**: Click green arrow next to any test function

4. **Navigate to definition**: `Cmd+Click` on any function/type

5. **Find usages**: `Option+F7` on any symbol

6. **Rename refactoring**: `Shift+F6` (renames everywhere, including imports!)

7. **Code completion**: Start typing, `Ctrl+Space` for suggestions

---

## 📚 Learn More

- [GoLand Documentation](https://www.jetbrains.com/help/go/)
- [Project README](../README.md)
- [Full Editor Setup Guide](./EDITOR_SETUP.md)
- [CI/CD Documentation](./CI_CD.md)

---

## ✅ Verification Checklist

Before starting development, verify:

- [ ] GoLand can find Go SDK (Settings → Go → GOROOT)
- [ ] Import grouping is configured (3 groups appear when formatting)
- [ ] Format on save works (save a file, imports reorganize)
- [ ] golangci-lint external tool works (right-click → External Tools)
- [ ] Tests run successfully (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Docker dev environment starts (`make dev-up`)

**All set? Start coding! 🎉**
