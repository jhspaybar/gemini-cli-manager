# UI Components

This package contains reusable UI components for the Gemini CLI Manager.

## TabBar Component

The `TabBar` component provides a flexible, themeable tab bar that can be used throughout the application.

### Features

- Seamless connection between active tab and content area
- Customizable styles for active and inactive tabs
- Automatic gap filling to extend the tab line
- Support for icons and IDs
- Clean border connections following Bubble Tea patterns

### Usage

```go
// Define your tabs
tabs := []components.Tab{
    {Title: "Extensions", Icon: "🧩", ID: "extensions"},
    {Title: "Profiles", Icon: "👤", ID: "profiles"},
    {Title: "Settings", Icon: "🔧", ID: "settings"},
    {Title: "Help", Icon: "❓", ID: "help"},
}

// Create tab bar
tabBar := components.NewTabBar(tabs, terminalWidth)

// Set styles
tabBar.SetStyles(activeTabStyle, inactiveTabStyle, borderColor)

// Set active tab
tabBar.SetActiveIndex(0)  // or tabBar.SetActiveByID("extensions")

// Render tabs only
tabRow := tabBar.Render()

// Or render with content
content := "Your content here..."
result := tabBar.RenderWithContent(content, contentHeight)
```

### Styling

The component expects lipgloss styles with specific border configurations:

```go
// Helper function for tab borders
tabBorderWithBottom := func(left, middle, right string) lipgloss.Border {
    border := lipgloss.RoundedBorder()
    border.BottomLeft = left
    border.Bottom = middle
    border.BottomRight = right
    return border
}

// Inactive tabs: connect to horizontal line
inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")

// Active tab: no bottom border to connect with content
activeTabBorder := tabBorderWithBottom("┘", " ", "└")
```

### Testing

See `test/adhoc/test_tabs.go` for a simple example and `test/adhoc/test_tabs_switching.go` for a demonstration of tab switching.

## Card Component

The `Card` component provides a flexible, reusable card UI element for displaying items like extensions and profiles.

### Features

- Multiple states: normal, selected, focused, and active
- Optional icon and subtitle support
- Automatic text truncation for long content
- Metadata display with optional icons
- Responsive width adjustment
- Compact rendering mode for grid layouts
- Theme-aware styling

### Usage

```go
// Create a card
card := components.NewCard(60)

// Basic card with title and description
card.SetTitle("Extension Name", "🧩").
    SetDescription("A helpful description of what this extension does")

// Card with subtitle (e.g., version)
card.SetTitle("Markdown Assistant", "🧩").
    SetSubtitle("v1.2.0").
    SetDescription("Format and preview Markdown documents")

// Add metadata
card.AddMetadata("MCP Servers", "2 servers", "⚡").
    AddMetadata("Size", "1.5MB", "💾")

// Set states
card.SetSelected(true)  // Thick border
card.SetActive(true)    // Shows active indicator (●)
card.SetFocused(true)   // Double border

// Render
output := card.Render()

// Compact render (for grids)
compact := card.RenderCompact()
```

### States

- **Normal**: Standard rounded border
- **Selected**: Thick border with accent color
- **Focused**: Double border with accent color
- **Active**: Success color border with bullet indicator (●)

### Customization

```go
// Custom styles
normalStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(theme.Border()).
    Padding(1, 2)

card.SetStyles(normalStyle, selectedStyle, focusedStyle, activeStyle)

// Dynamic width adjustment
card.SetWidth(80)
```

### Testing

- Unit tests: `internal/ui/components/card_test.go`
- Visual tests: Run `make card` in `test/adhoc/`

## Modal Component

The `Modal` component provides a centered, bordered container for dialogs, forms, and alerts.

### Features

- Centered display with automatic sizing
- Multiple preset configurations (form, alert, error, success)
- Customizable title with optional icon
- Content and footer sections
- Theme-aware styling
- Responsive width with max-width constraints
- Builder pattern for easy configuration

### Usage

```go
// Basic modal
modal := components.NewModal(terminalWidth, terminalHeight).
    SetTitle("Dialog Title", "📋").
    SetContent("Your content here\n\nCan be multiple lines").
    SetFooter("Press Enter to continue")

output := modal.Render()

// Form modal with larger width
formModal := components.NewModal(width, height).
    Form().  // Preset for forms (wider, focus border)
    SetTitle("Install Extension", "📦").
    SetContent(formContent).
    SetFooter("Enter: Submit • Esc: Cancel")

// Error modal with error styling
errorModal := components.NewModal(width, height).
    Error().  // Red border and title
    SetTitle("Error", "❌").
    SetContent("Something went wrong!").
    SetFooter("Press Enter to close")

// Success modal
successModal := components.NewModal(width, height).
    Success().  // Green border and title
    SetTitle("Success", "✅").
    SetContent("Operation completed!").
    SetFooter("Press Enter to continue")

// Custom width
modal.SetWidth(50).SetMaxWidth(80)
```

### Presets

- **Form()**: 70 width, focus border color
- **Alert()**: 50 width, warning colors
- **Error()**: 50 width, error colors
- **Success()**: 50 width, success colors
- **Large()**: 80 width, 100 max width
- **Small()**: 40 width, 50 max width

### Customization

```go
// Custom styles
modal.SetTitleStyle(lipgloss.NewStyle().Bold(true).Foreground(customColor)).
    SetContentStyle(lipgloss.NewStyle().Italic(true)).
    SetFooterStyle(lipgloss.NewStyle().Foreground(theme.TextSecondary())).
    SetBorderColor(theme.Primary())
```

### Testing

- Visual tests: `go run cmd/visual-tests/main.go modal`

## FormField Component

The `FormField` component provides consistent form field rendering with labels, help text, and validation.

### Features

- Multiple field types: text input and checkbox (more coming)
- Required field indicators (*)
- Placeholder text support
- Help text displayed when focused
- Built-in and custom validation
- Error message display
- Both inline and stacked layouts
- Theme-aware styling
- Configurable width

### Usage

```go
// Basic text input
nameField := components.NewFormField("Name", components.TextInput).
    SetPlaceholder("Enter your name").
    SetRequired(true).
    SetWidth(40)

// With help text
emailField := components.NewFormField("Email", components.TextInput).
    SetPlaceholder("user@example.com").
    SetHelpText("We'll never share your email").
    SetWidth(50)

// With validation
versionField := components.NewFormField("Version", components.TextInput).
    SetValidator(func(value string) error {
        if !strings.Contains(value, ".") {
            return fmt.Errorf("must be semantic version (e.g., 1.0.0)")
        }
        return nil
    })

// Checkbox field
agreeField := components.NewFormField("I agree to terms", components.Checkbox).
    SetChecked(true)

// Render stacked (label above field)
output := field.Render()

// Render inline (label beside field)
output := field.RenderInline(15) // 15 char label width
```

### Field Types

- **TextInput**: Single-line text input with full textinput.Model features
- **Checkbox**: Boolean checkbox with check/uncheck support

### Validation

```go
// Built-in required validation
field.SetRequired(true)
err := field.Validate() // Returns error if empty

// Custom validation
field.SetValidator(func(value string) error {
    if len(value) < 3 {
        return fmt.Errorf("must be at least 3 characters")
    }
    return nil
})
```

### State Management

```go
// Focus management
field.SetFocused(true)  // or field.Focus()
field.SetFocused(false) // or field.Blur()

// Value management
field.SetValue("initial value")
value := field.GetValue()

// Checkbox state
field.SetChecked(true)
isChecked := field.IsChecked()
```

### Bubble Tea Integration

```go
// In your Update method
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Update the focused field
    m.fields[m.focusIndex].Update(msg)
    
    // Handle focus navigation
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.fields[m.focusIndex].Blur()
            m.focusIndex = (m.focusIndex + 1) % len(m.fields)
            m.fields[m.focusIndex].Focus()
        }
    }
    
    return m, nil
}
```

### Testing

- Visual tests: `go run test/adhoc/test_form_field.go`

## StatusBar Component

The `StatusBar` component provides a three-section status bar with automatic layout and theme-aware styling.

### Features

- Three-section layout: left, middle (centered), and right (right-aligned)
- Configurable width proportions for responsive design
- Built-in support for profile status, error messages, and key bindings
- Helper functions for common status bar content
- Theme-aware styling with borders
- Support for error, info, and warning message types

### Usage

```go
// Basic status bar
statusBar := components.NewStatusBar(80)
statusBar.SetLeftItems(components.ProfileStatusItems("Production", 5, 12)).
    SetKeyBindings(components.CommonKeyBindings())

output := statusBar.Render()

// Status bar with error message
errorBar := components.NewStatusBar(80)
errorBar.SetLeftItems(components.ProfileStatusItems("Development", 8, 12)).
    SetErrorMessage(components.ErrorMessage{
        Type:    components.ErrorTypeError,
        Message: "Failed to install extension",
    }).
    SetKeyBindings([]components.KeyBinding{
        {"Enter", "Retry"},
        {"Esc", "Cancel"},
    })

// Custom content sections
customBar := components.NewStatusBar(80)
customBar.SetLeftContent("🔧 Settings Mode").
    SetMiddleContent("🎨 Applying theme changes...").
    SetRightContent("Press any key to continue")

// Adjust section proportions (left, middle, right)
proportionBar := components.NewStatusBar(80)
proportionBar.SetProportions(1, 4, 1). // Give more space to middle
    SetLeftContent("Left").
    SetMiddleContent("This is a much longer middle section").
    SetRightContent("Right")
```

### Content Types

#### Status Items
```go
// Profile and extension count
statusBar.SetLeftItems(components.ProfileStatusItems("Production", 5, 12))

// Custom status items
statusBar.SetLeftItems([]components.StatusItem{
    {"🏢", "", "Company Profile"},
    {"📦", "", "15 extensions"},
    {"🔧", "", "3 tools"},
})
```

#### Error Messages
```go
// Error message
statusBar.SetErrorMessage(components.ErrorMessage{
    Type:    components.ErrorTypeError,
    Message: "Failed to install extension",
})

// Info message with details
statusBar.SetErrorMessage(components.ErrorMessage{
    Type:    components.ErrorTypeInfo,
    Message: "Extension installed successfully",
    Details: "Restart required",
})

// Warning message
statusBar.SetErrorMessage(components.ErrorMessage{
    Type:    components.ErrorTypeWarning,
    Message: "No profile selected",
})
```

#### Key Bindings
```go
// Common key bindings
statusBar.SetKeyBindings(components.CommonKeyBindings())

// Custom key bindings
statusBar.SetKeyBindings([]components.KeyBinding{
    {"Ctrl+R", "Reload"},
    {"Ctrl+S", "Save"},
    {"Esc", "Exit"},
})
```

### Configuration

```go
// Set section proportions (default: 2, 2, 3)
statusBar.SetProportions(3, 2, 2) // More space for left section

// Set width
statusBar.SetWidth(100)

// Render content only (without border)
content := statusBar.RenderContent()

// Clear all sections
statusBar.Clear()
```

### Helper Functions

```go
// Common key bindings (Tab, L, ?, q)
bindings := components.CommonKeyBindings()

// Profile status items
items := components.ProfileStatusItems("Production", 5, 12)
// Creates: ["👤 Production", "🧩 5/12"]
```

### Testing

- Visual tests: `go run cmd/visual-tests/main.go status-bar`