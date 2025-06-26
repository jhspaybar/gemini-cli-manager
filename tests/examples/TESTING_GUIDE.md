# Testing Guide for Gemini CLI Manager

## Quick Start

1. **Setup Test Environment**
   ```bash
   ./scripts/setup-test-env.sh
   source ~/.gemini-test/test-env
   ```

2. **Run the Manager with Test Data**
   ```bash
   gemini-test  # Alias created by setup script
   # OR manually:
   GEMINI_CLI_PATH=~/.gemini-test/mock-gemini go run .
   ```

3. **Test Extension Installation**
   
   In the app, press `n` in Extensions view and try these paths:

   **Local Directories:**
   - `./tests/examples/extensions/simple-extension`
   - `./tests/examples/extensions/mcp-extension`
   - `./tests/examples/extensions/invalid-extension` (should fail)

   **Absolute Paths:**
   - `~/.gemini-test/extensions/simple-extension`
   - `/tmp/test-gemini-extension` (if you created it)

4. **Test Profile Management**
   
   Create profiles with these environment variables:
   ```
   NODE_ENV=development
   GEMINI_THEME=dark
   PYTHON_ENV=conda
   ```

5. **Test Launch Functionality**
   
   The mock gemini at `~/.gemini-test/mock-gemini` will:
   - Display the active profile
   - Show all GEMINI_ environment variables
   - Stay running until Ctrl+C

## Test Scenarios

### Extension Installation
- ✅ Install from local directory
- ✅ Install from home directory path (~/)
- ✅ Install from absolute path
- ✅ Fail on invalid extension
- ✅ Fail on non-existent path
- ✅ Handle duplicates gracefully

### Profile Management
- ✅ Create new profile
- ✅ Edit existing profile
- ✅ Delete profile (except default)
- ✅ Switch profiles with Ctrl+P
- ✅ Launch with profile environment

### Search Functionality
- ✅ Filter extensions by name
- ✅ Filter profiles by name
- ✅ Clear search with Esc

## Cleanup

```bash
./scripts/cleanup-test-env.sh
```

## Manual Testing Checklist

- [ ] Extension installation from local path
- [ ] Extension enable/disable toggle
- [ ] Extension deletion (move to trash)
- [ ] Profile creation with environment vars
- [ ] Profile editing and saving
- [ ] Profile quick switch (Ctrl+P)
- [ ] Launch with selected profile
- [ ] Search functionality (/)
- [ ] Navigation between views
- [ ] Error handling and recovery

## Debugging

1. **Check Installation Logs**
   ```bash
   tail -f /tmp/gemini-cli-manager-debug.log
   ```

2. **Verify Test Environment**
   ```bash
   ls -la ~/.gemini-test/
   ls -la ~/.gemini-test/extensions/
   ```

3. **Test Mock Gemini**
   ```bash
   GEMINI_PROFILE="Test" ~/.gemini-test/mock-gemini
   ```