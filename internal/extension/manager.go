package extension

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Manager handles extension lifecycle operations
type Manager struct {
	basePath      string
	extensionsDir string // Same as basePath, needed for installer
	extensions    map[string]*Extension
	mu            sync.RWMutex
	validator     *Validator
	loader        *Loader
}

// NewManager creates a new extension manager
func NewManager(basePath string) *Manager {
	return &Manager{
		basePath:      basePath,
		extensionsDir: basePath,
		extensions:    make(map[string]*Extension),
		validator:     NewValidator(),
		loader:        NewLoader(),
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
		
		// Skip hidden directories (like .trash)
		if strings.HasPrefix(entry.Name(), ".") {
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

	if ext.Status == StatusActive {
		return nil // Already active
	}

	// Load the extension
	if err := m.loader.Load(ext); err != nil {
		ext.Status = StatusError
		return fmt.Errorf("loading extension: %w", err)
	}

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

	if ext.Status != StatusActive {
		return nil // Not active
	}

	// Unload the extension
	if err := m.loader.Unload(ext); err != nil {
		return fmt.Errorf("unloading extension: %w", err)
	}

	ext.Status = StatusDisabled
	return nil
}

// Install adds a new extension from a source
func (m *Manager) Install(source string, isPath bool) (*Extension, error) {
	installer := NewInstaller(m.extensionsDir)
	
	// Install the extension
	ext, err := installer.Install(source, isPath)
	if err != nil {
		return nil, err
	}
	
	// Rescan to pick up the new extension
	if err := m.Scan(); err != nil {
		// Try to clean up
		os.RemoveAll(ext.Path)
		return nil, fmt.Errorf("rescanning after install: %w", err)
	}
	
	return ext, nil
}

// InstallWithProgress installs an extension with progress callback
func (m *Manager) InstallWithProgress(source string, isPath bool, callback InstallProgressCallback) (*Extension, error) {
	installer := NewInstaller(m.extensionsDir)
	
	// Install the extension with progress
	ext, err := installer.InstallWithProgress(source, isPath, callback)
	if err != nil {
		return nil, err
	}
	
	// Rescan to pick up the new extension
	if err := m.Scan(); err != nil {
		// Try to clean up
		os.RemoveAll(ext.Path)
		return nil, fmt.Errorf("rescanning after install: %w", err)
	}
	
	return ext, nil
}

// Remove deletes an extension
func (m *Manager) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ext, exists := m.extensions[id]
	if !exists {
		return fmt.Errorf("extension not found: %s", id)
	}

	// Disable if active
	if ext.Status == StatusActive {
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

// GetActiveExtensions returns all active extensions
func (m *Manager) GetActiveExtensions() []*Extension {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var active []*Extension
	for _, ext := range m.extensions {
		if ext.Status == StatusActive {
			active = append(active, ext)
		}
	}
	return active
}