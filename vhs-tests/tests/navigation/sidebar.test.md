# Sidebar Navigation Test Specification

## Test ID: NAV-001
## Component: Sidebar Navigation
## Priority: High

### Description
This test verifies that the sidebar navigation works correctly, including:
- Keyboard navigation (j/k keys)
- Selection with Enter and Space
- Visual feedback for current selection
- Focus management between sidebar and content

### Prerequisites
- Clean state directory (`/tmp/vhs-nav-test`)
- No existing profiles or extensions
- Default theme (github-dark)

### Test Steps & Expected Results

#### 1. Application Start
**Screenshot**: `nav-01-start.png`
- ✓ Application launches successfully
- ✓ Sidebar is visible on the left side
- ✓ "Extensions" (first item) is highlighted with selection color
- ✓ Extensions content panel is displayed
- ✓ All sidebar items show correct icons:
  - Extensions: ◆
  - Profiles: ▣
  - Settings: ⚙
  - Help: ?
- ✓ Status bar shows at bottom

#### 2. Navigate Down to Profiles
**Screenshot**: `nav-02-profiles.png`
- ✓ Selection moves from Extensions to Profiles
- ✓ Profiles item is now highlighted
- ✓ Extensions item is no longer highlighted
- ✓ Content area updates to show Profiles view
- ✓ Smooth transition without flicker

#### 3. Navigate Down to Settings
**Screenshot**: `nav-03-settings.png`
- ✓ Selection moves to Settings
- ✓ Settings content displays with theme selector
- ✓ Previous items not highlighted

#### 4. Navigate Down to Help
**Screenshot**: `nav-04-help.png`
- ✓ Selection moves to Help (last item)
- ✓ Help content displays
- ✓ Navigation wraps if trying to go further down (optional)

#### 5. Navigate Up
**Screenshot**: `nav-05-back-to-settings.png`
- ✓ Selection moves back up to Settings
- ✓ Proper highlight restoration

#### 6. Select with Enter
**Screenshot**: `nav-06-settings-view.png`
- ✓ Settings view is fully loaded
- ✓ Theme list is visible and interactive
- ✓ Focus may shift to content area
- ✓ Sidebar selection remains visible

#### 7. Tab to Sidebar
**Screenshot**: `nav-07-sidebar-focus.png`
- ✓ Tab key moves focus back to sidebar
- ✓ Visual indication of sidebar focus (border or highlight change)
- ✓ Settings remains the selected item

#### 8. Navigate to Extensions
**Screenshot**: `nav-08-extensions-view.png`
- ✓ Can navigate back to Extensions
- ✓ Extensions view loads correctly
- ✓ Empty state shown if no extensions

#### 9. Select with Space
**Screenshot**: `nav-09-profiles-via-space.png`
- ✓ Space bar also selects items
- ✓ Profiles view loads
- ✓ Same behavior as Enter key

### Pass Criteria
1. All navigation keys (j/k/Enter/Space) work as expected
2. Visual feedback is clear and immediate
3. No rendering artifacts or glitches
4. Focus management works correctly with Tab
5. Content updates match sidebar selection
6. Status bar remains stable during navigation

### Known Issues
- None

### Notes
- Pay attention to theme consistency
- Check for any lag in navigation response
- Verify borders and spacing remain consistent