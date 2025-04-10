# MCP Kubernetes

A Model Context Protocol (MCP) server for executing Kubernetes commands from Claude Desktop or any MCP-compatible client.

## Installation

You can install this tool directly using Go:

```bash
go install github.com/BhagyaAmarasinghe/mcp-kubernetes/cmd@latest
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

### Auto-starting the server

To automatically start the server when you log in:

#### macOS

Use the provided setup script:

```bash
chmod +x scripts/setup-claude-desktop.sh
./scripts/setup-claude-desktop.sh
```

Or manually set up a LaunchAgent:

```bash
cp scripts/com.bhagya.mcp-kubernetes.plist ~/Library/LaunchAgents/
launchctl load ~/Library/LaunchAgents/com.bhagya.mcp-kubernetes.plist
```

#### Linux (systemd)

Create a user service:

```bash
mkdir -p ~/.config/systemd/user/
cp scripts/mcp-kubernetes.service ~/.config/systemd/user/
systemctl --user enable mcp-kubernetes
systemctl --user start mcp-kubernetes
```

#### Windows

Add a shortcut to the executable in your startup folder:

```
Win+R → shell:startup → Create shortcut to mcp-kubernetes.exe
```

See the scripts/README.md file for more detailed instructions.

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

## Security

This MCP server executes kubectl commands directly on your machine, so it should only be used in trusted environments. It does not implement authentication or authorization controls by default.

## Requirements

- Go 1.18 or higher
- kubectl installed and in your PATH
- A valid kubeconfig file

## License

MIT
