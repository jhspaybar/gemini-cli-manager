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
	
	// Set up colors and styles
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
	
	// Create styles
	inactiveTabStyle := lipgloss.NewStyle().
		Border(tabBorderWithBottom("┴", "─", "┴"), true).
		BorderForeground(colorBorder).
		Padding(0, 2).
		Foreground(colorText)
	
	activeTabStyle := lipgloss.NewStyle().
		Border(tabBorderWithBottom("┘", " ", "└"), true).
		BorderForeground(colorAccent).
		Padding(0, 2).
		Foreground(colorAccent).
		Bold(true)
	
	// Define tabs
	tabs := []components.Tab{
		{Title: "Extensions", Icon: "🧩", ID: "extensions"},
		{Title: "Profiles", Icon: "👤", ID: "profiles"},
		{Title: "Settings", Icon: "🔧", ID: "settings"},
		{Title: "Help", Icon: "❓", ID: "help"},
	}
	
	// Content for each tab
	contents := map[string]string{
		"extensions": "Extensions\n\n2 extensions found\n\n- mcp-extension v2.1.0\n- simple-extension v1.0.0",
		"profiles":   "Profiles\n\nNo active profile\n\nPress 'n' to create a new profile",
		"settings":   "Settings\n\n🎨 Appearance\n  ▶ ✓ Solarized Light\n    Solarized Dark\n    GitHub\n    Nord",
		"help":       "Help\n\n🧭 Navigation\n  Tab    Next tab\n  ←/h    Previous tab\n  →/l    Next tab",
	}
	
	// Demonstrate each tab
	fmt.Println("Tab Bar Component Demo")
	fmt.Println("======================\n")
	
	for i, tab := range tabs {
		fmt.Printf("%d. %s %s Tab:\n", i+1, tab.Icon, tab.Title)
		fmt.Println(strings.Repeat("-", 40))
		
		// Create tab bar
		tabBar := components.NewTabBar(tabs, 80)
		tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
		tabBar.SetActiveByID(tab.ID)
		
		// Get content for this tab
		content := contents[tab.ID]
		
		// Render
		result := tabBar.RenderWithContent(content, 8)
		fmt.Println(result)
		fmt.Println()
	}
}

