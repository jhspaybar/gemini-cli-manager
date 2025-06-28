# VHS Visual Testing Suite

This directory contains visual tests for the Gemini CLI Manager using [VHS](https://github.com/charmbracelet/vhs).

## Overview

Each test consists of:
1. A `.tape` file that defines the interactions to record
2. A `.test.md` file that describes expected behaviors to verify in the resulting GIF
3. Generated `.gif` files showing the actual behavior

## Structure

```
vhs-tests/
├── README.md                    # This file
├── Makefile                     # Commands to run tests
├── config/                      # Shared VHS configuration
│   └── default.yml              # Default VHS settings
├── helpers/                     # Helper scripts
│   └── setup-test-env.sh        # Set up test environment
├── tests/                       # Test cases
│   ├── navigation/              # Navigation tests
│   │   ├── sidebar.tape         # Test sidebar navigation
│   │   └── sidebar.test.md      # Expected behaviors
│   ├── extensions/              # Extension management tests
│   │   ├── install.tape
│   │   └── install.test.md
│   ├── profiles/                # Profile management tests
│   │   ├── create.tape
│   │   └── create.test.md
│   └── state-dir/               # State directory tests
│       ├── multiple-setups.tape
│       └── multiple-setups.test.md
└── output/                      # Generated GIFs (gitignored)
```

## Running Tests

```bash
# Run all tests
make test-all

# Run specific test category
make test-navigation
make test-extensions
make test-profiles
make test-state-dir
make test-e2e

# Run single test
vhs tests/navigation/sidebar.tape

# Clean output
make clean
```

## Writing Tests

### Tape File Format

Example `sidebar.tape`:
```
# Configure VHS
Output output/sidebar.gif
Set FontSize 16
Set Width 1200
Set Height 800
Set Theme "Dracula"

# Start the application
Type "./gemini-cli-manager --state-dir /tmp/vhs-test"
Sleep 2s
Screenshot output/sidebar-start.png

# Test navigation
Type "j"  # Move down
Sleep 500ms
Screenshot output/sidebar-move-down.png

Type "Enter"  # Select
Sleep 1s
Screenshot output/sidebar-selected.png

Type "q"  # Quit
Sleep 500ms
```

### Test Specification Format

Example `sidebar.test.md`:
```markdown
# Sidebar Navigation Test

## Test ID: NAV-001
## Component: Sidebar Navigation

### Prerequisites
- Fresh state directory
- Default theme

### Test Steps & Expected Results

1. **Application Start**
   - Screenshot: `sidebar-start.png`
   - ✓ Sidebar visible on left
   - ✓ "Extensions" item highlighted (default)
   - ✓ Icons visible for all menu items

2. **Move Down**
   - Screenshot: `sidebar-move-down.png`
   - ✓ Highlight moves to "Profiles"
   - ✓ Previous item no longer highlighted
   - ✓ Content area updates

3. **Select Item**
   - Screenshot: `sidebar-selected.png`
   - ✓ Profiles view loads
   - ✓ Sidebar remains visible
   - ✓ "Profiles" stays highlighted

### Pass Criteria
- All navigation smooth without flicker
- Focus indicators clear and consistent
- No visual artifacts or rendering issues
```

## Best Practices

1. **Consistent Timing**
   - Use appropriate `Sleep` commands for animations
   - Take screenshots at key moments
   - Allow time for async operations

2. **Clean Environment**
   - Always use temporary state directories
   - Clean up after tests
   - Don't depend on existing data

3. **Clear Specifications**
   - Number each verification point
   - Reference specific screenshots
   - Be explicit about expected behavior

4. **Maintainable Tests**
   - Keep tests focused on one feature
   - Use descriptive names
   - Document any special setup

## CI Integration

Tests can be run in CI but GIF viewing requires manual inspection. Consider:
- Uploading GIFs as artifacts
- Using image diff tools for regression detection
- Creating a visual test report page