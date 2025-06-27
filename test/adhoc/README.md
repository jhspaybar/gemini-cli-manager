# Visual Tests for UI Components

This directory contains visual tests for UI components. These tests help developers see how components render and behave without running the full application.

## Running Tests

### Method 1: Using the Shell Script (Recommended)
```bash
# List all available tests
./run_visual_tests.sh list

# Run a specific test
./run_visual_tests.sh tabs

# Run all tests
./run_visual_tests.sh all

# Show help
./run_visual_tests.sh help
```

### Method 2: Using Make
```bash
# List available tests
make list

# Run specific test
make tabs
make emoji-width
make tabs-dynamic

# Run all tests
make all

# Run as Go tests
make test
```

### Method 3: Direct Go Command
```bash
# From test/adhoc directory
go run ../../cmd/visual-tests/main.go --list
go run ../../cmd/visual-tests/main.go tabs
go run ../../cmd/visual-tests/main.go all

# As Go tests
VISUAL_TESTS=true go test -v ./. -run TestVisual
```

## Available Tests

- **tabs** - Basic tab bar rendering with content area
- **tabs-dynamic** - Tab bar at different widths
- **tabs-overflow** - Tab bar with too many tabs and edge cases
- **tabs-switch** - Tab bar with different themes
- **emoji-width** - Tests emoji rendering and width calculations
- **gear-spacing** - Tests specific emoji spacing issues

## Adding New Visual Tests

To add a new visual test:

1. Add your test function to both:
   - `cmd/visual-tests/main.go` (for standalone runner)
   - `test/adhoc/visual_test.go` (for Go test runner)

2. For the standalone runner, register it in the `init()` function in `cmd/visual-tests/main.go`:
```go
func init() {
    TestRegistry = map[string]TestFunc{
        // ... existing tests ...
        "my-component": testMyComponent,
    }
}
```

3. For Go tests, add a test function in `test/adhoc/visual_test.go`:
```go
func TestVisualMyComponent(t *testing.T) {
    if !shouldRunVisualTests() {
        t.Skip("Skipping visual test. Set VISUAL_TESTS=true to run.")
    }
    
    output := captureOutput(func() {
        testMyComponent()
    })
    
    // Add assertions if needed
    t.Log("\n" + output)
}
```

4. Run your test:
```bash
# Standalone
./run_visual_tests.sh my-component
make my-component

# As Go test
VISUAL_TESTS=true go test -v -run TestVisualMyComponent
```

## Best Practices

1. **Always initialize theme** - Start each test with `theme.SetTheme("github-dark")`
2. **Test edge cases** - Empty states, overflow, minimum/maximum sizes
3. **Test with different themes** - Ensure components work with all themes
4. **Clear output** - Use headers and separators to make output readable
5. **Document behavior** - Add comments explaining what each test validates

## Example Test Structure

```go
func testComponentName() {
    fmt.Println("Component Name Test")
    fmt.Println("==================")
    fmt.Println()
    
    // Test case 1: Normal usage
    fmt.Println("1. Normal usage:")
    component := components.NewComponent()
    fmt.Println(component.Render())
    fmt.Println()
    
    // Test case 2: Edge case
    fmt.Println("2. Edge case (empty):")
    emptyComponent := components.NewComponent()
    fmt.Println(emptyComponent.Render())
    fmt.Println()
    
    // Test case 3: Different configurations
    fmt.Println("3. With custom configuration:")
    customComponent := components.NewComponent().
        SetWidth(40).
        SetStyle(customStyle)
    fmt.Println(customComponent.Render())
}
```

## Troubleshooting

### Tests not building
- Ensure you're in the `test/adhoc` directory
- Check that all imports are available: `go mod tidy`
- Verify component exists in `internal/ui/components`

### Visual output looks wrong
- Check terminal width (minimum 80 characters recommended)
- Verify theme is initialized
- Test in different terminals (iTerm2, Terminal.app, etc.)

### Adding test fails
- Ensure function name matches pattern: `testComponentName`
- Register in `init()` function with kebab-case name
- Check for typos in function name registration