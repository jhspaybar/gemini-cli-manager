# Gemini Extension Structure Analysis & Core Workflow

## Executive Summary

After analyzing the actual Gemini extension examples in our project, there are significant discrepancies between our specifications and the real Gemini extension format. This document provides the correct understanding and proposes our core workflow.

## Actual Gemini Extension Structure

Based on the examples in `examples/extensions/`, Gemini extensions use:

### 1. **gemini-extension.json** (NOT settings.json)
```json
{
  "name": "extension-name",
  "version": "1.0.0",
  "description": "Extension description",
  "mcpServers": {
    "server-name": {
      // Either URL-based:
      "url": "http://localhost:8080/mcp/endpoint",
      // Or command-based:
      "command": "node",
      "args": ["server.js"],
      "cwd": "./servers",
      "env": {
        "VAR_NAME": "$ENV_VAR"
      },
      "timeout": 60000,
      "trust": false
    }
  },
  "contextFileName": "CONTEXT.md"  // e.g., "GITHUB.md", "DATABASE.md"
}
```

### 2. **Context File** (e.g., GITHUB.md, DATABASE.md)
- Markdown file providing context about the extension
- Contains usage instructions and configuration details
- Referenced by `contextFileName` in the manifest

### Key Observations:
- Much simpler than our specs suggest
- No complex dependency system
- No tools/templates/resources directories
- Focus on MCP server configuration
- Environment variables use `$VAR_NAME` syntax for substitution

## Independent State Management

✅ **Confirmed**: Our application stores state independently from Gemini using:

```rust
// Platform-specific directories via the 'directories' crate
// macOS: ~/Library/Application Support/gemini-cli-manager/
// Linux: ~/.local/share/gemini-cli-manager/
// Windows: C:\Users\User\AppData\Local\gemini-cli-manager\

pub fn get_data_dir() -> PathBuf {
    ProjectDirs::from("com", "kdheepak", "gemini-cli-manager")
        .data_local_dir()
}

pub fn get_config_dir() -> PathBuf {
    ProjectDirs::from("com", "kdheepak", "gemini-cli-manager")
        .config_local_dir()
}
```

## Core Workflow Design

### Overview
```
[Import/Create Extensions] → [Bundle into Profiles] → [Launch Gemini with Profile]
```

### 1. Extension Management (Our Storage)

We'll store extension references in our data directory:

```
~/.local/share/gemini-cli-manager/
├── extensions/
│   ├── github-tools/
│   │   ├── manifest.json      # Copy of gemini-extension.json
│   │   ├── context.md         # Copy of GITHUB.md
│   │   └── metadata.json      # Our metadata (when imported, etc.)
│   ├── database-tools/
│   │   └── ...
│   └── custom-extension/
│       └── ...
├── profiles/
│   ├── web-dev.json
│   ├── data-science.json
│   └── devops.json
└── config.json
```

### 2. Profile Structure (Our Concept)

Profiles bundle multiple extensions with environment variables:

```json
{
  "id": "web-dev-profile",
  "name": "Web Development",
  "description": "Frontend and backend web development tools",
  "extensions": [
    {
      "id": "github-tools",
      "enabled": true,
      "env_overrides": {
        "GITHUB_TOKEN": "${GITHUB_TOKEN_WORK}"
      }
    },
    {
      "id": "database-tools",
      "enabled": true,
      "env_overrides": {
        "DATABASE_URL": "postgresql://localhost/myapp"
      }
    }
  ],
  "global_env": {
    "NODE_ENV": "development",
    "API_BASE_URL": "http://localhost:3000"
  },
  "launch_config": {
    "working_dir": "~/projects/current",
    "gemini_args": ["--model", "gemini-2.0-flash"]
  }
}
```

### 3. Launch Process

When launching Gemini with a profile:

1. **Prepare Environment**
   ```rust
   // Collect all environment variables
   let mut env = HashMap::new();
   
   // Add global profile env vars
   env.extend(profile.global_env);
   
   // Add extension-specific env vars
   for ext in profile.extensions {
       if ext.enabled {
           env.extend(ext.env_overrides);
       }
   }
   ```

2. **Generate Temporary Gemini Config**
   ```rust
   // Create a temporary directory with enabled extensions
   let temp_dir = create_temp_gemini_config();
   
   // Copy enabled extensions to temp directory
   for ext in profile.extensions {
       if ext.enabled {
           copy_extension_to(ext, temp_dir);
       }
   }
   ```

3. **Launch Gemini**
   ```rust
   Command::new("gemini")
       .env_clear()
       .envs(env)
       .arg("--extensions-dir").arg(temp_dir)
       .args(profile.launch_config.gemini_args)
       .current_dir(profile.launch_config.working_dir)
       .spawn()
   ```

## Key Implementation Tasks

### Phase 1: Extension Management
1. **Import Extension**
   - Parse `gemini-extension.json`
   - Validate MCP server configurations
   - Copy to our storage with metadata
   - Handle environment variable placeholders

2. **Create Extension**
   - Provide UI for creating `gemini-extension.json`
   - Support both URL and command-based MCP servers
   - Generate context markdown file

### Phase 2: Profile Management
1. **Create Profile**
   - Select extensions to include
   - Configure environment variables
   - Set launch parameters

2. **Edit Profile**
   - Enable/disable extensions
   - Override environment variables per extension
   - Configure global environment

### Phase 3: Launch Integration
1. **Profile Activation**
   - Validate all required environment variables
   - Check extension availability
   - Prepare launch environment

2. **Gemini Launch**
   - Create temporary extension directory
   - Set up environment variables
   - Execute Gemini with proper arguments
   - Monitor process status

## Corrected Data Models

### Extension Model
```rust
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GeminiExtension {
    pub name: String,
    pub version: String,
    pub description: String,
    pub mcp_servers: HashMap<String, McpServerConfig>,
    pub context_file_name: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct McpServerConfig {
    // URL-based server
    pub url: Option<String>,
    
    // Command-based server
    pub command: Option<String>,
    pub args: Option<Vec<String>>,
    pub cwd: Option<String>,
    
    // Common fields
    pub env: Option<HashMap<String, String>>,
    pub timeout: Option<u64>,
    pub trust: Option<bool>,
}

// Our metadata
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExtensionMetadata {
    pub id: String,
    pub imported_at: DateTime<Utc>,
    pub source: ExtensionSource,
    pub context_content: Option<String>,
}
```

### Profile Model
```rust
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Profile {
    pub id: String,
    pub name: String,
    pub description: Option<String>,
    pub extensions: Vec<ProfileExtension>,
    pub global_env: HashMap<String, String>,
    pub launch_config: LaunchConfig,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProfileExtension {
    pub extension_id: String,
    pub enabled: bool,
    pub env_overrides: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LaunchConfig {
    pub working_dir: Option<PathBuf>,
    pub gemini_args: Vec<String>,
    pub pre_launch_commands: Vec<String>,
}
```

## Summary

The core workflow is:
1. **Import/Create** Gemini extensions (storing them in our app's data directory)
2. **Bundle** multiple extensions into profiles with custom environment configurations
3. **Launch** Gemini with a selected profile, temporarily providing the enabled extensions

This approach:
- ✅ Maintains independent state from Gemini
- ✅ Allows multiple profile configurations
- ✅ Supports environment variable management
- ✅ Enables extension organization without modifying Gemini's actual extension directory
- ✅ Provides a clean, reversible way to manage configurations