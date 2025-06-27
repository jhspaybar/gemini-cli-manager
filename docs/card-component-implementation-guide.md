# Card Component Implementation Guide

This guide provides step-by-step instructions for implementing the Card component as our first refactored UI component.

## Step 1: Create the Component File

Create `internal/ui/components/card.go`:

```go
package components

import (
    "fmt"
    "strings"
    
    "github.com/charmbracelet/lipgloss"
    "github.com/jhspaybar/gemini-cli-manager/internal/theme"
    "github.com/muesli/reflow/truncate"
)
```

## Step 2: Define the Card Structure

Key fields needed:
- Width and height constraints
- Content (title, description, metadata)
- State (selected, focused)
- Styles for different states

## Step 3: Implement Core Methods

1. **Constructor**: `NewCard(width int) *Card`
2. **Content setters**: `SetTitle()`, `SetDescription()`, `AddMetadata()`
3. **State setters**: `SetSelected()`, `SetFocused()`
4. **Style customization**: `SetStyles()`
5. **Render methods**: `Render()`, `RenderCompact()`

## Step 4: Extract from Existing Code

Current implementations to refactor:
1. `internal/cli/view.go:renderExtensionCard()` - lines 268-319
2. `internal/cli/view.go:renderProfileCard()` - lines 421-477

Common patterns to preserve:
- Border changes for selection state
- Content truncation for long descriptions
- Metadata formatting with bullet separators
- Consistent padding and spacing

## Step 5: Create Tests

### Unit Tests (`internal/ui/components/card_test.go`):
```go
func TestCard_SetTitle(t *testing.T)
func TestCard_Truncation(t *testing.T)
func TestCard_StateChanges(t *testing.T)
func TestCard_EmptyContent(t *testing.T)
```

### Visual Test (`test/adhoc/test_card.go`):
- Test all states (normal, selected, focused)
- Test with different content lengths
- Test theme integration
- Test custom styles

## Step 6: Integration

Replace existing card rendering:

**Before:**
```go
// In view.go
func (m Model) renderExtensionCard(ext extension.Extension, selected bool) string {
    // 50+ lines of card rendering logic
}
```

**After:**
```go
// In view.go
func (m Model) renderExtensionCard(ext extension.Extension, selected bool) string {
    card := components.NewCard(m.getCardWidth()).
        SetTitle(ext.Name, ext.Icon).
        SetDescription(ext.Description).
        AddMetadata("Version", ext.Version).
        AddMetadata("Type", ext.Type).
        SetSelected(selected)
    
    return card.Render()
}
```

## Step 7: Update Other Views

Update all card usage:
1. Extension list view
2. Extension grid view
3. Profile list view
4. Profile grid view

## Step 8: Documentation

Update documentation:
1. Add to `internal/ui/components/README.md`
2. Update `internal/ui/components/CLAUDE.md`
3. Add usage examples

## Testing Checklist

- [ ] Cards render correctly in different terminal sizes
- [ ] Selection state is visually distinct
- [ ] Focus state is visually distinct
- [ ] Long content is truncated properly
- [ ] Theme colors are applied correctly
- [ ] Custom styles work as expected
- [ ] No visual regressions in existing views

## Common Pitfalls to Avoid

1. **Width calculations**: Remember to account for borders (2) and padding (4)
2. **State precedence**: Focused state should override selected state
3. **Theme consistency**: Always use theme colors, never hardcode
4. **Content overflow**: Always truncate to prevent layout breaks
5. **Empty content**: Handle missing titles/descriptions gracefully

## Performance Considerations

1. **Render caching**: Consider caching rendered output if content doesn't change
2. **Truncation**: Use efficient truncation library (muesli/reflow)
3. **Style reuse**: Create styles once, reuse for multiple cards

## Success Metrics

The Card component is successful when:
1. All existing card renders are replaced
2. Code duplication is eliminated (save ~200+ lines)
3. Visual appearance is identical or improved
4. Tests provide 90%+ coverage
5. Performance is equal or better

## Next Steps

After Card component is complete:
1. Use it as a template for Modal component
2. Apply learnings to FormField component
3. Continue with priority list from refactoring plan