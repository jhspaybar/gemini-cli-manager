package extension

import (
	"time"
)

// Loader handles extension loading and unloading
type Loader struct {
	// TODO: Add process manager for MCP servers
	// TODO: Add dependency resolver
}

// NewLoader creates a new extension loader
func NewLoader() *Loader {
	return &Loader{}
}

// Load activates an extension
func (l *Loader) Load(ext *Extension) error {
	ext.Status = StatusLoading

	// TODO: Resolve dependencies
	if err := l.resolveDependencies(ext); err != nil {
		return &LoadError{
			ExtensionID: ext.ID,
			Phase:       "dependency resolution",
			Err:         err,
		}
	}

	// TODO: Start MCP servers
	if err := l.startMCPServers(ext); err != nil {
		return &LoadError{
			ExtensionID: ext.ID,
			Phase:       "MCP server startup",
			Err:         err,
		}
	}

	// TODO: Register tools
	if err := l.registerTools(ext); err != nil {
		return &LoadError{
			ExtensionID: ext.ID,
			Phase:       "tool registration",
			Err:         err,
		}
	}

	ext.LoadedAt = time.Now()
	ext.Status = StatusActive
	return nil
}

// Unload deactivates an extension
func (l *Loader) Unload(ext *Extension) error {
	// TODO: Stop MCP servers
	if err := l.stopMCPServers(ext); err != nil {
		return &LoadError{
			ExtensionID: ext.ID,
			Phase:       "MCP server shutdown",
			Err:         err,
		}
	}

	// TODO: Unregister tools
	if err := l.unregisterTools(ext); err != nil {
		return &LoadError{
			ExtensionID: ext.ID,
			Phase:       "tool unregistration",
			Err:         err,
		}
	}

	ext.Status = StatusDisabled
	return nil
}

// resolveDependencies checks and resolves extension dependencies
func (l *Loader) resolveDependencies(ext *Extension) error {
	// Dependencies are handled by Gemini CLI
	// Nothing to do here
	return nil
}

// startMCPServers configures MCP servers for the extension
func (l *Loader) startMCPServers(ext *Extension) error {
	// Note: Gemini CLI handles MCP server lifecycle
	// We just need to ensure the configuration is present
	if len(ext.MCPServers) > 0 {
		// Configuration is present in gemini-extension.json
		// Gemini CLI will handle starting/stopping servers
		return nil
	}

	return nil
}

// stopMCPServers is called when disabling an extension
func (l *Loader) stopMCPServers(ext *Extension) error {
	// Note: Gemini CLI handles MCP server lifecycle
	// Nothing to do here as Gemini CLI will stop servers
	return nil
}

// registerTools is called when enabling an extension
func (l *Loader) registerTools(ext *Extension) error {
	// Tools are handled by Gemini CLI through extension config
	return nil
}

// unregisterTools is called when disabling an extension
func (l *Loader) unregisterTools(ext *Extension) error {
	// Tools are handled by Gemini CLI
	return nil
}
