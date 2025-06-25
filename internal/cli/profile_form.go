package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gemini-cli/manager/internal/profile"
)

// ProfileForm represents a form for creating/editing profiles
type ProfileForm struct {
	// Form data
	profileID   string
	isEdit      bool
	
	// Form inputs
	inputs      []textinput.Model
	focusIndex  int
	
	// Extensions selection
	availableExtensions []string
	selectedExtensions  map[string]bool
	extensionConfigs    map[string]map[string]interface{}
	extensionsCursor    int
	
	// UI state
	width       int
	height      int
	err         error
	
	// Callbacks
	onSave      func(*profile.Profile) tea.Cmd
	onCancel    func() tea.Cmd
}

// Form field indices
const (
	nameField = iota
	descriptionField
	totalFields
)

// NewProfileForm creates a new profile form
func NewProfileForm(p *profile.Profile, extensions []string, isEdit bool) ProfileForm {
	inputs := make([]textinput.Model, totalFields)
	
	// Name input
	inputs[nameField] = textinput.New()
	inputs[nameField].Placeholder = "My Development Profile"
	inputs[nameField].Focus()
	inputs[nameField].CharLimit = 50
	inputs[nameField].Width = 40
	inputs[nameField].Prompt = ""
	
	// Description input
	inputs[descriptionField] = textinput.New()
	inputs[descriptionField].Placeholder = "Profile for web development with React"
	inputs[descriptionField].CharLimit = 100
	inputs[descriptionField].Width = 40
	inputs[descriptionField].Prompt = ""
	
	// Pre-fill if editing
	if p != nil && isEdit {
		inputs[nameField].SetValue(p.Name)
		inputs[descriptionField].SetValue(p.Description)
	}
	
	// Initialize selected extensions
	selectedExts := make(map[string]bool)
	extConfigs := make(map[string]map[string]interface{})
	if p != nil {
		for _, extRef := range p.Extensions {
			selectedExts[extRef.ID] = extRef.Enabled
			if extRef.Config != nil {
				extConfigs[extRef.ID] = extRef.Config
			}
		}
	}
	
	form := ProfileForm{
		isEdit:              isEdit,
		inputs:              inputs,
		availableExtensions: extensions,
		selectedExtensions:  selectedExts,
		extensionConfigs:    extConfigs,
	}
	
	if p != nil {
		form.profileID = p.ID
	}
	
	return form
}

// Init initializes the form
func (f ProfileForm) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles form updates
func (f ProfileForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if f.onCancel != nil {
				return f, f.onCancel()
			}
			
		case "tab", "shift+tab":
			// Navigate between form sections
			if msg.String() == "tab" {
				f.focusIndex++
				if f.focusIndex > totalFields {
					f.focusIndex = 0
				}
			} else {
				f.focusIndex--
				if f.focusIndex < 0 {
					f.focusIndex = totalFields
				}
			}
			
			// Update focus
			for i := range f.inputs {
				if i == f.focusIndex {
					f.inputs[i].Focus()
				} else {
					f.inputs[i].Blur()
				}
			}
			
			return f, nil
			
		case "up", "k":
			if f.focusIndex == totalFields && f.extensionsCursor > 0 {
				f.extensionsCursor--
			}
			
		case "down", "j":
			if f.focusIndex == totalFields && f.extensionsCursor < len(f.availableExtensions)-1 {
				f.extensionsCursor++
			}
			
		case " ":
			// Toggle extension selection
			if f.focusIndex == totalFields && len(f.availableExtensions) > 0 {
				ext := f.availableExtensions[f.extensionsCursor]
				f.selectedExtensions[ext] = !f.selectedExtensions[ext]
			}
			
		case "enter":
			// Save only if we're on the last section or using Ctrl+S
			if f.focusIndex == totalFields {
				return f, f.save()
			}
			
		case "ctrl+s":
			// Save from anywhere
			return f, f.save()
		}
	}
	
	// Update text inputs
	var cmd tea.Cmd
	if f.focusIndex < totalFields {
		f.inputs[f.focusIndex], cmd = f.inputs[f.focusIndex].Update(msg)
	}
	
	return f, cmd
}

// View renders the form
func (f ProfileForm) View() string {
	// Form container
	formStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(2, 3).
		Width(60).
		MaxWidth(f.width - 4)
	
	// Title
	title := "Create New Profile"
	if f.isEdit {
		title = "Edit Profile"
	}
	titleStyle := h1Style.Copy().MarginBottom(1)
	
	// Build form content
	var b strings.Builder
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")
	
	// Name field
	b.WriteString(f.renderField("Name", nameField))
	b.WriteString("\n\n")
	
	// Description field
	b.WriteString(f.renderField("Description", descriptionField))
	b.WriteString("\n\n")
	
	// Extensions section
	b.WriteString(f.renderExtensions())
	b.WriteString("\n\n")
	
	// Help text
	helpText := []string{
		"Tab/Shift+Tab: Navigate fields",
		"Space: Toggle extension",
		"Ctrl+S: Save",
		"Esc: Cancel",
	}
	b.WriteString(helpDescStyle.Render(strings.Join(helpText, " • ")))
	
	// Error display
	if f.err != nil {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render("Error: " + f.err.Error()))
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
func (f ProfileForm) renderField(label string, index int) string {
	labelStyle := bodyStyle.Copy()
	if f.focusIndex == index {
		labelStyle = labelStyle.Foreground(colorPrimary).Bold(true)
	}
	
	fieldStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1)
		
	if f.focusIndex == index {
		fieldStyle = fieldStyle.BorderForeground(colorFocused)
	}
	
	return fmt.Sprintf("%s\n%s",
		labelStyle.Render(label+":"),
		fieldStyle.Render(f.inputs[index].View()),
	)
}

// renderExtensions renders the extensions selection
func (f ProfileForm) renderExtensions() string {
	labelStyle := bodyStyle.Copy()
	if f.focusIndex == totalFields {
		labelStyle = labelStyle.Foreground(colorPrimary).Bold(true)
	}
	
	var b strings.Builder
	b.WriteString(labelStyle.Render("Extensions:"))
	b.WriteString("\n")
	
	if len(f.availableExtensions) == 0 {
		b.WriteString(mutedStyle.Render("  No extensions available"))
		return b.String()
	}
	
	// Extension list box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1).
		Width(42)
		
	if f.focusIndex == totalFields {
		boxStyle = boxStyle.BorderForeground(colorFocused)
	}
	
	var extLines []string
	for i, ext := range f.availableExtensions {
		prefix := "  "
		if f.focusIndex == totalFields && i == f.extensionsCursor {
			prefix = "▶ "
		}
		
		checkbox := "☐"
		if f.selectedExtensions[ext] {
			checkbox = "☑"
		}
		
		line := fmt.Sprintf("%s%s %s", prefix, checkbox, ext)
		
		lineStyle := bodyStyle
		if f.focusIndex == totalFields && i == f.extensionsCursor {
			lineStyle = lineStyle.Bold(true)
		}
		
		extLines = append(extLines, lineStyle.Render(line))
	}
	
	b.WriteString(boxStyle.Render(strings.Join(extLines, "\n")))
	
	return b.String()
}

// save validates and saves the profile
func (f ProfileForm) save() tea.Cmd {
	// Validate
	name := strings.TrimSpace(f.inputs[nameField].Value())
	if name == "" {
		f.err = fmt.Errorf("profile name is required")
		return nil
	}
	
	// Validate name format (alphanumeric, hyphens, underscores)
	if !isValidProfileName(name) {
		f.err = fmt.Errorf("profile name must contain only letters, numbers, hyphens, and underscores")
		return nil
	}
	
	// Check for reserved names
	reservedNames := []string{"default", "system", "all", "none"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(name, reserved) && !f.isEdit {
			f.err = fmt.Errorf("'%s' is a reserved profile name", reserved)
			return nil
		}
	}
	
	// Build profile
	p := &profile.Profile{
		ID:          f.profileID,
		Name:        name,
		Description: strings.TrimSpace(f.inputs[descriptionField].Value()),
		Extensions:  []profile.ExtensionRef{},
	}
	
	// Generate ID if new
	if p.ID == "" {
		p.ID = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}
	
	// Add selected extensions
	for extID, selected := range f.selectedExtensions {
		if selected {
			extRef := profile.ExtensionRef{
				ID:      extID,
				Enabled: true,
			}
			if config, exists := f.extensionConfigs[extID]; exists {
				extRef.Config = config
			}
			p.Extensions = append(p.Extensions, extRef)
		}
	}
	
	// Environment and other fields would be edited separately
	p.Environment = make(map[string]string)
	p.MCPServers = make(map[string]profile.ServerConfig)
	
	// Call save callback
	if f.onSave != nil {
		return f.onSave(p)
	}
	
	return nil
}

// SetSize updates form dimensions
func (f *ProfileForm) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// SetCallbacks sets form callbacks
func (f *ProfileForm) SetCallbacks(onSave func(*profile.Profile) tea.Cmd, onCancel func() tea.Cmd) {
	f.onSave = onSave
	f.onCancel = onCancel
}

// isValidProfileName checks if a profile name is valid
func isValidProfileName(name string) bool {
	// Allow alphanumeric, hyphens, underscores, and spaces
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\- ]+$`, name)
	return matched
}