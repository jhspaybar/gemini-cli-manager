# Search Functionality Test Specification

## Test ID: NAV-002
## Component: Search Bar
## Priority: Medium

### Description
This test verifies the search functionality including:
- Search activation with `/` key
- Context-aware placeholder text
- Query input and clearing
- Visual feedback during search
- Search cancellation

### Prerequisites
- Clean state directory
- Default theme

### Test Steps & Expected Results

#### 1. Application Start
**Screenshot**: `search-01-start.png`
- ✓ Extensions view is active
- ✓ Search bar is visible but not active
- ✓ No cursor in search field

#### 2. Activate Search in Extensions
**Screenshot**: `search-02-extensions-search-active.png`
- ✓ `/` key activates search
- ✓ Search bar shows focus state (border color change)
- ✓ Cursor appears in search field
- ✓ Placeholder shows: "Search extensions by name or description..."

#### 3. Type Search Query
**Screenshot**: `search-03-extensions-query.png`
- ✓ Text "test" appears in search field
- ✓ Results filter in real-time (if items exist)
- ✓ Search remains active

#### 4. Clear Search
**Screenshot**: `search-04-extensions-cleared.png`
- ✓ Ctrl+U clears the search field
- ✓ Placeholder text returns
- ✓ Results reset to show all items
- ✓ Search remains active (cursor still visible)

#### 5. Cancel Search
**Screenshot**: `search-05-search-cancelled.png`
- ✓ Escape key deactivates search
- ✓ Search bar loses focus state
- ✓ Cursor disappears
- ✓ Can navigate normally again

#### 6. Navigate to Profiles
**Screenshot**: `search-06-profiles-view.png`
- ✓ Profiles view loads
- ✓ Search bar is visible
- ✓ Previous search is cleared

#### 7. Activate Search in Profiles
**Screenshot**: `search-07-profiles-search-active.png`
- ✓ `/` key activates search in profiles context
- ✓ Search bar shows focus state

#### 8. Profile Search Placeholder
**Screenshot**: `search-08-profiles-placeholder.png`
- ✓ Placeholder is context-aware
- ✓ Shows: "Search profiles by name or tags..."
- ✓ Different from extensions placeholder

#### 9. Profile Search Query
**Screenshot**: `search-09-profiles-query.png`
- ✓ Text "dev" appears in search field
- ✓ Search functionality works in profiles view

### Pass Criteria
1. Search activates/deactivates properly with keyboard shortcuts
2. Context-aware placeholders display correctly
3. Visual feedback clear for active/inactive states
4. Search queries can be typed and cleared
5. Navigation is blocked during active search
6. Each view maintains independent search state

### Edge Cases
- Very long search queries
- Special characters in search
- Search with no results
- Rapid activation/deactivation

### Notes
- Search is non-persistent between views
- Placeholder text is specific to current view
- Search should be case-insensitive