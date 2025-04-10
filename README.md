# MCP Kubernetes

A Model Context Protocol (MCP) server for executing Kubernetes commands from Claude Desktop or any MCP-compatible client.

## Installation

You can install this tool directly using Go:

```bash
go install github.com/BhagyaAmarasinghe/mcp-kubernetes/cmd@latest
```

If you encounter issues with the above command, you can build and install locally:

```bash
git clone https://github.com/BhagyaAmarasinghe/mcp-kubernetes.git
cd mcp-kubernetes
chmod +x build.sh
./build.sh
```

## Usage

### Starting the server

Run the MCP Kubernetes server:

```bash
mcp-kubernetes
```

By default, the server runs on port 3000. You can specify a different port using the `-port` flag:

```bash
mcp-kubernetes -port 8080
```

### Using with Claude Desktop

1. Configure Claude Desktop to use the MCP Kubernetes server at `ws://localhost:3000/ws`

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

## Dependencies

- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - Model Context Protocol implementation
- [github.com/gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [github.com/google/uuid](https://github.com/google/uuid) - UUID generation

## Security

This MCP server executes kubectl commands directly on your machine, so it should only be used in trusted environments. It does not implement authentication or authorization controls by default.

## Requirements

- Go 1.20 or higher
- kubectl installed and in your PATH
- A valid kubeconfig file

## License

MIT
