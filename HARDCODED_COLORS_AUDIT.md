# Hardcoded Colors Audit

## Summary

Found hardcoded colors in 6 component files that need to be replaced with theme functions.

## Files with Hardcoded Colors

### 1. `src/components/profile_list.rs`
- **Line 154**: `Color::Yellow` → `theme::highlight()`
- **Line 187**: `Color::Green` → `theme::success()`
- **Line 206**: `Color::Yellow` → `theme::highlight()` 
- **Line 209**: `Color::White` → `theme::text_primary()`
- **Line 213**: `Color::Blue` → `theme::primary()`
- **Line 224**: `Color::Gray` → `theme::text_secondary()`
- **Line 233**: `Color::DarkGray` → `theme::text_muted()`
- **Line 238**: `Color::Magenta` → `theme::secondary()`
- **Line 246-247**: `Color::Cyan` → `theme::info()`
- **Line 261**: `Color::DarkGray` → `theme::selection()`
- **Line 278**: `Color::DarkGray` → `theme::text_disabled()`

### 2. `src/components/profile_detail.rs`
- **Line 130**: `Color::Green` → `theme::success()`
- **Line 140, 147, 163, 174, 181**: `Color::Yellow` → `theme::highlight()`
- **Line 156**: `Color::Blue` → `theme::primary()`
- **Line 166**: `Color::Magenta` → `theme::secondary()`
- **Line 189**: `Color::Cyan` → `theme::info()`
- **Line 199**: `Color::Green` → `theme::success()`
- **Line 201**: `Color::DarkGray` → `theme::text_muted()`
- **Line 207**: `Color::Gray` → `theme::text_secondary()`
- **Line 217**: `Color::DarkGray` → `theme::text_muted()`
- **Line 229**: `Color::Cyan` → `theme::info()`
- **Line 249**: `Color::Yellow` → `theme::highlight()`
- **Line 259**: `Color::Cyan` → `theme::info()`
- **Line 282**: `Color::DarkGray` → `theme::text_disabled()`

### 3. `src/components/extension_form.rs`
- **Line 304**: `Color::Cyan` → `theme::info()`
- **Line 353, 372, 388, 405, 468, 488**: `Color::Yellow` → `theme::highlight()`
- **Line 479**: `Color::DarkGray` → `theme::text_muted()`
- **Line 541**: `Color::Green` → `theme::success()`
- **Line 574**: `Color::Yellow` → `theme::highlight()`
- **Line 617, 643**: `Color::DarkGray` → `theme::selection()`
- **Line 660**: `Color::DarkGray` → `theme::text_disabled()`

### 4. `src/components/extension_detail.rs`
- **Line 108**: `Color::Cyan` → `theme::info()`
- **Line 118, 126, 134, 145, 150**: `Color::Yellow` → `theme::highlight()`
- **Line 137**: `Color::Blue` → `theme::primary()`
- **Line 162**: `Color::Magenta` → `theme::secondary()`
- **Line 169**: `Color::Green` → `theme::success()`
- **Line 176, 183**: `Color::Blue` → `theme::primary()`
- **Line 201**: `Color::Yellow` → `theme::highlight()`
- **Line 214**: `Color::Green` → `theme::success()`
- **Line 214**: `Color::Red` → `theme::error()`
- **Line 233**: `Color::Magenta` → `theme::secondary()`
- **Line 258**: `Color::DarkGray` → `theme::text_disabled()`

### 5. `src/components/profile_form.rs`
- **Line 204**: `Color::Blue` → `theme::primary()`
- **Line 225, 249, 274, 299, 333**: `Color::Yellow` → `theme::highlight()`
- **Line 318**: `Color::DarkGray` → `theme::selection()`
- **Line 320**: `Color::Green` → `theme::success()`
- **Line 362**: `Color::DarkGray` → `theme::text_disabled()`

### 6. `src/components/confirm_dialog.rs`
- **Line 68**: `Color::Black` → `theme::background()`
- **Line 77**: `Color::Red` → `theme::error()`
- **Line 110**: `Color::Black` → `theme::background()`
- **Line 111**: `Color::White` → `theme::text_primary()`
- **Line 114**: `Color::DarkGray` → `theme::text_muted()`
- **Line 130**: `Color::Black` → `theme::background()`
- **Line 131**: `Color::Red` → `theme::error()`
- **Line 134**: `Color::Red` → `theme::error()`

## Theme Color Mapping

The theme provides these semantic color functions:

- **Text colors**:
  - `theme::text_primary()` - Main text color
  - `theme::text_secondary()` - Secondary/subdued text
  - `theme::text_muted()` - Muted/disabled text
  - `theme::text_disabled()` - Disabled UI elements

- **UI colors**:
  - `theme::primary()` - Primary brand color (blue)
  - `theme::secondary()` - Secondary brand color (mauve/purple)
  - `theme::accent()` - Accent color (pink)
  - `theme::highlight()` - Highlight/attention color (yellow)

- **Semantic colors**:
  - `theme::success()` - Success states (green)
  - `theme::error()` - Error states (red)
  - `theme::warning()` - Warning states (peach/orange)
  - `theme::info()` - Informational (sky/cyan)

- **Layout colors**:
  - `theme::background()` - Main background
  - `theme::surface()` - Surface/card backgrounds
  - `theme::overlay()` - Modal/overlay backgrounds
  - `theme::border()` - Default border color
  - `theme::border_focused()` - Focused border color
  - `theme::selection()` - Selected item background
  - `theme::cursor()` - Cursor color

## Action Items

1. Replace all `Color::` references with appropriate `theme::` functions
2. Ensure all components import the theme module
3. Test with different theme flavours to ensure proper contrast
4. Remove any remaining direct color references

## Note

The theme system uses Catppuccin color palettes which provide excellent contrast and accessibility across all theme variants (Latte, Frappe, Macchiato, Mocha).