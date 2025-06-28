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

	fmt.Println("Edge Case Layout Test")
	fmt.Println("====================")
	fmt.Println()

	// Test edge cases
	cases := []struct {
		name   string
		width  int
		height int
		padding int
	}{
		{"Initial render (0x0)", 0, 0, 3},
		{"Very small terminal", 40, 10, 3},
		{"Small terminal", 60, 15, 3},
		{"Normal terminal", 80, 24, 3},
		{"No padding space", 80, 24, 0},
	}

	for _, tc := range cases {
		fmt.Printf("\n%s (window=%dx%d, padding=%d):\n", tc.name, tc.width, tc.height, tc.padding)
		fmt.Println(strings.Repeat("-", 80))

		// Handle zero dimensions
		width := tc.width
		height := tc.height
		if width == 0 || height == 0 {
			width = 80
			height = 24
			fmt.Printf("(Using defaults: %dx%d)\n", width, height)
		}

		// Calculate content dimensions
		horizontalPadding := tc.padding
		verticalPadding := 2
		contentWidth := width - (horizontalPadding * 2)
		contentHeight := height - (verticalPadding * 2)
		
		// Ensure minimum dimensions
		if contentWidth < 20 {
			contentWidth = 20
		}
		if contentHeight < 5 {
			contentHeight = 5
		}

		fmt.Printf("Content area: %dx%d\n", contentWidth, contentHeight)

		// Create components
		tabs := []components.Tab{
			{Title: "Extensions", Icon: "🧩", ID: "extensions"},
			{Title: "Profiles", Icon: "👤", ID: "profiles"},
		}
		
		tabBar := components.NewTabBar(tabs, contentWidth)
		tabBar.SetStyles(
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Primary()).
				Padding(0, 1),
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Border()).
				Padding(0, 1),
			theme.Border(),
		)

		statusBar := components.NewStatusBar(contentWidth)
		statusBar.SetLeftItems(components.ProfileStatusItems("Dev", 3, 5))
		statusBar.SetKeyBindings([]components.KeyBinding{
			{"q", "Quit"},
			{"Tab", "Switch"},
		})

		// Calculate heights
		statusBarRendered := statusBar.Render()
		tabsRendered := tabBar.Render()
		
		statusHeight := lipgloss.Height(statusBarRendered)
		tabHeight := lipgloss.Height(tabsRendered)
		
		contentBoxHeight := contentHeight - statusHeight - tabHeight - 1
		
		fmt.Printf("Heights: tabs=%d, status=%d, contentBox=%d\n", 
			tabHeight, statusHeight, contentBoxHeight)
			
		// Check for negative heights
		if contentBoxHeight < 1 {
			fmt.Printf("WARNING: Content box height is %d (too small!)\n", contentBoxHeight)
			contentBoxHeight = 1
		}
		
		// Render a small sample
		content := "Test content"
		tabsAndContent := tabBar.RenderWithContent(content, contentBoxHeight)
		
		// Show first 5 lines
		lines := strings.Split(tabsAndContent, "\n")
		for i := 0; i < 5 && i < len(lines); i++ {
			if len(lines[i]) > 60 {
				fmt.Printf("Line %d: %s...\n", i+1, lines[i][:60])
			} else {
				fmt.Printf("Line %d: %s\n", i+1, lines[i])
			}
		}
	}
}