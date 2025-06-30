use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{
    action::Action,
    config::Config,
    models::Extension,
};

pub struct ExtensionList {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    extensions: Vec<Extension>,
    selected: usize,
}

impl Default for ExtensionList {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            extensions: Extension::mock_extensions(),
            selected: 0,
        }
    }
}

impl ExtensionList {
    pub fn new() -> Self {
        Self::default()
    }
    
    fn next(&mut self) {
        if !self.extensions.is_empty() {
            self.selected = (self.selected + 1) % self.extensions.len();
        }
    }
    
    fn previous(&mut self) {
        if !self.extensions.is_empty() {
            if self.selected > 0 {
                self.selected -= 1;
            } else {
                self.selected = self.extensions.len() - 1;
            }
        }
    }
    
    fn get_selected_extension(&self) -> Option<&Extension> {
        self.extensions.get(self.selected)
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
            _ => {}
        }
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Create a block for the extension list
        let block = Block::default()
            .title(" Extensions ")
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(Color::Cyan));
        
        // Create list items
        let items: Vec<ListItem> = self.extensions
            .iter()
            .enumerate()
            .map(|(i, ext)| {
                let is_selected = i == self.selected;
                
                // Build the display string
                let content = vec![
                    Line::from(vec![
                        Span::styled(
                            &ext.name,
                            if is_selected {
                                Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)
                            } else {
                                Style::default().fg(Color::White)
                            }
                        ),
                        Span::raw(" "),
                        Span::styled(
                            format!("v{}", ext.version),
                            Style::default().fg(Color::DarkGray)
                        ),
                    ]),
                    Line::from(vec![
                        Span::raw("  "),
                        Span::styled(
                            ext.description.as_deref().unwrap_or("No description"),
                            Style::default().fg(Color::Gray)
                        ),
                    ]),
                    Line::from(vec![
                        Span::raw("  "),
                        Span::styled(
                            format!("{} MCP servers", ext.mcp_servers.len()),
                            Style::default().fg(Color::Magenta)
                        ),
                        Span::raw(" | "),
                        Span::styled(
                            format!("{} tags", ext.metadata.tags.len()),
                            Style::default().fg(Color::Blue)
                        ),
                    ]),
                    Line::from(""), // Empty line for spacing
                ];
                
                ListItem::new(content)
            })
            .collect();
        
        // Create the list widget
        let list = List::new(items)
            .block(block)
            .highlight_style(Style::default().bg(Color::DarkGray))
            .highlight_symbol("│ ");
        
        // Create a stateful list to track selection
        let mut state = ListState::default();
        state.select(Some(self.selected));
        
        // Render the list
        frame.render_stateful_widget(list, area, &mut state);
        
        // Add help text at the bottom
        if area.height > 4 {
            let help_text = " ↑/↓: Navigate | Enter: View Details | i: Import | n: New | q: Quit ";
            let help_style = Style::default().fg(Color::DarkGray);
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
                help_area
            );
        }
        
        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::KeyCode;
        
        match event {
            Some(crate::tui::Event::Key(key)) => match key.code {
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
                KeyCode::Char('i') => Ok(Some(Action::ImportExtension)),
                KeyCode::Char('n') => Ok(Some(Action::CreateNewExtension)),
                KeyCode::Char('q') => Ok(Some(Action::Quit)),
                _ => Ok(None),
            },
            _ => Ok(None),
        }
    }
}