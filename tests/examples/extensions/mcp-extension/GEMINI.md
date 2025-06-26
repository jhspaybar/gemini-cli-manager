# MCP Test Extension

Extension for testing MCP (Model Context Protocol) server integration.

## Features

- **Echo Server**: Simple Node.js server that echoes requests
- **Math Server**: Python server that performs calculations

## MCP Servers

### test-echo-server
A Node.js server that echoes back any request it receives. Useful for testing basic MCP communication.

**Environment Variables:**
- `PORT`: Server port (default: 3000)
- `LOG_LEVEL`: Logging verbosity (debug|info|warn|error)

### test-math-server
A Python server that provides mathematical operations.

**Supported Operations:**
- Basic arithmetic (add, subtract, multiply, divide)
- Advanced functions (sqrt, pow, log)

## Configuration

- `mcp-test.enableLogging`: Enable/disable debug logging
- `mcp-test.timeout`: Server response timeout in milliseconds

## Testing Scenarios

1. **Multiple Servers**: Tests handling of multiple MCP servers in one extension
2. **Environment Variables**: Tests proper environment variable passing
3. **Different Runtimes**: Tests Node.js and Python servers
4. **Configuration**: Tests extension configuration handling