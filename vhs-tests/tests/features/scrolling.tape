# Test scrolling in small terminals
# Verifies that content scrolls properly when terminal is small

# VHS Configuration - Small terminal
Set FontSize 14
Set Width 600
Set Height 300  # Very small height to force scrolling
Set Theme "Builtin Solarized Light"
Set TypingSpeed 100ms

# Start the app
Type "./gemini-cli-manager --state-dir /tmp/vhs-scroll-test"
Enter
Sleep 2

# Take screenshot of initial state
Screenshot output/scroll-01-initial.png

# Navigate to Help tab (which has long content)
Tab
Tab
Tab
Sleep 1
Screenshot output/scroll-02-help-tab.png

# Try scrolling down
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/scroll-03-scrolled-down.png

# Page down
Type "pgdn"
Sleep 500ms
Screenshot output/scroll-04-page-down.png

# Go to bottom
Type "end"
Sleep 500ms
Screenshot output/scroll-05-bottom.png

# Go back to top
Type "home"
Sleep 500ms
Screenshot output/scroll-06-back-to-top.png

# Switch to Settings (another scrollable view)
Type "h"
Sleep 1s
Screenshot output/scroll-07-settings.png

# Scroll in settings
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/scroll-08-settings-scrolled.png

# Exit
Type "q"
Sleep 500ms