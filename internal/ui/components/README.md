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