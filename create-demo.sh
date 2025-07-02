#!/bin/bash
# Create a demo GIF using VHS (Video Home System)
# This script demonstrates the complete workflow of the Gemini CLI Manager

# First, clean up any existing data
echo "Cleaning up existing data..."
rm -rf ~/Library/Application\ Support/gemini-cli-manager

# Build and start the app first (outside of recording)
echo "Building the application..."
cargo build --release

# Now create and run the VHS script
cat > demo.tape << 'EOF'
Output demo.gif
Set FontSize 18
Set Width 1200
Set Height 1000
Set Theme "Catppuccin Mocha"

# Start the app (already built)
Type "./target/release/gemini-cli-manager"
Enter
Sleep 3s

# Press 'i' to import
Type "i"
Sleep 2s

# Navigate down 8 times to reach starter-extensions/
Down
Down
Down
Down
Down
Down
Down
Down
Sleep 1s

# Enter starter-extensions directory
Enter
Sleep 1s

# Navigate to rust-tools/ (3 downs from top)
Down
Down
Down
Sleep 1s

# Enter rust-tools directory
Enter
Sleep 1s

# Select gemini-extension.json
Enter
Sleep 3s

# Pause to show the extension details
Sleep 3s

# Press 'b' to go back to main view
Type "b"
Sleep 1s

# Now Tab to go to Profiles tab
Tab
Sleep 1s

# Create new profile with 'n'
Type "n"
Sleep 1s

# Type profile name
Type "Rust Development"
Sleep 500ms

# Tab to description
Tab
Type "Profile for Rust development with essential tools"
Sleep 500ms

# Tab to working directory
Tab
Type "~/rust-projects"
Sleep 500ms

# Tab to extensions selection
Tab

# Select the rust-tools extension (space to toggle)
Type " "
Sleep 500ms

# Tab to tags field
Tab

# Add a tag
Type "rust, development"
Sleep 500ms

# Save profile with Ctrl+S
Ctrl+S
Sleep 2s

# Pause to show the saved profile in the list
Sleep 2s

# Select the profile and view details with Enter
Enter
Sleep 1s

# Pause to show the profile details
Sleep 3s

# Launch Gemini with 'l'
Type "l"
Sleep 3s

# Quit Gemini immediately
Type "/quit"
Sleep 2s
Enter
Sleep 3s

# Back in the CLI manager - pause to show we're back
Sleep 2s

# Quit CLI manager
Type "q"
Sleep 1s
EOF

echo "Recording demo..."
vhs demo.tape

echo "Demo created: demo.gif"
rm demo.tape