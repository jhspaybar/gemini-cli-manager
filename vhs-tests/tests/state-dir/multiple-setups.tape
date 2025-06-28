# Multiple State Directories Test
# Tests running multiple independent setups

Output output/state-dir-multiple.gif

# VHS Configuration
Set FontSize 14
Set Width 1400
Set Height 1000
Set Theme "Builtin Solarized Light"
Set TypingSpeed 50ms

# Start from the project root directory
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

# Show the help to demonstrate the feature
Type "go run . --help"
Enter
Sleep 3s
Screenshot output/state-01-help.png

Type "clear"
Enter
Sleep 500ms

# First setup - Work profile
Type "# Setup 1: Work Environment"
Enter
Type "go run . --state-dir ~/vhs-work-setup"
Enter
Sleep 3s
Screenshot output/state-02-work-start.png

# Navigate to settings to show state directory
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "Enter"
Sleep 1s
Screenshot output/state-03-work-settings.png

# Create a profile
Type "Tab"
Sleep 300ms
Type "k"
Sleep 300ms
Type "Enter"
Sleep 1s
Type "Tab"
Sleep 300ms
Type "n"
Sleep 1s
Screenshot output/state-04-work-new-profile.png

Type "Work Profile"
Sleep 500ms
Type "Tab"
Sleep 300ms
Type "Development setup for work projects"
Sleep 500ms
Screenshot output/state-05-work-profile-filled.png

Type "Ctrl+S"
Sleep 1s
Type "Escape"
Sleep 1s
Screenshot output/state-06-work-profile-created.png

# Exit
Type "q"
Sleep 1s

Type "clear"
Enter
Sleep 500ms

# Second setup - Personal profile
Type "# Setup 2: Personal Environment"
Enter
Type "go run . --state-dir ~/vhs-personal-setup"
Enter
Sleep 3s
Screenshot output/state-07-personal-start.png

# Show it's a fresh setup
Type "j"
Sleep 300ms
Type "Enter"
Sleep 1s
Screenshot output/state-08-personal-empty-profiles.png

# Navigate to settings to show different state directory
Type "j"
Sleep 300ms
Type "Enter"
Sleep 1s
Screenshot output/state-09-personal-settings.png

# Exit
Type "q"
Sleep 1s

Type "clear"
Enter
Sleep 500ms

# Show both directories exist independently
Type "# Verify independent state directories:"
Enter
Type "ls -la ~/vhs-work-setup/"
Enter
Sleep 2s
Screenshot output/state-10-work-dir.png

Type "ls -la ~/vhs-personal-setup/"
Enter
Sleep 2s
Screenshot output/state-11-personal-dir.png

# Cleanup
Type "rm -rf ~/vhs-work-setup ~/vhs-personal-setup"
Enter
Sleep 500ms