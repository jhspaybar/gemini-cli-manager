# Gemini CLI Manager

> **⚠️ EXPERIMENTAL TOY PROJECT ⚠️**
> 
> This is an experimental project created to learn and explore TUI development with Go/Bubbletea and Rust/Ratatui. 
> **This is NOT production-ready software.** Do not use this project or its code for anything other than learning purposes.
> The code quality, architecture, and functionality are all experimental and should not be relied upon.

A Terminal User Interface (TUI) application for managing Gemini CLI extensions and profiles.

![Gemini CLI Manager Demo](demo.gif)

## Overview

Gemini CLI Manager provides a user-friendly interface to:
- Import and manage Gemini extensions
- Create profiles that combine extensions with custom settings
- Launch Gemini CLI with your configured profiles

## Features

### Extension Management
- Import extensions from local directories
- View extension details including MCP servers and context files
- Browse imported extensions in a searchable list

### Profile Management
- Create profiles with custom names and descriptions
- Select extensions to include in each profile
- Configure working directories and tags
- Launch Gemini CLI with a selected profile

### User Interface
- Intuitive keyboard navigation
- Tab-based interface (Extensions, Profiles, Settings)
- Real-time file browser for importing extensions
- Detailed views for extensions and profiles

## Technology Stack

- **Language**: Rust
- **TUI Framework**: [Ratatui](https://ratatui.rs/) v0.29.0
- **Backend**: Crossterm
- **Async Runtime**: Tokio
- **Architecture**: Component-based with message passing

## Installation

### Prerequisites
- Rust 1.70 or later
- Cargo (comes with Rust)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/jhspaybar/gemini-cli-manager.git
cd gemini-cli-manager

# Build and run
cargo build --release
./target/release/gemini-cli-manager

# Or run directly
cargo run --release
```

## Usage

### Navigation
- **Tab/Shift+Tab**: Switch between tabs
- **↑/↓ or j/k**: Navigate lists
- **Enter**: Select item
- **Space**: Toggle selection
- **Esc or b**: Go back
- **q**: Quit application

### Quick Start

1. **Import an Extension**:
   - Press `i` in the Extensions tab
   - Navigate to your extension directory
   - Select the extension file or directory

2. **Create a Profile**:
   - Switch to Profiles tab
   - Press `n` to create new profile
   - Fill in the details and select extensions
   - Save with `Ctrl+S`

3. **Launch Gemini**:
   - Select a profile
   - Press `l` to launch

## Starter Extensions

The `starter-extensions/` directory includes production-ready extensions:
- **rust-tools**: Comprehensive Rust development environment
- **filesystem-basic**: Safe filesystem operations
- **github-tools**: GitHub API integration

## Configuration

The application stores data in platform-specific directories:
- **macOS**: `~/Library/Application Support/gemini-cli-manager/`
- **Linux**: `~/.local/share/gemini-cli-manager/`
- **Windows**: `%APPDATA%\gemini-cli-manager\`

## Development

For development guidelines and architecture details, see [CLAUDE.md](CLAUDE.md).

## License

[License information to be added]

## Acknowledgments

Built with:
- [Ratatui](https://ratatui.rs/) - Terminal UI framework
- [Tokio](https://tokio.rs/) - Async runtime
- [Crossterm](https://github.com/crossterm-rs/crossterm) - Terminal manipulation