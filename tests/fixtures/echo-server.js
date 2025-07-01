#!/usr/bin/env node

/**
 * Simple Echo MCP Server for Testing
 * 
 * This server implements the Model Context Protocol (MCP) and provides
 * basic echo functionality for testing the Gemini CLI Manager.
 */

const readline = require('readline');

// Create interface for reading stdin
const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false
});

// MCP protocol implementation
class EchoMCPServer {
  constructor() {
    this.messageId = 0;
  }

  // Send a JSON-RPC response
  sendResponse(id, result) {
    const response = {
      jsonrpc: "2.0",
      id: id,
      result: result
    };
    console.log(JSON.stringify(response));
  }

  // Send an error response
  sendError(id, code, message) {
    const response = {
      jsonrpc: "2.0",
      id: id,
      error: {
        code: code,
        message: message
      }
    };
    console.log(JSON.stringify(response));
  }

  // Handle incoming requests
  handleRequest(request) {
    try {
      const { id, method, params } = request;

      switch (method) {
        case 'initialize':
          this.sendResponse(id, {
            protocolVersion: "1.0",
            serverInfo: {
              name: "echo-test-server",
              version: "1.0.0"
            },
            capabilities: {
              tools: {
                available: [
                  {
                    name: "echo",
                    description: "Echoes back the input message",
                    inputSchema: {
                      type: "object",
                      properties: {
                        message: {
                          type: "string",
                          description: "The message to echo"
                        }
                      },
                      required: ["message"]
                    }
                  },
                  {
                    name: "ping",
                    description: "Returns 'pong'",
                    inputSchema: {
                      type: "object",
                      properties: {}
                    }
                  }
                ]
              }
            }
          });
          break;

        case 'tools/call':
          const { name, arguments: args } = params;
          
          if (name === 'echo') {
            const message = args.message || '';
            const prefix = process.env.ECHO_PREFIX || 'Echo:';
            this.sendResponse(id, {
              output: `${prefix} ${message}`
            });
          } else if (name === 'ping') {
            this.sendResponse(id, {
              output: 'pong'
            });
          } else {
            this.sendError(id, -32601, `Unknown tool: ${name}`);
          }
          break;

        case 'shutdown':
          this.sendResponse(id, {});
          process.exit(0);
          break;

        default:
          this.sendError(id, -32601, `Method not found: ${method}`);
      }
    } catch (error) {
      this.sendError(request.id || null, -32603, `Internal error: ${error.message}`);
    }
  }
}

// Create server instance
const server = new EchoMCPServer();

// Buffer for incomplete JSON
let buffer = '';

// Process each line of input
rl.on('line', (line) => {
  buffer += line;
  
  try {
    // Try to parse as JSON
    const request = JSON.parse(buffer);
    buffer = ''; // Clear buffer on successful parse
    
    // Handle the request
    server.handleRequest(request);
  } catch (e) {
    // If JSON is incomplete, keep buffering
    if (e instanceof SyntaxError && buffer.length < 10000) {
      // Keep buffering
    } else {
      // Clear buffer and log error
      buffer = '';
      console.error(JSON.stringify({
        jsonrpc: "2.0",
        error: {
          code: -32700,
          message: "Parse error"
        }
      }));
    }
  }
});

// Handle process termination
process.on('SIGTERM', () => {
  process.exit(0);
});

process.on('SIGINT', () => {
  process.exit(0);
});

// Log that server is ready (to stderr so it doesn't interfere with JSON-RPC)
process.stderr.write('Echo MCP Server started\n');