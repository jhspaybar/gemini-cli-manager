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
	
	// Test cases with different numbers of tabs
	testCases := []struct {
		name string
		tabs []components.Tab
		activeIndex int
		width int
	}{
		{
			name: "Single Tab",
			tabs: []components.Tab{
				{Title: "Dashboard", Icon: "📊", ID: "dashboard"},
			},
			activeIndex: 0,
			width: 80,
		},
		{
			name: "Two Tabs",
			tabs: []components.Tab{
				{Title: "Files", Icon: "📁", ID: "files"},
				{Title: "Search", Icon: "🔍", ID: "search"},
			},
			activeIndex: 1,
			width: 80,
		},
		{
			name: "Three Tabs (Narrow)",
			tabs: []components.Tab{
				{Title: "Home", Icon: "🏠", ID: "home"},
				{Title: "Edit", Icon: "✏️", ID: "edit"},
				{Title: "Save", Icon: "💾", ID: "save"},
			},
			activeIndex: 0,
			width: 60,
		},
		{
			name: "Five Tabs",
			tabs: []components.Tab{
				{Title: "Code", Icon: "💻", ID: "code"},
				{Title: "Build", Icon: "🔨", ID: "build"},
				{Title: "Test", Icon: "🧪", ID: "test"},
				{Title: "Deploy", Icon: "🚀", ID: "deploy"},
				{Title: "Monitor", Icon: "📈", ID: "monitor"},
			},
			activeIndex: 2,
			width: 100,
		},
		{
			name: "Many Short Tabs",
			tabs: []components.Tab{
				{Title: "A", Icon: "🅰️", ID: "a"},
				{Title: "B", Icon: "🅱️", ID: "b"},
				{Title: "C", Icon: "©️", ID: "c"},
				{Title: "D", Icon: "🆔", ID: "d"},
				{Title: "E", Icon: "📧", ID: "e"},
				{Title: "F", Icon: "🎏", ID: "f"},
			},
			activeIndex: 3,
			width: 80,
		},
		{
			name: "Mixed Emoji Types",
			tabs: []components.Tab{
				{Title: "Fire", Icon: "🔥", ID: "fire"},
				{Title: "Water", Icon: "💧", ID: "water"},
				{Title: "Earth", Icon: "🌍", ID: "earth"},
				{Title: "Air", Icon: "💨", ID: "air"},
			},
			activeIndex: 1,
			width: 90,
		},
		{
			name: "Development Workflow",
			tabs: []components.Tab{
				{Title: "Issues", Icon: "🐛", ID: "issues"},
				{Title: "PRs", Icon: "🔄", ID: "prs"},
				{Title: "CI/CD", Icon: "♻️", ID: "cicd"},
				{Title: "Releases", Icon: "📦", ID: "releases"},
				{Title: "Metrics", Icon: "📊", ID: "metrics"},
			},
			activeIndex: 4,
			width: 120,
		},
		{
			name: "Unicode Mix",
			tabs: []components.Tab{
				{Title: "Star", Icon: "⭐", ID: "star"},
				{Title: "Heart", Icon: "❤️", ID: "heart"},
				{Title: "Check", Icon: "✅", ID: "check"},
				{Title: "Cross", Icon: "❌", ID: "cross"},
			},
			activeIndex: 2,
			width: 80,
		},
	}
	
	// Run test cases
	for _, tc := range testCases {
		fmt.Printf("\n%s\n", lipgloss.NewStyle().Bold(true).Render(tc.name))
		fmt.Println(strings.Repeat("=", len(tc.name)))
		fmt.Printf("Tabs: %d | Active: %s | Width: %d\n", 
			len(tc.tabs), 
			tc.tabs[tc.activeIndex].Title,
			tc.width)
		fmt.Println()
		
		// Create tab bar
		tabBar := components.NewTabBar(tc.tabs, tc.width)
		tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
		tabBar.SetActiveIndex(tc.activeIndex)
		
		// Generate content
		activeTab := tc.tabs[tc.activeIndex]
		content := fmt.Sprintf("%s %s\n\nThis is the %s tab content.\n\nTab %d of %d is active.",
			activeTab.Icon,
			activeTab.Title,
			activeTab.Title,
			tc.activeIndex + 1,
			len(tc.tabs))
		
		// Render
		result := tabBar.RenderWithContent(content, 8)
		
		// Add some padding for display
		final := lipgloss.NewStyle().Padding(1, 2).Render(result)
		fmt.Println(final)
		fmt.Println()
	}
	
	// Edge case: Empty tabs
	fmt.Printf("\n%s\n", lipgloss.NewStyle().Bold(true).Render("Edge Case: No Tabs"))
	fmt.Println("==================")
	emptyTabBar := components.NewTabBar([]components.Tab{}, 80)
	emptyTabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	result := emptyTabBar.RenderWithContent("No tabs defined", 5)
	fmt.Println(result)
}