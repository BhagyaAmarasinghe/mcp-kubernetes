#!/bin/bash

set -e

echo "Building MCP Kubernetes..."
cd "$(dirname "$0")"

# Run go mod tidy to update dependencies
echo "Updating dependencies..."
go mod tidy -v

# Determine Go path
if [ -z "$GOPATH" ]; then
  GOPATH=$(go env GOPATH)
fi

if [ -z "$GOPATH" ]; then
  # If GOPATH is still empty, use home directory
  GOPATH="$HOME/go"
fi

# Create bin directory if it doesn't exist
BIN_DIR="$GOPATH/bin"
mkdir -p "$BIN_DIR"

# Build and install the binary
echo "Building and installing to $BIN_DIR/mcp-kubernetes..."
go build -o "$BIN_DIR/mcp-kubernetes" ./cmd

echo "Installation complete!"
echo "You can now run: mcp-kubernetes"
echo ""
echo "Make sure $BIN_DIR is in your PATH. If not, add it with:"
echo "export PATH=\$PATH:$BIN_DIR"
