# IDE Setup Guide

## GoLand / IntelliJ IDEA

### Fixing "Cannot resolve import 'google/api/annotations.proto'" Error

GoLand needs to be configured to recognize the `third_party` directory as a proto import path.

#### Method 1: Configure Proto Paths (Recommended)

1. Open **Settings/Preferences** (`Cmd+,` on Mac, `Ctrl+Alt+S` on Windows/Linux)
2. Navigate to **Languages & Frameworks → Protocol Buffers**
3. In the **Import Paths** section, click the `+` button
4. Add the following path: `<project-root>/third_party`
5. Make sure **"Configure automatically"** is unchecked
6. Click **OK** to apply

After this, GoLand should resolve the imports correctly.

#### Method 2: Use Buf (Modern Approach)

The project includes a `buf.yaml` configuration file. Install Buf:

```bash
# macOS
brew install bufbuild/buf/buf

# Other platforms
# See: https://docs.buf.build/installation
```

Then Buf will automatically handle proto dependencies and GoLand should recognize them.

#### Method 3: Restart IDE with Cache Invalidation

Sometimes GoLand's cache needs to be cleared:

1. Go to **File → Invalidate Caches...**
2. Select **Invalidate and Restart**
3. After restart, wait for indexing to complete

### Verifying the Fix

Open `proto/delivery.proto` - the import errors should be gone. You should see:
- Green checkmarks on imports
- Auto-completion working for `google.api.http` options
- No red underlines

## VS Code

If using VS Code with the `vscode-proto3` extension:

1. Open **Settings** (`Cmd+,`)
2. Search for "protoc"
3. Add to `protoc` → `Path`:
   ```json
   "protoc.options": [
     "--proto_path=${workspaceRoot}/third_party",
     "--proto_path=${workspaceRoot}"
   ]
   ```

## Vim/Neovim

If using `vim-protobuf` or LSP:

Add to proto LSP config:
```lua
require'lspconfig'.protols.setup{
  cmd = {"protols"},
  root_dir = require'lspconfig'.util.root_pattern("buf.yaml", ".git"),
  settings = {
    proto = {
      paths = {"third_party"}
    }
  }
}
```

## Troubleshooting

### Issue: "Import not found" even after configuration

**Solution**: The proto files must be generated at least once:
```bash
make proto
```

### Issue: GoLand still shows errors after configuration

**Solutions**:
1. Check that `third_party/google/api/annotations.proto` exists
2. Re-sync the project: **File → Sync Project with Gradle Files** (if applicable)
3. Rebuild the project: **Build → Rebuild Project**
4. Clear the cache: **File → Invalidate Caches → Invalidate and Restart**

### Issue: Build works but IDE shows errors

This is usually a cosmetic issue. The build uses `protoc` with the correct paths specified in the Makefile. You can:
- Continue working (builds will succeed)
- Follow Method 1 above to fix IDE recognition
- Use `// noinspection` comments to suppress warnings (not recommended)

## Additional Resources

- [Protocol Buffers Style Guide](https://developers.google.com/protocol-buffers/docs/style)
- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Buf Documentation](https://docs.buf.build/)
