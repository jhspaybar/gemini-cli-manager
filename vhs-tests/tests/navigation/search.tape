# Search Functionality Test
# Tests the search bar behavior and context-aware placeholders

Output output/navigation-search.gif

# VHS Configuration
Set FontSize 16
Set Width 1200
Set Height 1000
Set Theme "Builtin Solarized Light"
Set TypingSpeed 75ms

# Start from the project root directory
Type "cd /Users/williamthurston/repos/gemini-cli-manager"
Enter
Sleep 1s

# Start with test directory using go run
Type "go run . --state-dir /tmp/vhs-search-test"
Enter
Sleep 3s
Screenshot output/search-01-start.png

# Activate search in Extensions view
Type "/"
Sleep 500ms
Screenshot output/search-02-extensions-search-active.png

# Type a search query
Type "test"
Sleep 1s
Screenshot output/search-03-extensions-query.png

# Clear search
Type "Ctrl+U"
Sleep 500ms
Screenshot output/search-04-extensions-cleared.png

# Cancel search
Type "Escape"
Sleep 500ms
Screenshot output/search-05-search-cancelled.png

# Navigate to Profiles
Type "j"
Sleep 500ms
Type "Enter"
Sleep 1s
Screenshot output/search-06-profiles-view.png

# Activate search in Profiles view
Type "/"
Sleep 500ms
Screenshot output/search-07-profiles-search-active.png

# Note the different placeholder text
Sleep 1s
Screenshot output/search-08-profiles-placeholder.png

# Type a profile search
Type "dev"
Sleep 1s
Screenshot output/search-09-profiles-query.png

# Exit
Type "Escape"
Sleep 500ms
Type "q"
Sleep 500ms