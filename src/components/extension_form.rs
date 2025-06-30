use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;
use chrono::Utc;
use std::collections::HashMap;
use tui_input::Input;
use tui_input::backend::crossterm::EventHandler;

use super::Component;
use crate::{action::Action, config::Config, models::{Extension, extension::{ExtensionMetadata, McpServerConfig}}, storage::Storage};

#[derive(Debug, Clone)]
enum FormField {
    Name,
    Version,
    Description,
    ContextFileName,
    ContextContent,
    McpServers,
    Tags,
}

pub struct ExtensionForm {
    command_tx: Option<UnboundedSender<Action>>,
    config: Config,
    storage: Storage,
    
    // Form state using tui-input
    name_input: Input,
    version_input: Input,
    description_input: Input,
    context_file_name_input: Input,
    context_content_input: Input,
    context_scroll_offset: u16,
    tags_input: Input,
    
    // MCP servers management
    mcp_servers: HashMap<String, McpServerConfig>,
    mcp_server_cursor: usize,
    editing_server: Option<String>,
    server_name_input: Input,
    server_command_input: Input,
    server_args_input: Input,
    server_env_input: Input,
    server_cwd_input: Input,
    server_timeout_input: Input,
    server_trust_input: bool,
    server_field_cursor: usize, // Which field in the server editor is active
    
    // Form navigation
    current_field: FormField,
    
    // Edit mode (if editing existing extension)
    edit_mode: bool,
    edit_extension_id: Option<String>,
}

impl ExtensionForm {
    pub fn new(storage: Storage) -> Self {
        Self {
            command_tx: None,
            config: Config::default(),
            storage,
            name_input: Input::default(),
            version_input: Input::from("1.0.0"),
            description_input: Input::default(),
            context_file_name_input: Input::from("GEMINI.md"),
            context_content_input: Input::default(),
            context_scroll_offset: 0,
            tags_input: Input::default(),
            mcp_servers: HashMap::new(),
            mcp_server_cursor: 0,
            editing_server: None,
            server_name_input: Input::default(),
            server_command_input: Input::default(),
            server_args_input: Input::default(),
            server_env_input: Input::default(),
            server_cwd_input: Input::default(),
            server_timeout_input: Input::default(),
            server_trust_input: false,
            server_field_cursor: 0,
            current_field: FormField::Name,
            edit_mode: false,
            edit_extension_id: None,
        }
    }
    
    pub fn with_extension(storage: Storage, extension: &Extension) -> Self {
        let name_input = Input::from(extension.name.clone());
        let version_input = Input::from(extension.version.clone());
        let description_input = Input::from(extension.description.clone().unwrap_or_default());
        let context_file_name_input = Input::from(
            extension.context_file_name.clone()
                .filter(|s| !s.trim().is_empty())  // Filter out empty strings
                .unwrap_or_else(|| "GEMINI.md".to_string())  // Use GEMINI.md for empty/missing
        );
        let context_content_input = Input::from(extension.context_content.clone().unwrap_or_default());
        let tags_input = Input::from(extension.metadata.tags.join(", "));
        
        Self {
            command_tx: None,
            config: Config::default(),
            storage,
            name_input,
            version_input,
            description_input,
            context_file_name_input,
            context_content_input,
            context_scroll_offset: 0,
            tags_input,
            mcp_servers: extension.mcp_servers.clone(),
            mcp_server_cursor: 0,
            editing_server: None,
            server_name_input: Input::default(),
            server_command_input: Input::default(),
            server_args_input: Input::default(),
            server_env_input: Input::default(),
            server_cwd_input: Input::default(),
            server_timeout_input: Input::default(),
            server_trust_input: false,
            server_field_cursor: 0,
            current_field: FormField::Name,
            edit_mode: true,
            edit_extension_id: Some(extension.id.clone()),
        }
    }
    
    fn save_extension(&self) -> Result<()> {
        let extension_id = if let Some(id) = &self.edit_extension_id {
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
        
        let extension = Extension {
            id: extension_id,
            name: self.name_input.value().to_string(),
            version: self.version_input.value().to_string(),
            description: if self.description_input.value().is_empty() {
                None
            } else {
                Some(self.description_input.value().to_string())
            },
            mcp_servers: self.mcp_servers.clone(),
            context_file_name: {
                let trimmed = self.context_file_name_input.value().trim();
                if trimmed.is_empty() || trimmed == "GEMINI.md" {
                    None  // Don't save default or empty values
                } else {
                    Some(trimmed.to_string())
                }
            },
            context_content: if self.context_content_input.value().is_empty() {
                None
            } else {
                Some(self.context_content_input.value().to_string())
            },
            metadata: ExtensionMetadata {
                imported_at: if self.edit_mode {
                    // Preserve original import date
                    self.storage.load_extension(&self.edit_extension_id.as_ref().unwrap())
                        .map(|e| e.metadata.imported_at)
                        .unwrap_or_else(|_| Utc::now())
                } else {
                    Utc::now()
                },
                source_path: None,
                tags,
            },
        };
        
        self.storage.save_extension(&extension)?;
        Ok(())
    }
    
    fn next_field(&mut self) {
        self.current_field = match self.current_field {
            FormField::Name => FormField::Version,
            FormField::Version => FormField::Description,
            FormField::Description => FormField::Tags,
            FormField::Tags => FormField::ContextFileName,
            FormField::ContextFileName => FormField::ContextContent,
            FormField::ContextContent => FormField::McpServers,
            FormField::McpServers => FormField::Name,
        };
    }
    
    fn previous_field(&mut self) {
        self.current_field = match self.current_field {
            FormField::Name => FormField::McpServers,
            FormField::Version => FormField::Name,
            FormField::Description => FormField::Version,
            FormField::Tags => FormField::Description,
            FormField::ContextFileName => FormField::Tags,
            FormField::ContextContent => FormField::ContextFileName,
            FormField::McpServers => FormField::ContextContent,
        };
    }
    
    fn start_add_server(&mut self) {
        self.editing_server = Some(String::new());
        self.server_name_input.reset();
        self.server_command_input.reset();
        self.server_args_input.reset();
        self.server_env_input.reset();
        self.server_cwd_input.reset();
        self.server_timeout_input.reset();
        self.server_trust_input = false;
        self.server_field_cursor = 0;
    }
    
    fn save_server(&mut self) {
        if let Some(_) = &self.editing_server {
            let name = self.server_name_input.value().to_string();
            let command = self.server_command_input.value().to_string();
            let args: Vec<String> = self.server_args_input.value()
                .split_whitespace()
                .map(|s| s.to_string())
                .collect();
            
            // Parse environment variables (format: KEY=VALUE,KEY2=VALUE2)
            let env: HashMap<String, String> = self.server_env_input.value()
                .split(',')
                .filter_map(|s| {
                    let parts: Vec<&str> = s.trim().splitn(2, '=').collect();
                    if parts.len() == 2 {
                        Some((parts[0].to_string(), parts[1].to_string()))
                    } else {
                        None
                    }
                })
                .collect();
            
            let cwd = self.server_cwd_input.value().to_string();
            let timeout = self.server_timeout_input.value().parse::<u64>().ok();
            
            if !name.is_empty() && !command.is_empty() {
                let server = McpServerConfig {
                    url: None, // MCP servers are command-based, not URL-based
                    command: Some(command),
                    args: if args.is_empty() { None } else { Some(args) },
                    cwd: if cwd.is_empty() { None } else { Some(cwd) },
                    env: if env.is_empty() { None } else { Some(env) },
                    timeout,
                    trust: if self.server_trust_input { Some(true) } else { None },
                };
                self.mcp_servers.insert(name, server);
                self.editing_server = None;
                self.server_field_cursor = 0;
            }
        }
    }
    
    fn delete_selected_server(&mut self) {
        let server_names: Vec<String> = self.mcp_servers.keys().cloned().collect();
        if let Some(name) = server_names.get(self.mcp_server_cursor) {
            self.mcp_servers.remove(name);
            if self.mcp_server_cursor > 0 && self.mcp_server_cursor >= self.mcp_servers.len() {
                self.mcp_server_cursor -= 1;
            }
        }
    }
}

impl Component for ExtensionForm {
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
            " Edit Extension "
        } else {
            " Create New Extension "
        };
        
        let block = Block::default()
            .title(title)
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Rounded)
            .border_style(Style::default().fg(Color::Cyan));
        
        let inner = block.inner(area);
        frame.render_widget(block, area);
        
        // Create main vertical layout
        let main_chunks = ratatui::layout::Layout::default()
            .direction(ratatui::layout::Direction::Vertical)
            .margin(1)
            .constraints([
                Constraint::Length(2),  // Name/Version row
                Constraint::Length(2),  // Description/Tags row  
                Constraint::Min(10),    // Editors area
                Constraint::Length(2),  // Help
            ])
            .split(inner);
        
        // First row: Name and Version side by side
        let first_row = ratatui::layout::Layout::default()
            .direction(ratatui::layout::Direction::Horizontal)
            .constraints([
                Constraint::Percentage(60), // Name
                Constraint::Percentage(40), // Version
            ])
            .split(main_chunks[0]);
        
        // Second row: Description and Tags side by side
        let second_row = ratatui::layout::Layout::default()
            .direction(ratatui::layout::Direction::Horizontal)
            .constraints([
                Constraint::Percentage(70), // Description
                Constraint::Percentage(30), // Tags
            ])
            .split(main_chunks[1]);
        
        // Editors area: Context file editor and MCP servers side by side
        let editors_area = ratatui::layout::Layout::default()
            .direction(ratatui::layout::Direction::Horizontal)
            .constraints([
                Constraint::Percentage(50), // Context file editor
                Constraint::Percentage(50), // MCP servers
            ])
            .split(main_chunks[2]);
        
        // Render compact fields
        // Name field (left side of first row)
        let name_style = if matches!(self.current_field, FormField::Name) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        let name_block = Block::default()
            .borders(Borders::NONE);
        let name_area = first_row[0];
        frame.render_widget(name_block, name_area);
        
        let name_label = "Name: ";
        frame.render_widget(
            Span::raw(name_label),
            Rect::new(name_area.x, name_area.y, name_label.len() as u16, 1)
        );
        frame.render_widget(
            Span::styled(self.name_input.value(), name_style),
            Rect::new(name_area.x + name_label.len() as u16, name_area.y, name_area.width - name_label.len() as u16, 1)
        );
        
        // Version field (right side of first row)
        let version_style = if matches!(self.current_field, FormField::Version) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        let version_area = first_row[1];
        let version_label = "Version: ";
        frame.render_widget(
            Span::raw(version_label),
            Rect::new(version_area.x + 1, version_area.y, version_label.len() as u16, 1)
        );
        frame.render_widget(
            Span::styled(self.version_input.value(), version_style),
            Rect::new(version_area.x + 1 + version_label.len() as u16, version_area.y, version_area.width - version_label.len() as u16 - 1, 1)
        );
        
        // Description field (left side of second row)
        let desc_style = if matches!(self.current_field, FormField::Description) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        let desc_area = second_row[0];
        let desc_label = "Description: ";
        frame.render_widget(
            Span::raw(desc_label),
            Rect::new(desc_area.x, desc_area.y, desc_label.len() as u16, 1)
        );
        frame.render_widget(
            Span::styled(self.description_input.value(), desc_style),
            Rect::new(desc_area.x + desc_label.len() as u16, desc_area.y, desc_area.width - desc_label.len() as u16, 1)
        );
        
        // Tags field (right side of second row)
        let tags_style = if matches!(self.current_field, FormField::Tags) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        let tags_area = second_row[1];
        let tags_label = "Tags: ";
        frame.render_widget(
            Span::raw(tags_label),
            Rect::new(tags_area.x + 1, tags_area.y, tags_label.len() as u16, 1)
        );
        frame.render_widget(
            Span::styled(self.tags_input.value(), tags_style),
            Rect::new(tags_area.x + 1 + tags_label.len() as u16, tags_area.y, tags_area.width - tags_label.len() as u16 - 1, 1)
        );
        
        // Set cursor for compact fields
        match self.current_field {
            FormField::Name => {
                let cursor_pos = self.name_input.visual_cursor();
                frame.set_cursor_position((
                    first_row[0].x + name_label.len() as u16 + cursor_pos as u16,
                    first_row[0].y
                ));
            }
            FormField::Version => {
                let cursor_pos = self.version_input.visual_cursor();
                frame.set_cursor_position((
                    first_row[1].x + 1 + version_label.len() as u16 + cursor_pos as u16,
                    first_row[1].y
                ));
            }
            FormField::Description => {
                let cursor_pos = self.description_input.visual_cursor();
                frame.set_cursor_position((
                    second_row[0].x + desc_label.len() as u16 + cursor_pos as u16,
                    second_row[0].y
                ));
            }
            FormField::Tags => {
                let cursor_pos = self.tags_input.visual_cursor();
                frame.set_cursor_position((
                    second_row[1].x + 1 + tags_label.len() as u16 + cursor_pos as u16,
                    second_row[1].y
                ));
            }
            _ => {}
        }
        
        // Context file editor (left side)
        let context_editor_area = editors_area[0];
        let context_chunks = ratatui::layout::Layout::default()
            .direction(ratatui::layout::Direction::Vertical)
            .constraints([
                Constraint::Length(3), // File name field
                Constraint::Min(5),    // Content editor
            ])
            .split(context_editor_area);
        
        // Context file name with default placeholder
        let context_name_style = if matches!(self.current_field, FormField::ContextFileName) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        let context_name_block = Block::default()
            .title("Context File")
            .borders(Borders::ALL)
            .border_style(context_name_style);
        frame.render_widget(context_name_block.clone(), context_chunks[0]);
        
        let context_name_inner = context_name_block.inner(context_chunks[0]);
        let display_text = self.context_file_name_input.value();
        let context_name_text = Paragraph::new(display_text)
            .style(if display_text == "GEMINI.md" && !matches!(self.current_field, FormField::ContextFileName) {
                Style::default().fg(Color::DarkGray).italic()
            } else {
                Style::default()
            });
        frame.render_widget(context_name_text, context_name_inner);
        
        // Context content editor
        let context_content_style = if matches!(self.current_field, FormField::ContextContent) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        let context_content_block = Block::default()
            .title("Context Content (Markdown)")
            .borders(Borders::ALL)
            .border_style(context_content_style);
        frame.render_widget(context_content_block.clone(), context_chunks[1]);
        
        let context_content_inner = context_content_block.inner(context_chunks[1]);
        let context_content_text = Paragraph::new(self.context_content_input.value())
            .wrap(Wrap { trim: false })
            .scroll((self.context_scroll_offset, 0));
        frame.render_widget(context_content_text, context_content_inner);
        
        // Set cursor for context fields
        if matches!(self.current_field, FormField::ContextFileName) {
            let cursor_pos = self.context_file_name_input.visual_cursor();
            frame.set_cursor_position((
                context_name_inner.x + cursor_pos as u16,
                context_name_inner.y
            ));
        } else if matches!(self.current_field, FormField::ContextContent) {
            let cursor_pos = self.context_content_input.visual_cursor();
            let content = self.context_content_input.value();
            let lines_before_cursor = content[..cursor_pos.min(content.len())].lines().count();
            let current_line_start = content[..cursor_pos.min(content.len())].rfind('\n').map(|p| p + 1).unwrap_or(0);
            let col = cursor_pos.saturating_sub(current_line_start);
            
            // Adjust for scrolling
            let visible_line = (lines_before_cursor as u16).saturating_sub(self.context_scroll_offset + 1);
            if visible_line < context_content_inner.height {
                frame.set_cursor_position((
                    context_content_inner.x + col as u16,
                    context_content_inner.y + visible_line
                ));
            }
        }
        
        // MCP Servers section (right side)
        let mcp_editor_area = editors_area[1];
        let mcp_style = if matches!(self.current_field, FormField::McpServers) {
            Style::default().fg(Color::Yellow)
        } else {
            Style::default()
        };
        
        if let Some(_) = &self.editing_server {
            // Show server edit form
            let server_block = Block::default()
                .title("Add MCP Server (Tab: Next field, Enter to save, Esc to cancel)")
                .borders(Borders::ALL)
                .border_style(Style::default().fg(Color::Green));
            
            let server_inner = server_block.inner(mcp_editor_area);
            frame.render_widget(server_block, mcp_editor_area);
            
            let server_chunks = ratatui::layout::Layout::default()
                .direction(ratatui::layout::Direction::Vertical)
                .constraints([
                    Constraint::Length(1), // Name
                    Constraint::Length(1), // Command
                    Constraint::Length(1), // Args
                    Constraint::Length(1), // Env
                    Constraint::Length(1), // CWD
                    Constraint::Length(1), // Timeout
                    Constraint::Length(1), // Trust
                    Constraint::Min(0),
                ])
                .split(server_inner);
            
            let fields = [
                ("Name: ", self.server_name_input.value(), false),
                ("Command: ", self.server_command_input.value(), false),
                ("Args: ", self.server_args_input.value(), false),
                ("Env (KEY=VALUE,KEY2=VALUE2): ", self.server_env_input.value(), false),
                ("Working Dir: ", self.server_cwd_input.value(), false),
                ("Timeout (ms): ", self.server_timeout_input.value(), false),
                ("Trust: ", if self.server_trust_input { "Yes" } else { "No" }, true),
            ];
            
            // Render each field
            for (i, (label, value, _is_bool)) in fields.iter().enumerate() {
                let style = if i == self.server_field_cursor {
                    Style::default().fg(Color::Yellow)
                } else {
                    Style::default()
                };
                
                let line = Line::from(vec![
                    Span::raw(*label),
                    Span::styled(*value, style),
                ]);
                frame.render_widget(line, server_chunks[i]);
            }
            
            // Set cursor position based on which field is being edited
            if self.server_field_cursor < 6 { // Not the boolean field
                let label_len = fields[self.server_field_cursor].0.len() as u16;
                let cursor_pos = match self.server_field_cursor {
                    0 => self.server_name_input.visual_cursor(),
                    1 => self.server_command_input.visual_cursor(),
                    2 => self.server_args_input.visual_cursor(),
                    3 => self.server_env_input.visual_cursor(),
                    4 => self.server_cwd_input.visual_cursor(),
                    5 => self.server_timeout_input.visual_cursor(),
                    _ => 0,
                };
                frame.set_cursor_position((
                    server_chunks[self.server_field_cursor].x + label_len + cursor_pos as u16,
                    server_chunks[self.server_field_cursor].y
                ));
            }
        } else {
            // Show server list
            let mcp_block = Block::default()
                .title("MCP Servers (↑/↓ to navigate, n: new, d: delete)")
                .borders(Borders::ALL)
                .border_style(mcp_style);
            
            let server_items: Vec<ListItem> = self.mcp_servers
                .iter()
                .enumerate()
                .map(|(i, (name, server))| {
                    let is_selected = i == self.mcp_server_cursor && matches!(self.current_field, FormField::McpServers);
                    
                    let style = if is_selected {
                        Style::default().bg(Color::DarkGray)
                    } else {
                        Style::default()
                    };
                    
                    let content = if let Some(url) = &server.url {
                        format!("{}: URL: {}", name, url)
                    } else if let Some(cmd) = &server.command {
                        format!("{}: {} {}", 
                            name, 
                            cmd,
                            server.args.as_ref().map(|a| a.join(" ")).unwrap_or_default()
                        )
                    } else {
                        format!("{}: (no configuration)", name)
                    };
                    
                    ListItem::new(content).style(style)
                })
                .collect();
            
            let mcp_list = if server_items.is_empty() {
                List::new(vec![ListItem::new("  No MCP servers configured")])
                    .block(mcp_block)
            } else {
                List::new(server_items)
                    .block(mcp_block)
                    .highlight_style(Style::default().bg(Color::DarkGray))
            };
            
            frame.render_widget(mcp_list, mcp_editor_area);
        }
        
        
        // Help text
        let help_text = match self.current_field {
            FormField::McpServers if self.editing_server.is_some() => 
                " Enter: Save server | Esc: Cancel | Tab: Next field ",
            FormField::McpServers => 
                " Tab: Next field | ↑/↓: Navigate | n: New server | d: Delete | Ctrl+S: Save | Esc: Cancel ",
            FormField::ContextContent => 
                " Tab: Next field | ↑/↓: Scroll | Type to edit | Ctrl+S: Save | Esc: Cancel ",
            _ => " Tab: Next field | Type to edit | Ctrl+S: Save | Esc: Cancel ",
        };
        let help_style = Style::default().fg(Color::DarkGray);
        frame.render_widget(
            Paragraph::new(help_text)
                .style(help_style)
                .alignment(Alignment::Center),
            main_chunks[3],
        );
        
        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::{KeyCode, KeyModifiers};

        match event {
            Some(crate::tui::Event::Key(key)) => {
                // Handle server editing mode separately
                if self.editing_server.is_some() {
                    match key.code {
                        KeyCode::Esc => {
                            self.editing_server = None;
                            return Ok(Some(Action::Render));
                        }
                        KeyCode::Enter => {
                            self.save_server();
                            return Ok(Some(Action::Render));
                        }
                        KeyCode::Tab => {
                            // Cycle through server fields
                            self.server_field_cursor = (self.server_field_cursor + 1) % 7;
                            return Ok(Some(Action::Render));
                        }
                        KeyCode::BackTab => {
                            // Cycle backwards through server fields
                            if self.server_field_cursor == 0 {
                                self.server_field_cursor = 6;
                            } else {
                                self.server_field_cursor -= 1;
                            }
                            return Ok(Some(Action::Render));
                        }
                        KeyCode::Char(' ') if self.server_field_cursor == 6 => {
                            // Toggle trust field
                            self.server_trust_input = !self.server_trust_input;
                            return Ok(Some(Action::Render));
                        }
                        _ => {
                            // Handle input for the current field
                            match self.server_field_cursor {
                                0 => {
                                    if self.server_name_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                        return Ok(Some(Action::Render));
                                    }
                                }
                                1 => {
                                    if self.server_command_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                        return Ok(Some(Action::Render));
                                    }
                                }
                                2 => {
                                    if self.server_args_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                        return Ok(Some(Action::Render));
                                    }
                                }
                                3 => {
                                    if self.server_env_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                        return Ok(Some(Action::Render));
                                    }
                                }
                                4 => {
                                    if self.server_cwd_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                        return Ok(Some(Action::Render));
                                    }
                                }
                                5 => {
                                    if self.server_timeout_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                        return Ok(Some(Action::Render));
                                    }
                                }
                                _ => {}
                            }
                        }
                    }
                    return Ok(None);
                }
                
                // Normal form handling
                match (key.code, key.modifiers) {
                    (KeyCode::Esc, _) => {
                        return Ok(Some(Action::NavigateBack));
                    }
                    (KeyCode::Char('s'), KeyModifiers::CONTROL) => {
                        // Save extension
                        if !self.name_input.value().is_empty() && !self.version_input.value().is_empty() {
                            match self.save_extension() {
                                Ok(_) => {
                                    // Send refresh action before navigating back
                                    if let Some(tx) = &self.command_tx {
                                        let _ = tx.send(Action::RefreshExtensions);
                                        let _ = tx.send(Action::Render);
                                    }
                                    return Ok(Some(Action::NavigateBack));
                                }
                                Err(e) => {
                                    return Ok(Some(Action::Error(format!("Failed to save extension: {}", e))));
                                }
                            }
                        } else {
                            return Ok(Some(Action::Error("Extension name and version are required".to_string())));
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
                            FormField::Version => {
                                if self.version_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                            FormField::Description => {
                                if self.description_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                            FormField::McpServers => {
                                match key.code {
                                    KeyCode::Up => {
                                        if self.mcp_server_cursor > 0 {
                                            self.mcp_server_cursor -= 1;
                                            return Ok(Some(Action::Render));
                                        }
                                    }
                                    KeyCode::Down => {
                                        if self.mcp_server_cursor < self.mcp_servers.len().saturating_sub(1) {
                                            self.mcp_server_cursor += 1;
                                            return Ok(Some(Action::Render));
                                        }
                                    }
                                    KeyCode::Char('n') => {
                                        self.start_add_server();
                                        return Ok(Some(Action::Render));
                                    }
                                    KeyCode::Char('d') => {
                                        if !self.mcp_servers.is_empty() {
                                            self.delete_selected_server();
                                            return Ok(Some(Action::Render));
                                        }
                                    }
                                    _ => {}
                                }
                            }
                            FormField::ContextFileName => {
                                if self.context_file_name_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                    return Ok(Some(Action::Render));
                                }
                            }
                            FormField::ContextContent => {
                                match key.code {
                                    KeyCode::Up => {
                                        if self.context_scroll_offset > 0 {
                                            self.context_scroll_offset = self.context_scroll_offset.saturating_sub(1);
                                            return Ok(Some(Action::Render));
                                        }
                                    }
                                    KeyCode::Down => {
                                        // Check if we need to scroll
                                        let lines_count = self.context_content_input.value().lines().count();
                                        self.context_scroll_offset = self.context_scroll_offset.saturating_add(1).min(lines_count.saturating_sub(1) as u16);
                                        return Ok(Some(Action::Render));
                                    }
                                    _ => {
                                        if self.context_content_input.handle_event(&crossterm::event::Event::Key(key)).is_some() {
                                            return Ok(Some(Action::Render));
                                        }
                                    }
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