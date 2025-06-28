# Full Workflow Test
# End-to-end test: Install extension, create profile, activate it, and launch Gemini

# VHS Configuration
Set FontSize 14
Set Width 1400
Set Height 1000
Set Theme "Builtin Solarized Light"
Set TypingSpeed 75ms

# Setup before recording (hidden)
Hide
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s
Type "./gemini-cli-manager --state-dir /tmp/vhs-e2e-full-$(date +%s)"
Enter
Show

# Start recording from here
Output output/e2e-full-workflow.gif

# Let the app fully load and show initial state
Sleep 3s
Screenshot output/e2e-01-start.png

# Step 1: Install an extension from the filesystem
Type "n"
Sleep 1s
Screenshot output/e2e-02-install-modal.png

# Enter the path to the example extension
Type "/Users/williamthurston/repos/gemini-cli-manager/examples/extensions/filesystem-enhanced"
Enter
Sleep 2s
Screenshot output/e2e-03-extension-installed.png

# Press Escape to dismiss any modal/message
Type "Escape"
Sleep 1s
Screenshot output/e2e-04-extensions-list.png

# Step 2: Navigate to Profiles
Tab
Sleep 1s
Screenshot output/e2e-05-profiles-view.png

# Step 3: Create a new profile
Type "n"
Sleep 1s
Screenshot output/e2e-06-new-profile-form.png

# Fill in profile details
Type "Development"
Tab
Sleep 500ms
Type "My development environment with filesystem tools"
Tab
Sleep 500ms
Screenshot output/e2e-07-profile-filled.png

# Save the profile
Ctrl+S
Sleep 2s
Screenshot output/e2e-08-profile-created.png

# Escape from form
Type "Escape"
Sleep 1s

# Step 4: Edit profile to add the extension
Type "e"
Sleep 1s
Screenshot output/e2e-09-edit-profile.png

# Navigate to extensions section (Tab to skip name/description)
Tab
Tab
Sleep 500ms

# Select the filesystem-enhanced extension
Space  # Toggle checkbox
Sleep 500ms
Screenshot output/e2e-10-extension-selected.png

# Save the profile
Ctrl+S
Sleep 2s

# Escape from edit form
Type "Escape"
Sleep 1s
Screenshot output/e2e-11-profile-with-extension.png

# Step 5: Activate the profile
Enter
Sleep 2s
Screenshot output/e2e-12-profile-activated.png

# Step 6: Launch Gemini
Type "L"
Sleep 3s
Screenshot output/e2e-13-launching-gemini.png

# Wait for Gemini to start (this might take a moment)
Sleep 5s
Screenshot output/e2e-14-gemini-launched.png

# Step 7: Exit from Gemini
Type "/quit"

# Stop recording - cleanup happens after
Hide
Enter
Sleep 2s
Screenshot output/e2e-15-back-to-manager.png
Type "q"
Sleep 1s
Type "rm -rf /tmp/vhs-e2e-full-*"
Enter
Sleep 500ms