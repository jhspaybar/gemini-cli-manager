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
	
	// Test what happens when tabs exceed available width
	fmt.Println("Tab Overflow Test")
	fmt.Println("=================")
	fmt.Println()
	
	// Many tabs with long names
	longTabs := []components.Tab{
		{Title: "Documentation", Icon: "📚", ID: "docs"},
		{Title: "Configuration", Icon: "⚙️", ID: "config"},
		{Title: "Performance", Icon: "⚡", ID: "perf"},
		{Title: "Security", Icon: "🔒", ID: "security"},
		{Title: "Analytics", Icon: "📊", ID: "analytics"},
		{Title: "Integration", Icon: "🔗", ID: "integration"},
		{Title: "Deployment", Icon: "🚀", ID: "deploy"},
	}
	
	// Test with different widths
	widths := []int{120, 100, 80, 60, 40}
	
	for _, width := range widths {
		fmt.Printf("Width: %d characters\n", width)
		fmt.Printf("-------------------\n")
		
		// Create tab bar
		tabBar := components.NewTabBar(longTabs, width)
		tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
		tabBar.SetActiveIndex(3) // Security tab
		
		// Just render the tab bar to see overflow behavior
		tabRow := tabBar.Render()
		
		// Show the result
		fmt.Println(tabRow)
		fmt.Printf("Actual width: %d\n", lipgloss.Width(tabRow))
		fmt.Println()
	}
	
	// Test with extremely long single tab
	fmt.Println("\nExtremely Long Tab Test")
	fmt.Println("=======================")
	fmt.Println()
	
	veryLongTab := []components.Tab{
		{Title: "This is an extremely long tab title that should definitely overflow", Icon: "🌟", ID: "long"},
		{Title: "Short", Icon: "📌", ID: "short"},
	}
	
	tabBar := components.NewTabBar(veryLongTab, 60)
	tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	tabBar.SetActiveIndex(0)
	
	result := tabBar.RenderWithContent("Content for the very long tab", 5)
	fmt.Println(result)
}