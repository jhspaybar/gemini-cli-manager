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

	fmt.Println("Full Layout Test")
	fmt.Println("================")
	fmt.Println()

	// Test different terminal sizes
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Small (80x24)", 80, 24},
		{"Medium (100x30)", 100, 30},
		{"Large (120x40)", 120, 40},
	}

	for _, size := range sizes {
		fmt.Printf("\n%s:\n", size.name)
		fmt.Println(strings.Repeat("-", size.width))

		// Create tabs
		tabs := []components.Tab{
			{Title: "Extensions", Icon: "🧩", ID: "extensions"},
			{Title: "Profiles", Icon: "👤", ID: "profiles"},
			{Title: "Settings", Icon: "🔧", ID: "settings"},
			{Title: "Help", Icon: "❓", ID: "help"},
		}

		tabBar := components.NewTabBar(tabs, size.width)
		// Create tab styles
		activeTabBorder := lipgloss.RoundedBorder()
		activeTabBorder.BottomLeft = "┘"
		activeTabBorder.Bottom = " "
		activeTabBorder.BottomRight = "└"
		
		inactiveTabBorder := lipgloss.RoundedBorder()
		inactiveTabBorder.BottomLeft = "┴"
		inactiveTabBorder.Bottom = "─"
		inactiveTabBorder.BottomRight = "┴"
		
		tabBar.SetStyles(
			lipgloss.NewStyle().
				Border(activeTabBorder).
				BorderForeground(theme.Primary()).
				Foreground(theme.TextPrimary()).
				Background(theme.Selection()).
				Padding(0, 1),
			lipgloss.NewStyle().
				Border(inactiveTabBorder).
				BorderForeground(theme.Border()).
				Foreground(theme.TextSecondary()).
				Padding(0, 1),
			theme.Border(),
		)
		tabBar.SetActiveByID("extensions")

		// Create status bar
		statusBar := components.NewStatusBar(size.width)
		statusBar.SetLeftItems(components.ProfileStatusItems("Production", 5, 12))
		statusBar.SetKeyBindings(components.CommonKeyBindings())

		// Create some dummy content
		content := fmt.Sprintf("Content area (%dx%d)\n", size.width, size.height-6)
		content += "This is where the main content would go.\n"
		content += "It should fill the available space between tabs and status bar."

		// Render status bar and get its height
		statusBarRendered := statusBar.Render()
		statusHeight := lipgloss.Height(statusBarRendered)
		
		// Render tabs to get height
		tabsRendered := tabBar.Render()
		tabHeight := lipgloss.Height(tabsRendered)
		
		// Calculate height for content box only
		contentBoxHeight := size.height - statusHeight - tabHeight - 1
		
		// Use TabBar's RenderWithContent for seamless connection
		tabsAndContent := tabBar.RenderWithContent(content, contentBoxHeight)
		
		// Combine using JoinVertical
		output := lipgloss.JoinVertical(
			lipgloss.Left,
			tabsAndContent,
			statusBarRendered,
		)
		fmt.Println(output)
	}
}

