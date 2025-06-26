package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
)

// ExtensionInstallForm represents a form for installing extensions
type ExtensionInstallForm struct {
	// Form data
	installingFromPath bool // true for path, false for URL
	
	// Form inputs
	inputs      []textinput.Model
	focusIndex  int
	
	// UI state
	width       int
	height      int
	err         error
	installing  bool
	progress    string
	
	// Callbacks
	onInstall   func(source string, isPath bool) tea.Cmd
	onCancel    func() tea.Cmd
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
	// Form container
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(2, 3).
		Width(70).
		MaxWidth(f.width - 4)
	
	// Title
	titleStyle := h1Style.Copy().MarginBottom(1)
	
	// Build form content
	var b strings.Builder
	b.WriteString(titleStyle.Render("Install Extension"))
	b.WriteString("\n")
	
	if f.installing {
		// Show installation progress
		b.WriteString("\n")
		b.WriteString(textStyle.Render("Installing..."))
		b.WriteString("\n\n")
		if f.progress != "" {
			b.WriteString(textMutedStyle.Render(f.progress))
			b.WriteString("\n")
		}
	} else {
		// Show form
		b.WriteString(textMutedStyle.Render("Install from a local path or remote URL"))
		b.WriteString("\n\n")
		
		// Source field
		b.WriteString(f.renderField("Source", sourceField))
		b.WriteString("\n\n")
		
		// Examples
		b.WriteString(helpStyle.Render("Examples:"))
		b.WriteString("\n")
		b.WriteString(textMutedStyle.Render("  • /Users/me/my-extension"))
		b.WriteString("\n")
		b.WriteString(textMutedStyle.Render("  • ~/Documents/extensions/my-tool"))
		b.WriteString("\n")
		b.WriteString(textMutedStyle.Render("  • https://github.com/user/gemini-extension"))
		b.WriteString("\n")
		b.WriteString(textMutedStyle.Render("  • git@github.com:user/gemini-extension.git"))
		b.WriteString("\n\n")
		
		// Help text
		helpText := []string{
			"Enter: Install",
			"Esc: Cancel",
		}
		b.WriteString(keyDescStyle.Render(strings.Join(helpText, " • ")))
		
		// Error display
		if f.err != nil {
			b.WriteString("\n\n")
			b.WriteString(errorStyle.Render("Error: " + f.err.Error()))
		}
	}
	
	// Center the form
	form := formStyle.Render(b.String())
	return lipgloss.Place(
		f.width, f.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
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