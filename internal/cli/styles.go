package cli

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// Color variables that map to theme colors
var (
	// Base colors
	colorBg        lipgloss.TerminalColor = theme.Background()
	colorBgLight   lipgloss.TerminalColor = theme.Selection()
	colorText      lipgloss.TerminalColor = theme.TextPrimary()
	colorTextDim   lipgloss.TerminalColor = theme.TextSecondary()
	colorTextMuted lipgloss.TerminalColor = theme.TextMuted()

	// Accent colors
	colorAccent    lipgloss.TerminalColor = theme.Primary()
	colorAccentDim lipgloss.TerminalColor = theme.Secondary()
	colorSuccess   lipgloss.TerminalColor = theme.Success()
	colorWarning   lipgloss.TerminalColor = theme.Warning()
	colorError     lipgloss.TerminalColor = theme.Error()

	// UI colors
	colorBorder      lipgloss.TerminalColor = theme.Border()
	colorBorderFocus lipgloss.TerminalColor = theme.BorderFocus()
	colorBorderDim   lipgloss.TerminalColor = theme.Border()

	// Special colors
	colorHighlight lipgloss.TerminalColor = theme.Selection()
)

// UpdateStyles refreshes all styles with current theme colors
func UpdateStyles() {
	// Update color variables
	colorBg = theme.Background()
	colorBgLight = theme.Selection()
	colorText = theme.TextPrimary()
	colorTextDim = theme.TextSecondary()
	colorTextMuted = theme.TextMuted()
	colorAccent = theme.Primary()
	colorAccentDim = theme.Secondary()
	colorSuccess = theme.Success()
	colorWarning = theme.Warning()
	colorError = theme.Error()
	colorBorder = theme.Border()
	colorBorderFocus = theme.BorderFocus()
	colorBorderDim = theme.Border()
	colorHighlight = theme.Selection()

	// Update all styles
	updateTextStyles()
	updateLayoutStyles()
	updateInteractiveStyles()
	updateTabStyles()
}

// Text styles
var (
	// Headers
	titleStyle lipgloss.Style
	h1Style    lipgloss.Style
	h2Style    lipgloss.Style
	// Body text
	textStyle      lipgloss.Style
	textDimStyle   lipgloss.Style
	textMutedStyle lipgloss.Style
	// Emphasis
	accentStyle  lipgloss.Style
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
)

func updateTextStyles() {
	// Headers
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorText)

	h1Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorText).
		MarginBottom(1)

	h2Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorText)

	// Body text
	textStyle = lipgloss.NewStyle().
		Foreground(colorText)

	textDimStyle = lipgloss.NewStyle().
		Foreground(colorTextDim)

	textMutedStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted)

	// Emphasis
	accentStyle = lipgloss.NewStyle().
		Foreground(colorAccent)

	successStyle = lipgloss.NewStyle().
		Foreground(colorSuccess)

	errorStyle = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)
}

// Layout components
var (
	// Sidebar
	sidebarStyle        lipgloss.Style
	sidebarFocusedStyle lipgloss.Style
	// Content
	contentStyle        lipgloss.Style
	contentFocusedStyle lipgloss.Style
	// Status bar
	statusBarStyle lipgloss.Style
)

func updateLayoutStyles() {
	// Sidebar
	sidebarStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(colorBorder).
		Padding(1)

	sidebarFocusedStyle = sidebarStyle.Copy().
		BorderForeground(colorBorderFocus).
		BorderStyle(lipgloss.ThickBorder())

	// Content
	contentStyle = lipgloss.NewStyle().
		Padding(1)

	contentFocusedStyle = contentStyle.Copy().
		Border(lipgloss.ThickBorder()).
		BorderForeground(colorBorderFocus)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(colorBorder).
		Foreground(colorTextDim).
		Padding(0, 1)
}

// Interactive elements
var (
	// Menu items
	menuItemStyle         lipgloss.Style
	menuItemSelectedStyle lipgloss.Style
	menuItemActiveStyle   lipgloss.Style
	// Modal
	modalStyle      lipgloss.Style
	modalTitleStyle lipgloss.Style
	// Input
	inputStyle        lipgloss.Style
	inputFocusedStyle lipgloss.Style
	// Labels
	labelStyle lipgloss.Style
	helpStyle  lipgloss.Style
	// Special
	spinnerStyle lipgloss.Style
	keyStyle     lipgloss.Style
	keyDescStyle lipgloss.Style
	formStyle    lipgloss.Style
)

func updateInteractiveStyles() {
	// Menu items
	menuItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	menuItemSelectedStyle = lipgloss.NewStyle().
		Foreground(colorText).
		Background(colorHighlight).
		PaddingLeft(1).
		PaddingRight(1)

	menuItemActiveStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		PaddingLeft(2)

	// Modal
	modalStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(2).
		Background(colorBg)

	modalTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorText).
		MarginBottom(1)

	// Input
	inputStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder)

	inputFocusedStyle = inputStyle.Copy().
		BorderForeground(colorBorderFocus)

	// Labels
	labelStyle = lipgloss.NewStyle().
		Foreground(colorTextDim)

	// Help text
	helpStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted).
		Italic(true)

	// Spinner
	spinnerStyle = lipgloss.NewStyle().
		Foreground(colorAccent)

	// Key bindings
	keyStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	keyDescStyle = lipgloss.NewStyle().
		Foreground(colorTextDim)

	// Form
	formStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(2)
}

// Helper functions
func renderKeyHelp(bindings [][2]string) string {
	var parts []string
	for _, binding := range bindings {
		part := keyStyle.Render(binding[0]) + " " + keyDescStyle.Render(binding[1])
		parts = append(parts, part)
	}
	return strings.Join(parts, " • ")
}

// Tab-specific styles
var (
	inactiveTabBorder lipgloss.Border
	activeTabBorder   lipgloss.Border
	inactiveTabStyle  lipgloss.Style
	activeTabStyle    lipgloss.Style
	tabGapStyle       lipgloss.Style
)

// tabBorderWithBottom creates a custom border for tabs
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func updateTabStyles() {
	// Define custom borders for tabs
	// Inactive tabs have special bottom corners for connection
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	
	// Active tab has no bottom border to connect with content
	activeTabBorder = tabBorderWithBottom("┘", " ", "└")

	// Inactive tab style
	inactiveTabStyle = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(colorBorder).
		Padding(0, 2).
		Foreground(colorText)

	// Active tab style
	activeTabStyle = lipgloss.NewStyle().
		Border(activeTabBorder, true).
		BorderForeground(colorAccent).
		Padding(0, 2).
		Foreground(colorAccent).
		Bold(true)

	// Style for the gap filler
	tabGapStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.Border{Bottom: "─"}).
		BorderForeground(colorBorder).
		BorderBottom(true)
}

// Initialize styles on import
func init() {
	UpdateStyles()
}
