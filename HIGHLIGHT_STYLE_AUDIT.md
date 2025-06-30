# Highlight Style and Background Color Audit

## Summary

After auditing all component files, I found the following issues with highlight styles and background colors:

### Critical Issues Found

1. **Using `theme::text_muted()` as a background color**
   - This is problematic because `text_muted` is meant for text, not backgrounds
   - Found in multiple files causing poor contrast

2. **Inconsistent highlight styles between components**
   - Different components use different highlight approaches
   - Some use `.bg()` while others use different methods

3. **Poor contrast combinations**
   - Using muted text colors as backgrounds reduces readability
   - Especially problematic in dark themes where the contrast is minimal

## Issues by File

### 1. `profile_form.rs`
- **Line 318**: `.bg(theme::text_muted())` - Using text color as background
- **Issue**: Poor contrast when highlighting selected extensions

### 2. `extension_form.rs`
- **Line 617**: `.bg(theme::text_muted())` - Using text color as background  
- **Line 643**: `.bg(theme::text_muted())` - Using text color as background
- **Issue**: Poor contrast for MCP server selection

### 3. `profile_list.rs`
- **Line 261**: `.bg(theme::selection())` - Correct usage ✓
- **Status**: This component uses the proper selection color

### 4. `extension_list.rs`
- **Line 248**: `.bg(theme::selection())` - Correct usage ✓
- **Status**: This component uses the proper selection color

### 5. `tab_bar.rs`
- **Line 106**: `.bg(theme::selection())` - Correct usage ✓
- **Status**: This component uses the proper selection color

### 6. `confirm_dialog.rs`
- Uses inverted colors for buttons (foreground/background swap)
- **Status**: Acceptable for button styling

## Recommended Fixes

### 1. Replace all instances of `.bg(theme::text_muted())` with `.bg(theme::selection())`

The `theme::selection()` color is specifically designed for highlighting and provides proper contrast.

### 2. Ensure consistency across all components

All list-like selections should use:
```rust
.bg(theme::selection())
```

### 3. For special cases needing different highlights

If a component needs a different highlight style, consider:
- Using `theme::surface()` or `theme::overlay()` for subtle highlights
- Using `theme::selection_bar()` for more prominent selections
- Never use text colors (`text_*`) as backgrounds

## Theme Color Reference

From `theme.rs`:
- `selection()`: Uses `surface1` - designed for selection backgrounds
- `selection_bar()`: Uses `sapphire` - for more prominent selections
- `text_muted()`: Uses `subtext0` - designed for muted text, NOT backgrounds
- `surface()`: Uses `surface0` - good for subtle backgrounds
- `overlay()`: Uses `surface1` - good for overlays and selections

## Action Items

1. ✓ Fixed `profile_form.rs` line 318 - Changed from `theme::text_muted()` to `theme::selection()`
2. ✓ Fixed `extension_form.rs` lines 617 and 643 - Changed from `theme::text_muted()` to `theme::selection()`
3. Consider adding a linting rule to prevent using `text_*` colors as backgrounds
4. Update development guidelines to clarify proper color usage

## Fixes Applied

All instances of `.bg(theme::text_muted())` have been replaced with `.bg(theme::selection())` to ensure:
- Proper contrast for highlighted items
- Consistency across all components
- Better readability in all theme variants

The application now uses the correct selection color (`surface1` from Catppuccin) for all highlight backgrounds, which provides appropriate contrast against both text colors and the base background.