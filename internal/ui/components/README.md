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