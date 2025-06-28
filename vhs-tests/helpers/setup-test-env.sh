#!/bin/bash
# Setup script for VHS tests

set -e

echo "Setting up VHS test environment..."

# Get the project root (parent of vhs-tests)
PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
VHS_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Check if gemini-cli-manager binary exists
if [ ! -f "$PROJECT_ROOT/gemini-cli-manager" ]; then
    echo "Building gemini-cli-manager..."
    cd "$PROJECT_ROOT"
    go build -o gemini-cli-manager
    cd "$VHS_ROOT"
else
    echo "✓ gemini-cli-manager binary found"
fi

# Create necessary directories
mkdir -p "$VHS_ROOT/output"
mkdir -p "$VHS_ROOT/temp"
mkdir -p "$VHS_ROOT/tests/navigation"
mkdir -p "$VHS_ROOT/tests/extensions"
mkdir -p "$VHS_ROOT/tests/profiles"
mkdir -p "$VHS_ROOT/tests/state-dir"

# Check if VHS is installed
if ! command -v vhs &> /dev/null; then
    echo "❌ VHS is not installed!"
    echo "Install with: brew install vhs"
    echo "Or see: https://github.com/charmbracelet/vhs#installation"
    exit 1
else
    echo "✓ VHS is installed: $(vhs --version)"
fi

# Create a gitignore for output
if [ ! -f "$VHS_ROOT/.gitignore" ]; then
    cat > "$VHS_ROOT/.gitignore" << EOF
# Generated files
output/
temp/
*.gif
*.png
*.mp4
*.webm

# Test state directories
/test-state-*/
EOF
    echo "✓ Created .gitignore"
fi

echo "✅ VHS test environment ready!"
echo ""
echo "Next steps:"
echo "1. Create test files in tests/ subdirectories"
echo "2. Run tests with: make test-all"
echo "3. View results in output/ directory"