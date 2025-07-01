# Starter Extensions for Gemini CLI Manager

This directory contains ready-to-use extensions for the Gemini CLI Manager. These extensions provide common functionality and serve as examples for creating your own extensions.

## Available Extensions

### 1. Rust Tools (`rust-tools`)
Comprehensive Rust development support including:
- Cargo command integration
- Rust analyzer for code analysis
- Best practices and patterns guide
- Performance optimization tips

**Context File**: `GEMINI.md` - Extensive Rust development guidelines

### 2. Filesystem Basic (`filesystem-basic`)
Safe filesystem operations with access controls:
- File read/write operations
- Directory management
- File search capabilities
- Security best practices

**Context File**: `FILESYSTEM.md` - Filesystem operations guide

### 3. GitHub Tools (`github-tools`)
Full GitHub API integration:
- Repository management
- Issues and pull requests
- GitHub Actions automation
- Gists and releases

**Context File**: `GITHUB.md` - GitHub workflow guidelines

## Installing Extensions

To install an extension in the Gemini CLI Manager:

1. Launch the Gemini CLI Manager:
   ```bash
   cargo run
   ```

2. Navigate to the Extensions tab

3. Use the import function (when implemented) to import from:
   ```
   ./starter-extensions/<extension-name>
   ```

## Extension Structure

Each extension follows this structure:
```
extension-name/
├── gemini-extension.json   # Extension manifest
└── CONTEXT.md             # Context file for Gemini AI
```

### Extension Manifest (`gemini-extension.json`)

Defines the extension metadata and MCP server configurations:

```json
{
  "name": "extension-name",
  "version": "1.0.0",
  "mcpServers": {
    "server-name": {
      "command": "command-to-run",
      "args": ["arguments"],
      "env": {
        "VAR_NAME": "$ENV_VAR"
      }
    }
  }
}
```

### Context Files

Context files (e.g., `GEMINI.md`, `GITHUB.md`) provide:
- Tool documentation
- Best practices
- Usage examples
- Troubleshooting guides

These files are loaded by Gemini AI to understand how to use the extension effectively.

## Creating Your Own Extensions

To create a new extension:

1. Create a new directory in `starter-extensions/`
2. Add a `gemini-extension.json` manifest
3. Create a context file (e.g., `MYEXT.md`) with:
   - Overview of functionality
   - Available tools and their parameters
   - Best practices
   - Examples
   - Troubleshooting

### Example Custom Extension

```bash
mkdir starter-extensions/my-extension
```

`starter-extensions/my-extension/gemini-extension.json`:
```json
{
  "name": "my-extension",
  "version": "1.0.0",
  "mcpServers": {
    "my-server": {
      "command": "node",
      "args": ["my-server.js"],
      "env": {
        "API_KEY": "$MY_API_KEY"
      }
    }
  }
}
```

`starter-extensions/my-extension/MYEXT.md`:
```markdown
# My Extension

This extension provides...

## Available Tools

- `my_tool`: Description
  - Parameters: `param1` (required), `param2` (optional)
  - Returns: Description of return value

## Best Practices

1. Always...
2. Never...
3. Consider...

## Examples

...
```

## Environment Variables

Extensions can reference environment variables using the `$VAR_NAME` syntax:

```json
"env": {
  "API_KEY": "$MY_API_KEY",
  "HOME_DIR": "$HOME"
}
```

Set these in your profile configuration when creating profiles in the Gemini CLI Manager.

## MCP Server Types

Extensions can use different types of MCP servers:

1. **NPX-based servers** (like GitHub and filesystem):
   ```json
   "command": "npx",
   "args": ["-y", "@modelcontextprotocol/server-name"]
   ```

2. **Direct executables** (like rust-analyzer):
   ```json
   "command": "rust-analyzer",
   "args": ["--stdio"]
   ```

3. **Script-based servers**:
   ```json
   "command": "python",
   "args": ["server.py"]
   ```

## Contributing

To contribute a new starter extension:

1. Create the extension following the structure above
2. Ensure the context file is comprehensive
3. Test the extension with Gemini CLI Manager
4. Submit a pull request

## License

These starter extensions are part of the Gemini CLI Manager project and follow the same license.