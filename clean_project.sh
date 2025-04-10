#!/bin/bash

# This script removes unnecessary files from the project
echo "Cleaning up MCP Kubernetes project..."

# Remove the to_delete directory if it exists
if [ -d "to_delete" ]; then
  rm -rf to_delete
  echo "Removed to_delete directory"
fi

# Remove the fix_dependencies.sh script if it exists
if [ -f "fix_dependencies.sh" ]; then
  rm -f fix_dependencies.sh
  echo "Removed fix_dependencies.sh"
fi

# Remove the debug_module.sh script if it exists
if [ -f "debug_module.sh" ]; then
  rm -f debug_module.sh
  echo "Removed debug_module.sh"
fi

echo "Cleanup complete!"
echo "You can now run: ./build.sh"
