use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Represents a Gemini CLI extension based on gemini-extension.json
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Extension {
    /// Our internal ID
    pub id: String,

    /// Extension name from gemini-extension.json
    pub name: String,

    /// Extension version
    pub version: String,

    /// Optional description (from our metadata)
    pub description: Option<String>,

    /// MCP servers defined in the extension
    pub mcp_servers: HashMap<String, McpServerConfig>,

    /// Context file name (e.g., "GITHUB.md", "DATABASE.md")
    pub context_file_name: Option<String>,

    /// Content of the context file
    pub context_content: Option<String>,

    /// Our metadata
    pub metadata: ExtensionMetadata,
}

/// MCP (Model Context Protocol) server configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct McpServerConfig {
    /// Command to execute for command-based servers
    pub command: Option<String>,

    /// Arguments for the command
    pub args: Option<Vec<String>>,

    /// Working directory
    pub cwd: Option<String>,

    /// Environment variables (with $VAR_NAME syntax)
    pub env: Option<HashMap<String, String>>,

    /// Timeout in milliseconds
    pub timeout: Option<u64>,

    /// Whether to trust this server
    pub trust: Option<bool>,
}

/// Our metadata for tracking extensions
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExtensionMetadata {
    /// When the extension was imported
    pub imported_at: DateTime<Utc>,

    /// Original source path
    pub source_path: Option<String>,

    /// User-defined tags
    pub tags: Vec<String>,
}

impl Extension {
    // Mock data methods removed - extensions should be imported from actual extension packages
}
