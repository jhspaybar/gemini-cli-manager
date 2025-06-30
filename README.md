# Gemini CLI Manager

A Terminal User Interface (TUI) application for managing Gemini CLI workspace configurations through profiles and extensions.

## Overview

Gemini CLI Manager is a **workspace configuration manager** that enables you to:
- Maintain a library of Gemini extensions
- Create profiles that bundle extensions with environment configurations
- Set up project-specific workspaces with selected extensions
- Launch Gemini CLI in properly configured environments

## How It Works

Since Gemini CLI loads extensions from the current directory's `.gemini/extensions/` folder, this manager:
1. **Stores** extensions in a central library
2. **Defines** profiles as workspace templates
3. **Sets up** workspaces by copying selected extensions
4. **Launches** Gemini CLI with the right configuration

## Features

### Core Features (Planned)
- **Extension Library**: Import and manage Gemini extensions in one place
- **Profile Management**: Create workspace configurations for different projects
- **Workspace Setup**: Automatically configure `.gemini/` directories
- **Environment Management**: Handle environment variables and settings
- **Quick Switching**: Fast context switching between projects

### User Interface
- **Keyboard-First Navigation**: Optimized for efficiency with vim-style keybindings
- **Tab-Based Interface**: Easy switching between Extensions, Profiles, Settings, and Help
- **Real-time Search**: Filter extensions and profiles as you type
- **Visual Feedback**: Clear indication of active profiles and workspace status

## Building the Application

### Prerequisites
- Rust 1.70 or later
- Cargo (comes with Rust)

### Build Instructions

```bash
# Clone the repository
git clone <repository-url>
cd gemini-cli-manager-ratatui-claude

# Build in debug mode (faster compilation)
cargo build

# Build in release mode (optimized)
cargo build --release

# Run the application
cargo run

# Run with debug logging
RUST_LOG=debug cargo run
```

### Development Setup

```bash
# Install development tools
cargo install cargo-watch

# Run with auto-reload on file changes
cargo watch -x run

# Run tests
cargo test

# Run tests with coverage
cargo install cargo-tarpaulin
cargo tarpaulin

# Format code
cargo fmt

# Lint code
cargo clippy
```

## Usage

### Basic Workflow
1. **Import Extensions**: Add existing Gemini extensions to your library
2. **Create Profile**: Define a workspace configuration with selected extensions
3. **Activate Profile**: Set up workspace and launch Gemini CLI

### Navigation
- **Tab/Shift+Tab**: Switch between tabs
- **↑/↓ or j/k**: Navigate lists
- **Enter**: Select item
- **Esc**: Go back
- **q**: Quit application
- **?**: Show help

### Extension Management
1. Navigate to Extensions tab
2. Press 'i' to import extension from directory
3. Press 'n' to create new extension
4. Press Enter to view extension details

### Profile Management
1. Navigate to Profiles tab
2. Press 'n' to create new profile
3. Select extensions to include
4. Configure environment variables
5. Press Space to activate profile (sets up workspace and launches Gemini)

## Configuration

### Config File Location
- **macOS/Linux**: `~/.config/gemini-cli-manager/config.json5`
- **Windows**: `%APPDATA%\gemini-cli-manager\config.json5`

### Example Configuration
```json5
{
  "gemini_cli_path": "/usr/local/bin/gemini",
  "extensions_dir": "~/.gemini/extensions",
  "default_profile": "default",
  "theme": "dark",
  "keybindings": {
    "quit": ["q", "Ctrl+c"],
    "help": ["?", "h"],
    "search": ["/", "Ctrl+f"]
  }
}
```

## Architecture

The application follows a component-based architecture with message passing:

```
┌─────────────────────────────────────────────────────────────┐
│                        Terminal UI                          │
├─────────────────────────────────────────────────────────────┤
│                    Component Layer                          │
│  (Extension Library, Profile Manager, Workspace Setup)      │
├─────────────────────────────────────────────────────────────┤
│                    Action Bus (MPSC)                        │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                            │
│  (Extension Service, Profile Service, Workspace Service)    │
├─────────────────────────────────────────────────────────────┤
│                    Storage Layer                            │
│  (Extension Library, Profile Storage, Workspace Tracking)   │
└─────────────────────────────────────────────────────────────┘
```

### Workspace Management Flow:
```
Profile Selection → Workspace Setup → Extension Copying → Launch Gemini
                           ↓
                  Create .gemini/extensions/
                  Generate settings.json
                  Create .env file
```

## Project Structure

```
src/
├── main.rs              # Application entry point
├── app.rs               # Main application orchestrator
├── action.rs            # Action types for message passing
├── components/          # UI components
│   ├── home.rs         # Home screen
│   ├── fps.rs          # FPS counter
│   └── ...             # Other components (to be implemented)
├── config.rs            # Configuration management
├── tui.rs               # Terminal UI setup
└── errors.rs            # Error handling

specs/                   # Feature specifications
├── product-spec-v1.md   # Main product specification
├── extension-architecture.md
├── profile-management.md
├── ui-ux-design.md
└── implementation-roadmap.md

docs/                    # Development documentation
├── ratatui-guide.md     # Ratatui framework guide
├── widget-cookbook.md   # Reusable UI components
└── gemini-cli-architecture.md
```

## Implementation Status

### ✅ Completed
- Basic TUI framework setup
- Component trait and registration system
- Action-based message passing
- Configuration file loading
- Independent state management

### 📋 Phase 1: Core Features (Next)
- Extension library management (import/create)
- Profile creation and editing
- Workspace setup functionality
- Basic Gemini CLI launch

### 📋 Phase 2: Enhanced Usability
- Extension search and filtering
- Profile templates
- Environment variable management
- Quick profile switching

### 🔮 Phase 3: Advanced Features
- Multi-workspace management
- Extension development tools
- Team collaboration features
- MCP server monitoring

See [docs/implementation-roadmap-revised.md](docs/implementation-roadmap-revised.md) for detailed planning.

## Contributing

1. Check the [implementation roadmap](docs/implementation-roadmap-revised.md) for planned features
2. Read [CLAUDE.md](CLAUDE.md) for development guidelines
3. Review [analysis summary](docs/ANALYSIS_SUMMARY.md) for key concepts
4. Follow the existing code patterns and architecture
5. Write tests for new functionality
6. Submit pull requests with clear descriptions

## License

[License information to be added]

## Acknowledgments

Built with:
- [Ratatui](https://ratatui.rs/) - Rust library for terminal user interfaces
- [Tokio](https://tokio.rs/) - Async runtime for Rust
- [Crossterm](https://github.com/crossterm-rs/crossterm) - Cross-platform terminal manipulation