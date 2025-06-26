#!/bin/bash

# Script to clean up old Gemini CLI Manager state from ~/.gemini
# This removes only the manager-specific files, preserving Gemini's own files

echo "Cleaning up old Gemini CLI Manager state from ~/.gemini..."
echo ""

# Check if ~/.gemini exists
if [ ! -d ~/.gemini ]; then
    echo "~/.gemini directory not found. Nothing to clean up."
    exit 0
fi

# Clean up profiles directory (this was created by the manager)
if [ -d ~/.gemini/profiles ]; then
    echo "Found old profiles directory at ~/.gemini/profiles"
    echo "Contents:"
    ls -la ~/.gemini/profiles/
    read -p "Remove this directory? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf ~/.gemini/profiles
        echo "✓ Removed ~/.gemini/profiles"
    else
        echo "⚠ Skipped ~/.gemini/profiles"
    fi
fi

# Clean up extension symlinks in ~/.gemini/extensions
if [ -d ~/.gemini/extensions ]; then
    echo ""
    echo "Checking ~/.gemini/extensions for manager-created symlinks..."
    
    # Find all symlinks that point to ~/.gemini-cli-manager
    found_symlinks=false
    for item in ~/.gemini/extensions/*; do
        if [ -L "$item" ]; then
            target=$(readlink "$item")
            if [[ "$target" == *".gemini-cli-manager"* ]]; then
                if [ "$found_symlinks" = false ]; then
                    echo "Found manager-created symlinks:"
                    found_symlinks=true
                fi
                echo "  $(basename "$item") -> $target"
            fi
        fi
    done
    
    if [ "$found_symlinks" = true ]; then
        read -p "Remove these symlinks? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            # Remove only symlinks pointing to our manager directory
            for item in ~/.gemini/extensions/*; do
                if [ -L "$item" ]; then
                    target=$(readlink "$item")
                    if [[ "$target" == *".gemini-cli-manager"* ]]; then
                        rm "$item"
                        echo "✓ Removed symlink: $(basename "$item")"
                    fi
                fi
            done
            
            # Remove extensions directory if it's now empty
            if [ -z "$(ls -A ~/.gemini/extensions)" ]; then
                rmdir ~/.gemini/extensions
                echo "✓ Removed empty ~/.gemini/extensions directory"
            fi
        else
            echo "⚠ Skipped extension symlinks"
        fi
    else
        # Check if extensions directory is empty
        if [ -z "$(ls -A ~/.gemini/extensions)" ]; then
            read -p "~/.gemini/extensions is empty. Remove it? (y/n) " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                rmdir ~/.gemini/extensions
                echo "✓ Removed empty ~/.gemini/extensions directory"
            fi
        else
            echo "~/.gemini/extensions contains non-manager files. Leaving as is."
        fi
    fi
fi

# Check for any other manager-specific files
echo ""
echo "Checking for other manager-specific files..."

# List files we know belong to Gemini (not the manager)
gemini_files=(
    "GEMINI.md"
    "oauth_creds.json"
    "settings.json"
    "user_id"
    "active_profile"
    "tmp"
)

# Show what's left
echo ""
echo "Remaining contents of ~/.gemini:"
ls -la ~/.gemini/
echo ""
echo "The following files belong to Gemini itself and were preserved:"
for file in "${gemini_files[@]}"; do
    if [ -e ~/.gemini/"$file" ]; then
        echo "  ✓ $file"
    fi
done

echo ""
echo "Cleanup complete!"
echo ""
echo "Your Gemini CLI Manager state is now at: ~/.gemini-cli-manager/"
ls -la ~/.gemini-cli-manager/ 2>/dev/null || echo "(Not yet created - will be created when you run the manager)"