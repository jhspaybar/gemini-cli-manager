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
	inputs[extNameField].SetValue(ext.Name)
	inputs[extNameField].CharLimit = 50
	inputs[extNameField].Focus()
	
	inputs[extVersionField] = textinput.New()
	inputs[extVersionField].Placeholder = "Version (e.g., 1.0.0)"
	inputs[extVersionField].SetValue(ext.Version)
	inputs[extVersionField].CharLimit = 20
	
	inputs[extDescriptionField] = textinput.New()
	inputs[extDescriptionField].Placeholder = "Brief description"
	inputs[extDescriptionField].SetValue(ext.Description)
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
func (f ExtensionEditForm) Init() tea.Cmd {
	// Return blink command for whichever component has focus
	if f.mode == EditModeConfig {
		return textinput.Blink
	}
	return textarea.Blink
}

// Update handles form updates
func (f ExtensionEditForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		// Update all inputs
		for i := range f.inputs {
			f.inputs[i], cmd = f.inputs[i].Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		
	case EditModeContext, EditModeJSON:
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
	
	return f, tea.Batch(cmds...)
}

// View renders the form
func (f ExtensionEditForm) View() string {
	// Use full width and height
	contentWidth := f.width
	contentHeight := f.height
	
	// Title with mode indicator
	modeStr := "Configuration"
	switch f.mode {
	case EditModeContext:
		modeStr = "Context File (GEMINI.md)"
	case EditModeJSON:
		modeStr = "Raw JSON"
	}
	
	title := h1Style.Render(fmt.Sprintf("✏️  Edit Extension - %s", modeStr))
	
	// Build content based on mode
	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")
	
	// Calculate available space for content
	availableHeight := contentHeight - 6 // Title, spacing, help text
	
	switch f.mode {
	case EditModeConfig:
		content.WriteString(f.renderConfigForm())
		
	case EditModeContext:
		if f.previewActive {
			content.WriteString(f.renderMarkdownPreview())
		} else {
			// Set textarea to use more space
			f.textarea.SetWidth(contentWidth - 8)
			f.textarea.SetHeight(availableHeight - 4)
			content.WriteString(f.renderTextarea())
		}
		
	case EditModeJSON:
		// Set textarea to use more space
		f.textarea.SetWidth(contentWidth - 8)
		f.textarea.SetHeight(availableHeight - 4)
		content.WriteString(f.renderTextarea())
	}
	
	// Mode switching help at bottom
	content.WriteString("\n\n")
	helpText := "Ctrl+T: Switch Mode • Ctrl+S: Save • Esc: Cancel"
	if f.mode == EditModeContext {
		helpText += " • Ctrl+P: Toggle Preview"
	}
	content.WriteString(keyDescStyle.Render(helpText))
	
	// Error display
	if f.err != nil {
		content.WriteString("\n")
		content.WriteString(errorStyle.Render(f.err.Error()))
	}
	
	// Use full screen with minimal padding
	formStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(contentWidth).
		Height(contentHeight)
	
	return formStyle.Render(content.String())
}

// renderConfigForm renders the configuration form fields
func (f ExtensionEditForm) renderConfigForm() string {
	var b strings.Builder
	
	fields := []struct {
		label string
		index int
	}{
		{"Name", extNameField},
		{"Version", extVersionField},
		{"Description", extDescriptionField},
	}
	
	for _, field := range fields {
		labelStyle := textStyle.Copy()
		if f.focusIndex == field.index {
			labelStyle = labelStyle.Foreground(colorAccent).Bold(true)
		}
		
		b.WriteString(labelStyle.Render(field.label + ":"))
		b.WriteString("\n")
		b.WriteString(f.inputs[field.index].View())
		b.WriteString("\n\n")
	}
	
	// MCP Servers section
	b.WriteString(h2Style.Render("MCP Servers"))
	b.WriteString("\n")
	
	if f.extension.MCPServers != nil && len(f.extension.MCPServers) > 0 {
		for name, server := range f.extension.MCPServers {
			serverInfo := fmt.Sprintf("• %s: %s", name, server.Command)
			b.WriteString(textDimStyle.Render(serverInfo))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(textDimStyle.Render("No MCP servers configured"))
		b.WriteString("\n")
	}
	
	b.WriteString("\n")
	b.WriteString(textDimStyle.Render("Note: Use JSON mode to edit MCP server configuration"))
	
	return b.String()
}

// renderTextarea renders the textarea for editing
func (f ExtensionEditForm) renderTextarea() string {
	return f.textarea.View()
}

// renderMarkdownPreview renders the markdown preview
func (f ExtensionEditForm) renderMarkdownPreview() string {
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


// save saves the extension changes
func (f ExtensionEditForm) save() tea.Cmd {
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