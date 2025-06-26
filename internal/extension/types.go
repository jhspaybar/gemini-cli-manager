package extension

import (
	"time"
)

// Extension represents a Gemini CLI extension
type Extension struct {
	// Core fields that match gemini-extension.json
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Version     string `json:"version"`
	Description string `json:"description"`
	
	// MCP configuration
	MCP *MCPConfig `json:"mcp,omitempty"`
	
	// Additional metadata
	Author      Author            `json:"author,omitempty"`
	Repository  Repository        `json:"repository,omitempty"`
	Categories  []string          `json:"categories,omitempty"`
	Keywords    []string          `json:"keywords,omitempty"`
	
	// Runtime information (not in JSON)
	ID        string    `json:"-"` // Derived from directory name
	Path      string    `json:"-"`
	Enabled   bool      `json:"-"` // Deprecated: profiles control which extensions are active
	LoadedAt  time.Time `json:"-"`
	Status    Status    `json:"-"`
}

// Author represents extension author information
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Repository represents the extension's source repository
type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// MCPConfig represents MCP configuration for an extension
type MCPConfig struct {
	Servers map[string]MCPServer `json:"servers"`
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// Tool represents a custom tool configuration
type Tool struct {
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Input       string   `json:"input,omitempty"`  // stdin, args, file
	Output      string   `json:"output,omitempty"` // stdout, file
}

// ConfigurationSchema defines extension configuration options
type ConfigurationSchema struct {
	Properties map[string]ConfigProperty `json:"properties"`
}

// ConfigProperty represents a single configuration property
type ConfigProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Minimum     *float64    `json:"minimum,omitempty"`
	Maximum     *float64    `json:"maximum,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

// Activation defines when the extension should be activated
type Activation struct {
	Events []string `json:"events"`
}

// Contributions defines what the extension contributes
type Contributions struct {
	Commands     []Command    `json:"commands,omitempty"`
	Keybindings  []Keybinding `json:"keybindings,omitempty"`
}

// Command represents a contributed command
type Command struct {
	Command string `json:"command"`
	Title   string `json:"title"`
}

// Keybinding represents a contributed keybinding
type Keybinding struct {
	Command string `json:"command"`
	Key     string `json:"key"`
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