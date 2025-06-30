use std::collections::HashMap;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

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
    /// URL-based server endpoint
    pub url: Option<String>,
    
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

/// Simplified MCP server info for display
#[derive(Debug, Clone)]
pub struct McpServer {
    pub name: String,
    pub server_type: McpServerType,
}

#[derive(Debug, Clone)]
pub enum McpServerType {
    Url(String),
    Command(String),
}

impl Extension {
    /// Create mock extensions for testing
    pub fn mock_extensions() -> Vec<Extension> {
        vec![
            Extension {
                id: "github-tools".to_string(),
                name: "github-tools".to_string(),
                version: "1.0.0".to_string(),
                description: Some("GitHub integration tools for Gemini CLI".to_string()),
                mcp_servers: {
                    let mut servers = HashMap::new();
                    servers.insert(
                        "github-api".to_string(),
                        McpServerConfig {
                            url: None,
                            command: Some("npx".to_string()),
                            args: Some(vec!["-y".to_string(), "@modelcontextprotocol/server-github".to_string()]),
                            cwd: None,
                            env: Some({
                                let mut env = HashMap::new();
                                env.insert("GITHUB_TOKEN".to_string(), "$GITHUB_TOKEN".to_string());
                                env
                            }),
                            timeout: None,
                            trust: Some(false),
                        },
                    );
                    servers
                },
                context_file_name: Some("GITHUB.md".to_string()),
                context_content: Some("# GitHub Tools\n\nProvides GitHub integration.".to_string()),
                metadata: ExtensionMetadata {
                    imported_at: Utc::now(),
                    source_path: Some("~/.gemini/extensions/github-tools".to_string()),
                    tags: vec!["vcs".to_string(), "github".to_string()],
                },
            },
            Extension {
                id: "database-tools".to_string(),
                name: "database-tools".to_string(),
                version: "2.0.0".to_string(),
                description: Some("Database query and management tools".to_string()),
                mcp_servers: {
                    let mut servers = HashMap::new();
                    servers.insert(
                        "postgres-server".to_string(),
                        McpServerConfig {
                            url: Some("http://localhost:8080/mcp/postgres".to_string()),
                            command: None,
                            args: None,
                            cwd: None,
                            env: Some({
                                let mut env = HashMap::new();
                                env.insert("DATABASE_URL".to_string(), "$DATABASE_URL".to_string());
                                env
                            }),
                            timeout: Some(60000),
                            trust: Some(false),
                        },
                    );
                    servers.insert(
                        "sqlite-server".to_string(),
                        McpServerConfig {
                            url: None,
                            command: Some("python".to_string()),
                            args: Some(vec!["-m".to_string(), "mcp_servers.sqlite".to_string()]),
                            cwd: Some("./servers".to_string()),
                            env: Some({
                                let mut env = HashMap::new();
                                env.insert("SQLITE_PATH".to_string(), "$HOME/data/sqlite".to_string());
                                env
                            }),
                            timeout: None,
                            trust: None,
                        },
                    );
                    servers
                },
                context_file_name: Some("DATABASE.md".to_string()),
                context_content: Some("# Database Tools\n\nDatabase management and query tools.".to_string()),
                metadata: ExtensionMetadata {
                    imported_at: Utc::now(),
                    source_path: Some("examples/extensions/database-tools".to_string()),
                    tags: vec!["database".to_string(), "sql".to_string()],
                },
            },
            Extension {
                id: "filesystem-enhanced".to_string(),
                name: "filesystem-enhanced".to_string(),
                version: "1.0.0".to_string(),
                description: Some("Enhanced file system operations".to_string()),
                mcp_servers: {
                    let mut servers = HashMap::new();
                    servers.insert(
                        "fs-server".to_string(),
                        McpServerConfig {
                            url: None,
                            command: Some("node".to_string()),
                            args: Some(vec!["fs-server.js".to_string()]),
                            cwd: None,
                            env: None,
                            timeout: None,
                            trust: Some(true),
                        },
                    );
                    servers
                },
                context_file_name: None,
                context_content: None,
                metadata: ExtensionMetadata {
                    imported_at: Utc::now(),
                    source_path: None,
                    tags: vec!["filesystem".to_string(), "tools".to_string()],
                },
            },
        ]
    }
    
    /// Get a list of MCP servers for display
    pub fn get_mcp_servers(&self) -> Vec<McpServer> {
        self.mcp_servers
            .iter()
            .map(|(name, config)| {
                let server_type = if let Some(url) = &config.url {
                    McpServerType::Url(url.clone())
                } else if let Some(cmd) = &config.command {
                    McpServerType::Command(cmd.clone())
                } else {
                    McpServerType::Command("unknown".to_string())
                };
                
                McpServer {
                    name: name.clone(),
                    server_type,
                }
            })
            .collect()
    }
}