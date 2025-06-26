package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// GeminiSettings represents the settings.json structure for Gemini CLI
type GeminiSettings struct {
	// Core settings
	ContextFileName string `json:"contextFileName,omitempty"`
	Theme          string `json:"theme,omitempty"`
	Sandbox        interface{} `json:"sandbox,omitempty"` // bool or string
	AutoAccept     bool   `json:"autoAccept,omitempty"`
	
	// Tool configuration
	CoreTools    []string `json:"coreTools,omitempty"`
	ExcludeTools []string `json:"excludeTools,omitempty"`
	
	// File filtering
	FileFiltering *FileFiltering `json:"fileFiltering,omitempty"`
	
	// MCP servers
	MCPServers map[string]extension.MCPServer `json:"mcpServers,omitempty"`
	
	// Other options
	Checkpointing          *Checkpointing `json:"checkpointing,omitempty"`
	PreferredEditor        string         `json:"preferredEditor,omitempty"`
	Telemetry             *Telemetry     `json:"telemetry,omitempty"`
	UsageStatisticsEnabled *bool          `json:"usageStatisticsEnabled,omitempty"`
}

// FileFiltering represents file filtering options
type FileFiltering struct {
	RespectGitIgnore           bool `json:"respectGitIgnore"`
	EnableRecursiveFileSearch  bool `json:"enableRecursiveFileSearch"`
}

// Checkpointing represents checkpointing configuration
type Checkpointing struct {
	Enabled bool `json:"enabled"`
}

// Telemetry represents telemetry configuration
type Telemetry struct {
	Enabled      bool   `json:"enabled"`
	Target       string `json:"target"`
	OTLPEndpoint string `json:"otlpEndpoint"`
	LogPrompts   bool   `json:"logPrompts"`
}

// SettingsGenerator generates Gemini CLI settings.json
type SettingsGenerator struct {
	profile    *profile.Profile
	extensions []*extension.Extension
}

// NewSettingsGenerator creates a new settings generator
func NewSettingsGenerator(p *profile.Profile, exts []*extension.Extension) *SettingsGenerator {
	return &SettingsGenerator{
		profile:    p,
		extensions: exts,
	}
}

// GenerateSettings creates a GeminiSettings structure from profile and extensions
func (sg *SettingsGenerator) GenerateSettings() (*GeminiSettings, error) {
	settings := &GeminiSettings{
		MCPServers: make(map[string]extension.MCPServer),
	}
	
	// Note: Profile structure doesn't have a Settings field currently
	// We could add it later if needed for profile-specific Gemini settings
	// For now, we'll just use environment variables and MCP servers
	
	// Merge MCP servers from all enabled extensions
	for _, ext := range sg.extensions {
		for name, server := range ext.MCPServers {
			// Check for conflicts
			if _, exists := settings.MCPServers[name]; exists {
				// Handle naming conflicts by prefixing with extension name
				name = fmt.Sprintf("%s_%s", ext.Name, name)
			}
			
			// Pass server config as-is - Gemini CLI handles env var expansion
			settings.MCPServers[name] = server
		}
		
		// Set context file name if specified and not already set
		if ext.ContextFileName != "" && settings.ContextFileName == "" {
			settings.ContextFileName = ext.ContextFileName
		}
	}
	
	return settings, nil
}

// WriteSettings writes the settings to the specified path
func (sg *SettingsGenerator) WriteSettings(path string) error {
	settings, err := sg.GenerateSettings()
	if err != nil {
		return fmt.Errorf("generating settings: %w", err)
	}
	
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating settings directory: %w", err)
	}
	
	// Marshal with pretty printing
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}
	
	// Write atomically
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("writing temporary settings file: %w", err)
	}
	
	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("renaming settings file: %w", err)
	}
	
	return nil
}


// Helper functions
func interfaceSliceToStringSlice(slice []interface{}) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func getBool(m map[string]interface{}, key string, defaultValue bool) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return defaultValue
}

func getString(m map[string]interface{}, key string, defaultValue string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultValue
}