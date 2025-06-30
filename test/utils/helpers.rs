use ratatui::{backend::TestBackend, Terminal, Frame};
use gemini_cli_manager::{
    models::{Extension, Profile, extension::ExtensionMetadata, profile::ProfileMetadata},
    storage::Storage,
};
use chrono::Utc;
use std::collections::HashMap;

/// Set up a test terminal with the specified dimensions
pub fn setup_test_terminal(width: u16, height: u16) -> Result<Terminal<TestBackend>, Box<dyn std::error::Error>> {
    let backend = TestBackend::new(width, height);
    Ok(Terminal::new(backend)?)
}

/// Render a component to a string for snapshot testing
pub fn render_to_string<F>(width: u16, height: u16, render_fn: F) -> Result<String, Box<dyn std::error::Error>>
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
            let cell = buffer.get(x, y);
            line.push_str(cell.symbol());
        }
        // Trim trailing whitespace from each line
        lines.push(line.trim_end().to_string());
    }
    // Remove trailing empty lines
    while lines.last().map_or(false, |line| line.is_empty()) {
        lines.pop();
    }
    lines.join("\n")
}

/// Create an in-memory storage for testing
pub fn create_test_storage() -> Storage {
    Storage::new_in_memory().expect("Failed to create in-memory storage")
}

/// Create a test extension with default values
pub fn create_test_extension(name: &str) -> Extension {
    Extension {
        id: name.to_lowercase().replace(' ', "-"),
        name: name.to_string(),
        version: "1.0.0".to_string(),
        description: Some(format!("Test extension: {}", name)),
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
        description: Some(format!("Test profile: {}", name)),
        extension_ids: vec![],
        environment_variables: HashMap::new(),
        working_directory: None,
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
        Extension {
            id: self.name.to_lowercase().replace(' ', "-"),
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
        Profile {
            id: self.name.to_lowercase().replace(' ', "-"),
            name: self.name,
            description: self.description,
            extension_ids: self.extension_ids,
            environment_variables: HashMap::new(),
            working_directory: None,
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
        "Expected to find '{}' in rendered output:\n{}",
        expected,
        content
    );
}

/// Assert that a terminal buffer does not contain specific text
pub fn assert_buffer_not_contains(terminal: &Terminal<TestBackend>, unexpected: &str) {
    let content = buffer_to_string(terminal.backend().buffer());
    assert!(
        !content.contains(unexpected),
        "Did not expect to find '{}' in rendered output:\n{}",
        unexpected,
        content
    );
}

/// Simulate a series of key events on an app
pub fn simulate_key_sequence<A>(app: &mut A, keys: Vec<crossterm::event::KeyCode>) 
where
    A: HandleKeyEvent,
{
    use crossterm::event::{KeyEvent, KeyEventKind};
    
    for key in keys {
        let event = KeyEvent::new(key, KeyEventKind::Press);
        app.handle_key_event(event);
    }
}

/// Trait for components that handle key events
pub trait HandleKeyEvent {
    fn handle_key_event(&mut self, event: crossterm::event::KeyEvent);
}