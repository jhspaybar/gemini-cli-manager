# Quick Workflow Test
# Streamlined end-to-end test focusing on key interactions

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
Type "./gemini-cli-manager --state-dir /tmp/vhs-e2e-quick-$(date +%s)"
Enter
Show

# Start recording from here
Output output/e2e-quick-workflow.gif

# Let the app fully load and show initial state
Sleep 3s

# Install extension
Type "n"
Sleep 1s
Type "/Users/williamthurston/repos/gemini-cli-manager/examples/extensions/filesystem-enhanced"
Enter
Sleep 2s
Type "Escape"
Sleep 500ms

# Go to profiles and create one
Tab
Sleep 500ms
Type "n"
Sleep 500ms
Type "Quick Demo"
Tab
Type "Demo profile for testing"
Tab
Sleep 500ms
Ctrl+S
Sleep 1s
Type "Escape"
Sleep 500ms

# Edit to add extension
Type "e"
Sleep 500ms
Tab
Tab
Space
Sleep 500ms
Ctrl+S
Sleep 1s
Type "Escape"
Sleep 500ms

# Activate profile
Enter
Sleep 1s

# Launch Gemini
Type "L"
Sleep 5s

# Quick interaction in Gemini and exit
Type "/quit"

# Stop recording - cleanup happens after
Hide
Enter
Sleep 2s
Type "q"
Sleep 500ms
Type "rm -rf /tmp/vhs-e2e-quick-*"
Enter