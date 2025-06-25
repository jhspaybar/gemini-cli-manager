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

### Core Architecture

```go
type model struct {
    // View state
    currentView  viewType
    windowWidth  int
    windowHeight int
    
    // Business logic
    extensions []Extension
    prompts    []Prompt
    tools      []Tool
    
    // UI components
    list       list.Model
    textInput  textinput.Model
    viewport   viewport.Model
    help       help.Model
    
    // State
    loading    bool
    err        error
}

func (m model) Init() tea.Cmd {
    return tea.Batch(
        m.loadInitialData(),
        textinput.Blink,
    )
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        return m.handleResize(msg)
    case dataLoadedMsg:
        return m.handleDataLoaded(msg)
    }
    return m, nil
}

func (m model) View() string {
    if m.loading {
        return m.loadingView()
    }
    
    switch m.currentView {
    case viewExtensions:
        return m.extensionsView()
    case viewPrompts:
        return m.promptsView()
    case viewTools:
        return m.toolsView()
    default:
        return m.mainMenuView()
    }
}
```

### State Management Pattern

```go
type viewType int

const (
    viewMainMenu viewType = iota
    viewExtensions
    viewPrompts
    viewTools
    viewDetails
    viewEdit
)

func (m model) switchView(view viewType) (model, tea.Cmd) {
    m.currentView = view
    
    // Initialize view-specific state
    switch view {
    case viewExtensions:
        return m, m.loadExtensions()
    case viewPrompts:
        return m, m.loadPrompts()
    case viewTools:
        return m, m.loadTools()
    }
    
    return m, nil
}
```

### Command Patterns for Async Operations

```go
// Message types
type extensionsLoadedMsg struct {
    extensions []Extension
    err        error
}

// Command functions
func loadExtensionsCmd(path string) tea.Cmd {
    return func() tea.Msg {
        extensions, err := loadExtensionsFromDisk(path)
        return extensionsLoadedMsg{
            extensions: extensions,
            err:        err,
        }
    }
}

// Handle in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case extensionsLoadedMsg:
        if msg.err != nil {
            m.err = msg.err
            return m, nil
        }
        m.extensions = msg.extensions
        m.loading = false
        return m, nil
    }
    // ... other cases
}
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