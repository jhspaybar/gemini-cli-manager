# Tab Navigation Test
# Tests basic navigation through tab menu items

Output output/navigation-tabs.gif

# VHS Configuration
Set FontSize 16
Set Width 1200
Set Height 1000
Set Theme "Builtin Solarized Light"
Set TypingSpeed 100ms
Set Shell "bash"
Set LetterSpacing 1
Set LineHeight 1.5
Set CursorBlink false

# Start from the project root directory
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

# Set VHS environment variable
Type "export VHS=true"
Enter
Sleep 500ms

# Start application with test state directory using go run
Type "go run . --state-dir /tmp/vhs-nav-test"
Enter
Sleep 5s  # Wait for app to fully load
Screenshot output/tabs-01-start.png
Sleep 1s  # Extra pause before starting navigation

# Test tab navigation with Tab key
Tab
Sleep 1s
Screenshot output/tabs-02-profiles.png

Tab
Sleep 1s
Screenshot output/tabs-03-settings.png

Tab
Sleep 1s
Screenshot output/tabs-04-help.png

Tab
Sleep 1s
Screenshot output/tabs-05-back-to-extensions.png

# Test navigation with arrow keys
Left
Sleep 1s
Screenshot output/tabs-06-help-via-left.png

Right
Sleep 1s
Screenshot output/tabs-07-extensions-via-right.png

# Test vim-style navigation
Type "l"
Sleep 1s
Screenshot output/tabs-08-profiles-via-l.png

Type "h"
Sleep 1s
Screenshot output/tabs-09-extensions-via-h.png

# Exit
Type "q"
Sleep 500ms