# Filesystem Operations Guide

## Core Principles

1. **Always validate paths** - Use absolute paths and verify they exist before operations
2. **Check permissions** - Ensure you have read/write access before attempting operations
3. **Handle errors gracefully** - All filesystem operations can fail, always use proper error handling
4. **Never trust user input** - Sanitize and validate all path inputs to prevent directory traversal

## Required Practices

### Path Validation
- Always resolve paths to absolute form before use
- Verify paths stay within allowed directories
- Check if path exists and is the expected type (file vs directory)

### Safe Operations
- Use atomic writes (write to temp file, then rename)
- Create parent directories before writing files
- Set appropriate permissions (usually 644 for files, 755 for directories)
- Always close file handles when done

### Error Handling
```bash
# Good practice
if [ -f "$file" ]; then
    cat "$file" || echo "Error: Failed to read file"
else
    echo "Error: File does not exist"
fi
```

## Common Tasks

### Reading Files
- Check file exists and is readable
- Handle encoding errors gracefully
- Set reasonable size limits

### Writing Files
- Create parent directories first
- Use atomic writes for safety
- Set correct permissions after creation

### Directory Operations
- Use `mkdir -p` for creating nested directories
- Check directory is empty before removal
- Never use `rm -rf` without explicit user confirmation

## What NOT to Do
- Never use relative paths in production code
- Never follow symbolic links without validation
- Never execute files without verifying source
- Never assume operations will succeed