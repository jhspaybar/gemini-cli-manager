# Gemini CLI Manager - Test Examples

This directory contains test examples and fixtures for testing the Gemini CLI Manager functionality.

## Directory Structure

```
tests/examples/
├── extensions/     # Example extensions for testing installation
├── profiles/       # Example profile configurations
└── scripts/        # Test and utility scripts
```

## Extensions

### Local Extension Examples

1. **simple-extension/** - A minimal valid extension
2. **mcp-extension/** - Extension with MCP server configuration
3. **invalid-extension/** - Extension with validation errors (for testing)
4. **complex-extension/** - Extension with all features

### Installation Test Commands

```bash
# Install from local path
/path/to/simple-extension

# Install from GitHub (when testing with real repos)
https://github.com/example/gemini-test-extension

# Install from archive
file:///path/to/extension.zip
```

## Profiles

Example profile configurations for testing profile management.

## Scripts

Utility scripts for:
- Setting up test environment
- Creating test extensions
- Cleaning up after tests
- Simulating the Gemini CLI for testing

## Usage

1. Run the setup script to prepare test environment:
   ```bash
   ./scripts/setup-test-env.sh
   ```

2. Use the examples when testing features:
   - Copy paths from this directory when testing extension installation
   - Import profile examples when testing profile management
   - Use mock scripts when testing launch functionality

## Contributing

When adding new test cases:
1. Create self-contained examples that don't require external dependencies
2. Document what each example tests
3. Include both positive and negative test cases
4. Keep examples minimal but realistic