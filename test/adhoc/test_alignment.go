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

	fmt.Println("Alignment Test")
	fmt.Println("==============")
	fmt.Println()

	width := 80
	height := 24

	// Simulate the actual layout
	horizontalPadding := 3
	verticalPadding := 1
	contentWidth := width - (horizontalPadding * 2)
	contentHeight := height - (verticalPadding * 2)

	fmt.Printf("Window: %dx%d, Content area: %dx%d\n", width, height, contentWidth, contentHeight)
	fmt.Println()

	// Create tabs
	tabs := []components.Tab{
		{Title: "Extensions", Icon: "🧩", ID: "extensions"},
		{Title: "Profiles", Icon: "👤", ID: "profiles"},
		{Title: "Settings", Icon: "🔧", ID: "settings"},
		{Title: "Help", Icon: "❓", ID: "help"},
	}

	tabBar := components.NewTabBar(tabs, contentWidth)
	
	// Set up tab styles
	activeStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary()).
		Foreground(theme.TextPrimary()).
		Background(theme.Selection()).
		Padding(0, 1)
		
	inactiveStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border()).
		Foreground(theme.TextSecondary()).
		Padding(0, 1)
		
	tabBar.SetStyles(activeStyle, inactiveStyle, theme.Border())
	tabBar.SetActiveByID("extensions")

	// Create status bar with same width as content
	statusBar := components.NewStatusBar(contentWidth)
	statusBar.SetLeftItems(components.ProfileStatusItems("Dev", 1, 2))
	statusBar.SetKeyBindings([]components.KeyBinding{
		{"Tab", "Switch"},
		{"L", "Launch"},
		{"q", "Quit"},
	})

	// Render components
	tabsRendered := tabBar.Render()
	statusBarRendered := statusBar.Render()
	
	// Create some content
	content := "Test content for alignment"
	
	// Calculate heights
	tabHeight := lipgloss.Height(tabsRendered)
	statusHeight := lipgloss.Height(statusBarRendered)
	contentBoxHeight := contentHeight - statusHeight - tabHeight - 1
	
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
	
	// Print the result
	fmt.Println("Rendered output:")
	fmt.Println(strings.Repeat("-", 80))
	
	// Show first 10 lines to check alignment
	lines := strings.Split(final, "\n")
	for i := 0; i < 10 && i < len(lines); i++ {
		// Add line numbers to check alignment
		fmt.Printf("%2d: %s\n", i+1, lines[i])
	}
	
	fmt.Println("...")
	
	// Show last 5 lines
	start := len(lines) - 5
	if start < 10 {
		start = 10
	}
	for i := start; i < len(lines); i++ {
		fmt.Printf("%2d: %s\n", i+1, lines[i])
	}
	
	// Measure widths and alignment
	fmt.Println("\nAlignment analysis:")
	
	// Find the content box line and status bar line
	for i, line := range lines {
		if strings.Contains(line, "│") && strings.Contains(line, "Test content") {
			fmt.Printf("Content line %d: starts at col %d, ends at col %d\n", 
				i+1, strings.Index(line, "│"), strings.LastIndex(line, "│"))
		}
		if strings.Contains(line, "👤 Dev") {
			// Find where the status bar content actually starts and ends
			trimmed := strings.TrimLeft(line, " ")
			startPos := len(line) - len(trimmed)
			endPos := startPos + lipgloss.Width(strings.TrimSpace(trimmed))
			fmt.Printf("Status line %d: content starts at col %d, ends at col %d\n", 
				i+1, startPos, endPos)
		}
	}
	
	// Visual ruler
	fmt.Println("\nRuler:")
	fmt.Println("    1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
	fmt.Println("    0        10        20        30        40        50        60        70        80        90       100")
}