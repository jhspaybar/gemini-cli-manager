package extension

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Manager handles extension lifecycle operations
type Manager struct {
	basePath   string
	extensions map[string]*Extension
	mu         sync.RWMutex
	validator  *Validator
	loader     *Loader
}

// NewManager creates a new extension manager
func NewManager(basePath string) *Manager {
	return &Manager{
		basePath:   basePath,
		extensions: make(map[string]*Extension),
		validator:  NewValidator(),
		loader:     NewLoader(),
	}
}

// BasePath returns the base path for extensions
func (m *Manager) BasePath() string {
	return m.basePath
}

// Scan discovers all extensions in the base path
func (m *Manager) Scan() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure base path exists
	if err := os.MkdirAll(m.basePath, 0755); err != nil {
		return fmt.Errorf("creating extensions directory: %w", err)
	}

	// Read directory entries
	entries, err := os.ReadDir(m.basePath)
	if err != nil {
		return fmt.Errorf("reading extensions directory: %w", err)
	}

	// Clear existing extensions
	m.extensions = make(map[string]*Extension)

	// Process each directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		extPath := filepath.Join(m.basePath, entry.Name())
		ext, err := m.loadExtension(extPath)
		if err != nil {
			// Log error but continue scanning
			fmt.Printf("Warning: failed to load extension at %s: %v\n", extPath, err)
			continue
		}

		m.extensions[ext.ID] = ext
	}

	return nil
}

// loadExtension loads a single extension from a directory
func (m *Manager) loadExtension(path string) (*Extension, error) {
	// Check for gemini-extension.json
	configPath := filepath.Join(path, "gemini-extension.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading gemini-extension.json: %w", err)
	}

	// Parse extension
	var ext Extension
	if err := json.Unmarshal(data, &ext); err != nil {
		return nil, fmt.Errorf("parsing gemini-extension.json: %w", err)
	}

	// Set runtime information
	ext.ID = filepath.Base(path) // ID is the directory name
	ext.Path = path
	ext.Status = StatusInstalled
	// Start enabled by default
	ext.Enabled = true

	// Validate extension
	if err := m.validator.Validate(&ext); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &ext, nil
}

// List returns all discovered extensions
func (m *Manager) List() []*Extension {
	m.mu.RLock()
	defer m.mu.RUnlock()

	extensions := make([]*Extension, 0, len(m.extensions))
	for _, ext := range m.extensions {
		extensions = append(extensions, ext)
	}
	return extensions
}

// Get returns an extension by ID
func (m *Manager) Get(id string) (*Extension, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ext, exists := m.extensions[id]
	if !exists {
		return nil, fmt.Errorf("extension not found: %s", id)
	}
	return ext, nil
}

// Enable activates an extension
func (m *Manager) Enable(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, exists := m.extensions[id]
	if !exists {
		return fmt.Errorf("extension not found: %s", id)
	}

	if ext.Enabled {
		return nil // Already enabled
	}

	// Load the extension
	if err := m.loader.Load(ext); err != nil {
		ext.Status = StatusError
		return fmt.Errorf("loading extension: %w", err)
	}

	ext.Enabled = true
	ext.Status = StatusActive
	return nil
}

// Disable deactivates an extension
func (m *Manager) Disable(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, exists := m.extensions[id]
	if !exists {
		return fmt.Errorf("extension not found: %s", id)
	}

	if !ext.Enabled {
		return nil // Already disabled
	}

	// Unload the extension
	if err := m.loader.Unload(ext); err != nil {
		return fmt.Errorf("unloading extension: %w", err)
	}

	ext.Enabled = false
	ext.Status = StatusDisabled
	return nil
}

// Install adds a new extension from a source
func (m *Manager) Install(source string) error {
	// TODO: Implement installation from various sources
	// - Local directory
	// - Git repository
	// - Archive file
	return fmt.Errorf("install not yet implemented")
}

// Remove deletes an extension
func (m *Manager) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, exists := m.extensions[id]
	if !exists {
		return fmt.Errorf("extension not found: %s", id)
	}

	// Disable if enabled
	if ext.Enabled {
		if err := m.loader.Unload(ext); err != nil {
			return fmt.Errorf("unloading extension: %w", err)
		}
	}

	// Move to trash instead of deleting
	trashPath := filepath.Join(m.basePath, ".trash", ext.ID)
	if err := os.MkdirAll(filepath.Dir(trashPath), 0755); err != nil {
		return fmt.Errorf("creating trash directory: %w", err)
	}

	if err := os.Rename(ext.Path, trashPath); err != nil {
		return fmt.Errorf("moving to trash: %w", err)
	}

	delete(m.extensions, id)
	return nil
}

// MoveToTrash is an alias for Remove
func (m *Manager) MoveToTrash(id string) error {
	return m.Remove(id)
}

// Update updates an extension to the latest version
func (m *Manager) Update(id string) error {
	// TODO: Implement update functionality
	return fmt.Errorf("update not yet implemented")
}

// GetEnabledExtensions returns all enabled extensions
func (m *Manager) GetEnabledExtensions() []*Extension {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var enabled []*Extension
	for _, ext := range m.extensions {
		if ext.Enabled {
			enabled = append(enabled, ext)
		}
	}
	return enabled
}