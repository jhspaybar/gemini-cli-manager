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

	fmt.Println("Precise Alignment Test")
	fmt.Println("======================")
	fmt.Println()

	// Test with actual window size
	width := 80
	height := 24
	horizontalPadding := 2
	verticalPadding := 0
	contentWidth := width - (horizontalPadding * 2)
	contentHeight := height - (verticalPadding * 2)

	fmt.Printf("Window: %dx%d\n", width, height)
	fmt.Printf("Padding: H=%d, V=%d\n", horizontalPadding, verticalPadding)
	fmt.Printf("Content area: %dx%d\n", contentWidth, contentHeight)
	fmt.Println()

	// Create components exactly as in the app
	tabs := []components.Tab{
		{Title: "Extensions", Icon: "🧩", ID: "extensions"},
		{Title: "Profiles", Icon: "👤", ID: "profiles"},
		{Title: "Settings", Icon: "🔧", ID: "settings"},
		{Title: "Help", Icon: "❓", ID: "help"},
	}

	tabBar := components.NewTabBar(tabs, contentWidth)
	
	// Set up tab styles
	activeTabBorder := lipgloss.RoundedBorder()
	activeTabBorder.BottomLeft = "┘"
	activeTabBorder.Bottom = " "
	activeTabBorder.BottomRight = "└"
	
	inactiveTabBorder := lipgloss.RoundedBorder()
	inactiveTabBorder.BottomLeft = "┴"
	inactiveTabBorder.Bottom = "─"
	inactiveTabBorder.BottomRight = "┴"
	
	activeStyle := lipgloss.NewStyle().
		Border(activeTabBorder).
		BorderForeground(theme.Primary()).
		Foreground(theme.TextPrimary()).
		Background(theme.Selection()).
		Padding(0, 1)
		
	inactiveStyle := lipgloss.NewStyle().
		Border(inactiveTabBorder).
		BorderForeground(theme.Border()).
		Foreground(theme.TextSecondary()).
		Padding(0, 1)
		
	tabBar.SetStyles(activeStyle, inactiveStyle, theme.Border())
	tabBar.SetActiveByID("extensions")

	// Create status bar
	statusBar := components.NewStatusBar(contentWidth)
	statusBar.SetLeftItems(components.ProfileStatusItems("Dev profile", 1, 2))
	statusBar.SetKeyBindings([]components.KeyBinding{
		{"Tab", "Switch"},
		{"L", "Launch"},
		{"q", "Quit"},
	})

	// Render components
	statusBarRendered := statusBar.Render()
	statusHeight := lipgloss.Height(statusBarRendered)
	
	tabsRendered := tabBar.Render()
	tabHeight := lipgloss.Height(tabsRendered)
	
	// Calculate content box height
	contentBoxHeight := contentHeight - statusHeight - tabHeight - 1
	
	fmt.Printf("Component heights: tabs=%d, status=%d, contentBox=%d\n", 
		tabHeight, statusHeight, contentBoxHeight)
	fmt.Println()
	
	// Create some content
	content := "Extensions\n\n2 extensions found"
	
	// Use TabBar's RenderWithContent
	tabsAndContent := tabBar.RenderWithContent(content, contentBoxHeight)
	
	// Combine everything
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabsAndContent,
		statusBarRendered,
	)
	
	// Apply padding
	final := lipgloss.NewStyle().
		Padding(verticalPadding, horizontalPadding).
		Render(fullContent)
	
	// Analyze each line
	fmt.Println("Line-by-line analysis:")
	fmt.Println("----------------------")
	lines := strings.Split(final, "\n")
	
	for i, line := range lines {
		runes := []rune(line)
		if len(runes) > 0 {
			// Find first and last non-space characters
			firstNonSpace := -1
			lastNonSpace := -1
			for j, r := range runes {
				if r != ' ' {
					if firstNonSpace == -1 {
						firstNonSpace = j
					}
					lastNonSpace = j
				}
			}
			
			if firstNonSpace >= 0 {
				fmt.Printf("Line %2d: First char at %2d '%c', Last char at %2d '%c', Width=%d\n", 
					i+1, firstNonSpace, runes[firstNonSpace], 
					lastNonSpace, runes[lastNonSpace], 
					lipgloss.Width(line))
			} else {
				fmt.Printf("Line %2d: (empty)\n", i+1)
			}
		}
	}
	
	// Show the actual output
	fmt.Println("\nRendered output:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Print(final)
}