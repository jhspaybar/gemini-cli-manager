package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the entire application UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// If showing modal, render it on top
	if m.showingModal && m.modal != nil {
		return m.modal.View()
	}

	// Calculate dimensions
	sidebarWidth := 20
	contentWidth := m.windowWidth - sidebarWidth - 1 // -1 for border
	contentHeight := m.windowHeight - 3 // -3 for status bar
	
	// Render components
	sidebar := m.renderSidebar(sidebarWidth, contentHeight)
	content := m.renderContent(contentWidth, contentHeight)
	statusBar := m.renderStatusBar()
	
	// Combine sidebar and content
	main := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		content,
	)
	
	// Add status bar
	return lipgloss.JoinVertical(
		lipgloss.Left,
		main,
		statusBar,
	)
}

// renderSidebar renders the navigation sidebar
func (m Model) renderSidebar(width, height int) string {
	var items []string
	
	// Title with focus indicator
	title := "Gemini CLI Manager"
	if m.focusedPane == PaneSidebar {
		title = "◀ " + title
	}
	items = append(items, h2Style.Render(title), "")
	
	// Profile info
	profileInfo := bodySmallStyle.Render(fmt.Sprintf("Profile: %s", m.getProfileName()))
	items = append(items, profileInfo, "")
	
	// Menu items
	for i, item := range m.sidebarItems {
		var style lipgloss.Style
		prefix := "  "
		
		// Show current view differently
		if item.View == m.currentView {
			if m.focusedPane == PaneSidebar && i == m.sidebarCursor {
				// Focused and selected
				style = focusedItemStyle
				prefix = "▶ "
			} else {
				// Current view but not focused
				style = menuItemStyle.Foreground(colorPrimary)
				prefix = "• "
			}
		} else if m.focusedPane == PaneSidebar && i == m.sidebarCursor {
			// Focused but not current
			style = menuItemStyle.Background(colorSelected)
			prefix = "  "
		} else {
			// Normal item
			style = menuItemStyle
		}
		
		line := fmt.Sprintf("%s%s %s", prefix, item.Icon, item.Title)
		items = append(items, style.Render(line))
	}
	
	// Launch button
	items = append(items, "", "")
	launchStyle := menuItemStyle
	if m.focusedPane == PaneSidebar && m.sidebarCursor >= len(m.sidebarItems) {
		launchStyle = focusedItemStyle
	}
	launchBtn := launchStyle.Render("  ▶ Launch Gemini")
	items = append(items, launchBtn)
	
	// Navigation help
	items = append(items, "")
	if m.focusedPane == PaneSidebar {
		items = append(items, helpDescStyle.Render("→/l: Focus content"))
		items = append(items, helpKeyStyle.Render("Enter: Select"))
	} else {
		items = append(items, helpDescStyle.Render("←/h: Back to sidebar"))
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	
	// Apply sidebar styling
	style := sidebarStyle
	if m.focusedPane == PaneSidebar {
		style = style.BorderForeground(colorFocused).
			BorderStyle(lipgloss.ThickBorder())
	}
	
	return style.
		Width(width).
		Height(height).
		Render(content)
}

// renderContent renders the main content area
func (m Model) renderContent(width, height int) string {
	var content string
	
	switch m.currentView {
	case ViewExtensions:
		content = m.renderExtensions(width, height)
	case ViewProfiles:
		content = m.renderProfiles(width, height)
	case ViewSettings:
		content = m.renderSettings(width, height)
	case ViewHelp:
		content = m.renderHelp(width, height)
	}
	
	style := contentStyle.
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderLeft(false).
		BorderRight(true).
		BorderBottom(true)
	
	if m.focusedPane == PaneContent {
		style = style.BorderForeground(colorFocused).
			BorderStyle(lipgloss.ThickBorder())
	} else {
		style = style.BorderForeground(colorBorder)
	}
	
	return style.
		Width(width).
		Height(height).
		Render(content)
}

// renderExtensions renders the extensions view
func (m Model) renderExtensions(width, height int) string {
	var lines []string
	
	// Header with focus indicator
	focusIndicator := ""
	if m.focusedPane == PaneContent {
		focusIndicator = " ▸"
	}
	header := h1Style.Render("Extensions" + focusIndicator)
	count := fmt.Sprintf("%d extensions found", len(m.extensions))
	lines = append(lines, header, bodySmallStyle.Render(count), "")
	
	if len(m.extensions) == 0 {
		lines = append(lines, mutedStyle.Render("No extensions installed."))
		lines = append(lines, "", bodyStyle.Render("Press 'n' to add a new extension."))
	} else {
		// Extension list
		for i, ext := range m.extensions {
			var line string
			cursor := "  "
			style := bodyStyle
			
			if i == m.extensionsCursor {
				cursor = "▶ "
				style = style.Bold(true)
			}
			
			// Status indicator
			status := "✗"
			statusStyle := disabledStyle
			if ext.Enabled {
				status = "✓"
				statusStyle = enabledStyle
			}
			
			line = fmt.Sprintf("%s%-30s %s", cursor, ext.Name, statusStyle.Render(status))
			lines = append(lines, style.Render(line))
			
			// Show description for selected item
			if i == m.extensionsCursor {
				desc := bodySmallStyle.Render("  " + ext.Description)
				lines = append(lines, desc)
			}
		}
		
		// Help text
		lines = append(lines, "")
		lines = append(lines, helpDescStyle.Render("Space: Toggle • Enter: Details • n: New • d: Delete"))
	}
	
	return strings.Join(lines, "\n")
}

// renderProfiles renders the profiles view
func (m Model) renderProfiles(width, height int) string {
	var lines []string
	
	// Header
	header := h1Style.Render("Profiles")
	current := fmt.Sprintf("Active: %s", m.getProfileName())
	lines = append(lines, header, bodySmallStyle.Render(current), "")
	
	if len(m.profiles) == 0 {
		lines = append(lines, mutedStyle.Render("No profiles configured."))
		lines = append(lines, "", bodyStyle.Render("Press 'n' to create a new profile."))
	} else {
		// Profile list
		for i, prof := range m.profiles {
			var line string
			cursor := "  "
			style := bodyStyle
			
			if i == m.profilesCursor {
				cursor = "▶ "
				style = style.Bold(true)
			}
			
			// Active indicator
			indicator := "  "
			if m.currentProfile != nil && prof.ID == m.currentProfile.ID {
				indicator = "● "
				style = style.Foreground(colorSuccess)
			}
			
			line = fmt.Sprintf("%s%s%s", cursor, indicator, prof.Name)
			lines = append(lines, style.Render(line))
			
			// Show description for selected item
			if i == m.profilesCursor && prof.Description != "" {
				desc := bodySmallStyle.Render("  " + prof.Description)
				lines = append(lines, desc)
			}
		}
		
		// Help text
		lines = append(lines, "")
		lines = append(lines, helpDescStyle.Render("Enter: Activate • n: New • e: Edit • d: Delete"))
	}
	
	return strings.Join(lines, "\n")
}

// renderSettings renders the settings view
func (m Model) renderSettings(width, height int) string {
	var lines []string
	
	header := h1Style.Render("Settings")
	lines = append(lines, header, "")
	
	// Settings items
	settings := []struct {
		label string
		value string
	}{
		{"Gemini CLI Path", "/usr/local/bin/gemini"},
		{"Extensions Directory", "~/.gemini/extensions"},
		{"Profiles Directory", "~/.gemini/profiles"},
		{"Auto-update Extensions", "Enabled"},
		{"Theme", "GitHub Dark"},
	}
	
	for _, setting := range settings {
		label := bodyStyle.Render(fmt.Sprintf("%-25s", setting.label))
		value := bodySmallStyle.Render(setting.value)
		lines = append(lines, fmt.Sprintf("%s %s", label, value))
	}
	
	lines = append(lines, "")
	lines = append(lines, helpDescStyle.Render("Press 'e' to edit settings"))
	
	return strings.Join(lines, "\n")
}

// renderHelp renders the help view
func (m Model) renderHelp(width, height int) string {
	var lines []string
	
	header := h1Style.Render("Help")
	lines = append(lines, header, "")
	
	// Navigation
	lines = append(lines, h2Style.Render("Navigation"), "")
	navHelp := []struct {
		key  string
		desc string
	}{
		{"←/h", "Focus sidebar"},
		{"→/l", "Focus content"},
		{"↑/k", "Move up in current pane"},
		{"↓/j", "Move down in current pane"},
		{"Enter", "Select/Navigate"},
		{"Esc", "Back to sidebar"},
		{"Tab", "Switch panes"},
		{"?", "Toggle this help"},
		{"q", "Quit application"},
	}
	
	for _, item := range navHelp {
		key := helpKeyStyle.Render(fmt.Sprintf("%-10s", item.key))
		desc := helpDescStyle.Render(item.desc)
		lines = append(lines, fmt.Sprintf("  %s %s", key, desc))
	}
	
	// View-specific
	lines = append(lines, "", h2Style.Render("View-Specific Actions"), "")
	viewHelp := []struct {
		key  string
		desc string
	}{
		{"Space", "Toggle extension on/off (Extensions view)"},
		{"n", "Create new item"},
		{"e", "Edit selected item"},
		{"d", "Delete selected item"},
		{"L", "Launch Gemini CLI (any view)"},
	}
	
	for _, item := range viewHelp {
		key := helpKeyStyle.Render(fmt.Sprintf("%-10s", item.key))
		desc := helpDescStyle.Render(item.desc)
		lines = append(lines, fmt.Sprintf("  %s %s", key, desc))
	}
	
	lines = append(lines, "", mutedStyle.Render("Press any key to return"))
	
	return strings.Join(lines, "\n")
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	var left []string
	var right []string
	
	// Left side - context info
	enabledCount := 0
	for _, ext := range m.extensions {
		if ext.Enabled {
			enabledCount++
		}
	}
	left = append(left, fmt.Sprintf("Extensions: %d/%d", enabledCount, len(m.extensions)))
	
	if m.currentProfile != nil {
		left = append(left, fmt.Sprintf("Profile: %s", m.currentProfile.Name))
	}
	
	// Right side - shortcuts
	if m.focusedPane == PaneSidebar {
		right = append(right, "→: Content")
		right = append(right, "Enter: Select")
	} else {
		right = append(right, "←: Sidebar")
		right = append(right, "Space: Toggle")
	}
	right = append(right, "?: Help")
	right = append(right, "q: Quit")
	
	// Error in the middle if present
	middle := ""
	if m.err != nil {
		middle = errorStyle.Render(" Error: " + m.err.Error() + " ")
	}
	
	leftStr := strings.Join(left, " • ")
	rightStr := strings.Join(right, " • ")
	
	// Calculate spacing
	totalWidth := m.windowWidth
	leftWidth := lipgloss.Width(leftStr)
	rightWidth := lipgloss.Width(rightStr)
	middleWidth := lipgloss.Width(middle)
	
	spacing := totalWidth - leftWidth - rightWidth - middleWidth - 2
	if spacing < 0 {
		spacing = 1
	}
	
	status := leftStr + strings.Repeat(" ", spacing/2) + middle + strings.Repeat(" ", spacing/2) + rightStr
	
	return statusBarStyle.
		Width(m.windowWidth).
		Render(status)
}

// Helper methods
func (m Model) getProfileName() string {
	if m.currentProfile != nil {
		return m.currentProfile.Name
	}
	return "None"
}

