use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, models::{Extension, Profile}};

pub struct ProfileDetail {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    profile: Option<Profile>,
    extensions: Vec<Extension>, // Full extension data for display
    scroll_offset: u16,
}

impl Default for ProfileDetail {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            profile: None,
            extensions: Vec::new(),
            scroll_offset: 0,
        }
    }
}

impl ProfileDetail {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn set_profile(&mut self, profile: Profile) {
        // Load the extensions that are part of this profile
        // In a real app, we'd fetch these from storage
        let all_extensions = Extension::mock_extensions();
        self.extensions = all_extensions
            .into_iter()
            .filter(|ext| profile.extension_ids.contains(&ext.id))
            .collect();
        
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
            Action::EditProfile(id) => {
                // In a real app, we'd fetch the profile by ID
                if let Some(profile) = Profile::mock_profiles()
                    .into_iter()
                    .find(|p| p.id == id)
                {
                    self.set_profile(profile);
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
            .border_style(Style::default().fg(Color::Green));

        let inner_area = block.inner(chunks[0]);

        // Build content
        let mut content = vec![];

        // Description
        if let Some(desc) = &profile.description {
            content.push(Line::from(vec![
                Span::styled("Description: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::raw(desc),
            ]));
            content.push(Line::from(""));
        }

        // ID
        content.push(Line::from(vec![
            Span::styled("ID: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
            Span::raw(&profile.id),
        ]));

        // Default status
        if profile.metadata.is_default {
            content.push(Line::from(vec![
                Span::styled("Status: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::styled("Default Profile", Style::default().fg(Color::Blue).add_modifier(Modifier::BOLD)),
            ]));
        }

        // Tags
        if !profile.metadata.tags.is_empty() {
            content.push(Line::from(vec![
                Span::styled("Tags: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::styled(
                    profile.metadata.tags.join(", "),
                    Style::default().fg(Color::Magenta),
                ),
            ]));
        }

        // Working directory
        if let Some(dir) = &profile.working_directory {
            content.push(Line::from(vec![
                Span::styled("Working Directory: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
                Span::raw(dir),
            ]));
        }

        // Creation date
        content.push(Line::from(vec![
            Span::styled("Created: ", Style::default().fg(Color::Yellow).add_modifier(Modifier::BOLD)),
            Span::raw(profile.metadata.created_at.format("%Y-%m-%d %H:%M:%S").to_string()),
        ]));
        content.push(Line::from(""));

        // Extensions section
        content.push(Line::from(Span::styled(
            "Extensions",
            Style::default().fg(Color::Cyan).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
        )));
        content.push(Line::from(""));

        if self.extensions.is_empty() {
            content.push(Line::from("  No extensions included"));
        } else {
            for ext in &self.extensions {
                content.push(Line::from(vec![
                    Span::raw("  "),
                    Span::styled(format!("• {}", ext.name), Style::default().fg(Color::Green)),
                    Span::raw(" "),
                    Span::styled(format!("v{}", ext.version), Style::default().fg(Color::DarkGray)),
                ]));
                
                if let Some(desc) = &ext.description {
                    content.push(Line::from(vec![
                        Span::raw("    "),
                        Span::styled(desc, Style::default().fg(Color::Gray)),
                    ]));
                }
                
                // Show MCP servers
                if !ext.mcp_servers.is_empty() {
                    content.push(Line::from(vec![
                        Span::raw("    "),
                        Span::styled(
                            format!("MCP Servers: {}", ext.mcp_servers.keys().cloned().collect::<Vec<_>>().join(", ")),
                            Style::default().fg(Color::DarkGray),
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
                Style::default().fg(Color::Cyan).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
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
                    Span::raw("  "),
                    Span::styled(key, Style::default().fg(Color::Yellow)),
                    Span::raw(" = "),
                    Span::raw(display_value),
                ]));
            }
            content.push(Line::from(""));
        }

        // Summary
        content.push(Line::from(Span::styled(
            "Summary",
            Style::default().fg(Color::Cyan).add_modifier(Modifier::BOLD | Modifier::UNDERLINED),
        )));
        content.push(Line::from(""));
        content.push(Line::from(format!("  • {} extensions", self.extensions.len())));
        content.push(Line::from(format!("  • {} environment variables", profile.environment_variables.len())));
        
        let total_mcp_servers: usize = self.extensions.iter()
            .map(|ext| ext.mcp_servers.len())
            .sum();
        content.push(Line::from(format!("  • {} total MCP servers", total_mcp_servers)));

        // Create scrollable paragraph
        let paragraph = Paragraph::new(content)
            .scroll((self.scroll_offset, 0));

        // Render main content
        frame.render_widget(block, chunks[0]);
        frame.render_widget(paragraph, inner_area);

        // Help bar
        let help_text = " ↑/↓: Scroll | Enter: Launch | e: Edit | b: Back | q: Quit ";
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