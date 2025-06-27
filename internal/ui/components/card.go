package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/muesli/reflow/truncate"
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
func NewCard(width int) *Card {
	// Default styles using theme
	borderColor := theme.Border()
	accentColor := theme.Primary()
	successColor := theme.Success()

	return &Card{
		width:    width,
		metadata: []MetadataItem{},

		normalStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(width - 2), // Account for borders

		selectedStyle: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			Width(width - 2), // Account for borders

		focusedStyle: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			Width(width - 2), // Account for borders

		activeStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Width(width - 2), // Account for borders
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
	// Update all styles with new width (accounting for borders)
	c.normalStyle = c.normalStyle.Width(width - 2)
	c.selectedStyle = c.selectedStyle.Width(width - 2)
	c.focusedStyle = c.focusedStyle.Width(width - 2)
	c.activeStyle = c.activeStyle.Width(width - 2)
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

	// Calculate content width (accounting for border and padding)
	contentWidth := c.width - 7 // 2 for border, 4 for padding, 1 for safety
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Build content
	var content []string

	// Title line with optional subtitle
	titleParts := []string{}
	
	// Add active indicator if active
	if c.active {
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
		
		// Calculate space used by icon and status
		prefixWidth := 0
		for _, part := range titleParts {
			prefixWidth += lipgloss.Width(part) + 1 // +1 for space
		}
		
		// Calculate available width for title
		availableWidth := contentWidth - lipgloss.Width(subtitleText) - prefixWidth - 2 // 2 for spacing between title and subtitle
		
		// Truncate title if needed
		truncatedTitle := truncate.String(c.title, uint(availableWidth))
		titleText := titleStyle.Render(truncatedTitle)
		titleParts = append(titleParts, titleText)
		
		// Join title parts with subtitle
		titleLine := strings.Join(titleParts, " ") + "  " + subtitleText
		content = append(content, titleLine)
	} else {
		// No subtitle, just render title
		titleText := titleStyle.Render(c.title)
		titleParts = append(titleParts, titleText)
		titleLine := strings.Join(titleParts, " ")
		content = append(content, truncate.String(titleLine, uint(contentWidth)))
	}

	// Description (if any)
	if c.description != "" {
		descStyle := lipgloss.NewStyle().Foreground(theme.TextSecondary())
		truncated := truncate.String(c.description, uint(contentWidth))
		content = append(content, descStyle.Render(truncated))
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
			content = append(content, metaStyle.Render(truncate.String(metaLine, uint(contentWidth))))
		}
	}

	// Join content with newlines
	return style.Render(strings.Join(content, "\n"))
}

// RenderCompact renders a more compact version of the card
func (c *Card) RenderCompact() string {
	// Use compact styles (less padding, maintain width adjustment)
	compactStyle := c.normalStyle.Copy().Padding(0, 1).Width(c.width - 2)
	if c.active {
		compactStyle = c.activeStyle.Copy().Padding(0, 1).Width(c.width - 2)
	} else if c.focused {
		compactStyle = c.focusedStyle.Copy().Padding(0, 1).Width(c.width - 2)
	} else if c.selected {
		compactStyle = c.selectedStyle.Copy().Padding(0, 1).Width(c.width - 2)
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
	contentWidth := c.width - 4 // Less padding in compact mode
	
	return compactStyle.Render(truncate.String(titleLine, uint(contentWidth)))
}

// Helper function to pluralize
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}