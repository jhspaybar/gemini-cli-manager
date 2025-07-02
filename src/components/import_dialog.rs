use chrono::Utc;
use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use ratatui_explorer::{FileExplorer, Theme as ExplorerTheme};
use serde::Deserialize;
use std::collections::HashMap;
use std::path::{Path, PathBuf};
use std::time::Instant;
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{
    action::Action,
    config::Config,
    models::extension::{Extension, ExtensionMetadata, McpServerConfig},
    storage::Storage,
    theme,
    tui::Event,
};

pub struct ImportDialog {
    action_tx: Option<UnboundedSender<Action>>,
    explorer: FileExplorer,
    storage: Storage,
    state: ImportState,
    state_timestamp: Option<Instant>,
}

#[derive(Debug, Clone, PartialEq)]
enum ImportState {
    Selecting,
    Importing,
    Error(String),
}

// Structure that matches the actual extension.json files
#[derive(Debug, Deserialize)]
struct ImportExtension {
    id: Option<String>,
    name: String,
    version: String,
    description: Option<String>,
    #[serde(rename = "mcpServers")]
    mcp_servers: Option<HashMap<String, McpServerConfig>>,
    context_file_name: Option<String>,
    context_content: Option<String>,
    // The metadata in the import files has a different structure than our internal one
    metadata: Option<ImportMetadata>,
}

#[derive(Debug, Deserialize)]
#[allow(dead_code)]
struct ImportMetadata {
    created_at: Option<String>,
    updated_at: Option<String>,
    tags: Option<Vec<String>>,
    is_builtin: Option<bool>,
    author: Option<String>,
}

impl ImportDialog {
    pub fn new(storage: Storage) -> Self {
        // Configure the explorer theme
        let explorer_theme = ExplorerTheme::default()
            .with_dir_style(Style::default().fg(theme::primary()))
            .with_item_style(Style::default().fg(theme::text_primary()))
            .with_highlight_dir_style(
                Style::default()
                    .bg(theme::selection())
                    .fg(theme::background()),
            )
            .with_highlight_item_style(
                Style::default()
                    .bg(theme::selection())
                    .fg(theme::background()),
            );

        // Initialize explorer with theme
        let explorer = FileExplorer::with_theme(explorer_theme).unwrap_or_else(|_| {
            // Fallback to default if theme fails
            FileExplorer::new().unwrap()
        });

        Self {
            action_tx: None,
            explorer,
            storage,
            state: ImportState::Selecting,
            state_timestamp: None,
        }
    }

    /// Reset the dialog state for a fresh import session
    pub fn reset(&mut self) {
        self.state = ImportState::Selecting;
        self.state_timestamp = None;
        // The explorer maintains its own state (current directory)
        // which is fine - users might want to stay in the same directory
    }

    fn import_extension(&mut self, path: PathBuf) -> Result<()> {
        // Check if it's a directory or a file
        if path.is_dir() {
            self.import_from_directory(path)?;
        } else if path.extension().and_then(|s| s.to_str()) == Some("json") {
            // Import JSON file directly
            self.import_from_file(path)?;
        } else if path.extension().and_then(|s| s.to_str()) == Some("md") {
            // Import MD file as a minimal extension
            self.import_from_md_file(path)?;
        } else {
            self.state =
                ImportState::Error("Please select a .json, .md file, or a directory".to_string());
            self.state_timestamp = Some(Instant::now());
        }

        Ok(())
    }

    fn import_from_directory(&mut self, dir_path: PathBuf) -> Result<()> {
        // Look for extension.json or gemini-extension.json
        let extension_path = dir_path.join("extension.json");
        let gemini_extension_path = dir_path.join("gemini-extension.json");

        let json_path = if extension_path.exists() {
            Some(extension_path)
        } else if gemini_extension_path.exists() {
            Some(gemini_extension_path)
        } else {
            None
        };

        // Look for context files
        let context_files = self.find_context_files(&dir_path);

        if json_path.is_none() && context_files.is_empty() {
            self.state = ImportState::Error(
                "No extension.json or context files found in directory".to_string(),
            );
            self.state_timestamp = Some(Instant::now());
            return Ok(());
        }

        if let Some(json_path) = json_path {
            // Import with JSON file
            self.import_from_file(json_path)?;
        } else if let Some((context_name, context_path)) = context_files.first() {
            // Import just the context file as a minimal extension
            self.import_context_as_extension(
                context_path.clone(),
                context_name.clone(),
                Some(dir_path),
            )?;
        }

        Ok(())
    }

    fn find_context_files(&self, dir_path: &Path) -> Vec<(String, PathBuf)> {
        let mut context_files = Vec::new();

        // Common context file patterns (for reference)
        // We'll look for: GEMINI.md, README.md, CONTEXT.md, or any .md file

        if let Ok(entries) = std::fs::read_dir(dir_path) {
            for entry in entries.flatten() {
                let path = entry.path();
                if path.is_file()
                    && let Some(name) = path.file_name().and_then(|n| n.to_str())
                    && name.ends_with(".md") && !name.starts_with(".") {
                        context_files.push((name.to_string(), path));
                    }
            }
        }

        // Sort to prioritize GEMINI.md, then extension-specific names, then generic names
        context_files.sort_by_key(|(name, _)| match name.as_str() {
            "GEMINI.md" => 0,
            name if name.ends_with(".md") && name != "README.md" && name != "CONTEXT.md" => 1,
            "CONTEXT.md" => 2,
            "README.md" => 3,
            _ => 4,
        });

        context_files
    }

    fn import_from_md_file(&mut self, path: PathBuf) -> Result<()> {
        if let Some(file_name) = path.file_name().and_then(|n| n.to_str()) {
            self.import_context_as_extension(
                path.clone(),
                file_name.to_string(),
                path.parent().map(|p| p.to_path_buf()),
            )?;
        } else {
            self.state = ImportState::Error("Invalid file path".to_string());
            self.state_timestamp = Some(Instant::now());
        }
        Ok(())
    }

    fn import_context_as_extension(
        &mut self,
        context_path: PathBuf,
        context_name: String,
        base_dir: Option<PathBuf>,
    ) -> Result<()> {
        self.state = ImportState::Importing;

        // Read the context file
        let context_content = std::fs::read_to_string(&context_path)?;

        // Generate a name from the file or directory
        let extension_name = if let Some(dir) = &base_dir {
            dir.file_name()
                .and_then(|n| n.to_str())
                .unwrap_or("imported-context")
                .to_string()
        } else {
            context_path
                .file_stem()
                .and_then(|n| n.to_str())
                .unwrap_or("imported-context")
                .to_string()
        };

        // Create a minimal extension with just the context
        let extension = Extension {
            id: uuid::Uuid::new_v4().to_string(),
            name: extension_name.clone(),
            version: "1.0.0".to_string(),
            description: Some(format!(
                "Context-only extension imported from {context_name}"
            )),
            mcp_servers: HashMap::new(),
            context_file_name: Some(context_name),
            context_content: Some(context_content),
            metadata: ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: Some(context_path.to_string_lossy().to_string()),
                tags: vec!["context-only".to_string()],
            },
        };

        self.storage.save_extension(&extension)?;

        // Send success message and navigate back
        if let Some(tx) = &self.action_tx {
            let _ = tx.send(Action::Success(format!(
                "Successfully imported context: {extension_name}"
            )));
            let _ = tx.send(Action::RefreshExtensions);
            let _ = tx.send(Action::NavigateBack);
        }

        Ok(())
    }

    fn import_from_file(&mut self, path: PathBuf) -> Result<()> {
        self.state = ImportState::Importing;

        // Read the file
        let content = std::fs::read_to_string(&path)?;

        // Parse as import extension first
        match serde_json::from_str::<ImportExtension>(&content) {
            Ok(import_ext) => {
                // Convert to our Extension format
                let mut extension = Extension {
                    id: import_ext
                        .id
                        .unwrap_or_else(|| uuid::Uuid::new_v4().to_string()),
                    name: import_ext.name,
                    version: import_ext.version,
                    description: import_ext.description,
                    mcp_servers: import_ext.mcp_servers.unwrap_or_default(),
                    context_file_name: import_ext.context_file_name,
                    context_content: import_ext.context_content,
                    metadata: ExtensionMetadata {
                        imported_at: Utc::now(),
                        source_path: Some(path.to_string_lossy().to_string()),
                        tags: import_ext.metadata.and_then(|m| m.tags).unwrap_or_default(),
                    },
                };

                // Always generate a new ID to avoid conflicts
                extension.id = uuid::Uuid::new_v4().to_string();

                // Save the extension
                self.storage.save_extension(&extension)?;

                // Check if there's a context file in the same directory
                if let Some(parent) = path.parent() {
                    // Look for common context file names
                    let potential_names = vec![
                        format!("{}.md", extension.name.to_uppercase()),
                        format!("{}.md", extension.name),
                        "GEMINI.md".to_string(),
                        "CONTEXT.md".to_string(),
                        "README.md".to_string(),
                    ];

                    for name in potential_names {
                        let context_path = parent.join(&name);
                        if context_path.exists()
                            && let Ok(context_content) = std::fs::read_to_string(&context_path) {
                                // Update extension with context
                                let mut updated = extension.clone();
                                // Store original filename for reference, but it will be written as GEMINI.md
                                updated.context_file_name = Some(name.clone());
                                updated.context_content = Some(context_content);
                                self.storage.save_extension(&updated)?;
                                break;
                            }
                    }
                }

                // Send success message and navigate back
                if let Some(tx) = &self.action_tx {
                    let _ = tx.send(Action::Success(format!(
                        "Successfully imported: {}",
                        extension.name
                    )));
                    let _ = tx.send(Action::RefreshExtensions);
                    let _ = tx.send(Action::NavigateBack);
                }
            }
            Err(e) => {
                self.state = ImportState::Error(format!("Failed to parse extension: {e}"));
                self.state_timestamp = Some(Instant::now());
            }
        }

        Ok(())
    }
}

impl Component for ImportDialog {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.action_tx = Some(tx);
        Ok(())
    }

    fn register_config_handler(&mut self, _config: Config) -> Result<()> {
        Ok(())
    }

    fn update(&mut self, action: Action) -> Result<Option<Action>> {
        match action {
            Action::ResetImportDialog => {
                self.reset();
            }
            _ => {
                // No other special handling needed - success navigates back immediately
                // Errors are handled by user key press in handle_events
            }
        }
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Create a centered popup area
        let popup_area = self.centered_rect(80, 80, area);

        // Clear the background
        frame.render_widget(Clear, popup_area);

        // Create the main block
        let block = Block::default()
            .title(" Import Extension ")
            .title_style(
                Style::default()
                    .fg(theme::primary())
                    .add_modifier(Modifier::BOLD),
            )
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_style(Style::default().fg(theme::border()))
            .style(Style::default().bg(theme::surface()));

        let inner_area = block.inner(popup_area);
        frame.render_widget(block, popup_area);

        match &self.state {
            ImportState::Selecting => {
                // Split area for instructions and explorer
                let chunks = Layout::default()
                    .direction(Direction::Vertical)
                    .constraints([
                        Constraint::Length(3), // Instructions
                        Constraint::Min(0),    // Explorer
                        Constraint::Length(3), // Help text
                    ])
                    .split(inner_area);

                // Instructions
                let instructions = Paragraph::new("Select a .json extension file, .md context file, or a directory containing either:")
                    .style(Style::default().fg(theme::text_primary()))
                    .alignment(Alignment::Center);
                frame.render_widget(instructions, chunks[0]);

                // File explorer
                frame.render_widget(&self.explorer.widget(), chunks[1]);

                // Help text
                let help = Paragraph::new(
                    "↑/↓: Navigate | Enter: Select | Esc: Cancel | h: Toggle hidden files",
                )
                .style(Style::default().fg(theme::text_secondary()))
                .alignment(Alignment::Center);
                frame.render_widget(help, chunks[2]);
            }
            ImportState::Importing => {
                let loading = Paragraph::new("Importing extension...")
                    .style(Style::default().fg(theme::primary()))
                    .alignment(Alignment::Center);
                frame.render_widget(loading, inner_area);
            }
            ImportState::Error(msg) => {
                // Split area for error message and instructions
                let chunks = Layout::default()
                    .direction(Direction::Vertical)
                    .constraints([
                        Constraint::Min(0),    // Error message
                        Constraint::Length(3), // Instructions
                    ])
                    .split(inner_area);

                let error = Paragraph::new(format!("✗ {msg}"))
                    .style(
                        Style::default()
                            .fg(theme::error())
                            .add_modifier(Modifier::BOLD),
                    )
                    .alignment(Alignment::Center)
                    .wrap(Wrap { trim: true });
                frame.render_widget(error, chunks[0]);

                let instructions = Paragraph::new("Press any key to continue")
                    .style(Style::default().fg(theme::text_secondary()))
                    .alignment(Alignment::Center);
                frame.render_widget(instructions, chunks[1]);
            }
        }

        Ok(())
    }

    fn handle_events(&mut self, event: Option<Event>) -> Result<Option<Action>> {
        use crossterm::event::{Event as CEvent, KeyCode};

        if let Some(Event::Key(key)) = event {
            match &self.state {
                ImportState::Selecting => {
                    match key.code {
                        KeyCode::Esc => return Ok(Some(Action::NavigateBack)),
                        KeyCode::Enter => {
                            let selected = self.explorer.current();
                            let path = selected.path().to_path_buf();

                            // Check if it's a directory that might contain an extension
                            if selected.is_dir() {
                                // Check if this directory contains an extension
                                let has_extension = path.join("extension.json").exists()
                                    || path.join("gemini-extension.json").exists();

                                if has_extension {
                                    // Import the extension from this directory
                                    if let Err(e) = self.import_extension(path) {
                                        self.state = ImportState::Error(e.to_string());
                                        self.state_timestamp = Some(Instant::now());
                                    }
                                } else {
                                    // Navigate into the directory
                                    self.explorer.handle(&CEvent::Key(key))?;
                                }
                            } else {
                                // It's a file, try to import it
                                if let Err(e) = self.import_extension(path) {
                                    self.state = ImportState::Error(e.to_string());
                                    self.state_timestamp = Some(Instant::now());
                                }
                            }
                        }
                        _ => {
                            // Let the explorer handle other keys
                            if let Some(Event::Key(key)) = event {
                                self.explorer.handle(&CEvent::Key(key))?;
                            }
                        }
                    }
                }
                ImportState::Error(_) => {
                    // Any key press returns to selecting state
                    self.state = ImportState::Selecting;
                    self.state_timestamp = None;
                }
                _ => {}
            }
        }

        Ok(None)
    }
}

impl ImportDialog {
    fn centered_rect(&self, percent_x: u16, percent_y: u16, area: Rect) -> Rect {
        let popup_layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Percentage((100 - percent_y) / 2),
                Constraint::Percentage(percent_y),
                Constraint::Percentage((100 - percent_y) / 2),
            ])
            .split(area);

        Layout::default()
            .direction(Direction::Horizontal)
            .constraints([
                Constraint::Percentage((100 - percent_x) / 2),
                Constraint::Percentage(percent_x),
                Constraint::Percentage((100 - percent_x) / 2),
            ])
            .split(popup_layout[1])[1]
    }
}
