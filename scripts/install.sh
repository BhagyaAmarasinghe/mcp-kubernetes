#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.23 or higher."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "kubectl is not installed. Please install kubectl."
    exit 1
fi

# Install the MCP Kubernetes server
go install github.com/BhagyaAmarasinghe/mcp-kubernetes@latest

echo "MCP Kubernetes server installed successfully."
echo "You can now configure Claude Desktop to use this server."
echo ""
echo "1. Open Claude Desktop"
echo "2. Add the following to your Claude Desktop configuration file:"
echo ""
echo '{
  "mcpServers": {
    "kubernetes": {
      "command": "mcp-kubernetes"
    }
  }
}'
echo ""
echo "3. Restart Claude Desktop"
echo ""
echo "You can now use Claude to execute Kubernetes commands!"
