package main

import (
	"fmt"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

func main() {
	// Initialize theme
	theme.SetTheme("github-dark")
	
	// Set up colors
	colorBorder := theme.Border()
	colorAccent := theme.Primary()
	colorText := theme.TextPrimary()
	
	// Helper function for tab borders
	tabBorderWithBottom := func(left, middle, right string) lipgloss.Border {
		border := lipgloss.RoundedBorder()
		border.BottomLeft = left
		border.Bottom = middle
		border.BottomRight = right
		return border
	}
	
	// Define tab borders
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	
	// Create tab styles
	inactiveTabStyle := lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(colorBorder).
		Padding(0, 2).
		Foreground(colorText)
	
	activeTabStyle := lipgloss.NewStyle().
		Border(activeTabBorder, true).
		BorderForeground(colorAccent).
		Padding(0, 2).
		Foreground(colorAccent).
		Bold(true)
	
	// Define tabs using the component
	tabs := []components.Tab{
		{Title: "Extensions", Icon: "🧩", ID: "extensions"},
		{Title: "Profiles", Icon: "👤", ID: "profiles"},
		{Title: "Settings", Icon: "🔧", ID: "settings"},
		{Title: "Help", Icon: "❓", ID: "help"},
	}
	
	// Create tab bar
	width := 100
	tabBar := components.NewTabBar(tabs, width)
	tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	tabBar.SetActiveIndex(0) // Extensions active
	
	// Render with content
	contentText := "Extensions\n\n2 extensions found\n\nThis is the content area that connects to the tabs above."
	result := tabBar.RenderWithContent(contentText, 10)
	
	// Add padding
	final := lipgloss.NewStyle().Padding(2, 3).Render(result)
	
	fmt.Println(final)
}