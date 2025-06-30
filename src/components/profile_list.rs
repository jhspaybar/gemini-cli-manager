use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, models::Profile};

pub struct ProfileList {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    profiles: Vec<Profile>,
    selected: usize,
}

impl Default for ProfileList {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            profiles: Profile::mock_profiles(),
            selected: 0,
        }
    }
}

impl ProfileList {
    pub fn new() -> Self {
        Self::default()
    }

    fn next(&mut self) {
        if !self.profiles.is_empty() {
            self.selected = (self.selected + 1) % self.profiles.len();
        }
    }

    fn previous(&mut self) {
        if !self.profiles.is_empty() {
            if self.selected > 0 {
                self.selected -= 1;
            } else {
                self.selected = self.profiles.len() - 1;
            }
        }
    }

    fn get_selected_profile(&self) -> Option<&Profile> {
        self.profiles.get(self.selected)
    }
}

impl Component for ProfileList {
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
        // Create a block for the profile list
        let block = Block::default()
            .title(" Profiles ")
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(Color::Green));

        // Create list items
        let items: Vec<ListItem> = self
            .profiles
            .iter()
            .enumerate()
            .map(|(i, profile)| {
                let is_selected = i == self.selected;
                let is_default = profile.metadata.is_default;

                // Build the display string
                let mut lines = vec![
                    Line::from(vec![
                        Span::styled(
                            profile.display_name(),
                            if is_selected {
                                Style::default()
                                    .fg(Color::Yellow)
                                    .add_modifier(Modifier::BOLD)
                            } else {
                                Style::default().fg(Color::White)
                            },
                        ),
                        if is_default {
                            Span::styled(" (default)", Style::default().fg(Color::Blue))
                        } else {
                            Span::raw("")
                        },
                    ]),
                ];

                // Add description
                if let Some(desc) = &profile.description {
                    lines.push(Line::from(vec![
                        Span::raw("  "),
                        Span::styled(desc, Style::default().fg(Color::Gray)),
                    ]));
                }

                // Add summary
                lines.push(Line::from(vec![
                    Span::raw("  "),
                    Span::styled(
                        profile.summary(),
                        Style::default().fg(Color::DarkGray),
                    ),
                    Span::raw(" | "),
                    Span::styled(
                        format!("{} tags", profile.metadata.tags.len()),
                        Style::default().fg(Color::Magenta),
                    ),
                ]));

                // Add working directory if specified
                if let Some(dir) = &profile.working_directory {
                    lines.push(Line::from(vec![
                        Span::raw("  "),
                        Span::styled("📂 ", Style::default().fg(Color::Cyan)),
                        Span::styled(dir, Style::default().fg(Color::Cyan)),
                    ]));
                }

                lines.push(Line::from("")); // Empty line for spacing

                ListItem::new(lines)
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
            let help_text = " ↑/↓: Navigate | Enter: Launch | e: Edit | n: New | d: Delete | x: Set Default | q: Quit ";
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
                help_area,
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
                    if let Some(profile) = self.get_selected_profile() {
                        Ok(Some(Action::LaunchWithProfile(profile.id.clone())))
                    } else {
                        Ok(None)
                    }
                }
                KeyCode::Char('e') => {
                    if let Some(profile) = self.get_selected_profile() {
                        Ok(Some(Action::EditProfile(profile.id.clone())))
                    } else {
                        Ok(None)
                    }
                }
                KeyCode::Char('n') => Ok(Some(Action::CreateProfile)),
                KeyCode::Char('d') => {
                    if let Some(profile) = self.get_selected_profile() {
                        Ok(Some(Action::DeleteProfile(profile.id.clone())))
                    } else {
                        Ok(None)
                    }
                }
                KeyCode::Char('x') => {
                    // TODO: Implement set default
                    Ok(None)
                }
                KeyCode::Char('q') => Ok(Some(Action::Quit)),
                KeyCode::Tab => Ok(Some(Action::NavigateToExtensions)),
                _ => Ok(None),
            },
            _ => Ok(None),
        }
    }
}