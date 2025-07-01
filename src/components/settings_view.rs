use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, theme, utils::KeybindingManager};

// NOTE: There's a complex module import resolution issue with the settings module
// The settings module compiles fine on its own, but importing from it causes circular 
// reference issues. For now using inline definitions until the issue is investigated.
// TODO: Investigate Rust module import resolution issue in components/settings_view.rs

#[derive(Debug, Clone)]
pub struct SettingsManager {
    settings: UserSettings,
    settings_path: std::path::PathBuf,
}

impl SettingsManager {
    pub fn new() -> color_eyre::Result<Self> {
        let settings_path = Self::get_settings_path()?;
        let settings = Self::load_settings(&settings_path)?;
        
        Ok(Self {
            settings,
            settings_path,
        })
    }
    
    pub fn get_settings(&self) -> &UserSettings {
        &self.settings
    }
    
    pub fn update_theme(&mut self, theme: String) -> color_eyre::Result<()> {
        self.settings.theme = theme;
        self.save()
    }
    
    pub fn reset_keybindings(&mut self) -> color_eyre::Result<()> {
        self.settings.keybindings = KeybindingConfig::default();
        self.save()
    }
    
    pub fn update_keybinding(&mut self, action: &str, keys: Vec<String>) -> color_eyre::Result<()> {
        match action {
            "up" => self.settings.keybindings.navigation.up = keys,
            "down" => self.settings.keybindings.navigation.down = keys,
            "left" => self.settings.keybindings.navigation.left = keys,
            "right" => self.settings.keybindings.navigation.right = keys,
            "back" => self.settings.keybindings.navigation.back = keys,
            "quit" => self.settings.keybindings.navigation.quit = keys,
            "edit" => self.settings.keybindings.actions.edit = keys,
            "delete" => self.settings.keybindings.actions.delete = keys,
            "create" => self.settings.keybindings.actions.create = keys,
            "import" => self.settings.keybindings.actions.import = keys,
            "launch" => self.settings.keybindings.actions.launch = keys,
            "select" => self.settings.keybindings.actions.select = keys,
            "search" => self.settings.keybindings.actions.search = keys,
            _ => return Err(color_eyre::eyre::eyre!("Unknown action: {}", action)),
        }
        self.save()
    }
    
    fn save(&self) -> color_eyre::Result<()> {
        let json = serde_json::to_string_pretty(&self.settings)?;
        if let Some(parent) = self.settings_path.parent() {
            std::fs::create_dir_all(parent)?;
        }
        std::fs::write(&self.settings_path, json)?;
        Ok(())
    }
    
    fn load_settings(path: &std::path::PathBuf) -> color_eyre::Result<UserSettings> {
        if path.exists() {
            let content = std::fs::read_to_string(path)?;
            let settings: UserSettings = serde_json::from_str(&content)?;
            Ok(settings)
        } else {
            Ok(UserSettings::default())
        }
    }
    
    fn get_settings_path() -> color_eyre::Result<std::path::PathBuf> {
        // Use data_local_dir for application data on all platforms
        // On macOS: ~/Library/Application Support
        // On Linux: ~/.local/share
        // On Windows: C:\Users\{username}\AppData\Local
        let data_dir = dirs::data_local_dir()
            .ok_or_else(|| color_eyre::eyre::eyre!("Could not find data directory"))?;
        Ok(data_dir.join("gemini-cli-manager").join("settings.json"))
    }
    
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct UserSettings {
    pub theme: String,
    pub keybindings: KeybindingConfig,
}

impl Default for UserSettings {
    fn default() -> Self {
        Self {
            theme: "mocha".to_string(),
            keybindings: KeybindingConfig::default(),
        }
    }
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct KeybindingConfig {
    pub navigation: NavigationKeys,
    pub actions: ActionKeys,
}

impl Default for KeybindingConfig {
    fn default() -> Self {
        Self {
            navigation: NavigationKeys::default(),
            actions: ActionKeys::default(),
        }
    }
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct NavigationKeys {
    pub up: Vec<String>,
    pub down: Vec<String>,
    pub left: Vec<String>,
    pub right: Vec<String>,
    pub back: Vec<String>,
    pub quit: Vec<String>,
}

impl Default for NavigationKeys {
    fn default() -> Self {
        Self {
            up: vec!["Up".to_string(), "k".to_string()],
            down: vec!["Down".to_string(), "j".to_string()],
            left: vec!["Left".to_string(), "h".to_string()],
            right: vec!["Right".to_string(), "l".to_string()],
            back: vec!["Esc".to_string(), "b".to_string()],
            quit: vec!["q".to_string(), "Ctrl+c".to_string()],
        }
    }
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct ActionKeys {
    pub edit: Vec<String>,
    pub delete: Vec<String>,
    pub create: Vec<String>,
    pub import: Vec<String>,
    pub launch: Vec<String>,
    pub select: Vec<String>,
    pub search: Vec<String>,
}

impl Default for ActionKeys {
    fn default() -> Self {
        Self {
            edit: vec!["e".to_string()],
            delete: vec!["d".to_string()],
            create: vec!["n".to_string()],
            import: vec!["i".to_string()],
            launch: vec!["l".to_string()],
            select: vec!["Enter".to_string(), "Space".to_string()],
            search: vec!["/".to_string()],
        }
    }
}

impl KeybindingConfig {
    pub fn get_keys_for_action(&self, action: &str) -> Vec<String> {
        // This will be replaced by KeybindingManager
        match action {
            "up" => self.navigation.up.clone(),
            "down" => self.navigation.down.clone(),
            "left" => self.navigation.left.clone(),
            "right" => self.navigation.right.clone(),
            "back" => self.navigation.back.clone(),
            "quit" => self.navigation.quit.clone(),
            "edit" => self.actions.edit.clone(),
            "delete" => self.actions.delete.clone(),
            "create" => self.actions.create.clone(),
            "import" => self.actions.import.clone(),
            "launch" => self.actions.launch.clone(),
            "select" => self.actions.select.clone(),
            "search" => self.actions.search.clone(),
            "tab" => vec!["Tab".to_string()], // Hardcoded for now
            "x" => vec!["x".to_string()], // Hardcoded for now
            "Space" => vec!["Space".to_string()], // Hardcoded for now
            "Ctrl+S" => vec!["Ctrl+S".to_string()], // Hardcoded for now
            "Type" => vec!["Type".to_string()], // Hardcoded for now - represents typing text
            "r" => vec!["r".to_string()], // Hardcoded for now - reset keybindings
            "e" => vec!["e".to_string()], // Hardcoded for now - export settings
            "i" => vec!["i".to_string()], // Hardcoded for now - import settings
            _ => vec![],
        }
    }
    
    /// Convert to HashMap for KeybindingManager
    #[allow(dead_code)]
    pub fn to_keybinding_config(&self) -> std::collections::HashMap<String, Vec<String>> {
        // TODO: Use KeyAction enum when keybindings module is available
        let mut config = std::collections::HashMap::new();
        
        config.insert("NavigateUp".to_string(), self.navigation.up.clone());
        config.insert("NavigateDown".to_string(), self.navigation.down.clone());
        config.insert("NavigateLeft".to_string(), self.navigation.left.clone());
        config.insert("NavigateRight".to_string(), self.navigation.right.clone());
        config.insert("NavigateBack".to_string(), self.navigation.back.clone());
        config.insert("Quit".to_string(), self.navigation.quit.clone());
        config.insert("Edit".to_string(), self.actions.edit.clone());
        config.insert("Delete".to_string(), self.actions.delete.clone());
        config.insert("Create".to_string(), self.actions.create.clone());
        config.insert("Import".to_string(), self.actions.import.clone());
        config.insert("LaunchProfile".to_string(), self.actions.launch.clone());
        config.insert("Select".to_string(), self.actions.select.clone());
        config.insert("StartSearch".to_string(), self.actions.search.clone());
        
        config
    }
}

#[derive(Debug, Clone)]
pub struct ThemeInfo {
    pub name: String,
    pub display_name: String,
    pub variant: String,
}

pub fn available_themes() -> Vec<ThemeInfo> {
    vec![
        ThemeInfo {
            name: "mocha".to_string(),
            display_name: "Mocha".to_string(),
            variant: "Dark".to_string(),
        },
        ThemeInfo {
            name: "macchiato".to_string(),
            display_name: "Macchiato".to_string(),
            variant: "Dark".to_string(),
        },
        ThemeInfo {
            name: "frappe".to_string(),
            display_name: "Frappé".to_string(),
            variant: "Dark".to_string(),
        },
        ThemeInfo {
            name: "latte".to_string(),
            display_name: "Latte".to_string(),
            variant: "Light".to_string(),
        },
    ]
}

#[derive(Debug, PartialEq)]
enum SettingsSection {
    Appearance,
    Keybindings,
}

#[derive(Debug, PartialEq)]
enum FocusedPane {
    Sections,
    Content,
}

pub struct Settings {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    settings_manager: Option<SettingsManager>,
    shared_settings: Option<std::sync::Arc<std::sync::RwLock<UserSettings>>>,
    keybinding_manager: Option<KeybindingManager>,
    
    // UI state
    current_section: SettingsSection,
    focused_pane: FocusedPane,
    selected_theme: usize,
    selected_keybinding: usize,
    editing_keybinding: bool,
    captured_keys: Vec<String>,
    
    // Data
    available_themes: Vec<ThemeInfo>,
    keybinding_actions: Vec<String>,
}

impl Default for Settings {
    fn default() -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            settings_manager: None,
            shared_settings: None,
            keybinding_manager: None,
            current_section: SettingsSection::Appearance,
            focused_pane: FocusedPane::Sections,
            selected_theme: 0,
            selected_keybinding: 0,
            editing_keybinding: false,
            captured_keys: Vec::new(),
            available_themes: available_themes(),
            keybinding_actions: vec![
                "up".to_string(),
                "down".to_string(),
                "left".to_string(),
                "right".to_string(),
                "back".to_string(),
                "quit".to_string(),
                "edit".to_string(),
                "delete".to_string(),
                "create".to_string(),
                "import".to_string(),
                "launch".to_string(),
                "select".to_string(),
                "search".to_string(),
            ],
        }
    }
}

impl Settings {
    pub fn new() -> Self {
        let mut settings = Self::default();
        
        // Always initialize settings manager, creating defaults if needed
        let manager = match SettingsManager::new() {
            Ok(m) => m,
            Err(_) => {
                // If loading fails, create a manager with default settings
                // This should always work unless there's a serious system issue
                let default_settings = UserSettings::default();
                let settings_path = dirs::data_local_dir()
                    .unwrap_or_else(|| std::path::PathBuf::from("."))
                    .join("gemini-cli-manager")
                    .join("settings.json");
                
                SettingsManager {
                    settings: default_settings,
                    settings_path,
                }
            }
        };
        
        // Set selected theme based on current setting
        let current_theme = &manager.get_settings().theme;
        if let Some(index) = settings.available_themes.iter()
            .position(|t| t.name == *current_theme) {
            settings.selected_theme = index;
        }
        
        settings.settings_manager = Some(manager);
        settings
    }

    fn get_sections() -> Vec<&'static str> {
        vec!["Appearance", "Keybindings"]
    }

    fn navigate_sections(&mut self, direction: isize) {
        let sections = Self::get_sections();
        let current_index = match self.current_section {
            SettingsSection::Appearance => 0,
            SettingsSection::Keybindings => 1,
        };
        
        let new_index = (current_index as isize + direction)
            .max(0)
            .min(sections.len() as isize - 1) as usize;
        
        self.current_section = match new_index {
            0 => SettingsSection::Appearance,
            1 => SettingsSection::Keybindings,
            _ => SettingsSection::Appearance,
        };
    }

    fn navigate_content(&mut self, direction: isize) {
        match self.current_section {
            SettingsSection::Appearance => {
                let len = self.available_themes.len();
                if len > 0 {
                    self.selected_theme = ((self.selected_theme as isize + direction).rem_euclid(len as isize)) as usize;
                }
            }
            SettingsSection::Keybindings => {
                let len = self.keybinding_actions.len();
                if len > 0 {
                    self.selected_keybinding = ((self.selected_keybinding as isize + direction).rem_euclid(len as isize)) as usize;
                }
            }
        }
    }

    fn apply_theme_change(&mut self) -> Result<()> {
        if let Some(theme) = self.available_themes.get(self.selected_theme) {
            // Apply theme immediately for live preview
            if let Err(e) = crate::theme::set_theme_by_name(&theme.name) {
                eprintln!("Error setting theme: {}", e);
                return Ok(());
            }

            // Update shared settings first
            if let Some(ref shared_settings) = self.shared_settings {
                if let Ok(mut settings_guard) = shared_settings.write() {
                    settings_guard.theme = theme.name.clone();
                }
            }

            // Then persist to disk
            if let Some(manager) = &mut self.settings_manager {
                match manager.update_theme(theme.name.clone()) {
                    Ok(()) => {
                        // Send success notification
                        if let Some(tx) = &self.command_tx {
                            let _ = tx.send(Action::Success(format!("Theme changed to {}", theme.display_name)));
                        }
                    }
                    Err(e) => {
                        // Send error notification
                        if let Some(tx) = &self.command_tx {
                            let _ = tx.send(Action::Error(format!("Failed to save theme: {}", e)));
                        }
                    }
                }
            }

            // Trigger re-render
            if let Some(tx) = &self.command_tx {
                let _ = tx.send(Action::Render);
            }
        }
        Ok(())
    }

    fn render_sections(&self, frame: &mut Frame, area: Rect) {
        let sections = Self::get_sections();
        let items: Vec<ListItem> = sections
            .iter()
            .enumerate()
            .map(|(i, &section)| {
                let style = if (i == 0 && self.current_section == SettingsSection::Appearance) ||
                             (i == 1 && self.current_section == SettingsSection::Keybindings) {
                    Style::default().fg(theme::primary()).add_modifier(Modifier::BOLD)
                } else {
                    Style::default().fg(theme::text_primary())
                };
                ListItem::new(section).style(style)
            })
            .collect();

        let list = List::new(items)
            .block(
                Block::default()
                    .title(" Settings ")
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(
                        if self.focused_pane == FocusedPane::Sections {
                            theme::border_focused()
                        } else {
                            theme::border()
                        }
                    ))
                    .border_type(BorderType::Rounded),
            )
            .style(Style::default().fg(theme::text_primary()));

        frame.render_widget(list, area);
    }

    fn render_appearance(&self, frame: &mut Frame, area: Rect) {
        let items: Vec<ListItem> = self.available_themes
            .iter()
            .enumerate()
            .map(|(i, theme)| {
                let mut content = vec![
                    Span::styled(&theme.display_name, Style::default().fg(theme::text_primary())),
                    Span::styled(" (", Style::default().fg(theme::text_muted())),
                    Span::styled(&theme.variant, Style::default().fg(theme::text_muted())),
                    Span::styled(")", Style::default().fg(theme::text_muted())),
                ];

                // Add preview colors
                content.push(Span::styled("  ", Style::default()));
                
                // Show current selection indicator
                if i == self.selected_theme {
                    content.insert(0, Span::styled("● ", Style::default().fg(theme::success())));
                } else {
                    content.insert(0, Span::styled("  ", Style::default()));
                }

                ListItem::new(Line::from(content))
            })
            .collect();

        let mut state = ListState::default();
        state.select(Some(self.selected_theme));

        let list = List::new(items)
            .block(
                Block::default()
                    .title(" Theme Selection ")
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(
                        if self.focused_pane == FocusedPane::Content && 
                           self.current_section == SettingsSection::Appearance {
                            theme::border_focused()
                        } else {
                            theme::border()
                        }
                    ))
                    .border_type(BorderType::Rounded),
            )
            .highlight_style(Style::default().bg(theme::selection()))
            .style(Style::default().fg(theme::text_primary()));

        frame.render_stateful_widget(list, area, &mut state);
    }

    fn render_keybindings(&self, frame: &mut Frame, area: Rect) {
        // Split area to make room for reset button
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Min(10),     // Keybindings list
                Constraint::Length(3),   // Reset button
            ])
            .split(area);

        let settings = self.settings_manager.as_ref()
            .map(|m| m.get_settings())
            .cloned()
            .unwrap_or_default();

        let items: Vec<ListItem> = self.keybinding_actions
            .iter()
            .enumerate()
            .map(|(idx, action)| {
                let is_editing = self.editing_keybinding && idx == self.selected_keybinding;
                
                let keys = if is_editing && !self.captured_keys.is_empty() {
                    self.captured_keys.clone()
                } else {
                    settings.keybindings.get_keys_for_action(action)
                };
                let keys_str = keys.join(", ");
                
                let mut content = vec![
                    Span::styled(format!("{:12}", action), Style::default().fg(theme::highlight())),
                    Span::styled(" → ", Style::default().fg(theme::text_muted())),
                ];
                
                if is_editing {
                    content.push(Span::styled(
                        if self.captured_keys.is_empty() {
                            "Press keys to capture... (Ctrl+S: save | Esc: cancel | Backspace: remove last)".to_string()
                        } else {
                            format!("{} (add more keys | Ctrl+S: save | Esc: cancel | Backspace: remove last)", keys_str)
                        },
                        Style::default().fg(theme::warning()).add_modifier(Modifier::ITALIC)
                    ));
                } else {
                    content.push(Span::styled(keys_str, Style::default().fg(theme::text_primary())));
                }

                ListItem::new(Line::from(content))
            })
            .collect();

        let mut state = ListState::default();
        state.select(Some(self.selected_keybinding));

        let list = List::new(items)
            .block(
                Block::default()
                    .title(" Keybindings ")
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(
                        if self.focused_pane == FocusedPane::Content && 
                           self.current_section == SettingsSection::Keybindings {
                            theme::border_focused()
                        } else {
                            theme::border()
                        }
                    ))
                    .border_type(BorderType::Rounded),
            )
            .highlight_style(Style::default().bg(theme::selection()))
            .style(Style::default().fg(theme::text_primary()));

        frame.render_stateful_widget(list, chunks[0], &mut state);

        // Render reset button
        let reset_button_text = if self.focused_pane == FocusedPane::Content {
            " Press 'r' to Reset to Defaults "
        } else {
            " Reset to Defaults (r) "
        };
        
        let reset_button = Paragraph::new(reset_button_text)
            .style(Style::default()
                .fg(theme::warning())
                .add_modifier(Modifier::BOLD))
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .border_style(Style::default().fg(theme::warning())),
            );

        frame.render_widget(reset_button, chunks[1]);
    }


    fn render_content(&self, frame: &mut Frame, area: Rect) {
        match self.current_section {
            SettingsSection::Appearance => self.render_appearance(frame, area),
            SettingsSection::Keybindings => self.render_keybindings(frame, area),
        }
    }
}

impl Component for Settings {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.command_tx = Some(tx);
        Ok(())
    }

    fn register_config_handler(&mut self, config: Config) -> Result<()> {
        self.config = config;
        Ok(())
    }

    fn register_settings_handler(&mut self, settings: std::sync::Arc<std::sync::RwLock<UserSettings>>) -> Result<()> {
        self.shared_settings = Some(settings.clone());
        self.keybinding_manager = Some(KeybindingManager::new(settings.clone()));
        
        // Initialize theme selection based on current settings
        if let Ok(settings_guard) = settings.read() {
            if let Some(index) = self.available_themes.iter()
                .position(|t| t.name == settings_guard.theme) {
                self.selected_theme = index;
            }
        }
        
        Ok(())
    }

    fn update(&mut self, action: Action) -> Result<Option<Action>> {
        match action {
            Action::ChangeTheme(theme_name) => {
                if let Some(index) = self.available_themes.iter()
                    .position(|t| t.name == theme_name) {
                    self.selected_theme = index;
                    self.apply_theme_change()?;
                }
                Ok(Some(Action::Render))
            }
            Action::SaveSettings => {
                if let Some(_manager) = &self.settings_manager {
                    // Settings are automatically saved when changes are made
                    // This action could trigger a "Settings saved" notification
                }
                Ok(None)
            }
            Action::ResetKeybindings => {
                if let Some(manager) = &mut self.settings_manager {
                    match manager.reset_keybindings() {
                        Ok(()) => {
                            // Send success notification
                            if let Some(tx) = &self.command_tx {
                                let _ = tx.send(Action::Success("Keybindings reset to defaults".to_string()));
                                let _ = tx.send(Action::Render);
                            }
                        }
                        Err(e) => {
                            // Send error notification
                            if let Some(tx) = &self.command_tx {
                                let _ = tx.send(Action::Error(format!("Failed to reset keybindings: {}", e)));
                            }
                        }
                    }
                }
                Ok(None)
            }
            _ => Ok(None),
        }
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Create layout with sections on left, content on right
        let chunks = Layout::default()
            .direction(Direction::Horizontal)
            .constraints([
                Constraint::Length(18), // Sections sidebar
                Constraint::Min(40),    // Content area
            ])
            .split(area);

        // Render sections
        self.render_sections(frame, chunks[0]);

        // Render content
        self.render_content(frame, chunks[1]);

        // Help text at the bottom
        let help_area = Rect {
            x: area.x,
            y: area.y + area.height.saturating_sub(3),
            width: area.width,
            height: 3,
        };

        use crate::utils::build_help_text;
        let help_text = match self.focused_pane {
            FocusedPane::Sections => {
                build_help_text(&[
                    ("up", "Navigate sections"),
                    ("down", "Navigate sections"),
                    ("right", "Enter section"),
                    ("tab", "Next tab"),
                    ("quit", "Quit"),
                ])
            }
            FocusedPane::Content => match self.current_section {
                SettingsSection::Appearance => {
                    build_help_text(&[
                        ("up", "Select theme"),
                        ("down", "Select theme"),
                        ("select", "Apply"),
                        ("left", "Back"),
                        ("tab", "Next tab"),
                        ("quit", "Quit"),
                    ])
                }
                SettingsSection::Keybindings => {
                    if self.editing_keybinding {
                        " Press any key to add | Backspace: Remove last | Ctrl+S: Save | Esc: Cancel ".to_string()
                    } else {
                        build_help_text(&[
                            ("up", "Select action"),
                            ("down", "Select action"),
                            ("select", "Edit keybinding"),
                            ("r", "Reset to defaults"),
                            ("left", "Back"),
                            ("tab", "Next tab"),
                            ("quit", "Quit"),
                        ])
                    }
                }
            },
        };

        let help_bar = Paragraph::new(help_text)
            .style(Style::default().fg(theme::text_muted()))
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .border_style(Style::default().fg(theme::border())),
            );

        frame.render_widget(help_bar, help_area);

        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::KeyCode;

        // Handle key capture mode first
        if self.editing_keybinding {
            if let Some(crate::tui::Event::Key(key)) = event {
                match (key.code, key.modifiers) {
                    (KeyCode::Esc, _) => {
                        // Cancel editing
                        self.editing_keybinding = false;
                        self.captured_keys.clear();
                        return Ok(Some(Action::Render));
                    }
                    (KeyCode::Char('s'), crossterm::event::KeyModifiers::CONTROL) => {
                        // Save the captured keys with Ctrl+S
                        if !self.captured_keys.is_empty() {
                            if let Some(action) = self.keybinding_actions.get(self.selected_keybinding) {
                                self.editing_keybinding = false;
                                let keys = self.captured_keys.clone();
                                self.captured_keys.clear();
                                
                                // Update shared settings first
                                if let Some(ref shared_settings) = self.shared_settings {
                                    if let Ok(mut settings_guard) = shared_settings.write() {
                                        // Update the appropriate keybinding in shared settings
                                        match action.as_str() {
                                            "up" => settings_guard.keybindings.navigation.up = keys.clone(),
                                            "down" => settings_guard.keybindings.navigation.down = keys.clone(),
                                            "left" => settings_guard.keybindings.navigation.left = keys.clone(),
                                            "right" => settings_guard.keybindings.navigation.right = keys.clone(),
                                            "back" => settings_guard.keybindings.navigation.back = keys.clone(),
                                            "quit" => settings_guard.keybindings.navigation.quit = keys.clone(),
                                            "edit" => settings_guard.keybindings.actions.edit = keys.clone(),
                                            "delete" => settings_guard.keybindings.actions.delete = keys.clone(),
                                            "create" => settings_guard.keybindings.actions.create = keys.clone(),
                                            "import" => settings_guard.keybindings.actions.import = keys.clone(),
                                            "launch" => settings_guard.keybindings.actions.launch = keys.clone(),
                                            "select" => settings_guard.keybindings.actions.select = keys.clone(),
                                            "search" => settings_guard.keybindings.actions.search = keys.clone(),
                                            _ => {}
                                        }
                                    }
                                }

                                // Then persist to disk
                                if let Some(manager) = &mut self.settings_manager {
                                    match manager.update_keybinding(action, keys.clone()) {
                                        Ok(()) => {
                                            if let Some(tx) = &self.command_tx {
                                                let _ = tx.send(Action::Success(format!("Keybinding for '{}' updated successfully", action)));
                                                let _ = tx.send(Action::Render);
                                            }
                                        }
                                        Err(e) => {
                                            if let Some(tx) = &self.command_tx {
                                                let _ = tx.send(Action::Error(format!("Failed to update keybinding: {}", e)));
                                            }
                                        }
                                    }
                                } else {
                                    if let Some(tx) = &self.command_tx {
                                        let _ = tx.send(Action::Error("Settings manager not initialized".to_string()));
                                    }
                                }
                                return Ok(Some(Action::Render));
                            }
                        }
                        return Ok(Some(Action::Render));
                    }
                    (KeyCode::Backspace, _) => {
                        // Remove the last captured key
                        if !self.captured_keys.is_empty() {
                            self.captured_keys.pop();
                        }
                        return Ok(Some(Action::Render));
                    }
                    _ => {
                        // Capture the key
                        let key_str = format_key_event(&key);
                        // Add to the list if not already present
                        if !self.captured_keys.contains(&key_str) {
                            self.captured_keys.push(key_str);
                        }
                        return Ok(Some(Action::Render));
                    }
                }
            }
            return Ok(None);
        }

        // Normal mode handling
        match event {
            Some(crate::tui::Event::Key(key)) => {
                // Use keybinding manager if available
                if let Some(ref kb_manager) = self.keybinding_manager {
                    // Check navigation keybindings
                    if kb_manager.matches(&key, "quit") {
                        return Ok(Some(Action::Quit));
                    } else if kb_manager.matches(&key, "back") {
                        return Ok(Some(Action::NavigateBack));
                    } else if key.code == KeyCode::Tab {
                        return Ok(Some(Action::NavigateToExtensions));
                    } else if kb_manager.matches(&key, "up") {
                        match self.focused_pane {
                            FocusedPane::Sections => self.navigate_sections(-1),
                            FocusedPane::Content => self.navigate_content(-1),
                        }
                        return Ok(Some(Action::Render));
                    } else if kb_manager.matches(&key, "down") {
                        match self.focused_pane {
                            FocusedPane::Sections => self.navigate_sections(1),
                            FocusedPane::Content => self.navigate_content(1),
                        }
                        return Ok(Some(Action::Render));
                    } else if kb_manager.matches(&key, "right") {
                        if self.focused_pane == FocusedPane::Sections {
                            self.focused_pane = FocusedPane::Content;
                            return Ok(Some(Action::Render));
                        }
                    } else if kb_manager.matches(&key, "left") {
                        if self.focused_pane == FocusedPane::Content {
                            self.focused_pane = FocusedPane::Sections;
                            return Ok(Some(Action::Render));
                        }
                    } else if kb_manager.matches(&key, "select") || key.code == KeyCode::Enter {
                        match self.current_section {
                            SettingsSection::Appearance => {
                                self.apply_theme_change()?;
                                return Ok(Some(Action::Render));
                            }
                            SettingsSection::Keybindings => {
                                // Start editing the selected keybinding
                                self.editing_keybinding = true;
                                self.captured_keys.clear();
                                return Ok(Some(Action::Render));
                            }
                        }
                    } else if key.code == KeyCode::Char('r') {
                        // Only handle reset when in keybindings section and content pane is focused
                        if self.current_section == SettingsSection::Keybindings && 
                           self.focused_pane == FocusedPane::Content &&
                           !self.editing_keybinding {
                            return Ok(Some(Action::ResetKeybindings));
                        }
                    }
                } else {
                    // Fallback to hardcoded keybindings if manager not available
                    match key.code {
                        KeyCode::Char('q') => return Ok(Some(Action::Quit)),
                        KeyCode::Esc => return Ok(Some(Action::NavigateBack)),
                        KeyCode::Tab => return Ok(Some(Action::NavigateToExtensions)),
                        
                        KeyCode::Up | KeyCode::Char('k') => {
                            match self.focused_pane {
                                FocusedPane::Sections => self.navigate_sections(-1),
                                FocusedPane::Content => self.navigate_content(-1),
                            }
                            return Ok(Some(Action::Render));
                        }
                        
                        KeyCode::Down | KeyCode::Char('j') => {
                            match self.focused_pane {
                                FocusedPane::Sections => self.navigate_sections(1),
                                FocusedPane::Content => self.navigate_content(1),
                            }
                            return Ok(Some(Action::Render));
                        }
                        
                        KeyCode::Right | KeyCode::Char('l') => {
                            if self.focused_pane == FocusedPane::Sections {
                                self.focused_pane = FocusedPane::Content;
                                return Ok(Some(Action::Render));
                            }
                        }
                        
                        KeyCode::Left | KeyCode::Char('h') => {
                            if self.focused_pane == FocusedPane::Content {
                                self.focused_pane = FocusedPane::Sections;
                                return Ok(Some(Action::Render));
                            }
                        }
                        
                        KeyCode::Enter => {
                            match self.current_section {
                                SettingsSection::Appearance => {
                                    self.apply_theme_change()?;
                                    return Ok(Some(Action::Render));
                        }
                                SettingsSection::Keybindings => {
                                    // Start editing the selected keybinding
                                    self.editing_keybinding = true;
                                    self.captured_keys.clear();
                                    return Ok(Some(Action::Render));
                                }
                            }
                        }
                        
                        KeyCode::Char('r') => {
                            // Only handle reset when in keybindings section and content pane is focused
                            if self.current_section == SettingsSection::Keybindings && 
                               self.focused_pane == FocusedPane::Content &&
                               !self.editing_keybinding {
                                return Ok(Some(Action::ResetKeybindings));
                            }
                        }
                        
                        _ => {}
                    }
                }
                Ok(None)
            }
            _ => Ok(None),
        }
    }
}

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