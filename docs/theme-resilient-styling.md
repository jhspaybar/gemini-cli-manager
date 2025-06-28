# Theme-Resilient Styling Guidelines

## Current Status

The Gemini CLI Manager uses `bubbletint` for theme-aware colors, which is good. However, there are areas where we can improve to ensure the app looks great regardless of terminal theme.

## Issues Observed

1. **Tab Visibility**: Active tabs are not sufficiently distinct from inactive tabs in some themes
2. **Contrast Issues**: Some color combinations may have poor contrast in certain themes
3. **Focus Indicators**: Focus states need to be more obvious

## Recommendations

### 1. Use Semantic Color Differences

Instead of relying only on color, use multiple visual cues:

```go
// Current approach (relies heavily on color)
activeTabStyle = lipgloss.NewStyle().
    BorderForeground(colorAccent).
    Foreground(colorAccent).
    Bold(true)

// Better approach (multiple visual cues)
activeTabStyle = lipgloss.NewStyle().
    Border(activeTabBorder, true).
    BorderForeground(colorAccent).
    Background(colorSelection).  // Add background
    Foreground(colorText).       // Use high contrast text
    Bold(true).
    Padding(0, 3)               // Slightly wider padding
```

### 2. Implement Adaptive Contrast

Create a contrast checker that adjusts colors based on the theme:

```go
// Add to theme package
func EnsureContrast(fg, bg lipgloss.TerminalColor) lipgloss.TerminalColor {
    // If contrast is too low, return an alternative
    // This would need actual implementation
    return fg
}
```

### 3. Visual Hierarchy Improvements

#### For Tabs
- **Active Tab**: Background color + bold + thicker border
- **Inactive Tab**: No background + normal weight + thin border
- **Hover/Focus**: Underline or different border style

#### For Lists
- **Selected Item**: Inverse colors (swap fg/bg)
- **Cursor Position**: Prefix with "▶" or similar indicator
- **Disabled Items**: Use TextMuted color

### 4. Theme Testing Strategy

Create visual tests for each supported theme:

```bash
# Test with different themes
for theme in "Solarized Light" "Solarized Dark" "Dracula" "Nord"; do
    vhs test-theme.tape --set-theme "$theme"
done
```

### 5. Specific Style Updates Needed

#### Tab Styles (`styles.go`)
```go
// Make active tabs more distinct
activeTabStyle = lipgloss.NewStyle().
    Border(activeTabBorder, true).
    BorderForeground(colorAccent).
    Background(colorSelection).     // Add background
    Foreground(colorText).          // Ensure readable text
    Bold(true).
    Padding(0, 3)                   // More padding

// Consider dimming inactive tabs more
inactiveTabStyle = lipgloss.NewStyle().
    Border(inactiveTabBorder, true).
    BorderForeground(colorBorderDim).
    Foreground(colorTextMuted).     // More muted
    Padding(0, 2)
```

#### Selection States
```go
// Better selection visibility
menuItemSelectedStyle = lipgloss.NewStyle().
    Foreground(colorBg).            // Inverse colors
    Background(colorAccent).        // Strong background
    PaddingLeft(1).
    PaddingRight(1).
    Bold(true)                      // Add bold
```

### 6. Add Visual Indicators

Beyond color, use Unicode characters for state:
- Active: `●`, `▶`, `◆`
- Inactive: `○`, `▷`, `◇`
- Selected: `✓`, `✗`
- Loading: `⣾⣽⣻⢿⡿⣟⣯⣷`

### 7. Testing Checklist

For each theme, verify:
- [ ] Active tab is clearly distinguishable
- [ ] Selected items are obvious
- [ ] Focus states are visible
- [ ] Error messages stand out
- [ ] All text is readable
- [ ] Borders are visible

## Implementation Priority

1. **High**: Fix tab visibility (active vs inactive)
2. **High**: Improve selection states in lists
3. **Medium**: Add non-color indicators
4. **Low**: Implement contrast checking

## Example Implementation

Here's how to update the tab rendering to be more theme-resilient:

```go
// In updateTabStyles()
func updateTabStyles() {
    // Inactive - subtle and recessed
    inactiveTabStyle = lipgloss.NewStyle().
        Border(inactiveTabBorder, true).
        BorderForeground(colorBorderDim).
        Foreground(colorTextMuted).
        Padding(0, 2)

    // Active - prominent with multiple indicators
    activeTabStyle = lipgloss.NewStyle().
        Border(activeTabBorder, true).
        BorderForeground(colorAccent).
        Background(colorSelection).      // Background differentiator
        Foreground(colorText).           // High contrast text
        Bold(true).                      // Weight differentiator
        Padding(0, 3)                    // Size differentiator

    // Focused but inactive - subtle highlight
    focusedInactiveTabStyle = lipgloss.NewStyle().
        Border(inactiveTabBorder, true).
        BorderForeground(colorBorder).
        Foreground(colorText).
        Underline(true).                 // Focus indicator
        Padding(0, 2)
}
```

## Conclusion

By implementing these changes, the app will:
1. Look good in any terminal theme
2. Be more accessible
3. Provide better visual feedback
4. Maintain a professional appearance

The key is to not rely solely on color for conveying information, but to use multiple visual dimensions: color, weight, size, borders, backgrounds, and symbols.