// Package main provides visual tests for UI components
// Run with: go run test/adhoc/visual_tests.go
// Or: go run test/adhoc/visual_tests.go tabs
// Or: go run test/adhoc/visual_tests.go --list
package main

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

// TestFunc represents a visual test function
type TestFunc func()

// TestRegistry holds all available tests
var TestRegistry map[string]TestFunc

func init() {
	TestRegistry = map[string]TestFunc{
		"tabs":          testTabs,
		"tabs-dynamic":  testTabsDynamic,
		"tabs-overflow": testTabsOverflow,
		"tabs-switch":   testTabsSwitching,
		"emoji-width":   testEmojiWidth,
		"gear-spacing":  testGearSpacing,
		"all":           runAllTests,
	}
}

func main() {
	// Initialize theme
	theme.SetTheme("github-dark")
	
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	
	arg := os.Args[1]
	
	// Handle special commands
	switch arg {
	case "--list", "-l":
		listTests()
		return
	case "--help", "-h":
		printUsage()
		return
	}
	
	// Run specific test
	if testFunc, exists := TestRegistry[arg]; exists {
		fmt.Printf("Running visual test: %s\n", arg)
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println()
		testFunc()
	} else {
		fmt.Printf("Unknown test: %s\n\n", arg)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: go run test/adhoc/visual_tests.go [test-name]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  --list, -l     List all available tests")
	fmt.Println("  --help, -h     Show this help message")
	fmt.Println("  all            Run all tests")
	fmt.Println()
	fmt.Println("Available tests:")
	for name := range TestRegistry {
		if name != "all" {
			fmt.Printf("  %s\n", name)
		}
	}
}

func listTests() {
	fmt.Println("Available visual tests:")
	fmt.Println()
	for name := range TestRegistry {
		fmt.Printf("  - %s\n", name)
	}
}

func runAllTests() {
	tests := []string{
		"tabs",
		"tabs-dynamic",
		"tabs-overflow", 
		"tabs-switch",
		"emoji-width",
		"gear-spacing",
	}
	
	for i, testName := range tests {
		if i > 0 {
			fmt.Println()
			fmt.Println(strings.Repeat("-", 80))
			fmt.Println()
		}
		
		fmt.Printf("Test %d/%d: %s\n", i+1, len(tests), testName)
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println()
		
		if testFunc, exists := TestRegistry[testName]; exists {
			testFunc()
		}
	}
	
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("✅ All %d visual tests completed\n", len(tests))
}

// Test implementations

func testTabs() {
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

func testTabsDynamic() {
	fmt.Println("Tab Bar Dynamic Width Test")
	fmt.Println()
	
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
	
	// Define tabs
	tabs := []components.Tab{
		{Title: "Home", Icon: "🏠", ID: "home"},
		{Title: "Work", Icon: "💼", ID: "work"},
		{Title: "Personal", Icon: "👤", ID: "personal"},
	}
	
	// Test different widths
	widths := []int{40, 60, 80, 100}
	
	for _, width := range widths {
		fmt.Printf("Width: %d\n", width)
		
		tabBar := components.NewTabBar(tabs, width)
		tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
		tabBar.SetActiveIndex(1) // Work active
		
		result := tabBar.Render()
		fmt.Println(result)
		fmt.Println()
	}
}

func testTabsOverflow() {
	fmt.Println("Tab Bar Overflow Test")
	fmt.Println()
	
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
	
	// Test with many tabs
	manyTabs := []components.Tab{
		{Title: "Dashboard", Icon: "📊", ID: "dashboard"},
		{Title: "Analytics", Icon: "📈", ID: "analytics"},
		{Title: "Reports", Icon: "📑", ID: "reports"},
		{Title: "Settings", Icon: "⚙️", ID: "settings"},
		{Title: "Profile", Icon: "👤", ID: "profile"},
		{Title: "Help", Icon: "❓", ID: "help"},
		{Title: "About", Icon: "ℹ️", ID: "about"},
		{Title: "Logout", Icon: "🚪", ID: "logout"},
	}
	
	// Test with narrow width (tabs overflow)
	fmt.Println("Many tabs in narrow space (80 chars):")
	tabBar := components.NewTabBar(manyTabs, 80)
	tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	tabBar.SetActiveIndex(3) // Settings active
	
	result := tabBar.Render()
	fmt.Println(result)
	fmt.Println()
	
	// Test edge cases
	fmt.Println("Edge case: No tabs")
	emptyTabBar := components.NewTabBar([]components.Tab{}, 80)
	emptyTabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	fmt.Println(emptyTabBar.Render())
	fmt.Println()
	
	fmt.Println("Edge case: Single tab")
	singleTab := []components.Tab{{Title: "Only One", Icon: "1️⃣", ID: "one"}}
	singleTabBar := components.NewTabBar(singleTab, 80)
	singleTabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	singleTabBar.SetActiveIndex(0)
	fmt.Println(singleTabBar.Render())
}

func testTabsSwitching() {
	fmt.Println("Tab Bar Theme Switching Test")
	fmt.Println()
	
	// Helper function for tab borders
	tabBorderWithBottom := func(left, middle, right string) lipgloss.Border {
		border := lipgloss.RoundedBorder()
		border.BottomLeft = left
		border.Bottom = middle
		border.BottomRight = right
		return border
	}
	
	// Define tabs
	tabs := []components.Tab{
		{Title: "Code", Icon: "📝", ID: "code"},
		{Title: "Terminal", Icon: "🖥️", ID: "terminal"},
		{Title: "Debug", Icon: "🐛", ID: "debug"},
	}
	
	// Test different themes
	themes := []string{"github-dark", "monokai", "solarized-dark", "one-dark"}
	
	for i, themeName := range themes {
		theme.SetTheme(themeName)
		
		// Get theme colors
		colorBorder := theme.Border()
		colorAccent := theme.Primary()
		colorText := theme.TextPrimary()
		
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
		
		fmt.Printf("Theme: %s\n", themeName)
		
		tabBar := components.NewTabBar(tabs, 80)
		tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
		tabBar.SetActiveIndex(i % len(tabs)) // Rotate active tab
		
		result := tabBar.RenderWithContent(fmt.Sprintf("Content with %s theme", themeName), 5)
		fmt.Println(result)
		fmt.Println()
	}
	
	// Reset to default theme
	theme.SetTheme("github-dark")
}

func testEmojiWidth() {
	// Test various emojis and their display widths
	emojis := []struct {
		emoji string
		desc  string
	}{
		{"🧩", "Puzzle piece"},
		{"👤", "User profile"},
		{"🔧", "Wrench"},
		{"❓", "Question mark"},
		{"⚙️", "Gear"},
		{"🏠", "House"},
		{"💼", "Briefcase"},
		{"📊", "Chart"},
		{"🖥️", "Computer"},
		{"🔍", "Magnifier"},
		{"✅", "Check mark"},
		{"❌", "X mark"},
		{"⭐", "Star"},
		{"🚀", "Rocket"},
		{"🎨", "Palette"},
	}
	
	fmt.Println("Emoji Width Test")
	fmt.Println("================")
	fmt.Println()
	
	// Test in a bordered box to see alignment
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(40)
	
	for _, e := range emojis {
		content := fmt.Sprintf("%s %s", e.emoji, e.desc)
		fmt.Println(boxStyle.Render(content))
	}
	
	// Test in tab-like structures
	fmt.Println("\nIn tab structures:")
	fmt.Println("==================")
	
	tabStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2)
	
	tabs := []string{
		"🧩 Extensions",
		"👤 Profiles",
		"⚙️ Settings",
		"❓ Help",
	}
	
	row := lipgloss.JoinHorizontal(lipgloss.Top, 
		tabStyle.Render(tabs[0]),
		tabStyle.Render(tabs[1]),
		tabStyle.Render(tabs[2]),
		tabStyle.Render(tabs[3]),
	)
	
	fmt.Println(row)
}

func testGearSpacing() {
	fmt.Println("Gear Emoji Spacing Test")
	fmt.Println("======================")
	fmt.Println()
	
	// Test the specific case of gear emoji with and without variation selector
	gearWithVS := "⚙️"  // With variation selector
	gearWithoutVS := "⚙" // Without variation selector
	
	// Create a box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(30)
	
	fmt.Println("With variation selector (⚙️):")
	fmt.Println(boxStyle.Render(gearWithVS + " Settings"))
	fmt.Println()
	
	fmt.Println("Without variation selector (⚙):")
	fmt.Println(boxStyle.Render(gearWithoutVS + " Settings"))
	fmt.Println()
	
	// Test in actual tab context
	colorBorder := theme.Border()
	colorAccent := theme.Primary()
	
	tabStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(0, 2)
	
	activeTabStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Padding(0, 2).
		Bold(true)
	
	fmt.Println("In tab bar context:")
	fmt.Println("==================")
	
	// Test both versions
	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		tabStyle.Render("🧩 Extensions"),
		tabStyle.Render("👤 Profiles"),
		activeTabStyle.Render(gearWithVS+" Settings"),
		tabStyle.Render("❓ Help"),
	)
	
	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		tabStyle.Render("🧩 Extensions"),
		tabStyle.Render("👤 Profiles"),
		activeTabStyle.Render(gearWithoutVS+" Settings"),
		tabStyle.Render("❓ Help"),
	)
	
	fmt.Println("With variation selector:")
	fmt.Println(row1)
	fmt.Println()
	
	fmt.Println("Without variation selector:")
	fmt.Println(row2)
	fmt.Println()
	
	// Test width measurements
	fmt.Println("Width measurements:")
	fmt.Println("==================")
	fmt.Printf("Width of '%s Settings': %d\n", gearWithVS, lipgloss.Width(gearWithVS+" Settings"))
	fmt.Printf("Width of '%s Settings': %d\n", gearWithoutVS, lipgloss.Width(gearWithoutVS+" Settings"))
	fmt.Printf("Width of gear with VS alone: %d\n", lipgloss.Width(gearWithVS))
	fmt.Printf("Width of gear without VS alone: %d\n", lipgloss.Width(gearWithoutVS))
}