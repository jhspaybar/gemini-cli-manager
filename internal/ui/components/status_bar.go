package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// StatusBar represents a three-section status bar component
type StatusBar struct {
	width int
	
	// Content sections
	leftContent   string
	middleContent string
	rightContent  string
	
	// Layout proportions (must add up to reasonable total)
	leftProportion   int // default 3
	middleProportion int // default 2  
	rightProportion  int // default 2
	
	// Styling
	style       lipgloss.Style
	errorStyle  lipgloss.Style
	keyStyle    lipgloss.Style
	keyDescStyle lipgloss.Style
}

// StatusItem represents a key-value item for the status bar
type StatusItem struct {
	Icon  string
	Label string
	Value string
}

// KeyBinding represents a keyboard shortcut
type KeyBinding struct {
	Key         string
	Description string
}

// ErrorMessage represents an error or info message
type ErrorMessage struct {
	Type    ErrorType
	Icon    string
	Message string
	Details string
}

// ErrorType represents the type of error message
type ErrorType int

const (
	ErrorTypeError ErrorType = iota
	ErrorTypeInfo
	ErrorTypeWarning
)

// NewStatusBar creates a new status bar component
func NewStatusBar(width int) *StatusBar {
	// Default styles using theme
	borderColor := theme.Border()
	textDim := theme.TextSecondary()
	accent := theme.Primary()
	errorColor := theme.Error()
	
	return &StatusBar{
		width:            width,
		leftProportion:   2,
		middleProportion: 2,
		rightProportion:  3,
		
		style: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(borderColor).
			Foreground(textDim).
			Padding(0, 1),
			
		errorStyle: lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true),
			
		keyStyle: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),
			
		keyDescStyle: lipgloss.NewStyle().
			Foreground(textDim),
	}
}

// SetProportions sets the width proportions for the three sections
func (sb *StatusBar) SetProportions(left, middle, right int) *StatusBar {
	sb.leftProportion = left
	sb.middleProportion = middle
	sb.rightProportion = right
	return sb
}

// SetLeftContent sets the left section content directly
func (sb *StatusBar) SetLeftContent(content string) *StatusBar {
	sb.leftContent = content
	return sb
}

// SetLeftItems sets the left section with status items
func (sb *StatusBar) SetLeftItems(items []StatusItem) *StatusBar {
	var parts []string
	for _, item := range items {
		if item.Icon != "" && item.Value != "" {
			parts = append(parts, fmt.Sprintf("%s %s", item.Icon, item.Value))
		} else if item.Icon != "" && item.Label != "" {
			parts = append(parts, fmt.Sprintf("%s %s", item.Icon, item.Label))
		} else if item.Label != "" && item.Value != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", item.Label, item.Value))
		}
	}
	sb.leftContent = strings.Join(parts, " • ")
	return sb
}

// SetMiddleContent sets the middle section content directly
func (sb *StatusBar) SetMiddleContent(content string) *StatusBar {
	sb.middleContent = content
	return sb
}

// SetErrorMessage sets the middle section with an error/info message
func (sb *StatusBar) SetErrorMessage(msg ErrorMessage) *StatusBar {
	if msg.Message == "" {
		sb.middleContent = ""
		return sb
	}
	
	var content string
	var style lipgloss.Style
	
	switch msg.Type {
	case ErrorTypeInfo:
		style = lipgloss.NewStyle().Foreground(theme.Primary())
		icon := msg.Icon
		if icon == "" {
			icon = "ℹ️"
		}
		content = fmt.Sprintf(" %s %s ", icon, msg.Message)
		if msg.Details != "" {
			content += lipgloss.NewStyle().Foreground(theme.TextSecondary()).Render(fmt.Sprintf(" - %s", msg.Details))
		}
	case ErrorTypeWarning:
		style = lipgloss.NewStyle().Foreground(theme.Warning())
		icon := msg.Icon
		if icon == "" {
			icon = "⚠️"
		}
		content = fmt.Sprintf(" %s %s ", icon, msg.Message)
	case ErrorTypeError:
		style = sb.errorStyle
		icon := msg.Icon
		if icon == "" {
			icon = "❌"
		}
		content = fmt.Sprintf(" %s %s ", icon, msg.Message)
	}
	
	sb.middleContent = style.Render(content)
	return sb
}

// SetRightContent sets the right section content directly  
func (sb *StatusBar) SetRightContent(content string) *StatusBar {
	sb.rightContent = content
	return sb
}

// SetKeyBindings sets the right section with key bindings
func (sb *StatusBar) SetKeyBindings(bindings []KeyBinding) *StatusBar {
	var parts []string
	for _, binding := range bindings {
		part := sb.keyStyle.Render(binding.Key) + " " + sb.keyDescStyle.Render(binding.Description)
		parts = append(parts, part)
	}
	sb.rightContent = strings.Join(parts, " • ")
	return sb
}

// SetWidth updates the status bar width
func (sb *StatusBar) SetWidth(width int) *StatusBar {
	sb.width = width
	return sb
}

// SetStyles allows customization of styles
func (sb *StatusBar) SetStyles(statusStyle, errorStyle, keyStyle, keyDescStyle lipgloss.Style) *StatusBar {
	sb.style = statusStyle
	sb.errorStyle = errorStyle
	sb.keyStyle = keyStyle
	sb.keyDescStyle = keyDescStyle
	return sb
}

// Render renders the status bar
func (sb *StatusBar) Render() string {
	if sb.width <= 0 {
		return ""
	}
	
	// Calculate section widths based on proportions
	total := sb.leftProportion + sb.middleProportion + sb.rightProportion
	leftWidth := (sb.width * sb.leftProportion) / total
	middleWidth := (sb.width * sb.middleProportion) / total
	rightWidth := sb.width - leftWidth - middleWidth // Ensure exact fit
	
	// Create sections with calculated widths
	leftSection := lipgloss.NewStyle().Width(leftWidth).Render(sb.leftContent)
	
	var middleSection string
	if sb.middleContent != "" {
		middleSection = lipgloss.NewStyle().Width(middleWidth).Align(lipgloss.Center).Render(sb.middleContent)
	} else {
		middleSection = lipgloss.NewStyle().Width(middleWidth).Render("")
	}
	
	var rightSection string
	if sb.rightContent != "" {
		rightSection = lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Right).Render(sb.rightContent)
	} else {
		rightSection = lipgloss.NewStyle().Width(rightWidth).Render("")
	}
	
	// Join sections horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, middleSection, rightSection)
	
	// Apply status bar styling and return
	return sb.style.Width(sb.width).Render(content)
}

// RenderContent renders just the content without the status bar border styling
func (sb *StatusBar) RenderContent() string {
	if sb.width <= 0 {
		return ""
	}
	
	// Calculate section widths based on proportions
	total := sb.leftProportion + sb.middleProportion + sb.rightProportion
	leftWidth := (sb.width * sb.leftProportion) / total
	middleWidth := (sb.width * sb.middleProportion) / total
	rightWidth := sb.width - leftWidth - middleWidth // Ensure exact fit
	
	// Create sections with calculated widths
	leftSection := lipgloss.NewStyle().Width(leftWidth).Render(sb.leftContent)
	
	var middleSection string
	if sb.middleContent != "" {
		middleSection = lipgloss.NewStyle().Width(middleWidth).Align(lipgloss.Center).Render(sb.middleContent)
	} else {
		middleSection = lipgloss.NewStyle().Width(middleWidth).Render("")
	}
	
	var rightSection string
	if sb.rightContent != "" {
		rightSection = lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Right).Render(sb.rightContent)
	} else {
		rightSection = lipgloss.NewStyle().Width(rightWidth).Render("")
	}
	
	// Join sections horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, leftSection, middleSection, rightSection)
}

// Clear clears all content sections
func (sb *StatusBar) Clear() *StatusBar {
	sb.leftContent = ""
	sb.middleContent = ""
	sb.rightContent = ""
	return sb
}

// Helper function to create common key bindings
func CommonKeyBindings() []KeyBinding {
	return []KeyBinding{
		{"Tab", "Switch"},
		{"L", "Launch"},
		{"?", "Help"},
		{"q", "Quit"},
	}
}

// Helper function to create profile status items
func ProfileStatusItems(profileName string, enabledCount, totalCount int) []StatusItem {
	items := []StatusItem{
		{"👤", "", profileName},
		{"🧩", "", fmt.Sprintf("%d/%d", enabledCount, totalCount)},
	}
	
	if profileName == "" {
		items[0] = StatusItem{"👤", "", "No Profile"}
	}
	
	return items
}