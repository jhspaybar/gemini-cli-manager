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
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
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

	// Render the entire app as a card with tabs
	return m.renderAppCard()
}

// renderAppCard renders the entire app as a card with tabs
func (m Model) renderAppCard() string {
	// Ensure we have valid dimensions
	width := m.windowWidth
	height := m.windowHeight
	if width == 0 || height == 0 {
		// Use reasonable defaults if window size not yet received
		width = 80
		height = 24
	}
	
	// Calculate dimensions with proper padding
	horizontalPadding := 2  // Reduced to better align components
	verticalPadding := 0  // No vertical padding to maximize space
	contentWidth := width - (horizontalPadding * 2)
	contentHeight := height - (verticalPadding * 2)
	
	LogDebug("renderAppCard: windowSize=%dx%d, contentSize=%dx%d", 
		width, height, contentWidth, contentHeight)

	// Create tabs and get the tab bar instance
	tabBar := m.createTabBar(contentWidth)
	
	// Render status bar with same width as content to ensure alignment
	statusBar := m.renderStatusBar(contentWidth)
	statusHeight := lipgloss.Height(statusBar)
	
	// Render tabs separately to get height
	tabs := tabBar.Render()
	tabHeight := lipgloss.Height(tabs)
	
	// Calculate height for content box
	contentBoxHeight := contentHeight - statusHeight - tabHeight - 1 // -1 for newline between content and status
	
	// Ensure we have at least some height for content
	if contentBoxHeight < 5 {
		contentBoxHeight = 5
	}
	
	LogDebug("Heights: content=%d, status=%d, tabs=%d, contentBox=%d", 
		contentHeight, statusHeight, tabHeight, contentBoxHeight)
	
	// Render content  
	content := m.renderContent(contentWidth - 4) // Account for borders and padding
	
	// Use TabBar's RenderWithContent for seamless connection
	tabsAndContent := tabBar.RenderWithContent(content, contentBoxHeight)
	
	// Combine tabs+content with status bar
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabsAndContent,
		statusBar,
	)

	return lipgloss.NewStyle().
		Padding(verticalPadding, horizontalPadding).
		Render(fullContent)
}

// createTabBar creates and configures the tab bar component
func (m Model) createTabBar(width int) *components.TabBar {
	// Define our tabs
	tabs := []components.Tab{
		{Title: "Extensions", Icon: "🧩", ID: "extensions"},
		{Title: "Profiles", Icon: "👤", ID: "profiles"},
		{Title: "Settings", Icon: "🔧", ID: "settings"},
		{Title: "Help", Icon: "❓", ID: "help"},
	}
	
	// Create tab bar
	tabBar := components.NewTabBar(tabs, width)
	tabBar.SetStyles(activeTabStyle, inactiveTabStyle, colorBorder)
	
	// Set active tab based on current view
	switch m.currentView {
	case ViewExtensions:
		tabBar.SetActiveByID("extensions")
	case ViewProfiles:
		tabBar.SetActiveByID("profiles")
	case ViewSettings:
		tabBar.SetActiveByID("settings")
	case ViewHelp:
		tabBar.SetActiveByID("help")
	}
	
	return tabBar
}


// renderStatusBar renders the status bar component
func (m Model) renderStatusBar(width int) string {
	// Create and populate status bar component
	statusBar := components.NewStatusBar(width)
	
	// Profile and extension count
	profileName := "No Profile"
	if m.currentProfile != nil {
		profileName = m.currentProfile.Name
	}
	enabledCount := 0
	if m.currentProfile != nil {
		for _, extRef := range m.currentProfile.Extensions {
			if extRef.Enabled {
				enabledCount++
			}
		}
	}
	statusBar.SetLeftItems(components.ProfileStatusItems(profileName, enabledCount, len(m.extensions)))
	
	// Error/info messages
	if m.err != nil {
		if uiErr, ok := m.err.(UIError); ok {
			var msgType components.ErrorType
			if uiErr.Type == ErrorTypeInfo {
				msgType = components.ErrorTypeInfo
			} else {
				msgType = components.ErrorTypeError
			}
			statusBar.SetErrorMessage(components.ErrorMessage{
				Type:    msgType,
				Message: uiErr.Message,
				Details: uiErr.Details,
			})
		} else {
			statusBar.SetErrorMessage(components.ErrorMessage{
				Type:    components.ErrorTypeError,
				Message: m.err.Error(),
			})
		}
	}
	
	// Key bindings
	statusBar.SetKeyBindings(components.CommonKeyBindings())
	
	return statusBar.Render()
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}


// renderTabBar renders the top navigation tabs (kept for compatibility)
func (m Model) renderTabBar() string {
	// Don't show tabs in detail view
	if m.currentView == ViewExtensionDetail {
		return ""
	}

	tabs := []struct {
		title string
		icon  string
		view  ViewType
	}{
		{"Extensions", "🧩", ViewExtensions},
		{"Profiles", "👤", ViewProfiles},
		{"Settings", "🔧", ViewSettings},
		{"Help", "❓", ViewHelp},
	}

	// Render individual tabs
	var renderedTabs []string
	for _, tab := range tabs {
		isActive := tab.view == m.currentView
		
		var style lipgloss.Style
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		
		content := fmt.Sprintf("%s %s", tab.icon, tab.title)
		renderedTabs = append(renderedTabs, style.Render(content))
	}
	
	// Join tabs horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// renderContent renders the main content area
func (m Model) renderContent(width int) string {
	LogDebug("renderContent called, view=%v, width=%d", m.currentView, width)

	var content string

	// Create a content container with padding - no fixed height
	contentContainer := lipgloss.NewStyle().
		Padding(2, 3).
		MaxWidth(width)
	
	// Calculate available width after padding
	availableWidth := width - 6  // 3 padding on each side
	
	switch m.currentView {
	case ViewExtensions:
		content = m.renderExtensions(availableWidth)
	case ViewProfiles:
		content = m.renderProfiles(availableWidth)
	case ViewSettings:
		content = m.renderSettings(availableWidth)
	case ViewHelp:
		content = m.renderHelp(availableWidth)
	case ViewExtensionDetail:
		LogDebug("Calling renderExtensionDetail")
		content = m.renderExtensionDetail(availableWidth)
	}

	LogDebug("renderContent returning, content length=%d", len(content))

	// Apply the container style
	return contentContainer.Render(content)
}

// renderExtensions renders the extensions view
func (m Model) renderExtensions(width int) string {
	var lines []string

	// Header
	header := h1Style.Render("Extensions")
	lines = append(lines, header)

	// Show search bar if active or has query
	if m.searchActive || m.searchBar.Value() != "" {
		lines = append(lines, "")
		// Set search bar width and placeholder for extensions
		searchWidth := min(60, width)
		m.searchBar.SetWidth(searchWidth).SetPlaceholder("Search extensions by name or description...")
		lines = append(lines, m.searchBar.Render())
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
		var emptyState *components.EmptyState
		if m.searchBar.Value() != "" {
			// Search with no results
			emptyState = components.NewEmptyState(width).
				NoItemsFound().
				SetAction("Press '/' to modify your search")
		} else {
			// No extensions at all
			emptyState = components.NewEmptyState(width).
				SetIcon("📦").
				SetTitle("No extensions installed").
				SetAction("Press 'n' to install your first extension")
		}
		
		lines = append(lines, "")
		lines = append(lines, emptyState.Render())
	} else {
		// Extension list with cards
		for i, ext := range m.filteredExtensions {
			isSelected := i == m.extensionsCursor
			card := m.RenderExtensionCard(ext, isSelected, width)
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

// RenderExtensionCard renders a single extension as a card using the Card component
// Exported for testing
func (m Model) RenderExtensionCard(ext *extension.Extension, isSelected bool, width int) string {
	// Clean extension data
	cleanName := stripANSI(ext.Name)
	cleanVersion := stripANSI(ext.Version)
	cleanDescription := stripANSI(ext.Description)

	// Create card
	card := components.NewCard(width).
		SetTitle(cleanName, "🧩").
		SetSubtitle(fmt.Sprintf("v%s", cleanVersion)).
		SetDescription(cleanDescription).
		SetSelected(isSelected)

	// Add MCP servers info if present
	if ext.MCPServers != nil && len(ext.MCPServers) > 0 {
		count := len(ext.MCPServers)
		card.AddMetadata("MCP Servers", fmt.Sprintf("%d server%s", count, pluralize(count)), "⚡")
	}

	return card.Render()
}

// renderProfiles renders the profiles view
func (m Model) renderProfiles(width int) string {
	var lines []string

	// Header
	header := h1Style.Render("Profiles")
	lines = append(lines, header)

	// Active profile badge
	activeProfile := "None"
	if m.currentProfile != nil {
		activeProfile = m.currentProfile.Name
	}
	badgeText := fmt.Sprintf("● Active: %s", activeProfile)
	activeBadge := lipgloss.NewStyle().
		Background(colorSuccess).
		Foreground(theme.Black()).
		Bold(true).
		Padding(0, 1).
		MaxWidth(width).
		Render(badgeText)
	lines = append(lines, "", activeBadge)

	// Show search bar if active or has query
	if m.searchActive || m.searchBar.Value() != "" {
		lines = append(lines, "")
		// Set search bar width and placeholder for profiles
		searchWidth := min(60, width)
		m.searchBar.SetWidth(searchWidth).SetPlaceholder("Search profiles by name or tags...")
		lines = append(lines, m.searchBar.Render())
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
		var emptyState *components.EmptyState
		if m.searchBar.Value() != "" {
			// Search with no results
			emptyState = components.NewEmptyState(width).
				NoItemsFound().
				SetAction("Press '/' to modify your search")
		} else {
			// No profiles at all
			emptyState = components.NewEmptyState(width).
				SetIcon("👤").
				SetTitle("No profiles configured").
				SetAction("Press 'n' to create your first profile")
		}
		
		lines = append(lines, "")
		lines = append(lines, emptyState.Render())
	} else {
		// Profile list with cards
		for i, prof := range m.filteredProfiles {
			isSelected := i == m.profilesCursor
			isActive := m.currentProfile != nil && prof.ID == m.currentProfile.ID
			card := m.RenderProfileCard(prof, isSelected, isActive, width)
			lines = append(lines, card)
			if i < len(m.filteredProfiles)-1 {
				lines = append(lines, "") // Spacing between cards
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

// RenderProfileCard renders a single profile as a card using the Card component
// Exported for testing
func (m Model) RenderProfileCard(prof *profile.Profile, isSelected, isActive bool, width int) string {
	// Create card
	card := components.NewCard(width).
		SetTitle(prof.Name, "").  // No icon for profiles
		SetDescription(prof.Description).
		SetSelected(isSelected).
		SetActive(isActive)

	// Add extension count
	if len(prof.Extensions) > 0 {
		card.AddMetadata("Extensions", fmt.Sprintf("%d extension%s", len(prof.Extensions), pluralize(len(prof.Extensions))), "📦")
	}

	return card.Render()
}

// renderSettings renders the settings view
func (m Model) renderSettings(width int) string {
	var lines []string

	header := h1Style.Render("Settings")
	lines = append(lines, header, "")

	// Theme section with selectable list
	lines = append(lines, h2Style.Render("🎨 Appearance"))
	lines = append(lines, "")

	// Get available themes
	themes := theme.GetAvailableThemes()
	currentTheme := theme.GetCurrentTheme()

	// Show theme list
	for i, themeName := range themes {
		var line string
		prefix := "  "
		style := textStyle

		if i == m.settingsCursor {
			prefix = "▶ "
			style = accentStyle.Bold(true)
		}

		// Add checkmark for current theme
		checkmark := "  "
		if themeName == currentTheme {
			checkmark = "✓ "
		}

		line = style.Render(fmt.Sprintf("%s%s%s", prefix, checkmark, themeName))
		lines = append(lines, line)
	}

	lines = append(lines, "", "")

	// Other settings sections (read-only for now)
	lines = append(lines, h2Style.Render("🔧 General"))
	lines = append(lines, textDimStyle.Render("  Gemini CLI Path:    /usr/local/bin/gemini"))
	lines = append(lines, textDimStyle.Render(fmt.Sprintf("  Config Directory:   %s", m.stateDir)))
	lines = append(lines, "")

	lines = append(lines, h2Style.Render("📦 Extensions"))
	lines = append(lines, textDimStyle.Render("  Extensions Directory: ~/.gemini/extensions"))
	lines = append(lines, textDimStyle.Render("  Auto-update:         Enabled"))
	lines = append(lines, "", "")

	// Help at bottom
	helpText := renderKeyHelp([][2]string{
		{"↑/↓", "Navigate"},
		{"Enter", "Apply Theme"},
		{"Tab", "Next"},
	})
	lines = append(lines, helpText)

	return strings.Join(lines, "\n")
}

// renderHelp renders the help view
func (m Model) renderHelp(width int) string {
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
func (m Model) renderExtensionDetail(width int) string {
	LogDebug("renderExtensionDetail called, width=%d", width)

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

	var lines []string

	// Header card using Card component - use full width
	headerCard := components.NewCard(width).
		SetTitle(cleanName, "📦").
		SetSubtitle(fmt.Sprintf("v%s", cleanVersion)).
		SetDescription(cleanDescription).
		SetFocused(true) // Use focused state for the header to make it stand out
	
	lines = append(lines, headerCard.Render())

	// Back navigation
	lines = append(lines, "")
	lines = append(lines, textDimStyle.Copy().MarginLeft(2).Render("← Press Esc to go back"))
	lines = append(lines, "")

	// Two-column layout for basic info and MCP servers
	// Use flexbox only for the horizontal layout
	columnsHfb := flexbox.NewHorizontal(width, 10) // Fixed height for columns section

	// Left column (basic info)
	leftCol := columnsHfb.NewColumn()
	leftCell := flexbox.NewCell(1, 1) // Equal width
	leftCell.SetContent(m.RenderExtDetailLeftColumn(ext, width/2))
	leftCol.AddCells(leftCell)

	// Right column (MCP servers)
	rightCol := columnsHfb.NewColumn()
	rightCell := flexbox.NewCell(1, 1) // Equal width
	rightCell.SetContent(m.RenderExtDetailRightColumn(ext, width/2))
	rightCol.AddCells(rightCell)

	columnsHfb.AddColumns([]*flexbox.Column{leftCol, rightCol})
	lines = append(lines, columnsHfb.Render())
	lines = append(lines, "")

	// Context file section
	lines = append(lines, m.RenderContextFileSection(ext, width))
	lines = append(lines, "")

	// Action bar
	lines = append(lines, m.renderExtDetailActions(width))

	return strings.Join(lines, "\n")
}

// RenderExtDetailLeftColumn renders the left column of extension details
// Exported for testing
func (m Model) RenderExtDetailLeftColumn(ext *extension.Extension, width int) string {
	var content strings.Builder

	// Basic Information card
	infoCard := components.NewCard(width).
		SetTitle("Basic Information", "📋")
	
	// Add metadata items
	infoCard.AddMetadata("ID", ext.ID, "🔑").
		AddMetadata("Type", "Extension", "🧩")
	
	// Handle long paths - truncate if needed
	pathValue := ext.Path
	if len(pathValue) > width-15 { // Leave room for icon and padding
		pathValue = "..." + pathValue[len(pathValue)-(width-18):]
	}
	infoCard.AddMetadata("Path", pathValue, "📁")
	
	content.WriteString(infoCard.Render())

	return content.String()
}

// RenderExtDetailRightColumn renders the right column of extension details
// Exported for testing
func (m Model) RenderExtDetailRightColumn(ext *extension.Extension, width int) string {
	// MCP Servers card
	mcpCard := components.NewCard(width).
		SetTitle("MCP Servers", "⚡")

	if ext.MCPServers != nil && len(ext.MCPServers) > 0 {
		// Build description with server details
		var serverDetails []string
		for name, config := range ext.MCPServers {
			// Server name and command
			cmdText := config.Command
			if len(cmdText) > width-20 { // Account for padding and name
				cmdText = cmdText[:width-23] + "..."
			}
			
			serverLine := fmt.Sprintf("%s\n📟 %s", 
				accentStyle.Bold(true).Render(name),
				cmdText)
			
			// Add args count if present
			if len(config.Args) > 0 {
				serverLine += fmt.Sprintf("\n   %d args", len(config.Args))
			}
			
			serverDetails = append(serverDetails, serverLine)
		}
		
		// Join all servers with spacing
		mcpCard.SetDescription(strings.Join(serverDetails, "\n\n"))
		
		// Use focused state to highlight this has content
		mcpCard.SetFocused(true)
	} else {
		// No servers - show empty state
		mcpCard.SetDescription("No MCP servers\nconfigured")
	}

	return mcpCard.Render()
}

// RenderContextFileSection renders the context file section
// Exported for testing
func (m Model) RenderContextFileSection(ext *extension.Extension, width int) string {
	contextFileName := ext.ContextFileName
	if contextFileName == "" {
		contextFileName = "GEMINI.md"
	}

	// Context File card
	contextCard := components.NewCard(width).
		SetTitle("Context File", "📄")

	// Try to read and display context file content
	contextPath := filepath.Join(ext.Path, contextFileName)
	LogDebug("Checking context file at: %s", contextPath)

	if fileContent, err := os.ReadFile(contextPath); err == nil && len(fileContent) > 0 {
		LogDebug("Context file found, size: %d bytes", len(fileContent))

		// Add file info as metadata
		contextCard.AddMetadata("File", fmt.Sprintf("%s (%d bytes)", contextFileName, len(fileContent)), "📝")

		// Render content
		var contentText string
		if m.markdownRenderer != nil {
			LogDebug("Using cached glamour renderer")
			rendered, err := m.markdownRenderer.Render(string(fileContent))
			if err == nil {
				contentText = rendered
			} else {
				// Fallback to plain text
				contentText = string(fileContent)
			}
		} else {
			// No renderer, show plain text
			contentText = string(fileContent)
		}
		
		// Truncate if too long
		if len(contentText) > 500 {
			contentText = contentText[:500] + "\n\n... (truncated)"
		}
		
		contextCard.SetDescription(contentText)
	} else {
		// No file found
		contextCard.AddMetadata("File", contextFileName, "📝")
		contextCard.SetDescription(fmt.Sprintf("No context file found\n\nCreate a %s file to add\ndocumentation for this extension", contextFileName))
	}

	return contextCard.Render()
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
	// Loading content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		"Loading extensions and profiles...",
		"",
		"Please wait",
	)

	// Use the Modal component
	modal := components.NewModal(m.windowWidth, m.windowHeight).
		SetTitle("Gemini CLI Manager", "🚀").
		SetContent(content).
		SetBorderColor(theme.Primary()).
		Small() // Use small preset for loading modal

	return modal.Render()
}

// renderStatusBarContent renders the content of the status bar
// Kept for backward compatibility - new code should use StatusBar component
func (m Model) renderStatusBarContent(width int) string {
	// Create StatusBar component
	statusBar := components.NewStatusBar(width)
	
	// Profile and extension count
	profileName := "No Profile"
	if m.currentProfile != nil {
		profileName = m.currentProfile.Name
	}
	enabledCount := 0
	if m.currentProfile != nil {
		for _, extRef := range m.currentProfile.Extensions {
			if extRef.Enabled {
				enabledCount++
			}
		}
	}
	statusBar.SetLeftItems(components.ProfileStatusItems(profileName, enabledCount, len(m.extensions)))
	
	// Error/info messages
	if m.err != nil {
		if uiErr, ok := m.err.(UIError); ok {
			var msgType components.ErrorType
			if uiErr.Type == ErrorTypeInfo {
				msgType = components.ErrorTypeInfo
			} else {
				msgType = components.ErrorTypeError
			}
			statusBar.SetErrorMessage(components.ErrorMessage{
				Type:    msgType,
				Message: uiErr.Message,
				Details: uiErr.Details,
			})
		} else {
			statusBar.SetErrorMessage(components.ErrorMessage{
				Type:    components.ErrorTypeError,
				Message: m.err.Error(),
			})
		}
	}
	
	// Key bindings
	statusBar.SetKeyBindings(components.CommonKeyBindings())
	
	// Return just the content without border styling
	return statusBar.RenderContent()
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
