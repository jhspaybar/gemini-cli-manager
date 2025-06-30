# Test in medium-sized terminal
# Verifies scrolling works in a typical small terminal

Set FontSize 14
Set Width 1000
Set Height 600  # Medium size - about 30 lines
Set Theme "Builtin Solarized Light"

Output output/medium-terminal.gif

Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

# Use a directory with existing extensions if available

Type "./gemini-cli-manager --state-dir /tmp/vhs-medium-test"
Enter
Sleep 3

# Initial view with extensions
Screenshot output/medium-01-extensions.png

# Navigate down through extensions
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/medium-02-extensions-scrolled.png

# Go to settings
Tab
Tab
Sleep 1s
Screenshot output/medium-03-settings.png

# Scroll in settings
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/medium-04-settings-scrolled.png

# Go to help
Tab
Sleep 1s
Screenshot output/medium-05-help.png

# Scroll in help
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/medium-06-help-scrolled.png

# Exit
Type "q"
Sleep 500ms

