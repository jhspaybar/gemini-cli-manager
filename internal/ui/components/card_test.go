package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

func TestCard_Basic(t *testing.T) {
	// Initialize theme for tests
	theme.SetTheme("github-dark")

	card := NewCard(50).
		SetTitle("Test Card", "🧪").
		SetDescription("A test description")

	output := card.Render()

	// Check that output contains expected elements
	if !strings.Contains(output, "Test Card") {
		t.Error("Expected card to contain title 'Test Card'")
	}
	if !strings.Contains(output, "🧪") {
		t.Error("Expected card to contain icon")
	}
	if !strings.Contains(output, "A test description") {
		t.Error("Expected card to contain description")
	}
}

func TestCard_States(t *testing.T) {
	theme.SetTheme("github-dark")

	tests := []struct {
		name     string
		setup    func(*Card) *Card
		contains string
	}{
		{
			name: "selected state",
			setup: func(c *Card) *Card {
				return c.SetSelected(true)
			},
			contains: "┏", // Thick border for selected
		},
		{
			name: "active state",
			setup: func(c *Card) *Card {
				return c.SetActive(true)
			},
			contains: "●", // Active indicator
		},
		{
			name: "focused state",
			setup: func(c *Card) *Card {
				return c.SetFocused(true)
			},
			contains: "╔", // Double border for focused
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := NewCard(50).SetTitle("Test", "")
			card = tt.setup(card)
			output := card.Render()

			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected output to contain '%s' for %s", tt.contains, tt.name)
			}
		})
	}
}

func TestCard_Subtitle(t *testing.T) {
	theme.SetTheme("github-dark")

	card := NewCard(60).
		SetTitle("Extension", "🧩").
		SetSubtitle("v1.2.3")

	output := card.Render()

	if !strings.Contains(output, "Extension") {
		t.Error("Expected card to contain title")
	}
	if !strings.Contains(output, "v1.2.3") {
		t.Error("Expected card to contain subtitle")
	}
}

func TestCard_Metadata(t *testing.T) {
	theme.SetTheme("github-dark")

	card := NewCard(50).
		SetTitle("Test", "").
		AddMetadata("Key1", "Value1", "").
		AddMetadata("Key2", "Value2", "🔑")

	output := card.Render()

	if !strings.Contains(output, "Key1: Value1") {
		t.Error("Expected card to contain metadata without icon")
	}
	if !strings.Contains(output, "🔑 Value2") {
		t.Error("Expected card to contain metadata with icon")
	}
}

func TestCard_WidthResize(t *testing.T) {
	theme.SetTheme("github-dark")

	card := NewCard(40).
		SetTitle("Resizable Card", "📏").
		SetDescription("This card can be resized")

	// Test initial width
	output1 := card.Render()
	width1 := calculateMaxLineWidth(output1)
	if width1 > 40 {
		t.Errorf("Card width exceeded limit: got %d, want <= 40", width1)
	}

	// Resize and test again
	card.SetWidth(60)
	output2 := card.Render()
	width2 := calculateMaxLineWidth(output2)
	if width2 > 60 {
		t.Errorf("Card width exceeded limit after resize: got %d, want <= 60", width2)
	}
}

func TestCard_Truncation(t *testing.T) {
	theme.SetTheme("github-dark")

	longTitle := "This is an extremely long title that should be truncated to fit within the card"
	longDesc := "This is a very long description that contains lots of text and should be truncated properly to ensure it doesn't break the layout of the card component"

	card := NewCard(40).
		SetTitle(longTitle, "📝").
		SetSubtitle("v1.0").
		SetDescription(longDesc)

	output := card.Render()

	// Check that output is properly bounded
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		width := lipgloss.Width(line)
		if width > 40 {
			t.Errorf("Line %d exceeded width limit: got %d, want <= 40", i, width)
		}
	}

	// Should contain truncated content
	if strings.Contains(output, longTitle) {
		t.Error("Long title should have been truncated")
	}
	if strings.Contains(output, longDesc) {
		t.Error("Long description should have been truncated")
	}
}

func TestCard_CompactRender(t *testing.T) {
	theme.SetTheme("github-dark")

	card := NewCard(30).
		SetTitle("Compact", "⚡").
		SetDescription("Should not appear in compact mode")

	compact := card.RenderCompact()

	// Compact should contain title and icon
	if !strings.Contains(compact, "Compact") {
		t.Error("Compact render should contain title")
	}
	if !strings.Contains(compact, "⚡") {
		t.Error("Compact render should contain icon")
	}

	// Compact should NOT contain description
	if strings.Contains(compact, "Should not appear") {
		t.Error("Compact render should not contain description")
	}

	// Compact should be smaller (less padding)
	normalHeight := len(strings.Split(card.Render(), "\n"))
	compactHeight := len(strings.Split(compact, "\n"))
	if compactHeight >= normalHeight {
		t.Error("Compact render should have fewer lines than normal render")
	}
}

func TestCard_EmptyContent(t *testing.T) {
	theme.SetTheme("github-dark")

	card := NewCard(40)
	output := card.Render()

	// Should still render a box
	if !strings.Contains(output, "╭") {
		t.Error("Empty card should still render borders")
	}

	// Check it has some height
	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		t.Error("Empty card should have at least 3 lines (top, content, bottom)")
	}
}

// Helper function to calculate the maximum line width in a multi-line string
func calculateMaxLineWidth(s string) int {
	lines := strings.Split(s, "\n")
	maxWidth := 0
	for _, line := range lines {
		width := lipgloss.Width(line)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}