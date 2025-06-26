package cli

import (
	"regexp"
	"strings"
)

// ansiRegex matches ANSI escape sequences
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][^\\]*\x07|\x1b[PX^_].*?\x1b\\|\x1b\[[0-9;]*[mGKHflSTu]|\x1b\[[0-9;]*;?[0-9]+[mGKHflSTuR]`)

// stripANSI removes ANSI escape sequences from a string
func stripANSI(s string) string {
	// Remove ANSI escape sequences
	s = ansiRegex.ReplaceAllString(s, "")
	
	// Also remove other control characters
	var result strings.Builder
	for _, r := range s {
		if r >= 32 && r != 127 { // printable characters
			result.WriteRune(r)
		}
	}
	
	return strings.TrimSpace(result.String())
}

// truncateString truncates a string to a maximum length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}