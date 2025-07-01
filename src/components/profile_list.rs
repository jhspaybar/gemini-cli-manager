use std::sync::{Arc, RwLock};

use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;
use tui_input::Input;
use tui_input::backend::crossterm::EventHandler;

use super::Component;
use crate::{action::Action, config::Config, models::Profile, storage::Storage, theme};
use crate::components::settings_view::UserSettings;

pub struct ProfileList {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    profiles: Vec<Profile>,
    filtered_profiles: Vec<usize>, // Indices of profiles that match filter
    selected: usize,
    storage: Option<Storage>,
    search_mode: bool,
    search_input: Input,
    settings: Option<Arc<RwLock<UserSettings>>>,
}

impl Default for ProfileList {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            profiles: Vec::new(),
            filtered_profiles: Vec::new(),
            selected: 0,
            storage: None,
            search_mode: false,
            search_input: Input::default(),
            settings: None,
        }
    }
}

impl ProfileList {
    pub fn with_storage(storage: Storage) -> Self {
        let mut list = Self::default();
        list.storage = Some(storage.clone());
        
        // Load profiles from storage
        if let Ok(profiles) = storage.list_profiles() {
            list.profiles = profiles;
            list.update_filter();
        }
        
        list
    }

    fn update_filter(&mut self) {
        let search_query = self.search_input.value();
        if search_query.is_empty() {
            // Show all profiles
            self.filtered_profiles = (0..self.profiles.len()).collect();
        } else {
            // Filter profiles based on search query
            let query = search_query.to_lowercase();
            self.filtered_profiles = self.profiles
                .iter()
                .enumerate()
                .filter(|(_, profile)| {
                    profile.name.to_lowercase().contains(&query) ||
                    profile.description.as_ref().map_or(false, |d| d.to_lowercase().contains(&query)) ||
                    profile.metadata.tags.iter().any(|tag| tag.to_lowercase().contains(&query))
                })
                .map(|(i, _)| i)
                .collect();
        }
        
        // Adjust selection if needed
        if self.selected >= self.filtered_profiles.len() && !self.filtered_profiles.is_empty() {
            self.selected = self.filtered_profiles.len() - 1;
        } else if self.filtered_profiles.is_empty() {
            self.selected = 0;
        }
    }

    fn next(&mut self) {
        if !self.filtered_profiles.is_empty() {
            self.selected = (self.selected + 1) % self.filtered_profiles.len();
        }
    }

    fn previous(&mut self) {
        if !self.filtered_profiles.is_empty() {
            if self.selected > 0 {
                self.selected -= 1;
            } else {
                self.selected = self.filtered_profiles.len() - 1;
            }
        }
    }

    fn get_selected_profile(&self) -> Option<&Profile> {
        self.filtered_profiles.get(self.selected)
            .and_then(|&idx| self.profiles.get(idx))
    }
    
    // Public methods for testing
    #[allow(dead_code)]
    pub fn selected_index(&self) -> usize {
        self.selected
    }
    
    #[allow(dead_code)]
    pub fn is_search_mode(&self) -> bool {
        self.search_mode
    }
    
    #[allow(dead_code)]
    pub fn search_query(&self) -> &str {
        self.search_input.value()
    }
    
    #[allow(dead_code)]
    pub fn filtered_count(&self) -> usize {
        self.filtered_profiles.len()
    }
    
    #[allow(dead_code)]
    pub fn total_count(&self) -> usize {
        self.profiles.len()
    }
    
    fn check_keybinding(&self, key: &crossterm::event::KeyEvent, action: &str) -> bool {
        if let Some(settings) = &self.settings {
            if let Ok(settings_lock) = settings.read() {
                let configured_keys = settings_lock.keybindings.get_keys_for_action(action);
                let key_str = format_key_event(key);
                return configured_keys.contains(&key_str);
            }
        }
        false
    }
}

// Helper function to format key events
fn format_key_event(key: &crossterm::event::KeyEvent) -> String {
    use crossterm::event::{KeyCode, KeyModifiers};
    
    let mut parts = Vec::new();
    
    // Add modifiers
    if key.modifiers.contains(KeyModifiers::CONTROL) {
        parts.push("Ctrl");
    }
    if key.modifiers.contains(KeyModifiers::ALT) {
        parts.push("Alt");
    }
    if key.modifiers.contains(KeyModifiers::SHIFT) {
        parts.push("Shift");
    }
    
    // Add the key itself
    let key_str = match key.code {
        KeyCode::Char(c) => {
            match c {
                ' ' => "Space".to_string(),
                c => {
                    if key.modifiers.contains(KeyModifiers::SHIFT) {
                        c.to_uppercase().to_string()
                    } else {
                        c.to_string()
                    }
                }
            }
        }
        KeyCode::F(n) => format!("F{}", n),
        KeyCode::Up => "Up".to_string(),
        KeyCode::Down => "Down".to_string(),
        KeyCode::Left => "Left".to_string(),
        KeyCode::Right => "Right".to_string(),
        KeyCode::Enter => "Enter".to_string(),
        KeyCode::Tab => "Tab".to_string(),
        KeyCode::BackTab => "BackTab".to_string(),
        KeyCode::Backspace => "Backspace".to_string(),
        KeyCode::Delete => "Delete".to_string(),
        KeyCode::Insert => "Insert".to_string(),
        KeyCode::Home => "Home".to_string(),
        KeyCode::End => "End".to_string(),
        KeyCode::PageUp => "PageUp".to_string(),
        KeyCode::PageDown => "PageDown".to_string(),
        KeyCode::Esc => "Esc".to_string(),
        _ => return "Unknown".to_string(),
    };
    
    parts.push(&key_str);
    
    if parts.len() > 1 && !parts[0].starts_with('F') {
        parts.join("+")
    } else {
        key_str
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
    
    fn register_settings_handler(&mut self, settings: Arc<RwLock<UserSettings>>) -> Result<()> {
        self.settings = Some(settings);
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
            Action::RefreshProfiles => {
                // Reload profiles from storage
                if let Some(storage) = &self.storage {
                    if let Ok(profiles) = storage.list_profiles() {
                        self.profiles = profiles;
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
                .border_style(Style::default().fg(theme::highlight()));
            
            // Use tui-input's widget with proper styling
            let input_widget = Paragraph::new(self.search_input.value())
                .style(Style::default().fg(theme::text_primary()))
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
        
        // Create a block for the profile list
        let title = if self.search_mode && !self.search_input.value().is_empty() {
            format!(" Profiles ({}/{}) ", self.filtered_profiles.len(), self.profiles.len())
        } else {
            " Profiles ".to_string()
        };
        
        let block = Block::default()
            .title(title)
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(theme::text_secondary()));

        // Create list items
        let items: Vec<ListItem> = self
            .filtered_profiles
            .iter()
            .enumerate()
            .filter_map(|(i, &profile_idx)| {
                self.profiles.get(profile_idx).map(|profile| {
                    let is_selected = i == self.selected;
                let is_default = profile.metadata.is_default;

                // Build the display string
                let mut lines = vec![
                    Line::from(vec![
                        Span::styled(
                            profile.display_name(),
                            if is_selected {
                                Style::default()
                                    .fg(theme::text_primary())
                                    .add_modifier(Modifier::BOLD)
                            } else {
                                Style::default().fg(theme::text_primary())
                            },
                        ),
                        if is_default {
                            Span::styled(" (default)", Style::default().fg(theme::primary()))
                        } else {
                            Span::styled("", Style::default().fg(theme::text_primary()))
                        },
                    ]),
                ];

                // Add description
                if let Some(desc) = &profile.description {
                    lines.push(Line::from(vec![
                        Span::styled("  ", Style::default().fg(theme::text_primary())),
                        Span::styled(desc, Style::default().fg(theme::text_secondary())),
                    ]));
                }

                // Add summary
                lines.push(Line::from(vec![
                    Span::styled("  ", Style::default().fg(theme::text_primary())),
                    Span::styled(
                        profile.summary(),
                        Style::default().fg(theme::text_muted()),
                    ),
                    Span::styled(" | ", Style::default().fg(theme::text_secondary())),
                    Span::styled(
                        format!("{} tags", profile.metadata.tags.len()),
                        Style::default().fg(theme::accent()),
                    ),
                ]));

                // Add working directory if specified
                if let Some(dir) = &profile.working_directory {
                    lines.push(Line::from(vec![
                        Span::styled("  ", Style::default().fg(theme::text_primary())),
                        Span::styled("📂 ", Style::default().fg(theme::info())),
                        Span::styled(dir, Style::default().fg(theme::info())),
                    ]));
                }

                lines.push(Line::from("")); // Empty line for spacing

                    ListItem::new(lines)
                })
            })
            .collect();

        // Check if list is empty
        if self.profiles.is_empty() || (self.search_mode && self.filtered_profiles.is_empty() && !self.search_input.value().is_empty()) {
            // Show empty state message
            let empty_msg = if self.search_mode && !self.search_input.value().is_empty() {
                vec![
                    "No profiles match your search",
                    "",
                    "Try a different search term"
                ]
            } else {
                vec![
                    "No profiles found",
                    "",
                    "Press 'n' to create your first profile"
                ]
            };
            
            let empty_widget = Paragraph::new(empty_msg.join("\n"))
                .style(Style::default().fg(theme::text_secondary()))
                .alignment(Alignment::Center)
                .block(block);
            
            frame.render_widget(empty_widget, list_area);
        } else {
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
        }

        // Add help text at the bottom
        if list_area.height > 4 {
            let help_text = if self.search_mode {
                // Use dynamic help text for search mode too
                use crate::utils::build_help_text;
                build_help_text(&[
                    ("Type", "Search"),
                    ("back", "Close search"),
                    ("up", "Navigate results"),
                ])
            } else {
                // Use dynamic help text based on current keybindings
                use crate::utils::build_help_text;
                build_help_text(&[
                    ("up", "Navigate"),
                    ("down", "Navigate"),
                    ("select", "View"),
                    ("edit", "Edit"),
                    ("launch", "Launch"),
                    ("create", "New"),
                    ("delete", "Delete"),
                    ("search", "Search"),
                    ("tab", "Extensions"),
                    ("quit", "Quit"),
                ])
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
                    if self.check_keybinding(&key, "back") {
                        self.search_mode = false;
                        self.search_input.reset();
                        self.update_filter();
                        return Ok(Some(Action::Render));
                    }
                    
                    if self.check_keybinding(&key, "up") {
                        self.previous();
                        return Ok(Some(Action::Render));
                    }
                    
                    if self.check_keybinding(&key, "down") {
                        self.next();
                        return Ok(Some(Action::Render));
                    }
                    
                    if self.check_keybinding(&key, "select") {
                        if let Some(profile) = self.get_selected_profile() {
                            return Ok(Some(Action::ViewProfileDetails(profile.id.clone())));
                        }
                        return Ok(None);
                    }
                    
                    match key.code {
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
                    // Normal mode - use configured keybindings
                    if self.check_keybinding(&key, "up") {
                        self.previous();
                        return Ok(Some(Action::Render));
                    }
                    
                    if self.check_keybinding(&key, "down") {
                        self.next();
                        return Ok(Some(Action::Render));
                    }
                    
                    if self.check_keybinding(&key, "select") {
                        if let Some(profile) = self.get_selected_profile() {
                            return Ok(Some(Action::ViewProfileDetails(profile.id.clone())));
                        }
                        return Ok(None);
                    }
                    
                    if self.check_keybinding(&key, "edit") {
                        if let Some(profile) = self.get_selected_profile() {
                            return Ok(Some(Action::EditProfile(profile.id.clone())));
                        }
                        return Ok(None);
                    }
                    
                    if self.check_keybinding(&key, "launch") {
                        if let Some(profile) = self.get_selected_profile() {
                            return Ok(Some(Action::LaunchWithProfile(profile.id.clone())));
                        }
                        return Ok(None);
                    }
                    
                    if self.check_keybinding(&key, "create") {
                        return Ok(Some(Action::CreateProfile));
                    }
                    
                    if self.check_keybinding(&key, "delete") {
                        if let Some(profile) = self.get_selected_profile() {
                            return Ok(Some(Action::DeleteProfile(profile.id.clone())));
                        }
                        return Ok(None);
                    }
                    if self.check_keybinding(&key, "search") {
                        // Enter search mode
                        self.search_mode = true;
                        self.search_input.reset();
                        return Ok(Some(Action::Render));
                    }
                    
                    if self.check_keybinding(&key, "quit") {
                        return Ok(Some(Action::Quit));
                    }
                    
                    // Handle special keys that aren't customizable
                    match key.code {
                        KeyCode::Char('x') => {
                            if let Some(profile) = self.get_selected_profile() {
                                let profile_id = profile.id.clone();
                                // Set this profile as default
                                for p in &mut self.profiles {
                                    p.metadata.is_default = p.id == profile_id;
                                }
                                // Save the updated profiles
                                if let Some(storage) = &self.storage {
                                    for p in &self.profiles {
                                        let _ = storage.save_profile(p);
                                    }
                                }
                                Ok(Some(Action::Render))
                            } else {
                                Ok(None)
                            }
                        }
                        KeyCode::Tab => Ok(Some(Action::NavigateToSettings)),
                        _ => Ok(None),
                    }
                }
            }
            _ => Ok(None),
        }
    }
}