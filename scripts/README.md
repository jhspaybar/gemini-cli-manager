# Scripts

This directory contains utility scripts for the Gemini CLI Manager.

## cleanup-old-state.sh

This script helps clean up old manager state from the `~/.gemini` directory after we migrated to storing our state in `~/.gemini-cli-manager`.

### What it does:
- Removes the old `profiles` directory from `~/.gemini/profiles`
- Removes extension symlinks that point to `~/.gemini-cli-manager`
- Preserves all of Gemini's own files (oauth_creds.json, settings.json, etc.)

### Usage:
```bash
./scripts/cleanup-old-state.sh
```

The script will show you what it finds and ask for confirmation before removing anything.

### When to run:
- After updating to the new version that uses `~/.gemini-cli-manager`
- If you see duplicate state in both directories
- To clean up after testing/development

The script is safe to run multiple times - it will only remove manager-specific files and will always preserve Gemini's own configuration.