# [Test Name]
# [Brief description of what this test covers]

Output output/[category]-[test-name].gif

# VHS Configuration
Set FontSize 16
Set Width 1200
Set Height 800
Set Theme "github-dark"
Set TypingSpeed 100ms

# Navigate to project root
Type "cd ../.."
Enter
Sleep 1s

# Start application with test state directory
Type "./gemini-cli-manager --state-dir /tmp/vhs-[test-name]"
Enter
Sleep 3s
Screenshot output/[test]-01-start.png

# [Test Step 1 Description]
Type "[command]"
Sleep 500ms
Screenshot output/[test]-02-[description].png

# [Test Step 2 Description]
Type "[command]"
Sleep 500ms
Screenshot output/[test]-03-[description].png

# Add more test steps as needed...

# Exit application
Type "q"
Sleep 500ms

# Cleanup (optional)
# Type "rm -rf /tmp/vhs-[test-name]"
# Enter