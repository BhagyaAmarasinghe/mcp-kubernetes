# Autostart Scripts

This directory contains scripts for automatically starting the MCP Kubernetes server.

## macOS LaunchAgent

The file `com.bhagya.mcp-kubernetes.plist` is a LaunchAgent configuration that will start the MCP Kubernetes server automatically when you log in.

### Installation

1. First ensure that the MCP Kubernetes server is installed:
   ```bash
   go install github.com/bhagya/mcp-kubernetes/cmd@latest
   ```

2. Copy the LaunchAgent file to your LaunchAgents directory:
   ```bash
   cp com.bhagya.mcp-kubernetes.plist ~/Library/LaunchAgents/
   ```

3. Load the LaunchAgent:
   ```bash
   launchctl load ~/Library/LaunchAgents/com.bhagya.mcp-kubernetes.plist
   ```

4. The server should now be running. You can verify with:
   ```bash
   curl http://localhost:3000/health
   ```

### Manual Control

- To start the service:
  ```bash
  launchctl start com.bhagya.mcp-kubernetes
  ```

- To stop the service:
  ```bash
  launchctl stop com.bhagya.mcp-kubernetes
  ```

- To remove the service:
  ```bash
  launchctl unload ~/Library/LaunchAgents/com.bhagya.mcp-kubernetes.plist
  rm ~/Library/LaunchAgents/com.bhagya.mcp-kubernetes.plist
  ```

### Logs

The standard output and error logs are saved to:
- `/Users/bhagya/Library/Logs/mcp-kubernetes.log`
- `/Users/bhagya/Library/Logs/mcp-kubernetes.err`

## Linux Systemd Service

For Linux systems using systemd, you can create a user service:

1. Create the file `~/.config/systemd/user/mcp-kubernetes.service` with the following content:
   ```
   [Unit]
   Description=MCP Kubernetes Server
   After=network.target

   [Service]
   ExecStart=/home/username/go/bin/mcp-kubernetes
   Restart=on-failure

   [Install]
   WantedBy=default.target
   ```

2. Enable and start the service:
   ```bash
   systemctl --user enable mcp-kubernetes
   systemctl --user start mcp-kubernetes
   ```

## Windows Startup

For Windows:

1. Create a shortcut to the `mcp-kubernetes.exe` file
2. Press Win+R and type `shell:startup` to open the Startup folder
3. Copy the shortcut to this folder
