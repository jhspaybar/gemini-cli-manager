# Ratatui Widget Cookbook

A collection of ready-to-use widget implementations for common UI patterns.

## Table of Contents
1. [Form Components](#form-components)
2. [Navigation Components](#navigation-components)
3. [Data Display](#data-display)
4. [Feedback Components](#feedback-components)
5. [Layout Helpers](#layout-helpers)

## Form Components

### Text Input with Validation

```rust
use ratatui::prelude::*;
use crossterm::event::{KeyCode, KeyEvent};

pub struct ValidatedInput {
    value: String,
    cursor: usize,
    validator: Box<dyn Fn(&str) -> Result<(), String>>,
    error: Option<String>,
    focused: bool,
}

impl ValidatedInput {
    pub fn new<F>(validator: F) -> Self 
    where 
        F: Fn(&str) -> Result<(), String> + 'static
    {
        Self {
            value: String::new(),
            cursor: 0,
            validator: Box::new(validator),
            error: None,
            focused: false,
        }
    }

    pub fn handle_key(&mut self, key: KeyEvent) {
        match key.code {
            KeyCode::Char(c) => {
                self.value.insert(self.cursor, c);
                self.cursor += 1;
                self.validate();
            }
            KeyCode::Backspace => {
                if self.cursor > 0 {
                    self.value.remove(self.cursor - 1);
                    self.cursor -= 1;
                    self.validate();
                }
            }
            KeyCode::Left => {
                self.cursor = self.cursor.saturating_sub(1);
            }
            KeyCode::Right => {
                self.cursor = (self.cursor + 1).min(self.value.len());
            }
            _ => {}
        }
    }

    fn validate(&mut self) {
        self.error = match (self.validator)(&self.value) {
            Ok(()) => None,
            Err(e) => Some(e),
        };
    }

    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let chunks = Layout::vertical([
            Constraint::Length(3),
            Constraint::Length(1),
        ]).split(area);

        let border_color = match (&self.error, self.focused) {
            (Some(_), _) => Color::Red,
            (None, true) => Color::Cyan,
            (None, false) => Color::Gray,
        };

        let input = Paragraph::new(self.value.as_str())
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(border_color))
            );

        frame.render_widget(input, chunks[0]);

        if let Some(error) = &self.error {
            let error_text = Paragraph::new(error.as_str())
                .style(Style::default().fg(Color::Red));
            frame.render_widget(error_text, chunks[1]);
        }

        if self.focused {
            frame.set_cursor_position(Position::new(
                chunks[0].x + self.cursor as u16 + 1,
                chunks[0].y + 1,
            ));
        }
    }
}

// Usage example
let email_input = ValidatedInput::new(|s| {
    if s.contains('@') {
        Ok(())
    } else {
        Err("Invalid email format".to_string())
    }
});
```

### Multi-Select List

```rust
use std::collections::HashSet;

pub struct MultiSelectList<T> {
    items: Vec<T>,
    selected: HashSet<usize>,
    cursor: usize,
}

impl<T: Display> MultiSelectList<T> {
    pub fn new(items: Vec<T>) -> Self {
        Self {
            items,
            selected: HashSet::new(),
            cursor: 0,
        }
    }

    pub fn toggle_selection(&mut self) {
        if self.selected.contains(&self.cursor) {
            self.selected.remove(&self.cursor);
        } else {
            self.selected.insert(self.cursor);
        }
    }

    pub fn move_up(&mut self) {
        self.cursor = self.cursor.saturating_sub(1);
    }

    pub fn move_down(&mut self) {
        self.cursor = (self.cursor + 1).min(self.items.len().saturating_sub(1));
    }

    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let items: Vec<ListItem> = self.items
            .iter()
            .enumerate()
            .map(|(i, item)| {
                let checkbox = if self.selected.contains(&i) { "☑" } else { "☐" };
                let content = format!("{} {}", checkbox, item);
                let style = if i == self.cursor {
                    Style::default().add_modifier(Modifier::REVERSED)
                } else {
                    Style::default()
                };
                ListItem::new(content).style(style)
            })
            .collect();

        let list = List::new(items)
            .block(Block::default().borders(Borders::ALL).title("Select Items"));

        frame.render_widget(list, area);
    }

    pub fn get_selected(&self) -> Vec<&T> {
        self.selected
            .iter()
            .filter_map(|&i| self.items.get(i))
            .collect()
    }
}
```

### Form Layout

```rust
pub struct Form {
    fields: Vec<FormField>,
    active_field: usize,
}

pub struct FormField {
    label: String,
    input: ValidatedInput,
}

impl Form {
    pub fn render(&mut self, frame: &mut Frame, area: Rect) {
        let field_height = 4; // Label + input + error + spacing
        let constraints: Vec<Constraint> = self.fields
            .iter()
            .map(|_| Constraint::Length(field_height))
            .collect();

        let chunks = Layout::vertical(constraints).split(area);

        for (i, (field, chunk)) in self.fields.iter_mut().zip(chunks.iter()).enumerate() {
            field.input.focused = i == self.active_field;
            
            let field_chunks = Layout::vertical([
                Constraint::Length(1),
                Constraint::Length(3),
            ]).split(*chunk);

            // Render label
            let label = Paragraph::new(field.label.as_str())
                .style(Style::default().add_modifier(Modifier::BOLD));
            frame.render_widget(label, field_chunks[0]);

            // Render input
            field.input.render(frame, field_chunks[1]);
        }
    }

    pub fn next_field(&mut self) {
        self.active_field = (self.active_field + 1) % self.fields.len();
    }

    pub fn previous_field(&mut self) {
        if self.active_field > 0 {
            self.active_field -= 1;
        } else {
            self.active_field = self.fields.len() - 1;
        }
    }
}
```

## Navigation Components

### Tab Bar with Icons

```rust
pub struct IconTab {
    icon: &'static str,
    title: &'static str,
}

pub struct IconTabBar {
    tabs: Vec<IconTab>,
    selected: usize,
}

impl IconTabBar {
    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let tab_width = area.width / self.tabs.len() as u16;
        let mut current_x = area.x;

        for (i, tab) in self.tabs.iter().enumerate() {
            let is_selected = i == self.selected;
            let style = if is_selected {
                Style::default()
                    .fg(Color::Yellow)
                    .add_modifier(Modifier::BOLD)
            } else {
                Style::default().fg(Color::Gray)
            };

            let content = format!("{} {}", tab.icon, tab.title);
            let width = content.len() as u16;
            let x = current_x + (tab_width - width) / 2;

            let span = Span::styled(content, style);
            frame.render_widget(
                Paragraph::new(span),
                Rect::new(x, area.y, width, 1)
            );

            if is_selected {
                // Render underline
                let underline = "─".repeat(tab_width as usize);
                frame.render_widget(
                    Paragraph::new(underline).style(style),
                    Rect::new(current_x, area.y + 1, tab_width, 1)
                );
            }

            current_x += tab_width;
        }
    }
}

// Usage
let tabs = IconTabBar {
    tabs: vec![
        IconTab { icon: "📦", title: "Extensions" },
        IconTab { icon: "👤", title: "Profiles" },
        IconTab { icon: "⚙️", title: "Settings" },
    ],
    selected: 0,
};
```

### Breadcrumb Navigation

```rust
pub struct Breadcrumbs {
    path: Vec<String>,
    separator: String,
}

impl Breadcrumbs {
    pub fn new() -> Self {
        Self {
            path: vec![],
            separator: " › ".to_string(),
        }
    }

    pub fn push(&mut self, segment: String) {
        self.path.push(segment);
    }

    pub fn pop(&mut self) -> Option<String> {
        self.path.pop()
    }

    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let mut spans = vec![];
        
        for (i, segment) in self.path.iter().enumerate() {
            if i > 0 {
                spans.push(Span::raw(&self.separator));
            }
            
            let style = if i == self.path.len() - 1 {
                Style::default().add_modifier(Modifier::BOLD)
            } else {
                Style::default().fg(Color::Gray)
            };
            
            spans.push(Span::styled(segment, style));
        }

        let breadcrumbs = Paragraph::new(Line::from(spans))
            .block(Block::default().borders(Borders::BOTTOM));

        frame.render_widget(breadcrumbs, area);
    }
}
```

### Sidebar Menu

```rust
pub struct MenuItem {
    icon: &'static str,
    label: &'static str,
    id: &'static str,
}

pub struct Sidebar {
    items: Vec<MenuItem>,
    selected: usize,
    expanded: bool,
}

impl Sidebar {
    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let items: Vec<ListItem> = self.items
            .iter()
            .enumerate()
            .map(|(i, item)| {
                let content = if self.expanded {
                    format!("{} {}", item.icon, item.label)
                } else {
                    item.icon.to_string()
                };

                let style = if i == self.selected {
                    Style::default()
                        .bg(Color::DarkGray)
                        .add_modifier(Modifier::BOLD)
                } else {
                    Style::default()
                };

                ListItem::new(content).style(style)
            })
            .collect();

        let list = List::new(items)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
            );

        frame.render_widget(list, area);
    }

    pub fn toggle_expanded(&mut self) {
        self.expanded = !self.expanded;
    }
}
```

## Data Display

### Sortable Table

```rust
use std::cmp::Ordering;

#[derive(Clone)]
pub struct TableColumn {
    header: String,
    width: Constraint,
}

pub struct SortableTable<T> {
    columns: Vec<TableColumn>,
    rows: Vec<T>,
    sort_column: Option<usize>,
    sort_ascending: bool,
    selected: Option<usize>,
}

impl<T> SortableTable<T> {
    pub fn sort_by_column<F>(&mut self, column: usize, extractor: F)
    where
        F: Fn(&T) -> String,
    {
        let ascending = if Some(column) == self.sort_column {
            !self.sort_ascending
        } else {
            true
        };

        self.rows.sort_by(|a, b| {
            let a_val = extractor(a);
            let b_val = extractor(b);
            if ascending {
                a_val.cmp(&b_val)
            } else {
                b_val.cmp(&a_val)
            }
        });

        self.sort_column = Some(column);
        self.sort_ascending = ascending;
    }

    pub fn render<F>(&self, frame: &mut Frame, area: Rect, row_mapper: F)
    where
        F: Fn(&T) -> Vec<String>,
    {
        // Create header with sort indicators
        let header_cells: Vec<Cell> = self.columns
            .iter()
            .enumerate()
            .map(|(i, col)| {
                let mut header = col.header.clone();
                if Some(i) == self.sort_column {
                    header.push_str(if self.sort_ascending { " ▲" } else { " ▼" });
                }
                Cell::from(header).style(Style::default().add_modifier(Modifier::BOLD))
            })
            .collect();

        let header = Row::new(header_cells).height(1).bottom_margin(1);

        // Create rows
        let rows: Vec<Row> = self.rows
            .iter()
            .map(|item| {
                let cells = row_mapper(item);
                Row::new(cells)
            })
            .collect();

        let widths: Vec<Constraint> = self.columns
            .iter()
            .map(|col| col.width.clone())
            .collect();

        let table = Table::new(rows, &widths)
            .header(header)
            .block(Block::default().borders(Borders::ALL))
            .highlight_style(Style::default().add_modifier(Modifier::REVERSED));

        let mut state = TableState::default();
        state.select(self.selected);

        frame.render_stateful_widget(table, area, &mut state);
    }
}
```

### Tree View

```rust
pub struct TreeNode<T> {
    data: T,
    children: Vec<TreeNode<T>>,
    expanded: bool,
}

pub struct TreeView<T> {
    root: Vec<TreeNode<T>>,
    selected_path: Vec<usize>,
}

impl<T: Display> TreeView<T> {
    fn render_node(
        &self,
        node: &TreeNode<T>,
        depth: usize,
        lines: &mut Vec<String>,
        path: &[usize],
        current_path: &mut Vec<usize>,
    ) {
        let indent = "  ".repeat(depth);
        let prefix = if node.children.is_empty() {
            "  "
        } else if node.expanded {
            "▼ "
        } else {
            "▶ "
        };

        let is_selected = current_path == &self.selected_path[..current_path.len()];
        let line = format!("{}{}{}", indent, prefix, node.data);
        
        lines.push(line);

        if node.expanded {
            for (i, child) in node.children.iter().enumerate() {
                current_path.push(i);
                self.render_node(child, depth + 1, lines, path, current_path);
                current_path.pop();
            }
        }
    }

    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let mut lines = vec![];
        let mut current_path = vec![];

        for (i, node) in self.root.iter().enumerate() {
            current_path.push(i);
            self.render_node(node, 0, &mut lines, &self.selected_path, &mut current_path);
            current_path.clear();
        }

        let items: Vec<ListItem> = lines
            .into_iter()
            .enumerate()
            .map(|(i, line)| {
                let style = if i == self.get_selected_line_index() {
                    Style::default().add_modifier(Modifier::REVERSED)
                } else {
                    Style::default()
                };
                ListItem::new(line).style(style)
            })
            .collect();

        let list = List::new(items)
            .block(Block::default().borders(Borders::ALL).title("Tree View"));

        frame.render_widget(list, area);
    }

    fn get_selected_line_index(&self) -> usize {
        // Calculate which line number corresponds to selected path
        // Implementation depends on tree structure
        0
    }
}
```

## Feedback Components

### Progress Dialog

```rust
pub struct ProgressDialog {
    title: String,
    message: String,
    progress: f64,
    is_indeterminate: bool,
}

impl ProgressDialog {
    pub fn render(&self, frame: &mut Frame, parent_area: Rect) {
        let area = centered_rect(60, 20, parent_area);
        
        // Clear background
        frame.render_widget(Clear, area);

        let chunks = Layout::vertical([
            Constraint::Length(3),
            Constraint::Length(2),
            Constraint::Length(3),
            Constraint::Min(0),
        ]).split(area);

        // Border
        let block = Block::default()
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .title(self.title.as_str());
        frame.render_widget(block, area);

        // Message
        let message = Paragraph::new(self.message.as_str())
            .alignment(Alignment::Center);
        frame.render_widget(message, chunks[1]);

        // Progress bar
        if self.is_indeterminate {
            let spinner_frames = vec!["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
            let current_frame = (self.progress * 10.0) as usize % spinner_frames.len();
            let spinner = Paragraph::new(spinner_frames[current_frame])
                .alignment(Alignment::Center);
            frame.render_widget(spinner, chunks[2]);
        } else {
            let gauge = Gauge::default()
                .gauge_style(Style::default().fg(Color::Cyan))
                .percent((self.progress * 100.0) as u16)
                .label(format!("{:.0}%", self.progress * 100.0));
            frame.render_widget(gauge, chunks[2]);
        }
    }
}
```

### Notification Toast

```rust
#[derive(Clone)]
pub struct Notification {
    id: usize,
    message: String,
    level: NotificationLevel,
    created_at: Instant,
}

#[derive(Clone)]
pub enum NotificationLevel {
    Info,
    Success,
    Warning,
    Error,
}

pub struct NotificationStack {
    notifications: Vec<Notification>,
    next_id: usize,
    duration: Duration,
}

impl NotificationStack {
    pub fn new() -> Self {
        Self {
            notifications: vec![],
            next_id: 0,
            duration: Duration::from_secs(3),
        }
    }

    pub fn push(&mut self, message: String, level: NotificationLevel) {
        self.notifications.push(Notification {
            id: self.next_id,
            message,
            level,
            created_at: Instant::now(),
        });
        self.next_id += 1;
    }

    pub fn update(&mut self) {
        let now = Instant::now();
        self.notifications.retain(|n| now.duration_since(n.created_at) < self.duration);
    }

    pub fn render(&self, frame: &mut Frame, area: Rect) {
        let notification_height = 3;
        let spacing = 1;
        let max_notifications = 5;

        for (i, notification) in self.notifications.iter().rev().take(max_notifications).enumerate() {
            let y = area.y + (i as u16 * (notification_height + spacing));
            let notification_area = Rect::new(
                area.x,
                y,
                area.width,
                notification_height,
            );

            let (border_color, icon) = match notification.level {
                NotificationLevel::Info => (Color::Blue, "ℹ"),
                NotificationLevel::Success => (Color::Green, "✓"),
                NotificationLevel::Warning => (Color::Yellow, "⚠"),
                NotificationLevel::Error => (Color::Red, "✗"),
            };

            let content = format!("{} {}", icon, notification.message);
            let notification_widget = Paragraph::new(content)
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .border_style(Style::default().fg(border_color))
                )
                .wrap(Wrap { trim: true });

            frame.render_widget(Clear, notification_area);
            frame.render_widget(notification_widget, notification_area);
        }
    }
}
```

### Confirmation Dialog

```rust
pub struct ConfirmDialog {
    title: String,
    message: String,
    confirm_text: String,
    cancel_text: String,
    selected_button: bool, // false = cancel, true = confirm
}

impl ConfirmDialog {
    pub fn new(title: String, message: String) -> Self {
        Self {
            title,
            message,
            confirm_text: "Confirm".to_string(),
            cancel_text: "Cancel".to_string(),
            selected_button: false,
        }
    }

    pub fn toggle_selection(&mut self) {
        self.selected_button = !self.selected_button;
    }

    pub fn render(&self, frame: &mut Frame, parent_area: Rect) {
        let area = centered_rect(50, 20, parent_area);
        frame.render_widget(Clear, area);

        let chunks = Layout::vertical([
            Constraint::Length(1),
            Constraint::Min(0),
            Constraint::Length(3),
        ]).split(area);

        // Border and title
        let block = Block::default()
            .borders(Borders::ALL)
            .border_type(BorderType::Double)
            .title(self.title.as_str())
            .title_alignment(Alignment::Center);
        frame.render_widget(block, area);

        // Message
        let message = Paragraph::new(self.message.as_str())
            .alignment(Alignment::Center)
            .wrap(Wrap { trim: true });
        frame.render_widget(message, chunks[1]);

        // Buttons
        let button_chunks = Layout::horizontal([
            Constraint::Percentage(50),
            Constraint::Percentage(50),
        ]).split(chunks[2]);

        let cancel_style = if !self.selected_button {
            Style::default().bg(Color::Red).fg(Color::White)
        } else {
            Style::default()
        };

        let confirm_style = if self.selected_button {
            Style::default().bg(Color::Green).fg(Color::White)
        } else {
            Style::default()
        };

        let cancel_button = Paragraph::new(self.cancel_text.as_str())
            .style(cancel_style)
            .alignment(Alignment::Center);

        let confirm_button = Paragraph::new(self.confirm_text.as_str())
            .style(confirm_style)
            .alignment(Alignment::Center);

        frame.render_widget(cancel_button, button_chunks[0]);
        frame.render_widget(confirm_button, button_chunks[1]);
    }
}
```

## Layout Helpers

### Centered Rectangle

```rust
pub fn centered_rect(percent_x: u16, percent_y: u16, area: Rect) -> Rect {
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

### Grid Layout

```rust
pub struct Grid {
    columns: u16,
    rows: u16,
    column_spacing: u16,
    row_spacing: u16,
}

impl Grid {
    pub fn new(columns: u16, rows: u16) -> Self {
        Self {
            columns,
            rows,
            column_spacing: 1,
            row_spacing: 1,
        }
    }

    pub fn split(&self, area: Rect) -> Vec<Vec<Rect>> {
        let total_column_spacing = self.column_spacing * (self.columns - 1);
        let total_row_spacing = self.row_spacing * (self.rows - 1);

        let cell_width = (area.width - total_column_spacing) / self.columns;
        let cell_height = (area.height - total_row_spacing) / self.rows;

        let mut grid = vec![];

        for row in 0..self.rows {
            let mut row_cells = vec![];
            let y = area.y + row * (cell_height + self.row_spacing);

            for col in 0..self.columns {
                let x = area.x + col * (cell_width + self.column_spacing);
                row_cells.push(Rect::new(x, y, cell_width, cell_height));
            }

            grid.push(row_cells);
        }

        grid
    }
}

// Usage
let grid = Grid::new(3, 2);
let cells = grid.split(area);
// cells[0][0] is top-left, cells[1][2] is bottom-right
```

### Responsive Container

```rust
pub struct ResponsiveContainer {
    breakpoints: Vec<(u16, Layout)>,
}

impl ResponsiveContainer {
    pub fn new() -> Self {
        Self {
            breakpoints: vec![],
        }
    }

    pub fn add_breakpoint(mut self, min_width: u16, layout: Layout) -> Self {
        self.breakpoints.push((min_width, layout));
        self.breakpoints.sort_by_key(|&(w, _)| w);
        self
    }

    pub fn split(&self, area: Rect) -> Vec<Rect> {
        let layout = self.breakpoints
            .iter()
            .rev()
            .find(|&&(min_width, _)| area.width >= min_width)
            .map(|(_, layout)| layout)
            .unwrap_or(&Layout::default());

        layout.split(area).to_vec()
    }
}

// Usage
let container = ResponsiveContainer::new()
    .add_breakpoint(0, Layout::vertical([Constraint::Percentage(100)]))
    .add_breakpoint(80, Layout::horizontal([
        Constraint::Percentage(30),
        Constraint::Percentage(70),
    ]))
    .add_breakpoint(120, Layout::horizontal([
        Constraint::Length(20),
        Constraint::Min(0),
        Constraint::Length(20),
    ]));
```

This cookbook provides ready-to-use components that can be easily integrated into any Ratatui application. Each component is self-contained and can be customized to fit specific needs.