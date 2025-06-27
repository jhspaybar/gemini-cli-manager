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