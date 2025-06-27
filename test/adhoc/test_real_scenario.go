package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

func main() {
	theme.SetTheme("github-dark")

	// Simulate a realistic terminal width
	terminalWidth := 120
	
	fmt.Printf("Simulating real app scenario with terminal width: %d\n", terminalWidth)
	fmt.Println(strings.Repeat("=", terminalWidth))
	fmt.Println()
	
	// Step 1: renderAppCard
	appPadding := lipgloss.NewStyle().Padding(2, 3) // vertical: 2, horizontal: 3
	contentWidth := terminalWidth
	contentHeight := 40
	
	// Step 2: renderTabsWithContent passes width to renderContent
	// Step 3: renderContent applies its own padding
	contentContainer := lipgloss.NewStyle().
		Padding(2, 3).
		MaxWidth(contentWidth).
		MaxHeight(contentHeight)
	
	// Calculate available space after content padding
	availableWidth := contentWidth - 6  // 3 padding on each side
	// availableHeight := contentHeight - 4 // 2 padding on top and bottom
	
	fmt.Printf("After content padding, available width: %d\n", availableWidth)
	
	// Step 4: renderExtensions gets availableWidth
	// Step 5: Cards are created with this width
	
	// Create two cards as they appear in the extension list
	card1 := components.NewCard(availableWidth).
		SetTitle("mcp-extension", "🧩").
		SetSubtitle("v2.1.0").
		SetDescription("Extension with MCP server configuration for testing").
		AddMetadata("Servers", "2 servers", "⚡").
		SetSelected(true)
	
	card2 := components.NewCard(availableWidth).
		SetTitle("simple-extension", "🧩").
		SetSubtitle("v1.0.0").
		SetDescription("A minimal extension for testing basic installation")
	
	// Build the extension list content
	var lines []string
	lines = append(lines, "Extensions")
	lines = append(lines, "")
	lines = append(lines, "2 extensions found")
	lines = append(lines, "")
	lines = append(lines, card1.Render())
	lines = append(lines, "")
	lines = append(lines, card2.Render())
	
	extensionContent := strings.Join(lines, "\n")
	
	// Apply the content container styling
	styledContent := contentContainer.Render(extensionContent)
	
	// Apply the app padding
	finalOutput := appPadding.Render(styledContent)
	
	fmt.Println("Final rendered output:")
	fmt.Println(finalOutput)
	
	// Show what's happening with a single card at each step
	fmt.Println("\n\nStep-by-step single card rendering:")
	fmt.Println(strings.Repeat("-", terminalWidth))
	
	singleCard := components.NewCard(availableWidth).
		SetTitle("Debug Card", "🔍").
		SetSubtitle("v1.0.0").
		SetDescription("Testing width calculations")
	
	fmt.Printf("\n1. Raw card (width %d):\n", availableWidth)
	fmt.Println(singleCard.Render())
	
	fmt.Println("\n2. With content container padding:")
	fmt.Println(contentContainer.Render(singleCard.Render()))
	
	fmt.Println("\n3. With app padding:")
	fmt.Println(appPadding.Render(contentContainer.Render(singleCard.Render())))
}