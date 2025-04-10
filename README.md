# Kubernetes MCP Server

A Model Context Protocol (MCP) server for executing Kubernetes commands from Claude Desktop or any MCP-compatible client.

## Features
- Execute kubectl commands securely
- List available Kubernetes contexts
- Switch between contexts
- Get current context
- View command history
- Safety features to limit allowed commands
- Configure allowed command set for security

**Note**: The server will only allow execution of kubectl commands specified via the allowedCommands parameter or all commands if configured with "*".

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

By default, the server runs on port 3000 and allows all kubectl commands. You can specify a different port and restrict allowed commands using flags:

```bash
# Allow only specific commands
mcp-kubernetes --allowed-commands=get,describe,config --port=3000
```

### Available MCP Commands

The MCP Kubernetes server supports the following MCP requests:

#### execute

Executes a kubectl command:

```json
{
  "command": "get pods -n default"
}
```

Response:
```json
{
  "success": true,
  "output": "NAME                     READY   STATUS    RESTARTS   AGE\n...",
  "execution_time": "152.4ms"
}
```

#### get-contexts

Retrieves a list of available Kubernetes contexts:

```json
{}
```

Response:
```json
{
  "contexts": ["minikube", "docker-desktop", "my-cluster"]
}
```

#### current-context

Gets the current Kubernetes context:

```json
{}
```

Response:
```json
{
  "context": "minikube"
}
```

#### set-context

Sets the current Kubernetes context:

```json
{
  "context": "docker-desktop"
}
```

Response:
```json
{
  "success": true
}
```

#### list-recent-commands

Lists recently executed commands:

```json
{
  "limit": 5
}
```

Response:
```json
{
  "commands": [
    {
      "command": "get pods",
      "timestamp": "2023-04-15T14:32:15Z",
      "success": true
    },
    {
      "command": "get nodes",
      "timestamp": "2023-04-15T14:30:05Z",
      "success": true
    }
  ]
}
```

#### list-allowed-commands

Lists all commands that the server is allowed to execute:

```json
{}
```

Response:
```json
{
  "allowed_commands": ["get", "describe", "config"]
}
```

Or if all commands are allowed:
```json
{
  "allowed_commands": "*"
}
```

## Usage with Claude Desktop

Add this to your `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "kubernetes": {
      "command": "mcp-kubernetes",
      "args": [
        "--allowed-commands=get,describe,config",
        "--port=3000"
      ]
    }
  }
}
```

To allow all kubectl commands (use with caution):
```json
{
  "mcpServers": {
    "kubernetes": {
      "command": "mcp-kubernetes",
      "args": [
        "--allowed-commands=*"
      ]
    }
  }
}
```

## Security Considerations

When using this MCP server, please consider:
1. Only allow kubectl commands you trust - a restrictive allowlist is recommended
2. Avoid allowing commands that could modify cluster settings or access sensitive data
3. The server runs with the permissions of the user running Claude Desktop
4. Command output is sent back to the LLM, so be mindful of sensitive information

## Requirements

- Go 1.20 or higher
- kubectl installed and in your PATH
- A valid kubeconfig file

## License

MIT
