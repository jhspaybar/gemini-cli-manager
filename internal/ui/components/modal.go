package components

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// Modal represents a reusable modal container component
type Modal struct {
	title        string
	icon         string
	content      string
	footer       string
	width        int
	maxWidth     int
	windowWidth  int
	windowHeight int
	borderColor  lipgloss.TerminalColor
	titleStyle   lipgloss.Style
	contentStyle lipgloss.Style
	footerStyle  lipgloss.Style
}

// NewModal creates a new modal with default styling
func NewModal(windowWidth, windowHeight int) *Modal {
	return &Modal{
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		width:        60,
		maxWidth:     80,
		borderColor:  theme.Border(),
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary()).
			MarginBottom(1),
		contentStyle: lipgloss.NewStyle().
			Foreground(theme.TextPrimary()),
		footerStyle: lipgloss.NewStyle().
			Foreground(theme.TextSecondary()).
			MarginTop(1),
	}
}

// SetTitle sets the modal title with optional icon
func (m *Modal) SetTitle(title, icon string) *Modal {
	m.title = title
	m.icon = icon
	return m
}

// SetContent sets the main content of the modal
func (m *Modal) SetContent(content string) *Modal {
	m.content = content
	return m
}

// SetFooter sets the footer content (usually help text or actions)
func (m *Modal) SetFooter(footer string) *Modal {
	m.footer = footer
	return m
}

// SetWidth sets the modal width
func (m *Modal) SetWidth(width int) *Modal {
	m.width = width
	return m
}

// SetMaxWidth sets the maximum width
func (m *Modal) SetMaxWidth(maxWidth int) *Modal {
	m.maxWidth = maxWidth
	return m
}

// SetBorderColor sets the border color
func (m *Modal) SetBorderColor(color lipgloss.TerminalColor) *Modal {
	m.borderColor = color
	return m
}

// SetTitleStyle sets a custom title style
func (m *Modal) SetTitleStyle(style lipgloss.Style) *Modal {
	m.titleStyle = style
	return m
}

// SetContentStyle sets a custom content style
func (m *Modal) SetContentStyle(style lipgloss.Style) *Modal {
	m.contentStyle = style
	return m
}

// SetFooterStyle sets a custom footer style
func (m *Modal) SetFooterStyle(style lipgloss.Style) *Modal {
	m.footerStyle = style
	return m
}

// Render renders the modal centered on screen
func (m *Modal) Render() string {
	// Calculate actual width
	actualWidth := m.width
	if actualWidth > m.maxWidth {
		actualWidth = m.maxWidth
	}
	if actualWidth > m.windowWidth-4 {
		actualWidth = m.windowWidth - 4
	}
	
	// Build modal content
	var parts []string
	
	// Title with icon
	if m.title != "" {
		titleText := m.title
		if m.icon != "" {
			titleText = m.icon + " " + m.title
		}
		parts = append(parts, m.titleStyle.Render(titleText))
	}
	
	// Main content
	if m.content != "" {
		parts = append(parts, m.contentStyle.Render(m.content))
	}
	
	// Footer
	if m.footer != "" {
		parts = append(parts, m.footerStyle.Render(m.footer))
	}
	
	// Join all parts
	modalContent := strings.Join(parts, "\n")
	
	// Create modal container
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.borderColor).
		Padding(2, 3).
		Width(actualWidth - 2). // Account for border
		MaxWidth(actualWidth - 2)
	
	// Render with border
	modal := modalStyle.Render(modalContent)
	
	// Center in viewport
	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

// RenderWithDimmedBackground renders the modal with a dimmed background
// The background parameter should be the full screen content to dim
func (m *Modal) RenderWithDimmedBackground(background string) string {
	// Since we can't truly overlay in terminal, we'll just return the modal
	// In a real implementation, the calling code would need to handle not rendering
	// the background when a modal is shown
	return m.Render()
}

// Builder pattern convenience methods for common configurations

// Alert creates a modal configured for alert messages
func (m *Modal) Alert() *Modal {
	return m.SetBorderColor(theme.Warning()).
		SetWidth(50).
		SetTitleStyle(m.titleStyle.Foreground(theme.Warning()))
}

// Error creates a modal configured for error messages
func (m *Modal) Error() *Modal {
	return m.SetBorderColor(theme.Error()).
		SetWidth(50).
		SetTitleStyle(m.titleStyle.Foreground(theme.Error()))
}

// Success creates a modal configured for success messages
func (m *Modal) Success() *Modal {
	return m.SetBorderColor(theme.Success()).
		SetWidth(50).
		SetTitleStyle(m.titleStyle.Foreground(theme.Success()))
}

// Form creates a modal configured for forms
func (m *Modal) Form() *Modal {
	return m.SetWidth(70).
		SetBorderColor(theme.BorderFocus())
}

// Large creates a larger modal for detailed content
func (m *Modal) Large() *Modal {
	return m.SetWidth(80).
		SetMaxWidth(100)
}

// Small creates a smaller modal for simple messages
func (m *Modal) Small() *Modal {
	return m.SetWidth(40).
		SetMaxWidth(50)
}