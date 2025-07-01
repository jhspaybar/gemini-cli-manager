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

Run the import script from the project root:
```bash
./import-starter-extensions.sh
```

Or manually copy extension.json files to your data directory.

## Creating Your Own Extensions

Use these starter extensions as templates for creating your own. Each extension should include:
- `extension.json` - Extension configuration
- Context file (e.g., `GEMINI.md`) - Instructions for the AI model
- MCP server configuration - Tools and capabilities

See the Gemini CLI documentation for more details on extension format.