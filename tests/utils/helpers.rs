use chrono::Utc;
use gemini_cli_manager::{
    models::{
        Extension, Profile,
        extension::ExtensionMetadata,
        profile::{LaunchConfig, ProfileMetadata},
    },
    storage::Storage,
};
use ratatui::{Frame, Terminal, backend::TestBackend};
use std::collections::HashMap;

/// Set up a test terminal with the specified dimensions
pub fn setup_test_terminal(
    width: u16,
    height: u16,
) -> Result<Terminal<TestBackend>, Box<dyn std::error::Error>> {
    let backend = TestBackend::new(width, height);
    Ok(Terminal::new(backend)?)
}

/// Render a component to a string for snapshot testing
pub fn render_to_string<F>(
    width: u16,
    height: u16,
    render_fn: F,
) -> Result<String, Box<dyn std::error::Error>>
where
    F: FnOnce(&mut Frame),
{
    let mut terminal = setup_test_terminal(width, height)?;
    terminal.draw(render_fn)?;

    let buffer = terminal.backend().buffer();
    Ok(buffer_to_string(buffer))
}

/// Convert a buffer to a string representation
pub fn buffer_to_string(buffer: &ratatui::buffer::Buffer) -> String {
    let mut lines = vec![];
    for y in 0..buffer.area.height {
        let mut line = String::new();
        for x in 0..buffer.area.width {
            let cell = &buffer[(x, y)];
            line.push_str(cell.symbol());
        }
        // Trim trailing whitespace from each line
        lines.push(line.trim_end().to_string());
    }
    // Remove trailing empty lines
    while lines.last().is_some_and(|line| line.is_empty()) {
        lines.pop();
    }
    lines.join("\n")
}

/// Create a test extension with default values
pub fn create_test_extension(name: &str) -> Extension {
    Extension {
        id: name.to_lowercase().replace(' ', "-"),
        name: name.to_string(),
        version: "1.0.0".to_string(),
        description: Some(format!("Test extension: {name}")),
        mcp_servers: HashMap::new(),
        context_file_name: None,
        context_content: None,
        metadata: ExtensionMetadata {
            imported_at: Utc::now(),
            source_path: None,
            tags: vec!["test".to_string()],
        },
    }
}

/// Create a test profile with default values
pub fn create_test_profile(name: &str) -> Profile {
    Profile {
        id: name.to_lowercase().replace(' ', "-"),
        name: name.to_string(),
        description: Some(format!("Test profile: {name}")),
        extension_ids: vec![],
        environment_variables: HashMap::new(),
        working_directory: None,
        launch_config: LaunchConfig::default(),
        metadata: ProfileMetadata {
            created_at: Utc::now(),
            updated_at: Utc::now(),
            tags: vec!["test".to_string()],
            is_default: false,
            icon: None,
        },
    }
}

/// Builder for creating test extensions with specific properties
pub struct ExtensionBuilder {
    name: String,
    version: String,
    description: Option<String>,
    tags: Vec<String>,
    mcp_servers: HashMap<String, gemini_cli_manager::models::extension::McpServerConfig>,
}

impl ExtensionBuilder {
    pub fn new(name: &str) -> Self {
        Self {
            name: name.to_string(),
            version: "1.0.0".to_string(),
            description: None,
            tags: vec![],
            mcp_servers: HashMap::new(),
        }
    }

    pub fn with_version(mut self, version: &str) -> Self {
        self.version = version.to_string();
        self
    }

    pub fn with_description(mut self, description: &str) -> Self {
        self.description = Some(description.to_string());
        self
    }

    pub fn with_tags(mut self, tags: Vec<&str>) -> Self {
        self.tags = tags.into_iter().map(|s| s.to_string()).collect();
        self
    }

    pub fn build(self) -> Extension {
        // Generate ID using the same logic as the application
        let id = self
            .name
            .to_lowercase()
            .chars()
            .map(|c| {
                if c.is_alphanumeric() {
                    c
                } else if c == ' ' || c == '-' || c == '_' || c == '.' {
                    '-'
                } else {
                    // Remove other special characters
                    '\0'
                }
            })
            .filter(|c| *c != '\0')
            .collect::<String>()
            .split('-')
            .filter(|s| !s.is_empty())
            .collect::<Vec<_>>()
            .join("-");

        Extension {
            id,
            name: self.name,
            version: self.version,
            description: self.description,
            mcp_servers: self.mcp_servers,
            context_file_name: None,
            context_content: None,
            metadata: ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: self.tags,
            },
        }
    }
}

/// Builder for creating test profiles with specific properties
pub struct ProfileBuilder {
    name: String,
    description: Option<String>,
    extension_ids: Vec<String>,
    tags: Vec<String>,
    is_default: bool,
}

impl ProfileBuilder {
    pub fn new(name: &str) -> Self {
        Self {
            name: name.to_string(),
            description: None,
            extension_ids: vec![],
            tags: vec![],
            is_default: false,
        }
    }

    pub fn with_description(mut self, description: &str) -> Self {
        self.description = Some(description.to_string());
        self
    }

    pub fn with_extensions(mut self, extension_ids: Vec<&str>) -> Self {
        self.extension_ids = extension_ids.into_iter().map(|s| s.to_string()).collect();
        self
    }

    pub fn with_tags(mut self, tags: Vec<&str>) -> Self {
        self.tags = tags.into_iter().map(|s| s.to_string()).collect();
        self
    }

    pub fn as_default(mut self) -> Self {
        self.is_default = true;
        self
    }

    pub fn build(self) -> Profile {
        // Generate ID using the same logic as the application
        let id = self
            .name
            .to_lowercase()
            .chars()
            .map(|c| {
                if c.is_alphanumeric() {
                    c
                } else if c == ' ' || c == '-' || c == '_' || c == '.' {
                    '-'
                } else {
                    // Remove other special characters
                    '\0'
                }
            })
            .filter(|c| *c != '\0')
            .collect::<String>()
            .split('-')
            .filter(|s| !s.is_empty())
            .collect::<Vec<_>>()
            .join("-");

        Profile {
            id,
            name: self.name,
            description: self.description,
            extension_ids: self.extension_ids,
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: LaunchConfig::default(),
            metadata: ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: self.tags,
                is_default: self.is_default,
                icon: None,
            },
        }
    }
}

/// Assert that a terminal buffer contains specific text
pub fn assert_buffer_contains(terminal: &Terminal<TestBackend>, expected: &str) {
    let content = buffer_to_string(terminal.backend().buffer());
    assert!(
        content.contains(expected),
        "Expected to find '{expected}' in rendered output:\n{content}"
    );
}

/// Assert that a terminal buffer does not contain specific text
pub fn assert_buffer_not_contains(terminal: &Terminal<TestBackend>, unexpected: &str) {
    let content = buffer_to_string(terminal.backend().buffer());
    assert!(
        !content.contains(unexpected),
        "Did not expect to find '{unexpected}' in rendered output:\n{content}"
    );
}

/// Simulate a series of key events on an app
pub fn simulate_key_sequence<A>(app: &mut A, keys: Vec<crossterm::event::KeyCode>)
where
    A: HandleKeyEvent,
{
    use crossterm::event::{KeyEvent, KeyEventKind};

    for key in keys {
        use crossterm::event::KeyModifiers;
        let event = KeyEvent {
            code: key,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        };
        app.handle_key_event(event);
    }
}

/// Trait for components that handle key events
pub trait HandleKeyEvent {
    fn handle_key_event(&mut self, event: crossterm::event::KeyEvent);
}

/// Helper for verifying workspace setup
pub struct WorkspaceVerifier;

impl WorkspaceVerifier {
    /// Verify that a workspace directory has the correct structure
    pub fn verify_workspace_structure(workspace_dir: &std::path::Path) -> Result<(), String> {
        // Check workspace exists
        if !workspace_dir.exists() {
            return Err("Workspace directory does not exist".to_string());
        }

        // Check .gemini directory
        let gemini_dir = workspace_dir.join(".gemini");
        if !gemini_dir.exists() {
            return Err(".gemini directory not found".to_string());
        }

        // Check extensions directory
        let extensions_dir = gemini_dir.join("extensions");
        if !extensions_dir.exists() {
            return Err(".gemini/extensions directory not found".to_string());
        }

        Ok(())
    }

    /// Verify an extension is properly installed
    pub fn verify_extension_installed(
        workspace_dir: &std::path::Path,
        extension_id: &str,
    ) -> Result<(), String> {
        use std::fs;

        let ext_dir = workspace_dir
            .join(".gemini")
            .join("extensions")
            .join(extension_id);

        if !ext_dir.exists() {
            return Err(format!("Extension directory '{extension_id}' not found"));
        }

        // Check for gemini-extension.json
        let config_file = ext_dir.join("gemini-extension.json");
        if !config_file.exists() {
            return Err("gemini-extension.json not found".to_string());
        }

        // Validate JSON content
        let content =
            fs::read_to_string(&config_file).map_err(|e| format!("Failed to read config: {e}"))?;

        let json: serde_json::Value =
            serde_json::from_str(&content).map_err(|e| format!("Invalid JSON: {e}"))?;

        // Check required fields
        if json.get("name").is_none() {
            return Err("Missing 'name' field in extension config".to_string());
        }

        if json.get("version").is_none() {
            return Err("Missing 'version' field in extension config".to_string());
        }

        Ok(())
    }

    /// Verify context file exists and has content
    pub fn verify_context_file(
        workspace_dir: &std::path::Path,
        extension_id: &str,
        filename: &str,
    ) -> Result<String, String> {
        use std::fs;

        let context_file = workspace_dir
            .join(".gemini")
            .join("extensions")
            .join(extension_id)
            .join(filename);

        if !context_file.exists() {
            return Err(format!("Context file '{filename}' not found"));
        }

        let content = fs::read_to_string(&context_file)
            .map_err(|e| format!("Failed to read context file: {e}"))?;

        if content.trim().is_empty() {
            return Err("Context file is empty".to_string());
        }

        Ok(content)
    }
}

/// Create a test storage backed by a temporary directory
pub fn create_test_storage() -> Storage {
    let temp_dir = tempfile::TempDir::new().expect("Failed to create temp dir");
    let storage = Storage::with_data_dir(temp_dir.path().to_path_buf());
    // Create directories without initializing mock data
    std::fs::create_dir_all(temp_dir.path().join("extensions"))
        .expect("Failed to create extensions dir");
    std::fs::create_dir_all(temp_dir.path().join("profiles"))
        .expect("Failed to create profiles dir");
    // We leak the temp dir here so it persists for the test duration
    // In a real test framework, we'd manage this lifecycle better
    std::mem::forget(temp_dir);
    storage
}

/// Create a temporary test storage with file system backing and return the temp dir
pub fn create_temp_storage() -> (Storage, tempfile::TempDir) {
    let temp_dir = tempfile::TempDir::new().expect("Failed to create temp dir");
    let storage = Storage::with_data_dir(temp_dir.path().to_path_buf());
    // Create directories without initializing mock data
    std::fs::create_dir_all(temp_dir.path().join("extensions"))
        .expect("Failed to create extensions dir");
    std::fs::create_dir_all(temp_dir.path().join("profiles"))
        .expect("Failed to create profiles dir");
    (storage, temp_dir)
}
