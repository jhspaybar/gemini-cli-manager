use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, models::{Extension, Profile}, storage::Storage, theme};

pub struct ProfileDetail {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    storage: Option<Storage>,
    profile: Option<Profile>,
    extensions: Vec<Extension>, // Full extension data for display
    scroll_offset: u16,
}

impl Default for ProfileDetail {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            storage: None,
            profile: None,
            extensions: Vec::new(),
            scroll_offset: 0,
        }
    }
}

impl ProfileDetail {
    pub fn with_storage(storage: Storage) -> Self {
        let mut detail = Self::default();
        detail.storage = Some(storage);
        detail
    }

    pub fn set_profile(&mut self, profile: Profile) {
        // Load the extensions from storage
        if let Some(storage) = &self.storage {
            self.extensions = profile.extension_ids
                .iter()
                .filter_map(|ext_id| storage.load_extension(ext_id).ok())
                .collect();
        } else {
            self.extensions = Vec::new();
        }
        
        self.profile = Some(profile);
        self.scroll_offset = 0;
    }

    fn scroll_up(&mut self) {
        if self.scroll_offset > 0 {
            self.scroll_offset = self.scroll_offset.saturating_sub(1);
        }
    }

    fn scroll_down(&mut self) {
        self.scroll_offset = self.scroll_offset.saturating_add(1);
    }
}

impl Component for ProfileDetail {
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
            Action::ViewProfileDetails(id) => {
                // Load the profile from storage
                if let Some(storage) = &self.storage {
                    if let Ok(profile) = storage.load_profile(&id) {
                        self.set_profile(profile);
                    }
                }
            }
            Action::RefreshProfiles => {
                // Reload the current profile if we have one
                if let Some(current_profile) = &self.profile {
                    let profile_id = current_profile.id.clone();
                    if let Some(storage) = &self.storage {
                        if let Ok(profile) = storage.load_profile(&profile_id) {
                            self.set_profile(profile);
                        }
                    }
                }
            }
            _ => {}
        }
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        let Some(profile) = &self.profile else {
            // Show empty state
            let block = Block::default()
                .title(" Profile Details ")
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded);
            
            let text = Paragraph::new("No profile selected")
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
            .title(format!(" {} ", profile.display_name()))
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(theme::success()));

        let inner_area = block.inner(chunks[0]);

        // Build content
        let mut content = vec![];

        // Description
        if let Some(desc) = &profile.description {
            content.push(Line::from(vec![
                Span::styled("Description: ", Style::default().fg(theme::highlight()).add_modifier(Modifier::BOLD)),
                Span::styled(desc, Style::default().fg(theme::text_primary())),
            ]));
            content.push(Line::from(""));
        }

        // ID
        content.push(Line::from(vec![
            Span::styled("ID: ", Style::default().fg(theme::highlight()).add_modifier(Modifier::BOLD)),
            Span::styled(&profile.id, Style::default().fg(theme::text_primary())),
        ]));

        // Default status
        if profile.metadata.is_default {
            content.push(Line::from(vec![
                Span::styled("Status: ", Style::default().fg(theme::highlight()).add_modifier(Modifier::BOLD)),
                Span::styled("Default Profile", Style::default().fg(theme::primary()).add_modifier(Modifier::BOLD)),
            ]));
        }

        // Tags
        if !profile.metadata.tags.is_empty() {
            content.push(Line::from(vec![
                Span::styled("Tags: ", Style::default().fg(theme::highlight()).add_modifier(Modifier::BOLD)),
                Span::styled(
                    profile.metadata.tags.join(", "),
                    Style::default().fg(theme::accent()),
                ),
            ]));
        }

        // Working directory
        if let Some(dir) = &profile.working_directory {
            content.push(Line::from(vec![
                Span::styled("Working Directory: ", Style::default().fg(theme::highlight()).add_modifier(Modifier::BOLD)),
                Span::styled(dir, Style::default().fg(theme::text_primary())),
            ]));
        }

        // Creation date
        content.push(Line::from(vec![
            Span::styled("Created: ", Style::default().fg(theme::highlight()).add_modifier(Modifier::BOLD)),
            Span::styled(profile.metadata.created_at.format("%Y-%m-%d %H:%M:%S").to_string(), Style::default().fg(theme::text_primary())),
        ]));
        content.push(Line::from(""));

        // Extensions section
        content.push(Line::from(Span::styled(
            "Extensions",
            Style::default().fg(theme::info()).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
        )));
        content.push(Line::from(""));

        if self.extensions.is_empty() {
            content.push(Line::from(Span::styled("  No extensions included", Style::default().fg(theme::text_muted()))));
        } else {
            for ext in &self.extensions {
                content.push(Line::from(vec![
                    Span::styled("  ", Style::default().fg(theme::text_primary())),
                    Span::styled(format!("• {}", ext.name), Style::default().fg(theme::success())),
                    Span::styled(" ", Style::default().fg(theme::text_primary())),
                    Span::styled(format!("v{}", ext.version), Style::default().fg(theme::text_muted())),
                ]));
                
                if let Some(desc) = &ext.description {
                    content.push(Line::from(vec![
                        Span::styled("    ", Style::default().fg(theme::text_primary())),
                        Span::styled(desc, Style::default().fg(theme::text_secondary())),
                    ]));
                }
                
                // Show MCP servers
                if !ext.mcp_servers.is_empty() {
                    content.push(Line::from(vec![
                        Span::styled("    ", Style::default().fg(theme::text_primary())),
                        Span::styled(
                            format!("MCP Servers: {}", ext.mcp_servers.keys().cloned().collect::<Vec<_>>().join(", ")),
                            Style::default().fg(theme::text_muted()),
                        ),
                    ]));
                }
                
                content.push(Line::from(""));
            }
        }

        // Environment Variables section
        if !profile.environment_variables.is_empty() {
            content.push(Line::from(Span::styled(
                "Environment Variables",
                Style::default().fg(theme::info()).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
            )));
            content.push(Line::from(""));

            for (key, value) in &profile.environment_variables {
                // Mask sensitive values
                let display_value = if key.contains("TOKEN") || key.contains("KEY") || key.contains("SECRET") {
                    let len = value.len();
                    if len > 8 {
                        format!("{}...{}", &value[..4], &value[len-4..])
                    } else {
                        "***".to_string()
                    }
                } else {
                    value.clone()
                };
                
                content.push(Line::from(vec![
                    Span::styled("  ", Style::default().fg(theme::text_primary())),
                    Span::styled(key, Style::default().fg(theme::highlight())),
                    Span::styled(" = ", Style::default().fg(theme::text_secondary())),
                    Span::styled(display_value, Style::default().fg(theme::text_primary())),
                ]));
            }
            content.push(Line::from(""));
        }

        // Summary
        content.push(Line::from(Span::styled(
            "Summary",
            Style::default().fg(theme::info()).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
        )));
        content.push(Line::from(""));
        content.push(Line::from(Span::styled(format!("  • {} extensions", self.extensions.len()), Style::default().fg(theme::text_primary()))));
        content.push(Line::from(Span::styled(format!("  • {} environment variables", profile.environment_variables.len()), Style::default().fg(theme::text_primary()))));
        
        let total_mcp_servers: usize = self.extensions.iter()
            .map(|ext| ext.mcp_servers.len())
            .sum();
        content.push(Line::from(Span::styled(format!("  • {} MCP servers total", total_mcp_servers), Style::default().fg(theme::text_primary()))));

        // Create scrollable paragraph
        let paragraph = Paragraph::new(content)
            .scroll((self.scroll_offset, 0));

        // Render main content
        frame.render_widget(block, chunks[0]);
        frame.render_widget(paragraph, inner_area);

        // Help bar
        let help_text = " ↑/↓: Scroll | Enter: Launch | e: Edit | b: Back | q: Quit ";
        let help_bar = Paragraph::new(help_text)
            .style(Style::default().fg(theme::text_muted()))
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
                KeyCode::Enter => {
                    if let Some(profile) = &self.profile {
                        Ok(Some(Action::LaunchWithProfile(profile.id.clone())))
                    } else {
                        Ok(None)
                    }
                }
                KeyCode::Char('b') | KeyCode::Esc => Ok(Some(Action::NavigateBack)),
                KeyCode::Char('e') => {
                    if let Some(profile) = &self.profile {
                        Ok(Some(Action::EditProfile(profile.id.clone())))
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