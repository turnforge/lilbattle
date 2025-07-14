#!/bin/bash

# Build script for WeeWar WASM modules
set -e

# Copy wasm_exec.js from Go installation
echo "Copying wasm_exec.js..."
GO_ROOT=$(go env GOROOT)
WASM_EXEC_PATH=$(find "$GO_ROOT" -name "wasm_exec.js" 2>/dev/null | head -1)
if [ -n "$WASM_EXEC_PATH" ]; then
    cp "$WASM_EXEC_PATH" wasm/
else
    echo "Warning: wasm_exec.js not found in Go installation"
fi

echo "WASM build complete!"
echo "Output files:"
echo "  wasm/weewar-cli.wasm - Main CLI WASM module"
echo "  wasm/editor.wasm     - Map Editor WASM module"
echo "  wasm/wasm_exec.js    - Go WASM runtime"
echo ""
echo "File sizes:"
du -h wasm/*.wasm wasm/*.js
