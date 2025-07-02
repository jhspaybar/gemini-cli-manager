# Starter Extensions

This directory contains production-ready extensions to help you get started with Gemini CLI Manager.

## Available Extensions

### rust-tools
Comprehensive Rust development support with detailed GEMINI.md guide covering cargo best practices, error handling, memory management, async programming, testing, and security.

### filesystem-basic
Safe filesystem operations with security-focused patterns. Includes path validation, atomic file operations, permission management, and resource cleanup.

### github-tools
Complete GitHub API integration toolkit with authentication, rate limiting, repository management, issues, pull requests, and GitHub Actions support.

## Usage

To import these extensions:

1. **Using the UI (recommended)**:
   - Launch the Gemini CLI Manager
   - Press 'i' in the Extensions tab
   - Navigate to the starter-extensions directory
   - Select the extension you want to import

2. **Manual copy**:
   Copy the extension directories to your platform's data directory:
   - **macOS**: `~/Library/Application Support/gemini-cli-manager/extensions/`
   - **Linux**: `~/.local/share/gemini-cli-manager/extensions/`
   - **Windows**: `%APPDATA%\gemini-cli-manager\extensions\`

## Creating Your Own Extensions

Use these starter extensions as templates for creating your own. Each extension should include:
- `extension.json` - Extension configuration
- Context file (e.g., `GEMINI.md`) - Instructions for the AI model
- MCP server configuration - Tools and capabilities

See the Gemini CLI documentation for more details on extension format.