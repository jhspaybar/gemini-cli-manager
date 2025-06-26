# Gemini CLI Analysis and Integration Guide

This document provides a comprehensive analysis of the Gemini CLI codebase, its extension system, and integration points relevant to our Gemini CLI Manager.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Extension System](#extension-system)
3. [MCP Server Integration](#mcp-server-integration)
4. [Configuration System](#configuration-system)
5. [Context Files (GEMINI.md)](#context-files-geminimd)
6. [Key Integration Points](#key-integration-points)
7. [Implementation Recommendations](#implementation-recommendations)

## Architecture Overview

The Gemini CLI follows a modular architecture with clear separation of concerns:

### Core Components

1. **CLI Package** (`packages/cli/`)
   - User interface layer using React + Ink for terminal UI
   - Command processing (slash commands, at commands)
   - Theme management
   - History management
   - Configuration loading

2. **Core Package** (`packages/core/`)
   - Backend logic and API communication
   - Tool registration and execution
   - MCP server management
   - State management
   - Prompt construction

3. **Tools System** (`packages/core/src/tools/`)
   - Built-in tools (file operations, shell, web)
   - MCP client integration
   - Tool registry and discovery

### Request Flow

```
User Input → CLI Package → Core Package → Gemini API
                ↓                ↓
           UI Rendering    Tool Execution
                            (if requested)
```

## Extension System

### Extension Structure

Extensions are directories containing a `gemini-extension.json` file:

```
<workspace>/.gemini/extensions/my-extension/
├── gemini-extension.json
├── GEMINI.md (optional context file)
└── ... (other extension files)
```

### Extension Loading

1. **Discovery Locations** (in order of precedence):
   - `<workspace>/.gemini/extensions/` (project-specific)
   - `<home>/.gemini/extensions/` (user-global)

2. **Conflict Resolution**:
   - Workspace extensions override home directory extensions
   - First extension wins when multiple extensions have same name

### gemini-extension.json Format

```json
{
  "name": "my-extension",
  "version": "1.0.0",
  "mcpServers": {
    "server-name": {
      "command": "node",
      "args": ["server.js"],
      "env": {
        "API_KEY": "$MY_API_KEY"
      },
      "cwd": "./servers",
      "timeout": 30000,
      "trust": false
    }
  },
  "contextFileName": "GEMINI.md"
}
```

**Key Fields**:
- `name`: Must match the extension directory name
- `version`: Extension version (informational)
- `mcpServers`: MCP server configurations (merged with settings.json)
- `contextFileName`: Override default context file name (defaults to GEMINI.md)

## MCP Server Integration

### MCP Server Configuration

MCP servers can be configured in three places:
1. Global `~/.gemini/settings.json`
2. Project `.gemini/settings.json`
3. Extension `gemini-extension.json`

Configuration precedence: settings.json > extension

### Server Configuration Properties

```json
{
  "mcpServers": {
    "serverName": {
      // Required (one of):
      "command": "path/to/executable",     // Stdio transport
      "url": "http://localhost:8080/sse",  // SSE transport
      "httpUrl": "http://localhost:3000",  // HTTP streaming
      
      // Optional:
      "args": ["--arg1", "value1"],        // Command arguments
      "env": {                              // Environment variables
        "API_KEY": "$ENV_VAR"               // Supports $VAR expansion
      },
      "cwd": "./server-directory",          // Working directory
      "timeout": 30000,                     // Request timeout (ms)
      "trust": false                        // Bypass confirmations
    }
  }
}
```

### MCP Discovery Process

1. **Connection Phase**:
   - Iterate through configured servers
   - Establish transport connections (Stdio/SSE/HTTP)
   - Track connection status

2. **Tool Discovery**:
   - List available tools from each server
   - Validate tool schemas
   - Sanitize tool names (63 char limit, alphanumeric + underscore/dot/hyphen)

3. **Conflict Resolution**:
   - First server gets unprefixed tool names
   - Subsequent servers get prefixed names: `serverName__toolName`

4. **Schema Processing**:
   - Remove `$schema` properties
   - Strip `additionalProperties`
   - Handle Vertex AI compatibility issues

### Trust Model

- `trust: false` (default): User must confirm each tool execution
- `trust: true`: Bypass all confirmations for this server
- Dynamic allow-listing: Users can "always allow" specific tools/servers

## Configuration System

### Configuration Layers (in precedence order)

1. Default values
2. User settings (`~/.gemini/settings.json`)
3. Project settings (`.gemini/settings.json`)
4. Environment variables
5. Command-line arguments

### Key Configuration Options

```json
{
  // Core settings
  "contextFileName": "GEMINI.md",           // Or array of names
  "theme": "Default",
  "sandbox": false,                         // Or "docker", custom command
  "autoAccept": false,                      // Auto-accept safe tools
  
  // Tool configuration
  "coreTools": ["ReadFileTool", "GlobTool"], // Limit available tools
  "excludeTools": ["run_shell_command"],      // Exclude specific tools
  
  // File filtering
  "fileFiltering": {
    "respectGitIgnore": true,
    "enableRecursiveFileSearch": true
  },
  
  // MCP servers
  "mcpServers": { /* ... */ },
  
  // Other options
  "checkpointing": {"enabled": false},
  "preferredEditor": "vscode",
  "telemetry": {
    "enabled": false,
    "target": "local",
    "otlpEndpoint": "http://localhost:4317",
    "logPrompts": true
  },
  "usageStatisticsEnabled": true
}
```

### Environment Variable Support

- String values in JSON can reference environment variables
- Syntax: `$VAR_NAME` or `${VAR_NAME}`
- Example: `"apiKey": "$MY_API_TOKEN"`

## Context Files (GEMINI.md)

### Hierarchical Loading

Context files are loaded from multiple locations and concatenated:

1. **Global**: `~/.gemini/GEMINI.md`
2. **Project Root & Ancestors**: Search up to `.git` directory or home
3. **Subdirectories**: Below current working directory

### Context File Purpose

- Provide project-specific instructions
- Define coding conventions
- Document component behavior
- Set AI behavior guidelines

### Example Context File

```markdown
# Project: My TypeScript Library

## General Instructions:
- Follow existing coding style
- All functions need JSDoc comments
- Use functional programming patterns

## Coding Style:
- 2 spaces for indentation
- Interface names prefixed with 'I'
- Use strict equality operators

## Component Notes:
- src/api/client.ts handles all API requests
- Use fetchWithRetry for GET requests
```

## Key Integration Points

### 1. Extension Management

Our CLI Manager should:
- Install extensions to `~/.gemini/extensions/` or project `.gemini/extensions/`
- Validate `gemini-extension.json` format
- Handle version conflicts between workspace/home extensions
- Manage MCP server configurations

### 2. Profile System Integration

Map our profile concepts to Gemini's configuration:
- **Extensions**: Copy enabled extensions to active directories
- **Environment Variables**: Set via MCP server `env` config
- **MCP Servers**: Merge configurations from profiles
- **Settings**: Generate appropriate settings.json

### 3. MCP Server Lifecycle

Key considerations:
- Servers are started when Gemini CLI launches
- Persistent connections maintained during session
- Trust settings critical for user experience
- Tool name conflicts handled via prefixing

### 4. Configuration Generation

Generate proper configuration files:
- `settings.json` with merged MCP servers
- Extension `gemini-extension.json` files
- Context files for profile-specific instructions

## Implementation Recommendations

### 1. Extension Installation

```typescript
interface ExtensionInstallation {
  // Validate extension structure
  validateExtension(path: string): boolean
  
  // Install to appropriate directory
  installExtension(extension: Extension, target: 'user' | 'project'): void
  
  // Handle MCP server merging
  mergeMCPServers(profile: Profile): MCPServerConfig
}
```

### 2. Profile Activation

When activating a profile:
1. Clear existing symlinks in `~/.gemini/extensions/`
2. Create symlinks for profile's enabled extensions
3. Each extension's `gemini-extension.json` contains its MCP servers
4. Gemini CLI automatically discovers and loads these configurations
5. No need to generate settings.json - extensions are self-contained!

### 3. MCP Server Configuration

```typescript
interface MCPServerManager {
  // Convert our extension format to Gemini format
  convertToGeminiFormat(extension: Extension): MCPServerConfig
  
  // Handle environment variable expansion
  expandEnvironmentVars(config: MCPServerConfig): MCPServerConfig
  
  // Validate server configuration
  validateServerConfig(config: MCPServerConfig): ValidationResult
}
```

### 4. Extension Management

Extensions are self-contained:
- Each extension has its own `gemini-extension.json`
- MCP servers are defined in the extension's config
- Gemini CLI automatically loads these when discovering extensions
- No need for our manager to merge or generate configurations

### 5. Security Considerations

- **Trust Settings**: Default to `trust: false` for all MCP servers
- **Environment Variables**: Sanitize and validate before expansion
- **Command Execution**: Validate MCP server commands/paths
- **Extension Validation**: Verify extension integrity

### 6. User Experience

- Show which extensions are active in current profile
- Display MCP server connection status
- Provide clear error messages for conflicts
- Allow easy profile switching

## Technical Implementation Details

### Extension Symlink Strategy

```bash
# Profile activation
~/.gemini/extensions/
├── ext1 -> ~/.gemini-cli-manager/extensions/ext1
├── ext2 -> ~/.gemini-cli-manager/extensions/ext2
└── ext3 -> ~/.gemini-cli-manager/extensions/ext3
```

### Configuration Merging

```typescript
function mergeConfigurations(profile: Profile): GeminiConfig {
  const config: GeminiConfig = {
    mcpServers: {},
    ...profile.settings
  }
  
  // Merge MCP servers from extensions
  for (const ext of profile.extensions) {
    if (ext.mcpServers) {
      Object.assign(config.mcpServers, ext.mcpServers)
    }
  }
  
  // Apply environment variable overrides
  for (const [key, value] of Object.entries(profile.environment)) {
    // Inject into MCP server configs
  }
  
  return config
}
```

### Validation Requirements

1. **Extension Names**: Must match directory name
2. **Tool Names**: Alphanumeric + underscore/dot/hyphen, max 63 chars
3. **MCP Commands**: Must be executable
4. **Environment Variables**: Validate against shell injection

## Conclusion

The Gemini CLI provides a flexible extension system that is self-contained and elegant. Our CLI Manager integrates by:

1. Managing extension installations
2. Creating/removing symlinks in `~/.gemini/extensions/` based on active profile
3. Letting Gemini CLI's extension discovery system handle the rest

Key insights:
- Extensions are self-contained with their own `gemini-extension.json`
- MCP servers are defined within each extension
- Gemini CLI automatically discovers and loads extension configurations
- No need for complex configuration merging or generation

This simple approach leverages Gemini CLI's design perfectly - we just manage which extensions are "visible" through symlinks, and Gemini CLI handles all the complexity of loading and configuring them!