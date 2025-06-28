# Simple Workflow Test
# Minimal end-to-end test with better error handling

# VHS Configuration
Set FontSize 16
Set Width 1200
Set Height 1000
Set Theme "Builtin Solarized Light"
Set TypingSpeed 100ms

# Setup before recording (hidden)
Hide
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s
Type "./gemini-cli-manager --state-dir /tmp/vhs-simple-$(date +%s)"
Enter
Show

# Start recording from here
Output output/e2e-simple-workflow.gif

# Let the app fully load and show initial state
Sleep 3s

# Install extension
Type "n"
Sleep 1s
Type "/Users/williamthurston/repos/gemini-cli-manager/examples/extensions/filesystem-enhanced"
Enter
Sleep 3s

# Dismiss any modal with multiple escapes
Escape
Sleep 500ms
Escape
Sleep 500ms

# Go to profiles
Tab
Sleep 1s

# Create profile
Type "n"
Sleep 1s
Type "Test Profile"
Tab
Type "Simple test profile"
Tab
Sleep 500ms

# Now check the filesystem-enhanced extension (should be first/only one)
Space
Sleep 500ms

# Save profile
Ctrl+S
Sleep 2s

# Back to profiles list
Escape
Sleep 1s

# Activate the profile (should be selected already)
Enter
Sleep 2s

# Launch Gemini with activated profile
Type "L"
Sleep 5s

# Exit Gemini
Type "/quit"

# Stop recording - cleanup happens after
Hide
Enter
Sleep 2s
Type "q"
Sleep 500ms
Type "rm -rf /tmp/vhs-simple-*"
Enter