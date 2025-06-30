# Test viewport scrolling functionality
# Simple test to verify scrolling works in small terminals

# VHS Configuration - Small terminal
Set FontSize 14
Set Width 800
Set Height 400  # Small height to test scrolling
Set Theme "Builtin Solarized Light"

# Start recording
Output output/viewport-test.gif

# Navigate to project root and start the app
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s
Type "./gemini-cli-manager --state-dir /tmp/vhs-viewport-$(date +%s)"
Enter
Sleep 3

# The app should load with Extensions tab active
Screenshot output/viewport-01-extensions.png

# Navigate to Help tab which has long content
Tab
Tab  
Tab
Sleep 1
Screenshot output/viewport-02-help-initial.png

# Try scrolling down with j key
Type "j"
Sleep 500ms
Screenshot output/viewport-03-help-scrolled.png

# Exit
Type "q"
Sleep 500ms