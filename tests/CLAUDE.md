# Testing Guide for Gemini CLI Manager

This guide provides comprehensive information on testing our Ratatui-based TUI application.

## Table of Contents
1. [Testing Philosophy](#testing-philosophy)
2. [Test Structure](#test-structure)
3. [Testing Tools & Setup](#testing-tools--setup)
4. [Types of Tests](#types-of-tests)
5. [Testing Patterns](#testing-patterns)
6. [Test Plan](#test-plan)
7. [Running Tests](#running-tests)
8. [Best Practices](#best-practices)

## Testing Philosophy

Our testing approach follows these principles:
1. **Separation of Concerns**: Test UI rendering, state management, and event handling separately
2. **Component Isolation**: Test individual components before integration
3. **Snapshot Testing**: Use snapshots for complex UI layouts
4. **Real User Flows**: Integration tests should mirror actual user interactions
5. **Fast Feedback**: Unit tests should run quickly for rapid development

## Test Structure

```
test/
├── CLAUDE.md                 # This file
├── unit/                     # Unit tests for individual components
│   ├── components/           # UI component tests
│   │   ├── extension_form_test.rs
│   │   ├── profile_form_test.rs
│   │   ├── extension_list_test.rs
│   │   ├── profile_list_test.rs
│   │   └── tab_bar_test.rs
│   ├── models/              # Model and data structure tests
│   │   ├── extension_test.rs
│   │   └── profile_test.rs
│   └── storage/             # Storage layer tests
│       └── storage_test.rs
├── integration/             # Integration tests
│   ├── navigation_test.rs   # Test navigation flows
│   ├── crud_operations_test.rs # Test CRUD workflows
│   └── search_test.rs       # Test search functionality
├── e2e/                     # End-to-end tests
│   └── user_journeys_test.rs
├── snapshots/               # Snapshot test files
│   └── __snapshots__/
└── utils/                   # Test utilities
    ├── mod.rs
    └── helpers.rs
```

## Testing Tools & Setup

### Required Dependencies

Add to `Cargo.toml`:
```toml
[dev-dependencies]
# Snapshot testing
insta = { version = "1.39", features = ["yaml"] }
# Async testing
tokio-test = "0.4"
# Property-based testing (optional)
proptest = "1.4"
# Test data generation
fake = "2.9"
rstest = "0.18"
```

### Install Tools
```bash
# Install insta CLI for snapshot management
cargo install cargo-insta

# Install test coverage tool (optional)
cargo install cargo-tarpaulin
```

## Types of Tests

### 1. Unit Tests

Test individual components in isolation:

```rust
// test/unit/components/extension_form_test.rs
use gemini_cli_manager::components::ExtensionForm;
use ratatui::{backend::TestBackend, Terminal};

#[test]
fn test_extension_form_initial_state() {
    let form = ExtensionForm::new(storage);
    assert_eq!(form.current_field(), FormField::Name);
    assert!(form.name_input().is_empty());
}

#[test]
fn test_extension_form_rendering() {
    let backend = TestBackend::new(80, 30);
    let mut terminal = Terminal::new(backend)?;
    let form = ExtensionForm::new(storage);
    
    terminal.draw(|f| {
        form.render(f, f.area());
    })?;
    
    // Verify form title is rendered
    let buffer = terminal.backend().buffer();
    assert!(buffer.content().contains("Create New Extension"));
}
```

### 2. Event Handling Tests

Test keyboard and input handling:

```rust
// test/unit/components/extension_form_test.rs
#[test]
fn test_tab_navigation() {
    let mut form = ExtensionForm::new(storage);
    
    // Initial field should be Name
    assert_eq!(form.current_field(), FormField::Name);
    
    // Tab should move to Version
    form.handle_key_event(KeyCode::Tab.into());
    assert_eq!(form.current_field(), FormField::Version);
    
    // Tab again should move to Description
    form.handle_key_event(KeyCode::Tab.into());
    assert_eq!(form.current_field(), FormField::Description);
}

#[test]
fn test_form_input() {
    let mut form = ExtensionForm::new(storage);
    
    // Type into name field
    for ch in "Test Extension".chars() {
        form.handle_key_event(KeyCode::Char(ch).into());
    }
    
    assert_eq!(form.name_input().value(), "Test Extension");
}
```

### 3. State Management Tests

Test application state transitions:

```rust
// test/unit/state/app_state_test.rs
#[test]
fn test_view_navigation() {
    let mut app = App::default();
    
    // Should start at extension list
    assert_eq!(app.current_view(), ViewType::ExtensionList);
    
    // Tab should switch to profiles
    app.handle_action(Action::SwitchTab);
    assert_eq!(app.current_view(), ViewType::ProfileList);
    
    // Enter should go to detail view
    app.handle_action(Action::Select);
    assert_eq!(app.current_view(), ViewType::ProfileDetail);
}
```

### 4. Snapshot Tests

Use for complex UI layouts:

```rust
// test/integration/snapshots_test.rs
use insta::assert_snapshot;

#[test]
fn test_extension_list_empty_state() {
    let app = App::new_with_empty_storage();
    let output = render_to_string(80, 24, |f| {
        app.render(f, f.area());
    });
    
    assert_snapshot!(output);
}

#[test]
fn test_extension_form_all_fields() {
    let mut form = ExtensionForm::new(storage);
    form.set_name("Test Extension");
    form.set_version("1.0.0");
    form.set_description("A test extension");
    
    let output = render_to_string(80, 30, |f| {
        form.render(f, f.area());
    });
    
    assert_snapshot!(output);
}
```

### 5. Integration Tests

Test complete workflows:

```rust
// test/integration/crud_operations_test.rs
#[test]
fn test_create_extension_workflow() {
    let mut app = App::new_with_test_storage();
    
    // Navigate to create form
    app.handle_key(KeyCode::Char('n').into());
    assert_eq!(app.current_view(), ViewType::ExtensionForm);
    
    // Fill out form
    app.input_text("My Extension");
    app.handle_key(KeyCode::Tab.into());
    app.input_text("1.0.0");
    app.handle_key(KeyCode::Tab.into());
    app.input_text("Description");
    
    // Save
    app.handle_key(KeyCode::Char('s').with_ctrl());
    
    // Should return to list with new extension
    assert_eq!(app.current_view(), ViewType::ExtensionList);
    assert_eq!(app.extension_count(), 1);
}
```

## Testing Patterns

### 1. Test Utilities

Create reusable test helpers:

```rust
// test/utils/helpers.rs
pub mod test_utils {
    use ratatui::{backend::TestBackend, Terminal, Frame};
    
    pub fn setup_test_terminal(width: u16, height: u16) -> Terminal<TestBackend> {
        let backend = TestBackend::new(width, height);
        Terminal::new(backend).unwrap()
    }
    
    pub fn render_to_string<F>(width: u16, height: u16, render_fn: F) -> String
    where
        F: FnOnce(&mut Frame),
    {
        let mut terminal = setup_test_terminal(width, height);
        terminal.draw(render_fn).unwrap();
        terminal.backend().buffer().content()
    }
    
    pub fn create_test_storage() -> Storage {
        Storage::new_in_memory()
    }
    
    pub fn create_test_extension(name: &str) -> Extension {
        Extension {
            id: name.to_lowercase().replace(' ', "-"),
            name: name.to_string(),
            version: "1.0.0".to_string(),
            description: Some(format!("Test extension: {}", name)),
            ..Default::default()
        }
    }
}
```

### 2. Parameterized Tests

Use rstest for parameterized testing:

```rust
use rstest::rstest;

#[rstest]
#[case(80, 24)]
#[case(120, 40)]
#[case(40, 20)]
fn test_responsive_layout(#[case] width: u16, #[case] height: u16) {
    let mut terminal = setup_test_terminal(width, height);
    let app = App::default();
    
    // Should render without panic at any size
    terminal.draw(|f| {
        app.render(f, f.area());
    }).unwrap();
}
```

### 3. Async Testing

For async operations:

```rust
#[tokio::test]
async fn test_extension_installation() {
    let storage = create_test_storage();
    let installer = ExtensionInstaller::new(storage);
    
    let result = installer.install_from_url("https://example.com/ext.json").await;
    assert!(result.is_ok());
    
    let installed = storage.get_extension("example-ext");
    assert!(installed.is_some());
}
```

## Test Plan

### Phase 1: Unit Tests (Priority: High)

1. **Component Tests**
   - [ ] `ExtensionForm`: All field validations, navigation, input handling
   - [ ] `ProfileForm`: Field validations, extension selection, navigation
   - [ ] `ExtensionList`: Filtering, sorting, selection, keyboard navigation
   - [ ] `ProfileList`: Search, selection, default profile handling
   - [ ] `TabBar`: Tab switching, visual state
   - [ ] `ConfirmDialog`: Button selection, escape handling

2. **Model Tests**
   - [ ] Extension: Serialization, validation, MCP server config
   - [ ] Profile: Validation, environment variables, metadata

3. **Storage Tests**
   - [ ] CRUD operations for extensions and profiles
   - [ ] Search and filtering
   - [ ] Error handling

### Phase 2: Integration Tests (Priority: High)

1. **Navigation Flows**
   - [ ] Tab switching between extensions and profiles
   - [ ] List → Detail → Edit flow
   - [ ] Direct edit with 'e' key
   - [ ] Back navigation with escape

2. **CRUD Workflows**
   - [ ] Create new extension with all fields
   - [ ] Edit existing extension
   - [ ] Delete with confirmation
   - [ ] Create profile with extension selection

3. **Search Functionality**
   - [ ] Real-time search filtering
   - [ ] Search mode activation/deactivation
   - [ ] Result count display

### Phase 3: Rendering Tests (Priority: Medium)

1. **Theme Tests**
   - [ ] All components render with correct theme colors
   - [ ] Focus states show proper highlighting
   - [ ] Text contrast meets accessibility standards

2. **Layout Tests**
   - [ ] Responsive to different terminal sizes
   - [ ] Proper text truncation
   - [ ] Border rendering

3. **Snapshot Tests**
   - [ ] Each component in various states
   - [ ] Empty states
   - [ ] Error states
   - [ ] Full forms with data

### Phase 4: Error Handling (Priority: Medium)

1. **Validation Tests**
   - [ ] Required field validation
   - [ ] Version format validation
   - [ ] Duplicate ID handling

2. **Error Display**
   - [ ] Error popup rendering
   - [ ] Error message timeout
   - [ ] User recovery flows

### Phase 5: End-to-End Tests (Priority: Low)

1. **Complete User Journeys**
   - [ ] Install extension → Create profile → Launch
   - [ ] Search → Edit → Save workflow
   - [ ] Batch operations

## Running Tests

```bash
# Run all tests
cargo test

# Run specific test file
cargo test --test extension_form_test

# Run with output
cargo test -- --nocapture

# Run snapshot tests
cargo insta test

# Review snapshot changes
cargo insta review

# Generate test coverage
cargo tarpaulin --out Html

# Run only unit tests
cargo test --test unit::*

# Run only integration tests  
cargo test --test integration::*
```

## Best Practices

### 1. Test Naming
```rust
// Good: Descriptive and specific
#[test]
fn test_extension_form_validates_required_name_field() { }

// Bad: Too generic
#[test]
fn test_form() { }
```

### 2. Test Organization
- Group related tests in modules
- Use descriptive module names
- Keep test files focused and manageable

### 3. Test Data
```rust
// Create builder patterns for test data
struct ExtensionBuilder {
    name: String,
    version: String,
}

impl ExtensionBuilder {
    fn new(name: &str) -> Self {
        Self {
            name: name.to_string(),
            version: "1.0.0".to_string(),
        }
    }
    
    fn with_version(mut self, version: &str) -> Self {
        self.version = version.to_string();
        self
    }
    
    fn build(self) -> Extension {
        Extension {
            name: self.name,
            version: self.version,
            ..Default::default()
        }
    }
}
```

### 4. Assertions
```rust
// Use custom assertions for better error messages
fn assert_renders_text(terminal: &Terminal<TestBackend>, expected: &str) {
    let content = terminal.backend().buffer().content();
    assert!(
        content.contains(expected),
        "Expected to find '{}' in rendered output:\n{}",
        expected,
        content
    );
}
```

### 5. Test Independence
- Each test should be independent
- Use fresh storage for each test
- Clean up any side effects

### 6. Performance
- Keep unit tests fast (< 100ms)
- Use smaller terminal sizes for faster rendering
- Mock expensive operations

## Next Steps

1. Set up test infrastructure with helper utilities
2. Start with high-priority unit tests
3. Add snapshot tests for UI components
4. Implement integration tests for critical paths
5. Set up CI to run tests on every commit
6. Aim for >80% code coverage

Remember: Good tests enable confident refactoring and catch regressions early!