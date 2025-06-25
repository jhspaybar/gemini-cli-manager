package cli

import (
	"github.com/charmbracelet/lipgloss"
)

// Neutral, readable color scheme
var (
	// Base colors - softer and more neutral
	colorBackground   = lipgloss.Color("235") // Dark gray
	colorSurface      = lipgloss.Color("237") // Slightly lighter gray
	colorBorder       = lipgloss.Color("240") // Subtle border
	
	// Text colors - better contrast
	colorTextPrimary   = lipgloss.Color("252") // Almost white
	colorTextSecondary = lipgloss.Color("246") // Light gray
	colorTextMuted     = lipgloss.Color("241") // Muted gray
	
	// Accent colors - more subtle
	colorPrimary  = lipgloss.Color("39")  // Soft blue
	colorSuccess  = lipgloss.Color("42")  // Soft green
	colorWarning  = lipgloss.Color("214") // Soft orange
	colorError    = lipgloss.Color("203") // Soft red
	
	// State colors
	colorSelected = lipgloss.Color("240") // Subtle highlight
	colorFocused  = lipgloss.Color("33")  // Blue for focused items
)

// Styles - cleaner and more minimal
var (
	// Base styles
	baseStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary)
	
	// Headers - less aggressive
	h1Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorTextPrimary).
		MarginBottom(1)
	
	h2Style = lipgloss.NewStyle().
		Foreground(colorTextPrimary)
	
	// Body text
	bodyStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary)
	
	bodySmallStyle = lipgloss.NewStyle().
		Foreground(colorTextSecondary)
	
	mutedStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted)
	
	// UI elements - cleaner borders
	sidebarStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(colorBorder).
		Padding(1)
	
	contentStyle = lipgloss.NewStyle().
		Padding(1)
	
	statusBarStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(colorBorder).
		Foreground(colorTextSecondary).
		Padding(0, 1)
	
	// Interactive elements - clearer selection
	menuItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)
	
	selectedMenuItemStyle = lipgloss.NewStyle().
		Foreground(colorTextPrimary).
		Background(colorSelected).
		PaddingLeft(1).
		PaddingRight(1)
	
	focusedItemStyle = lipgloss.NewStyle().
		Foreground(colorFocused).
		Bold(true)
	
	// Status indicators
	enabledStyle = lipgloss.NewStyle().
		Foreground(colorSuccess)
	
	disabledStyle = lipgloss.NewStyle().
		Foreground(colorTextMuted)
	
	// Help text
	helpKeyStyle = lipgloss.NewStyle().
		Foreground(colorPrimary)
	
	helpDescStyle = lipgloss.NewStyle().
		Foreground(colorTextSecondary)
	
	// Error style
	errorStyle = lipgloss.NewStyle().
		Foreground(colorError)
)