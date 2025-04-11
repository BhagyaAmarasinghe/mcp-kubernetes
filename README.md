# MCP Kubernetes

A Model Context Protocol (MCP) server for executing Kubernetes commands from Claude Desktop or any MCP-compatible client.

## Installation

You can install this tool directly using Go:

```bash
go install github.com/BhagyaAmarasinghe/mcp-kubernetes/cmd@latest
```

## Usage

### Configuration Options

When configuring Claude Desktop to use the MCP Kubernetes server, you can provide several command-line options:

* `-port`: Specify a port for the server (default: 3000, only relevant for direct usage)
* `-allowed-contexts`: Comma-separated list of allowed Kubernetes contexts (default: all contexts are allowed)
* `-namespace`: Default namespace for commands that don't specify one

Examples:

```bash
# Restrict to specific contexts
-allowed-contexts=minikube,docker-desktop

# Set default namespace
-namespace=default
```

### Using with Claude Desktop

1. Open your Claude Desktop App configuration at:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

2. Add the MCP Kubernetes server to your configuration:

```json
{
  "mcpServers": {
    "kubernetes": {
      "command": "mcp-kubernetes",
      "args": [
        "-allowed-contexts=minikube,docker-desktop",
        "-namespace=default"
      ]
    }
  }
}
```

3. You can now use Claude to execute Kubernetes commands, such as:
   - "Show me the pods in the default namespace"
   - "List all services across all namespaces"
   - "Check the status of my deployment named my-app"

### Available Commands

The MCP Kubernetes server supports the following MCP requests:

#### execute

Executes a kubectl command:

```json
{
  "command": "get pods -n default"
}
```

#### get-contexts

Retrieves a list of available Kubernetes contexts:

```json
{}
```

#### current-context

Gets the current Kubernetes context:

```json
{}
```

#### set-context

Sets the current Kubernetes context:

```json
{
  "context": "minikube"
}
```

## Security

This MCP server executes kubectl commands directly on your machine, so it should only be used in trusted environments. It does not implement authentication or authorization controls by default.

You can increase security by:

1. Using the `-allowed-contexts` flag to restrict which Kubernetes contexts can be used
2. Using the `-namespace` flag to set a default namespace for commands that don't specify one
3. Running the server with a user that has limited permissions

## Requirements

- Go 1.20 or higher
- kubectl installed and in your PATH
- A valid kubeconfig file

## License

MIT
