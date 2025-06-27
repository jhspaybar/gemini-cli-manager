package extension

import (
	"time"
)

// Extension represents a Gemini CLI extension
type Extension struct {
	// Core fields that match gemini-extension.json
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`

	// MCP configuration (changed from mcp.servers to mcpServers)
	MCPServers map[string]MCPServer `json:"mcpServers,omitempty"`

	// Context file name override
	ContextFileName string `json:"contextFileName,omitempty"`

	// Runtime information (not in JSON)
	ID       string    `json:"-"` // Derived from directory name
	Path     string    `json:"-"`
	LoadedAt time.Time `json:"-"`
	Status   Status    `json:"-"`
}

// MCPServer represents an MCP server configuration matching Gemini spec
type MCPServer struct {
	// Transport options (one required)
	Command string `json:"command,omitempty"` // Stdio transport
	URL     string `json:"url,omitempty"`     // SSE transport
	HTTPUrl string `json:"httpUrl,omitempty"` // HTTP streaming transport

	// Optional configuration
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	CWD     string            `json:"cwd,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // milliseconds
	Trust   bool              `json:"trust"`             // bypass confirmations
}

// Status represents the current status of an extension
type Status int

const (
	StatusUnknown Status = iota
	StatusInstalled
	StatusLoading
	StatusActive
	StatusError
	StatusDisabled
)

func (s Status) String() string {
	switch s {
	case StatusInstalled:
		return "Installed"
	case StatusLoading:
		return "Loading"
	case StatusActive:
		return "Active"
	case StatusError:
		return "Error"
	case StatusDisabled:
		return "Disabled"
	default:
		return "Unknown"
	}
}

// ValidationError represents an extension validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// LoadError represents an extension loading error
type LoadError struct {
	ExtensionID string
	Phase       string
	Err         error
}

func (e LoadError) Error() string {
	return "failed to load extension " + e.ExtensionID + " during " + e.Phase + ": " + e.Err.Error()
}
