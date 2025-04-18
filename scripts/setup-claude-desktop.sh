#!/bin/bash

# Check if MCP Kubernetes is installed
if ! command -v mcp-kubernetes &> /dev/null; then
    echo "MCP Kubernetes is not installed. Installing now..."
    go install github.com/BhagyaAmarasinghe/mcp-kubernetes@latest
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
    cat > "$PLIST_FILE" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.bhagya.mcp-kubernetes</string>
    <key>ProgramArguments</key>
    <array>
        <string>$(which mcp-kubernetes)</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>$HOME/Library/Logs/mcp-kubernetes.err</string>
    <key>StandardOutPath</key>
    <string>$HOME/Library/Logs/mcp-kubernetes.log</string>
</dict>
</plist>
EOF
    
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

echo ""
echo "======================================================"
echo "MCP Kubernetes is now running"
echo ""
echo "To use with Claude Desktop:"
echo "1. Open Claude Desktop"
echo "2. Go to Settings > Developer"
echo "3. Edit Config and add the following:"
echo ""
echo '{
  "mcpServers": {
    "kubernetes": {
      "command": "mcp-kubernetes"
    }
  }
}'
echo ""
echo "4. Save the settings and restart Claude Desktop"
echo ""
echo "You can now use Claude to execute Kubernetes commands!"
echo "======================================================"
