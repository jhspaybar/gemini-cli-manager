package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// FieldType represents the type of form field
type FieldType int

const (
	// TextInput is a single-line text input
	TextInput FieldType = iota
	// TextArea is a multi-line text input (future enhancement)
	TextArea
	// Checkbox is a boolean checkbox
	Checkbox
	// Select is a dropdown selection (future enhancement)
	Select
)

// FormField represents a reusable form field component
type FormField struct {
	label       string
	value       string
	helpText    string
	placeholder string
	fieldType   FieldType
	focused     bool
	required    bool
	width       int
	
	// For text input
	textInput textinput.Model
	
	// For checkbox
	checked bool
	
	// Validation
	validator func(string) error
	error     error
	
	// Styles
	labelStyle        lipgloss.Style
	focusedLabelStyle lipgloss.Style
	fieldStyle        lipgloss.Style
	focusedFieldStyle lipgloss.Style
	helpStyle         lipgloss.Style
	errorStyle        lipgloss.Style
}

// NewFormField creates a new form field component
func NewFormField(label string, fieldType FieldType) *FormField {
	// Create text input model
	ti := textinput.New()
	ti.Prompt = ""
	
	// Default styles using theme
	borderColor := theme.Border()
	focusColor := theme.BorderFocus()
	errorColor := theme.Error()
	
	return &FormField{
		label:     label,
		fieldType: fieldType,
		textInput: ti,
		width:     40, // Default width
		
		labelStyle: lipgloss.NewStyle().
			Foreground(theme.TextPrimary()),
		
		focusedLabelStyle: lipgloss.NewStyle().
			Foreground(theme.Primary()).
			Bold(true),
		
		fieldStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor).
			Padding(0, 1),
		
		focusedFieldStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(focusColor).
			Padding(0, 1),
		
		helpStyle: lipgloss.NewStyle().
			Foreground(theme.TextSecondary()),
		
		errorStyle: lipgloss.NewStyle().
			Foreground(errorColor),
	}
}

// SetValue sets the field value
func (f *FormField) SetValue(value string) *FormField {
	f.value = value
	if f.fieldType == TextInput {
		f.textInput.SetValue(value)
	}
	return f
}

// GetValue returns the current value
func (f *FormField) GetValue() string {
	if f.fieldType == TextInput {
		return f.textInput.Value()
	}
	return f.value
}

// SetPlaceholder sets the placeholder text
func (f *FormField) SetPlaceholder(placeholder string) *FormField {
	f.placeholder = placeholder
	if f.fieldType == TextInput {
		f.textInput.Placeholder = placeholder
	}
	return f
}

// SetHelpText sets the help text shown below the field
func (f *FormField) SetHelpText(helpText string) *FormField {
	f.helpText = helpText
	return f
}

// SetRequired marks the field as required
func (f *FormField) SetRequired(required bool) *FormField {
	f.required = required
	return f
}

// SetWidth sets the field width
func (f *FormField) SetWidth(width int) *FormField {
	f.width = width
	if f.fieldType == TextInput {
		f.textInput.Width = width - 2 // Account for padding
	}
	return f
}

// SetFocused sets the focus state
func (f *FormField) SetFocused(focused bool) *FormField {
	f.focused = focused
	if f.fieldType == TextInput {
		if focused {
			f.textInput.Focus()
		} else {
			f.textInput.Blur()
		}
	}
	return f
}

// SetValidator sets a validation function
func (f *FormField) SetValidator(validator func(string) error) *FormField {
	f.validator = validator
	return f
}

// Validate validates the current value
func (f *FormField) Validate() error {
	value := f.GetValue()
	
	// Check required
	if f.required && strings.TrimSpace(value) == "" {
		f.error = fmt.Errorf("This field is required")
		return f.error
	}
	
	// Run custom validator
	if f.validator != nil {
		f.error = f.validator(value)
		return f.error
	}
	
	f.error = nil
	return nil
}

// SetChecked sets the checkbox state (for checkbox fields)
func (f *FormField) SetChecked(checked bool) *FormField {
	if f.fieldType == Checkbox {
		f.checked = checked
		f.value = "false"
		if checked {
			f.value = "true"
		}
	}
	return f
}

// IsChecked returns the checkbox state
func (f *FormField) IsChecked() bool {
	return f.checked
}

// Render renders the form field
func (f *FormField) Render() string {
	var lines []string
	
	// Label with required indicator
	label := f.label
	if f.required {
		label += " *"
	}
	
	labelStyle := f.labelStyle
	if f.focused {
		labelStyle = f.focusedLabelStyle
	}
	
	lines = append(lines, labelStyle.Render(label+":"))
	
	// Field content based on type
	fieldStyle := f.fieldStyle
	if f.focused {
		fieldStyle = f.focusedFieldStyle
	}
	
	switch f.fieldType {
	case TextInput:
		// Apply width constraint to the style
		fieldStyle = fieldStyle.Width(f.width)
		lines = append(lines, fieldStyle.Render(f.textInput.View()))
		
	case Checkbox:
		checkbox := "[ ]"
		if f.checked {
			checkbox = "[✓]"
		}
		checkboxStyle := lipgloss.NewStyle()
		if f.focused {
			checkboxStyle = checkboxStyle.Foreground(theme.Primary())
		}
		lines = append(lines, checkboxStyle.Render(checkbox))
	}
	
	// Help text
	if f.helpText != "" && f.focused {
		lines = append(lines, f.helpStyle.Render(f.helpText))
	}
	
	// Error message
	if f.error != nil {
		lines = append(lines, f.errorStyle.Render("Error: "+f.error.Error()))
	}
	
	return strings.Join(lines, "\n")
}

// RenderInline renders the field in a single line (label and field side by side)
func (f *FormField) RenderInline(labelWidth int) string {
	// Label
	label := f.label
	if f.required {
		label += " *"
	}
	
	labelStyle := f.labelStyle.Width(labelWidth).Align(lipgloss.Right)
	if f.focused {
		labelStyle = f.focusedLabelStyle.Width(labelWidth).Align(lipgloss.Right)
	}
	
	// Field
	fieldStyle := f.fieldStyle
	if f.focused {
		fieldStyle = f.focusedFieldStyle
	}
	
	var fieldContent string
	switch f.fieldType {
	case TextInput:
		fieldContent = fieldStyle.Render(f.textInput.View())
	case Checkbox:
		checkbox := "[ ]"
		if f.checked {
			checkbox = "[✓]"
		}
		fieldContent = fieldStyle.Render(checkbox)
	}
	
	// Join horizontally
	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(label+":"),
		"  ", // Gap
		fieldContent,
	)
	
	// Add help text or error below if needed
	var additionalLines []string
	if f.helpText != "" && f.focused {
		padding := strings.Repeat(" ", labelWidth+2)
		additionalLines = append(additionalLines, padding+f.helpStyle.Render(f.helpText))
	}
	if f.error != nil {
		padding := strings.Repeat(" ", labelWidth+2)
		additionalLines = append(additionalLines, padding+f.errorStyle.Render("Error: "+f.error.Error()))
	}
	
	if len(additionalLines) > 0 {
		return line + "\n" + strings.Join(additionalLines, "\n")
	}
	
	return line
}

// Update handles input events (for use in Bubble Tea Update method)
func (f *FormField) Update(msg interface{}) {
	if f.fieldType == TextInput {
		f.textInput, _ = f.textInput.Update(msg)
	}
}

// Focus focuses the field
func (f *FormField) Focus() {
	f.SetFocused(true)
}

// Blur removes focus from the field
func (f *FormField) Blur() {
	f.SetFocused(false)
}