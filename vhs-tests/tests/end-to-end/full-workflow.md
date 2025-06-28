# Full Workflow Test Specification

## Test: full-workflow.tape

### Purpose
Validates the complete end-to-end workflow of installing an extension, creating a profile, activating it, and launching Gemini.

### Prerequisites
- Gemini CLI Manager can be run with `go run .`
- The filesystem-enhanced example extension exists at the expected path
- Gemini CLI is available in the system PATH

### Test Steps

1. **Start Application**
   - Launch with custom state directory `/tmp/vhs-e2e-test`
   - Verify app starts successfully

2. **Install Extension**
   - Press 'n' to open install modal
   - Enter path to filesystem-enhanced example
   - Verify extension is installed successfully
   - Dismiss any confirmation modal

3. **Create Profile**
   - Navigate to Profiles tab
   - Press 'n' to create new profile
   - Enter name: "Development"
   - Enter description: "My development environment with filesystem tools"
   - Save with Ctrl+S

4. **Add Extension to Profile**
   - Press 'e' to edit profile
   - Navigate to extensions section
   - Toggle filesystem-enhanced extension
   - Save changes

5. **Activate Profile**
   - Press Enter to activate the profile
   - Verify profile becomes active

6. **Launch Gemini**
   - Press 'L' to launch Gemini
   - Wait for Gemini to start
   - Verify Gemini TUI appears

7. **Exit Gemini**
   - Type `/quit` command
   - Verify return to CLI Manager

8. **Cleanup**
   - Exit CLI Manager with 'q'
   - Remove test state directory

### Expected Results

- Extension installation completes without errors
- Profile creation and editing works correctly
- Extension can be added to profile
- Profile activation changes the active profile indicator
- Gemini launches with the configured profile
- Gemini accepts `/quit` command and exits cleanly
- Application returns to manager after Gemini exits

### Visual Validation Points

The following screenshots capture key states:
- `e2e-01-start.png` - Initial empty state
- `e2e-04-extensions-list.png` - Extension installed and visible
- `e2e-08-profile-created.png` - Profile successfully created
- `e2e-10-extension-selected.png` - Extension selected in profile
- `e2e-12-profile-activated.png` - Profile shows as active
- `e2e-14-gemini-launched.png` - Gemini TUI is running
- `e2e-15-back-to-manager.png` - Clean return to manager

### Notes

- The test uses filesystem-enhanced extension as it has MCP server configuration
- Test state is isolated in `/tmp/vhs-e2e-test` and cleaned up after
- Timing allows for realistic user interaction speeds
- The test validates the core user journey from setup to usage