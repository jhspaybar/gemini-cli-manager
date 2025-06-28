# Screenshot Workflow Test
# Captures screenshots at key moments to verify the workflow

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
Type "./gemini-cli-manager --state-dir /tmp/vhs-screenshot-$(date +%s)"
Enter
Show

# Start recording from here
Output output/e2e-screenshot-workflow.gif

# Let the app fully load and show initial state
Sleep 3s
Screenshot output/e2e-ss-01-empty-extensions.png

# Install extension
Type "n"
Sleep 1s
Type "/Users/williamthurston/repos/gemini-cli-manager/examples/extensions/filesystem-enhanced"
Enter
Sleep 3s
Screenshot output/e2e-ss-02-extension-installed.png

# Dismiss modal
Escape
Sleep 500ms

# Go to profiles
Tab
Sleep 1s
Screenshot output/e2e-ss-03-empty-profiles.png

# Create profile with extension
Type "n"
Sleep 1s
Type "Demo Profile"
Tab
Type "Profile with filesystem extension"
Tab
Sleep 500ms
Space  # Select the extension
Sleep 500ms
Screenshot output/e2e-ss-04-profile-form.png

# Save profile
Ctrl+S
Sleep 2s
Escape
Sleep 1s
Screenshot output/e2e-ss-05-profile-created.png

# Activate the profile
Enter
Sleep 2s
Screenshot output/e2e-ss-06-profile-activated.png

# Check status bar shows active profile
Sleep 1s

# Launch Gemini
Type "L"
Sleep 3s
Screenshot output/e2e-ss-07-launching.png
Sleep 3s
Screenshot output/e2e-ss-08-gemini-running.png

# Exit Gemini
Type "/quit"

# Stop recording - cleanup happens after
Hide
Enter
Sleep 2s
Screenshot output/e2e-ss-09-back-to-manager.png
Type "q"
Sleep 500ms
Type "rm -rf /tmp/vhs-screenshot-*"
Enter