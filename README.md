# MCP Kubernetes

A Model Context Protocol (MCP) server for executing Kubernetes commands from Claude Desktop or any MCP-compatible client.

## Installation

You can install this tool directly using Go:

```bash
go install github.com/BhagyaAmarasinghe/mcp-kubernetes@latest
```

## Usage

### Starting the server

Run the MCP Kubernetes server:

```bash
mcp-kubernetes
```

By default, the server runs using stdio transport. You can specify a different port for HTTP transport with the `-port` flag:

```bash
mcp-kubernetes -port 8080
```

### Using with Claude Desktop

1. Configure Claude Desktop to use the MCP Kubernetes server by adding it to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "kubernetes": {
      "command": "mcp-kubernetes"
    }
  }
}
```

2. You can now use Claude to execute Kubernetes commands, such as:
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

## Requirements

- Go 1.23 or higher
- kubectl installed and in your PATH
- A valid kubeconfig file

## License

MIT
