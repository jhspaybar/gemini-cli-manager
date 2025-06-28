# Multiple State Directories Test Specification

## Test ID: STATE-001
## Component: State Directory Management
## Priority: High

### Description
This test verifies that the `--state-dir` flag allows running multiple independent instances of the application with separate data storage.

### Prerequisites
- No existing directories at test paths
- Clean environment

### Test Steps & Expected Results

#### 1. Show Help Documentation
**Screenshot**: `state-01-help.png`
- ✓ Help text displays `--state-dir` flag
- ✓ Examples show how to use custom directories
- ✓ Default behavior is documented

#### 2. Launch Work Setup
**Screenshot**: `state-02-work-start.png`
- ✓ Application starts with custom state directory
- ✓ No errors about missing directories
- ✓ Creates directory if it doesn't exist
- ✓ Shows empty/default state

#### 3. Work Settings View
**Screenshot**: `state-03-work-settings.png`
- ✓ Settings shows correct config directory: `~/vhs-work-setup`
- ✓ Path is expanded (no ~)
- ✓ Absolute path is displayed

#### 4. Create Work Profile - New Form
**Screenshot**: `state-04-work-new-profile.png`
- ✓ Profile creation form opens
- ✓ Modal displays correctly

#### 5. Work Profile - Filled
**Screenshot**: `state-05-work-profile-filled.png`
- ✓ Form accepts input
- ✓ Name: "Work Profile"
- ✓ Description: "Development setup for work projects"

#### 6. Work Profile Created
**Screenshot**: `state-06-work-profile-created.png`
- ✓ Profile saved successfully
- ✓ Profile appears in list
- ✓ Returns to profile view

#### 7. Launch Personal Setup
**Screenshot**: `state-07-personal-start.png`
- ✓ Fresh instance starts
- ✓ No data from work setup
- ✓ Independent state

#### 8. Personal Empty Profiles
**Screenshot**: `state-08-personal-empty-profiles.png`
- ✓ No profiles exist (confirming isolation)
- ✓ "Work Profile" is NOT present
- ✓ Shows empty state message

#### 9. Personal Settings View  
**Screenshot**: `state-09-personal-settings.png`
- ✓ Shows different config directory: `~/vhs-personal-setup`
- ✓ Confirms using different state

#### 10. Work Directory Contents
**Screenshot**: `state-10-work-dir.png`
- ✓ Directory exists at `~/vhs-work-setup`
- ✓ Contains `profiles/` subdirectory
- ✓ Contains `extensions/` subdirectory
- ✓ Has expected structure

#### 11. Personal Directory Contents
**Screenshot**: `state-11-personal-dir.png`
- ✓ Directory exists at `~/vhs-personal-setup`
- ✓ Has same structure as work
- ✓ Is independent (different inode/timestamps)

### Pass Criteria
1. Each state directory is completely independent
2. No data leakage between instances
3. Directories are created automatically
4. Settings correctly show active state directory
5. All features work normally with custom directories
6. Tilde expansion works correctly

### Edge Cases to Consider
- Directory permissions
- Non-existent parent directories
- Symbolic links
- Network paths (if supported)

### Notes
- Test demonstrates primary use case for the feature
- Cleanup is performed at end of test
- Could extend to test with --debug flag combination