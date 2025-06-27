package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// View renders the entire application UI
func (m Model) View() string {
	LogDebug("Model.View called, ready=%v, loading=%v, showingModal=%v", m.ready, m.loading, m.showingModal)
	
	if !m.ready {
		return "\n  Initializing..."
	}
	
	// Show loading screen while data is being loaded
	if m.loading {
		return m.renderLoading()
	}

	// If showing modal, render it on top
	if m.showingModal && m.modal != nil {
		LogDebug("View: Rendering modal, type: %T", m.modal)
		modalView := m.modal.View()
		LogDebug("View: Modal rendered, length: %d", len(modalView))
		return modalView
	}

	// Create main flexbox layout
	fb := flexbox.New(m.windowWidth, m.windowHeight)
	
	// Tab bar row
	tabRow := fb.NewRow()
	tabCell := flexbox.NewCell(1, 1)
	tabCell.SetContent(m.renderTabBar())
	tabRow.AddCells(tabCell)
	
	// Content row (flexible)
	contentRow := fb.NewRow()
	contentCell := flexbox.NewCell(1, 10) // Takes most vertical space
	contentCell.SetContent(m.renderContent(m.windowWidth, m.windowHeight-4))
	contentRow.AddCells(contentCell)
	
	// Status bar row
	statusRow := fb.NewRow()
	statusCell := flexbox.NewCell(1, 1)
	statusCell.SetContent(m.renderStatusBar())
	statusRow.AddCells(statusCell)
	
	fb.AddRows([]*flexbox.Row{tabRow, contentRow, statusRow})
	return fb.Render()
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
	
	// Create flexbox for tab bar
	fb := flexbox.NewHorizontal(m.windowWidth, 2)
	row := fb.NewColumn()
	
	// Create cells for each tab with equal ratios
	for i, tab := range tabs {
		cell := flexbox.NewCell(1, 1) // Equal width for all tabs
		
		var tabStyle lipgloss.Style
		if tab.view == m.currentView {
			// Active tab
			tabStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent).
				Background(lipgloss.Color("236")).
				Padding(0, 2).
				Align(lipgloss.Center)
		} else {
			// Inactive tab
			tabStyle = lipgloss.NewStyle().
				Foreground(colorTextDim).
				Background(lipgloss.Color("235")).
				Padding(0, 2).
				Align(lipgloss.Center)
		}
		
		// Add left margin for first tab
		if i == 0 {
			tabStyle = tabStyle.MarginLeft(1)
		}
		
		tabContent := fmt.Sprintf("%s %s", tab.icon, tab.title)
		cell.SetContent(tabStyle.Render(tabContent))
		row.AddCells(cell)
	}
	
	fb.AddColumns([]*flexbox.Column{row})
	
	// Wrap with bottom border
	return lipgloss.NewStyle().
		Width(m.windowWidth).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(colorBorder).
		Render(fb.Render())
}

// renderContent renders the main content area
func (m Model) renderContent(width, height int) string {
	LogDebug("renderContent called, view=%v, width=%d, height=%d", m.currentView, width, height)
	
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
		LogDebug("Calling renderExtensionDetail")
		content = m.renderExtensionDetail(innerWidth, height)
	}
	
	LogDebug("renderContent returning, content length=%d", len(content))
	
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
	// Create main flexbox
	fb := flexbox.New(width, height)
	
	// Header row
	headerRow := fb.NewRow()
	headerCell := flexbox.NewCell(1, 1)
	headerCell.SetContent(h1Style.Render("Extensions"))
	headerRow.AddCells(headerCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Search bar row (if active)
	if m.searchActive || m.searchBar.Value() != "" {
		searchRow := fb.NewRow()
		searchCell := flexbox.NewCell(1, 1)
		searchBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorderFocus).
			Padding(0, 1).
			Width(60).
			Render(m.searchBar.View())
		searchCell.SetContent(searchBox)
		searchRow.AddCells(searchCell)
		// Fixed height rows not supported - use cell ratios instead
		fb.AddRows([]*flexbox.Row{searchRow})
	}
	
	// Count row
	countRow := fb.NewRow()
	countCell := flexbox.NewCell(1, 1)
	var count string
	if m.searchBar.Value() != "" {
		count = fmt.Sprintf("%d of %d extensions (filtered)", len(m.filteredExtensions), len(m.extensions))
	} else {
		count = fmt.Sprintf("%d extensions found", len(m.filteredExtensions))
	}
	countCell.SetContent(textDimStyle.Render(count))
	countRow.AddCells(countCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Content area (flexible)
	contentRow := fb.NewRow()
	contentCell := flexbox.NewCell(1, 10) // Takes most space
	
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
		contentCell.SetContent(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(emptyBox))
	} else {
		// Extension list with cards
		var cards []string
		for i, ext := range m.filteredExtensions {
			isSelected := i == m.extensionsCursor
			card := m.renderExtensionCard(ext, isSelected, width-4)
			cards = append(cards, card)
			if i < len(m.filteredExtensions)-1 {
				cards = append(cards, "") // Spacing between cards
			}
		}
		contentCell.SetContent(strings.Join(cards, "\n"))
	}
	contentRow.AddCells(contentCell)
	
	// Help text row
	if len(m.filteredExtensions) > 0 {
		helpRow := fb.NewRow()
		helpCell := flexbox.NewCell(1, 1)
		helpText := renderKeyHelp([][2]string{
			{"↵", "Details"},
			{"n", "Install"},
			{"d", "Delete"},
			{"/", "Search"},
			{"Tab", "Next"},
		})
		helpCell.SetContent(helpText)
		helpRow.AddCells(helpCell)
		// Fixed height rows not supported - use cell ratios instead
		fb.AddRows([]*flexbox.Row{helpRow})
	}
	
	// Add all rows to flexbox
	rows := []*flexbox.Row{headerRow, countRow, contentRow}
	fb.AddRows(rows)
	
	return fb.Render()
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
	
	// Clean extension data
	cleanName := stripANSI(ext.Name)
	cleanVersion := stripANSI(ext.Version)
	cleanDescription := stripANSI(ext.Description)
	
	// Extension name and version on same line
	nameStyle := textStyle
	if isSelected {
		nameStyle = nameStyle.Bold(true).Foreground(colorAccent)
	}
	
	// Calculate available width for text (accounting for padding and borders)
	textWidth := width - 6  // 2 for borders, 4 for padding
	
	// Calculate space for name and version
	versionText := fmt.Sprintf("v%s", cleanVersion)
	nameWidth := textWidth - lipgloss.Width(versionText) - 2 // 2 for spacing
	
	name := nameStyle.MaxWidth(nameWidth).Render(cleanName)
	version := textDimStyle.Render(versionText)
	header := lipgloss.JoinHorizontal(lipgloss.Top, name, "  ", version)
	
	// Description on second line
	desc := textDimStyle.MaxWidth(textWidth).Render(cleanDescription)
	
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
	// Create main flexbox
	fb := flexbox.New(width, height)
	
	// Header row
	headerRow := fb.NewRow()
	headerCell := flexbox.NewCell(1, 1)
	headerCell.SetContent(h1Style.Render("Profiles"))
	headerRow.AddCells(headerCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Active profile badge row
	badgeRow := fb.NewRow()
	badgeCell := flexbox.NewCell(1, 1)
	activeProfile := "None"
	if m.currentProfile != nil {
		activeProfile = m.currentProfile.Name
	}
	badgeText := fmt.Sprintf("● Active: %s", activeProfile)
	activeBadge := lipgloss.NewStyle().
		Background(colorSuccess).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		MaxWidth(width - 4).
		Render(badgeText)
	badgeCell.SetContent(activeBadge)
	badgeRow.AddCells(badgeCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Search bar row (if active)
	if m.searchActive || m.searchBar.Value() != "" {
		searchRow := fb.NewRow()
		searchCell := flexbox.NewCell(1, 1)
		searchBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorderFocus).
			Padding(0, 1).
			Width(60).
			Render(m.searchBar.View())
		searchCell.SetContent(searchBox)
		searchRow.AddCells(searchCell)
		// Fixed height rows not supported - use cell ratios instead
		fb.AddRows([]*flexbox.Row{searchRow})
	}
	
	// Count row
	countRow := fb.NewRow()
	countCell := flexbox.NewCell(1, 1)
	var count string
	if m.searchBar.Value() != "" {
		count = fmt.Sprintf("%d of %d profiles (filtered)", len(m.filteredProfiles), len(m.profiles))
	} else {
		count = fmt.Sprintf("%d profiles", len(m.profiles))
	}
	countCell.SetContent(textDimStyle.Render(count))
	countRow.AddCells(countCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Content area (flexible)
	contentRow := fb.NewRow()
	contentCell := flexbox.NewCell(1, 10) // Takes most space
	
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
		contentCell.SetContent(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(emptyBox))
	} else {
		// Profile list
		var cards []string
		for i, prof := range m.filteredProfiles {
			isSelected := i == m.profilesCursor
			isActive := m.currentProfile != nil && prof.ID == m.currentProfile.ID
			card := m.renderProfileCard(prof, isSelected, isActive, width-4)
			cards = append(cards, card)
			if i < len(m.filteredProfiles)-1 {
				cards = append(cards, "")
			}
		}
		contentCell.SetContent(strings.Join(cards, "\n"))
	}
	contentRow.AddCells(contentCell)
	
	// Help text row
	if len(m.filteredProfiles) > 0 {
		helpRow := fb.NewRow()
		helpCell := flexbox.NewCell(1, 1)
		helpText := renderKeyHelp([][2]string{
			{"↵", "Activate"},
			{"n", "New"},
			{"e", "Edit"},
			{"d", "Delete"},
			{"/", "Search"},
		})
		helpCell.SetContent(helpText)
		helpRow.AddCells(helpCell)
		// Fixed height rows not supported - use cell ratios instead
		fb.AddRows([]*flexbox.Row{helpRow})
	}
	
	// Add all rows to flexbox
	rows := []*flexbox.Row{headerRow, badgeRow, countRow, contentRow}
	fb.AddRows(rows)
	
	return fb.Render()
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
	LogDebug("renderExtensionDetail called, width=%d, height=%d", width, height)
	
	if m.selectedExtension == nil {
		LogDebug("No extension selected")
		return "No extension selected"
	}
	
	ext := m.selectedExtension
	LogDebug("Rendering detail for extension: %s", ext.Name)
	
	// Clean all extension data of ANSI sequences
	cleanName := stripANSI(ext.Name)
	cleanVersion := stripANSI(ext.Version)
	cleanDescription := stripANSI(ext.Description)
	
	// Create main flexbox
	fb := flexbox.New(width, height)
	
	// Header row
	headerRow := fb.NewRow()
	headerCell := flexbox.NewCell(1, 1)
	headerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Padding(1).
		Width(width-2)
	
	headerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		h1Style.Copy().MarginBottom(0).Render(fmt.Sprintf("📦 %s", cleanName)),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			accentStyle.Bold(true).Render(fmt.Sprintf("v%s", cleanVersion)),
			textDimStyle.Render(" • "),
			textDimStyle.Render(cleanDescription),
		),
	)
	headerCell.SetContent(headerBox.Render(headerContent))
	headerRow.AddCells(headerCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Back navigation row
	backRow := fb.NewRow()
	backCell := flexbox.NewCell(1, 1)
	backCell.SetContent(textDimStyle.Copy().MarginLeft(2).Render("← Press Esc to go back"))
	backRow.AddCells(backCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Two-column layout row for basic info and MCP servers
	columnsRow := fb.NewRow()
	
	// Create horizontal flexbox for columns
	columnsHfb := flexbox.NewHorizontal(width, 0) // Height will be inherited
	
	// Left column (basic info)
	leftCol := columnsHfb.NewColumn()
	leftCell := flexbox.NewCell(1, 1) // Equal width
	leftCell.SetContent(m.renderExtDetailLeftColumn(ext, 0)) // Width will be calculated by flexbox
	leftCol.AddCells(leftCell)
	
	// Right column (MCP servers)
	rightCol := columnsHfb.NewColumn()
	rightCell := flexbox.NewCell(1, 1) // Equal width
	rightCell.SetContent(m.renderExtDetailRightColumn(ext, 0)) // Width will be calculated by flexbox
	rightCol.AddCells(rightCell)
	
	columnsHfb.AddColumns([]*flexbox.Column{leftCol, rightCol})
	
	// Add the horizontal flexbox content to the columns row
	columnsCell := flexbox.NewCell(1, 1)
	columnsCell.SetContent(columnsHfb.Render())
	columnsRow.AddCells(columnsCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Context file section row
	contextRow := fb.NewRow()
	contextCell := flexbox.NewCell(1, 3) // Takes more vertical space
	contextCell.SetContent(m.renderContextFileSection(ext, width-2))
	contextRow.AddCells(contextCell)
	
	// Action bar row
	actionRow := fb.NewRow()
	actionCell := flexbox.NewCell(1, 1)
	actionCell.SetContent(m.renderExtDetailActions(width-2))
	actionRow.AddCells(actionCell)
	// Fixed height rows not supported - use cell ratios instead
	
	// Add all rows
	fb.AddRows([]*flexbox.Row{headerRow, backRow, columnsRow, contextRow, actionRow})
	
	return fb.Render()
}

// renderExtDetailLeftColumn renders the left column of extension details
func (m Model) renderExtDetailLeftColumn(ext *extension.Extension, width int) string {
	var content strings.Builder
	
	// Basic Information
	content.WriteString(h2Style.Render("📋 Basic Information"))
	content.WriteString("\n\n")
	
	infoItems := []struct {
		label string
		value string
		icon  string
	}{
		{"ID", ext.ID, "🔑"},
		{"Path", ext.Path, "📁"},
		{"Type", "Extension", "🧩"},
	}
	
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1).
		Width(width)
	
	var infoContent strings.Builder
	for i, item := range infoItems {
		if i > 0 {
			infoContent.WriteString("\n\n")
		}
		
		// Label
		infoContent.WriteString(textDimStyle.Render(item.label))
		infoContent.WriteString("\n")
		
		// Value with icon
		valueText := fmt.Sprintf("%s %s", item.icon, item.value)
		// Truncate long paths
		if item.label == "Path" && len(valueText) > width-4 {
			valueText = "..." + valueText[len(valueText)-(width-7):]
		}
		infoContent.WriteString(textStyle.Render(valueText))
	}
	
	content.WriteString(infoBox.Render(infoContent.String()))
	
	return content.String()
}

// renderExtDetailRightColumn renders the right column of extension details
func (m Model) renderExtDetailRightColumn(ext *extension.Extension, width int) string {
	var content strings.Builder
	
	// MCP Servers
	content.WriteString(h2Style.Render("⚡ MCP Servers"))
	content.WriteString("\n\n")
	
	if ext.MCPServers != nil && len(ext.MCPServers) > 0 {
		mcpBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1).
			Width(width)
		
		var mcpContent strings.Builder
		serverIdx := 0
		for name, config := range ext.MCPServers {
			if serverIdx > 0 {
				mcpContent.WriteString("\n\n")
			}
			
			// Server name
			mcpContent.WriteString(accentStyle.Bold(true).Render(name))
			mcpContent.WriteString("\n")
			
			// Command
			cmdText := fmt.Sprintf("📟 %s", config.Command)
			if len(cmdText) > width-6 {
				cmdText = cmdText[:width-9] + "..."
			}
			mcpContent.WriteString(textStyle.Render(cmdText))
			
			// Show args count if present
			if len(config.Args) > 0 {
				mcpContent.WriteString("\n")
				mcpContent.WriteString(textDimStyle.Render(fmt.Sprintf("   %d args", len(config.Args))))
			}
			
			serverIdx++
		}
		
		content.WriteString(mcpBox.Render(mcpContent.String()))
	} else {
		// No servers box
		noServersBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1).
			Width(width).
			Align(lipgloss.Center)
		
		noServersContent := lipgloss.JoinVertical(
			lipgloss.Center,
			textDimStyle.Render("No MCP servers"),
			textDimStyle.Render("configured"),
		)
		
		content.WriteString(noServersBox.Render(noServersContent))
	}
	
	return content.String()
}

// renderContextFileSection renders the context file section
func (m Model) renderContextFileSection(ext *extension.Extension, width int) string {
	var content strings.Builder
	
	content.WriteString(h2Style.Render("📄 Context File"))
	content.WriteString("\n\n")
	
	contextFileName := ext.ContextFileName
	if contextFileName == "" {
		contextFileName = "GEMINI.md"
	}
	
	// Try to read and display context file content
	contextPath := filepath.Join(ext.Path, contextFileName)
	LogDebug("Checking context file at: %s", contextPath)
	
	contextBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1).
		Width(width).
		MaxHeight(12) // Limit height
	
	if fileContent, err := os.ReadFile(contextPath); err == nil && len(fileContent) > 0 {
		LogDebug("Context file found, size: %d bytes", len(fileContent))
		
		// File info header
		fileInfo := fmt.Sprintf("📝 %s (%d bytes)", contextFileName, len(fileContent))
		content.WriteString(textDimStyle.Render(fileInfo))
		content.WriteString("\n\n")
		
		// Render content
		if m.markdownRenderer != nil {
			LogDebug("Using cached glamour renderer")
			rendered, err := m.markdownRenderer.Render(string(fileContent))
			if err == nil {
				content.WriteString(contextBox.Render(rendered))
			} else {
				// Fallback to plain text
				plainText := string(fileContent)
				if len(plainText) > 500 {
					plainText = plainText[:500] + "\n\n... (truncated)"
				}
				content.WriteString(contextBox.Render(plainText))
			}
		} else {
			// No renderer, show plain text
			plainText := string(fileContent)
			if len(plainText) > 500 {
				plainText = plainText[:500] + "\n\n... (truncated)"
			}
			content.WriteString(contextBox.Render(plainText))
		}
	} else {
		// No file found
		content.WriteString(textDimStyle.Render(fmt.Sprintf("📝 %s", contextFileName)))
		content.WriteString("\n\n")
		
		noFileContent := lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			textDimStyle.Render("No context file found"),
			"",
			textDimStyle.Render("Create a "+contextFileName+" file to add"),
			textDimStyle.Render("documentation for this extension"),
			"",
		)
		
		content.WriteString(contextBox.Copy().Align(lipgloss.Center).Render(noFileContent))
	}
	
	return content.String()
}

// renderExtDetailActions renders the action bar for extension details
func (m Model) renderExtDetailActions(width int) string {
	actions := []struct {
		key   string
		label string
		style lipgloss.Style
	}{
		{"Esc", "Back", keyStyle},
		{"e", "Edit", keyStyle},
		{"d", "Delete", keyStyle.Copy().Foreground(colorError)},
	}
	
	var actionItems []string
	for _, action := range actions {
		item := lipgloss.JoinHorizontal(
			lipgloss.Top,
			action.style.Render(action.key),
			textDimStyle.Render(": "),
			textDimStyle.Render(action.label),
		)
		actionItems = append(actionItems, item)
	}
	
	actionBar := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(strings.Join(actionItems, "  •  "))
	
	return actionBar
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
	// Create flexbox for status bar
	fb := flexbox.NewHorizontal(m.windowWidth, 1)
	row := fb.NewColumn()
	
	// Left section (profile and extension count)
	var leftParts []string
	
	// Profile indicator
	if m.currentProfile != nil {
		leftParts = append(leftParts, fmt.Sprintf("👤 %s", m.currentProfile.Name))
	} else {
		leftParts = append(leftParts, "👤 No Profile")
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
	leftParts = append(leftParts, fmt.Sprintf("🧩 %d/%d", enabledCount, len(m.extensions)))
	
	leftCell := flexbox.NewCell(3, 1) // Takes 3/7 of width
	leftCell.SetContent(strings.Join(leftParts, " • "))
	
	// Middle section (error/info messages)
	middleCell := flexbox.NewCell(2, 1) // Takes 2/7 of width
	
	if m.err != nil {
		var errorMsg string
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
		middleCell.SetContent(lipgloss.NewStyle().Align(lipgloss.Center).Render(errorMsg))
	} else {
		middleCell.SetContent("")
	}
	
	// Right section (key hints)
	var hints []string
	hints = append(hints, "Tab: Switch")
	hints = append(hints, "L: Launch")
	hints = append(hints, "?: Help")
	hints = append(hints, "q: Quit")
	
	rightCell := flexbox.NewCell(2, 1) // Takes 2/7 of width
	rightCell.SetContent(lipgloss.NewStyle().Align(lipgloss.Right).Render(strings.Join(hints, " • ")))
	
	// Add cells to column
	row.AddCells(leftCell, middleCell, rightCell)
	fb.AddColumns([]*flexbox.Column{row})
	
	// Wrap with status bar style
	return statusBarStyle.
		Width(m.windowWidth).
		Render(fb.Render())
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