#!/bin/bash

# Build script for WeeWar WASM modules
set -e

echo "Building WeeWar WASM modules..."

# Create output directory
mkdir -p wasm

# Build main CLI WASM
echo "Building weewar-cli WASM..."
cd cmd/weewar-wasm
GOOS=js GOARCH=wasm go build -o ../../wasm/weewar-cli.wasm
cd ../..

# Build map editor WASM
echo "Building map editor WASM..."
cd cmd/editor-wasm
GOOS=js GOARCH=wasm go build -o ../../wasm/editor.wasm
cd ../..

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