#!/bin/bash

# Setup script for Gemini CLI Manager test environment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXAMPLES_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$(dirname "$EXAMPLES_DIR")")"

echo "Setting up Gemini CLI Manager test environment..."
echo "Project root: $PROJECT_ROOT"
echo "Examples dir: $EXAMPLES_DIR"

# Create test directories
echo "Creating test directories..."
mkdir -p ~/.gemini-cli-manager-test/extensions
mkdir -p ~/.gemini-cli-manager-test/profiles
mkdir -p ~/.gemini-test  # For the mock gemini

# Copy example extensions
echo "Installing example extensions..."
cp -r "$EXAMPLES_DIR/extensions/simple-extension" ~/.gemini-cli-manager-test/extensions/
cp -r "$EXAMPLES_DIR/extensions/mcp-extension" ~/.gemini-cli-manager-test/extensions/

# Make scripts executable
chmod +x ~/.gemini-cli-manager-test/extensions/mcp-extension/servers/echo-server.js 2>/dev/null || true

# Copy example profiles (convert JSON to YAML)
echo "Setting up example profiles..."
# For now, we'll create simple YAML profiles manually
cat > ~/.gemini-cli-manager-test/profiles/default.yaml << 'EOF'
id: default
name: Default Profile
description: Default profile for Gemini CLI
environment:
  GEMINI_ENV: default
extensions: []
created_at: 2024-01-01T00:00:00Z
updated_at: 2024-01-01T00:00:00Z
EOF

# Create a mock gemini CLI for testing
cat > ~/.gemini-test/mock-gemini << 'EOF'
#!/bin/bash
echo "Mock Gemini CLI v1.0.0"
echo "Profile: $GEMINI_PROFILE"
echo "Profile ID: $GEMINI_PROFILE_ID"
echo "Environment variables:"
env | grep GEMINI_ | sort
echo ""
echo "Arguments: $@"
echo ""
echo "Press Ctrl+C to exit..."

# Keep running to simulate a CLI
trap 'echo "Exiting..."; exit 0' INT TERM
while true; do
    sleep 1
done
EOF

chmod +x ~/.gemini-test/mock-gemini

# Create environment file
cat > ~/.gemini-cli-manager-test/test-env << EOF
# Test environment for Gemini CLI Manager
export GEMINI_CLI_PATH=~/.gemini-test/mock-gemini
export HOME_BACKUP=\$HOME
export HOME=~/.gemini-cli-manager-test

# Run the manager with test configuration
alias gemini-test='GEMINI_CLI_PATH=~/.gemini-test/mock-gemini HOME=~/.gemini-cli-manager-test go run $PROJECT_ROOT'
alias gemini-test-reset='export HOME=\$HOME_BACKUP'
EOF

echo ""
echo "Test environment setup complete!"
echo ""
echo "To use the test environment:"
echo "  source ~/.gemini-cli-manager-test/test-env"
echo "  gemini-test"
echo ""
echo "To restore normal HOME:"
echo "  gemini-test-reset"
echo ""
echo "Manager extensions in: ~/.gemini-cli-manager-test/extensions/"
echo "Manager profiles in: ~/.gemini-cli-manager-test/profiles/"
echo "Mock Gemini CLI at: ~/.gemini-test/mock-gemini"
echo ""
echo "Note: The manager will create symlinks in ~/.gemini/extensions/ when launching"