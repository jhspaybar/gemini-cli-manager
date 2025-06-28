package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// Card represents a reusable card UI component
type Card struct {
	width       int
	title       string
	titleIcon   string
	subtitle    string
	description string
	metadata    []MetadataItem
	selected    bool
	focused     bool
	active      bool // For states like "active profile"

	// Styles
	normalStyle   lipgloss.Style
	selectedStyle lipgloss.Style
	focusedStyle  lipgloss.Style
	activeStyle   lipgloss.Style
}

// MetadataItem represents a key-value pair in the card
type MetadataItem struct {
	Key   string
	Value string
	Icon  string // Optional icon
}

// NewCard creates a new card component
// The width parameter represents the maximum outer width the card can occupy
func NewCard(width int) *Card {
	// Default styles using theme
	borderColor := theme.Border()
	accentColor := theme.Primary()
	successColor := theme.Success()
	
	// Calculate inner width by accounting for borders (2) and padding (4)
	// This ensures the total card width doesn't exceed the given width
	borderWidth := 2  // 1 char on each side
	paddingWidth := 4 // 2 chars on each side
	innerWidth := width - borderWidth - paddingWidth
	if innerWidth < 10 {
		innerWidth = 10 // Minimum width to prevent display issues
	}

	return &Card{
		width:    width,
		metadata: []MetadataItem{},

		normalStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(innerWidth),

		selectedStyle: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(accentColor).
			Background(theme.Selection()).    // Add subtle background
			Padding(1, 2).
			Width(innerWidth),

		focusedStyle: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			Width(innerWidth),

		activeStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Width(innerWidth),
	}
}

// SetTitle sets the card title with optional icon
func (c *Card) SetTitle(title, icon string) *Card {
	c.title = title
	c.titleIcon = icon
	return c
}

// SetSubtitle sets a subtitle (e.g., version number)
func (c *Card) SetSubtitle(subtitle string) *Card {
	c.subtitle = subtitle
	return c
}

// SetDescription sets the card description
func (c *Card) SetDescription(description string) *Card {
	c.description = description
	return c
}

// AddMetadata adds a key-value metadata item with optional icon
func (c *Card) AddMetadata(key, value, icon string) *Card {
	c.metadata = append(c.metadata, MetadataItem{Key: key, Value: value, Icon: icon})
	return c
}

// SetSelected sets the selected state
func (c *Card) SetSelected(selected bool) *Card {
	c.selected = selected
	return c
}

// SetFocused sets the focused state
func (c *Card) SetFocused(focused bool) *Card {
	c.focused = focused
	return c
}

// SetActive sets the active state (e.g., for active profile)
func (c *Card) SetActive(active bool) *Card {
	c.active = active
	return c
}

// SetStyles allows custom style configuration
func (c *Card) SetStyles(normal, selected, focused, active lipgloss.Style) *Card {
	c.normalStyle = normal
	c.selectedStyle = selected
	c.focusedStyle = focused
	c.activeStyle = active
	return c
}

// SetWidth updates the card width
func (c *Card) SetWidth(width int) *Card {
	c.width = width
	// Calculate inner width accounting for borders and padding
	borderWidth := 2  // 1 char on each side
	paddingWidth := 4 // 2 chars on each side
	innerWidth := width - borderWidth - paddingWidth
	if innerWidth < 10 {
		innerWidth = 10
	}
	// Update all styles with new inner width
	c.normalStyle = c.normalStyle.Width(innerWidth)
	c.selectedStyle = c.selectedStyle.Width(innerWidth)
	c.focusedStyle = c.focusedStyle.Width(innerWidth)
	c.activeStyle = c.activeStyle.Width(innerWidth)
	return c
}

// Render produces the card's visual output
func (c *Card) Render() string {
	// Select appropriate style
	style := c.normalStyle
	if c.active {
		style = c.activeStyle
	} else if c.focused {
		style = c.focusedStyle
	} else if c.selected {
		style = c.selectedStyle
	}

	// Content will be automatically constrained by MaxWidth on the style

	// Build content
	var content []string

	// Title line with optional subtitle
	titleParts := []string{}
	
	// Add selection indicator
	if c.selected {
		titleParts = append(titleParts, "▶")
	} else if c.active {
		titleParts = append(titleParts, "●")
	}
	
	// Add icon if present
	if c.titleIcon != "" {
		titleParts = append(titleParts, c.titleIcon)
	}
	
	// Add title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.TextPrimary())
	if c.selected || c.focused {
		titleStyle = titleStyle.Foreground(theme.Primary())
	} else if c.active {
		titleStyle = titleStyle.Foreground(theme.Success())
	}
	
	// Handle subtitle (like version)
	if c.subtitle != "" {
		subtitleStyle := lipgloss.NewStyle().Foreground(theme.TextMuted())
		subtitleText := subtitleStyle.Render(c.subtitle)
		
		// Build the title parts
		prefix := strings.Join(titleParts, " ")
		if prefix != "" {
			prefix += " "
		}
		
		// Create a title line with subtitle
		titleText := titleStyle.Render(c.title)
		fullTitle := prefix + titleText + "  " + subtitleText
		content = append(content, fullTitle)
	} else {
		// No subtitle, just render title
		titleText := titleStyle.Render(c.title)
		titleParts = append(titleParts, titleText)
		titleLine := strings.Join(titleParts, " ")
		content = append(content, titleLine)
	}

	// Description (if any)
	if c.description != "" {
		descStyle := lipgloss.NewStyle().Foreground(theme.TextSecondary())
		content = append(content, descStyle.Render(c.description))
	}

	// Metadata section
	if len(c.metadata) > 0 {
		metaStyle := lipgloss.NewStyle().Foreground(theme.Primary())
		for _, meta := range c.metadata {
			var metaLine string
			if meta.Icon != "" {
				metaLine = fmt.Sprintf("%s %s", meta.Icon, meta.Value)
			} else {
				metaLine = fmt.Sprintf("%s: %s", meta.Key, meta.Value)
			}
			content = append(content, metaStyle.Render(metaLine))
		}
	}

	// Join content with newlines
	return style.Render(strings.Join(content, "\n"))
}

// RenderCompact renders a more compact version of the card
func (c *Card) RenderCompact() string {
	// Calculate inner width for compact mode (less padding)
	borderWidth := 2  // 1 char on each side
	paddingWidth := 2 // 1 char on each side (compact padding)
	innerWidth := c.width - borderWidth - paddingWidth
	if innerWidth < 10 {
		innerWidth = 10
	}
	
	// Use compact styles (less padding)
	compactStyle := c.normalStyle.Copy().Padding(0, 1).Width(innerWidth)
	if c.active {
		compactStyle = c.activeStyle.Copy().Padding(0, 1).Width(innerWidth)
	} else if c.focused {
		compactStyle = c.focusedStyle.Copy().Padding(0, 1).Width(innerWidth)
	} else if c.selected {
		compactStyle = c.selectedStyle.Copy().Padding(0, 1).Width(innerWidth)
	}

	// Build compact title line
	parts := []string{}
	if c.active {
		parts = append(parts, "●")
	}
	if c.titleIcon != "" {
		parts = append(parts, c.titleIcon)
	}
	parts = append(parts, c.title)
	
	titleLine := strings.Join(parts, " ")
	
	return compactStyle.Render(titleLine)
}

// Helper function to pluralize
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}