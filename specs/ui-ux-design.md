# UI/UX Design Specification

## Design Philosophy

The Gemini CLI Manager follows these core principles:
- **Keyboard-First**: Every action accessible via keyboard
- **Progressive Disclosure**: Show complexity only when needed
- **Contextual Help**: Guidance available at every step
- **Visual Clarity**: Clear hierarchy and focus indicators
- **Responsive Feedback**: Immediate visual response to actions

## Visual Design System

### Color Palette

```go
type ColorScheme struct {
    // Base colors
    Background      Color // #0d1117 (GitHub dark)
    Surface         Color // #161b22
    Border          Color // #30363d
    
    // Text colors
    TextPrimary     Color // #c9d1d9
    TextSecondary   Color // #8b949e
    TextMuted       Color // #484f58
    
    // Accent colors
    Primary         Color // #58a6ff (blue)
    Success         Color // #3fb950 (green)
    Warning         Color // #d29922 (yellow)
    Error           Color // #f85149 (red)
    Info            Color // #58a6ff (blue)
    
    // State colors
    Selected        Color // #1f6feb
    Hover           Color // #30363d
    Disabled        Color // #21262d
}
```

### Typography

```go
type Typography struct {
    // Headers
    H1 lipgloss.Style // Bold, Primary color
    H2 lipgloss.Style // Bold, Secondary color
    H3 lipgloss.Style // Regular, Primary color
    
    // Body
    Body      lipgloss.Style // Regular, Primary color
    BodySmall lipgloss.Style // Regular, Secondary color
    
    // Special
    Code      lipgloss.Style // Monospace, Muted background
    Keyboard  lipgloss.Style // Monospace, Border, Padding
    Label     lipgloss.Style // Small, Muted color
}
```

## Layout Architecture

### Main Layout Structure

```
┌─────────────────────────────────────────────────────────────┐
│ Header Bar                                                  │
├─────────────────────────────────────────────────────────────┤
│ Navigation │ Content Area                                   │
│            │                                                │
│ Extensions │ ┌───────────────────────────────────────────┐ │
│ Profiles   │ │ Detail View                               │ │
│ Settings   │ │                                           │ │
│ Help       │ │ Current content based on navigation      │ │
│            │ │                                           │ │
│            │ └───────────────────────────────────────────┘ │
├────────────┴────────────────────────────────────────────────┤
│ Status Bar                                                  │
└─────────────────────────────────────────────────────────────┘
```

### Responsive Behavior

```go
type LayoutManager struct {
    minWidth  int // 80 columns
    minHeight int // 24 rows
}

func (lm *LayoutManager) CalculateLayout(width, height int) Layout {
    if width < 100 {
        // Compact mode: Hide sidebar, use full width
        return CompactLayout{
            ShowSidebar: false,
            ContentWidth: width - 2,
        }
    }
    
    // Normal mode: Sidebar + Content
    return NormalLayout{
        SidebarWidth: 20,
        ContentWidth: width - 22,
    }
}
```

## Component Library

### List Component

```go
type ListComponent struct {
    items       []ListItem
    selected    int
    viewport    viewport.Model
    filter      string
    showIcons   bool
    multiSelect bool
}

// Visual representation
/*
┌─ Extensions (12) ──────────┐
│ ▶ 󰊢 typescript-tools    ✓ │
│   󰊢 react-helper        ✓ │
│   󰊢 database-tools      ✗ │
│ ▶ 󰊢 testing-suite       ✓ │
│   󰊢 linting-config      ✓ │
└────────────────────────────┘
*/
```

### Form Component

```go
type FormComponent struct {
    fields      []FormField
    activeField int
    validation  map[string]ValidationResult
}

// Visual representation
/*
┌─ Edit Profile ─────────────────────┐
│ Name: [Web Development         ]   │
│       ─────────────────────────    │
│                                    │
│ Description:                       │
│ [Full-stack web development    ]   │
│ [environment with React and    ]   │
│ [Node.js                       ]   │
│                                    │
│ Color: [●] #61DAFB                 │
│                                    │
│ [Save] [Cancel]                    │
└────────────────────────────────────┘
*/
```

### Modal Dialog

```go
type ModalDialog struct {
    title       string
    content     tea.Model
    buttons     []Button
    width       int
    height      int
}

// Visual representation
/*
╭─ Confirm Delete ───────────────────╮
│                                    │
│ Are you sure you want to delete    │
│ the "Web Development" profile?     │
│                                    │
│ This action cannot be undone.      │
│                                    │
│        [Delete] [Cancel]           │
╰────────────────────────────────────╯
*/
```

## Navigation Patterns

### Focus Management

```go
type FocusManager struct {
    zones       []FocusZone
    activeZone  int
    activeItem  map[int]int
}

type FocusZone struct {
    ID          string
    Items       []Focusable
    Navigation  NavigationType
}

const (
    NavVertical NavigationType = iota
    NavHorizontal
    NavGrid
)

// Focus flow
/*
Tab/Shift+Tab: Move between zones
Arrow keys: Navigate within zone
Enter: Activate focused item
Esc: Exit current zone/cancel
*/
```

### Keyboard Shortcuts

```go
var globalKeyMap = KeyMap{
    // Navigation
    NextZone:     key.NewBinding(key.WithKeys("tab")),
    PrevZone:     key.NewBinding(key.WithKeys("shift+tab")),
    
    // Commands
    CommandPalette: key.NewBinding(key.WithKeys("ctrl+k")),
    QuickSwitch:    key.NewBinding(key.WithKeys("ctrl+p")),
    Search:         key.NewBinding(key.WithKeys("/")),
    
    // Actions
    New:    key.NewBinding(key.WithKeys("n")),
    Edit:   key.NewBinding(key.WithKeys("e")),
    Delete: key.NewBinding(key.WithKeys("d")),
    Toggle: key.NewBinding(key.WithKeys("space")),
    
    // Application
    Quit: key.NewBinding(key.WithKeys("ctrl+q", "ctrl+c")),
    Help: key.NewBinding(key.WithKeys("?")),
}
```

## Screen Designs

### Main Menu

```
╭─ Gemini CLI Manager ─────────────────────────────────────────╮
│                                                              │
│  Welcome to Gemini CLI Manager                               │
│                                                              │
│  Current Profile: Web Development                            │
│  Extensions: 5 active, 2 available                           │
│                                                              │
│  Quick Actions:                                              │
│                                                              │
│  → Extensions     Manage extensions and MCP servers          │
│    Profiles       Switch or create profiles                  │
│    Launch         Start Gemini CLI                          │
│    Settings       Configure application                      │
│                                                              │
│  Press ? for help · ↑↓ to navigate · Enter to select        │
│                                                              │
╰──────────────────────────────────────────────────────────────╯
```

### Extensions View

```
╭─ Extensions ─────────────────────────────────────────────────╮
│ Extensions │ typescript-tools                    [Enabled] ✓ │
│ ────────── ├─────────────────────────────────────────────────┤
│ Profiles   │ TypeScript Language Tools                       │
│ Settings   │                                                 │
│ Help       │ Provides comprehensive TypeScript support       │
│            │ including IntelliSense, refactoring, and       │
│ Search: /  │ advanced type checking.                         │
│            │                                                 │
│ [n] New    │ MCP Servers:                                    │
│ [i] Import │ • typescript-lsp (active)                       │
│ [r] Reload │   - Port: 6009                                  │
│            │   - Memory: 125MB                               │
│            │                                                 │
│ ▶ Active   │ Commands:                                       │
│   ✓ ts-tools  • ts:check-types                              │
│   ✓ react     • ts:refactor                                 │
│   ✗ vue       • ts:organize-imports                         │
│   ✓ testing                                                 │
│            │ [Space] Toggle  [e] Edit  [d] Delete  [→] Logs │
╰────────────┴─────────────────────────────────────────────────╯
```

### Profile Switcher

```
╭─ Switch Profile ─────────────────────────────────────────────╮
│                                                              │
│  Select a profile:                      [esc] Cancel         │
│                                                              │
│  ▶ 🌐 Web Development              Used 2 hours ago         │
│    🐍 Data Science                 Used yesterday           │
│    🚀 Go Backend                   Used 3 days ago          │
│    📱 Mobile Development           Used last week           │
│    ➕ Create New Profile...                                 │
│                                                              │
│  Recently used:                                              │
│  [1] Web Development  [2] Data Science  [3] Go Backend      │
│                                                              │
╰──────────────────────────────────────────────────────────────╯
```

### Settings View

```
╭─ Settings ───────────────────────────────────────────────────╮
│ General    │ Theme              [Dark ▼]                     │
│ ────────   │                                                 │
│ Extensions │ Auto-update        [✓] Enabled                  │
│ Profiles   │                                                 │
│ Advanced   │ Launch on startup  [✗] Disabled                 │
│            │                                                 │
│            │ Default Profile    [Web Development ▼]          │
│            │                                                 │
│            │ Gemini CLI Path                                 │
│            │ [/usr/local/bin/gemini                     ]    │
│            │                                                 │
│            │ Extension Directory                             │
│            │ [~/.gemini/extensions                      ]    │
│            │                                                 │
│            │                           [Apply] [Reset]        │
╰────────────┴─────────────────────────────────────────────────╯
```

## Interaction Patterns

### Loading States

```go
type LoadingIndicator struct {
    spinner  spinner.Model
    message  string
    progress *ProgressBar
}

// Visual representations
/*
Simple spinner:  ⠋ Loading extensions...
Progress bar:    [████████░░░░░░░] 53% Installing typescript-tools
Indeterminate:   [≈≈≈≈≈≈≈≈≈≈≈≈≈≈≈] Connecting to MCP server...
*/
```

### Error States

```go
type ErrorDisplay struct {
    level   ErrorLevel
    title   string
    message string
    actions []Action
}

// Visual representation
/*
╭─ ⚠️  Warning ────────────────────────────────────────────────╮
│                                                              │
│  Extension 'typescript-tools' requires update                │
│                                                              │
│  A newer version (2.1.0) is available with important        │
│  security fixes.                                             │
│                                                              │
│  [Update Now] [Remind Later] [View Changes]                 │
│                                                              │
╰──────────────────────────────────────────────────────────────╯
*/
```

### Success Feedback

```go
type SuccessNotification struct {
    message  string
    duration time.Duration
}

// Visual representation
/*
✓ Profile 'Web Development' activated successfully
*/
```

## Animation and Transitions

### Smooth Transitions

```go
type Transition struct {
    duration time.Duration
    easing   EasingFunction
}

var transitions = map[string]Transition{
    "focus":    {50 * time.Millisecond, EaseInOut},
    "expand":   {150 * time.Millisecond, EaseOut},
    "collapse": {100 * time.Millisecond, EaseIn},
    "fade":     {200 * time.Millisecond, Linear},
}
```

### Loading Animations

```go
var spinnerStyles = []string{
    "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏", // Dots
    "◐◓◑◒",           // Circle
    "▁▂▃▄▅▆▇█▇▆▅▄▃▂", // Bar
}
```

## Accessibility Features

### Screen Reader Support

```go
type Accessible struct {
    Role        string
    Label       string
    Description string
    State       map[string]bool
}

func (a Accessible) AriaLabel() string {
    // Generate appropriate ARIA-like labels
}
```

### High Contrast Mode

```go
type HighContrastTheme struct {
    ColorScheme
}

func (hct *HighContrastTheme) Apply() {
    hct.Background = Black
    hct.TextPrimary = White
    hct.Border = White
    hct.Selected = Yellow
}
```

## Help System

### Contextual Help

```go
type HelpSystem struct {
    context  string
    shortcuts []Shortcut
    tips     []string
}

// Visual representation
/*
╭─ Help ───────────────────────────────────────────────────────╮
│                                                              │
│  Extension Management                                        │
│                                                              │
│  Keyboard Shortcuts:                                         │
│  Space     Toggle extension on/off                           │
│  e         Edit extension settings                           │
│  d         Delete extension (with confirmation)              │
│  n         Add new extension                                 │
│  /         Search extensions                                 │
│  Tab       Switch to profiles                                │
│                                                              │
│  Tips:                                                       │
│  • Disabled extensions remain installed but inactive         │
│  • Use profiles to group related extensions                  │
│  • MCP servers start automatically when enabled              │
│                                                              │
│                                    [Close] Press ? or Esc   │
╰──────────────────────────────────────────────────────────────╯
*/
```

## Performance Considerations

### Rendering Optimization

```go
type RenderOptimizer struct {
    cache       map[string]string
    dirtyRegions []Region
}

func (ro *RenderOptimizer) ShouldRender(component Component) bool {
    // Only render if component has changed
    newHash := component.Hash()
    oldHash, exists := ro.cache[component.ID()]
    
    if exists && newHash == oldHash {
        return false
    }
    
    ro.cache[component.ID()] = newHash
    return true
}
```

### Viewport Management

```go
type ViewportManager struct {
    content    []string
    viewport   viewport.Model
    bufferSize int
}

func (vm *ViewportManager) OptimizeForSize(width, height int) {
    // Adjust buffer size based on terminal dimensions
    if height > 50 {
        vm.bufferSize = 1000
    } else {
        vm.bufferSize = 500
    }
}
```

This comprehensive UI/UX design specification provides a consistent, intuitive, and efficient interface for the Gemini CLI Manager.