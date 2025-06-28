# Documentation Update Plan

## Overview

This document outlines specific updates needed to our CLAUDE.md documentation files to ensure they accurately reflect current best practices and the actual state of the codebase.

## Required Updates

### 1. Main CLAUDE.md Updates

#### Section: Bubble Tea TUI Framework - View Best Practices

**Current Issue**: The BAD example on lines 402-410 doesn't show the corresponding GOOD example.

**Update needed**:
```go
// ❌ BAD: Manual layout calculations
func (m Model) View() string {
    sidebarWidth := 30
    contentWidth := m.windowWidth - sidebarWidth - 3
    // Don't do manual width calculations!
    return lipgloss.JoinHorizontal(...)
}

// ✅ GOOD: Use flexbox or lipgloss features
func (m Model) View() string {
    // For flexbox setup, calculations are OK:
    fb := flexbox.New(m.windowWidth, m.windowHeight)
    
    // Or use lipgloss with MaxWidth:
    sidebar := lipgloss.NewStyle().
        MaxWidth(30).
        Render(sidebarContent)
    
    content := lipgloss.NewStyle().
        MaxWidth(m.windowWidth - 30).
        Margin(0, 1). // Use margin instead of manual spacing
        Render(mainContent)
    
    return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
}
```

#### Section: Project Structure

**Update needed**: Add SearchBar to the components list:
```
│   ├── components/      # Reusable components
│   │   ├── list.go
│   │   ├── form.go
│   │   ├── modal.go
│   │   ├── card.go
│   │   ├── status_bar.go
│   │   ├── tabs.go
│   │   ├── form_field.go
│   │   └── search_bar.go  # ADD THIS
```

#### New Section: When Manual Calculations Are Acceptable

Add after line 509:
```markdown
#### When Manual Calculations Are Acceptable

While we generally avoid manual width/height calculations, there are specific cases where they're necessary:

1. **Initial Flexbox Setup**:
   ```go
   // ✅ OK: Setting up flexbox dimensions
   contentWidth := windowWidth - (horizontalPadding * 2)
   fb := flexbox.New(contentWidth, contentHeight)
   ```

2. **Component Creation**:
   ```go
   // ✅ OK: Creating components with specific dimensions
   tabBar := components.NewTabBar(tabs, contentWidth)
   statusBar := components.NewStatusBar(contentWidth)
   ```

3. **Never for Content Adjustment**:
   ```go
   // ❌ BAD: Manually adjusting for borders/padding
   textWidth := boxWidth - 6  // Don't do this!
   
   // ✅ GOOD: Let lipgloss handle it
   style := lipgloss.NewStyle().
       Border(lipgloss.RoundedBorder()).
       Padding(1, 2).
       MaxWidth(boxWidth)
   ```
```

### 2. Components CLAUDE.md Updates

No updates needed - this file is comprehensive and accurate.

### 3. Code Migration Tasks

Based on the documentation review, these code changes should be made:

1. **Move SearchBar Component**:
   - From: `internal/cli/search.go`
   - To: `internal/ui/components/search_bar.go`
   - Update imports in all files using SearchBar

2. **Fix Hardcoded Color**:
   - File: `internal/cli/view.go:379`
   - Change: `lipgloss.Color("0")` → `theme.Background()`

3. **Refactor Search Rendering**:
   - Replace manual search box rendering in `view.go` with SearchBar component
   - Lines: 275-280, 392-399

4. **Update Empty States**:
   - Replace manual empty state boxes with Card component
   - Lines: 296-315, 416-434

## Priority Order

1. **High**: Fix hardcoded color (immediate visual impact)
2. **High**: Update CLAUDE.md with GOOD example for manual calculations
3. **Medium**: Move SearchBar to components package
4. **Medium**: Update search rendering to use component
5. **Low**: Update empty states to use Card component

## Validation Checklist

After updates:
- [ ] All examples in CLAUDE.md compile and work
- [ ] No contradictions between documentation and code
- [ ] All current components are documented
- [ ] Best practices are clearly explained with good/bad examples
- [ ] Project structure in docs matches actual structure

## Timeline

- Documentation updates: 1 hour
- Code refactoring: 2-3 hours
- Testing and validation: 1 hour

Total estimated time: 4-5 hours