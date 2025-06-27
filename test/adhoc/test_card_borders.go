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

	fmt.Println("Testing Card Borders - Ensuring no cutoffs")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Test 1: Card at terminal width
	fmt.Println("Test 1: Card using full terminal width (80)")
	fmt.Println(strings.Repeat("-", 80))
	
	card1 := components.NewCard(80).
		SetTitle("Full Width Card", "🖥️").
		SetSubtitle("v1.0.0").
		SetDescription("This card should use the full terminal width without cutoffs")
	
	fmt.Println(card1.Render())
	
	// Test 2: Card with margin simulation (what view.go should do)
	fmt.Println("\n\nTest 2: Card with proper margins (simulating view.go)")
	fmt.Println(strings.Repeat("-", 80))
	
	// This simulates what should happen in view.go
	containerStyle := lipgloss.NewStyle().
		Margin(0, 2). // 2 chars margin on each side
		MaxWidth(80)
	
	card2 := components.NewCard(76). // 80 - 4 for margins
		SetTitle("Card with Container Margins", "📦").
		SetSubtitle("v2.0.0").
		SetDescription("This card is inside a container with margins")
	
	fmt.Println(containerStyle.Render(card2.Render()))
	
	// Test 3: Multiple cards side by side
	fmt.Println("\n\nTest 3: Two cards side by side")
	fmt.Println(strings.Repeat("-", 80))
	
	leftCard := components.NewCard(38). // Half width minus gap
		SetTitle("Left Card", "⬅️").
		SetDescription("Left side content")
	
	rightCard := components.NewCard(38).
		SetTitle("Right Card", "➡️").
		SetDescription("Right side content")
	
	// Join with gap
	sideBySide := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftCard.Render(),
		"  ", // 2 char gap
		rightCard.Render(),
	)
	
	fmt.Println(sideBySide)
	
	// Test 4: Card in a narrow space
	fmt.Println("\n\nTest 4: Card in narrow space (40 chars)")
	fmt.Println(strings.Repeat("-", 80))
	
	narrowCard := components.NewCard(40).
		SetTitle("Narrow Card", "📐").
		SetDescription("Content in a narrow space")
	
	// Center it
	centered := lipgloss.Place(80, 5, lipgloss.Center, lipgloss.Center, narrowCard.Render())
	fmt.Println(centered)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Visual check:")
	fmt.Println("1. All borders should be complete (no cutoffs)")
	fmt.Println("2. Text should wrap properly within cards")
	fmt.Println("3. Cards should respect their container boundaries")
}