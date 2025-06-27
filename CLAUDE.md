# CLAUDE.md - Go Development Guide for Gemini CLI

This comprehensive guide provides best practices, patterns, and conventions for developing our Gemini CLI tool using Go and Bubble Tea.

## Table of Contents
1. [Project Overview](#project-overview)
2. [Go Development Standards](#go-development-standards)
3. [Bubble Tea TUI Framework](#bubble-tea-tui-framework)
4. [Project Structure](#project-structure)
5. [Development Workflow](#development-workflow)
6. [Testing Strategy](#testing-strategy)
7. [Security Considerations](#security-considerations)
8. [Performance Guidelines](#performance-guidelines)

## Project Overview

We're building a CLI tool to manage extensions, prompts, and tools for Gemini in a first-party context. The application will use:
- **Language**: Go
- **TUI Framework**: Bubble Tea (charmbracelet/bubbletea)
- **Architecture**: Model-Update-View (MUV) pattern

## Go Development Standards

### Code Formatting and Style

**Always run before committing:**
```bash
go fmt ./...
go vet ./...
```

**Naming Conventions:**
- Package names: lowercase, single-word (e.g., `config`, `prompt`)
- Exported identifiers: PascalCase (e.g., `ExtensionManager`)
- Private identifiers: camelCase (e.g., `loadConfig`)
- Interfaces: verb + "er" suffix (e.g., `Reader`, `Validator`)

**Example:**
```go
// Package gemini provides CLI functionality for managing Gemini extensions.
package gemini

// ExtensionManager handles extension lifecycle operations.
type ExtensionManager struct {
    registry map[string]Extension
    config   *Config
}

// LoadExtension loads an extension from the specified path.
func (em *ExtensionManager) LoadExtension(path string) error {
    // Implementation
}
```

### Error Handling

**Always wrap errors with context:**
```go
func LoadPrompt(path string) (*Prompt, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading prompt file %s: %w", path, err)
    }
    
    var prompt Prompt
    if err := json.Unmarshal(data, &prompt); err != nil {
        return nil, fmt.Errorf("parsing prompt JSON: %w", err)
    }
    
    return &prompt, nil
}
```

**Define sentinel errors for known conditions:**
```go
var (
    ErrExtensionNotFound = errors.New("extension not found")
    ErrInvalidPrompt     = errors.New("invalid prompt format")
    ErrToolNotSupported  = errors.New("tool not supported")
)
```

### Go Module Management

```bash
# Initialize the project
go mod init github.com/your-org/gemini-cli

# Add dependencies
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss

# Keep dependencies clean
go mod tidy
```

## Bubble Tea TUI Framework

### 🎯 Core Principles

Bubble Tea follows **The Elm Architecture** with three core components:
1. **Model** - Your application state
2. **Update** - Handle messages and update state
3. **View** - Render the UI based on the model

**Critical Rules:**
- ⚠️ **NEVER use goroutines** within a Bubble Tea program
- ⚠️ **Use commands for ALL I/O operations**
- ⚠️ **Keep Update() and View() methods fast** - offload expensive work to commands
- ⚠️ **State is immutable** - always return a new model from Update()
- ⚠️ **ALWAYS use theme colors** - never hardcode color values like `lipgloss.Color("87")`

### 🏗️ Architecture Patterns

#### 1. Multi-View Applications

For complex applications with multiple screens, use one of these patterns:

##### Pattern A: Root Navigator Pattern (Recommended)
```go
// Root model manages view switching
type RootModel struct {
    currentView ViewType
    views       map[ViewType]tea.Model
    shared      *SharedState // Shared between views
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle global messages first
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.shared.Width = msg.Width
        m.shared.Height = msg.Height
        // Broadcast to all views
        var cmds []tea.Cmd
        for _, view := range m.views {
            _, cmd := view.Update(msg)
            cmds = append(cmds, cmd)
        }
        return m, tea.Batch(cmds...)
    
    case SwitchViewMsg:
        m.currentView = msg.View
        return m, m.views[m.currentView].Init()
    }
    
    // Route to current view
    newView, cmd := m.views[m.currentView].Update(msg)
    m.views[m.currentView] = newView
    return m, cmd
}

func (m RootModel) View() string {
    return m.views[m.currentView].View()
}
```

##### Pattern B: State Machine Pattern
```go
type AppState int

const (
    StateMenu AppState = iota
    StateExtensionList
    StateExtensionDetail
    StateProfileList
)

type Model struct {
    state    AppState
    previous AppState
    
    // Sub-models for each state
    menu      MenuModel
    extList   ExtensionListModel
    extDetail ExtensionDetailModel
    profiles  ProfileListModel
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateMenu:
        menu, cmd := m.menu.Update(msg)
        m.menu = menu
        // Handle state transitions
        if menu.Selected != "" {
            m.previous = m.state
            m.state = // ... next state based on selection
        }
        return m, cmd
    // ... other states
    }
}
```

#### 2. Component Composition

Each component should be self-contained:

```go
// Component interface
type Component interface {
    Init() tea.Cmd
    Update(tea.Msg) (Component, tea.Cmd)
    View() string
}

// Example: Reusable list component
type ListComponent struct {
    items    []string
    cursor   int
    selected map[int]bool
    focused  bool
}

func (l ListComponent) Update(msg tea.Msg) (Component, tea.Cmd) {
    if !l.focused {
        return l, nil
    }
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if l.cursor > 0 {
                l.cursor--
            }
        case "down", "j":
            if l.cursor < len(l.items)-1 {
                l.cursor++
            }
        case "enter", " ":
            l.selected[l.cursor] = !l.selected[l.cursor]
            return l, SelectedMsg{Index: l.cursor}
        }
    }
    return l, nil
}
```

#### 3. Message Passing

Define clear message types for communication:

```go
// Navigation messages
type (
    // Navigation
    BackMsg struct{}
    SwitchViewMsg struct{ View ViewType }
    
    // Data events
    ItemSelectedMsg struct{ ID string }
    ItemDeletedMsg struct{ ID string }
    DataLoadedMsg struct{ 
        Data interface{}
        Err  error
    }
    
    // UI events
    FocusChangedMsg struct{ Component string }
)

// Message routing in parent
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // First, try to handle in parent
    switch msg := msg.(type) {
    case BackMsg:
        return m.navigateBack()
    }
    
    // Then route to child
    child, cmd := m.currentView.Update(msg)
    m.currentView = child
    
    // Check if child emitted a message we need to handle
    if cmd != nil {
        // You can inspect the command's message here if needed
    }
    
    return m, cmd
}
```

### 📋 Command Best Practices

#### DO: Use Commands for All I/O
```go
// ✅ GOOD: Async file operation
func loadConfigCmd() tea.Cmd {
    return func() tea.Msg {
        data, err := os.ReadFile("config.json")
        if err != nil {
            return ErrorMsg{err}
        }
        
        var config Config
        if err := json.Unmarshal(data, &config); err != nil {
            return ErrorMsg{err}
        }
        
        return ConfigLoadedMsg{config}
    }
}

// ❌ BAD: Blocking I/O in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    data, err := os.ReadFile("config.json") // BLOCKS THE UI!
    // ...
}
```

#### DON'T: Use Commands Just for Messages
```go
// ❌ ANTI-PATTERN: Command that just returns a message
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, func() tea.Msg {
        return SomeMsg{} // Just use the message directly!
    }
}

// ✅ BETTER: Return message directly or update state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.someField = newValue
    return m, nil
}
```

#### Batch Operations
```go
// Run multiple commands in parallel
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        loadConfigCmd(),
        loadExtensionsCmd(),
        checkForUpdatesCmd(),
    )
}

// Sequential commands
func (m Model) performInstall() tea.Cmd {
    return tea.Sequence(
        downloadCmd(),
        extractCmd(),
        validateCmd(),
        installCmd(),
    )
}
```

### 🎨 View Best Practices

#### 1. Keep Views Pure
```go
// ✅ GOOD: Pure view function
func (m Model) View() string {
    if m.width < 40 {
        return m.compactView()
    }
    return m.fullView()
}

// ❌ BAD: Side effects in view
func (m Model) View() string {
    m.updateSomething() // NO! Views should not modify state
    go fetchData()      // NO! No goroutines!
    return m.render()
}
```

#### 2. ALWAYS Use Flexbox for Layouts

**IMPORTANT**: We use the [stickers](https://github.com/76creates/stickers) flexbox library for ALL layout needs. This ensures responsive, maintainable layouts.

```go
import "github.com/76creates/stickers/flexbox"

// ✅ GOOD: Using flexbox for layout
func (m Model) View() string {
    fb := flexbox.New(m.windowWidth, m.windowHeight)
    
    // Header
    headerRow := fb.NewRow()
    headerCell := flexbox.NewCell(1, 1)
    headerCell.SetContent(m.renderHeader())
    headerRow.AddCells(headerCell)
    headerRow.LockHeight(3)
    
    // Body with two columns
    bodyRow := fb.NewRow()
    sidebarCell := flexbox.NewCell(1, 1) // 1/3 width
    contentCell := flexbox.NewCell(2, 1) // 2/3 width
    sidebarCell.SetContent(m.renderSidebar())
    contentCell.SetContent(m.renderContent())
    bodyRow.AddCells(sidebarCell, contentCell)
    
    fb.AddRows([]*flexbox.Row{headerRow, bodyRow})
    return fb.Render()
}

// ❌ BAD: Manual layout calculations
func (m Model) View() string {
    sidebarWidth := 30
    contentWidth := m.windowWidth - sidebarWidth - 3
    // Don't do manual width calculations!
    return lipgloss.JoinHorizontal(...)
}
```

**See [docs/flexbox-guide.md](../docs/flexbox-guide.md) for comprehensive flexbox usage patterns.**

#### 3. Use Lipgloss with Theme Colors

**MANDATORY**: Always use theme variables for colors. Never hardcode color values.

```go
// ❌ NEVER DO THIS - Hardcoded colors
var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("87"))  // BAD!

// ✅ ALWAYS DO THIS - Theme colors
import "github.com/jhspaybar/gemini-cli-manager/internal/theme"

var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(theme.Primary())  // GOOD!

var selectedStyle = lipgloss.NewStyle().
    Background(theme.Selection()).
    Foreground(theme.TextPrimary())

func (m Model) renderItem(item string, selected bool) string {
    if selected {
        return selectedStyle.Render("> " + item)
    }
    return "  " + item
}
```

**See [docs/theming-guide.md](../docs/theming-guide.md) for comprehensive theming patterns.**

### 🐛 Debugging Techniques

#### 1. Message Logging
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Log all messages to a file
    if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
        defer f.Close()
        fmt.Fprintf(f, "[%s] %T: %+v\n", time.Now().Format("15:04:05"), msg, msg)
    }
    
    // Regular update logic...
}
```

#### 2. State Inspection
```go
// Add debug view toggle
func (m Model) View() string {
    if m.debugMode {
        return m.debugView()
    }
    return m.normalView()
}

func (m Model) debugView() string {
    return fmt.Sprintf(`
DEBUG MODE
==========
State: %v
Cursor: %d
Selected: %v
Error: %v

Press 'd' to toggle debug mode
`, m.state, m.cursor, m.selected, m.err)
}
```

### ⚡ Performance Guidelines

#### 1. Efficient Updates
```go
// ✅ GOOD: Only update what changed
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case CursorMovedMsg:
        m.cursor = msg.Position
        // Only update cursor, not entire state
        return m, nil
    }
}

// ❌ BAD: Recreating everything
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Don't reload everything on every update!
    m.items = loadAllItems()
    m.profiles = loadAllProfiles()
    return m, nil
}
```

#### 2. View Caching
```go
type Model struct {
    // Cache rendered components
    cachedHeader string
    headerDirty  bool
    
    items       []Item
    itemsDirty  bool
    cachedItems string
}

func (m Model) View() string {
    if m.headerDirty || m.cachedHeader == "" {
        m.cachedHeader = m.renderHeader()
        m.headerDirty = false
    }
    
    if m.itemsDirty || m.cachedItems == "" {
        m.cachedItems = m.renderItems()
        m.itemsDirty = false
    }
    
    return m.cachedHeader + "\n" + m.cachedItems
}
```

### 🚨 Common Pitfalls

1. **Modifying receivers**: Use pointer receivers when needed
```go
// ❌ Value receiver can't modify model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.count++ // This change is lost!
    return m, nil
}

// ✅ Return modified copy
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.count++ // This works because we return m
    return m, nil
}
```

2. **Message ordering**: Messages may not arrive in order
```go
// ❌ Assuming order
case Step1CompleteMsg:
    m.step = 2 // Step2CompleteMsg might arrive first!

// ✅ Check state
case StepCompleteMsg:
    if msg.Step == m.expectedStep {
        m.expectedStep++
    }
```

3. **Terminal state**: Always restore terminal on exit
```go
func main() {
    p := tea.NewProgram(Model{})
    
    // Restore terminal on panic
    defer func() {
        if r := recover(); r != nil {
            p.ReleaseTerminal()
            panic(r)
        }
    }()
    
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 📚 Testing Strategies

```go
// Use teatest for component testing
func TestListNavigation(t *testing.T) {
    tm := teatest.NewTestModel(t, 
        NewListModel([]string{"a", "b", "c"}),
    )
    
    // Send key events
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
    
    // Wait for result
    teatest.WaitFor(t, tm, func(b []byte) bool {
        return bytes.Contains(b, []byte("selected: c"))
    })
}
```

### 🎯 When to Use What Pattern

| Scenario | Pattern | Why |
|----------|---------|-----|
| Simple list selection | Single Model | Low complexity |
| Multiple screens/views | Root Navigator | Clean separation |
| Complex workflows | State Machine | Clear transitions |
| Reusable UI elements | Component Pattern | Modularity |
| Modal dialogs | Overlay Pattern | Temporary state |

### 📁 Recommended Project Structure

```
internal/
├── tui/
│   ├── app.go           # Root application model
│   ├── views/           # View-specific models
│   │   ├── menu.go
│   │   ├── extensions.go
│   │   └── profiles.go
│   ├── components/      # Reusable components
│   │   ├── list.go
│   │   ├── form.go
│   │   └── modal.go
│   ├── layouts/         # Flexbox layout definitions
│   │   ├── main.go      # Main app layout
│   │   ├── modal.go     # Modal layouts
│   │   └── forms.go     # Form layouts
│   ├── styles/          # Lipgloss styles
│   │   └── theme.go
│   ├── theme/           # Theme management (bubbletint)
│   │   ├── theme.go     # Theme registry and functions
│   │   └── custom.go    # Custom theme definitions
│   └── messages/        # Message definitions
│       └── types.go
```

## Project Structure

```
gemini-cli/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   ├── cli/
│   │   ├── model.go       # Main TUI model
│   │   ├── update.go      # Update logic
│   │   ├── view.go        # View rendering
│   │   ├── commands.go    # Async commands
│   │   └── keys.go        # Key bindings
│   ├── extension/
│   │   ├── extension.go   # Extension types
│   │   ├── loader.go      # Extension loading
│   │   ├── manager.go     # Extension management
│   │   └── validator.go   # Extension validation
│   ├── prompt/
│   │   ├── prompt.go      # Prompt types
│   │   ├── parser.go      # Prompt parsing
│   │   └── store.go       # Prompt storage
│   ├── tool/
│   │   ├── tool.go        # Tool types
│   │   ├── registry.go    # Tool registry
│   │   └── executor.go    # Tool execution
│   └── config/
│       ├── config.go      # Configuration types
│       └── loader.go      # Config loading
├── pkg/                   # Public packages (if needed)
├── cmd/                   # Additional commands (if needed)
└── scripts/
    ├── build.sh
    └── test.sh
```

## Development Workflow

### Critical Build Verification

**IMPORTANT: After writing significant code changes, ALWAYS run a build to ensure the project remains in a compilable state:**

```bash
# Quick build check (fastest)
go build ./...

# Full build with binary
go build -o gemini-cli-manager

# Build and run tests
go test ./...

# Complete verification (run before committing)
make verify  # or use the script below
```

Create a `Makefile` or `scripts/verify.sh`:
```bash
#!/bin/bash
# scripts/verify.sh - Complete build verification

set -e

echo "🔨 Running build verification..."

# Format check
echo "📝 Checking formatting..."
if ! go fmt ./... | grep -q .; then
    echo "✅ Format check passed"
else
    echo "❌ Format issues found. Run: go fmt ./..."
    exit 1
fi

# Vet check
echo "🔍 Running go vet..."
go vet ./...
echo "✅ Vet check passed"

# Build all packages
echo "🏗️  Building all packages..."
go build ./...
echo "✅ Build successful"

# Run tests
echo "🧪 Running tests..."
go test ./... -short
echo "✅ Tests passed"

# Build binary
echo "📦 Building binary..."
go build -o gemini-cli-manager
echo "✅ Binary built successfully"

echo "✨ All verification checks passed!"
```

### Building and Running

```bash
# Development build
go build -o gemini-cli-manager

# Production build with optimizations
go build -ldflags="-s -w" -o gemini-cli-manager

# Run with debug logging
./gemini-cli-manager --debug

# Live reload during development
# Install: go install github.com/cosmtrek/air@latest
air
```

### When to Run Build Verification

Run build verification:
- **After implementing new features** - Before moving to the next task
- **After refactoring** - Ensure nothing broke
- **Before committing** - Keep the main branch green
- **After resolving merge conflicts** - Verify integration
- **When switching between branches** - Ensure clean state

### Debugging Bubble Tea Apps

```go
// Enable debug logging
func main() {
    if len(os.Args) > 1 && os.Args[1] == "--debug" {
        f, err := tea.LogToFile("debug.log", "debug")
        if err != nil {
            fmt.Println("fatal:", err)
            os.Exit(1)
        }
        defer f.Close()
    }
    
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
```

## Testing Strategy

### Unit Tests with Table-Driven Approach

```go
func TestExtensionValidation(t *testing.T) {
    tests := []struct {
        name      string
        extension Extension
        wantErr   bool
        errMsg    string
    }{
        {
            name: "valid extension",
            extension: Extension{
                Name:    "test-ext",
                Version: "1.0.0",
                Type:    "prompt",
            },
            wantErr: false,
        },
        {
            name: "missing name",
            extension: Extension{
                Version: "1.0.0",
                Type:    "prompt",
            },
            wantErr: true,
            errMsg:  "extension name is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateExtension(tt.extension)
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error, got nil")
                } else if !strings.Contains(err.Error(), tt.errMsg) {
                    t.Errorf("expected error containing %q, got %q", 
                        tt.errMsg, err.Error())
                }
            } else if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### Testing Bubble Tea Components

```go
import "github.com/charmbracelet/x/exp/teatest"

func TestMainMenu(t *testing.T) {
    tm := teatest.NewTestModel(
        t, 
        initialModel(),
        teatest.WithInitialTermSize(80, 24),
    )
    
    // Navigate to extensions
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
    
    // Wait for view change
    teatest.WaitFor(
        t, tm,
        func(b []byte) bool {
            return bytes.Contains(b, []byte("Extensions"))
        },
        teatest.WithDuration(3 * time.Second),
    )
    
    // Verify final state
    finalModel := tm.FinalModel(t).(model)
    if finalModel.currentView != viewExtensions {
        t.Errorf("expected viewExtensions, got %v", finalModel.currentView)
    }
}
```

## Security Considerations

### Input Validation

```go
func validateExtensionPath(path string) error {
    // Prevent directory traversal
    cleanPath := filepath.Clean(path)
    if strings.Contains(cleanPath, "..") {
        return errors.New("invalid path: directory traversal detected")
    }
    
    // Ensure path is within allowed directory
    absPath, err := filepath.Abs(cleanPath)
    if err != nil {
        return fmt.Errorf("resolving path: %w", err)
    }
    
    if !strings.HasPrefix(absPath, allowedBaseDir) {
        return errors.New("path outside allowed directory")
    }
    
    return nil
}
```

### Safe File Operations

```go
func loadExtensionSafely(path string) (*Extension, error) {
    if err := validateExtensionPath(path); err != nil {
        return nil, err
    }
    
    // Limit file size to prevent DoS
    info, err := os.Stat(path)
    if err != nil {
        return nil, fmt.Errorf("stat file: %w", err)
    }
    
    if info.Size() > maxExtensionSize {
        return nil, errors.New("extension file too large")
    }
    
    // Read with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return loadExtensionWithContext(ctx, path)
}
```

## Performance Guidelines

### Efficient TUI Rendering

```go
// Cache complex computations
type model struct {
    items         []Item
    filteredItems []Item // Cache filtered results
    filterQuery   string
    needsRefilter bool
}

func (m *model) updateFilter(query string) {
    if m.filterQuery != query {
        m.filterQuery = query
        m.needsRefilter = true
    }
}

func (m *model) getFilteredItems() []Item {
    if m.needsRefilter {
        m.filteredItems = filterItems(m.items, m.filterQuery)
        m.needsRefilter = false
    }
    return m.filteredItems
}
```

### Memory Management

```go
// Use sync.Pool for temporary objects
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func renderExtension(ext Extension) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    // Use buffer for rendering
    buf.WriteString(ext.Name)
    buf.WriteString(" v")
    buf.WriteString(ext.Version)
    
    return buf.String()
}
```

## Key Commands and Shortcuts

Define consistent keyboard shortcuts across the application:

```go
var defaultKeyMap = keyMap{
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "help"),
    ),
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "down"),
    ),
    Select: key.NewBinding(
        key.WithKeys("enter", "space"),
        key.WithHelp("enter", "select"),
    ),
    Back: key.NewBinding(
        key.WithKeys("esc", "backspace"),
        key.WithHelp("esc", "back"),
    ),
}
```

## Build Commands

### Quick Reference
```bash
# Most common during development
go build ./...              # Quick syntax check
go test ./... -short        # Fast test run
./scripts/verify.sh         # Full verification

# Before pushing code
go fmt ./...                # Format code
go mod tidy                 # Clean dependencies
./scripts/verify.sh         # Full verification
```

### Makefile for Convenience
Create a `Makefile` in the project root:
```makefile
.PHONY: build test verify fmt clean run

# Default target
all: verify

# Quick build check
build:
	@echo "🏗️  Building..."
	@go build ./...

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./... -v

# Run short tests (faster)
test-short:
	@echo "🧪 Running short tests..."
	@go test ./... -short

# Format code
fmt:
	@echo "📝 Formatting code..."
	@go fmt ./...

# Full verification
verify: fmt
	@echo "🔨 Running full verification..."
	@./scripts/verify.sh

# Build and run
run: build
	@./gemini-cli-manager

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f gemini-cli-manager
	@go clean ./...

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
```

## Additional Considerations

### 1. Configuration Management
- Use environment variables for configuration
- Support config files in standard locations
- Validate all configuration on startup

### 2. Logging and Monitoring
- Use structured logging (e.g., zerolog)
- Log important operations for debugging
- Consider telemetry for usage patterns

### 3. Distribution
- Build for multiple platforms
- Consider using goreleaser for releases
- Provide installation scripts

### 4. Documentation
- Keep code comments up to date
- Document public APIs thoroughly
- Maintain user-facing documentation

### 5. Continuous Integration
```yaml
# .github/workflows/ci.yml example
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - run: make verify
```

This guide should be updated as the project evolves and new patterns emerge.