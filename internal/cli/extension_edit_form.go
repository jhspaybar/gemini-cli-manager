package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
)

// EditMode represents what we're editing
type EditMode int

const (
	EditModeConfig EditMode = iota
	EditModeContext
	EditModeJSON
)

// ExtensionEditForm is a form for editing extensions
type ExtensionEditForm struct {
	extension      *extension.Extension
	width          int
	height         int
	err            error
	mode           EditMode
	
	// Config mode fields
	inputs         []textinput.Model
	focusIndex     int
	
	// Context/JSON mode fields
	textarea       textarea.Model
	contextContent string
	jsonContent    string
	
	// Markdown preview
	renderer       *glamour.TermRenderer
	previewActive  bool
	
	// Callbacks
	onSave   func(*extension.Extension) tea.Cmd
	onCancel func() tea.Cmd
}

const (
	extNameField = iota
	extVersionField
	extDescriptionField
	totalExtConfigFields
)

// NewExtensionEditForm creates a new extension edit form
func NewExtensionEditForm(ext *extension.Extension) ExtensionEditForm {
	// Create text inputs for config fields
	inputs := make([]textinput.Model, totalExtConfigFields)
	
	inputs[extNameField] = textinput.New()
	inputs[extNameField].Placeholder = "Extension name"
	inputs[extNameField].SetValue(stripANSI(ext.Name))
	inputs[extNameField].CharLimit = 50
	inputs[extNameField].Focus()
	
	inputs[extVersionField] = textinput.New()
	inputs[extVersionField].Placeholder = "Version (e.g., 1.0.0)"
	inputs[extVersionField].SetValue(stripANSI(ext.Version))
	inputs[extVersionField].CharLimit = 20
	
	inputs[extDescriptionField] = textinput.New()
	inputs[extDescriptionField].Placeholder = "Brief description"
	inputs[extDescriptionField].SetValue(stripANSI(ext.Description))
	inputs[extDescriptionField].CharLimit = 200
	
	// Create textarea for editing context/JSON
	ta := textarea.New()
	ta.Placeholder = "Enter content..."
	ta.CharLimit = 50000
	ta.SetWidth(80)
	ta.SetHeight(20)
	
	// Create markdown renderer
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	
	form := ExtensionEditForm{
		extension:  ext,
		inputs:     inputs,
		mode:       EditModeConfig,
		textarea:   ta,
		renderer:   renderer,
	}
	
	// Load context file content
	contextPath := filepath.Join(ext.Path, ext.ContextFileName)
	if ext.ContextFileName == "" {
		contextPath = filepath.Join(ext.Path, "GEMINI.md")
	}
	if content, err := os.ReadFile(contextPath); err == nil {
		form.contextContent = string(content)
	}
	
	// Load JSON content
	jsonPath := filepath.Join(ext.Path, "gemini-extension.json")
	if content, err := os.ReadFile(jsonPath); err == nil {
		form.jsonContent = string(content)
	}
	
	return form
}

// Init initializes the form
func (f *ExtensionEditForm) Init() tea.Cmd {
	LogDebug("ExtensionEditForm.Init called, mode=%v", f.mode)
	
	// Return blink command for whichever component has focus
	if f.mode == EditModeConfig {
		return textinput.Blink
	}
	return textarea.Blink
}

// Update handles form updates
func (f *ExtensionEditForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	LogMessage("ExtensionEditForm.Update", msg)
	
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if f.onCancel != nil {
				return f, f.onCancel()
			}
			return f, nil
			
		case "tab":
			// Switch between modes
			if f.mode == EditModeConfig {
				// Move to next input field
				f.focusIndex++
				if f.focusIndex >= totalExtConfigFields {
					f.focusIndex = 0
				}
				for i := range f.inputs {
					if i == f.focusIndex {
						f.inputs[i].Focus()
					} else {
						f.inputs[i].Blur()
					}
				}
				return f, textinput.Blink
			}
			
		case "ctrl+t":
			// Toggle between config/context/json modes
			switch f.mode {
			case EditModeConfig:
				f.mode = EditModeContext
				f.textarea.SetValue(f.contextContent)
				f.textarea.Focus()
				return f, textarea.Blink
			case EditModeContext:
				f.mode = EditModeJSON
				f.textarea.SetValue(f.jsonContent)
				return f, textarea.Blink
			case EditModeJSON:
				f.mode = EditModeConfig
				f.inputs[f.focusIndex].Focus()
				return f, textinput.Blink
			}
			
		case "ctrl+p":
			// Toggle preview in context mode
			if f.mode == EditModeContext {
				f.previewActive = !f.previewActive
			}
			return f, nil
			
		case "ctrl+s":
			// Save changes
			return f, f.save()
		}
		
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		// Update textarea size
		f.textarea.SetWidth(f.width - 8)
		f.textarea.SetHeight(f.height - 10)
	}
	
	// Always update the appropriate component based on mode
	switch f.mode {
	case EditModeConfig:
		LogDebug("Updating inputs in config mode")
		// Update all inputs
		for i := range f.inputs {
			f.inputs[i], cmd = f.inputs[i].Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case EditModeContext, EditModeJSON:
		LogDebug("Updating textarea in %v mode", f.mode)
		// Update textarea
		f.textarea, cmd = f.textarea.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		
		// Save content back to appropriate field
		if f.mode == EditModeContext {
			f.contextContent = f.textarea.Value()
		} else {
			f.jsonContent = f.textarea.Value()
		}
	}
	
	result := tea.Batch(cmds...)
	LogDebug("ExtensionEditForm.Update returning with %d commands", len(cmds))
	return f, result
}

// View renders the form
func (f *ExtensionEditForm) View() string {
	LogDebug("ExtensionEditForm.View called, mode=%v, width=%d, height=%d", f.mode, f.width, f.height)
	
	// Use full width and height
	contentWidth := f.width
	contentHeight := f.height
	
	// Clean extension name
	cleanName := stripANSI(f.extension.Name)
	
	// Create header with title and mode tabs on same line
	titleText := fmt.Sprintf("✏️  Edit: %s", cleanName)
	titleStyle := h1Style.Copy().MarginBottom(0)
	title := titleStyle.Render(titleText)
	
	// Mode tabs
	tabs := f.renderModeTabs()
	
	// Create a two-column layout for the header
	headerLeft := title
	headerRight := tabs
	
	// Calculate spacing for header
	titleWidth := lipgloss.Width(headerLeft)
	tabsWidth := lipgloss.Width(headerRight)
	availableWidth := contentWidth - 4 // Account for padding
	spacing := availableWidth - titleWidth - tabsWidth
	if spacing < 2 {
		spacing = 2
	}
	
	header := headerLeft + strings.Repeat(" ", spacing) + headerRight
	
	// Build content area
	var contentArea strings.Builder
	
	// Calculate available space for content
	availableHeight := contentHeight - 8 // Header, spacing, help text
	
	// Create a bordered content area
	contentBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1).
		Width(availableWidth).
		Height(availableHeight)
	
	var innerContent string
	
	switch f.mode {
	case EditModeConfig:
		innerContent = f.renderConfigForm()
		
	case EditModeContext:
		if f.previewActive {
			innerContent = f.renderMarkdownPreview()
		} else {
			// Set textarea to use available space inside the box
			f.textarea.SetWidth(availableWidth - 4) // Account for border and padding
			f.textarea.SetHeight(availableHeight - 3)
			innerContent = f.renderTextarea()
		}
		
	case EditModeJSON:
		// Set textarea to use available space inside the box
		f.textarea.SetWidth(availableWidth - 4)
		f.textarea.SetHeight(availableHeight - 3)
		innerContent = f.renderTextarea()
	}
	
	contentArea.WriteString(contentBoxStyle.Render(innerContent))
	
	// Help text at bottom
	helpStyle := keyDescStyle.Copy().
		Width(availableWidth).
		Align(lipgloss.Center).
		MarginTop(1)
	
	helpText := "Ctrl+T: Switch Mode • Ctrl+S: Save • Esc: Cancel"
	if f.mode == EditModeContext {
		helpText += " • Ctrl+P: Toggle Preview"
	}
	help := helpStyle.Render(helpText)
	
	// Error display
	var errorDisplay string
	if f.err != nil {
		errorDisplay = errorStyle.Copy().
			Width(availableWidth).
			Align(lipgloss.Center).
			MarginTop(1).
			Render(f.err.Error())
	}
	
	// Combine all elements
	var fullContent strings.Builder
	fullContent.WriteString(header)
	fullContent.WriteString("\n\n")
	fullContent.WriteString(contentArea.String())
	fullContent.WriteString("\n")
	fullContent.WriteString(help)
	if errorDisplay != "" {
		fullContent.WriteString("\n")
		fullContent.WriteString(errorDisplay)
	}
	
	// Use full screen with minimal padding
	formStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(contentWidth).
		Height(contentHeight)
	
	return formStyle.Render(fullContent.String())
}

// renderConfigForm renders the configuration form fields
func (f *ExtensionEditForm) renderConfigForm() string {
	var b strings.Builder
	
	// Calculate available width
	availableWidth := f.width - 8 // Account for outer padding and border
	
	// Basic Information section
	basicInfoTitle := h2Style.Copy().MarginBottom(1).Render("📋 Basic Information")
	b.WriteString(basicInfoTitle)
	b.WriteString("\n\n")
	
	fields := []struct {
		label string
		index int
		help  string
		icon  string
	}{
		{"Name", extNameField, "The display name of your extension", "📝"},
		{"Version", extVersionField, "Semantic version (e.g., 1.0.0)", "🏷️"},
		{"Description", extDescriptionField, "Brief description of what your extension does", "📄"},
	}
	
	for i, field := range fields {
		// Create field container
		fieldStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1).
			Width(min(availableWidth, 70))
		
		if f.focusIndex == field.index {
			fieldStyle = fieldStyle.BorderForeground(colorAccent)
		}
		
		var fieldContent strings.Builder
		
		// Label with icon
		labelStyle := textStyle.Copy()
		if f.focusIndex == field.index {
			labelStyle = labelStyle.Foreground(colorAccent).Bold(true)
		}
		
		fieldContent.WriteString(fmt.Sprintf("%s %s", field.icon, labelStyle.Render(field.label)))
		fieldContent.WriteString("\n")
		
		// Input field (already has its own styling)
		inputView := f.inputs[field.index].View()
		fieldContent.WriteString(inputView)
		
		// Help text
		if f.focusIndex == field.index {
			fieldContent.WriteString("\n")
			fieldContent.WriteString(textDimStyle.Render(field.help))
		}
		
		b.WriteString(fieldStyle.Render(fieldContent.String()))
		
		if i < len(fields)-1 {
			b.WriteString("\n")
		}
	}
	
	b.WriteString("\n\n")
	
	// MCP Servers section
	mcpTitle := h2Style.Copy().MarginBottom(1).Render("⚡ MCP Servers")
	b.WriteString(mcpTitle)
	b.WriteString("\n\n")
	
	if f.extension.MCPServers != nil && len(f.extension.MCPServers) > 0 {
		// MCP servers box
		mcpBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1).
			Width(min(availableWidth, 70))
		
		var mcpContent strings.Builder
		
		serverCount := 0
		for name, server := range f.extension.MCPServers {
			if serverCount > 0 {
				mcpContent.WriteString("\n")
			}
			
			// Server name
			mcpContent.WriteString(accentStyle.Bold(true).Render(name))
			mcpContent.WriteString("\n")
			
			// Server details with proper indentation
			detailStyle := textStyle.Copy().MarginLeft(2)
			mcpContent.WriteString(detailStyle.Render(fmt.Sprintf("Command: %s", server.Command)))
			
			if len(server.Args) > 0 {
				mcpContent.WriteString("\n")
				mcpContent.WriteString(detailStyle.Render(fmt.Sprintf("Args: %v", server.Args)))
			}
			
			if len(server.Env) > 0 {
				mcpContent.WriteString("\n")
				mcpContent.WriteString(detailStyle.Render(fmt.Sprintf("Env: %v", server.Env)))
			}
			
			serverCount++
		}
		
		b.WriteString(mcpBox.Render(mcpContent.String()))
		b.WriteString("\n")
		
		// Tip about editing
		tipBox := lipgloss.NewStyle().
			Foreground(colorTextDim).
			MarginLeft(2)
		b.WriteString(tipBox.Render("💡 Tip: Switch to JSON mode to edit MCP server configuration"))
	} else {
		// No MCP servers message
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Align(lipgloss.Center).
			Width(min(availableWidth, 60))
		
		emptyContent := lipgloss.JoinVertical(
			lipgloss.Center,
			textDimStyle.Render("No MCP servers configured"),
			"",
			textDimStyle.Render("Switch to JSON mode to add MCP servers"),
		)
		
		b.WriteString(emptyBox.Render(emptyContent))
	}
	
	return b.String()
}

// renderTextarea renders the textarea for editing
func (f *ExtensionEditForm) renderTextarea() string {
	return f.textarea.View()
}

// renderMarkdownPreview renders the markdown preview
func (f *ExtensionEditForm) renderMarkdownPreview() string {
	if f.renderer == nil {
		return "Preview not available"
	}
	
	preview, err := f.renderer.Render(f.contextContent)
	if err != nil {
		return fmt.Sprintf("Preview error: %v", err)
	}
	
	// Wrap in a scrollable container
	return lipgloss.NewStyle().
		MaxHeight(f.height - 10).
		MaxWidth(min(f.width-10, 80)).
		Render(preview)
}

// renderModeTabs renders the mode selection tabs
func (f *ExtensionEditForm) renderModeTabs() string {
	modes := []struct {
		mode EditMode
		name string
		icon string
	}{
		{EditModeConfig, "Config", "⚙️"},
		{EditModeContext, "Context", "📝"},
		{EditModeJSON, "JSON", "{ }"},
	}
	
	var tabs []string
	for i, m := range modes {
		style := lipgloss.NewStyle().
			Padding(0, 1).
			MarginRight(0)
		
		if m.mode == f.mode {
			// Active tab
			style = style.
				Bold(true).
				Foreground(colorAccent).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(colorAccent)
		} else {
			// Inactive tab
			style = style.
				Foreground(colorTextDim).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(colorBorder)
		}
		
		tabContent := fmt.Sprintf("%s %s", m.icon, m.name)
		
		// Add separators between tabs
		if i > 0 {
			tabs = append(tabs, textDimStyle.Render(" "))
		}
		
		tabs = append(tabs, style.Render(tabContent))
	}
	
	// Create tab container with underline
	tabContainer := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(colorBorder).
		PaddingBottom(0)
	
	return tabContainer.Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
}


// save saves the extension changes
func (f *ExtensionEditForm) save() tea.Cmd {
	// Update extension from form fields
	f.extension.Name = f.inputs[extNameField].Value()
	f.extension.Version = f.inputs[extVersionField].Value()
	f.extension.Description = f.inputs[extDescriptionField].Value()
	
	// Save JSON if in JSON mode or if it was edited
	if f.mode == EditModeJSON || f.jsonContent != "" {
		// Parse and validate JSON
		var config extension.Extension
		if err := json.Unmarshal([]byte(f.jsonContent), &config); err != nil {
			f.err = fmt.Errorf("Invalid JSON: %v", err)
			return nil
		}
		
		// Update extension with parsed values
		f.extension.Name = config.Name
		f.extension.Version = config.Version
		f.extension.Description = config.Description
		f.extension.MCPServers = config.MCPServers
		f.extension.ContextFileName = config.ContextFileName
	}
	
	// Save files
	return func() tea.Msg {
		// Save gemini-extension.json
		jsonPath := filepath.Join(f.extension.Path, "gemini-extension.json")
		data, err := json.MarshalIndent(struct {
			Name            string                       `json:"name"`
			Version         string                       `json:"version"`
			Description     string                       `json:"description,omitempty"`
			MCPServers      map[string]extension.MCPServer `json:"mcpServers,omitempty"`
			ContextFileName string                       `json:"contextFileName,omitempty"`
		}{
			Name:            f.extension.Name,
			Version:         f.extension.Version,
			Description:     f.extension.Description,
			MCPServers:      f.extension.MCPServers,
			ContextFileName: f.extension.ContextFileName,
		}, "", "  ")
		
		if err != nil {
			return UIError{
				Type:    ErrorTypeFileSystem,
				Message: "Failed to marshal extension",
				Details: err.Error(),
			}
		}
		
		if err := os.WriteFile(jsonPath, data, 0644); err != nil {
			return UIError{
				Type:    ErrorTypeFileSystem,
				Message: "Failed to save extension config",
				Details: err.Error(),
			}
		}
		
		// Save context file if it was edited
		if f.contextContent != "" {
			contextPath := filepath.Join(f.extension.Path, f.extension.ContextFileName)
			if f.extension.ContextFileName == "" {
				contextPath = filepath.Join(f.extension.Path, "GEMINI.md")
			}
			
			if err := os.WriteFile(contextPath, []byte(f.contextContent), 0644); err != nil {
				return UIError{
					Type:    ErrorTypeFileSystem,
					Message: "Failed to save context file",
					Details: err.Error(),
				}
			}
		}
		
		// Call the save callback
		if f.onSave != nil {
			return extensionSavedMsg{extension: f.extension}
		}
		
		return nil
	}
}

// SetSize updates the form dimensions
func (f *ExtensionEditForm) SetSize(width, height int) {
	f.width = width
	f.height = height
	// Set textarea to use most of the available space
	f.textarea.SetWidth(width - 8)
	f.textarea.SetHeight(height - 10)
}

// SetCallbacks sets the form callbacks
func (f *ExtensionEditForm) SetCallbacks(onSave func(*extension.Extension) tea.Cmd, onCancel func() tea.Cmd) {
	f.onSave = onSave
	f.onCancel = onCancel
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}