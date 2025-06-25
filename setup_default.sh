#!/bin/bash
# Quick setup script for default profile and extension

echo "Setting up default Gemini CLI Manager configuration..."

# Create directories
mkdir -p ~/.gemini/profiles
mkdir -p ~/.gemini/extensions/example-extension

# Create default profile
cat > ~/.gemini/profiles/default.yaml << 'EOF'
id: default
name: Default
description: Default development profile
extensions:
  - example-extension
environment:
  GEMINI_ENV: development
mcp_servers: {}
created_at: 2025-06-25T11:02:24.591403-07:00
updated_at: 2025-06-25T11:02:24.591406-07:00
usage_count: 0
EOF

# Create example extension
cat > ~/.gemini/extensions/example-extension/gemini-extension.json << 'EOF'
{
  "name": "example-extension",
  "displayName": "Example Extension",
  "version": "1.0.0",
  "description": "A sample extension for testing",
  "mcp": {
    "servers": {
      "example-server": {
        "command": "echo",
        "args": ["MCP server would start here"],
        "env": {
          "NODE_ENV": "production"
        }
      }
    }
  }
}
EOF

# Create GEMINI.md for the extension
cat > ~/.gemini/extensions/example-extension/GEMINI.md << 'EOF'
# Example Extension

This is a sample extension for testing the Gemini CLI Manager.

## Features
- Example MCP server configuration
- Testing capabilities
EOF

echo "✅ Setup complete!"
echo ""
echo "You now have:"
echo "  - A default profile at ~/.gemini/profiles/default.yaml"
echo "  - An example extension at ~/.gemini/extensions/example-extension/"
echo ""
echo "Run the Gemini CLI Manager with:"
echo "  ./gemini-cli-manager"