# Tab Design Simplification Plan

## Overview
After struggling with complex manual border manipulation, we're simplifying our tab implementation to use lipgloss's built-in capabilities more effectively.

## Key Principles
1. **Use lipgloss's border system** - No manual character placement
2. **Leverage built-in methods** like `UnsetBorderTop()` and `JoinHorizontal/Vertical`
3. **Keep state simple** - Just track which tab is active
4. **Let lipgloss handle the rendering**

## Implementation Plan

### 1. Simplified Tab Structure
```go
// Simple tab definition
type Tab struct {
    Title string
    Icon  string
    View  ViewType
}

// Tabs are just styled boxes with conditional bottom borders
```

### 2. Border Strategy
Instead of complex border manipulation, use lipgloss's custom border definitions:

```go
// Define two border types
var (
    // Inactive tab: full border
    inactiveTabBorder = lipgloss.RoundedBorder()
    
    // Active tab: no bottom border to connect with content
    activeTabBorder = lipgloss.Border{
        Top:         "─",
        Bottom:      " ",  // Space = no border
        Left:        "│",
        Right:       "│", 
        TopLeft:     "╭",
        TopRight:    "╮",
        BottomLeft:  "│",  // Continue side border
        BottomRight: "│",  // Continue side border
    }
)
```

### 3. Content Window Connection
The content area uses a standard border with the top removed:

```go
// Option A: Using UnsetBorderTop()
contentStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(theme.Border()).
    UnsetBorderTop().  // Removes top border
    Padding(1, 2)

// Option B: Using BorderTop(false) - often cleaner
contentStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(theme.Border()).
    BorderTop(false).  // Don't show top border
    Padding(1, 2)

// Option C: Selective borders from the start
contentStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder(), false, true, true, true).  // top, right, bottom, left
    BorderForeground(theme.Border()).
    Padding(1, 2)
```

### 4. Simple Rendering Logic
```go
func (m Model) renderTabs() string {
    var tabs []string
    
    for i, tab := range m.tabs {
        style := inactiveTabStyle
        if i == m.activeTab {
            style = activeTabStyle
        }
        
        content := fmt.Sprintf("%s %s", tab.Icon, tab.Title)
        tabs = append(tabs, style.Render(content))
    }
    
    // Join tabs with small gaps
    return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) View() string {
    tabs := m.renderTabs()
    content := m.renderActiveContent()
    
    // Content box with no top border
    contentBox := contentStyle.
        Width(m.width).
        Height(m.height - 3).
        Render(content)
    
    // Simple vertical join
    return lipgloss.JoinVertical(lipgloss.Left, tabs, contentBox)
}
```

### 5. Benefits of This Approach
- **No manual border drawing** - lipgloss handles all border rendering
- **Clean separation** - Tabs and content are separate components
- **Easy to understand** - Just two border styles and one `UnsetBorderTop()`
- **Maintainable** - Changes to styling don't require border character updates
- **Theme-friendly** - All colors come from the theme system

### 6. Visual Result
```
╭─────────────╮ ╭─────────────╮ ╭─────────────╮
│ Extensions  │ │  Profiles   │ │  Settings   │
│             │ └─────────────┘ └─────────────┘
│                                               │
│  Content appears here                         │
│  Connected to the active tab                  │
│                                               │
╰───────────────────────────────────────────────╯
```

### 7. Migration Steps
1. Remove all manual border character manipulation
2. Define simple active/inactive tab styles
3. Use `UnsetBorderTop()` on content area
4. Let lipgloss handle all rendering
5. Test with different terminal sizes

## Key Insight: Gap Filling
From the official Bubble Tea tabs example, here's how they handle the space after tabs:

```go
// Calculate the width of all tabs
row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

// Create a gap filler with bottom border
tabGap := lipgloss.NewStyle().
    BorderStyle(lipgloss.Border{Bottom: "─"}).
    BorderBottom(true).
    BorderForeground(theme.Border())

// Fill the remaining space
gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))

// Join tabs and gap
row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
```

This creates a clean line that extends from the last tab to the right edge.

## Alternative: Ultra-Simple Approach
If even this is too complex, consider:
- No borders on tabs, just background colors
- Use whitespace for visual separation
- Focus on functionality over aesthetics initially

## Implementation Order
1. Start with the ultra-simple approach (no borders)
2. Add basic borders to tabs and content
3. Implement the active tab connection
4. Add the gap filler for polish

## References
- [Bubble Tea Tabs Example](https://github.com/charmbracelet/bubbletea/tree/main/examples/tabs)
- [Lipgloss Border Documentation](https://github.com/charmbracelet/lipgloss#borders)
- [Lipgloss Layout Example](https://github.com/charmbracelet/lipgloss/tree/main/examples/layout)