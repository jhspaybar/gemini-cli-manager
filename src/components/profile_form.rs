use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;
use chrono::Utc;
use std::collections::HashMap;
use tui_input::Input;
use tui_input::backend::crossterm::EventHandler;

use super::Component;
use crate::{action::Action, config::Config, models::{Extension, Profile, profile::ProfileMetadata}, storage::Storage, theme};

#[derive(Debug, Clone)]
enum FormField {
    Name,
    Description,
    WorkingDirectory,
    Extensions,
    Tags,
}

pub struct ProfileForm {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    storage: Storage,
    
    // Form state using tui-input
    name_input: Input,
    description_input: Input,
    working_directory_input: Input,
    tags_input: Input,
    selected_extensions: Vec<String>,
    
    // Available extensions
    available_extensions: Vec<Extension>,
    extension_cursor: usize,
    
    // Form navigation
    current_field: FormField,
    
    // Edit mode (if editing existing profile)
    edit_mode: bool,
    edit_profile_id: Option<String>,
}

impl ProfileForm {
    pub fn new(storage: Storage) -> Self {
        let available_extensions = storage.list_extensions().unwrap_or_default();
        
        Self {
            command_tx: None,
            config: Config::default(),
            storage,
            name_input: Input::default(),
            description_input: Input::default(),
            working_directory_input: Input::default(),
            tags_input: Input::default(),
            selected_extensions: Vec::new(),
            available_extensions,
            extension_cursor: 0,
            current_field: FormField::Name,
            edit_mode: false,
            edit_profile_id: None,
        }
    }
    
    pub fn with_profile(storage: Storage, profile: &Profile) -> Self {
        let available_extensions = storage.list_extensions().unwrap_or_default();
        
        let name_input = Input::from(profile.name.clone());
        let description_input = Input::from(profile.description.clone().unwrap_or_default());
        let working_directory_input = Input::from(profile.working_directory.clone().unwrap_or_default());
        let tags_input = Input::from(profile.metadata.tags.join(", "));
        
        Self {
            command_tx: None,
            config: Config::default(),
            storage,
            name_input,
            description_input,
            working_directory_input,
            tags_input,
            selected_extensions: profile.extension_ids.clone(),
            available_extensions,
            extension_cursor: 0,
            current_field: FormField::Name,
            edit_mode: true,
            edit_profile_id: Some(profile.id.clone()),
        }
    }
    
    fn save_profile(&self) -> Result<()> {
        let profile_id = if let Some(id) = &self.edit_profile_id {
            id.clone()
        } else {
            // Generate a simple ID from the name
            self.name_input.value().to_lowercase().replace(' ', "-")
        };
        
        let tags: Vec<String> = self.tags_input.value()
            .split(',')
            .map(|s| s.trim().to_string())
            .filter(|s| !s.is_empty())
            .collect();
        
        let profile = Profile {
            id: profile_id,
            name: self.name_input.value().to_string(),
            description: if self.description_input.value().is_empty() {
                None
            } else {
                Some(self.description_input.value().to_string())
            },
            extension_ids: self.selected_extensions.clone(),
            environment_variables: HashMap::new(), // TODO: Add env var editor
            working_directory: if self.working_directory_input.value().is_empty() {
                None
            } else {
                Some(self.working_directory_input.value().to_string())
            },
            metadata: ProfileMetadata {
                created_at: if self.edit_mode {
                    // Preserve original creation date
                    self.storage.load_profile(&self.edit_profile_id.as_ref().unwrap())
                        .map(|p| p.metadata.created_at)
                        .unwrap_or_else(|_| Utc::now())
                } else {
                    Utc::now()
                },
                updated_at: Utc::now(),
                tags,
                is_default: false,
                icon: None,
            },
        };
        
        self.storage.save_profile(&profile)?;
        Ok(())
    }
    
    fn toggle_extension(&mut self) {
        if let Some(ext) = self.available_extensions.get(self.extension_cursor) {
            let ext_id = &ext.id;
            if let Some(pos) = self.selected_extensions.iter().position(|id| id == ext_id) {
                self.selected_extensions.remove(pos);
            } else {
                self.selected_extensions.push(ext_id.clone());
            }
        }
    }
    
    fn next_field(&mut self) {
        self.current_field = match self.current_field {
            FormField::Name => FormField::Description,
            FormField::Description => FormField::WorkingDirectory,
            FormField::WorkingDirectory => FormField::Extensions,
            FormField::Extensions => FormField::Tags,
            FormField::Tags => FormField::Name,
        };
    }
    
    fn previous_field(&mut self) {
        self.current_field = match self.current_field {
            FormField::Name => FormField::Tags,
            FormField::Description => FormField::Name,
            FormField::WorkingDirectory => FormField::Description,
            FormField::Extensions => FormField::WorkingDirectory,
            FormField::Tags => FormField::Extensions,
        };
    }
}

impl Component for ProfileForm {
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
            Action::Tick => {}
            Action::Render => {}
            _ => {}
        }
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        let title = if self.edit_mode {
            " Edit Profile "
        } else {
            " Create New Profile "
        };
        
        let block = Block::default()
            .title(title)
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(theme::primary()));
        
        let inner = block.inner(area);
        frame.render_widget(block, area);
        
        // Create layout for form fields
        let chunks = ratatui::layout::Layout::default()
            .direction(ratatui::layout::Direction::Vertical)
            .margin(1)
            .constraints([
                Constraint::Length(3), // Name
                Constraint::Length(3), // Description
                Constraint::Length(3), // Working Directory
                Constraint::Min(5),    // Extensions
                Constraint::Length(3), // Tags
                Constraint::Length(3), // Help
            ])
            .split(inner);
        
        // Name field
        let name_style = if matches!(self.current_field, FormField::Name) {
            Style::default().fg(theme::highlight())
        } else {
            Style::default()
        };
        let name_block = Block::default()
            .title("Name")
            .borders(Borders::ALL)
            .border_style(name_style);
        frame.render_widget(name_block.clone(), chunks[0]);
        
        let name_inner = name_block.inner(chunks[0]);
        let name_text = Paragraph::new(self.name_input.value())
            .style(Style::default().fg(theme::text_primary()));
        frame.render_widget(name_text, name_inner);
        
        // Set cursor for name field if it's active
        if matches!(self.current_field, FormField::Name) {
            let cursor_pos = self.name_input.visual_cursor();
            frame.set_cursor_position((
                name_inner.x + cursor_pos as u16,
                name_inner.y
            ));
        }
        
        // Description field
        let desc_style = if matches!(self.current_field, FormField::Description) {
            Style::default().fg(theme::highlight())
        } else {
            Style::default()
        };
        let desc_block = Block::default()
            .title("Description (optional)")
            .borders(Borders::ALL)
            .border_style(desc_style);
        frame.render_widget(desc_block.clone(), chunks[1]);
        
        let desc_inner = desc_block.inner(chunks[1]);
        let desc_text = Paragraph::new(self.description_input.value())
            .style(Style::default().fg(theme::text_primary()));
        frame.render_widget(desc_text, desc_inner);
        
        // Set cursor for description field if it's active
        if matches!(self.current_field, FormField::Description) {
            let cursor_pos = self.description_input.visual_cursor();
            frame.set_cursor_position((
                desc_inner.x + cursor_pos as u16,
                desc_inner.y
            ));
        }
        
        // Working Directory field
        let dir_style = if matches!(self.current_field, FormField::WorkingDirectory) {
            Style::default().fg(theme::highlight())
        } else {
            Style::default()
        };
        let dir_block = Block::default()
            .title("Working Directory (optional)")
            .borders(Borders::ALL)
            .border_style(dir_style);
        frame.render_widget(dir_block.clone(), chunks[2]);
        
        let dir_inner = dir_block.inner(chunks[2]);
        let dir_text = Paragraph::new(self.working_directory_input.value())
            .style(Style::default().fg(theme::text_primary()));
        frame.render_widget(dir_text, dir_inner);
        
        // Set cursor for working directory field if it's active
        if matches!(self.current_field, FormField::WorkingDirectory) {
            let cursor_pos = self.working_directory_input.visual_cursor();
            frame.set_cursor_position((
                dir_inner.x + cursor_pos as u16,
                dir_inner.y
            ));
        }
        
        // Extensions selection
        let ext_style = if matches!(self.current_field, FormField::Extensions) {
            Style::default().fg(theme::highlight())
        } else {
            Style::default()
        };
        let ext_block = Block::default()
            .title("Extensions (↑/↓ to navigate, Space to toggle)")
            .borders(Borders::ALL)
            .border_style(ext_style);
        
        let ext_items: Vec<ListItem> = self.available_extensions
            .iter()
            .enumerate()
            .map(|(i, ext)| {
                let is_selected = self.selected_extensions.contains(&ext.id);
                let is_cursor = i == self.extension_cursor && matches!(self.current_field, FormField::Extensions);
                
                let prefix = if is_selected { "[✓] " } else { "[ ] " };
                let style = if is_cursor {
                    Style::default().bg(theme::selection()).fg(theme::text_primary())
                } else if is_selected {
                    Style::default().fg(theme::success())
                } else {
                    Style::default().fg(theme::text_primary())
                };
                
                ListItem::new(format!("{}{}", prefix, ext.name)).style(style)
            })
            .collect();
        
        let ext_list = List::new(ext_items).block(ext_block);
        frame.render_widget(ext_list, chunks[3]);
        
        // Tags field
        let tags_style = if matches!(self.current_field, FormField::Tags) {
            Style::default().fg(theme::highlight())
        } else {
            Style::default()
        };
        let tags_block = Block::default()
            .title("Tags (comma-separated)")
            .borders(Borders::ALL)
            .border_style(tags_style);
        frame.render_widget(tags_block.clone(), chunks[4]);
        
        let tags_inner = tags_block.inner(chunks[4]);
        let tags_text = Paragraph::new(self.tags_input.value())
            .style(Style::default().fg(theme::text_primary()));
        frame.render_widget(tags_text, tags_inner);
        
        // Set cursor for tags field if it's active
        if matches!(self.current_field, FormField::Tags) {
            let cursor_pos = self.tags_input.visual_cursor();
            frame.set_cursor_position((
                tags_inner.x + cursor_pos as u16,
                tags_inner.y
            ));
        }
        
        // Help text
        let help_text = match self.current_field {
            FormField::Extensions => " Tab: Next field | ↑/↓: Navigate | Space: Toggle | Ctrl+S: Save | Esc: Cancel ",
            _ => " Tab: Next field | Type to edit | Ctrl+S: Save | Esc: Cancel ",
        };
        let help_style = Style::default().fg(theme::text_muted());
        frame.render_widget(
            Paragraph::new(help_text)
                .style(help_style)
                .alignment(Alignment::Center),
            chunks[5],
        );
        
        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::{KeyCode, KeyModifiers};

        match event {
            Some(crate::tui::Event::Key(key)) => {
                match (key.code, key.modifiers) {
                    (KeyCode::Esc, _) => {
                        return Ok(Some(Action::NavigateBack));
                    }
                    (KeyCode::Char('s'), KeyModifiers::CONTROL) => {
                        // Save profile
                        if !self.name_input.value().is_empty() {
                            match self.save_profile() {
                                Ok(_) => {
                                    // Send refresh action before navigating back
                                    if let Some(tx) = &self.command_tx {
                                        let _ = tx.send(Action::RefreshProfiles);
                                        let _ = tx.send(Action::Render);
                                    }
                                    return Ok(Some(Action::NavigateBack));
                                }
                                Err(e) => {
                                    return Ok(Some(Action::Error(format!("Failed to save profile: {}", e))));
                                }
                            }
                        } else {
                            return Ok(Some(Action::Error("Profile name is required".to_string())));
                        }
                    }
                    (KeyCode::Tab, _) => {
                        self.next_field();
                        return Ok(Some(Action::Render));
                    }
                    (KeyCode::BackTab, _) => {
                        self.previous_field();
                        return Ok(Some(Action::Render));
                    }
                    _ => {
                        // Handle field-specific input
                        match self.current_field {
                            FormField::Name => {
                                if self.name_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                            FormField::Description => {
                                if self.description_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                            FormField::WorkingDirectory => {
                                if self.working_directory_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                            FormField::Extensions => {
                                match key.code {
                                    KeyCode::Up => {
                                        if self.extension_cursor > 0 {
                                            self.extension_cursor -= 1;
                                            return Ok(Some(Action::Render));
                                        }
                                    }
                                    KeyCode::Down => {
                                        if self.extension_cursor < self.available_extensions.len() - 1 {
                                            self.extension_cursor += 1;
                                            return Ok(Some(Action::Render));
                                        }
                                    }
                                    KeyCode::Char(' ') => {
                                        self.toggle_extension();
                                        return Ok(Some(Action::Render));
                                    }
                                    _ => {}
                                }
                            }
                            FormField::Tags => {
                                if self.tags_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                        }
                    }
                }
            }
            _ => {}
        }
        Ok(None)
    }
}