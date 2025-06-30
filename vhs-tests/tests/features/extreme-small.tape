# Test in extremely small terminal
# Verifies UI remains functional in very constrained space

Set FontSize 14
Set Width 600
Set Height 250  # Extremely small
Set Theme "Builtin Solarized Light"

Output output/extreme-small.gif

Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

Type "./gemini-cli-manager --state-dir /tmp/vhs-extreme-$(date +%s)"
Enter
Sleep 3

# Initial view
Screenshot output/extreme-01-initial.png

# Navigate to settings which has many items
Tab
Tab
Sleep 1s
Screenshot output/extreme-02-settings.png

# Try scrolling down
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/extreme-03-scrolled.png

# Go to help
Tab
Sleep 1s
Screenshot output/extreme-04-help.png

# Scroll in help
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/extreme-05-help-scrolled.png

# Exit
Type "q"
Sleep 500ms