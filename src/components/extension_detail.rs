use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, models::Extension};

pub struct ExtensionDetail {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    extension: Option<Extension>,
    scroll_offset: u16,
}

impl Default for ExtensionDetail {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            extension: None,
            scroll_offset: 0,
        }
    }
}

impl ExtensionDetail {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn set_extension(&mut self, extension: Extension) {
        self.extension = Some(extension);
        self.scroll_offset = 0; // Reset scroll when setting new extension
    }

    fn scroll_up(&mut self) {
        if self.scroll_offset > 0 {
            self.scroll_offset = self.scroll_offset.saturating_sub(1);
        }
    }

    fn scroll_down(&mut self) {
        // We'll calculate max scroll based on content height in draw
        self.scroll_offset = self.scroll_offset.saturating_add(1);
    }
}

impl Component for ExtensionDetail {
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
            Action::ViewExtensionDetails(id) => {
                // In a real app, we'd fetch the extension by ID
                // For now, we'll use mock data
                if let Some(ext) = Extension::mock_extensions()
                    .into_iter()
                    .find(|e| e.id == id)
                {
                    self.set_extension(ext);
                }
            }
            _ => {}
        }
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        let Some(extension) = &self.extension else {
            // Show empty state
            let block = Block::default()
                .title(" Extension Details ")
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded);
            
            let text = Paragraph::new("No extension selected")
                .alignment(Alignment::Center)
                .block(block);
            
            frame.render_widget(text, area);
            return Ok(());
        };

        // Create layout
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Length(area.height.saturating_sub(3)), // Main content
                Constraint::Length(3),                               // Help bar
            ])
            .split(area);

        // Main content block
        let block = Block::default()
            .title(format!(" {} v{} ", extension.name, extension.version))
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(Color::Cyan));

        let inner_area = block.inner(chunks[0]);

        // Build content
        let mut content = vec![];

        // Description
        if let Some(desc) = &extension.description {
            content.push(Line::from(vec![
                Span::styled("Description: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::raw(desc),
            ]));
            content.push(Line::from(""));
        }

        // ID
        content.push(Line::from(vec![
            Span::styled("ID: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
            Span::raw(&extension.id),
        ]));
        content.push(Line::from(""));

        // Tags
        if !extension.metadata.tags.is_empty() {
            content.push(Line::from(vec![
                Span::styled("Tags: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::styled(
                    extension.metadata.tags.join(", "),
                    Style::default().fg(Color::Blue),
                ),
            ]));
            content.push(Line::from(""));
        }

        // Import date
        content.push(Line::from(vec![
            Span::styled("Imported: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
            Span::raw(extension.metadata.imported_at.format("%Y-%m-%d %H:%M:%S").to_string()),
        ]));
        
        // Source path
        if let Some(path) = &extension.metadata.source_path {
            content.push(Line::from(vec![
                Span::styled("Source: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::raw(path),
            ]));
        }
        content.push(Line::from(""));

        // MCP Servers section
        if !extension.mcp_servers.is_empty() {
            content.push(Line::from(Span::styled(
                "MCP Servers",
                Style::default().fg(Color::Magenta).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
            )));
            content.push(Line::from(""));

            for (name, config) in &extension.mcp_servers {
                content.push(Line::from(vec![
                    Span::raw("  "),
                    Span::styled(format!("• {}", name), Style::default().fg(Color::Green)),
                ]));

                // Server type
                if let Some(url) = &config.url {
                    content.push(Line::from(vec![
                        Span::raw("    Type: "),
                        Span::styled("URL", Style::default().fg(Color::Blue)),
                        Span::raw(" - "),
                        Span::raw(url),
                    ]));
                } else if let Some(cmd) = &config.command {
                    content.push(Line::from(vec![
                        Span::raw("    Type: "),
                        Span::styled("Command", Style::default().fg(Color::Blue)),
                        Span::raw(" - "),
                        Span::raw(cmd),
                    ]));
                    
                    if let Some(args) = &config.args {
                        content.push(Line::from(vec![
                            Span::raw("    Args: "),
                            Span::raw(args.join(" ")),
                        ]));
                    }
                }

                // Environment variables
                if let Some(env) = &config.env {
                    for (key, value) in env {
                        content.push(Line::from(vec![
                            Span::raw("    Env: "),
                            Span::styled(key, Style::default().fg(Color::Yellow)),
                            Span::raw(" = "),
                            Span::raw(value),
                        ]));
                    }
                }

                // Trust status
                if let Some(trust) = config.trust {
                    content.push(Line::from(vec![
                        Span::raw("    Trust: "),
                        Span::styled(
                            if trust { "Yes" } else { "No" },
                            Style::default().fg(if trust { Color::Green } else { Color::Red }),
                        ),
                    ]));
                }

                content.push(Line::from(""));
            }
        }

        // Context file section
        if let Some(filename) = &extension.context_file_name {
            content.push(Line::from(Span::styled(
                format!("Context File: {}", filename),
                Style::default().fg(Color::Magenta).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
            )));
            content.push(Line::from(""));
            
            if let Some(content_text) = &extension.context_content {
                // Add context file content with proper indentation
                for line in content_text.lines() {
                    content.push(Line::from(vec![
                        Span::raw("  "),
                        Span::raw(line),
                    ]));
                }
            }
        }

        // Create scrollable paragraph
        let paragraph = Paragraph::new(content)
            .scroll((self.scroll_offset, 0));

        // Render main content
        frame.render_widget(block, chunks[0]);
        frame.render_widget(paragraph, inner_area);

        // Help bar
        let help_text = " ↑/↓: Scroll | b: Back | e: Edit | d: Delete | q: Quit ";
        let help_bar = Paragraph::new(help_text)
            .style(Style::default().fg(Color::DarkGray))
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded),
            );

        frame.render_widget(help_bar, chunks[1]);

        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::KeyCode;

        match event {
            Some(crate::tui::Event::Key(key)) => match key.code {
                KeyCode::Up | KeyCode::Char('k') => {
                    self.scroll_up();
                    Ok(Some(Action::Render))
                }
                KeyCode::Down | KeyCode::Char('j') => {
                    self.scroll_down();
                    Ok(Some(Action::Render))
                }
                KeyCode::Char('b') | KeyCode::Esc => Ok(Some(Action::NavigateBack)),
                KeyCode::Char('e') => {
                    if let Some(ext) = &self.extension {
                        Ok(Some(Action::EditExtension(ext.id.clone())))
                    } else {
                        Ok(None)
                    }
                }
                KeyCode::Char('d') => {
                    if let Some(ext) = &self.extension {
                        Ok(Some(Action::DeleteExtension(ext.id.clone())))
                    } else {
                        Ok(None)
                    }
                }
                KeyCode::Char('q') => Ok(Some(Action::Quit)),
                _ => Ok(None),
            },
            _ => Ok(None),
        }
    }
}