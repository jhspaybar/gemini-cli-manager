# Test list scrolling with multiple extensions
# Verifies that lists scroll properly when there are many items

# VHS Configuration - Small terminal
Set FontSize 14
Set Width 800
Set Height 300  # Very small to force scrolling
Set Theme "Builtin Solarized Light"

# Start recording
Output output/list-scrolling.gif

# Navigate to project root
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

# Setup multiple test extensions
Type "mkdir -p /tmp/vhs-list-test/extensions"
Enter
Sleep 500ms

# Create some dummy extensions for testing
Type "for i in {1..10}; do mkdir -p /tmp/vhs-list-test/extensions/test-ext-$i; echo '{\"name\":\"Test Extension $i\",\"version\":\"1.0.$i\",\"description\":\"Test extension number $i for scrolling test\"}' > /tmp/vhs-list-test/extensions/test-ext-$i/extension.json; done"
Enter
Sleep 1s

# Start the app with test directory
Type "./gemini-cli-manager --state-dir /tmp/vhs-list-test"
Enter
Sleep 3

# Take screenshot of initial state with many extensions
Screenshot output/list-01-many-extensions.png

# Navigate down through the list
Type "j"
Sleep 500ms
Type "j"  
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/list-02-scrolled.png

# Navigate to bottom
Type "G"
Sleep 500ms
Screenshot output/list-03-bottom.png

# Navigate back to top
Type "g"
Sleep 500ms
Screenshot output/list-04-top.png

# Exit
Type "q"
Sleep 500ms

# Cleanup
Type "rm -rf /tmp/vhs-list-test"
Enter