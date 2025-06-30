use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;
use tui_input::Input;
use tui_input::backend::crossterm::EventHandler;

use super::Component;
use crate::{action::Action, config::Config, models::Extension, storage::Storage, theme};

pub struct ExtensionList {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    extensions: Vec<Extension>,
    filtered_extensions: Vec<usize>, // Indices of extensions that match filter
    selected: usize,
    storage: Option<Storage>,
    search_mode: bool,
    search_input: Input,
}

impl Default for ExtensionList {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            extensions: Vec::new(),
            filtered_extensions: Vec::new(),
            selected: 0,
            storage: None,
            search_mode: false,
            search_input: Input::default(),
        }
    }
}

impl ExtensionList {
    pub fn with_storage(storage: Storage) -> Self {
        let mut list = Self::default();
        list.storage = Some(storage.clone());
        
        // Load extensions from storage
        if let Ok(extensions) = storage.list_extensions() {
            list.extensions = extensions;
            list.update_filter();
        }
        
        list
    }
    
    fn update_filter(&mut self) {
        let search_query = self.search_input.value();
        if search_query.is_empty() {
            // Show all extensions
            self.filtered_extensions = (0..self.extensions.len()).collect();
        } else {
            // Filter extensions based on search query
            let query = search_query.to_lowercase();
            self.filtered_extensions = self.extensions
                .iter()
                .enumerate()
                .filter(|(_, ext)| {
                    ext.name.to_lowercase().contains(&query) ||
                    ext.description.as_ref().map_or(false, |d| d.to_lowercase().contains(&query)) ||
                    ext.metadata.tags.iter().any(|tag| tag.to_lowercase().contains(&query)) ||
                    ext.mcp_servers.keys().any(|server_name| server_name.to_lowercase().contains(&query))
                })
                .map(|(i, _)| i)
                .collect();
        }
        
        // Adjust selection if needed
        if self.selected >= self.filtered_extensions.len() && !self.filtered_extensions.is_empty() {
            self.selected = self.filtered_extensions.len() - 1;
        } else if self.filtered_extensions.is_empty() {
            self.selected = 0;
        }
    }

    fn next(&mut self) {
        if !self.filtered_extensions.is_empty() {
            self.selected = (self.selected + 1) % self.filtered_extensions.len();
        }
    }

    fn previous(&mut self) {
        if !self.filtered_extensions.is_empty() {
            if self.selected > 0 {
                self.selected -= 1;
            } else {
                self.selected = self.filtered_extensions.len() - 1;
            }
        }
    }

    fn get_selected_extension(&self) -> Option<&Extension> {
        self.filtered_extensions.get(self.selected)
            .and_then(|&idx| self.extensions.get(idx))
    }
}

impl Component for ExtensionList {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.command_tx = Some(tx);
        Ok(())
    }

    fn register_config_handler(&mut self, config: Config) -> Result<()> {
        self.config = config;
        Ok(())
    }

    fn update(&mut self, action: Action) -> Result<Option<Action>> {
        match action {
            Action::Tick => {
                // No tick logic needed for now
            }
            Action::Render => {
                // No render-specific logic needed
            }
            Action::RefreshExtensions => {
                // Reload extensions from storage
                if let Some(storage) = &self.storage {
                    if let Ok(extensions) = storage.list_extensions() {
                        self.extensions = extensions;
                        self.update_filter();
                    }
                }
            }
            _ => {}
        }
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Split area if in search mode
        let (search_area, list_area) = if self.search_mode {
            let chunks = ratatui::layout::Layout::default()
                .direction(ratatui::layout::Direction::Vertical)
                .constraints([
                    ratatui::layout::Constraint::Length(3), // Search bar
                    ratatui::layout::Constraint::Min(3),    // List
                ])
                .split(area);
            (Some(chunks[0]), chunks[1])
        } else {
            (None, area)
        };
        
        // Draw search bar if in search mode
        if let Some(search_area) = search_area {
            let search_block = Block::default()
                .title(" Search (Esc to close) ")
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded)
                .border_style(Style::default().fg(theme::warning()));
            
            // Use tui-input's widget with proper styling
            let input_widget = Paragraph::new(self.search_input.value())
                .style(Style::default())
                .block(search_block);
            
            frame.render_widget(input_widget, search_area);
            
            // Set cursor position using tui-input's cursor position
            if self.search_mode {
                // Get the cursor position from the input
                let cursor_pos = self.search_input.visual_cursor();
                // Account for the border (1 char) and block padding
                frame.set_cursor_position((
                    search_area.x + cursor_pos as u16 + 1,
                    search_area.y + 1
                ));
            }
        }
        
        // Create a block for the extension list
        let title = if self.search_mode && !self.search_input.value().is_empty() {
            format!(" Extensions ({}/{}) ", self.filtered_extensions.len(), self.extensions.len())
        } else {
            " Extensions ".to_string()
        };
        
        let block = Block::default()
            .title(title)
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(theme::primary()));

        // Create list items
        let items: Vec<ListItem> = self
            .filtered_extensions
            .iter()
            .enumerate()
            .filter_map(|(i, &ext_idx)| {
                self.extensions.get(ext_idx).map(|ext| {
                    let is_selected = i == self.selected;

                // Build the display string
                let content = vec![
                    Line::from(vec![
                        Span::styled(
                            &ext.name,
                            if is_selected {
                                Style::default()
                                    .fg(theme::text_primary())
                                    .add_modifier(Modifier::BOLD)
                            } else {
                                Style::default().fg(theme::text_primary())
                            },
                        ),
                        Span::raw(" "),
                        Span::styled(
                            format!("v{}", ext.version),
                            Style::default().fg(theme::text_muted()),
                        ),
                    ]),
                    Line::from(vec![
                        Span::raw("  "),
                        Span::styled(
                            ext.description.as_deref().unwrap_or("No description"),
                            Style::default().fg(theme::text_secondary()),
                        ),
                    ]),
                    Line::from(vec![
                        Span::raw("  "),
                        Span::styled(
                            format!("{} MCP servers", ext.mcp_servers.len()),
                            Style::default().fg(theme::accent()),
                        ),
                        Span::raw(" | "),
                        Span::styled(
                            format!("{} tags", ext.metadata.tags.len()),
                            Style::default().fg(theme::primary()),
                        ),
                    ]),
                    Line::from(""), // Empty line for spacing
                ];

                    ListItem::new(content)
                })
            })
            .collect();

        // Create the list widget
        let list = List::new(items)
            .block(block)
            .highlight_style(Style::default().bg(theme::selection()))
            .highlight_symbol("│ ");

        // Create a stateful list to track selection
        let mut state = ListState::default();
        state.select(Some(self.selected));

        // Render the list
        frame.render_stateful_widget(list, list_area, &mut state);

        // Add help text at the bottom
        if list_area.height > 4 {
            let help_text = if self.search_mode {
                " Type to search | Esc: Close search | ↑/↓: Navigate results "
            } else {
                " ↑/↓: Navigate | Enter: View | e: Edit | n: New | d: Delete | /: Search | Tab: Profiles | q: Quit "
            };
            let help_style = Style::default().fg(theme::text_muted());
            let help_area = Rect {
                x: area.x + 1,
                y: area.y + area.height - 1,
                width: area.width.saturating_sub(2),
                height: 1,
            };
            frame.render_widget(
                Paragraph::new(help_text)
                    .style(help_style)
                    .alignment(Alignment::Center),
                help_area,
            );
        }

        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::KeyCode;

        match event {
            Some(crate::tui::Event::Key(key)) => {
                if self.search_mode {
                    // Handle search mode input
                    match key.code {
                        KeyCode::Esc => {
                            self.search_mode = false;
                            self.search_input.reset();
                            self.update_filter();
                            Ok(Some(Action::Render))
                        }
                        KeyCode::Up => {
                            self.previous();
                            Ok(Some(Action::Render))
                        }
                        KeyCode::Down => {
                            self.next();
                            Ok(Some(Action::Render))
                        }
                        KeyCode::Enter => {
                            if let Some(ext) = self.get_selected_extension() {
                                Ok(Some(Action::ViewExtensionDetails(ext.id.clone())))
                            } else {
                                Ok(None)
                            }
                        }
                        _ => {
                            // Let tui-input handle the key event
                            if self.search_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                self.update_filter();
                                Ok(Some(Action::Render))
                            } else {
                                Ok(None)
                            }
                        }
                    }
                } else {
                    // Normal mode
                    match key.code {
                        KeyCode::Up | KeyCode::Char('k') => {
                            self.previous();
                            Ok(Some(Action::Render))
                        }
                        KeyCode::Down | KeyCode::Char('j') => {
                            self.next();
                            Ok(Some(Action::Render))
                        }
                        KeyCode::Enter => {
                            if let Some(ext) = self.get_selected_extension() {
                                Ok(Some(Action::ViewExtensionDetails(ext.id.clone())))
                            } else {
                                Ok(None)
                            }
                        }
                        KeyCode::Char('n') => Ok(Some(Action::CreateNewExtension)),
                        KeyCode::Char('e') => {
                            if let Some(ext) = self.get_selected_extension() {
                                Ok(Some(Action::EditExtension(ext.id.clone())))
                            } else {
                                Ok(None)
                            }
                        }
                        KeyCode::Char('d') => {
                            if let Some(ext) = self.get_selected_extension() {
                                Ok(Some(Action::DeleteExtension(ext.id.clone())))
                            } else {
                                Ok(None)
                            }
                        }
                        KeyCode::Char('/') => {
                            // Enter search mode
                            self.search_mode = true;
                            self.search_input.reset();
                            Ok(Some(Action::Render))
                        }
                        KeyCode::Char('q') => Ok(Some(Action::Quit)),
                        KeyCode::Tab => Ok(Some(Action::NavigateToProfiles)),
                        _ => Ok(None),
                    }
                }
            }
            _ => Ok(None),
        }
    }
}
