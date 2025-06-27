package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func main() {
	// Test emoji widths
	emojis := []struct {
		name  string
		emoji string
	}{
		{"Extensions", "🧩"},
		{"Profiles", "👤"},
		{"Settings", "⚙️"},
		{"Help", "❓"},
	}
	
	fmt.Println("Emoji width analysis:")
	fmt.Println("=====================")
	for _, e := range emojis {
		emojiWidth := runewidth.StringWidth(e.emoji)
		textWidth := runewidth.StringWidth(e.name)
		combined := fmt.Sprintf("%s %s", e.emoji, e.name)
		combinedWidth := runewidth.StringWidth(combined)
		lipglossWidth := lipgloss.Width(combined)
		
		fmt.Printf("%s (%s):\n", e.name, e.emoji)
		fmt.Printf("  Emoji width: %d\n", emojiWidth)
		fmt.Printf("  Text width: %d\n", textWidth)
		fmt.Printf("  Combined width (runewidth): %d\n", combinedWidth)
		fmt.Printf("  Combined width (lipgloss): %d\n", lipglossWidth)
		fmt.Printf("  Expected: %d (emoji + space + text)\n", emojiWidth + 1 + textWidth)
		fmt.Println()
	}
	
	// Test with actual tab rendering
	fmt.Println("\nTab rendering test:")
	fmt.Println("===================")
	
	border := lipgloss.RoundedBorder()
	style := lipgloss.NewStyle().
		Border(border).
		Padding(0, 2)
	
	for _, e := range emojis {
		content := fmt.Sprintf("%s %s", e.emoji, e.name)
		rendered := style.Render(content)
		fmt.Printf("%-12s: Width=%2d | %s\n", e.name, lipgloss.Width(rendered), rendered)
	}
}