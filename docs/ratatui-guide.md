# Ratatui Complete Guide

This guide covers everything you need to know about using Ratatui for the Gemini CLI Manager.

## Table of Contents
1. [Introduction to Ratatui](#introduction-to-ratatui)
2. [Core Concepts](#core-concepts)
3. [Widgets Deep Dive](#widgets-deep-dive)
4. [Layout System](#layout-system)
5. [Styling and Colors](#styling-and-colors)
6. [Event Handling](#event-handling)
7. [State Management Patterns](#state-management-patterns)
8. [Performance Optimization](#performance-optimization)
9. [Common Recipes](#common-recipes)
10. [Troubleshooting](#troubleshooting)

## Introduction to Ratatui

Ratatui is a Rust library for building Terminal User Interfaces (TUIs). It's inspired by the JavaScript library `blessed-contrib` and the Go library `termui`.

### Key Features
- **Immediate mode rendering**: Redraw the entire UI each frame
- **Backend agnostic**: Works with crossterm, termion, or termwiz
- **No runtime overhead**: Zero-cost abstractions
- **Flexible layouts**: Constraint-based layout system
- **Rich widget library**: Tables, charts, gauges, and more

### When to Use Ratatui
- Building CLI tools with interactive interfaces
- Creating system monitoring dashboards
- Developing terminal-based games
- Building developer tools with rich output

## Core Concepts

### 1. The Frame

The `Frame` is your canvas for drawing:

```rust
fn draw<B: Backend>(frame: &mut Frame<B>) {
    // All drawing happens through the frame
    frame.render_widget(widget, area);
}
```

### 2. Immediate Mode Rendering

Unlike retained mode GUIs, Ratatui redraws everything each frame:

```rust
// This runs 60 times per second
terminal.draw(|frame| {
    // Clear and redraw entire UI
    ui.draw(frame);
})?;
```

### 3. The Backend

Ratatui abstracts terminal operations through backends:

```rust
// Crossterm backend (recommended)
let backend = CrosstermBackend::new(stdout());
let mut terminal = Terminal::new(backend)?;

// Enable raw mode
terminal::enable_raw_mode()?;
```

### 4. Widgets

Widgets are the building blocks of your UI:

```rust
// Widgets implement the Widget trait
pub trait Widget {
    fn render(self, area: Rect, buf: &mut Buffer);
}
```

## Widgets Deep Dive

### Block

The fundamental container widget:

```rust
let block = Block::default()
    .title("My Widget")
    .title_alignment(Alignment::Center)
    .borders(Borders::ALL)
    .border_type(BorderType::Rounded)
    .border_style(Style::default().fg(Color::Cyan))
    .style(Style::default().bg(Color::Black));
```

### Paragraph

For displaying text:

```rust
let text = vec![
    Line::from("First line"),
    Line::from(vec![
        Span::raw("Normal "),
        Span::styled("styled", Style::default().fg(Color::Red)),
        Span::raw(" text"),
    ]),
];

let paragraph = Paragraph::new(text)
    .block(Block::default().borders(Borders::ALL))
    .style(Style::default().fg(Color::White))
    .alignment(Alignment::Left)
    .wrap(Wrap { trim: true });
```

### List

For selectable lists:

```rust
// Stateless list
let items = vec![
    ListItem::new("Item 1"),
    ListItem::new("Item 2").style(Style::default().fg(Color::Yellow)),
    ListItem::new("Item 3"),
];

let list = List::new(items)
    .block(Block::default().borders(Borders::ALL).title("List"))
    .highlight_style(Style::default().add_modifier(Modifier::REVERSED))
    .highlight_symbol("> ");

// Stateful rendering
let mut list_state = ListState::default();
list_state.select(Some(1));

frame.render_stateful_widget(list, area, &mut list_state);
```

### Table

For structured data:

```rust
let header = Row::new(vec!["Name", "Age", "City"])
    .style(Style::default().fg(Color::Yellow))
    .height(1)
    .bottom_margin(1);

let rows = vec![
    Row::new(vec!["Alice", "30", "New York"]),
    Row::new(vec!["Bob", "25", "San Francisco"]),
];

let widths = &[
    Constraint::Length(10),
    Constraint::Length(5),
    Constraint::Min(10),
];

let table = Table::new(rows, widths)
    .header(header)
    .block(Block::default().borders(Borders::ALL).title("Table"))
    .highlight_style(Style::default().add_modifier(Modifier::REVERSED));

// With state for selection
let mut table_state = TableState::default();
table_state.select(Some(0));

frame.render_stateful_widget(table, area, &mut table_state);
```

### Gauge

For progress indicators:

```rust
let gauge = Gauge::default()
    .block(Block::default().borders(Borders::ALL).title("Progress"))
    .gauge_style(Style::default().fg(Color::Green))
    .percent(75)
    .label(format!("{}%", 75));
```

### Chart

For data visualization:

```rust
let datasets = vec![
    Dataset::default()
        .name("data1")
        .marker(symbols::Marker::Dot)
        .graph_type(GraphType::Scatter)
        .style(Style::default().fg(Color::Cyan))
        .data(&[(0.0, 5.0), (1.0, 6.0), (2.0, 7.0)]),
];

let chart = Chart::new(datasets)
    .block(Block::default().title("Chart"))
    .x_axis(
        Axis::default()
            .title("X Axis")
            .bounds([0.0, 10.0])
            .labels(vec!["0", "5", "10"]),
    )
    .y_axis(
        Axis::default()
            .title("Y Axis")
            .bounds([0.0, 10.0])
            .labels(vec!["0", "5", "10"]),
    );
```

### Tabs

For navigation:

```rust
let titles = vec!["Tab1", "Tab2", "Tab3"];
let tabs = Tabs::new(titles)
    .block(Block::default().borders(Borders::ALL).title("Tabs"))
    .select(0)
    .style(Style::default().fg(Color::White))
    .highlight_style(Style::default().fg(Color::Yellow));
```

### Custom Widgets

Create your own widgets:

```rust
pub struct MyWidget {
    text: String,
}

impl Widget for MyWidget {
    fn render(self, area: Rect, buf: &mut Buffer) {
        // Custom rendering logic
        buf.set_string(area.x, area.y, &self.text, Style::default());
    }
}
```

## Layout System

### Basic Layouts

```rust
// Vertical split
let chunks = Layout::default()
    .direction(Direction::Vertical)
    .constraints([
        Constraint::Length(3),      // Fixed height
        Constraint::Min(0),         // Take remaining space
        Constraint::Percentage(25), // 25% of available space
    ])
    .split(area);
```

### Advanced Layouts

```rust
// Nested layouts
let main_chunks = Layout::default()
    .direction(Direction::Vertical)
    .constraints([Constraint::Length(3), Constraint::Min(0)])
    .split(area);

let body_chunks = Layout::default()
    .direction(Direction::Horizontal)
    .constraints([Constraint::Percentage(30), Constraint::Percentage(70)])
    .split(main_chunks[1]);

// Using the new const array syntax
let [header, body, footer] = Layout::vertical([
    Constraint::Length(3),
    Constraint::Min(0),
    Constraint::Length(3),
]).areas(area);
```

### Centering Content

```rust
fn centered_rect(percent_x: u16, percent_y: u16, area: Rect) -> Rect {
    let popup_layout = Layout::default()
        .direction(Direction::Vertical)
        .constraints([
            Constraint::Percentage((100 - percent_y) / 2),
            Constraint::Percentage(percent_y),
            Constraint::Percentage((100 - percent_y) / 2),
        ])
        .split(area);

    Layout::default()
        .direction(Direction::Horizontal)
        .constraints([
            Constraint::Percentage((100 - percent_x) / 2),
            Constraint::Percentage(percent_x),
            Constraint::Percentage((100 - percent_x) / 2),
        ])
        .split(popup_layout[1])[1]
}
```

### Responsive Design

```rust
fn get_constraints(width: u16) -> Vec<Constraint> {
    if width < 80 {
        // Mobile layout
        vec![Constraint::Percentage(100)]
    } else if width < 120 {
        // Tablet layout
        vec![Constraint::Percentage(40), Constraint::Percentage(60)]
    } else {
        // Desktop layout
        vec![
            Constraint::Length(20),  // Sidebar
            Constraint::Min(80),     // Main content
            Constraint::Length(30),  // Right panel
        ]
    }
}
```

## Styling and Colors

### Style Basics

```rust
let style = Style::default()
    .fg(Color::White)
    .bg(Color::Black)
    .add_modifier(Modifier::BOLD | Modifier::ITALIC);
```

### Color Options

```rust
// Basic colors
Color::Black, Color::Red, Color::Green, Color::Yellow,
Color::Blue, Color::Magenta, Color::Cyan, Color::White,
Color::Gray, Color::DarkGray, Color::LightRed, // ... etc

// RGB colors
Color::Rgb(255, 0, 0)  // Red

// Indexed colors (0-255)
Color::Indexed(200)
```

### Modifiers

```rust
Modifier::BOLD
Modifier::DIM
Modifier::ITALIC
Modifier::UNDERLINED
Modifier::SLOW_BLINK
Modifier::RAPID_BLINK
Modifier::REVERSED
Modifier::HIDDEN
Modifier::CROSSED_OUT
```

### Theme System

```rust
#[derive(Clone)]
pub struct Theme {
    pub base: Style,
    pub highlight: Style,
    pub error: Style,
    pub success: Style,
    pub warning: Style,
}

impl Theme {
    pub fn dark() -> Self {
        Self {
            base: Style::default().fg(Color::White).bg(Color::Black),
            highlight: Style::default().fg(Color::Black).bg(Color::Cyan),
            error: Style::default().fg(Color::Red),
            success: Style::default().fg(Color::Green),
            warning: Style::default().fg(Color::Yellow),
        }
    }
}
```

## Event Handling

### Basic Event Loop

```rust
loop {
    terminal.draw(|frame| ui.draw(frame))?;
    
    if event::poll(Duration::from_millis(16))? {
        if let Event::Key(key) = event::read()? {
            match key.code {
                KeyCode::Char('q') => break,
                KeyCode::Up => ui.previous(),
                KeyCode::Down => ui.next(),
                _ => {}
            }
        }
    }
}
```

### Advanced Event Handling

```rust
pub fn handle_events(&mut self) -> Result<()> {
    match event::read()? {
        Event::Key(key) => self.handle_key(key),
        Event::Mouse(mouse) => self.handle_mouse(mouse),
        Event::Resize(width, height) => self.handle_resize(width, height),
        Event::FocusGained => self.handle_focus_gained(),
        Event::FocusLost => self.handle_focus_lost(),
        Event::Paste(data) => self.handle_paste(data),
    }
    Ok(())
}
```

### Key Combinations

```rust
fn handle_key(&mut self, key: KeyEvent) {
    match (key.code, key.modifiers) {
        (KeyCode::Char('c'), KeyModifiers::CONTROL) => self.copy(),
        (KeyCode::Char('v'), KeyModifiers::CONTROL) => self.paste(),
        (KeyCode::Tab, KeyModifiers::NONE) => self.next_widget(),
        (KeyCode::BackTab, KeyModifiers::SHIFT) => self.previous_widget(),
        _ => {}
    }
}
```

## State Management Patterns

### Component State Pattern

```rust
pub struct AppState {
    pub active_tab: TabIndex,
    pub extension_list: ExtensionListState,
    pub profile_manager: ProfileManagerState,
}

impl AppState {
    pub fn handle_action(&mut self, action: Action) -> Option<Action> {
        match self.active_tab {
            TabIndex::Extensions => self.extension_list.handle_action(action),
            TabIndex::Profiles => self.profile_manager.handle_action(action),
        }
    }
}
```

### Redux-like Pattern

```rust
#[derive(Clone)]
pub struct Store {
    state: AppState,
    subscribers: Vec<Sender<AppState>>,
}

impl Store {
    pub fn dispatch(&mut self, action: Action) {
        self.state = reducer(&self.state, action);
        self.notify_subscribers();
    }
}

fn reducer(state: &AppState, action: Action) -> AppState {
    match action {
        Action::SelectExtension(id) => AppState {
            selected_extension: Some(id),
            ..state.clone()
        },
        // ... other actions
    }
}
```

## Performance Optimization

### Minimize Allocations

```rust
pub struct OptimizedList {
    // Pre-allocate and reuse
    items: Vec<ListItem<'static>>,
    buffer: String,
}

impl OptimizedList {
    pub fn update_items(&mut self, data: &[&str]) {
        self.items.clear();
        for item in data {
            self.buffer.clear();
            self.buffer.push_str(item);
            self.items.push(ListItem::new(self.buffer.clone()));
        }
    }
}
```

### Efficient Rendering

```rust
pub struct CachedWidget {
    content: String,
    last_size: (u16, u16),
    cached_lines: Vec<String>,
}

impl CachedWidget {
    pub fn render(&mut self, area: Rect) -> Paragraph {
        let size = (area.width, area.height);
        if size != self.last_size {
            self.recalculate_lines(size);
            self.last_size = size;
        }
        
        Paragraph::new(self.cached_lines.join("\n"))
    }
}
```

### Conditional Rendering

```rust
pub struct SmartComponent {
    needs_redraw: bool,
    cached_widget: Option<Paragraph<'static>>,
}

impl SmartComponent {
    pub fn render(&mut self) -> &Paragraph {
        if self.needs_redraw {
            self.cached_widget = Some(self.build_widget());
            self.needs_redraw = false;
        }
        self.cached_widget.as_ref().unwrap()
    }
}
```

## Common Recipes

### Modal Dialog

```rust
pub fn render_modal(frame: &mut Frame, title: &str, content: &str) {
    let area = centered_rect(60, 20, frame.area());
    
    // Clear the background
    let clear = Clear;
    frame.render_widget(clear, area);
    
    // Render the modal
    let block = Block::default()
        .title(title)
        .borders(Borders::ALL)
        .border_type(BorderType::Rounded);
    
    let paragraph = Paragraph::new(content)
        .block(block)
        .wrap(Wrap { trim: true });
    
    frame.render_widget(paragraph, area);
}
```

### Loading Indicator

```rust
pub struct Spinner {
    frames: Vec<&'static str>,
    current: usize,
}

impl Spinner {
    pub fn new() -> Self {
        Self {
            frames: vec!["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"],
            current: 0,
        }
    }
    
    pub fn tick(&mut self) {
        self.current = (self.current + 1) % self.frames.len();
    }
    
    pub fn frame(&self) -> &str {
        self.frames[self.current]
    }
}
```

### Scrollable Content

```rust
pub struct ScrollableList {
    items: Vec<String>,
    state: ListState,
    offset: usize,
}

impl ScrollableList {
    pub fn render(&mut self, frame: &mut Frame, area: Rect) {
        let visible_items = (area.height as usize).saturating_sub(2); // Account for borders
        
        let items: Vec<ListItem> = self.items
            .iter()
            .skip(self.offset)
            .take(visible_items)
            .map(|i| ListItem::new(i.as_str()))
            .collect();
        
        let list = List::new(items)
            .block(Block::default().borders(Borders::ALL))
            .highlight_style(Style::default().add_modifier(Modifier::REVERSED));
        
        frame.render_stateful_widget(list, area, &mut self.state);
    }
    
    pub fn scroll_down(&mut self) {
        if self.offset + 1 < self.items.len() {
            self.offset += 1;
        }
    }
}
```

### Input Field

```rust
pub struct InputField {
    input: String,
    cursor_position: usize,
}

impl InputField {
    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let input = Paragraph::new(self.input.as_str())
            .block(Block::default().borders(Borders::ALL).title("Input"));
        
        frame.render_widget(input, area);
        
        // Show cursor
        frame.set_cursor_position(Position::new(
            area.x + self.cursor_position as u16 + 1,
            area.y + 1,
        ));
    }
    
    pub fn handle_input(&mut self, key: KeyEvent) {
        match key.code {
            KeyCode::Char(c) => {
                self.input.insert(self.cursor_position, c);
                self.cursor_position += 1;
            }
            KeyCode::Backspace => {
                if self.cursor_position > 0 {
                    self.input.remove(self.cursor_position - 1);
                    self.cursor_position -= 1;
                }
            }
            KeyCode::Left => {
                self.cursor_position = self.cursor_position.saturating_sub(1);
            }
            KeyCode::Right => {
                if self.cursor_position < self.input.len() {
                    self.cursor_position += 1;
                }
            }
            _ => {}
        }
    }
}
```

## Troubleshooting

### Common Issues

1. **Flickering**: Ensure you're not calling `terminal.clear()` manually
2. **Cut-off content**: Check your constraints sum to 100% or leave one as `Min(0)`
3. **Unresponsive UI**: Make sure event polling timeout is reasonable (16ms for 60fps)
4. **Memory leaks**: Clear collections instead of recreating them

### Debugging Tips

```rust
// Debug layout areas
fn debug_layout(frame: &mut Frame, area: Rect, label: &str) {
    let debug = Block::default()
        .title(format!("{}: {}x{} at ({},{})", 
            label, area.width, area.height, area.x, area.y))
        .borders(Borders::ALL)
        .border_style(Style::default().fg(Color::Red));
    frame.render_widget(debug, area);
}

// Log to file while TUI is running
use std::fs::OpenOptions;
use std::io::Write;

fn debug_log(msg: &str) {
    let mut file = OpenOptions::new()
        .create(true)
        .append(true)
        .open("debug.log")
        .unwrap();
    writeln!(file, "{}: {}", chrono::Local::now(), msg).unwrap();
}
```

### Performance Profiling

```rust
use std::time::Instant;

let start = Instant::now();
terminal.draw(|frame| ui.draw(frame))?;
let elapsed = start.elapsed();

if elapsed.as_millis() > 16 {
    debug_log(&format!("Slow frame: {:?}", elapsed));
}
```

Remember: Ratatui is just a rendering library. It doesn't dictate how you structure your application. Use the patterns that make sense for your use case!