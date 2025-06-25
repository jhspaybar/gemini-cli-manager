#!/bin/bash

# Test script to verify launch functionality

echo "Testing Gemini CLI Manager Launch"
echo "================================="
echo ""

# Clean up any previous debug log
rm -f /tmp/gemini-cli-manager-debug.log

# Create a simple mock gemini executable for testing
cat > /tmp/test-gemini << 'EOF'
#!/bin/bash
echo "Mock Gemini CLI started!"
echo "Profile: $GEMINI_PROFILE"
echo "Profile ID: $GEMINI_PROFILE_ID"
echo "Args: $@"
echo "Press any key to exit..."
read -n 1
EOF

chmod +x /tmp/test-gemini

# Set the test gemini path
export GEMINI_CLI_PATH=/tmp/test-gemini

echo "1. Created mock gemini executable at /tmp/test-gemini"
echo "2. Set GEMINI_CLI_PATH=/tmp/test-gemini"
echo ""
echo "To test the launch functionality:"
echo "  1. Run: ./gemini-cli-manager"
echo "  2. Press 'L' to launch"
echo "  3. The TUI should quit and exec into the mock gemini"
echo ""
echo "Debug logs will be written to: /tmp/gemini-cli-manager-debug.log"