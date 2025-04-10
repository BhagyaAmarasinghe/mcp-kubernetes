#!/bin/bash

# Check if MCP Kubernetes is installed
if ! command -v mcp-kubernetes &> /dev/null; then
    echo "MCP Kubernetes is not installed. Installing now..."
    go install github.com/BhagyaAmarasinghe/mcp-kubernetes/cmd@latest
    if [ $? -ne 0 ]; then
        echo "Failed to install MCP Kubernetes"
        exit 1
    fi
    echo "MCP Kubernetes installed successfully"
else
    echo "MCP Kubernetes is already installed"
fi

# Set up as a LaunchAgent if on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    LAUNCH_AGENTS_DIR="$HOME/Library/LaunchAgents"
    PLIST_FILE="$LAUNCH_AGENTS_DIR/com.bhagya.mcp-kubernetes.plist"
    
    # Create LaunchAgents directory if it doesn't exist
    mkdir -p "$LAUNCH_AGENTS_DIR"
    
    # Copy the plist file
    cp "$(dirname "$0")/com.bhagya.mcp-kubernetes.plist" "$PLIST_FILE"
    
    # Update the path in the plist file if needed
    GO_BIN=$(which mcp-kubernetes)
    sed -i '' "s|/Users/bhagya/go/bin/mcp-kubernetes|$GO_BIN|g" "$PLIST_FILE"
    
    # Unload if it exists
    launchctl unload "$PLIST_FILE" 2>/dev/null
    
    # Load the LaunchAgent
    launchctl load "$PLIST_FILE"
    
    echo "LaunchAgent installed and loaded"
fi

# Start the server in the background
nohup mcp-kubernetes > /dev/null 2>&1 &

# Wait for server to start
echo "Waiting for MCP Kubernetes server to start..."
sleep 2

# Check if server is running
if curl -s http://localhost:3000/health > /dev/null; then
    echo "MCP Kubernetes server is running"
else
    echo "Failed to start MCP Kubernetes server"
    exit 1
fi

echo ""
echo "======================================================"
echo "MCP Kubernetes is now running at ws://localhost:3000/ws"
echo ""
echo "To use with Claude Desktop:"
echo "1. Open Claude Desktop"
echo "2. Go to Settings > Tools"
echo "3. Add a new MCP Tool with the URL: ws://localhost:3000/ws"
echo "4. Name it 'Kubernetes'"
echo "5. Save the settings"
echo ""
echo "You can now use Claude to execute Kubernetes commands!"
echo "======================================================"
