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

	fmt.Println("Width Calculation Debug")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Simulate what view.go is doing
	terminalWidth := 80
	
	// renderAppCard adds padding
	horizontalPadding := 3
	contentWidth := terminalWidth - (horizontalPadding * 2) // 74
	
	fmt.Printf("Terminal width: %d\n", terminalWidth)
	fmt.Printf("After renderAppCard padding (3 each side): %d\n", contentWidth)
	
	// renderContent adds more padding
	contentPadding := 3
	availableWidth := contentWidth - (contentPadding * 2) // 68
	
	fmt.Printf("After renderContent padding (3 each side): %d\n", availableWidth)
	fmt.Println()
	
	// Create cards with different widths
	fmt.Println("Test 1: Card with availableWidth (68)")
	fmt.Println(strings.Repeat("-", 80))
	card1 := components.NewCard(availableWidth).
		SetTitle("Extension with MCP server configuration for testing", "🧩").
		SetSubtitle("v2.1.0").
		SetDescription("Extension with MCP server configuration for testing").
		AddMetadata("Servers", "2 servers", "⚡")
	
	fmt.Println(card1.Render())
	
	// Test what the card looks like with padding applied
	fmt.Println("\n\nTest 2: Same card with renderContent padding applied")
	fmt.Println(strings.Repeat("-", 80))
	
	contentContainer := lipgloss.NewStyle().
		Padding(2, 3).
		MaxWidth(contentWidth)
	
	fmt.Println(contentContainer.Render(card1.Render()))
	
	// Test with full terminal width
	fmt.Println("\n\nTest 3: Card with terminal width (80)")
	fmt.Println(strings.Repeat("-", 80))
	card2 := components.NewCard(terminalWidth).
		SetTitle("Extension with MCP server configuration for testing", "🧩").
		SetSubtitle("v2.1.0").
		SetDescription("Extension with MCP server configuration for testing").
		AddMetadata("Servers", "2 servers", "⚡")
	
	fmt.Println(card2.Render())
	
	// Show the issue with double padding
	fmt.Println("\n\nTest 4: Card with availableWidth inside container (double padding)")
	fmt.Println(strings.Repeat("-", 80))
	
	// This simulates the current behavior
	appContainer := lipgloss.NewStyle().Padding(0, horizontalPadding)
	fullRender := appContainer.Render(contentContainer.Render(card1.Render()))
	
	fmt.Println(fullRender)
}