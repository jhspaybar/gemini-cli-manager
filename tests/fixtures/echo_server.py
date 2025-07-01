#!/usr/bin/env python3
"""
Simple Echo MCP Server for Testing (Python version)

This server implements the Model Context Protocol (MCP) and provides
basic echo functionality for testing the Gemini CLI Manager.
"""

import json
import sys
import os
from typing import Dict, Any

class EchoMCPServer:
    def __init__(self):
        self.message_id = 0
    
    def send_response(self, request_id: Any, result: Dict[str, Any]) -> None:
        """Send a JSON-RPC response"""
        response = {
            "jsonrpc": "2.0",
            "id": request_id,
            "result": result
        }
        print(json.dumps(response), flush=True)
    
    def send_error(self, request_id: Any, code: int, message: str) -> None:
        """Send an error response"""
        response = {
            "jsonrpc": "2.0",
            "id": request_id,
            "error": {
                "code": code,
                "message": message
            }
        }
        print(json.dumps(response), flush=True)
    
    def handle_request(self, request: Dict[str, Any]) -> None:
        """Handle incoming requests"""
        try:
            request_id = request.get("id")
            method = request.get("method")
            params = request.get("params", {})
            
            if method == "initialize":
                self.send_response(request_id, {
                    "protocolVersion": "1.0",
                    "serverInfo": {
                        "name": "echo-test-server-python",
                        "version": "1.0.0"
                    },
                    "capabilities": {
                        "tools": {
                            "available": [
                                {
                                    "name": "echo",
                                    "description": "Echoes back the input message",
                                    "inputSchema": {
                                        "type": "object",
                                        "properties": {
                                            "message": {
                                                "type": "string",
                                                "description": "The message to echo"
                                            }
                                        },
                                        "required": ["message"]
                                    }
                                },
                                {
                                    "name": "ping",
                                    "description": "Returns 'pong'",
                                    "inputSchema": {
                                        "type": "object",
                                        "properties": {}
                                    }
                                },
                                {
                                    "name": "get_env",
                                    "description": "Get an environment variable",
                                    "inputSchema": {
                                        "type": "object",
                                        "properties": {
                                            "name": {
                                                "type": "string",
                                                "description": "Environment variable name"
                                            }
                                        },
                                        "required": ["name"]
                                    }
                                }
                            ]
                        }
                    }
                })
            
            elif method == "tools/call":
                tool_name = params.get("name")
                args = params.get("arguments", {})
                
                if tool_name == "echo":
                    message = args.get("message", "")
                    prefix = os.environ.get("ECHO_PREFIX", "[ECHO]")
                    self.send_response(request_id, {
                        "output": f"{prefix} {message}"
                    })
                
                elif tool_name == "ping":
                    self.send_response(request_id, {
                        "output": "pong"
                    })
                
                elif tool_name == "get_env":
                    var_name = args.get("name", "")
                    value = os.environ.get(var_name, f"<{var_name} not set>")
                    self.send_response(request_id, {
                        "output": f"{var_name}={value}"
                    })
                
                else:
                    self.send_error(request_id, -32601, f"Unknown tool: {tool_name}")
            
            elif method == "shutdown":
                self.send_response(request_id, {})
                sys.exit(0)
            
            else:
                self.send_error(request_id, -32601, f"Method not found: {method}")
        
        except Exception as e:
            self.send_error(
                request.get("id"),
                -32603,
                f"Internal error: {str(e)}"
            )

def main():
    """Main entry point"""
    server = EchoMCPServer()
    buffer = ""
    
    # Log startup to stderr
    sys.stderr.write("Python Echo MCP Server started\n")
    sys.stderr.flush()
    
    try:
        while True:
            line = sys.stdin.readline()
            if not line:
                break
            
            buffer += line
            
            # Try to parse JSON
            try:
                request = json.loads(buffer)
                buffer = ""  # Clear buffer on successful parse
                server.handle_request(request)
            except json.JSONDecodeError:
                # Keep buffering if incomplete
                if len(buffer) > 10000:
                    # Buffer too large, clear and send error
                    buffer = ""
                    print(json.dumps({
                        "jsonrpc": "2.0",
                        "error": {
                            "code": -32700,
                            "message": "Parse error"
                        }
                    }), flush=True)
    
    except KeyboardInterrupt:
        sys.exit(0)
    except Exception as e:
        sys.stderr.write(f"Server error: {str(e)}\n")
        sys.exit(1)

if __name__ == "__main__":
    main()