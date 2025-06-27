// Package adhoc provides visual tests for UI components
// These tests output visual representations to help verify component rendering
// Run with: go test -v ./test/adhoc -run TestVisual
package adhoc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

// TestVisualTabs tests the tab component rendering
func TestVisualTabs(t *testing.T) {
	if !shouldRunVisualTests() {
		t.Skip("Skipping visual test. Set VISUAL_TESTS=true to run.")
	}
	
	output := captureOutput(func() {
		testTabs()
	})
	
	// Basic assertions
	if !strings.Contains(output, "Extensions") {
		t.Error("Expected to see 'Extensions' in tab output")
	}
	if !strings.Contains(output, "🧩") {
		t.Error("Expected to see extension icon in output")
	}
	
	// Print for visual inspection when running with -v
	t.Log("\n" + output)
}

// TestVisualTabsDynamic tests tabs at different widths
func TestVisualTabsDynamic(t *testing.T) {
	if !shouldRunVisualTests() {
		t.Skip("Skipping visual test. Set VISUAL_TESTS=true to run.")
	}
	
	output := captureOutput(func() {
		testTabsDynamic()
	})
	
	// Verify we see different widths
	if !strings.Contains(output, "Width: 40") {
		t.Error("Expected to see Width: 40")
	}
	if !strings.Contains(output, "Width: 100") {
		t.Error("Expected to see Width: 100")
	}
	
	t.Log("\n" + output)
}

// TestVisualEmojiWidth tests emoji rendering
func TestVisualEmojiWidth(t *testing.T) {
	if !shouldRunVisualTests() {
		t.Skip("Skipping visual test. Set VISUAL_TESTS=true to run.")
	}
	
	output := captureOutput(func() {
		testEmojiWidth()
	})
	
	// Check for various emojis
	emojis := []string{"🧩", "👤", "🔧", "❓", "⚙️"}
	for _, emoji := range emojis {
		if !strings.Contains(output, emoji) {
			t.Errorf("Expected to see emoji %s in output", emoji)
		}
	}
	
	t.Log("\n" + output)
}

// TestVisualAll runs all visual tests
func TestVisualAll(t *testing.T) {
	if !shouldRunVisualTests() {
		t.Skip("Skipping visual test. Set VISUAL_TESTS=true to run.")
	}
	
	tests := []struct {
		name string
		fn   func()
	}{
		{"Tabs", testTabs},
		{"TabsDynamic", testTabsDynamic},
		{"TabsOverflow", testTabsOverflow},
		{"TabsSwitch", testTabsSwitching},
		{"EmojiWidth", testEmojiWidth},
		{"GearSpacing", testGearSpacing},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(test.fn)
			t.Log("\n" + output)
		})
	}
}

// Helper function to check if visual tests should run
func shouldRunVisualTests() bool {
	return os.Getenv("VISUAL_TESTS") == "true"
}

// Helper function to capture output
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	fn()
	
	w.Close()
	os.Stdout = old
	
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Test implementations (same as in visual_tests.go)

func testTabs() {
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

func testTabsDynamic() {
	fmt.Println("Tab Bar Dynamic Width Test")
	fmt.Println()
	
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