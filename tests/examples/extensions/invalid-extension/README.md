# Invalid Extension

This extension contains various validation errors for testing error handling.

## Validation Errors

1. **Invalid ID**: Will be derived from directory name "invalid-extension" but the manifest has issues
2. **Invalid Name**: Contains special characters (!)
3. **Invalid Version**: Not semantic versioning (1.0.0.0)
4. **Empty Description**: Required field is empty
5. **Empty Author Name**: Required field is empty
6. **Invalid MCP Server**: Missing required "command" field
7. **Missing GEMINI.md**: No documentation file

## Test Scenarios

- Tests validation error reporting
- Tests installation failure handling
- Tests error message clarity
- Tests cleanup after failed installation