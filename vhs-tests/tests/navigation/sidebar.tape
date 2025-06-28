# Sidebar Navigation Test
# Tests basic navigation through sidebar menu items

Output output/navigation-sidebar.gif

# VHS Configuration
Set FontSize 16
Set Width 1200
Set Height 800
Set Theme "github-dark"
Set TypingSpeed 100ms

# Start application with test state directory
Type "cd ../.."
Enter
Sleep 1s

Type "./gemini-cli-manager --state-dir /tmp/vhs-nav-test"
Enter
Sleep 3s
Screenshot output/nav-01-start.png

# Test moving down through sidebar
Type "j"
Sleep 500ms
Screenshot output/nav-02-profiles.png

Type "j"
Sleep 500ms
Screenshot output/nav-03-settings.png

Type "j"
Sleep 500ms
Screenshot output/nav-04-help.png

# Test moving up
Type "k"
Sleep 500ms
Screenshot output/nav-05-back-to-settings.png

# Test selection with Enter
Type "Enter"
Sleep 1s
Screenshot output/nav-06-settings-view.png

# Go back to sidebar navigation
Type "Tab"
Sleep 500ms
Screenshot output/nav-07-sidebar-focus.png

# Navigate to Extensions
Type "k"
Sleep 500ms
Type "k"
Sleep 500ms
Type "Enter"
Sleep 1s
Screenshot output/nav-08-extensions-view.png

# Test space bar selection
Type "Tab"
Sleep 500ms
Type "j"
Sleep 500ms
Type " "  # Space bar
Sleep 1s
Screenshot output/nav-09-profiles-via-space.png

# Exit
Type "q"
Sleep 500ms