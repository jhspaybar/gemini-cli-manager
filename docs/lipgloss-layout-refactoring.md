# Lipgloss Layout Refactoring Plan

## Overview

After fixing the Modal component to use lipgloss's `Margin()` instead of manual width calculations, we should refactor other components and views to follow the same pattern. This will make the code more maintainable and less error-prone.

## Current Issues

1. **Manual width calculations everywhere**: `width - 2`, `width - 4`, `width - 8` etc.
2. **Inconsistent spacing**: Different components use different calculations
3. **Border cutoff issues**: Manual calculations sometimes get it wrong
4. **Hard to maintain**: Changes to padding/borders require updating calculations

## Refactoring Strategy

### Phase 1: Update Core Components

#### Card Component (`internal/ui/components/card.go`)
**Current issues:**
- Manual width calculations: `Width(width - 2)` for borders
- Content width calculation: `contentWidth := c.width - 7`

**Refactor to:**
```go
// Instead of:
normalStyle: lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Padding(1, 2).
    Width(width - 2), // Manual calculation

// Use:
normalStyle: lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Padding(1, 2).
    MaxWidth(width), // Let lipgloss handle it
```

#### TabBar Component (`internal/ui/components/tabs.go`)
**Current issues:**
- Complex width calculations for gap filling
- Manual border joining

**Refactor approach:**
- This component may need manual calculations for precise tab alignment
- Document why manual calculations are necessary
- Consider using flexbox for more reliable layout

### Phase 2: Update View Components

#### Main View (`internal/cli/view.go`)
**Current issues:**
- `contentWidth := m.windowWidth - (horizontalPadding * 2)`
- `searchWidth := min(60, width-4)`
- Many `Width(width-4)` patterns

**Refactor to:**
- Use margins on container elements
- Use MaxWidth for content constraints
- Let lipgloss handle overflow

### Phase 3: Update Form Components

#### Search Bar (`internal/cli/search.go`)
**Current issues:**
- Manual width adjustments

#### Forms (ProfileForm, ExtensionEditForm)
**Current issues:**
- Complex width calculations for form fields
- Manual spacing calculations

**Refactor to:**
- Use container styles with proper margins
- Let form fields use MaxWidth

## Implementation Guidelines

1. **Start with leaf components** (Card, Modal) before containers
2. **Test each change visually** using the visual test suite
3. **Document exceptions** where manual calculations are truly needed
4. **Update tests** to verify proper rendering

## Benefits

1. **More reliable rendering** - Lipgloss handles edge cases better
2. **Easier maintenance** - No need to update calculations when changing styles
3. **Consistent spacing** - Margins provide uniform spacing
4. **Better terminal compatibility** - Lipgloss adapts to different terminals

## Example Refactoring

### Before (Manual Calculation):
```go
func (c *Card) Render() string {
    // Calculate content width manually
    contentWidth := c.width - 7 // 2 for border, 4 for padding, 1 for safety
    
    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        Width(c.width - 2) // Account for borders
        
    // Truncate content manually
    if len(title) > contentWidth {
        title = title[:contentWidth-3] + "..."
    }
    
    return style.Render(content)
}
```

### After (Using Lipgloss Features):
```go
func (c *Card) Render() string {
    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2).
        MaxWidth(c.width) // Let lipgloss handle sizing
    
    // Use lipgloss for text truncation
    titleStyle := lipgloss.NewStyle().
        MaxWidth(c.width - 4). // Only if absolutely needed
        Ellipsis("...")
    
    return style.Render(titleStyle.Render(title) + "\n" + content)
}
```

## Testing Strategy

1. Create visual tests comparing old vs new rendering
2. Test with various terminal sizes (60x20, 80x24, 120x40)
3. Verify no border cutoffs or overflow
4. Check that content is properly constrained

## Priority Order

1. **High Priority** (causes visible issues):
   - Card component (used extensively)
   - View rendering functions
   - Form modals

2. **Medium Priority** (works but could be cleaner):
   - TabBar component
   - Search bar
   - Status bar

3. **Low Priority** (mostly working):
   - Helper functions
   - Test utilities