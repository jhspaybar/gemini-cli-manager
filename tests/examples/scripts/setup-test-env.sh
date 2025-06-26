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
mkdir -p ~/.gemini-test/extensions
mkdir -p ~/.gemini-test/profiles
mkdir -p ~/.gemini-test/config

# Copy example extensions
echo "Installing example extensions..."
cp -r "$EXAMPLES_DIR/extensions/simple-extension" ~/.gemini-test/extensions/
cp -r "$EXAMPLES_DIR/extensions/mcp-extension" ~/.gemini-test/extensions/

# Make scripts executable
chmod +x ~/.gemini-test/extensions/mcp-extension/servers/echo-server.js 2>/dev/null || true

# Copy example profiles
echo "Setting up example profiles..."
cp "$EXAMPLES_DIR/profiles"/*.json ~/.gemini-test/profiles/

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
cat > ~/.gemini-test/test-env << EOF
# Test environment for Gemini CLI Manager
export GEMINI_CLI_PATH=~/.gemini-test/mock-gemini
export GEMINI_BASE_PATH=~/.gemini-test

# Run the manager with test configuration
alias gemini-test='GEMINI_CLI_PATH=~/.gemini-test/mock-gemini go run $PROJECT_ROOT'
EOF

echo ""
echo "Test environment setup complete!"
echo ""
echo "To use the test environment:"
echo "  source ~/.gemini-test/test-env"
echo "  gemini-test"
echo ""
echo "Test extensions installed in: ~/.gemini-test/extensions/"
echo "Test profiles installed in: ~/.gemini-test/profiles/"
echo "Mock Gemini CLI at: ~/.gemini-test/mock-gemini"