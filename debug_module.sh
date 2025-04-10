#!/bin/bash

# Print debugging information for Go modules

echo "========== Go Environment =========="
go env GOPATH
go env GOMODCACHE
go version

echo -e "\n========== Go Module Info =========="
cd /Users/bhagya/work/personal/mcp-kubernetes
echo "Current directory: $(pwd)"
echo "Module path in go.mod: $(grep '^module' go.mod)"

echo -e "\n========== Go Module Files =========="
find . -name "go.mod" -o -name "go.sum"

echo -e "\n========== Try Local Build =========="
cd /Users/bhagya/work/personal/mcp-kubernetes
go build -v ./cmd
echo "Build exit code: $?"

echo -e "\n=== Alternative Installation Method ==="
echo "Instead of using 'go install github.com/BhagyaAmarasinghe/mcp-kubernetes/cmd@latest', try:"
echo "1. cd /Users/bhagya/work/personal/mcp-kubernetes"
echo "2. go build -o \$GOPATH/bin/mcp-kubernetes ./cmd"
echo "This should build and install directly from your local code."
