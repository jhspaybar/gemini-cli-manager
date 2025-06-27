# Flexbox Layout Guide for Gemini CLI Manager

This guide explains how to use the [stickers](https://github.com/76creates/stickers) flexbox library for creating responsive, maintainable layouts in our Bubble Tea TUI application.

## Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Core Concepts](#core-concepts)
4. [Basic Usage](#basic-usage)
5. [Advanced Patterns](#advanced-patterns)
6. [Best Practices](#best-practices)
7. [Migration Guide](#migration-guide)

## Introduction

The stickers flexbox library provides a CSS flexbox-inspired layout system for terminal applications. It allows us to create responsive, scalable layouts using ratios instead of fixed dimensions.

### Why Flexbox?

- **Responsive**: Automatically adjusts to terminal size changes
- **Maintainable**: Clear separation between layout and content
- **Flexible**: Easy to create complex layouts with proper alignment
- **Consistent**: Predictable behavior across different terminal sizes

## Installation

```bash
go get github.com/76creates/stickers/flexbox
```

## Core Concepts

### 1. FlexBox Types

- **FlexBox** (Vertical): Stacks rows vertically
- **HorizontalFlexBox**: Stacks columns horizontally

### 2. Container Types

- **Row**: Contains cells arranged horizontally
- **Column**: Contains cells arranged vertically

### 3. Cell

The basic building block that holds content and styling.

### 4. Ratios

Elements are sized using ratios, not fixed dimensions. A cell with ratio (2, 1) will be twice as wide as a cell with ratio (1, 1) in the same row.

## Basic Usage

### Creating a Simple Layout

```go
import (
    "github.com/76creates/stickers/flexbox"
    "github.com/charmbracelet/lipgloss"
)

func createSimpleLayout(width, height int) string {
    // Create main container
    fb := flexbox.New(width, height)
    
    // Create a row
    row := fb.NewRow()
    
    // Create cells with ratios
    leftCell := flexbox.NewCell(1, 1)  // Takes 1/3 of width
    rightCell := flexbox.NewCell(2, 1) // Takes 2/3 of width
    
    // Set content
    leftCell.SetContent("Sidebar")
    rightCell.SetContent("Main Content")
    
    // Add cells to row
    row.AddCells(leftCell, rightCell)
    
    // Add row to flexbox
    fb.AddRows([]*flexbox.Row{row})
    
    return fb.Render()
}
```

### Creating a Header-Body-Footer Layout

```go
func createAppLayout(width, height int) string {
    fb := flexbox.New(width, height)
    
    // Header (fixed height)
    headerRow := fb.NewRow()
    headerCell := flexbox.NewCell(1, 1)
    headerCell.SetContent(h1Style.Render("🚀 Gemini CLI Manager"))
    headerRow.AddCells(headerCell)
    headerRow.LockHeight(3) // Fixed 3 lines height
    
    // Body (flexible)
    bodyRow := fb.NewRow()
    bodyCell := flexbox.NewCell(1, 3) // Takes most vertical space
    bodyCell.SetContent(renderMainContent())
    bodyRow.AddCells(bodyCell)
    
    // Footer (fixed height)
    footerRow := fb.NewRow()
    footerCell := flexbox.NewCell(1, 1)
    footerCell.SetContent(renderStatusBar())
    footerRow.AddCells(footerCell)
    footerRow.LockHeight(1) // Fixed 1 line height
    
    fb.AddRows([]*flexbox.Row{headerRow, bodyRow, footerRow})
    return fb.Render()
}
```

### Creating a Two-Column Layout with Sidebar

```go
func createTwoColumnLayout(width, height int) string {
    fb := flexbox.New(width, height)
    
    mainRow := fb.NewRow()
    
    // Sidebar (25% width)
    sidebarCell := flexbox.NewCell(1, 1)
    sidebarCell.SetStyle(sidebarStyle)
    sidebarCell.SetContent(renderSidebar())
    
    // Content area (75% width)
    contentCell := flexbox.NewCell(3, 1)
    contentCell.SetStyle(contentStyle)
    contentCell.SetContent(renderContent())
    
    mainRow.AddCells(sidebarCell, contentCell)
    fb.AddRows([]*flexbox.Row{mainRow})
    
    return fb.Render()
}
```

## Advanced Patterns

### 1. Nested Flexboxes

```go
func createComplexLayout(width, height int) string {
    // Main container
    mainFb := flexbox.New(width, height)
    
    // Header row
    headerRow := mainFb.NewRow()
    headerCell := flexbox.NewCell(1, 1)
    headerCell.SetContent("Header")
    headerRow.AddCells(headerCell)
    headerRow.LockHeight(3)
    
    // Body row with nested flexbox
    bodyRow := mainFb.NewRow()
    bodyCell := flexbox.NewCell(1, 5)
    
    // Create nested flexbox for body content
    bodyFb := flexbox.NewHorizontal(0, 0) // Size will be inherited
    
    // Left panel
    leftCol := bodyFb.NewColumn()
    leftCell := flexbox.NewCell(1, 1)
    leftCell.SetContent(renderExtensionsList())
    leftCol.AddCells(leftCell)
    
    // Right panel with vertical split
    rightCol := bodyFb.NewColumn()
    rightFb := flexbox.New(0, 0)
    
    // Top section
    topRow := rightFb.NewRow()
    topCell := flexbox.NewCell(1, 2)
    topCell.SetContent(renderDetails())
    topRow.AddCells(topCell)
    
    // Bottom section
    bottomRow := rightFb.NewRow()
    bottomCell := flexbox.NewCell(1, 1)
    bottomCell.SetContent(renderActions())
    bottomRow.AddCells(bottomCell)
    
    rightFb.AddRows([]*flexbox.Row{topRow, bottomRow})
    rightCell := flexbox.NewCell(2, 1)
    rightCell.SetContent(rightFb.Render())
    rightCol.AddCells(rightCell)
    
    bodyFb.AddColumns([]*flexbox.Column{leftCol, rightCol})
    bodyCell.SetContent(bodyFb.Render())
    bodyRow.AddCells(bodyCell)
    
    mainFb.AddRows([]*flexbox.Row{headerRow, bodyRow})
    return mainFb.Render()
}
```

### 2. Dynamic Content with Scrolling

```go
func createScrollableList(items []string, width, height int) string {
    fb := flexbox.New(width, height)
    
    // Calculate visible items
    visibleItems := height - 2 // Account for borders
    startIdx := max(0, cursorPos - visibleItems/2)
    endIdx := min(len(items), startIdx + visibleItems)
    
    for i := startIdx; i < endIdx; i++ {
        row := fb.NewRow()
        cell := flexbox.NewCell(1, 1)
        
        // Highlight selected item
        content := items[i]
        if i == cursorPos {
            cell.SetStyle(selectedStyle)
            content = "> " + content
        } else {
            content = "  " + content
        }
        
        cell.SetContent(content)
        row.AddCells(cell)
        fb.AddRows([]*flexbox.Row{row})
    }
    
    return fb.Render()
}
```

### 3. Modal Dialogs

```go
func createModal(title, message string, width, height int) string {
    // Create centered modal
    modalWidth := min(60, width-10)
    modalHeight := min(20, height-4)
    
    fb := flexbox.New(modalWidth, modalHeight)
    fb.SetStyle(modalStyle)
    
    // Title row
    titleRow := fb.NewRow()
    titleCell := flexbox.NewCell(1, 1)
    titleCell.SetStyle(modalTitleStyle)
    titleCell.SetContent(title)
    titleRow.AddCells(titleCell)
    titleRow.LockHeight(2)
    
    // Content row
    contentRow := fb.NewRow()
    contentCell := flexbox.NewCell(1, 3)
    contentCell.SetContent(lipgloss.NewStyle().
        Padding(1).
        Render(message))
    contentRow.AddCells(contentCell)
    
    // Button row
    buttonRow := fb.NewRow()
    cancelCell := flexbox.NewCell(1, 1)
    cancelCell.SetContent(buttonStyle.Render("Cancel"))
    okCell := flexbox.NewCell(1, 1)
    okCell.SetContent(primaryButtonStyle.Render("OK"))
    buttonRow.AddCells(cancelCell, okCell)
    buttonRow.LockHeight(3)
    
    fb.AddRows([]*flexbox.Row{titleRow, contentRow, buttonRow})
    
    // Center in viewport
    return lipgloss.Place(
        width, height,
        lipgloss.Center, lipgloss.Center,
        fb.Render(),
    )
}
```

## When NOT to Use Flexbox

While flexbox is powerful for creating responsive layouts, there are cases where it's not appropriate:

### 1. Vertical Content Lists

**Problem**: Flexbox with ratios distributes content across the entire available height, causing unwanted spacing.

```go
// ❌ Bad: Using flexbox for vertical content lists
fb := flexbox.New(width, height)
for _, item := range items {
    row := fb.NewRow()
    cell := flexbox.NewCell(1, 1)
    cell.SetContent(item)
    row.AddCells(cell)
    fb.AddRows([]*flexbox.Row{row})
}
// This spreads items across the full height!

// ✅ Good: Simple string concatenation for natural flow
var lines []string
for _, item := range items {
    lines = append(lines, renderItem(item))
}
return strings.Join(lines, "\n")
```

### 2. Main Application Layout

**Problem**: Using flexbox for the entire app view can cause content spreading.

```go
// ❌ Bad: Flexbox for main app view with ratios
func (m Model) View() string {
    fb := flexbox.New(m.width, m.height)
    headerRow := fb.NewRow()
    contentRow := fb.NewRow()
    contentCell := flexbox.NewCell(1, 10) // 10:1 ratio causes spreading!
    // ...
}

// ✅ Good: lipgloss.JoinVertical for main layout
func (m Model) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.renderHeader(),
        m.renderContent(),
        m.renderFooter(),
    )
}
```

### 3. Dynamic Content Areas

**When to avoid flexbox:**
- Content height is unknown or variable
- You want content to take only the space it needs
- Scrollable areas where content exceeds viewport
- Simple vertical stacking of elements

**When flexbox IS appropriate:**
- Two-column or multi-column layouts
- Responsive sidebars
- Fixed-size modals and forms
- Tab bars and navigation
- Status bars with multiple sections

### 4. Rule of Thumb

Use flexbox for **structural layout** (columns, sidebars, navigation), but use simple string joining for **content flow** (lists, paragraphs, dynamic content).

## Best Practices

### 1. Use Ratios, Not Fixed Sizes

```go
// ❌ Bad: Fixed pixel values
cell.SetMinWidth(80)

// ✅ Good: Ratio-based sizing
leftCell := flexbox.NewCell(1, 1)   // 1/3 width
rightCell := flexbox.NewCell(2, 1)  // 2/3 width
```

### 2. Lock Heights/Widths Sparingly

```go
// Only lock for fixed elements like headers/footers
headerRow.LockHeight(3)  // Header is always 3 lines
// Don't lock content areas - let them be responsive
```

### 3. Separate Layout from Content

```go
// ❌ Bad: Mixing layout and content logic
func renderView() string {
    // Complex layout logic mixed with content...
}

// ✅ Good: Clear separation
func createLayout(width, height int) *flexbox.FlexBox {
    // Pure layout structure
    fb := flexbox.New(width, height)
    // ...
    return fb
}

func renderContent() string {
    // Pure content rendering
}

func (m Model) View() string {
    layout := createLayout(m.width, m.height)
    // Populate layout with content
    return layout.Render()
}
```

### 4. Reuse Layout Components

```go
// Create reusable layout functions
func createCard(title, content string) *flexbox.Row {
    row := flexbox.NewRow()
    cell := flexbox.NewCell(1, 1)
    cell.SetStyle(cardStyle)
    cell.SetContent(
        lipgloss.JoinVertical(
            lipgloss.Left,
            h2Style.Render(title),
            content,
        ),
    )
    row.AddCells(cell)
    return row
}

// Use in multiple places
fb.AddRows([]*flexbox.Row{
    createCard("Extensions", extensionList),
    createCard("Profiles", profileList),
})
```

### 5. Handle Terminal Resize

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Layout will automatically adjust on next View()
    }
    return m, nil
}
```

## Migration Guide

### Converting from Manual Layout to Flexbox

#### Before (Manual Layout):
```go
func (m Model) View() string {
    // Complex width calculations
    sidebarWidth := 30
    contentWidth := m.windowWidth - sidebarWidth - 3 // borders + padding
    
    // Manual joining with padding
    sidebar := lipgloss.NewStyle().
        Width(sidebarWidth).
        Padding(1).
        Border(lipgloss.NormalBorder()).
        Render(m.renderSidebar())
    
    content := lipgloss.NewStyle().
        Width(contentWidth).
        Padding(1).
        Render(m.renderContent())
    
    return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
}
```

#### After (Flexbox):
```go
func (m Model) View() string {
    fb := flexbox.New(m.windowWidth, m.windowHeight)
    row := fb.NewRow()
    
    // Sidebar - 30% width
    sidebarCell := flexbox.NewCell(3, 1)
    sidebarCell.SetStyle(sidebarStyle)
    sidebarCell.SetContent(m.renderSidebar())
    
    // Content - 70% width
    contentCell := flexbox.NewCell(7, 1)
    contentCell.SetStyle(contentStyle)
    contentCell.SetContent(m.renderContent())
    
    row.AddCells(sidebarCell, contentCell)
    fb.AddRows([]*flexbox.Row{row})
    
    return fb.Render()
}
```

### Key Migration Steps

1. **Replace lipgloss.JoinVertical/JoinHorizontal** with FlexBox rows/columns
2. **Replace width calculations** with ratios
3. **Replace manual padding/margins** with cell styling
4. **Use nested FlexBoxes** for complex layouts
5. **Let FlexBox handle** responsive behavior

## Common Patterns for Gemini CLI

### Extension List with Details Panel

```go
func (m Model) renderExtensionsView() string {
    fb := flexbox.New(m.windowWidth, m.windowHeight-4) // Account for tabs and status
    
    row := fb.NewRow()
    
    // Extensions list (40%)
    listCell := flexbox.NewCell(2, 1)
    listCell.SetStyle(listStyle)
    listCell.SetContent(m.renderExtensionsList())
    
    // Details panel (60%)
    detailsCell := flexbox.NewCell(3, 1)
    detailsCell.SetStyle(detailsStyle)
    if m.selectedExtension != nil {
        detailsCell.SetContent(m.renderExtensionDetails())
    } else {
        detailsCell.SetContent(emptyStateStyle.Render("Select an extension"))
    }
    
    row.AddCells(listCell, detailsCell)
    fb.AddRows([]*flexbox.Row{row})
    
    return fb.Render()
}
```

### Form Layout

```go
func createFormLayout(fields []FormField, width, height int) string {
    fb := flexbox.New(width, height)
    
    for _, field := range fields {
        row := fb.NewRow()
        
        // Label (30%)
        labelCell := flexbox.NewCell(3, 1)
        labelCell.SetContent(labelStyle.Render(field.Label))
        
        // Input (70%)
        inputCell := flexbox.NewCell(7, 1)
        inputCell.SetContent(field.Input.View())
        
        row.AddCells(labelCell, inputCell)
        row.LockHeight(3) // Standard height for form fields
        fb.AddRows([]*flexbox.Row{row})
    }
    
    return fb.Render()
}
```

## Performance Tips

1. **Pre-calculate layouts** when possible instead of recreating on every render
2. **Use style passing** to avoid setting styles on every cell
3. **Minimize nested FlexBoxes** - flatten when possible
4. **Cache rendered content** for static sections

## Troubleshooting

### Common Issues

1. **Content overflow**: Use `SetMinWidth()` or `SetMinHeight()` on cells
2. **Uneven spacing**: Check your ratios add up correctly
3. **Style not applying**: Enable style passing with `StylePassing(true)`
4. **Performance issues**: Profile your render methods, cache static content

## Summary

The stickers flexbox library provides a powerful, flexible way to create responsive TUI layouts. By following these patterns and best practices, we can create maintainable, beautiful terminal interfaces that work across different terminal sizes.

Key takeaways:
- Think in ratios, not pixels
- Separate layout from content
- Use nested FlexBoxes for complex layouts
- Let the library handle responsive behavior
- Keep it simple and reusable