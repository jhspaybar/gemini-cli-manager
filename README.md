# Gemini CLI Manager

A Terminal User Interface (TUI) application for managing Gemini CLI extensions, MCP servers, and profiles.

## Features

- 🔧 **Extension Management** - Install, enable/disable, and configure extensions
- 👤 **Profile Management** - Create and switch between different configuration profiles
- 🚀 **Quick Launch** - Start Gemini CLI with the right profile and extensions
- ⚙️ **Settings** - Configure paths and preferences
- ⌨️ **Keyboard-First** - Fully navigable with keyboard shortcuts

## Installation

```bash
# Clone the repository
git clone https://github.com/gemini-cli/manager.git
cd gemini-cli-manager

# Build the application
make build

# Or directly with go
/usr/local/go/bin/go build -o gemini-cli-manager
```

## Usage

```bash
# Run the TUI
./gemini-cli-manager

# Run with debug logging
./gemini-cli-manager --debug
```

### Keyboard Shortcuts

- `↑/k` - Move up
- `↓/j` - Move down
- `Enter/Space` - Select
- `Esc` - Go back
- `?` - Toggle help
- `q` - Quit

## Development

### Building

```bash
# Quick build check
make build

# Run tests
make test

# Full verification (format, vet, test, build)
make verify

# Clean build artifacts
make clean
```

### Project Structure

```
gemini-cli-manager/
├── main.go              # Application entry point
├── internal/
│   ├── cli/            # TUI components
│   ├── extension/      # Extension management
│   ├── profile/        # Profile management
│   ├── launcher/       # Gemini CLI launcher
│   └── config/         # Configuration
├── specs/              # Detailed specifications
└── scripts/            # Build and utility scripts
```

## Contributing

Please read the [CLAUDE.md](CLAUDE.md) file for development guidelines and best practices.

## License

[License information to be added]