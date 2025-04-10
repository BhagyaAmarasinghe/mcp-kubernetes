#!/bin/bash

# This script removes unnecessary files from the project
echo "Cleaning up MCP Kubernetes project..."

# Remove the to_delete directory
rm -rf to_delete

# Remove this script itself
rm -f clean_project.sh

echo "Cleanup complete!"
