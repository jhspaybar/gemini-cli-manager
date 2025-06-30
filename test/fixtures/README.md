# Test Fixtures

This directory contains test fixtures for the Gemini CLI Manager tests.

## MCP Echo Servers

### echo-server.js
A Node.js implementation of an MCP echo server that provides:
- `echo` tool: Returns the input message with a configurable prefix
- `ping` tool: Always returns "pong"

Environment variables:
- `ECHO_PREFIX`: Custom prefix for echo responses (default: "Echo:")

### echo_server.py
A Python implementation of an MCP echo server that provides:
- `echo` tool: Returns the input message with a configurable prefix
- `ping` tool: Always returns "pong"
- `get_env` tool: Returns the value of an environment variable

Environment variables:
- `ECHO_PREFIX`: Custom prefix for echo responses (default: "[ECHO]")

## Usage in Tests

These servers are used to test:
1. MCP server configuration in extensions
2. Workspace setup and extension installation
3. Environment variable passing
4. Server launch validation

Example test usage:
```rust
let server_config = McpServerConfig {
    command: Some("node".to_string()),
    args: Some(vec!["test/fixtures/echo-server.js".to_string()]),
    env: Some(HashMap::from([
        ("ECHO_PREFIX".to_string(), "Test: ".to_string())
    ])),
    timeout: Some(30000),
    trust: Some(true),
};
```