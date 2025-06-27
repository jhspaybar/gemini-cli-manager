# Theming Guide for Gemini CLI Manager

This guide covers how to implement and use themes in our TUI application using `bubbletint` for consistent, maintainable, and customizable color schemes.

## Table of Contents

1. [Overview](#overview)
2. [Installing bubbletint](#installing-bubbletint)
3. [Setting Up Themes](#setting-up-themes)
4. [Using Theme Colors](#using-theme-colors)
5. [Creating Custom Themes](#creating-custom-themes)
6. [Best Practices](#best-practices)
7. [Migration Guide](#migration-guide)

## Overview

Theming is essential for creating professional, accessible, and customizable TUI applications. Using `bubbletint` provides:

- **280+ pre-built themes** from popular terminal emulators
- **Consistent color usage** across the entire application
- **Runtime theme switching** for user preferences
- **Accessibility support** with high-contrast themes
- **Easy maintenance** with centralized color definitions

## Installing bubbletint

Add bubbletint to your project:

```bash
go get -u github.com/lrstanley/bubbletint@latest
```

## Setting Up Themes

### 1. Create a Theme Package

Create `internal/theme/theme.go`:

```go
package theme

import (
    "github.com/lrstanley/bubbletint"
    "github.com/charmbracelet/lipgloss"
)

// Global theme registry
var registry *bubbletint.Registry

// Initialize with default theme
func init() {
    registry = bubbletint.NewRegistry(
        bubbletint.TintDraculaPlus,  // Default theme
        bubbletint.TintGithubDark,
        bubbletint.TintNord,
        bubbletint.TintGruvboxDark,
        bubbletint.TintSolarizedDark,
        bubbletint.TintTomorrowNight,
    )
    registry.SetTint(bubbletint.TintDraculaPlus)
}

// SetTheme changes the active theme
func SetTheme(theme string) error {
    return registry.SetTintString(theme)
}

// GetAvailableThemes returns all available theme names
func GetAvailableThemes() []string {
    return registry.TintNames()
}

// GetCurrentTheme returns the active theme name
func GetCurrentTheme() string {
    return registry.Tint().Name()
}
```

### 2. Define Semantic Color Functions

Add semantic color accessors to `theme.go`:

```go
// Core colors from theme
func Fg() lipgloss.Color           { return registry.Fg() }
func Bg() lipgloss.Color           { return registry.Bg() }
func Black() lipgloss.Color        { return registry.Black() }
func Red() lipgloss.Color          { return registry.Red() }
func Green() lipgloss.Color        { return registry.Green() }
func Yellow() lipgloss.Color       { return registry.Yellow() }
func Blue() lipgloss.Color         { return registry.Blue() }
func Magenta() lipgloss.Color      { return registry.Magenta() }
func Cyan() lipgloss.Color         { return registry.Cyan() }
func White() lipgloss.Color        { return registry.White() }
func BrightBlack() lipgloss.Color  { return registry.BrightBlack() }
func BrightRed() lipgloss.Color    { return registry.BrightRed() }
func BrightGreen() lipgloss.Color  { return registry.BrightGreen() }
func BrightYellow() lipgloss.Color { return registry.BrightYellow() }
func BrightBlue() lipgloss.Color   { return registry.BrightBlue() }
func BrightMagenta() lipgloss.Color{ return registry.BrightMagenta() }
func BrightCyan() lipgloss.Color   { return registry.BrightCyan() }
func BrightWhite() lipgloss.Color  { return registry.BrightWhite() }

// Semantic colors for our application
func Primary() lipgloss.Color      { return Cyan() }
func Secondary() lipgloss.Color    { return Blue() }
func Success() lipgloss.Color      { return Green() }
func Warning() lipgloss.Color      { return Yellow() }
func Error() lipgloss.Color        { return Red() }
func Info() lipgloss.Color         { return BrightBlue() }

// UI element colors
func TextPrimary() lipgloss.Color  { return Fg() }
func TextSecondary() lipgloss.Color{ return BrightBlack() }
func TextMuted() lipgloss.Color    { return Black() }
func Border() lipgloss.Color       { return BrightBlack() }
func BorderFocus() lipgloss.Color  { return Primary() }
func Selection() lipgloss.Color    { return BrightBlack() }
```

## Using Theme Colors

### 1. Update Styles to Use Theme Colors

Replace hardcoded colors with theme functions:

```go
// ❌ Bad: Hardcoded colors
var titleStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("87")).
    Background(lipgloss.Color("235"))

// ✅ Good: Theme colors
var titleStyle = lipgloss.NewStyle().
    Foreground(theme.Primary()).
    Background(theme.Bg())
```

### 2. Create a Styles Package

Create `internal/cli/styles.go`:

```go
package cli

import (
    "github.com/charmbracelet/lipgloss"
    "github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// UpdateStyles refreshes all styles with current theme
func UpdateStyles() {
    // Text styles
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(theme.TextPrimary())
    
    h1Style = lipgloss.NewStyle().
        Bold(true).
        Foreground(theme.TextPrimary()).
        MarginBottom(1)
    
    h2Style = lipgloss.NewStyle().
        Bold(true).
        Foreground(theme.TextPrimary())
    
    textStyle = lipgloss.NewStyle().
        Foreground(theme.TextPrimary())
    
    textDimStyle = lipgloss.NewStyle().
        Foreground(theme.TextSecondary())
    
    textMutedStyle = lipgloss.NewStyle().
        Foreground(theme.TextMuted())
    
    // Accent styles
    accentStyle = lipgloss.NewStyle().
        Foreground(theme.Primary())
    
    successStyle = lipgloss.NewStyle().
        Foreground(theme.Success())
    
    warningStyle = lipgloss.NewStyle().
        Foreground(theme.Warning())
    
    errorStyle = lipgloss.NewStyle().
        Foreground(theme.Error())
    
    // UI element styles
    selectedStyle = lipgloss.NewStyle().
        Background(theme.Selection()).
        Foreground(theme.TextPrimary())
    
    borderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.Border())
    
    focusBorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.BorderFocus())
    
    // Component styles
    cardStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.Border()).
        Padding(1).
        Background(theme.Bg())
    
    activeCardStyle = cardStyle.Copy().
        BorderForeground(theme.Primary()).
        Background(theme.Selection())
    
    buttonStyle = lipgloss.NewStyle().
        Padding(0, 3).
        Background(theme.Selection()).
        Foreground(theme.TextPrimary())
    
    primaryButtonStyle = buttonStyle.Copy().
        Background(theme.Primary()).
        Foreground(theme.Bg()).
        Bold(true)
}

// Initialize styles with default theme
func init() {
    UpdateStyles()
}
```

### 3. Dynamic Theme Switching

Add theme switching capability:

```go
// In your model
type Model struct {
    // ... other fields
    currentTheme string
}

// Handle theme change
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+t":
            // Cycle through themes
            themes := theme.GetAvailableThemes()
            currentIdx := 0
            for i, t := range themes {
                if t == m.currentTheme {
                    currentIdx = i
                    break
                }
            }
            nextIdx := (currentIdx + 1) % len(themes)
            
            // Apply new theme
            theme.SetTheme(themes[nextIdx])
            UpdateStyles() // Refresh all styles
            m.currentTheme = themes[nextIdx]
            
            return m, nil
        }
    }
}
```

## Creating Custom Themes

### 1. Define a Custom Theme

```go
// internal/theme/custom.go
package theme

import (
    "github.com/lrstanley/bubbletint"
    "image/color"
)

var GeminiTheme = bubbletint.Tint{
    ID:          "gemini",
    Name:        "Gemini",
    Bg:          color.RGBA{R: 0x1a, G: 0x1b, B: 0x26, A: 0xff},
    Fg:          color.RGBA{R: 0xd5, G: 0xd8, B: 0xda, A: 0xff},
    Black:       color.RGBA{R: 0x28, G: 0x2a, B: 0x36, A: 0xff},
    Red:         color.RGBA{R: 0xed, G: 0x86, B: 0x96, A: 0xff},
    Green:       color.RGBA{R: 0xc3, G: 0xe8, B: 0x8d, A: 0xff},
    Yellow:      color.RGBA{R: 0xfa, G: 0xc8, B: 0x63, A: 0xff},
    Blue:        color.RGBA{R: 0x82, G: 0xaa, B: 0xff, A: 0xff},
    Magenta:     color.RGBA{R: 0xc7, G: 0x92, B: 0xea, A: 0xff},
    Cyan:        color.RGBA{R: 0x89, G: 0xdd, B: 0xff, A: 0xff},
    White:       color.RGBA{R: 0xd5, G: 0xd8, B: 0xda, A: 0xff},
    BrightBlack: color.RGBA{R: 0x41, G: 0x45, B: 0x51, A: 0xff},
    // ... define all bright colors
}

// Register custom theme
func init() {
    bubbletint.Register(GeminiTheme)
}
```

### 2. Theme Configuration File

Support loading theme preferences:

```go
// internal/config/theme.go
type ThemeConfig struct {
    DefaultTheme string   `json:"default_theme"`
    CustomThemes []string `json:"custom_themes"`
}

func LoadThemeConfig() (*ThemeConfig, error) {
    // Load from ~/.gemini-cli/theme.json
}
```

## Best Practices

### 1. Always Use Theme Variables

```go
// ❌ Never hardcode colors
style.Foreground(lipgloss.Color("87"))

// ✅ Always use theme functions
style.Foreground(theme.Primary())
```

### 2. Create Semantic Color Functions

```go
// ❌ Using base colors directly
style.Foreground(theme.Red())

// ✅ Use semantic functions
style.Foreground(theme.Error())
```

### 3. Group Related Styles

```go
// Card component styles
var (
    cardStyle       lipgloss.Style
    cardTitleStyle  lipgloss.Style
    cardBodyStyle   lipgloss.Style
)

func updateCardStyles() {
    base := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.Border())
    
    cardStyle = base.Copy().
        Padding(1)
    
    cardTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(theme.Primary())
    
    cardBodyStyle = lipgloss.NewStyle().
        Foreground(theme.TextPrimary())
}
```

### 4. Consider Accessibility

```go
// Provide high contrast theme options
var accessibleThemes = []string{
    "high-contrast",
    "solarized-light",
    "github-light",
}

// Check contrast ratios
func isHighContrast(theme string) bool {
    return contains(accessibleThemes, theme)
}
```

### 5. Cache Style Calculations

```go
type StyleCache struct {
    mu     sync.RWMutex
    styles map[string]lipgloss.Style
}

func (c *StyleCache) Get(key string, builder func() lipgloss.Style) lipgloss.Style {
    c.mu.RLock()
    if style, ok := c.styles[key]; ok {
        c.mu.RUnlock()
        return style
    }
    c.mu.RUnlock()
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    style := builder()
    c.styles[key] = style
    return style
}
```

## Migration Guide

### Step 1: Install Dependencies

```bash
go get -u github.com/lrstanley/bubbletint@latest
```

### Step 2: Create Theme Package

1. Create `internal/theme/` directory
2. Copy the theme setup code from above
3. Define your semantic colors

### Step 3: Update Existing Styles

```go
// Before
var (
    colorAccent = lipgloss.Color("87")
    titleStyle = lipgloss.NewStyle().
        Foreground(colorAccent)
)

// After
import "github.com/jhspaybar/gemini-cli-manager/internal/theme"

var titleStyle = lipgloss.NewStyle().
    Foreground(theme.Primary())
```

### Step 4: Add Theme Switching

1. Add theme field to your model
2. Add keyboard shortcut for theme switching
3. Update help text to include theme controls

### Step 5: Test All Themes

```go
func TestAllThemes(t *testing.T) {
    for _, themeName := range theme.GetAvailableThemes() {
        t.Run(themeName, func(t *testing.T) {
            theme.SetTheme(themeName)
            UpdateStyles()
            // Run visual tests
        })
    }
}
```

## Theme Gallery

Here are some recommended themes for different preferences:

### Dark Themes (Default)
- **Dracula Plus** - High contrast with vibrant colors
- **Nord** - Muted, comfortable colors
- **Gruvbox Dark** - Warm, retro feel
- **Tomorrow Night** - Modern, balanced

### Light Themes
- **GitHub Light** - Clean, professional
- **Solarized Light** - Easy on the eyes
- **Tomorrow** - Minimal, clear

### High Contrast
- **High Contrast** - Maximum readability
- **Zenburn** - Low strain, high contrast

## Troubleshooting

### Colors Not Updating
```go
// Ensure you call UpdateStyles after theme change
theme.SetTheme(newTheme)
UpdateStyles() // Don't forget this!
```

### Terminal Compatibility
```go
// Check terminal color support
if !term.IsTerminal(int(os.Stdout.Fd())) {
    // Fallback to basic colors
}
```

### Performance Issues
```go
// Cache complex style calculations
var styleCache = make(map[string]lipgloss.Style)

func getCachedStyle(key string) lipgloss.Style {
    if style, ok := styleCache[key]; ok {
        return style
    }
    // Build and cache style
}
```

## Resources

- [bubbletint GitHub](https://github.com/lrstanley/bubbletint)
- [lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Terminal Color Schemes](https://windowsterminalthemes.dev/)
- [Color Contrast Checker](https://webaim.org/resources/contrastchecker/)