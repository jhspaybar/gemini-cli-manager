package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/76creates/stickers/flexbox"
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
	extension *extension.Extension
	width     int
	height    int
	err       error
	mode      EditMode

	// Config mode fields
	inputs     []textinput.Model
	focusIndex int

	// Context/JSON mode fields
	textarea       textarea.Model
	contextContent string
	jsonContent    string

	// Markdown preview
	renderer      *glamour.TermRenderer
	previewActive bool

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
	// Don't set default size here - it will be set dynamically

	// Don't create renderer here - it's slow and we have a cached one in Model

	form := ExtensionEditForm{
		extension: ext,
		inputs:    inputs,
		mode:      EditModeConfig,
		textarea:  ta,
		renderer:  nil, // Will be set by parent model
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
		// Update textarea size with proper constraints
		if f.width > 10 && f.height > 10 {
			f.textarea.SetWidth(min(f.width, 120))   // Max width of 120
			f.textarea.SetHeight(min(f.height/2, 30)) // Max height of 30
		}
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

	if f.width == 0 || f.height == 0 {
		return "Loading..."
	}

	// Create main flexbox
	fb := flexbox.New(f.width, f.height)

	// Clean extension name
	cleanName := stripANSI(f.extension.Name)

	// Header row
	headerRow := fb.NewRow()
	headerCell := flexbox.NewCell(1, 1)
	headerCell.SetContent(h1Style.Render(fmt.Sprintf("Edit: %s", cleanName)))
	headerRow.AddCells(headerCell)

	// Tabs row
	tabsRow := fb.NewRow()
	tabsCell := flexbox.NewCell(1, 1)
	tabsCell.SetContent(f.renderModeTabs())
	tabsRow.AddCells(tabsCell)

	// Content row (takes most space)
	contentRow := fb.NewRow()
	contentCell := flexbox.NewCell(1, 5) // 5:1 ratio for content

	// Content based on mode
	var content string
	switch f.mode {
	case EditModeConfig:
		content = f.renderConfigForm()
	case EditModeContext:
		if f.previewActive {
			content = f.renderMarkdownPreview()
		} else {
			// Let flexbox handle sizing - just set max dimensions
			f.textarea.SetHeight(min(30, f.height/2)) // Use proportional height
			f.textarea.SetWidth(min(120, f.width))     // Let textarea use available width
			content = f.renderTextarea()
		}
	case EditModeJSON:
		// Let flexbox handle sizing - just set max dimensions
		f.textarea.SetHeight(min(30, f.height/2)) // Use proportional height
		f.textarea.SetWidth(min(120, f.width))     // Let textarea use available width
		content = f.renderTextarea()
	}
	contentCell.SetContent(content)
	contentRow.AddCells(contentCell)

	// Help row
	helpRow := fb.NewRow()
	helpCell := flexbox.NewCell(1, 1)
	helpText := "Ctrl+T: Mode • Ctrl+S: Save • Esc: Cancel"
	if f.mode == EditModeContext {
		helpText += " • Ctrl+P: Preview"
	}
	helpCell.SetContent(keyDescStyle.Render(helpText))
	helpRow.AddCells(helpCell)

	// Add all rows
	fb.AddRows([]*flexbox.Row{headerRow, tabsRow, contentRow, helpRow})

	// Add error row if needed
	if f.err != nil {
		errorRow := fb.NewRow()
		errorCell := flexbox.NewCell(1, 1)
		errorCell.SetContent(errorStyle.Render(f.err.Error()))
		errorRow.AddCells(errorCell)
		fb.AddRows([]*flexbox.Row{errorRow})
	}

	return fb.Render()
}

// renderConfigForm renders the configuration form fields
func (f *ExtensionEditForm) renderConfigForm() string {
	// Create flexbox for form layout - let parent determine size
	fb := flexbox.New(0, 0) // Size inherited from parent

	fields := []struct {
		label string
		index int
		help  string
	}{
		{"Name", extNameField, "Extension display name"},
		{"Version", extVersionField, "Semantic version (e.g., 1.0.0)"},
		{"Description", extDescriptionField, "Brief description"},
	}

	// Add form fields
	for _, field := range fields {
		fieldRow := fb.NewRow()

		// Label cell (30% width)
		labelCell := flexbox.NewCell(3, 1)
		labelStyle := textDimStyle
		if f.focusIndex == field.index {
			labelStyle = accentStyle
		}
		labelCell.SetContent(labelStyle.Render(field.label))

		// Input cell (70% width)
		inputCell := flexbox.NewCell(7, 1)
		inputContent := f.inputs[field.index].View()
		if f.focusIndex == field.index && field.help != "" {
			inputContent += "\n" + textDimStyle.Render(field.help)
		}
		inputCell.SetContent(inputContent)

		fieldRow.AddCells(labelCell, inputCell)
		fb.AddRows([]*flexbox.Row{fieldRow})

		// Add spacing row
		if field.index < len(fields)-1 {
			spacerRow := fb.NewRow()
			spacerCell := flexbox.NewCell(1, 1)
			spacerCell.SetContent(" ")
			spacerRow.AddCells(spacerCell)
			fb.AddRows([]*flexbox.Row{spacerRow})
		}
	}

	// MCP Servers section
	if f.extension.MCPServers != nil && len(f.extension.MCPServers) > 0 {
		// Spacer
		spacerRow := fb.NewRow()
		spacerCell := flexbox.NewCell(1, 1)
		spacerCell.SetContent("")
		spacerRow.AddCells(spacerCell)
		fb.AddRows([]*flexbox.Row{spacerRow})

		// MCP Header
		mcpHeaderRow := fb.NewRow()
		mcpHeaderCell := flexbox.NewCell(1, 1)
		mcpHeaderCell.SetContent(h2Style.Render("MCP Servers"))
		mcpHeaderRow.AddCells(mcpHeaderCell)
		fb.AddRows([]*flexbox.Row{mcpHeaderRow})

		// MCP Server list
		for name, server := range f.extension.MCPServers {
			serverRow := fb.NewRow()
			serverCell := flexbox.NewCell(1, 1)
			serverCell.SetContent(fmt.Sprintf("• %s: %s", accentStyle.Render(name), server.Command))
			serverRow.AddCells(serverCell)
			fb.AddRows([]*flexbox.Row{serverRow})
		}

		// Tip
		tipRow := fb.NewRow()
		tipCell := flexbox.NewCell(1, 1)
		tipCell.SetContent(textDimStyle.Render("Use JSON mode to edit servers"))
		tipRow.AddCells(tipCell)
		fb.AddRows([]*flexbox.Row{tipRow})
	}

	return fb.Render()
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

	// Wrap in a scrollable container with proper constraints
	return lipgloss.NewStyle().
		MaxHeight(f.height/2). // Use proportional height
		MaxWidth(80).          // Fixed max width for readability
		Render(preview)
}

// renderModeTabs renders the mode selection tabs
func (f *ExtensionEditForm) renderModeTabs() string {
	modes := []struct {
		mode EditMode
		name string
	}{
		{EditModeConfig, "Config"},
		{EditModeContext, "Context"},
		{EditModeJSON, "JSON"},
	}

	var tabs []string
	for _, m := range modes {
		style := textStyle
		if m.mode == f.mode {
			style = accentStyle.Copy().Bold(true).Underline(true)
		}
		tabs = append(tabs, style.Render(m.name))
	}

	return strings.Join(tabs, "  ")
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
			Name            string                         `json:"name"`
			Version         string                         `json:"version"`
			Description     string                         `json:"description,omitempty"`
			MCPServers      map[string]extension.MCPServer `json:"mcpServers,omitempty"`
			ContextFileName string                         `json:"contextFileName,omitempty"`
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
	// Set textarea with reasonable constraints
	if width > 10 && height > 10 {
		f.textarea.SetWidth(min(width, 120))   // Max width of 120
		f.textarea.SetHeight(min(height/2, 30)) // Max height of 30
	}
}

// SetCallbacks sets the form callbacks
func (f *ExtensionEditForm) SetCallbacks(onSave func(*extension.Extension) tea.Cmd, onCancel func() tea.Cmd) {
	f.onSave = onSave
	f.onCancel = onCancel
}

// SetRenderer sets the markdown renderer
func (f *ExtensionEditForm) SetRenderer(renderer *glamour.TermRenderer) {
	f.renderer = renderer
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
