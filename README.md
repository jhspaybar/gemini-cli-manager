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
# Run the TUI with default state directory (~/.gemini-cli-manager)
./gemini-cli-manager

# Run with custom state directory
./gemini-cli-manager --state-dir /path/to/state

# Run with debug logging
./gemini-cli-manager --debug

# Show help
./gemini-cli-manager --help
```

### State Directory

By default, the application stores all data (extensions, profiles, settings) in `~/.gemini-cli-manager`. You can override this with the `--state-dir` flag to:

- Run multiple independent setups
- Test without affecting your main configuration
- Store data on different volumes

Examples:
```bash
# Separate work and personal setups
./gemini-cli-manager --state-dir ~/gemini-work
./gemini-cli-manager --state-dir ~/gemini-personal

# Testing
./gemini-cli-manager --state-dir /tmp/gemini-test
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