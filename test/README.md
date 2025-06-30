# Gemini CLI Manager Tests

This directory contains all tests for the Gemini CLI Manager application.

## Quick Start

```bash
# Run all tests
cargo test

# Run with output
cargo test -- --nocapture

# Run specific test module
cargo test tab_bar_test

# Run snapshot tests
cargo insta test
cargo insta review  # Review snapshot changes
```

## Test Organization

- `unit/` - Unit tests for individual components
  - `components/` - UI component tests  
  - `models/` - Data model tests
  - `storage/` - Storage layer tests
- `integration/` - Integration tests for workflows
- `e2e/` - End-to-end user journey tests
- `snapshots/` - Snapshot test files
- `utils/` - Shared test utilities

## Writing Tests

See [CLAUDE.md](./CLAUDE.md) for comprehensive testing guidelines and patterns.

## Key Testing Tools

- **ratatui TestBackend** - For rendering tests without a real terminal
- **insta** - For snapshot testing complex UI layouts
- **rstest** - For parameterized tests
- **tokio-test** - For async operation testing

## Coverage Goals

We aim for:
- >90% coverage for business logic
- >80% coverage for UI components
- 100% coverage for critical paths (CRUD operations)