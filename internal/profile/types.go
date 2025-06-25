package profile

import (
	"time"
)

// Profile represents a configuration profile
type Profile struct {
	ID          string                    `yaml:"id"`
	Name        string                    `yaml:"name"`
	Description string                    `yaml:"description"`
	Icon        string                    `yaml:"icon,omitempty"`
	Color       string                    `yaml:"color,omitempty"`
	
	// Configuration
	Extensions  []ExtensionRef            `yaml:"extensions"`
	Environment map[string]string         `yaml:"environment"`
	MCPServers  map[string]ServerConfig   `yaml:"mcp_servers"`
	
	// Metadata
	CreatedAt  time.Time                  `yaml:"created_at"`
	UpdatedAt  time.Time                  `yaml:"updated_at"`
	LastUsed   *time.Time                 `yaml:"last_used,omitempty"`
	UsageCount int                        `yaml:"usage_count"`
	
	// Advanced
	Inherits   []string                   `yaml:"inherits,omitempty"`
	Tags       []string                   `yaml:"tags,omitempty"`
	AutoDetect *AutoDetectRules           `yaml:"auto_detect,omitempty"`
}

// ExtensionRef represents a reference to an extension in a profile
type ExtensionRef struct {
	ID      string                 `yaml:"id"`
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:"config,omitempty"`
}

// ServerConfig represents MCP server configuration in a profile
type ServerConfig struct {
	Enabled  bool                   `yaml:"enabled"`
	Settings map[string]interface{} `yaml:"settings,omitempty"`
}

// AutoDetectRules defines rules for automatic profile detection
type AutoDetectRules struct {
	Patterns []string `yaml:"patterns"`
	Priority int      `yaml:"priority"`
}

// Template represents a profile template
type Template struct {
	ID          string
	Name        string
	Description string
	Extensions  []string
	Environment map[string]string
	MCPServers  []string
	Tags        []string
}

// ValidationError represents a profile validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}