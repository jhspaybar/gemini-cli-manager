package extension

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Validator validates extension configurations
type Validator struct {
	idPattern      *regexp.Regexp
	versionPattern *regexp.Regexp
}

// NewValidator creates a new extension validator
func NewValidator() *Validator {
	return &Validator{
		idPattern:      regexp.MustCompile(`^[a-z0-9-]+$`),
		versionPattern: regexp.MustCompile(`^\d+\.\d+\.\d+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`),
	}
}

// Validate performs comprehensive validation on an extension
func (v *Validator) Validate(ext *Extension) error {
	if ext == nil {
		return &ValidationError{Field: "extension", Message: "extension cannot be nil"}
	}
	
	// Validate required fields
	if ext.ID == "" {
		return &ValidationError{Field: "id", Message: "extension ID is required"}
	}
	if !v.idPattern.MatchString(ext.ID) {
		return &ValidationError{Field: "id", Message: "invalid extension ID format (use lowercase letters, numbers, and hyphens)"}
	}
	if len(ext.ID) > 64 {
		return &ValidationError{Field: "id", Message: "extension ID must be 64 characters or less"}
	}

	if ext.Name == "" {
		return &ValidationError{Field: "name", Message: "extension name is required"}
	}
	// Validate name matches directory name (Gemini requirement)
	if ext.ID != "" && ext.Name != ext.ID {
		return &ValidationError{
			Field: "name", 
			Message: fmt.Sprintf("extension name '%s' must match directory name '%s'", ext.Name, ext.ID),
		}
	}
	if len(ext.Name) > 100 {
		return &ValidationError{Field: "name", Message: "extension name must be 100 characters or less"}
	}

	// Version is required and must be semantic
	if ext.Version == "" {
		return &ValidationError{Field: "version", Message: "extension version is required"}
	}
	if !v.versionPattern.MatchString(ext.Version) {
		return &ValidationError{Field: "version", Message: "invalid version format (use semantic versioning)"}
	}
	
	// Description is optional in Gemini spec
	// But we still recommend it
	if ext.Description == "" {
		fmt.Printf("Warning: extension %s has no description\n", ext.Name)
	}

	// Validate author
	// Author is optional in gemini-extension.json

	// Validate file structure if path is set
	if ext.Path != "" {
		if err := v.validateFileStructure(ext); err != nil {
			return err
		}
	}

	// Validate MCP servers
	for name, server := range ext.MCPServers {
		if err := v.validateMCPServer(name, server); err != nil {
			return err
		}
	}

	// Tools validation removed - not part of current structure

	return nil
}

// validateFileStructure checks required files exist
func (v *Validator) validateFileStructure(ext *Extension) error {
	// Check gemini-extension.json exists
	configPath := filepath.Join(ext.Path, "gemini-extension.json")
	info, err := os.Stat(configPath)
	if err != nil {
		return &ValidationError{Field: "path", Message: "gemini-extension.json not found"}
	}
	if info.IsDir() {
		return &ValidationError{Field: "path", Message: "gemini-extension.json must be a file, not a directory"}
	}

	// Check for context file (GEMINI.md or custom contextFileName)
	contextFileName := "GEMINI.md"
	if ext.ContextFileName != "" {
		contextFileName = ext.ContextFileName
	}
	contextPath := filepath.Join(ext.Path, contextFileName)
	if _, err := os.Stat(contextPath); err != nil {
		// Just a warning, not an error
		fmt.Printf("Warning: context file %s not found for extension %s\n", contextFileName, ext.Name)
	}

	return nil
}

// validateMCPServer validates an MCP server configuration
func (v *Validator) validateMCPServer(name string, server MCPServer) error {
	// Validate transport (one required)
	if server.Command == "" && server.URL == "" && server.HTTPUrl == "" {
		return &ValidationError{
			Field:   fmt.Sprintf("mcpServers.%s", name),
			Message: "at least one transport (command, url, or httpUrl) is required",
		}
	}
	
	// Validate extension name must match directory name
	// This is enforced by Gemini CLI
	
	return nil
}



// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}