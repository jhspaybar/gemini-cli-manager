package cli

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

// Simple, clean color scheme
var (
	// Base colors
	colorBg        = lipgloss.Color("235")
	colorText      = lipgloss.Color("252")
	colorTextDim   = lipgloss.Color("246")
	colorTextMuted = lipgloss.Color("241")
	
	// Accent colors
	colorAccent  = lipgloss.Color("39")
	colorSuccess = lipgloss.Color("42")
	colorError   = lipgloss.Color("196")
	
	// UI colors
	colorBorder      = lipgloss.Color("240")
	colorBorderFocus = lipgloss.Color("33")
)

// Text styles
var (
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
)

// Layout components
var (
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
)

// Interactive elements
var (
	// Menu items
	menuItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	
	menuItemSelectedStyle = lipgloss.NewStyle().
		Foreground(colorText).
		Background(lipgloss.Color("237")).
		PaddingLeft(1).
		PaddingRight(1)
	
	menuItemActiveStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		PaddingLeft(2)
)

// Form elements
var (
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
)

// Special components
var (
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
)

// Helper functions
func renderKeyHelp(bindings [][2]string) string {
	var parts []string
	for _, binding := range bindings {
		part := keyStyle.Render(binding[0]) + " " + keyDescStyle.Render(binding[1])
		parts = append(parts, part)
	}
	return strings.Join(parts, " • ")
}