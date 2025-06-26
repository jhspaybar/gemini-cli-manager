#!/usr/bin/env node

/**
 * Simple Echo MCP Server for testing
 * Echoes back any request it receives
 */

const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false
});

console.error(`Echo server started on PID ${process.pid}`);
console.error(`Environment: PORT=${process.env.PORT}, LOG_LEVEL=${process.env.LOG_LEVEL}`);

// Send initialization message
console.log(JSON.stringify({
  jsonrpc: "2.0",
  method: "initialized",
  params: {
    serverInfo: {
      name: "test-echo-server",
      version: "1.0.0"
    }
  }
}));

// Handle incoming messages
rl.on('line', (line) => {
  try {
    const message = JSON.parse(line);
    
    if (process.env.LOG_LEVEL === 'debug') {
      console.error('Received:', message);
    }
    
    // Echo the message back with a response wrapper
    const response = {
      jsonrpc: "2.0",
      id: message.id || null,
      result: {
        echo: message,
        timestamp: new Date().toISOString(),
        pid: process.pid
      }
    };
    
    console.log(JSON.stringify(response));
  } catch (error) {
    console.error('Error processing message:', error.message);
  }
});

// Handle shutdown
process.on('SIGTERM', () => {
  console.error('Echo server shutting down...');
  process.exit(0);
});