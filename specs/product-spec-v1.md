# Gemini CLI Manager - Product Specification v1.0

## Executive Summary

The Gemini CLI Manager is a Terminal User Interface (TUI) application designed to simplify the management of Gemini CLI extensions, MCP servers, and custom tools. It provides a user-friendly interface for configuring, organizing, and launching Gemini CLI with different profiles for various use cases.

## Core Design Principles

1. **Developer-First**: Optimized for keyboard navigation and efficiency
2. **Safe by Default**: Non-destructive operations with clear confirmations
3. **Progressive Disclosure**: Simple for beginners, powerful for experts
4. **Offline-First**: Full functionality without internet connectivity
5. **Git-Friendly**: All configurations in version-controllable formats

## User Personas

### Primary: Power Developer
- Uses multiple MCP servers daily
- Switches between client projects frequently
- Values efficiency and keyboard shortcuts
- Needs reproducible environments

### Secondary: Team Lead
- Manages shared configurations
- Onboards new team members
- Ensures consistency across team
- Monitors extension usage

### Tertiary: Casual User
- Occasional Gemini CLI user
- Prefers GUI over command line
- Needs clear guidance
- Values stability over features

## Feature Specifications

### 1. Extension Management

#### Extension Sources
- **Local Folders**: Primary source at `~/.gemini/extensions/`
- **Git Repositories**: Clone directly from GitHub/GitLab URLs
- **Extension Packs**: Curated collections of related extensions

#### Extension Operations
- **Install**: Copy/clone to extensions directory with validation
- **Enable/Disable**: Toggle without removing files
- **Update**: Git pull for repository-based extensions
- **Remove**: Move to trash (recoverable) rather than delete
- **Validate**: Check settings.json and GEMINI.md integrity

#### Extension Discovery
- **List View**: Show all available extensions with status indicators
- **Search**: Filter by name, description, or MCP server type
- **Categories**: Auto-categorize by functionality (AI, Development, Productivity)
- **Health Check**: Validate configurations and dependencies

### 2. Profile Management

#### Profile Structure
```yaml
name: "Web Development"
description: "Full-stack web development environment"
extensions:
  - id: "typescript-tools"
    enabled: true
  - id: "react-helper"
    enabled: true
  - id: "database-tools"
    enabled: false
environment:
  NODE_ENV: "development"
  API_ENDPOINT: "http://localhost:3000"
mcp_servers:
  - name: "code-analyzer"
    settings:
      model: "gemini-2.0-flash"
created: "2024-01-15T10:00:00Z"
last_modified: "2024-01-15T10:00:00Z"
```

#### Profile Features
- **Quick Switch**: Hotkey (Ctrl+P) to change profiles
- **Templates**: Pre-configured profiles for common scenarios
- **Inheritance**: Base profiles with overrides
- **Validation**: Ensure all referenced extensions exist

#### Profile Storage
- Location: `~/.gemini/profiles/`
- Format: YAML for readability and Git compatibility
- Backup: Automatic daily backups of active profiles

### 3. MCP Server Configuration

#### Configuration Interface
- **Visual Editor**: Form-based editing with validation
- **JSON Mode**: Direct editing for power users
- **Environment Variables**: Secure credential management
- **Test Connection**: Verify server connectivity

#### Supported Settings
```json
{
  "mcp_servers": {
    "code-analyzer": {
      "command": "node",
      "args": ["path/to/server.js"],
      "env": {
        "API_KEY": "${GEMINI_API_KEY}",
        "MODEL": "gemini-2.0-flash"
      },
      "enabled": true
    }
  }
}
```

### 4. User Interface Design

#### Layout Structure
```
┌─────────────────────────────────────────────────────┐
│ Gemini CLI Manager             [Profile: Development]│
├───────────────┬─────────────────────────────────────┤
│ Extensions (5)│ typescript-tools          [Enabled]  │
│ Profiles   (3)│ ──────────────────────────────────── │
│ Settings      │ MCP Servers:                         │
│ Help          │ - language-server (active)           │
│               │ - code-analyzer (active)             │
│               │                                      │
│ [Tab] Switch  │ Description:                         │
│ [Space] Toggle│ Provides TypeScript language support │
│ [Enter] Edit  │ with IntelliSense and refactoring   │
│               │                                      │
│               │ [Space] Toggle  [E] Edit  [D] Delete │
└───────────────┴─────────────────────────────────────┘
```

#### Navigation Model
- **List/Detail Split**: Navigate items on left, view details on right
- **Modal Dialogs**: For create/edit operations
- **Command Palette**: Ctrl+K for quick actions
- **Breadcrumbs**: Show current location in navigation

#### Key Bindings
```
Global:
  Ctrl+Q    - Quit application
  Ctrl+P    - Profile switcher
  Ctrl+K    - Command palette
  ?         - Show help
  
Navigation:
  ↑/↓, j/k  - Move selection
  ←/→, h/l  - Switch panels
  Tab       - Next section
  Enter     - Select/Edit
  
Actions:
  Space     - Toggle enable/disable
  n         - New item
  d         - Delete (with confirmation)
  r         - Refresh
  /         - Search/Filter
```

### 5. Launcher Integration

#### Launch Script Features
```bash
#!/bin/bash
# gemini-launcher

# Auto-detect profile based on current directory
# Launch with specified profile
# Handle environment setup
# Start Gemini CLI with proper configuration
```

#### Launch Modes
- **Interactive**: Select profile from TUI
- **Direct**: `gemini-launcher --profile web-dev`
- **Auto**: Detect from .gemini-profile file
- **Last Used**: Remember previous selection

### 6. Security & Safety

#### Extension Security
- **Checksum Verification**: Validate extension integrity
- **Permission Warnings**: Alert for file system access
- **Sandboxing**: Run in isolated environments where possible
- **Update Notifications**: Alert for security updates

#### Credential Management
- **OS Keychain**: Use system credential stores
- **Environment Files**: .env files excluded from Git
- **Masked Input**: Hide sensitive data in UI
- **Session Tokens**: Temporary credentials when possible

### 7. Data Management

#### Configuration Locations
```
~/.gemini/
├── extensions/          # Extension files
├── profiles/           # Profile definitions
├── config/            # Global settings
├── cache/             # Temporary data
├── logs/              # Application logs
└── backups/           # Automatic backups
```

#### Backup Strategy
- **Automatic**: Daily backups of profiles and settings
- **Before Changes**: Snapshot before major operations
- **Export/Import**: Manual backup to custom locations
- **Version History**: Keep last 7 days of changes

### 8. Error Handling & Recovery

#### Error Categories
- **Validation Errors**: Clear messages with fix suggestions
- **Runtime Errors**: Graceful degradation, continue working
- **Connection Errors**: Offline mode with cached data
- **Configuration Errors**: Rollback to last working state

#### Recovery Features
- **Undo**: Last operation reversal
- **Reset**: Return to default configuration
- **Repair**: Auto-fix common issues
- **Logs**: Detailed logs for debugging

## Implementation Phases

### Phase 1: Core Functionality (MVP)
- Basic extension listing and toggling
- Simple profile creation and switching
- Launch script integration
- File-based configuration

### Phase 2: Enhanced Management
- Git integration for extensions
- Advanced profile features
- MCP server configuration UI
- Search and filtering

### Phase 3: Team Features
- Profile sharing mechanisms
- Extension marketplace integration
- Analytics and monitoring
- Advanced security features

## Success Metrics

1. **Time to Switch Profiles**: < 3 seconds
2. **Extension Toggle Time**: < 1 second
3. **Launch Time**: < 2 seconds
4. **User Error Rate**: < 5% of operations
5. **Configuration Corruption**: 0 incidents

## Technical Requirements

- **Go Version**: 1.21+
- **Terminal**: 80x24 minimum, 256 color support
- **OS Support**: macOS, Linux, Windows (WSL)
- **Dependencies**: Minimal, vendor all dependencies

## Inspiration & References

- **LazyGit**: Navigation patterns and keyboard shortcuts
- **K9s**: Resource management and real-time updates
- **Homebrew**: Simple, clear command structure
- **VS Code**: Extension management patterns

## Appendix: Example User Flows

### First-Time Setup
1. Launch gemini-cli-manager
2. Guided setup wizard
3. Create first profile
4. Add essential extensions
5. Test launch Gemini CLI

### Daily Workflow
1. Open terminal
2. Run `gcm` (alias)
3. See current profile
4. Switch if needed (Ctrl+P)
5. Launch Gemini CLI

### Team Onboarding
1. Clone team repository
2. Import shared profiles
3. Install required extensions
4. Verify configuration
5. Start development

---

*This specification is a living document and will be updated based on user feedback and implementation discoveries.*