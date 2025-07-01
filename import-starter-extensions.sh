#!/bin/bash

# Import starter extensions into Gemini CLI Manager
# This script imports the provided starter extensions to your local storage

set -e

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
STARTER_DIR="$SCRIPT_DIR/starter-extensions"

# Determine the storage directory
if [[ -d "$HOME/Library/Application Support/gemini-cli-manager" ]]; then
    # macOS
    STORAGE_DIR="$HOME/Library/Application Support/gemini-cli-manager"
elif [[ -d "$HOME/.local/share/gemini-cli-manager" ]]; then
    # Linux
    STORAGE_DIR="$HOME/.local/share/gemini-cli-manager"
elif [[ -d "$APPDATA/gemini-cli-manager" ]]; then
    # Windows
    STORAGE_DIR="$APPDATA/gemini-cli-manager"
else
    echo "Error: Could not find Gemini CLI Manager storage directory"
    echo "Please run the Gemini CLI Manager at least once to create the storage directory"
    exit 1
fi

EXTENSIONS_DIR="$STORAGE_DIR/extensions"

# Create extensions directory if it doesn't exist
mkdir -p "$EXTENSIONS_DIR"

echo "Importing starter extensions to: $EXTENSIONS_DIR"
echo

# Import each extension
for ext_dir in "$STARTER_DIR"/*; do
    if [[ -d "$ext_dir" ]]; then
        ext_name=$(basename "$ext_dir")
        ext_json="$ext_dir/extension.json"
        
        if [[ -f "$ext_json" ]]; then
            echo "Importing $ext_name..."
            
            # Read the extension ID from the JSON file
            ext_id=$(grep -o '"id"[[:space:]]*:[[:space:]]*"[^"]*"' "$ext_json" | sed 's/.*: *"\([^"]*\)"/\1/')
            
            if [[ -n "$ext_id" ]]; then
                # Copy the extension.json file to storage
                cp "$ext_json" "$EXTENSIONS_DIR/${ext_id}.json"
                echo "  ✓ Imported $ext_name (ID: $ext_id)"
            else
                echo "  ✗ Error: Could not read extension ID from $ext_json"
            fi
        else
            echo "  ✗ Skipping $ext_name - no extension.json found"
        fi
    fi
done

echo
echo "Import complete!"
echo
echo "Available extensions:"
for ext_file in "$EXTENSIONS_DIR"/*.json; do
    if [[ -f "$ext_file" ]]; then
        ext_name=$(grep -o '"name"[[:space:]]*:[[:space:]]*"[^"]*"' "$ext_file" | sed 's/.*: *"\([^"]*\)"/\1/')
        echo "  - $ext_name"
    fi
done

echo
echo "You can now:"
echo "1. Run './gemini-cli-manager' to launch the TUI"
echo "2. View the imported extensions in the Extensions tab"
echo "3. Create profiles that use these extensions"