# Editor Setup Guide

## VS Code Setup (Recommended)

### Required Extensions

Install the Go extension:
```bash
code --install-extension golang.go
```

Or install all recommended extensions at once:
1. Open project in VS Code
2. Press `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
3. Type "Extensions: Show Recommended Extensions"
4. Click "Install All"

### Configuration

The project includes `.vscode/settings.json` with:

✅ **Auto-format on save** with `goimports`
✅ **Auto-sort imports** (stdlib → 3rd-party → local)
✅ **Linting** with `golangci-lint`
✅ **Race detection** in tests
✅ **Inlay hints** for better code understanding

### Import Sorting

Imports are automatically sorted into groups:

```go
package mypackage

import (
    // Standard library (built-in Go packages)
    "context"
    "fmt"
    "time"

    // Third-party packages
    "github.com/google/uuid"
    "go.uber.org/zap"
    "google.golang.org/grpc"

    // Local/company packages
    "github.com/company/order-delivery-service/internal/domain"
    "github.com/company/order-delivery-service/pkg/logger"
)
```

**How it works:**
- `goimports` with `-local github.com/company/order-delivery-service`
- Runs automatically on save
- Can be triggered manually: `Cmd+Shift+I` or right-click → "Format Document"

### Keyboard Shortcuts

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| Format Document | `Shift+Option+F` | `Shift+Alt+F` |
| Organize Imports | `Cmd+Shift+I` | `Ctrl+Shift+I` |
| Go to Definition | `F12` | `F12` |
| Find References | `Shift+F12` | `Shift+F12` |
| Rename Symbol | `F2` | `F2` |
| Run Tests | `Cmd+Shift+T` | `Ctrl+Shift+T` |

### Running Commands from VS Code

Open integrated terminal: `` Ctrl+` `` (backtick)

Then run:
```bash
make dev-up        # Start development
make test          # Run tests
make lint          # Run linter
```

## GoLand/IntelliJ IDEA Setup

GoLand has excellent built-in Go support. Follow these steps to enable automatic import sorting and formatting.

### Step 1: Configure Import Grouping

1. Open **Settings/Preferences** (`Cmd+,` on Mac, `Ctrl+Alt+S` on Windows/Linux)
2. Navigate to: **Editor** → **Code Style** → **Go** → **Imports** tab
3. **Enable "Group stdlib imports"** (checkbox)
4. **Enable "Group current project imports"** (checkbox)
5. In **"Sorting type"**: Select **"goimports"**
6. In **"Add import with project local prefix"**: Enter `github.com/company/order-delivery-service`
7. Click **Apply**

This will create 3 groups:
- Group 1: Standard library imports
- Group 2: Third-party imports
- Group 3: Local project imports (github.com/company/order-delivery-service/*)

### Step 2: Enable Format on Save

1. In **Settings/Preferences** (`Cmd+,` or `Ctrl+Alt+S`)
2. Navigate to: **Tools** → **Actions on Save**
3. Enable the following checkboxes:
   - ✅ **Reformat code** (formats according to Go standards)
   - ✅ **Optimize imports** (removes unused, sorts groups)
   - ✅ **Run code cleanup** (optional, applies code inspections)
4. Click **Apply** and **OK**

Now every time you save a file (`Cmd+S`), GoLand will automatically:
- Format the code
- Sort imports into 3 groups
- Remove unused imports

### Step 3: Configure gofmt/goimports

1. In **Settings/Preferences** (`Cmd+,` or `Ctrl+Alt+S`)
2. Navigate to: **Go** → **Gofmt**
3. Ensure **"On code reformat"** is enabled
4. **Tool to run on reformat**: Select **"goimports"** (NOT gofmt)
5. Click **Apply** and **OK**

### Manual Format Shortcuts

If you want to format manually without saving:

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| Reformat Code | `Cmd+Option+L` | `Ctrl+Alt+L` |
| Optimize Imports | `Ctrl+Option+O` | `Ctrl+Alt+O` |
| Reformat File | `Cmd+Option+Shift+L` | `Ctrl+Alt+Shift+L` |

### External Tools

#### golangci-lint

1. **Settings** → **Tools** → **External Tools** → **+**
2. Configure:
   - **Name:** `golangci-lint`
   - **Program:** `$HOME/go/bin/golangci-lint`
   - **Arguments:** `run $FileDir$`
   - **Working directory:** `$ProjectFileDir$`

#### goimports

1. **Settings** → **Tools** → **File Watchers** → **+** → **Go fmt**
2. Change:
   - **Program:** `$GOPATH$/bin/goimports`
   - **Arguments:** `-local github.com/company/order-delivery-service -w $FilePath$`

## Neovim/Vim Setup

### Using vim-go

Add to your config:

```vim
" Install vim-go
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

" Auto-format on save
let g:go_fmt_command = "goimports"
let g:go_fmt_options = {
  \ 'goimports': '-local github.com/company/order-delivery-service',
  \ }

" Linting
let g:go_metalinter_command = "golangci-lint"
let g:go_metalinter_autosave = 1

" Syntax highlighting
let g:go_highlight_functions = 1
let g:go_highlight_methods = 1
let g:go_highlight_structs = 1
let g:go_highlight_operators = 1
let g:go_highlight_build_constraints = 1

" Auto-import packages
let g:go_fmt_autosave = 1
let g:go_imports_autosave = 1
```

### Using Neovim with LSP

```lua
-- Using nvim-lspconfig
require'lspconfig'.gopls.setup{
  settings = {
    gopls = {
      gofumpt = true,
      analyses = {
        unusedparams = true,
        shadow = true,
      },
      staticcheck = true,
      ["local"] = "github.com/company/order-delivery-service",
    },
  },
}

-- Auto-format on save
vim.api.nvim_create_autocmd("BufWritePre", {
  pattern = "*.go",
  callback = function()
    vim.lsp.buf.format({ async = false })
  end,
})
```

## Command Line Tools

### Install goimports

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

### Manual Import Sorting

Format a single file:
```bash
~/go/bin/goimports -local github.com/company/order-delivery-service -w file.go
```

Format all files using Make:
```bash
make format
```

Or manually:
```bash
find . -name "*.go" -not -path "./vendor/*" -not -path "*/.pb.go" | \
  xargs ~/go/bin/goimports -local github.com/company/order-delivery-service -w
```

### golangci-lint

Run linter:
```bash
make lint
# or manually:
~/go/bin/golangci-lint run ./...
```

Auto-fix issues:
```bash
make lint-fix
# or manually:
~/go/bin/golangci-lint run --fix ./...
```

Run specific linters:
```bash
~/go/bin/golangci-lint run --disable-all --enable=goimports ./...
```

## Git Hooks (Pre-commit)

Automatically format and lint before commits:

### Using pre-commit framework

Create `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: goimports
        name: goimports
        entry: goimports
        args: [-local, github.com/company/order-delivery-service, -w]
        language: system
        files: \.go$

      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint
        args: [run, --fix]
        language: system
        files: \.go$
        pass_filenames: false
```

Install hooks:
```bash
pip install pre-commit
pre-commit install
```

### Manual Git Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/sh

# Format code
echo "Running goimports..."
goimports -local github.com/company/order-delivery-service -w $(find . -name "*.go" -not -path "./vendor/*" -not -path "*.pb.go")

# Run linter
echo "Running golangci-lint..."
golangci-lint run --fix ./...

# Add formatted files
git add -u

exit 0
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Troubleshooting

### "goimports not found"

Install it:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

Ensure `~/go/bin` is in PATH:
```bash
export PATH=$PATH:$HOME/go/bin
```

### Auto-format not working in VS Code

1. **Verify goimports is installed**:
   ```bash
   ~/go/bin/goimports -version
   ```

2. **Check Go extension is installed**:
   - Open Extensions: `Cmd+Shift+X`
   - Search for "Go" by Go Team at Google
   - Click "Install" if not installed

3. **Reload VS Code window**:
   - `Cmd+Shift+P` → "Developer: Reload Window"

4. **Check Go extension is using goimports**:
   - `Cmd+Shift+P` → "Go: Locate Configured Go Tools"
   - Verify `goimports` is listed and path is correct

5. **Manual format to test**:
   - Open a .go file
   - `Shift+Option+F` (Mac) or `Shift+Alt+F` (Windows/Linux)
   - Check if imports get organized

6. **Check Output panel for errors**:
   - `Cmd+Shift+U` → Select "Go" from dropdown
   - Look for any error messages

### "golangci-lint not found"

Install with PostgreSQL support:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### VS Code not auto-formatting

1. Check Go extension is installed
2. Verify settings: `Cmd+,` → search "format on save"
3. Ensure `editor.formatOnSave` is enabled for Go files
4. Reload window: `Cmd+Shift+P` → "Reload Window"

### Imports not being sorted (VS Code)

Check the format tool:
```json
{
  "go.formatTool": "goimports",  // Should be goimports, not gofmt
  "go.formatFlags": [
    "-local",
    "github.com/company/order-delivery-service"
  ]
}
```

### Imports not being sorted (GoLand)

1. **Check goimports is selected**:
   - Settings → Go → Gofmt
   - Verify "Tool to run on reformat" is set to **goimports**

2. **Check import grouping is enabled**:
   - Settings → Editor → Code Style → Go → Imports
   - ✅ "Group stdlib imports" should be checked
   - ✅ "Group current project imports" should be checked
   - "Add import with project local prefix" should be: `github.com/company/order-delivery-service`

3. **Verify Actions on Save**:
   - Settings → Tools → Actions on Save
   - ✅ "Reformat code" should be enabled
   - ✅ "Optimize imports" should be enabled

4. **Test manually**:
   - Open a .go file
   - Press `Cmd+Option+L` (Mac) or `Ctrl+Alt+L` (Windows/Linux)
   - Imports should reorganize into 3 groups

5. **Invalidate caches** (if still not working):
   - File → Invalidate Caches → Invalidate and Restart

### Linter taking too long

Use `--fast` flag:
```bash
golangci-lint run --fast ./...
```

Or in VS Code settings:
```json
{
  "go.lintFlags": ["--fast"]
}
```

## Linter Configuration

See `.golangci.yml` for current configuration.

**Key settings:**
- Import grouping with `goimports`
- Misspelling exceptions (e.g., "cancelled")
- Security checks with `gosec`
- Code complexity limits
- Style checking with `stylecheck`

**Excluded checks:**
- Package comments (too verbose for internal packages)
- Some security warnings for internal/metrics endpoints
- Integer overflow warnings (we validate bounds)

## CI Integration

The CI pipeline automatically:
1. Checks formatting with `goimports`
2. Runs `golangci-lint`
3. Fails build if issues found

Test locally before pushing:
```bash
make ci-lint
```

## Best Practices

1. **Always format before committing**
   - Use editor auto-format or run `make lint-fix`

2. **Keep imports organized**
   - 3 groups: stdlib, 3rd-party, local
   - Blank line between groups

3. **Run linter locally**
   - Faster feedback than waiting for CI
   - `make lint` or `golangci-lint run ./...`

4. **Use editor features**
   - Go to definition (F12)
   - Find references (Shift+F12)
   - Rename refactoring (F2)

5. **Configure once, use everywhere**
   - `.golangci.yml` - shared linter config
   - `.vscode/settings.json` - editor config
   - Works for all team members

## Resources

- [goimports documentation](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [golangci-lint](https://golangci-lint.run/)
- [VS Code Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
- [vim-go](https://github.com/fatih/vim-go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
