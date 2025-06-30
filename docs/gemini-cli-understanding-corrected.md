# Corrected Understanding of Gemini CLI & Updated Implementation Strategy

## Key Discoveries from Official Documentation

After analyzing the official Gemini CLI documentation, several critical insights have emerged that require us to adjust our approach:

### 1. Gemini CLI is an Interactive REPL

**NOT a traditional CLI with command-line arguments!**

- Gemini CLI is launched by simply running `gemini` or `npx @google/gemini-cli`
- It provides an interactive REPL (Read-Eval-Print Loop) environment
- No command-line flags for specifying extensions, profiles, or configurations
- Configuration is done through files in specific locations

### 2. Extension System Architecture

Extensions are loaded from two locations (in priority order):
1. **Workspace**: `<current_directory>/.gemini/extensions/`
2. **Home**: `~/.gemini/extensions/`

Extension structure:
```
.gemini/extensions/
└── extension-name/
    ├── gemini-extension.json    # Required: Extension manifest
    └── GEMINI.md               # Optional: Context file (or custom name)
```

### 3. Extension Configuration

**gemini-extension.json** structure:
```json
{
  "name": "extension-name",
  "version": "1.0.0",
  "mcpServers": {
    "server-name": {
      // Command-based server:
      "command": "node",
      "args": ["server.js"],
      "cwd": "./servers",
      "env": {
        "API_KEY": "$MY_API_KEY"  // Environment variable substitution
      },
      "timeout": 60000,
      "trust": false
    }
  },
  "contextFileName": "CUSTOM.md"  // Optional, defaults to "GEMINI.md"
}
```

### 4. Settings Override System

**settings.json** (in `.gemini/settings.json`) can override MCP server configurations:

```json
{
  "mcpServers": {
    "extension-name/server-name": {
      "env": {
        "API_KEY": "$WORK_API_KEY"  // Override extension's env var
      }
    }
  }
}
```

### 5. Environment Variables

- Automatically loaded from `.env` files in current directory and parents
- Referenced using `$VAR_NAME` or `${VAR_NAME}` syntax
- Required: `GEMINI_API_KEY` for authentication

## Revised Implementation Strategy

Given these insights, our Gemini CLI Manager needs to work differently:

### Core Concept: Profile-Based Workspace Management

Since Gemini CLI loads extensions from the working directory, our profiles will manage **workspace configurations** rather than launching Gemini with different arguments.

## New Workflow Design

```
[Create/Import Extensions] → [Define Profiles] → [Setup Workspace] → [Launch Gemini]
```

### 1. Extension Management

Store extensions in our app's data directory:
```
~/.local/share/gemini-cli-manager/
├── extensions/
│   ├── github-tools/
│   │   ├── gemini-extension.json
│   │   ├── GITHUB.md
│   │   └── metadata.json          # Our tracking data
│   └── database-tools/
│       └── ...
└── profiles/
    ├── web-dev.json
    └── data-science.json
```

### 2. Profile Structure (Updated)

```json
{
  "id": "web-dev-profile",
  "name": "Web Development",
  "description": "Frontend and backend web development setup",
  "workspace_template": {
    "extensions": [
      {
        "id": "github-tools",
        "enabled": true,
        "server_overrides": {
          "github-api": {
            "env": {
              "GITHUB_TOKEN": "$GITHUB_TOKEN_WORK"
            }
          }
        }
      },
      {
        "id": "database-tools",
        "enabled": true
      }
    ],
    "settings": {
      "theme": "GitHub Dark",
      "telemetry": {
        "enabled": false
      },
      "mcpServers": {
        // Additional MCP servers not in extensions
        "custom-server": {
          "command": "python",
          "args": ["custom_server.py"]
        }
      }
    },
    "env_file_content": [
      "GITHUB_TOKEN_WORK=ghp_xxx",
      "DATABASE_URL=postgresql://localhost/myapp",
      "NODE_ENV=development"
    ]
  },
  "launch_directory": "~/projects/current"
}
```

### 3. Profile Activation Process

When activating a profile:

```rust
pub fn activate_profile(profile: &Profile) -> Result<()> {
    // 1. Determine target directory
    let workspace_dir = expand_path(&profile.launch_directory);
    
    // 2. Create .gemini directory structure
    let gemini_dir = workspace_dir.join(".gemini");
    let extensions_dir = gemini_dir.join("extensions");
    fs::create_dir_all(&extensions_dir)?;
    
    // 3. Copy enabled extensions
    for ext_ref in &profile.workspace_template.extensions {
        if ext_ref.enabled {
            let ext = load_extension(&ext_ref.id)?;
            let target = extensions_dir.join(&ext.name);
            
            // Copy extension files
            copy_extension(&ext, &target)?;
            
            // Apply any server overrides
            if !ext_ref.server_overrides.is_empty() {
                apply_overrides(&target, &ext_ref.server_overrides)?;
            }
        }
    }
    
    // 4. Create settings.json
    let settings_path = gemini_dir.join("settings.json");
    fs::write(&settings_path, serde_json::to_string_pretty(&profile.workspace_template.settings)?)?;
    
    // 5. Create .env file
    let env_path = workspace_dir.join(".env");
    let env_content = profile.workspace_template.env_file_content.join("\n");
    fs::write(&env_path, env_content)?;
    
    // 6. Launch Gemini CLI
    Command::new("gemini")
        .current_dir(&workspace_dir)
        .spawn()?;
    
    Ok(())
}
```

### 4. Quick Switch Feature

For rapidly switching between profiles:

```rust
pub struct WorkspaceManager {
    active_workspaces: HashMap<String, PathBuf>,
}

impl WorkspaceManager {
    pub fn quick_switch(&mut self, profile_id: &str) -> Result<()> {
        // Check if workspace already set up
        if let Some(workspace) = self.active_workspaces.get(profile_id) {
            // Just launch Gemini in existing workspace
            Command::new("gemini")
                .current_dir(workspace)
                .spawn()?;
        } else {
            // Set up new workspace
            let profile = load_profile(profile_id)?;
            let workspace = setup_workspace(&profile)?;
            self.active_workspaces.insert(profile_id.to_string(), workspace.clone());
            
            Command::new("gemini")
                .current_dir(&workspace)
                .spawn()?;
        }
        Ok(())
    }
}
```

## Updated Core Features

### Phase 1: Extension & Profile Management
1. **Extension Import/Create**
   - Parse `gemini-extension.json`
   - Store in app's extension library
   - Track metadata (source, import date, etc.)

2. **Profile Creation**
   - Select extensions to include
   - Configure environment variables
   - Set workspace directory
   - Define settings.json overrides

3. **Workspace Setup**
   - Create `.gemini` directory structure
   - Copy selected extensions
   - Generate `settings.json`
   - Create `.env` file

### Phase 2: Advanced Features
1. **Workspace Templates**
   - Pre-configured directory structures
   - Project scaffolding
   - Git initialization

2. **Multi-Workspace Management**
   - Track multiple active workspaces
   - Quick switching between setups
   - Workspace cleanup utilities

3. **Extension Development**
   - Create new extensions
   - Test MCP servers
   - Validate configurations

### Phase 3: Team Features
1. **Profile Sharing**
   - Export profiles as shareable configs
   - Import team profiles
   - Version control integration

2. **Extension Registry**
   - Browse available extensions
   - One-click installation
   - Dependency management

## Key Differences from Original Specs

1. **No command-line arguments** - Gemini CLI doesn't accept configuration via CLI args
2. **Workspace-based** - Extensions must be in workspace `.gemini` directory
3. **Interactive REPL** - Not a traditional command-line tool
4. **Settings override** - Use `settings.json` for environment-specific configs
5. **Simple extension format** - Much simpler than our original specs suggested

## Benefits of This Approach

1. **Aligns with Gemini's design** - Works with how Gemini actually loads extensions
2. **Profile isolation** - Each profile can have its own workspace
3. **Quick switching** - Fast context switching between projects
4. **Team friendly** - Easy to share and version control profiles
5. **Non-invasive** - Doesn't modify global Gemini installation

## Summary

Our Gemini CLI Manager will be a **workspace configuration manager** that:
- Maintains a library of extensions
- Defines profiles as workspace templates
- Sets up project-specific `.gemini` directories
- Manages environment variables via `.env` files
- Launches Gemini CLI in configured workspaces

This approach respects Gemini CLI's architecture while providing the profile-based workflow we want.