#!/bin/bash

# Cleanup script for Gemini CLI Manager test environment

echo "Cleaning up Gemini CLI Manager test environment..."

# Remove test directory
if [ -d ~/.gemini-test ]; then
    echo "Removing ~/.gemini-test directory..."
    rm -rf ~/.gemini-test
fi

# Remove any test archives created
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXAMPLES_DIR="$(dirname "$SCRIPT_DIR")"

cd "$EXAMPLES_DIR/extensions"
rm -f *.zip *.tar.gz

echo "Test environment cleaned up!"