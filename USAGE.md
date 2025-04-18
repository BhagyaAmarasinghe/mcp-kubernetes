# MCP Kubernetes Usage Guide

This document provides detailed usage instructions for the MCP Kubernetes server.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- kubectl installed and properly configured
- A valid kubeconfig file

### Installation

Install the MCP Kubernetes server using Go:

```bash
go install github.com/BhagyaAmarasinghe/mcp-kubernetes@latest
```

Or use the installation script:

```bash
./scripts/install.sh
```

## Server Modes

### Stdio Mode (Default)

When run without arguments, the server operates in stdio mode, which is ideal for integration with Claude Desktop:

```bash
mcp-kubernetes
```

### HTTP Mode

For integration with other MCP clients, you can run in HTTP mode by specifying a port:

```bash
mcp-kubernetes -port 8080
```

This will start an HTTP server on the specified port.

## Tool Documentation

### execute

Executes any kubectl command.

**Parameters:**
- `command` (string, required): The kubectl command to execute without the "kubectl" prefix

**Example:**
```json
{
  "command": "get pods -n default"
}
```

### get-contexts

Lists all available Kubernetes contexts from your kubeconfig.

**Parameters:** None

**Example:**
```json
{}
```

### current-context

Shows the currently active Kubernetes context.

**Parameters:** None

**Example:**
```json
{}
```

### set-context

Switches to a different Kubernetes context.

**Parameters:**
- `context` (string, required): Name of the context to set as current

**Example:**
```json
{
  "context": "minikube"
}
```

## Using with Claude Desktop

1. Install Claude Desktop from https://claude.ai/download
2. Configure Claude Desktop to use the MCP Kubernetes server:
   - On macOS: Edit `~/Library/Application Support/Claude/claude_desktop_config.json`
   - On Windows: Edit `%APPDATA%\Claude\claude_desktop_config.json`
   - Add the following configuration:
   ```json
   {
     "mcpServers": {
       "kubernetes": {
         "command": "mcp-kubernetes"
       }
     }
   }
   ```
3. Restart Claude Desktop
4. You can now ask Claude to execute Kubernetes commands

## Example Prompts

Here are some example prompts you can use with Claude:

- "Show me the pods in my current namespace"
- "List all services across all namespaces"
- "Switch to the minikube context and then show me the deployments"
- "Get the logs from the pod named nginx in the default namespace"
- "Check the status of the deployment named frontend"

## Troubleshooting

### Server Doesn't Start

- Ensure kubectl is properly installed and in your PATH
- Verify your kubeconfig file exists and is valid
- Check for appropriate permissions to access your kubeconfig

### Claude Cannot Connect to Server

- Verify the server is running
- Check the configuration in `claude_desktop_config.json`
- Restart Claude Desktop after making configuration changes

### Command Execution Fails

- Ensure you have the necessary permissions for the Kubernetes operations
- Verify the command syntax is correct
- Check if your Kubernetes cluster is accessible
