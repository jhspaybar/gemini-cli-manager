# Debug viewport rendering
# Simple test to check if content is rendering

Set FontSize 14
Set Width 1000
Set Height 600
Set Theme "Builtin Solarized Light"

Output output/debug-viewport.gif

Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

Type "./gemini-cli-manager --state-dir /tmp/vhs-debug-$(date +%s)"
Enter
Sleep 3

# Check extensions view
Screenshot output/debug-01-extensions.png

# Try help view  
Type "?"
Sleep 2s
Screenshot output/debug-02-help.png

# Try settings
Type "?"  # Back to extensions
Sleep 500ms
Tab
Tab
Sleep 1s
Screenshot output/debug-03-settings.png

# Exit
Type "q"
Sleep 500ms