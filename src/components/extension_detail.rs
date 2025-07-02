use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, models::Extension, storage::Storage, theme};

#[derive(Default)]
pub struct ExtensionDetail {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    storage: Option<Storage>,
    extension: Option<Extension>,
    scroll_offset: u16,
}


impl ExtensionDetail {
    pub fn with_storage(storage: Storage) -> Self {
        Self {
            storage: Some(storage),
            ..Self::default()
        }
    }

    #[allow(dead_code)]
    pub fn new(storage: Storage, extension_id: String) -> Self {
        let mut detail = Self::with_storage(storage.clone());
        if let Ok(extension) = storage.load_extension(&extension_id) {
            detail.set_extension(extension);
        }
        detail
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
        if let Action::ViewExtensionDetails(id) = action {
            // Load the extension from storage
            if let Some(storage) = &self.storage {
                if let Ok(extension) = storage.load_extension(&id) {
                    self.set_extension(extension);
                }
            }
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
                Constraint::Length(3),                             // Help bar
            ])
            .split(area);

        // Main content block
        let block = Block::default()
            .title(format!(" {} v{} ", extension.name, extension.version))
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(theme::text_secondary()));

        let inner_area = block.inner(chunks[0]);

        // Build content
        let mut content = vec![];

        // Description
        if let Some(desc) = &extension.description {
            content.push(Line::from(vec![
                Span::styled(
                    "Description: ",
                    Style::default()
                        .fg(theme::highlight())
                        .add_modifier(Modifier::BOLD),
                ),
                Span::styled(desc, Style::default().fg(theme::text_primary())),
            ]));
            content.push(Line::from(""));
        }

        // ID
        content.push(Line::from(vec![
            Span::styled(
                "ID: ",
                Style::default()
                    .fg(theme::highlight())
                    .add_modifier(Modifier::BOLD),
            ),
            Span::styled(&extension.id, Style::default().fg(theme::text_primary())),
        ]));
        content.push(Line::from(""));

        // Tags
        if !extension.metadata.tags.is_empty() {
            content.push(Line::from(vec![
                Span::styled(
                    "Tags: ",
                    Style::default()
                        .fg(theme::highlight())
                        .add_modifier(Modifier::BOLD),
                ),
                Span::styled(
                    extension.metadata.tags.join(", "),
                    Style::default().fg(theme::primary()),
                ),
            ]));
            content.push(Line::from(""));
        }

        // Import date
        content.push(Line::from(vec![
            Span::styled(
                "Imported: ",
                Style::default()
                    .fg(theme::highlight())
                    .add_modifier(Modifier::BOLD),
            ),
            Span::styled(
                extension
                    .metadata
                    .imported_at
                    .format("%Y-%m-%d %H:%M:%S")
                    .to_string(),
                Style::default().fg(theme::text_primary()),
            ),
        ]));

        // Source path
        if let Some(path) = &extension.metadata.source_path {
            content.push(Line::from(vec![
                Span::styled(
                    "Source: ",
                    Style::default()
                        .fg(theme::highlight())
                        .add_modifier(Modifier::BOLD),
                ),
                Span::styled(path, Style::default().fg(theme::text_primary())),
            ]));
        }
        content.push(Line::from(""));

        // MCP Servers section
        if !extension.mcp_servers.is_empty() {
            content.push(Line::from(Span::styled(
                "MCP Servers",
                Style::default()
                    .fg(theme::accent())
                    .add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
            )));
            content.push(Line::from(""));

            for (name, config) in &extension.mcp_servers {
                content.push(Line::from(vec![
                    Span::styled("  ", Style::default().fg(theme::text_primary())),
                    Span::styled(format!("• {name}"), Style::default().fg(theme::success())),
                ]));

                // Server type - MCP servers can be URL-based or command-based
                if let Some(url) = &config.url {
                    content.push(Line::from(vec![
                        Span::styled("    Type: ", Style::default().fg(theme::text_secondary())),
                        Span::styled("URL (SSE)", Style::default().fg(theme::primary())),
                    ]));
                    content.push(Line::from(vec![
                        Span::styled("    URL: ", Style::default().fg(theme::text_secondary())),
                        Span::styled(url, Style::default().fg(theme::text_primary())),
                    ]));
                } else if let Some(cmd) = &config.command {
                    content.push(Line::from(vec![
                        Span::styled("    Type: ", Style::default().fg(theme::text_secondary())),
                        Span::styled("Command", Style::default().fg(theme::primary())),
                        Span::styled(" - ", Style::default().fg(theme::text_secondary())),
                        Span::styled(cmd, Style::default().fg(theme::text_primary())),
                    ]));

                    if let Some(args) = &config.args {
                        content.push(Line::from(vec![
                            Span::styled(
                                "    Args: ",
                                Style::default().fg(theme::text_secondary()),
                            ),
                            Span::styled(
                                args.join(" "),
                                Style::default().fg(theme::text_primary()),
                            ),
                        ]));
                    }
                }

                // Environment variables
                if let Some(env) = &config.env {
                    for (key, value) in env {
                        content.push(Line::from(vec![
                            Span::styled("    Env: ", Style::default().fg(theme::text_secondary())),
                            Span::styled(key, Style::default().fg(theme::highlight())),
                            Span::styled(" = ", Style::default().fg(theme::text_secondary())),
                            Span::styled(value, Style::default().fg(theme::text_primary())),
                        ]));
                    }
                }

                // Trust status
                if let Some(trust) = config.trust {
                    content.push(Line::from(vec![
                        Span::styled("    Trust: ", Style::default().fg(theme::text_secondary())),
                        Span::styled(
                            if trust { "Yes" } else { "No" },
                            Style::default().fg(if trust {
                                theme::success()
                            } else {
                                theme::error()
                            }),
                        ),
                    ]));
                }

                content.push(Line::from(""));
            }
        }

        // Context file section
        if let Some(content_text) = &extension.context_content {
            // Determine filename - use provided or default to GEMINI.md
            let filename = extension
                .context_file_name
                .as_ref()
                .filter(|s| !s.trim().is_empty())
                .map(|s| s.as_str())
                .unwrap_or("GEMINI.md");

            content.push(Line::from(Span::styled(
                format!("Context File: {filename}"),
                Style::default()
                    .fg(theme::accent())
                    .add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
            )));
            content.push(Line::from(""));

            // Add context file content with proper indentation
            for line in content_text.lines() {
                content.push(Line::from(vec![
                    Span::styled("  ", Style::default().fg(theme::text_primary())),
                    Span::styled(line, Style::default().fg(theme::text_primary())),
                ]));
            }
            content.push(Line::from(""));
        }

        // Create scrollable paragraph
        let paragraph = Paragraph::new(content).scroll((self.scroll_offset, 0));

        // Render main content
        frame.render_widget(block, chunks[0]);
        frame.render_widget(paragraph, inner_area);

        // Help bar
        use crate::utils::build_help_text;
        let help_text = build_help_text(&[
            ("up", "Scroll"),
            ("down", "Scroll"),
            ("back", "Back"),
            ("edit", "Edit"),
            ("delete", "Delete"),
            ("quit", "Quit"),
        ]);
        let help_bar = Paragraph::new(help_text)
            .style(Style::default().fg(theme::text_muted()))
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .border_style(Style::default().fg(theme::text_secondary())),
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

// Test helper methods
impl ExtensionDetail {
    /// Test helper method - returns current section
    #[doc(hidden)]
    #[allow(dead_code)]
    pub fn current_section(&self) -> usize {
        // For now, we don't have sections, but this could track scroll position
        0
    }
}
