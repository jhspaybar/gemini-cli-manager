# Revised Implementation Roadmap for Gemini CLI Manager

## Overview

Based on our corrected understanding of Gemini CLI, this roadmap outlines the implementation phases for a **workspace configuration manager** that enables profile-based workflows.

## Core Concept

Since Gemini CLI loads extensions from the current working directory's `.gemini/extensions/` folder, our application will:
1. Maintain a library of extensions
2. Define profiles as workspace configurations
3. Set up workspaces with selected extensions
4. Launch Gemini CLI in configured workspaces

## Phase 1: Foundation (MVP)

### 1.1 Extension Library Management
**Priority**: Critical  
**Status**: Not Started

#### Features:
- Import extensions from existing directories
- Create new extensions with UI
- Store extensions in app's data directory
- Parse and validate `gemini-extension.json`
- Display extension details

#### Data Model:
```rust
pub struct Extension {
    pub id: String,                    // Internal UUID
    pub name: String,                  // From gemini-extension.json
    pub version: String,
    pub mcp_servers: HashMap<String, McpServer>,
    pub context_file_name: Option<String>,
    pub context_content: Option<String>,
    pub metadata: ExtensionMetadata,   // Our tracking data
}

pub struct ExtensionMetadata {
    pub imported_at: DateTime<Utc>,
    pub source_path: Option<PathBuf>,
    pub tags: Vec<String>,
    pub description: Option<String>,
}
```

### 1.2 Profile Definition
**Priority**: Critical  
**Status**: Not Started

#### Features:
- Create profiles with selected extensions
- Configure environment variables
- Set workspace directory
- Basic profile management (create, edit, delete)

#### Data Model:
```rust
pub struct Profile {
    pub id: String,
    pub name: String,
    pub description: Option<String>,
    pub workspace_dir: PathBuf,
    pub extensions: Vec<ProfileExtension>,
    pub env_vars: HashMap<String, String>,
    pub settings_overrides: Option<Value>, // JSON for settings.json
}

pub struct ProfileExtension {
    pub extension_id: String,
    pub enabled: bool,
    pub server_env_overrides: HashMap<String, HashMap<String, String>>,
}
```

### 1.3 Workspace Setup & Launch
**Priority**: Critical  
**Status**: Not Started

#### Features:
- Create `.gemini` directory structure
- Copy extensions to workspace
- Generate `settings.json` with overrides
- Create `.env` file
- Launch Gemini CLI in workspace

#### Implementation:
```rust
pub fn activate_profile(profile: &Profile) -> Result<()> {
    // 1. Ensure workspace directory exists
    // 2. Create .gemini/extensions/
    // 3. Copy enabled extensions
    // 4. Apply environment overrides
    // 5. Generate settings.json
    // 6. Create .env file
    // 7. Launch gemini
}
```

## Phase 2: Enhanced Usability

### 2.1 Extension Discovery & Search
**Priority**: High  
**Status**: Not Started

#### Features:
- Search extensions by name, tag, or MCP server
- Filter by capabilities
- Quick preview of extension details
- Bulk import from directory scan

### 2.2 Profile Templates
**Priority**: High  
**Status**: Not Started

#### Built-in Templates:
```yaml
web-development:
  extensions: [github-tools, node-tools, browser-tools]
  env_vars:
    NODE_ENV: development
    
data-science:
  extensions: [python-tools, jupyter-tools, database-tools]
  env_vars:
    PYTHON_ENV: analysis
    
devops:
  extensions: [kubernetes-tools, docker-tools, cloud-tools]
  env_vars:
    CLOUD_ENV: staging
```

### 2.3 Quick Profile Switching
**Priority**: High  
**Status**: Not Started

#### Features:
- Hotkey for profile switcher
- Recent profiles list
- Workspace status indicator
- Option to reuse existing workspace

### 2.4 Environment Variable Management
**Priority**: Medium  
**Status**: Not Started

#### Features:
- Secure credential storage (OS keychain)
- Variable templates with placeholders
- Import from existing .env files
- Validation and error checking

## Phase 3: Advanced Features

### 3.1 Workspace Management
**Priority**: Medium  
**Status**: Not Started

#### Features:
- Multiple workspace support
- Workspace cleanup utilities
- Backup and restore workspaces
- Git integration for workspace versioning

### 3.2 Extension Development Tools
**Priority**: Low  
**Status**: Not Started

#### Features:
- Extension creation wizard
- MCP server testing interface
- Context file editor with preview
- Extension validation and linting

### 3.3 Team Collaboration
**Priority**: Low  
**Status**: Not Started

#### Features:
- Export/import profiles
- Shared extension library (file/git based)
- Profile versioning
- Team settings synchronization

### 3.4 Advanced MCP Server Management
**Priority**: Low  
**Status**: Not Started

#### Features:
- MCP server health monitoring
- Log viewing for running servers
- Resource usage tracking
- Server restart/debug capabilities

## Phase 4: Polish & Integration

### 4.1 User Experience
**Priority**: Medium  
**Status**: Not Started

#### Features:
- Onboarding wizard
- Interactive tutorials
- Contextual help system
- Keyboard shortcut customization

### 4.2 Integration Features
**Priority**: Low  
**Status**: Not Started

#### Features:
- Shell integration (launch from terminal)
- IDE plugins (VS Code, IntelliJ)
- CI/CD integration examples
- Docker containerization support

## Implementation Guidelines

### State Storage Structure
```
~/.local/share/gemini-cli-manager/
├── extensions/
│   ├── {extension-id}/
│   │   ├── gemini-extension.json
│   │   ├── context.md
│   │   └── metadata.json
│   └── ...
├── profiles/
│   ├── {profile-id}.json
│   └── ...
├── workspaces/
│   ├── {workspace-id}/
│   │   └── .cleanup-tracker.json
│   └── ...
└── config.json
```

### Key Technical Decisions

1. **Independent State**: All data stored in platform-specific app directories
2. **Non-invasive**: Never modify global Gemini installation
3. **Workspace-centric**: Each profile activation creates/updates a workspace
4. **Secure by default**: Credentials in OS keychain, not plain text
5. **Git-friendly**: All configs in JSON for version control

### Testing Strategy

1. **Unit Tests**: Model validation, serialization
2. **Integration Tests**: Workspace setup, file operations
3. **UI Tests**: Component interactions, navigation
4. **End-to-end Tests**: Full profile activation flow

### Success Metrics

- Profile activation time < 2 seconds
- Zero data loss on crashes
- Support 100+ extensions without performance degradation
- < 50MB total storage for typical usage

## Migration from Previous Design

Since our understanding has changed significantly:

1. **Remove**: Command-line argument handling for Gemini
2. **Remove**: Complex extension dependency system
3. **Add**: Workspace management functionality
4. **Add**: Settings.json generation
5. **Refocus**: From "launcher" to "workspace manager"

## Next Steps

1. Update existing code to remove incorrect assumptions
2. Implement Phase 1.1 (Extension Library) first
3. Create UI mockups for new workflow
4. Write integration tests for workspace setup
5. Document the new user workflow

This revised roadmap aligns with how Gemini CLI actually works while still providing the profile-based workflow that makes managing multiple configurations easy.