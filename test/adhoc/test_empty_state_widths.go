package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

func main() {
	// Initialize theme
	theme.SetTheme("github-dark")

	fmt.Println("EmptyState Width Rendering Test")
	fmt.Println("===============================")
	fmt.Println()
	fmt.Println("Testing that borders render completely at various widths")
	fmt.Println()

	// Test at different container widths
	widths := []int{40, 50, 60, 70, 80, 90, 100}
	
	for _, width := range widths {
		fmt.Printf("Container Width: %d\n", width)
		fmt.Println(strings.Repeat("-", width))
		
		emptyState := components.NewEmptyState(width).
			NoItemsFound().
			SetAction("Press '/' to modify your search")
		
		rendered := emptyState.Render()
		fmt.Println(rendered)
		
		// Check actual rendered width
		actualWidth := lipgloss.Width(rendered)
		fmt.Printf("Actual rendered width: %d\n", actualWidth)
		
		// Verify borders by checking first and last characters of each line
		lines := strings.Split(rendered, "\n")
		borderIssues := false
		for i, line := range lines {
			if len(line) > 0 {
				runes := []rune(line)
				if len(runes) > 0 {
					// Check if line appears truncated
					if width > 60 && len(runes) < width-10 {
						// For lines that should span more of the width
						lastChar := runes[len(runes)-1]
						if lastChar != ' ' && lastChar != '╮' && lastChar != '│' && lastChar != '╯' {
							fmt.Printf("  Line %d might be truncated, ends with: '%c'\n", i+1, lastChar)
							borderIssues = true
						}
					}
				}
			}
		}
		
		if !borderIssues {
			fmt.Println("  ✅ All borders appear complete")
		}
		
		fmt.Println()
	}
}