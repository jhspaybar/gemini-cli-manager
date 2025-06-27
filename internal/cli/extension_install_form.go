package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

// ExtensionInstallForm represents a form for installing extensions
type ExtensionInstallForm struct {
	// Form data
	installingFromPath bool // true for path, false for URL

	// Form inputs
	inputs     []textinput.Model
	focusIndex int

	// UI state
	width      int
	height     int
	err        error
	installing bool
	progress   string

	// Callbacks
	onInstall func(source string, isPath bool) tea.Cmd
	onCancel  func() tea.Cmd
}

// Form field indices
const (
	sourceField = iota
	totalInstallFields
)

// NewExtensionInstallForm creates a new extension installation form
func NewExtensionInstallForm() ExtensionInstallForm {
	inputs := make([]textinput.Model, totalInstallFields)

	// Source input (path or URL)
	inputs[sourceField] = textinput.New()
	inputs[sourceField].Placeholder = "Enter path or URL to extension (e.g., /path/to/extension or https://github.com/...)"
	inputs[sourceField].Focus()
	inputs[sourceField].CharLimit = 200
	inputs[sourceField].Width = 50
	inputs[sourceField].Prompt = ""

	return ExtensionInstallForm{
		inputs:     inputs,
		focusIndex: 0,
	}
}

// Init initializes the form
func (f ExtensionInstallForm) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles form updates
func (f ExtensionInstallForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if f.installing {
		// Handle installation progress messages
		switch msg := msg.(type) {
		case installProgressMsg:
			f.progress = msg.message
			return f, nil
		case installCompleteMsg:
			f.installing = false
			if msg.err != nil {
				f.err = msg.err
				return f, nil
			}
			// Success - close form
			if f.onCancel != nil {
				return f, f.onCancel()
			}
		}
		return f, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if f.onCancel != nil {
				return f, f.onCancel()
			}
			return f, nil

		case "enter":
			// Start installation
			cmd := f.install()
			return f, cmd

		case "ctrl+v":
			// Paste support would be nice but requires clipboard access
			return f, nil

		default:
			// Let text input handle other keys
			var cmd tea.Cmd
			f.inputs[f.focusIndex], cmd = f.inputs[f.focusIndex].Update(msg)
			return f, cmd
		}
	}

	// For non-key messages, pass to text input
	var cmd tea.Cmd
	f.inputs[f.focusIndex], cmd = f.inputs[f.focusIndex].Update(msg)

	return f, cmd
}

// View renders the form
func (f ExtensionInstallForm) View() string {
	// Build form content
	var content strings.Builder

	if f.installing {
		// Show installation progress
		content.WriteString(textStyle.Render("Installing..."))
		content.WriteString("\n\n")
		if f.progress != "" {
			content.WriteString(textMutedStyle.Render(f.progress))
		}
	} else {
		// Show form
		content.WriteString(textMutedStyle.Render("Install from a local path or remote URL"))
		content.WriteString("\n\n")

		// Source field
		content.WriteString(f.renderField("Source", sourceField))
		content.WriteString("\n\n")

		// Examples
		content.WriteString(helpStyle.Render("Examples:"))
		content.WriteString("\n")
		content.WriteString(textMutedStyle.Render("  • /Users/me/my-extension"))
		content.WriteString("\n")
		content.WriteString(textMutedStyle.Render("  • ~/Documents/extensions/my-tool"))
		content.WriteString("\n")
		content.WriteString(textMutedStyle.Render("  • https://github.com/user/gemini-extension"))
		content.WriteString("\n")
		content.WriteString(textMutedStyle.Render("  • git@github.com:user/gemini-extension.git"))

		// Error display
		if f.err != nil {
			content.WriteString("\n\n")
			content.WriteString(errorStyle.Render("Error: " + f.err.Error()))
		}
	}

	// Use the Modal component
	modal := components.NewModal(f.width, f.height).
		Form().
		SetTitle("Install Extension", "📦").
		SetContent(content.String()).
		SetFooter("Enter: Install • Esc: Cancel")
		
	// If there's an error, use error styling
	if f.err != nil {
		modal = modal.Error()
	}

	return modal.Render()
}

// renderField renders a form field
func (f ExtensionInstallForm) renderField(label string, index int) string {
	labelStyle := textStyle.Copy()
	if f.focusIndex == index {
		labelStyle = labelStyle.Foreground(colorAccent).Bold(true)
	}

	fieldStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1)

	if f.focusIndex == index {
		fieldStyle = fieldStyle.BorderForeground(colorBorderFocus)
	}

	return fmt.Sprintf("%s\n%s",
		labelStyle.Render(label+":"),
		fieldStyle.Render(f.inputs[index].View()),
	)
}

// install validates and starts the installation
func (f *ExtensionInstallForm) install() tea.Cmd {
	source := strings.TrimSpace(f.inputs[sourceField].Value())
	if source == "" {
		f.err = fmt.Errorf("source path or URL is required")
		return nil
	}

	// Basic validation for URLs
	if strings.HasPrefix(source, "http://") {
		f.err = fmt.Errorf("HTTP URLs are not allowed for security reasons, please use HTTPS")
		return nil
	}

	// Validate GitHub URLs
	if strings.Contains(source, "github.com") {
		if !strings.HasPrefix(source, "https://github.com/") && !strings.HasPrefix(source, "git@github.com:") {
			f.err = fmt.Errorf("invalid GitHub URL format")
			return nil
		}
	}

	// Determine if it's a path or URL
	isPath := !strings.HasPrefix(source, "http://") &&
		!strings.HasPrefix(source, "https://") &&
		!strings.HasPrefix(source, "git@")

	// Expand home directory if needed
	if isPath && strings.HasPrefix(source, "~/") {
		// This would need proper home directory expansion
		// For now, we'll pass it as-is and let the installer handle it
	}

	f.installing = true
	f.err = nil
	f.progress = "Starting installation..."

	// Call install callback
	if f.onInstall != nil {
		return f.onInstall(source, isPath)
	}

	return nil
}

// SetSize updates form dimensions
func (f *ExtensionInstallForm) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// SetCallbacks sets form callbacks
func (f *ExtensionInstallForm) SetCallbacks(onInstall func(source string, isPath bool) tea.Cmd, onCancel func() tea.Cmd) {
	f.onInstall = onInstall
	f.onCancel = onCancel
}

// installCompleteMsg is used for installation completion
type installCompleteMsg struct {
	extension *extension.Extension
	err       error
}
