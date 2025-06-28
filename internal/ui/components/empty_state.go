package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// EmptyState represents a component for displaying empty/no-data states
type EmptyState struct {
	width       int
	icon        string
	title       string
	description string
	action      string
	centered    bool
	style       lipgloss.Style
	iconStyle   lipgloss.Style
	titleStyle  lipgloss.Style
	descStyle   lipgloss.Style
	actionStyle lipgloss.Style
}

// NewEmptyState creates a new empty state component
func NewEmptyState(width int) *EmptyState {
	return &EmptyState{
		width:    width,
		centered: true,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border()).
			Padding(2, 4).
			Align(lipgloss.Center),
		iconStyle: lipgloss.NewStyle().
			Foreground(theme.TextSecondary()),
		titleStyle: lipgloss.NewStyle().
			Foreground(theme.TextPrimary()).
			Bold(true),
		descStyle: lipgloss.NewStyle().
			Foreground(theme.TextSecondary()),
		actionStyle: lipgloss.NewStyle().
			Foreground(theme.Primary()).
			Italic(true),
	}
}

// SetIcon sets the emoji/icon for the empty state
func (es *EmptyState) SetIcon(icon string) *EmptyState {
	es.icon = icon
	return es
}

// SetTitle sets the main message
func (es *EmptyState) SetTitle(title string) *EmptyState {
	es.title = title
	return es
}

// SetDescription sets the descriptive text
func (es *EmptyState) SetDescription(description string) *EmptyState {
	es.description = description
	return es
}

// SetAction sets the action hint (e.g., "Press 'n' to create")
func (es *EmptyState) SetAction(action string) *EmptyState {
	es.action = action
	return es
}

// SetCentered controls whether the empty state is centered
func (es *EmptyState) SetCentered(centered bool) *EmptyState {
	es.centered = centered
	return es
}

// SetWidth updates the component width
func (es *EmptyState) SetWidth(width int) *EmptyState {
	es.width = width
	return es
}

// SetStyles allows customization of all styles
func (es *EmptyState) SetStyles(boxStyle, iconStyle, titleStyle, descStyle, actionStyle lipgloss.Style) *EmptyState {
	es.style = boxStyle
	es.iconStyle = iconStyle
	es.titleStyle = titleStyle
	es.descStyle = descStyle
	es.actionStyle = actionStyle
	return es
}

// Render renders the empty state
func (es *EmptyState) Render() string {
	// Build content parts
	var parts []string
	
	if es.icon != "" {
		parts = append(parts, es.iconStyle.Render(es.icon))
		parts = append(parts, "") // Empty line after icon
	}
	
	if es.title != "" {
		parts = append(parts, es.titleStyle.Render(es.title))
	}
	
	if es.description != "" {
		parts = append(parts, es.descStyle.Render(es.description))
	}
	
	if es.action != "" {
		parts = append(parts, es.actionStyle.Render(es.action))
	}
	
	// Join parts vertically with center alignment
	content := lipgloss.JoinVertical(lipgloss.Center, parts...)
	
	// Calculate box width - constrain to reasonable size
	boxWidth := es.width
	if boxWidth > 60 {
		boxWidth = 60 // Max width for readability
	}
	if boxWidth < 30 {
		boxWidth = 30 // Min width
	}
	
	// Apply box styling
	// Use MaxWidth to constrain content, let lipgloss handle total width with borders
	box := es.style.
		MaxWidth(boxWidth).
		Render(content)
	
	// Center the box if requested
	if es.centered {
		// Get actual rendered width including borders
		actualWidth := lipgloss.Width(box)
		if es.width > actualWidth {
			return lipgloss.Place(es.width, lipgloss.Height(box), lipgloss.Center, lipgloss.Center, box)
		}
	}
	
	return box
}

// Preset empty states for common scenarios

// NoItemsFound creates an empty state for search with no results
func (es *EmptyState) NoItemsFound() *EmptyState {
	return es.
		SetIcon("🔍").
		SetTitle("No items found").
		SetDescription("Try adjusting your search or filters")
}

// NoData creates a generic no data empty state
func (es *EmptyState) NoData(itemType string) *EmptyState {
	return es.
		SetIcon("📭").
		SetTitle("No " + itemType).
		SetDescription("Nothing to display yet")
}

// ComingSoon creates an empty state for features not yet implemented
func (es *EmptyState) ComingSoon() *EmptyState {
	return es.
		SetIcon("🚧").
		SetTitle("Coming Soon").
		SetDescription("This feature is under construction")
}