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
		versionPattern: regexp.MustCompile(`^\d+\.\d+\.\d+(-\w+)?$`),
	}
}

// Validate performs comprehensive validation on an extension
func (v *Validator) Validate(ext *Extension) error {
	// Validate required fields
	if ext.ID == "" {
		return &ValidationError{Field: "id", Message: "extension ID is required"}
	}
	if !v.idPattern.MatchString(ext.ID) {
		return &ValidationError{Field: "id", Message: "invalid extension ID format (use lowercase letters, numbers, and hyphens)"}
	}

	if ext.Name == "" {
		return &ValidationError{Field: "name", Message: "extension name is required"}
	}

	// Version is required and must be semantic
	if ext.Version != "" && !v.versionPattern.MatchString(ext.Version) {
		return &ValidationError{Field: "version", Message: "invalid version format (use semantic versioning)"}
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
	if ext.MCP != nil {
		for name, server := range ext.MCP.Servers {
			if err := v.validateMCPServer(name, server); err != nil {
				return err
			}
		}
	}

	// Tools validation removed - not part of current structure

	return nil
}

// validateFileStructure checks required files exist
func (v *Validator) validateFileStructure(ext *Extension) error {
	// Check gemini-extension.json exists
	configPath := filepath.Join(ext.Path, "gemini-extension.json")
	if _, err := os.Stat(configPath); err != nil {
		return &ValidationError{Field: "path", Message: "gemini-extension.json not found"}
	}

	// GEMINI.md is optional but recommended
	readmePath := filepath.Join(ext.Path, "GEMINI.md")
	if _, err := os.Stat(readmePath); err != nil {
		// Just a warning, not an error
		fmt.Printf("Warning: GEMINI.md not found for extension %s\n", ext.Name)
	}

	return nil
}

// validateMCPServer validates an MCP server configuration
func (v *Validator) validateMCPServer(name string, server MCPServer) error {
	if server.Command == "" {
		return &ValidationError{
			Field:   fmt.Sprintf("mcp.servers.%s.command", name),
			Message: "command is required",
		}
	}

	return nil
}

// validateTool validates a tool configuration
func (v *Validator) validateTool(name string, tool Tool) error {
	if tool.DisplayName == "" {
		return &ValidationError{
			Field:   fmt.Sprintf("tools.%s.displayName", name),
			Message: "display name is required",
		}
	}

	if tool.Command == "" {
		return &ValidationError{
			Field:   fmt.Sprintf("tools.%s.command", name),
			Message: "command is required",
		}
	}

	// Validate input/output types
	validInput := []string{"", "stdin", "args", "file"}
	validOutput := []string{"", "stdout", "file"}

	if !contains(validInput, tool.Input) {
		return &ValidationError{
			Field:   fmt.Sprintf("tools.%s.input", name),
			Message: "invalid input type",
		}
	}

	if !contains(validOutput, tool.Output) {
		return &ValidationError{
			Field:   fmt.Sprintf("tools.%s.output", name),
			Message: "invalid output type",
		}
	}

	return nil
}

// isBuiltinCommand checks if a command is a built-in system command
func isBuiltinCommand(cmd string) bool {
	builtins := []string{"node", "python", "python3", "ruby", "go", "java", "sh", "bash"}
	for _, builtin := range builtins {
		if cmd == builtin {
			return true
		}
	}
	return false
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