use gemini_cli_manager::models::extension::{Extension, ExtensionMetadata, McpServerConfig};
use chrono::Utc;
use std::collections::HashMap;

/// Test fixtures for MCP server configurations
pub struct McpFixtures;

impl McpFixtures {
    /// Create a simple echo server configuration
    pub fn echo_server_simple() -> McpServerConfig {
        McpServerConfig {
            command: Some("node".to_string()),
            args: Some(vec!["echo-server.js".to_string()]),
            cwd: None,
            env: None,
            timeout: None,
            trust: Some(true),
            url: None,
        }
    }
    
    /// Create a Docker-based echo server
    pub fn echo_server_docker() -> McpServerConfig {
        McpServerConfig {
            command: Some("docker".to_string()),
            args: Some(vec![
                "run".to_string(),
                "-i".to_string(),
                "--rm".to_string(),
                "aadversteeg/echo-mcp-server:latest".to_string(),
            ]),
            cwd: None,
            env: Some({
                let mut env = HashMap::new();
                env.insert("MessageFormat".to_string(), "Echo: {message}".to_string());
                env
            }),
            timeout: Some(30000),
            trust: Some(false),
            url: None,
        }
    }
    
    /// Create a Python echo server
    pub fn echo_server_python() -> McpServerConfig {
        McpServerConfig {
            command: Some("python".to_string()),
            args: Some(vec![
                "-m".to_string(),
                "mcp_echo_server".to_string(),
            ]),
            cwd: Some("./servers".to_string()),
            env: Some({
                let mut env = HashMap::new();
                env.insert("ECHO_PREFIX".to_string(), "[ECHO]".to_string());
                env.insert("PYTHONPATH".to_string(), "$PYTHONPATH:./lib".to_string());
                env
            }),
            timeout: Some(60000),
            trust: Some(true),
            url: None,
        }
    }
    
    /// Create a command-based server with environment configuration
    pub fn env_configured_server() -> McpServerConfig {
        McpServerConfig {
            command: Some("node".to_string()),
            args: Some(vec!["server.js".to_string()]),
            cwd: None,
            env: Some({
                let mut env = HashMap::new();
                env.insert("API_KEY".to_string(), "$MCP_API_KEY".to_string());
                env
            }),
            timeout: Some(10000),
            trust: Some(false),
            url: None,
        }
    }
    
    /// Create an extension with echo server
    pub fn echo_extension() -> Extension {
        Extension {
            id: "echo-test".to_string(),
            name: "Echo Test Extension".to_string(),
            version: "1.0.0".to_string(),
            description: Some("Test extension with echo MCP server".to_string()),
            mcp_servers: {
                let mut servers = HashMap::new();
                servers.insert("echo".to_string(), Self::echo_server_simple());
                servers
            },
            context_file_name: Some("GEMINI.md".to_string()),
            context_content: Some(Self::echo_context_content()),
            metadata: ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: Some("/test/extensions/echo-test".to_string()),
                tags: vec!["test".to_string(), "echo".to_string()],
            },
        }
    }
    
    /// Create an extension with multiple servers
    pub fn multi_server_extension() -> Extension {
        Extension {
            id: "multi-server-test".to_string(),
            name: "Multi-Server Test".to_string(),
            version: "2.0.0".to_string(),
            description: Some("Extension with multiple MCP servers".to_string()),
            mcp_servers: {
                let mut servers = HashMap::new();
                servers.insert("echo".to_string(), Self::echo_server_simple());
                servers.insert("python-echo".to_string(), Self::echo_server_python());
                servers.insert("api-server".to_string(), Self::env_configured_server());
                servers
            },
            context_file_name: Some("MULTI_SERVER.md".to_string()),
            context_content: Some(Self::multi_server_context()),
            metadata: ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec!["test".to_string(), "multi-server".to_string()],
            },
        }
    }
    
    /// Create an extension with only context (no servers)
    pub fn context_only_extension() -> Extension {
        Extension {
            id: "context-only".to_string(),
            name: "Context Only Extension".to_string(),
            version: "1.0.0".to_string(),
            description: Some("Extension providing only context instructions".to_string()),
            mcp_servers: HashMap::new(),
            context_file_name: Some("INSTRUCTIONS.md".to_string()),
            context_content: Some(Self::context_only_content()),
            metadata: ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec!["test".to_string(), "context".to_string()],
            },
        }
    }
    
    /// Create a fully populated extension
    pub fn full_featured_extension() -> Extension {
        Extension {
            id: "full-featured".to_string(),
            name: "Full Featured Extension".to_string(),
            version: "3.1.4".to_string(),
            description: Some("Extension demonstrating all features".to_string()),
            mcp_servers: {
                let mut servers = HashMap::new();
                servers.insert("main-server".to_string(), McpServerConfig {
                    command: Some("npx".to_string()),
                    args: Some(vec![
                        "-y".to_string(),
                        "@example/mcp-server".to_string(),
                        "--config".to_string(),
                        "production.json".to_string(),
                    ]),
                    cwd: Some("$HOME/.mcp/servers".to_string()),
                    env: Some({
                        let mut env = HashMap::new();
                        env.insert("NODE_ENV".to_string(), "production".to_string());
                        env.insert("LOG_LEVEL".to_string(), "info".to_string());
                        env.insert("API_TOKEN".to_string(), "$MCP_API_TOKEN".to_string());
                        env
                    }),
                    timeout: Some(120000),
                    trust: Some(true),
                    url: None,
                });
                servers
            },
            context_file_name: Some("ADVANCED.md".to_string()),
            context_content: Some(Self::advanced_context_content()),
            metadata: ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: Some("/opt/extensions/full-featured".to_string()),
                tags: vec![
                    "test".to_string(),
                    "advanced".to_string(),
                    "production".to_string(),
                ],
            },
        }
    }
    
    /// Echo server context content
    fn echo_context_content() -> String {
        r#"# Echo Test Extension

This extension provides a simple echo server for testing MCP functionality.

## Available Commands

### echo
Returns the input message prefixed with "Echo: "

Example:
- Input: "Hello, World!"
- Output: "Echo: Hello, World!"

### ping
Always returns "pong"

## Configuration

No additional configuration required.

## Testing

Use this extension to verify that MCP server communication is working correctly."#.to_string()
    }
    
    /// Multi-server context content
    fn multi_server_context() -> String {
        r#"# Multi-Server Test Extension

This extension demonstrates multiple MCP servers working together.

## Servers

### echo
Basic echo functionality using Node.js

### python-echo
Python-based echo server with prefix configuration

### api-server
Command-based server with environment configuration

## Environment Variables

- `ECHO_PREFIX`: Prefix for Python echo responses
- `API_KEY`: Authentication for API server
- `PYTHONPATH`: Python module search path

## Usage

Each server operates independently and can be called separately."#.to_string()
    }
    
    /// Context-only content
    fn context_only_content() -> String {
        r#"# Context Only Extension

This extension provides instructions without any MCP servers.

## Purpose

Sometimes you just need to provide context or instructions to Gemini without
adding any additional tools or servers.

## Instructions

When working with this project:
1. Follow the coding standards in CONTRIBUTING.md
2. Run tests before committing
3. Update documentation for any API changes

## Project Structure

```
src/
├── models/      # Data models
├── components/  # UI components
├── storage/     # Persistence layer
└── theme/       # Theme management
```

## Best Practices

- Always use semantic commit messages
- Keep functions small and focused
- Write tests for new functionality"#.to_string()
    }
    
    /// Advanced context content
    fn advanced_context_content() -> String {
        r#"# Advanced Full-Featured Extension

This extension demonstrates all available MCP features.

## Features

### Production-Ready Server
- Configurable via environment variables
- Timeout management
- Trust settings for security
- Custom working directory

### Environment Configuration
- `NODE_ENV`: Set to "production" for optimized performance
- `LOG_LEVEL`: Controls server logging verbosity
- `API_TOKEN`: Secure authentication token

### Advanced Commands

#### analyze
Performs deep analysis of provided data

#### transform
Transforms data between formats (JSON, YAML, XML)

#### validate
Validates data against schemas

## Security

This server runs with trust enabled, allowing full system access.
Use with caution in production environments.

## Performance

- Timeout: 120 seconds for long-running operations
- Optimized for production workloads
- Supports concurrent requests"#.to_string()
    }
}

/// Helper to create test gemini-extension.json content
pub fn create_extension_json(extension: &Extension) -> serde_json::Value {
    serde_json::json!({
        "name": extension.name,
        "version": extension.version,
        "mcpServers": extension.mcp_servers,
    })
}

/// Verify that an extension's JSON is valid for Gemini
pub fn validate_extension_json(extension: &Extension) -> Result<(), String> {
    // Check required fields
    if extension.name.is_empty() {
        return Err("Extension name is required".to_string());
    }
    
    if extension.version.is_empty() {
        return Err("Extension version is required".to_string());
    }
    
    // Validate MCP servers
    for (name, server) in &extension.mcp_servers {
        // Must have command
        if server.command.is_none() {
            return Err(format!("Server '{}' must have a command", name));
        }
        
        // If command-based, args should be present (but can be empty)
        if server.command.is_some() && server.args.is_none() {
            return Err(format!("Server '{}' with command should have args field", name));
        }
    }
    
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::{McpFixtures, validate_extension_json};
    
    #[test]
    fn test_echo_extension_valid() {
        let ext = McpFixtures::echo_extension();
        assert!(validate_extension_json(&ext).is_ok());
    }
    
    #[test]
    fn test_multi_server_extension_valid() {
        let ext = McpFixtures::multi_server_extension();
        assert!(validate_extension_json(&ext).is_ok());
        assert_eq!(ext.mcp_servers.len(), 3);
    }
    
    #[test]
    fn test_extension_json_format() {
        let ext = McpFixtures::echo_extension();
        let json = super::create_extension_json(&ext);
        
        assert_eq!(json["name"], "Echo Test Extension");
        assert_eq!(json["version"], "1.0.0");
        assert!(json["mcpServers"]["echo"].is_object());
    }
}