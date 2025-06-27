# UI Component Refactoring Plan

This document outlines the plan to refactor existing UI elements into reusable components, following the pattern established with the `TabBar` component.

## Overview

After analyzing the codebase, we've identified 12 UI patterns that are candidates for componentization. This refactoring will improve code reuse, consistency, and testability across the application.

## Component Priorities

### 🔴 Priority 1: High Impact, High Reuse

#### 1. Card Component
**Current Locations:**
- `internal/cli/view.go:268-319` (renderExtensionCard)
- `internal/cli/view.go:421-477` (renderProfileCard)

**Features to Include:**
- Border styles (normal, selected, focused)
- Title with optional icon
- Description with truncation
- Metadata section
- Selection state handling
- Configurable padding and width

**Example API:**
```go
card := components.NewCard(width).
    SetTitle("Extension Name", "🧩").
    SetDescription(description).
    AddMetadata("Version", "1.0.0").
    AddMetadata("Type", "prompt").
    SetSelected(true)

output := card.Render()
```

#### 2. Modal Container Component
**Current Locations:**
- `internal/cli/launch_modal_simple.go:186-260`
- `internal/cli/profile_quick_switch.go:156-168`
- All form files use similar patterns

**Features to Include:**
- Centered positioning
- Consistent border and padding
- Optional title bar
- Width/height configuration
- Content area with proper spacing

**Example API:**
```go
modal := components.NewModal(width, height).
    SetTitle("Install Extension").
    SetContent(formContent)

output := modal.Render()
```

#### 3. Form Field Component
**Current Locations:**
- `internal/cli/extension_install_form.go:190-209`
- `internal/cli/profile_form.go:284-303`
- `internal/cli/extension_edit_form.go:319-396`

**Features to Include:**
- Label rendering
- Input field with focus state
- Validation state/messages
- Help text
- Required field indicator
- Multi-line support

**Example API:**
```go
field := components.NewFormField("Name", value).
    SetRequired(true).
    SetFocused(true).
    SetError("Name is required").
    SetHelp("Enter a unique name")

output := field.Render()
```

### 🟡 Priority 2: Medium Impact

#### 4. Status Bar Component
**Current Location:**
- `internal/cli/view.go:934-1002` (renderStatusBarContent)

**Features to Include:**
- Three-section layout (left, center, right)
- Profile info display
- Message/status area
- Key hints
- Responsive width handling

**Example API:**
```go
statusBar := components.NewStatusBar(width).
    SetLeft("Profile: Production").
    SetCenter("Extension installed successfully").
    SetRight("Press ? for help")

output := statusBar.Render()
```

#### 5. Empty State Component
**Current Locations:**
- `internal/cli/view.go:219-240` (extensions)
- `internal/cli/view.go:371-393` (profiles)

**Features to Include:**
- Centered container
- Icon support
- Title and description
- Call-to-action text
- Consistent styling

**Example API:**
```go
empty := components.NewEmptyState().
    SetIcon("📦").
    SetTitle("No extensions found").
    SetDescription("Install extensions to enhance Gemini").
    SetAction("Press 'i' to install")

output := empty.Render()
```

#### 6. Progress Steps Component
**Current Location:**
- `internal/cli/launch_modal_simple.go:203-222`

**Features to Include:**
- Step list with icons
- Current step highlighting
- Completed/pending states
- Optional step descriptions

**Example API:**
```go
steps := components.NewProgressSteps([]Step{
    {Icon: "🔍", Title: "Searching for updates"},
    {Icon: "📥", Title: "Downloading"},
    {Icon: "✅", Title: "Installing"},
}).SetCurrentStep(1)

output := steps.Render()
```

### 🟢 Priority 3: Lower Complexity

#### 7. Key Help Component
**Current Location:**
- `internal/cli/styles.go:239-246` (renderKeyHelp)

**Features to Include:**
- Key/description pairs
- Consistent formatting
- Separator handling
- Theme-aware styling

#### 8. List Item Component
**Features to Include:**
- Cursor/selection indicator
- Icon support
- Consistent padding
- Focus state

#### 9. Section Header Component
**Features to Include:**
- H1/H2/H3 styles
- Optional icon
- Consistent spacing
- Theme integration

#### 10. Info Box Component
**Current Location:**
- `internal/cli/view.go:733-796`

**Features to Include:**
- Title bar
- Key-value pairs
- Border styling
- Width handling

#### 11. Search Bar Component
**Status:** Already exists in `internal/cli/search.go`
**Action:** Move to `internal/ui/components` package

#### 12. Textarea with Preview Component
**Current Location:**
- `internal/cli/extension_edit_form.go:267-285,399-419`

**Features to Include:**
- Text editing area
- Preview toggle
- Markdown rendering
- Split view support

## Implementation Plan

### Phase 1: Foundation (Week 1)
1. Create Card component with tests
2. Create Modal component with tests
3. Create FormField component with tests
4. Update existing code to use new components

### Phase 2: Core Components (Week 2)
1. Create StatusBar component
2. Create EmptyState component
3. Create ProgressSteps component
4. Move SearchBar to components package

### Phase 3: Supporting Components (Week 3)
1. Create remaining simple components
2. Create Textarea with Preview component
3. Full integration testing
4. Documentation updates

## Testing Strategy

For each component:
1. **Unit tests** in `internal/ui/components/*_test.go`
2. **Visual tests** in `test/adhoc/test_*.go`
3. **Integration tests** to ensure components work together

## Success Criteria

- [ ] All identified components extracted and tested
- [ ] No code duplication for UI patterns
- [ ] Consistent theming across all components
- [ ] Visual tests for all components
- [ ] Updated documentation
- [ ] All existing functionality preserved

## Benefits

1. **Consistency**: Uniform UI behavior and appearance
2. **Maintainability**: Fix bugs in one place
3. **Testability**: Components can be tested in isolation
4. **Development Speed**: Faster to build new features
5. **Theme Support**: Centralized theme application

## Next Steps

1. Review and approve this plan
2. Create GitHub issues for each component
3. Begin with Priority 1 components
4. Regular testing and integration
5. Update CLAUDE.md with new patterns

## Notes

- Follow the pattern established in `internal/ui/components/tabs.go`
- Always use theme colors from `internal/theme`
- Include flexbox support where appropriate
- Consider accessibility in all components
- Maintain backward compatibility during migration