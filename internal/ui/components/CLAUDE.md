# UI Components Development Guide

This guide provides best practices and patterns for developing reusable UI components in the Gemini CLI Manager.

## Overview

The UI components package (`internal/ui/components`) contains reusable, testable UI elements that can be shared across the application. This approach promotes:

- **Code reuse** - Write once, use everywhere
- **Consistency** - Uniform UI behavior across the app
- **Testability** - Components can be tested in isolation
- **Maintainability** - Changes in one place affect all uses

## Component Design Principles

### 1. Self-Contained
Each component should be a complete, independent unit:
```go
type Component struct {
    // Internal state
    data      []Item
    selected  int
    
    // Configuration
    width     int
    height    int
    
    // Styling
    styles    ComponentStyles
}
```

### 2. Configurable
Components should expose configuration methods:
```go
// Good: Chainable configuration
tabBar := components.NewTabBar(tabs, width).
    SetStyles(activeStyle, inactiveStyle, borderColor).
    SetActiveIndex(0)

// Good: Individual setters for runtime changes
tabBar.SetWidth(newWidth)
tabBar.SetActiveByID("settings")
```

### 3. Theme-Aware
Always use theme colors, never hardcode:
```go
// ❌ BAD
style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

// ✅ GOOD
style := lipgloss.NewStyle().Foreground(theme.Primary())
```

### 4. Flexible Rendering
Provide multiple rendering options:
```go
// Basic render
output := component.Render()

// Render with additional context
output := component.RenderWithContent(content, height)

// Render specific parts
header := component.RenderHeader()
body := component.RenderBody()
```

## Creating a New Component

### 1. Define the Structure
```go
// Package declaration with documentation
package components

// ComponentName represents a reusable UI element
type ComponentName struct {
    // Public fields (if any)
    
    // Private fields
    items    []Item
    selected int
    width    int
    styles   Styles
}

// Item represents a single item in the component
type Item struct {
    ID    string
    Label string
    Icon  string
}
```

### 2. Provide Constructor
```go
// NewComponentName creates a new component instance
func NewComponentName(items []Item, width int) *ComponentName {
    return &ComponentName{
        items: items,
        width: width,
        // Set sensible defaults
        selected: 0,
        styles:   DefaultStyles(),
    }
}
```

### 3. Add Configuration Methods
```go
// SetWidth updates the component width
func (c *ComponentName) SetWidth(width int) {
    c.width = width
}

// SetStyles allows style customization
func (c *ComponentName) SetStyles(styles Styles) {
    c.styles = styles
}

// SetSelected sets the selected item by index
func (c *ComponentName) SetSelected(index int) error {
    if index < 0 || index >= len(c.items) {
        return fmt.Errorf("index out of bounds: %d", index)
    }
    c.selected = index
    return nil
}
```

### 4. Implement Render Methods
```go
// Render produces the component's visual output
func (c *ComponentName) Render() string {
    var output strings.Builder
    
    // Build the component
    for i, item := range c.items {
        style := c.styles.Normal
        if i == c.selected {
            style = c.styles.Selected
        }
        
        line := fmt.Sprintf("%s %s", item.Icon, item.Label)
        output.WriteString(style.Render(line))
        output.WriteString("\n")
    }
    
    return output.String()
}
```

## Testing Components

### 1. Unit Tests
Test component logic in isolation:
```go
func TestComponentName_SetSelected(t *testing.T) {
    items := []Item{
        {ID: "1", Label: "First"},
        {ID: "2", Label: "Second"},
    }
    comp := NewComponentName(items, 80)
    
    // Test valid selection
    err := comp.SetSelected(1)
    assert.NoError(t, err)
    assert.Equal(t, 1, comp.GetSelected())
    
    // Test invalid selection
    err = comp.SetSelected(5)
    assert.Error(t, err)
}
```

### 2. Visual Tests
Create test programs to verify appearance:
```go
// test/adhoc/test_component.go
func main() {
    theme.SetTheme("github-dark")
    
    items := generateTestItems()
    comp := components.NewComponentName(items, 80)
    
    // Test different states
    fmt.Println("Normal state:")
    fmt.Println(comp.Render())
    
    fmt.Println("\nWith selection:")
    comp.SetSelected(2)
    fmt.Println(comp.Render())
}
```

### 3. Integration Tests
Test how components work together:
```go
func TestTabBarWithList(t *testing.T) {
    tabs := []Tab{{Title: "Items", ID: "items"}}
    tabBar := NewTabBar(tabs, 80)
    
    list := NewList(items, 80)
    
    // Render together
    output := tabBar.RenderWithContent(list.Render(), 20)
    
    // Verify output
    assert.Contains(t, output, "Items")
    assert.Contains(t, output, "First item")
}
```

## Component Patterns

### 1. State Management
Keep state internal but provide access methods:
```go
type List struct {
    items    []Item
    cursor   int
    selected map[int]bool
}

// Public getters
func (l *List) GetCursor() int { return l.cursor }
func (l *List) IsSelected(index int) bool { return l.selected[index] }

// State modifiers return new state (immutable style)
func (l *List) MoveCursor(delta int) *List {
    newList := *l
    newList.cursor = clamp(l.cursor+delta, 0, len(l.items)-1)
    return &newList
}
```

### 2. Event Handling
For interactive components, define clear interfaces:
```go
// Event types
type SelectEvent struct{ Index int }
type ChangeEvent struct{ OldIndex, NewIndex int }

// Handler interface
type EventHandler interface {
    OnSelect(SelectEvent)
    OnChange(ChangeEvent)
}

// Component with handlers
func (c *Component) SetHandler(h EventHandler) {
    c.handler = h
}
```

### 3. Composition
Build complex components from simpler ones:
```go
type Form struct {
    title  *Text
    fields []*Input
    submit *Button
}

func (f *Form) Render() string {
    parts := []string{
        f.title.Render(),
        "",  // Empty line
    }
    
    for _, field := range f.fields {
        parts = append(parts, field.Render())
    }
    
    parts = append(parts, "", f.submit.Render())
    
    return strings.Join(parts, "\n")
}
```

## Best Practices

### 1. Width Calculations
Always account for borders and padding:
```go
func (c *Component) Render() string {
    // Account for borders (2) and padding (4)
    contentWidth := c.width - 6
    
    // Ensure minimum width
    if contentWidth < 10 {
        contentWidth = 10
    }
    
    // Use MaxWidth to prevent overflow
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        Width(c.width).
        MaxWidth(c.width).
        Render(content)
}
```

### 2. Error Handling
Fail gracefully with sensible defaults:
```go
func (c *Component) Render() string {
    if len(c.items) == 0 {
        return c.renderEmptyState()
    }
    
    if c.width < MinimumWidth {
        c.width = MinimumWidth
    }
    
    // Normal rendering...
}
```

### 3. Performance
Cache expensive computations:
```go
type Table struct {
    data         [][]string
    columnWidths []int
    dirty        bool
}

func (t *Table) calculateColumnWidths() {
    if !t.dirty {
        return  // Use cached values
    }
    
    // Expensive calculation...
    t.dirty = false
}
```

### 4. Documentation
Document public APIs thoroughly:
```go
// TabBar renders a horizontal bar of selectable tabs.
// 
// The active tab is visually distinguished and connects seamlessly
// with content rendered below. Tabs can have icons and IDs for
// programmatic selection.
//
// Example:
//
//	tabs := []Tab{
//	    {Title: "Home", Icon: "🏠", ID: "home"},
//	    {Title: "Settings", Icon: "⚙️", ID: "settings"},
//	}
//	tabBar := NewTabBar(tabs, 80)
//	tabBar.SetActiveByID("home")
//	output := tabBar.Render()
type TabBar struct {
    // ...
}
```

## Component Checklist

When creating a new component, ensure:

- [ ] **Constructor** with sensible defaults
- [ ] **Configuration methods** for runtime changes
- [ ] **Theme integration** using theme package
- [ ] **Width handling** with proper calculations
- [ ] **Error handling** for edge cases
- [ ] **Documentation** with examples
- [ ] **Unit tests** for logic
- [ ] **Visual tests** in test/adhoc
- [ ] **README entry** in components/README.md

## Current Components

### TabBar
- **Purpose**: Renders horizontal tabs with content area
- **Features**: Dynamic tabs, icon support, seamless borders
- **Test files**: test/adhoc/test_tabs*.go

### Card
- **Purpose**: Renders cards for extensions, profiles, and other items
- **Features**: Multiple states (normal, selected, focused, active), icon/subtitle support, metadata display, truncation, responsive width
- **Test files**: 
  - Unit tests: internal/ui/components/card_test.go
  - Visual test: Run `make card` in test/adhoc/
- **Usage example**:
```go
card := components.NewCard(60).
    SetTitle("Extension Name", "🧩").
    SetSubtitle("v1.2.0").
    SetDescription("Extension description").
    AddMetadata("MCP Servers", "2 servers", "⚡").
    SetSelected(isSelected)

output := card.Render()
```

### (Future components)
- **Modal** - Centered modal dialogs
- **FormField** - Consistent form inputs
- **StatusBar** - Three-section status bar
- **EmptyState** - "No items" displays
- **SearchBar** - Already exists, needs to be moved to components

## Contributing

When adding new components:

1. Check if similar component exists
2. Discuss design in issue/PR
3. Follow patterns in this guide
4. Add comprehensive tests
5. Update documentation
6. Consider accessibility (screen reader friendly output)

Remember: Components should be **reusable**, **testable**, and **maintainable**.