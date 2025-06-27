package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

func main() {
	theme.SetTheme("github-dark")

	fmt.Println("FormField Component Visual Test")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Test 1: Basic text input
	fmt.Println("Test 1: Basic Text Input")
	fmt.Println(strings.Repeat("-", 80))
	
	nameField := components.NewFormField("Name", components.TextInput).
		SetPlaceholder("Enter your name").
		SetRequired(true).
		SetWidth(40)
	
	fmt.Println("Unfocused:")
	fmt.Println(nameField.Render())
	
	fmt.Println("\nFocused:")
	nameField.SetFocused(true)
	fmt.Println(nameField.Render())
	
	// Test 2: Text input with help text
	fmt.Println("\n\nTest 2: Text Input with Help Text")
	fmt.Println(strings.Repeat("-", 80))
	
	emailField := components.NewFormField("Email", components.TextInput).
		SetPlaceholder("user@example.com").
		SetHelpText("We'll never share your email").
		SetWidth(50).
		SetFocused(true)
	
	fmt.Println(emailField.Render())
	
	// Test 3: Field with validation error
	fmt.Println("\n\nTest 3: Field with Validation Error")
	fmt.Println(strings.Repeat("-", 80))
	
	versionField := components.NewFormField("Version", components.TextInput).
		SetValue("invalid-version").
		SetValidator(func(value string) error {
			if value != "" && !strings.Contains(value, ".") {
				return fmt.Errorf("must be semantic version (e.g., 1.0.0)")
			}
			return nil
		}).
		SetWidth(30)
	
	versionField.Validate()
	fmt.Println(versionField.Render())
	
	// Test 4: Checkbox field
	fmt.Println("\n\nTest 4: Checkbox Field")
	fmt.Println(strings.Repeat("-", 80))
	
	agreeField := components.NewFormField("I agree to the terms", components.Checkbox)
	
	fmt.Println("Unchecked:")
	fmt.Println(agreeField.Render())
	
	fmt.Println("\nChecked and focused:")
	agreeField.SetChecked(true).SetFocused(true)
	fmt.Println(agreeField.Render())
	
	// Test 5: Inline rendering
	fmt.Println("\n\nTest 5: Inline Field Rendering")
	fmt.Println(strings.Repeat("-", 80))
	
	fields := []*components.FormField{
		components.NewFormField("First Name", components.TextInput).
			SetPlaceholder("John").
			SetWidth(30),
		components.NewFormField("Last Name", components.TextInput).
			SetPlaceholder("Doe").
			SetWidth(30),
		components.NewFormField("Age", components.TextInput).
			SetPlaceholder("25").
			SetWidth(10),
	}
	
	// Focus the second field
	fields[1].SetFocused(true).SetHelpText("Family name")
	
	for _, field := range fields {
		fmt.Println(field.RenderInline(15))
		fmt.Println() // Spacing between fields
	}
	
	// Test 6: Form layout simulation
	fmt.Println("\n\nTest 6: Complete Form Layout")
	fmt.Println(strings.Repeat("-", 80))
	
	formFields := []*components.FormField{
		components.NewFormField("Extension Name", components.TextInput).
			SetRequired(true).
			SetValue("My Extension").
			SetWidth(40),
		components.NewFormField("Version", components.TextInput).
			SetRequired(true).
			SetValue("1.0.0").
			SetWidth(20),
		components.NewFormField("Description", components.TextInput).
			SetPlaceholder("Brief description of your extension").
			SetWidth(60),
		components.NewFormField("Enable Auto-update", components.Checkbox).
			SetChecked(true),
	}
	
	// Focus the third field
	formFields[2].SetFocused(true)
	
	// Render with border
	formContent := []string{}
	for i, field := range formFields {
		formContent = append(formContent, field.Render())
		if i < len(formFields)-1 {
			formContent = append(formContent, "") // Spacing
		}
	}
	
	// Add a border around the form
	formBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border()).
		Padding(2, 3).
		Width(70).
		Render(strings.Join(formContent, "\n"))
	
	fmt.Println(formBox)
	
	// Test 7: Different widths
	fmt.Println("\n\nTest 7: Fields with Different Widths")
	fmt.Println(strings.Repeat("-", 80))
	
	widthTests := []int{20, 40, 60}
	for _, width := range widthTests {
		field := components.NewFormField("Test Field", components.TextInput).
			SetPlaceholder("Type here...").
			SetWidth(width)
		fmt.Printf("\nWidth %d:\n", width)
		fmt.Println(field.Render())
	}
}