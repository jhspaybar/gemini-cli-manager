package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// View renders the entire application UI
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	
	// Show loading screen while data is being loaded
	if m.loading {
		return m.renderLoading()
	}

	// If showing modal, render it on top
	if m.showingModal && m.modal != nil {
		return m.modal.View()
	}

	// Calculate dimensions
	contentHeight := m.windowHeight - 4 // -4 for tab bar and status bar
	
	// Render components
	tabBar := m.renderTabBar()
	content := m.renderContent(m.windowWidth, contentHeight)
	statusBar := m.renderStatusBar()
	
	// Combine all elements
	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabBar,
		content,
		statusBar,
	)
}

// renderTabBar renders the top navigation tabs
func (m Model) renderTabBar() string {
	// Don't show tabs in detail view
	if m.currentView == ViewExtensionDetail {
		return lipgloss.NewStyle().
			Width(m.windowWidth).
			Height(1).
			Render("")
	}
	
	tabs := []struct {
		title string
		icon  string
		view  ViewType
	}{
		{"Extensions", "🧩", ViewExtensions},
		{"Profiles", "👥", ViewProfiles},
		{"Settings", "⚙️", ViewSettings},
		{"Help", "❓", ViewHelp},
	}
	
	var tabStrings []string
	
	for i, tab := range tabs {
		var tabStyle lipgloss.Style
		
		if tab.view == m.currentView {
			// Active tab - clean look with bottom highlight
			tabStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent).
				Background(lipgloss.Color("236")).
				Padding(0, 3).
				MarginRight(1)
		} else {
			// Inactive tab
			tabStyle = lipgloss.NewStyle().
				Foreground(colorTextDim).
				Background(lipgloss.Color("235")).
				Padding(0, 3).
				MarginRight(1)
		}
		
		// Add left margin for first tab
		if i == 0 {
			tabStyle = tabStyle.MarginLeft(2)
		}
		
		tabContent := fmt.Sprintf("%s %s", tab.icon, tab.title)
		tabStrings = append(tabStrings, tabStyle.Render(tabContent))
	}
	
	// Create tab row
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, tabStrings...)
	
	// Create full-width container with bottom border
	return lipgloss.NewStyle().
		Width(m.windowWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(colorBorder).
		Render(tabRow)
}

// renderContent renders the main content area
func (m Model) renderContent(width, height int) string {
	var content string
	
	// Calculate inner width accounting for padding
	innerWidth := width - 4  // 4 for left and right padding (2 each side)
	
	switch m.currentView {
	case ViewExtensions:
		content = m.renderExtensions(innerWidth, height)
	case ViewProfiles:
		content = m.renderProfiles(innerWidth, height)
	case ViewSettings:
		content = m.renderSettings(innerWidth, height)
	case ViewHelp:
		content = m.renderHelp(innerWidth, height)
	case ViewExtensionDetail:
		content = m.renderExtensionDetail(innerWidth, height)
	}
	
	// Simple content styling without borders
	return lipgloss.NewStyle().
		Padding(1, 2).
		Width(width).
		Height(height).
		MaxWidth(width).
		Render(content)
}

// renderExtensions renders the extensions view
func (m Model) renderExtensions(width, height int) string {
	var lines []string
	
	// Header
	header := h1Style.Render("Extensions")
	lines = append(lines, header)
	
	// Show search bar if active or has query
	if m.searchActive || m.searchBar.Value() != "" {
		lines = append(lines, "")
		searchBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorderFocus).
			Padding(0, 1).
			Width(60).
			Render(m.searchBar.View())
		lines = append(lines, searchBox)
	}
	
	lines = append(lines, "")
	
	// Show count
	var count string
	if m.searchBar.Value() != "" {
		count = fmt.Sprintf("%d of %d extensions (filtered)", len(m.filteredExtensions), len(m.extensions))
	} else {
		count = fmt.Sprintf("%d extensions found", len(m.filteredExtensions))
	}
	lines = append(lines, textDimStyle.Render(count), "")
	
	if len(m.filteredExtensions) == 0 {
		// Empty state
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(2, 4).
			Width(50).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"📦",
					"",
					"No extensions installed",
					textDimStyle.Render("Press 'n' to install your first extension"),
				),
			)
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Width(width-4).Align(lipgloss.Center).Render(emptyBox))
	} else {
		// Extension list with cards
		for i, ext := range m.filteredExtensions {
			isSelected := i == m.extensionsCursor
			card := m.renderExtensionCard(ext, isSelected, width-4)
			lines = append(lines, card)
			if i < len(m.filteredExtensions)-1 {
				lines = append(lines, "") // Spacing between cards
			}
		}
		
		// Help text
		lines = append(lines, "", "")
		helpText := renderKeyHelp([][2]string{
			{"↵", "Details"},
			{"n", "Install"},
			{"d", "Delete"},
			{"/", "Search"},
			{"Tab", "Next"},
		})
		lines = append(lines, helpText)
	}
	
	return strings.Join(lines, "\n")
}

// renderExtensionCard renders a single extension as a card
func (m Model) renderExtensionCard(ext *extension.Extension, isSelected bool, width int) string {
	// Card styling - more compact
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1).
		Width(width)
	
	if isSelected {
		cardStyle = cardStyle.
			BorderForeground(colorAccent).
			BorderStyle(lipgloss.ThickBorder())
	}
	
	// Extension name and version on same line
	nameStyle := textStyle
	if isSelected {
		nameStyle = nameStyle.Bold(true).Foreground(colorAccent)
	}
	
	// Calculate available width for text (accounting for padding and borders)
	textWidth := width - 6  // 2 for borders, 4 for padding
	
	// Calculate space for name and version
	versionText := fmt.Sprintf("v%s", ext.Version)
	nameWidth := textWidth - lipgloss.Width(versionText) - 2 // 2 for spacing
	
	name := nameStyle.MaxWidth(nameWidth).Render(ext.Name)
	version := textDimStyle.Render(versionText)
	header := lipgloss.JoinHorizontal(lipgloss.Top, name, "  ", version)
	
	// Description on second line
	desc := textDimStyle.MaxWidth(textWidth).Render(ext.Description)
	
	// MCP servers info if present
	var content []string
	content = append(content, header)
	content = append(content, desc)
	
	if ext.MCPServers != nil && len(ext.MCPServers) > 0 {
		count := len(ext.MCPServers)
		info := accentStyle.Render(fmt.Sprintf("⚡ %d MCP server%s", count, pluralize(count)))
		content = append(content, info)
	}
	
	return cardStyle.Render(strings.Join(content, "\n"))
}

// renderProfiles renders the profiles view
func (m Model) renderProfiles(width, height int) string {
	var lines []string
	
	// Header
	header := h1Style.Render("Profiles")
	lines = append(lines, header)
	
	// Active profile indicator
	activeProfile := "None"
	if m.currentProfile != nil {
		activeProfile = m.currentProfile.Name
	}
	// Constrain badge width to prevent overflow
	badgeText := fmt.Sprintf("● Active: %s", activeProfile)
	activeBadge := lipgloss.NewStyle().
		Background(colorSuccess).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		MaxWidth(width - 4).
		Render(badgeText)
	lines = append(lines, "", activeBadge)
	
	// Show search bar if active
	if m.searchActive || m.searchBar.Value() != "" {
		lines = append(lines, "")
		searchBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorderFocus).
			Padding(0, 1).
			Width(60).
			Render(m.searchBar.View())
		lines = append(lines, searchBox)
	}
	
	lines = append(lines, "")
	
	// Show count
	var count string
	if m.searchBar.Value() != "" {
		count = fmt.Sprintf("%d of %d profiles (filtered)", len(m.filteredProfiles), len(m.profiles))
	} else {
		count = fmt.Sprintf("%d profiles", len(m.profiles))
	}
	lines = append(lines, textDimStyle.Render(count), "")
	
	if len(m.filteredProfiles) == 0 {
		// Empty state
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(2, 4).
			Width(50).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"👤",
					"",
					"No profiles configured",
					textDimStyle.Render("Press 'n' to create your first profile"),
				),
			)
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Width(width-4).Align(lipgloss.Center).Render(emptyBox))
	} else {
		// Profile list
		for i, prof := range m.filteredProfiles {
			isSelected := i == m.profilesCursor
			isActive := m.currentProfile != nil && prof.ID == m.currentProfile.ID
			card := m.renderProfileCard(prof, isSelected, isActive, width-4)
			lines = append(lines, card)
			if i < len(m.filteredProfiles)-1 {
				lines = append(lines, "")
			}
		}
		
		// Help text
		lines = append(lines, "", "")
		helpText := renderKeyHelp([][2]string{
			{"↵", "Activate"},
			{"n", "New"},
			{"e", "Edit"},
			{"d", "Delete"},
			{"/", "Search"},
		})
		lines = append(lines, helpText)
	}
	
	return strings.Join(lines, "\n")
}

// renderProfileCard renders a single profile as a card
func (m Model) renderProfileCard(prof *profile.Profile, isSelected, isActive bool, width int) string {
	// Card styling - more compact
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1).
		Width(width)
	
	if isActive {
		cardStyle = cardStyle.BorderForeground(colorSuccess)
	} else if isSelected {
		cardStyle = cardStyle.
			BorderForeground(colorAccent).
			BorderStyle(lipgloss.ThickBorder())
	}
	
	// Profile name with active indicator
	nameStyle := textStyle
	if isSelected {
		nameStyle = nameStyle.Bold(true)
	}
	if isActive {
		nameStyle = nameStyle.Foreground(colorSuccess)
	}
	
	statusIcon := "  "
	if isActive {
		statusIcon = "● "
	}
	
	// Build content lines
	var content []string
	
	// Calculate available width for text (accounting for padding and borders)
	textWidth := width - 6  // 2 for borders, 4 for padding
	
	// First line: status + name
	// Apply style only to the name, not the status icon
	styledName := nameStyle.MaxWidth(textWidth - lipgloss.Width(statusIcon)).Render(prof.Name)
	nameText := statusIcon + styledName
	content = append(content, nameText)
	
	// Second line: description (if exists)
	if prof.Description != "" {
		desc := textDimStyle.MaxWidth(textWidth).Render(prof.Description)
		content = append(content, desc)
	}
	
	// Third line: extension count
	if len(prof.Extensions) > 0 {
		extInfo := accentStyle.MaxWidth(textWidth).Render(fmt.Sprintf("📦 %d extension%s", len(prof.Extensions), pluralize(len(prof.Extensions))))
		content = append(content, extInfo)
	}
	
	return cardStyle.Render(strings.Join(content, "\n"))
}

// renderSettings renders the settings view
func (m Model) renderSettings(width, height int) string {
	var lines []string
	
	header := h1Style.Render("Settings")
	lines = append(lines, header, "")
	
	// Settings in a more compact table-like format
	settings := []struct {
		section string
		icon    string
		items   []struct {
			label string
			value string
		}
	}{
		{
			"General",
			"🔧",
			[]struct {
				label string
				value string
			}{
				{"Gemini CLI Path", "/usr/local/bin/gemini"},
				{"Config Directory", "~/.gemini-cli-manager"},
			},
		},
		{
			"Extensions",
			"📦",
			[]struct {
				label string
				value string
			}{
				{"Extensions Directory", "~/.gemini/extensions"},
				{"Auto-update", "Enabled"},
			},
		},
		{
			"Appearance",
			"🎨",
			[]struct {
				label string
				value string
			}{
				{"Theme", "Dark Modern"},
				{"Tab Position", "Top"},
			},
		},
	}
	
	// Calculate column widths
	labelWidth := 25
	
	for _, group := range settings {
		// Section header with icon
		sectionHeader := fmt.Sprintf("%s %s", group.icon, group.section)
		lines = append(lines, h2Style.Render(sectionHeader))
		
		// Items in the section
		for _, item := range group.items {
			label := textDimStyle.Width(labelWidth).Render(item.label)
			value := accentStyle.Render(item.value)
			line := fmt.Sprintf("  %s  %s", label, value)
			lines = append(lines, line)
		}
		lines = append(lines, "") // Space between sections
	}
	
	// Help at bottom
	lines = append(lines, keyDescStyle.Render("Press 'e' to edit settings"))
	
	return strings.Join(lines, "\n")
}

// renderHelp renders the help view
func (m Model) renderHelp(width, height int) string {
	var lines []string
	
	header := h1Style.Render("Help")
	lines = append(lines, header, "")
	
	// Two-column layout for keyboard shortcuts
	shortcuts := []struct {
		category string
		icon     string
		items    [][2]string
	}{
		{
			"Navigation",
			"🧭",
			[][2]string{
				{"Tab", "Next tab"},
				{"←/h", "Previous tab"},
				{"→/l", "Next tab"},
				{"↑/k", "Move up"},
				{"↓/j", "Move down"},
				{"Enter", "Select"},
			},
		},
		{
			"Actions",
			"⚡",
			[][2]string{
				{"n", "New item"},
				{"e", "Edit item"},
				{"d", "Delete item"},
				{"/", "Search"},
				{"L", "Launch Gemini"},
				{"?", "Toggle help"},
			},
		},
		{
			"Global",
			"🌐",
			[][2]string{
				{"Ctrl+P", "Quick switch profile"},
				{"Esc", "Cancel/Back"},
				{"q", "Quit"},
			},
		},
	}
	
	// Render shortcuts in columns
	for _, section := range shortcuts {
		lines = append(lines, h2Style.Render(fmt.Sprintf("%s %s", section.icon, section.category)))
		
		for _, item := range section.items {
			key := keyStyle.Width(10).Render(item[0])
			desc := textDimStyle.Render(item[1])
			lines = append(lines, fmt.Sprintf("  %s  %s", key, desc))
		}
		lines = append(lines, "")
	}
	
	// Tips section
	lines = append(lines, h2Style.Render("💡 Tips"))
	tips := []string{
		"• Use Tab to quickly cycle through all views",
		"• Press / to search in any list view",
		"• Create profiles for different workflows",
		"• Extensions can be installed from URLs or local paths",
	}
	for _, tip := range tips {
		lines = append(lines, textDimStyle.Render(tip))
	}
	
	return strings.Join(lines, "\n")
}

// renderExtensionDetail renders the detailed view of an extension
func (m Model) renderExtensionDetail(width, height int) string {
	if m.selectedExtension == nil {
		return "No extension selected"
	}
	
	var lines []string
	ext := m.selectedExtension
	
	// Header with back navigation hint
	header := h1Style.Render(fmt.Sprintf("📦 %s", ext.Name))
	lines = append(lines, header)
	lines = append(lines, textDimStyle.Render("Press Esc to go back"))
	lines = append(lines, "")
	
	// Basic info section
	lines = append(lines, h2Style.Render("📋 Basic Information"))
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1).
		Width(width-2).
		Render(strings.Join([]string{
			fmt.Sprintf("ID:          %s", ext.ID),
			fmt.Sprintf("Name:        %s", ext.Name),
			fmt.Sprintf("Version:     %s", ext.Version),
			fmt.Sprintf("Description: %s", ext.Description),
			fmt.Sprintf("Path:        %s", ext.Path),
		}, "\n"))
	lines = append(lines, infoBox)
	lines = append(lines, "")
	
	// MCP Servers section if present
	if ext.MCPServers != nil && len(ext.MCPServers) > 0 {
		lines = append(lines, h2Style.Render("⚡ MCP Servers"))
		for name, config := range ext.MCPServers {
			serverBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(1).
				Width(width-2).
				Render(strings.Join([]string{
					fmt.Sprintf("Server: %s", name),
					fmt.Sprintf("Command: %s", config.Command),
					fmt.Sprintf("Args: %v", config.Args),
					fmt.Sprintf("Env: %v", config.Env),
				}, "\n"))
			lines = append(lines, serverBox)
		}
		lines = append(lines, "")
	}
	
	// TODO: Add context file support when Extension struct includes it
	// For now, just show that context files can exist
	lines = append(lines, h2Style.Render("📄 Context File"))
	lines = append(lines, textDimStyle.Render("Context file name: " + func() string {
		if ext.ContextFileName != "" {
			return ext.ContextFileName
		}
		return "GEMINI.md (default)"
	}()))
	lines = append(lines, "")
	
	// Help text
	lines = append(lines, keyDescStyle.Render("Esc: Back • d: Delete • e: Edit (coming soon)"))
	
	return strings.Join(lines, "\n")
}

// renderLoading renders the loading screen
func (m Model) renderLoading() string {
	width := m.windowWidth
	height := m.windowHeight
	
	// Center the loading message
	loadingBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Padding(2, 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			"🚀 Gemini CLI Manager",
			"",
			"Loading extensions and profiles...",
			"",
			"Please wait",
		))
	
	// Center in viewport
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		loadingBox,
	)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	var parts []string
	
	// Profile indicator
	if m.currentProfile != nil {
		parts = append(parts, fmt.Sprintf("👤 %s", m.currentProfile.Name))
	} else {
		parts = append(parts, "👤 No Profile")
	}
	
	// Extension count
	enabledCount := 0
	if m.currentProfile != nil {
		for _, extRef := range m.currentProfile.Extensions {
			if extRef.Enabled {
				enabledCount++
			}
		}
	}
	parts = append(parts, fmt.Sprintf("🧩 %d/%d", enabledCount, len(m.extensions)))
	
	// Key hints based on context
	var hints []string
	hints = append(hints, "Tab: Switch")
	hints = append(hints, "L: Launch")
	hints = append(hints, "?: Help")
	hints = append(hints, "q: Quit")
	
	// Error/Info display
	errorMsg := ""
	if m.err != nil {
		if uiErr, ok := m.err.(UIError); ok {
			if uiErr.Type == ErrorTypeInfo {
				// Info message - use different styling
				errorMsg = accentStyle.Render(fmt.Sprintf(" ℹ️  %s ", uiErr.Message))
				if uiErr.Details != "" {
					errorMsg += textDimStyle.Render(fmt.Sprintf(" - %s", uiErr.Details))
				}
			} else {
				errorMsg = errorStyle.Render(fmt.Sprintf(" ❌ %s ", uiErr.Message))
			}
		} else {
			errorMsg = errorStyle.Render(fmt.Sprintf(" ❌ %s ", m.err.Error()))
		}
	}
	
	// Build status bar
	left := strings.Join(parts, " • ")
	right := strings.Join(hints, " • ")
	
	// Calculate spacing
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	errorWidth := lipgloss.Width(errorMsg)
	availableWidth := m.windowWidth - 2 // Account for padding
	
	// If content is too wide, truncate the left side (profile name)
	totalContentWidth := leftWidth + rightWidth + errorWidth + 4 // min spacing
	if totalContentWidth > availableWidth {
		// Truncate left side to fit
		maxLeftWidth := availableWidth - rightWidth - errorWidth - 4
		if maxLeftWidth > 0 {
			left = lipgloss.NewStyle().MaxWidth(maxLeftWidth).Render(left)
			leftWidth = lipgloss.Width(left)
		}
	}
	
	spacingTotal := availableWidth - leftWidth - rightWidth - errorWidth
	if spacingTotal < 2 {
		spacingTotal = 2
	}
	
	spacing1 := spacingTotal / 2
	spacing2 := spacingTotal - spacing1
	
	content := left + strings.Repeat(" ", spacing1) + errorMsg + strings.Repeat(" ", spacing2) + right
	
	return statusBarStyle.
		Width(m.windowWidth).
		Render(content)
}

// Helper methods
func (m Model) getProfileName() string {
	if m.currentProfile != nil {
		return m.currentProfile.Name
	}
	return "None"
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}