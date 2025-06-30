# Test in small but usable terminal
# Verifies scrolling works when terminal is small but still usable

Set FontSize 14
Set Width 800
Set Height 400  # Small but reasonable - about 20 lines
Set Theme "Builtin Solarized Light"

Output output/small-terminal.gif

Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

Type "./gemini-cli-manager --state-dir /tmp/vhs-small-$(date +%s)"
Enter
Sleep 3

# Initial view
Screenshot output/small-01-initial.png

# Navigate to settings which has many items
Tab
Tab
Sleep 1s
Screenshot output/small-02-settings.png

# Try scrolling down in settings
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/small-03-settings-scrolled.png

# More scrolling
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/small-04-settings-more-scroll.png

# Go to help
Tab
Sleep 1s
Screenshot output/small-05-help.png

# Scroll in help
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Screenshot output/small-06-help-scrolled.png

# Page down
Type "pgdn"
Sleep 500ms
Screenshot output/small-07-help-pagedown.png

# Back to top
Type "home"
Sleep 500ms
Screenshot output/small-08-help-top.png

# Exit
Type "q"
Sleep 500ms