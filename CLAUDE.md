# CLAUDE.md - Rust/Ratatui Development Guide for Gemini CLI Manager

This comprehensive guide provides patterns, best practices, and conventions for developing the Gemini CLI Manager using Rust and Ratatui.

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture Guidelines](#architecture-guidelines)
3. [Ratatui Patterns](#ratatui-patterns)
4. [Component Development](#component-development)
5. [State Management](#state-management)
6. [Event Handling](#event-handling)
7. [Layout Best Practices](#layout-best-practices)
8. [Styling and Theming](#styling-and-theming)
9. [Testing Strategy](#testing-strategy)
10. [Performance Guidelines](#performance-guidelines)
11. [Common Pitfalls](#common-pitfalls)

## Project Overview

The Gemini CLI Manager is a Terminal User Interface (TUI) application for managing Gemini extensions, profiles, and tools. It uses:
- **Language**: Rust
- **TUI Framework**: Ratatui (0.29.0)
- **Backend**: Crossterm
- **Async Runtime**: Tokio
- **Architecture**: Component-based with message passing

### Gemini Extension Format

#### MCP Servers Configuration

Extensions can configure connections to Model-Context Protocol (MCP) servers for discovering and using custom tools. The `mcp_servers` field in an extension configures these connections.

**Key Points:**
- Gemini CLI attempts to connect to each configured MCP server to discover available tools
- If multiple MCP servers expose a tool with the same name, tool names are prefixed with the server alias (e.g., `serverAlias__actualToolName`)
- The system might strip certain schema properties from MCP tool definitions for compatibility
- Note: In our storage model, we use `McpServerConfig` which differs slightly from the JSON format (e.g., no URL field since MCP servers are command-based)

**MCP Server Properties:**
- `command` (string, required): The command to execute to start the MCP server
- `args` (array of strings, optional): Arguments to pass to the command
- `env` (object, optional): Environment variables to set for the server process (values can use `$VAR_NAME` syntax)
- `cwd` (string, optional): The working directory in which to start the server
- `timeout` (number, optional): Timeout in milliseconds for requests to this MCP server
- `trust` (boolean, optional): Trust this server and bypass all tool call confirmations

**Example:**
```json
"mcp_servers": {
  "myPythonServer": {
    "command": "python",
    "args": ["mcp_server.py", "--port", "8080"],
    "cwd": "./mcp_tools/python",
    "timeout": 5000
  },
  "myNodeServer": {
    "command": "node",
    "args": ["mcp_server.js"],
    "cwd": "./mcp_tools/node"
  },
  "myDockerServer": {
    "command": "docker",
    "args": ["run", "-i", "--rm", "-e", "API_KEY", "ghcr.io/foo/bar"],
    "env": {
      "API_KEY": "$MY_API_TOKEN"
    }
  }
}
```

### Context Files (Hierarchical Instructional Context)

Extensions can include context files (typically named after the extension, e.g., `GITHUB.md` for github-tools) that provide instructions, guidelines, or context for the Gemini model. This powerful feature allows extensions to give specific instructions, coding style guides, or relevant background information to the AI.

**Key Points:**
- Context files are Markdown files containing instructions for the Gemini model
- The system manages instructional context hierarchically
- Multiple context files are concatenated with separators indicating their origin
- The CLI displays the count of loaded context files in the footer

**Context File Properties in Extensions:**
- `context_file_name` (string, optional): The name of the context file (e.g., "GITHUB.md")
- `context_content` (string, optional): The actual content of the context file

**Example Context File Content:**
```markdown
# GitHub Tools Extension

## General Instructions:
- When interacting with GitHub APIs, always check rate limits
- Prefer GraphQL API over REST API when possible for efficiency
- Include error handling for common GitHub API errors

## Authentication:
- This extension uses the GITHUB_TOKEN environment variable
- Ensure the token has appropriate permissions for requested operations

## Available Tools:
- `github_create_issue`: Creates a new issue in a repository
- `github_list_prs`: Lists pull requests with various filters
- `github_review_pr`: Reviews and comments on pull requests

## Best Practices:
- Always validate repository ownership before destructive operations
- Cache responses when appropriate to reduce API calls
- Use pagination for large result sets
```

### Current Architecture

```
src/
├── main.rs         # Entry point, initializes Tokio runtime
├── app.rs          # Main application orchestrator
├── tui.rs          # Terminal I/O and event loop management
├── components.rs   # Component trait and registry
├── action.rs       # Action types for message passing
├── config.rs       # Configuration management
├── cli.rs          # Command-line interface
├── errors.rs       # Error handling setup
└── components/     # Individual UI components
    ├── home.rs     # Home screen component
    └── fps.rs      # FPS counter component
```

### Dependency Management

**Always use `cargo add` to add new dependencies:**

```bash
# Add a new dependency
cargo add ratatui

# Add with specific features
cargo add tokio --features full

# Add a specific version
cargo add tui-input@0.10.1  # Compatible with ratatui 0.29.0

# Add as a dev dependency
cargo add --dev pretty_assertions
```

**Why use `cargo add`:**
- Automatically selects the latest compatible version
- Properly formats the Cargo.toml entry
- Handles feature flags correctly
- Updates the lock file automatically

**Common dependencies for this project:**
```bash
# TUI and rendering
cargo add ratatui
cargo add crossterm
cargo add tui-input  # For text input widgets

# Async runtime
cargo add tokio --features full

# Error handling
cargo add color-eyre
cargo add thiserror

# Serialization
cargo add serde --features derive
cargo add serde_json

# File system
cargo add dirs
cargo add tempfile --dev
```

## Architecture Guidelines

### Component-Based Design

Every UI element should be a component implementing the `Component` trait:

```rust
pub trait Component {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()>;
    fn register_config_handler(&mut self, config: Config) -> Result<()>;
    fn update(&mut self, action: Action) -> Result<Option<Action>>;
    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()>;
    fn handle_events(&mut self, event: Option<Event>) -> Result<Option<Action>>;
}
```

### Message Passing Pattern

Components communicate through `Action` messages via Tokio channels:

```rust
// Define actions in action.rs
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Action {
    Tick,
    Render,
    Resize(u16, u16),
    Navigate(Route),
    ExtensionSelected(String),
    ProfileSwitch(String),
    // Add new actions as needed
}
```

### Async Event Loop

The application uses separate loops for:
- **Ticks**: Logic updates at configurable rate (default: 4 Hz)
- **Renders**: UI updates at configurable FPS (default: 60 FPS)
- **Events**: Terminal events (keyboard, mouse, resize)

## Ratatui Patterns

### 1. Immediate Mode Rendering

Ratatui uses immediate mode rendering - the entire UI is redrawn each frame:

```rust
fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
    // Clear and redraw everything
    let chunks = Layout::default()
        .direction(Direction::Vertical)
        .constraints([
            Constraint::Length(3),    // Header
            Constraint::Min(0),       // Content
            Constraint::Length(3),    // Footer
        ])
        .split(area);

    // Render each section
    self.draw_header(frame, chunks[0])?;
    self.draw_content(frame, chunks[1])?;
    self.draw_footer(frame, chunks[2])?;
    
    Ok(())
}
```

### 2. Widget Composition

Build complex UIs by composing simple widgets:

```rust
// Use Block for borders and titles
let block = Block::default()
    .borders(Borders::ALL)
    .title("Extensions")
    .title_style(Style::default().fg(Color::Cyan));

// Wrap content in blocks
let paragraph = Paragraph::new(content)
    .block(block)
    .wrap(Wrap { trim: true });

frame.render_widget(paragraph, area);
```

### 3. Stateful Widgets

For widgets with state (like lists), maintain state in the component:

```rust
pub struct ExtensionList {
    items: Vec<Extension>,
    state: ListState,
    selected: Option<usize>,
}

impl ExtensionList {
    fn next(&mut self) {
        let i = match self.state.selected() {
            Some(i) => (i + 1) % self.items.len(),
            None => 0,
        };
        self.state.select(Some(i));
    }
}
```

### 4. Third-Party Widgets

Leverage existing widgets for common UI patterns:

#### Text Input with tui-input

For text input fields, use `tui-input` instead of implementing from scratch:

```rust
use tui_input::Input;
use tui_input::backend::crossterm::EventHandler;

pub struct SearchableList {
    search_input: Input,
    search_mode: bool,
    items: Vec<String>,
}

impl SearchableList {
    fn handle_search_input(&mut self, key: KeyEvent) -> bool {
        // Let tui-input handle the event
        if self.search_input.handle_event(&Event::Key(key)).is_some() {
            self.update_filter();
            true
        } else {
            false
        }
    }
    
    fn draw_search(&self, frame: &mut Frame, area: Rect) {
        let input_widget = Paragraph::new(self.search_input.value())
            .style(Style::default())
            .block(Block::default()
                .borders(Borders::ALL)
                .title("Search"));
        
        frame.render_widget(input_widget, area);
        
        // Set cursor position
        if self.search_mode {
            let cursor_pos = self.search_input.visual_cursor();
            frame.set_cursor_position((
                area.x + cursor_pos as u16 + 1, // +1 for border
                area.y + 1
            ));
        }
    }
}
```

**Benefits of using tui-input:**
- Proper cursor management
- Text selection support
- Clipboard operations (copy/paste)
- Unicode support
- Input validation capabilities

## Component Development

### Creating a New Component

1. Create a new file in `src/components/`:

```rust
// src/components/extension_list.rs
use crate::{action::Action, components::Component};
use color_eyre::Result;
use ratatui::prelude::*;
use tokio::sync::mpsc::UnboundedSender;

#[derive(Default)]
pub struct ExtensionList {
    action_tx: Option<UnboundedSender<Action>>,
    items: Vec<String>,
    selected: usize,
}

impl Component for ExtensionList {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.action_tx = Some(tx);
        Ok(())
    }

    fn update(&mut self, action: Action) -> Result<Option<Action>> {
        match action {
            Action::NavigateDown => {
                self.selected = (self.selected + 1) % self.items.len();
                Ok(Some(Action::Render))
            }
            _ => Ok(None),
        }
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Implement rendering logic
        Ok(())
    }
}
```

2. Register the component in `app.rs`:

```rust
self.components.insert(
    "extension_list".to_string(),
    Box::new(ExtensionList::default()),
);
```

### Component Guidelines

1. **Single Responsibility**: Each component should handle one UI concern
2. **Self-Contained State**: Components manage their own state
3. **Action-Based Updates**: State changes only through actions
4. **Error Propagation**: Always propagate errors with `?`

## State Management

### Application State

Global application state belongs in `App`:

```rust
pub struct App {
    pub running: bool,
    pub current_screen: Screen,
    pub selected_profile: Option<Profile>,
    pub extensions: Vec<Extension>,
    // ...
}
```

### Component State

Component-specific state stays within the component:

```rust
pub struct ProfileManager {
    profiles: Vec<Profile>,
    selected_index: usize,
    edit_mode: bool,
    form_state: FormState,
}
```

### State Updates

State changes follow this flow:
1. User input → Event
2. Event → Action
3. Action → State update
4. State update → Render

## Event Handling

### Keyboard Events

Handle keyboard input in components:

```rust
fn handle_events(&mut self, event: Option<Event>) -> Result<Option<Action>> {
    match event {
        Some(Event::Key(key)) => match key.code {
            KeyCode::Char('q') => Ok(Some(Action::Quit)),
            KeyCode::Up => Ok(Some(Action::NavigateUp)),
            KeyCode::Down => Ok(Some(Action::NavigateDown)),
            KeyCode::Enter => Ok(Some(Action::Select)),
            _ => Ok(None),
        },
        Some(Event::Mouse(mouse)) => {
            // Handle mouse events if needed
            Ok(None)
        },
        _ => Ok(None),
    }
}
```

### Custom Keybindings

Use the configuration system for customizable keybindings:

```rust
// In config.json5
{
    "keybindings": {
        "Home": {
            "<q>": "Quit",
            "<Ctrl-c>": "Quit",
            "<Up>": "NavigateUp",
            "<Down>": "NavigateDown",
            "<Enter>": "Select"
        }
    }
}
```

## Layout Best Practices

### Responsive Layouts

Use constraints that adapt to terminal size:

```rust
let chunks = Layout::default()
    .direction(Direction::Horizontal)
    .constraints([
        Constraint::Percentage(30),  // Sidebar
        Constraint::Min(50),         // Main content
    ])
    .split(area);
```

### Nested Layouts

Build complex layouts by nesting:

```rust
// Main vertical split
let main_chunks = Layout::vertical([
    Constraint::Length(3),   // Header
    Constraint::Min(0),      // Body
    Constraint::Length(3),   // Footer
]).split(area);

// Split body horizontally
let body_chunks = Layout::horizontal([
    Constraint::Percentage(25),  // Sidebar
    Constraint::Min(0),          // Content
]).split(main_chunks[1]);
```

### Layout Guidelines

1. **Use `Min(0)` for flexible areas** that should take remaining space
2. **Use `Length(n)` for fixed-size** elements like headers/footers
3. **Use `Percentage(n)` for proportional** layouts
4. **Avoid hardcoded dimensions** - always consider different terminal sizes

## Styling and Theming

### Consistent Styling

Define styles in a central location:

```rust
// src/theme.rs
pub struct Theme {
    pub primary: Color,
    pub secondary: Color,
    pub accent: Color,
    pub error: Color,
    pub success: Color,
    pub warning: Color,
    pub text: Color,
    pub text_dim: Color,
    pub background: Color,
}

impl Default for Theme {
    fn default() -> Self {
        Self {
            primary: Color::Cyan,
            secondary: Color::Magenta,
            accent: Color::Yellow,
            error: Color::Red,
            success: Color::Green,
            warning: Color::Yellow,
            text: Color::White,
            text_dim: Color::Gray,
            background: Color::Black,
        }
    }
}
```

### Style Application

Apply styles consistently:

```rust
let title_style = Style::default()
    .fg(theme.primary)
    .add_modifier(Modifier::BOLD);

let selected_style = Style::default()
    .bg(theme.primary)
    .fg(theme.background)
    .add_modifier(Modifier::BOLD);
```

### Styling Rules

1. **Use semantic colors** (primary, error) not raw colors
2. **Support light/dark themes** through configuration
3. **Ensure readable contrast** between fg/bg colors
4. **Use modifiers sparingly** (BOLD, ITALIC, UNDERLINED)

## Testing Strategy

### Unit Tests

Test component logic separately from rendering:

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_navigation() {
        let mut list = ExtensionList::new(vec!["a", "b", "c"]);
        assert_eq!(list.selected, 0);
        
        list.next();
        assert_eq!(list.selected, 1);
        
        list.next();
        list.next();
        assert_eq!(list.selected, 0); // Wraps around
    }
}
```

### Integration Tests

Test component interactions:

```rust
#[tokio::test]
async fn test_extension_selection() {
    let (tx, mut rx) = mpsc::unbounded_channel();
    let mut app = App::new();
    
    app.handle_action(Action::ExtensionSelected("test".into()));
    
    let action = rx.recv().await.unwrap();
    assert_eq!(action, Action::Navigate(Route::ExtensionDetail));
}
```

### TUI Testing

For testing the actual TUI output, use the `TestBackend`:

```rust
use ratatui::backend::TestBackend;
use ratatui::Terminal;

#[test]
fn test_render_output() {
    let backend = TestBackend::new(80, 24);
    let mut terminal = Terminal::new(backend).unwrap();
    
    terminal.draw(|frame| {
        app.draw(frame, frame.area()).unwrap();
    }).unwrap();
    
    let buffer = terminal.backend().buffer();
    assert!(buffer.content().contains("Extensions"));
}
```

## Performance Guidelines

### Efficient Rendering

1. **Minimize allocations** in the draw loop:
```rust
// Bad: Allocates every frame
fn draw(&mut self, frame: &mut Frame, area: Rect) {
    let items: Vec<ListItem> = self.items.iter()
        .map(|i| ListItem::new(i.clone()))
        .collect();
}

// Good: Reuse allocations
fn draw(&mut self, frame: &mut Frame, area: Rect) {
    // Cache items in component state
    if self.items_changed {
        self.cached_items = self.items.iter()
            .map(|i| ListItem::new(i.as_str()))
            .collect();
        self.items_changed = false;
    }
}
```

2. **Avoid complex calculations** in draw:
```rust
// Calculate once in update(), not in draw()
fn update(&mut self, action: Action) -> Result<Option<Action>> {
    if self.data_changed {
        self.processed_data = self.calculate_expensive_thing();
        self.data_changed = false;
    }
    Ok(None)
}
```

### Memory Management

1. **Use `&str` over `String`** where possible
2. **Prefer `Vec` capacity hints** for known sizes
3. **Clear rather than reallocate** collections

## Common Pitfalls

### 1. Blocking the Event Loop

**Wrong:**
```rust
fn update(&mut self, action: Action) -> Result<Option<Action>> {
    // This blocks the UI!
    let data = std::fs::read_to_string("large_file.txt")?;
    Ok(None)
}
```

**Right:**
```rust
fn update(&mut self, action: Action) -> Result<Option<Action>> {
    // Spawn async task
    if let Some(tx) = &self.action_tx {
        let tx = tx.clone();
        tokio::spawn(async move {
            let data = tokio::fs::read_to_string("large_file.txt").await?;
            tx.send(Action::DataLoaded(data))?;
        });
    }
    Ok(None)
}
```

### 2. State Mutation Outside Update

**Wrong:**
```rust
fn draw(&mut self, frame: &mut Frame, area: Rect) {
    // Don't mutate state in draw!
    self.counter += 1;
}
```

**Right:**
```rust
fn update(&mut self, action: Action) -> Result<Option<Action>> {
    match action {
        Action::Tick => {
            self.counter += 1;
            Ok(Some(Action::Render))
        }
        _ => Ok(None),
    }
}
```

### 3. Forgetting Error Context

**Wrong:**
```rust
let file = std::fs::read_to_string(path)?;
```

**Right:**
```rust
let file = std::fs::read_to_string(&path)
    .wrap_err_with(|| format!("Failed to read file: {}", path))?;
```

### 4. Hardcoded Dimensions

**Wrong:**
```rust
let chunks = Layout::default()
    .constraints([Constraint::Length(80)]) // Assumes 80-char terminal!
    .split(area);
```

**Right:**
```rust
let chunks = Layout::default()
    .constraints([Constraint::Percentage(100)])
    .split(area);
```

## Development Workflow

### Adding New Features

1. **Define Actions** in `action.rs`
2. **Create Component** in `src/components/`
3. **Register Component** in `app.rs`
4. **Add Routes** if needed
5. **Update Config** for new keybindings
6. **Write Tests**
7. **Update Documentation**

### Code Organization

```
src/
├── components/       # UI components
│   ├── mod.rs       # Component exports
│   ├── extension/   # Extension-related components
│   ├── profile/     # Profile-related components
│   └── common/      # Shared components
├── models/          # Data structures
├── services/        # Business logic
├── utils/           # Helper functions
└── widgets/         # Custom Ratatui widgets
```

### Git Workflow

1. Create feature branch from `main`
2. Implement feature with tests
3. Run `cargo fmt` and `cargo clippy`
4. Update CHANGELOG.md
5. Create PR with description

## Quick Reference

### Essential Imports

```rust
use color_eyre::Result;
use ratatui::{
    prelude::*,
    widgets::{Block, Borders, List, ListItem, Paragraph},
};
use tokio::sync::mpsc::UnboundedSender;
```

### Common Patterns

```rust
// Create a bordered block
let block = Block::default()
    .borders(Borders::ALL)
    .title("Title");

// Create a list with selection
let items: Vec<ListItem> = items.iter().map(|i| ListItem::new(i.as_str())).collect();
let list = List::new(items)
    .block(block)
    .highlight_style(Style::default().add_modifier(Modifier::REVERSED));

// Center content
let area = centered_rect(50, 50, area); // 50% width, 50% height

// Split area
let [header, content, footer] = Layout::vertical([
    Constraint::Length(3),
    Constraint::Min(0),
    Constraint::Length(3),
]).areas(area);
```

### Test Coverage Commands

**Always run coverage with the same output directory for consistency:**

```bash
# Generate HTML coverage report
cargo tarpaulin --out Html --output-dir coverage-report

# View the report
open coverage-report/tarpaulin-report.html
```

### Debugging Tips

1. Use `tracing` for logging:
```rust
tracing::debug!("Action received: {:?}", action);
```

2. Enable file logging in debug mode:
```rust
RUST_LOG=debug cargo run
```

3. Use `dbg!()` macro for quick debugging:
```rust
dbg!(&self.state);
```

4. Test with small terminal sizes:
```rust
TERM=xterm-256color cargo run -- --tick-rate 1000
```

Remember: Ratatui is immediate mode - think in terms of "what to draw now" rather than "how to update the widget". Keep components focused, use the message-passing system for communication, and always handle errors gracefully.