package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jhspaybar/gemini-cli-manager/internal/state"
	"gopkg.in/yaml.v3"
)

// Manager handles profile operations
type Manager struct {
	basePath     string
	profiles     map[string]*Profile
	activeID     string
	mu           sync.RWMutex
	validator    *Validator
	stateManager *state.Manager
}

// NewManager creates a new profile manager
func NewManager(basePath string) *Manager {
	// Get parent directory for state
	parentDir := filepath.Dir(basePath)

	return &Manager{
		basePath:     basePath,
		profiles:     make(map[string]*Profile),
		validator:    NewValidator(),
		stateManager: state.NewManager(parentDir),
	}
}

// Initialize sets up the profile directory and loads profiles
func (m *Manager) Initialize() error {
	// Ensure profiles directory exists
	if err := os.MkdirAll(m.basePath, 0755); err != nil {
		return fmt.Errorf("creating profiles directory: %w", err)
	}

	// Create default profile if it doesn't exist
	defaultPath := filepath.Join(m.basePath, "default.yaml")
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		if err := m.createDefaultProfile(); err != nil {
			return fmt.Errorf("creating default profile: %w", err)
		}
	}

	// Load all profiles
	if err := m.LoadProfiles(); err != nil {
		return err
	}

	// Load saved active profile
	savedActiveID, err := m.stateManager.GetActiveProfile()
	if err != nil {
		// Log but don't fail initialization
		fmt.Printf("Warning: failed to load saved active profile: %v\n", err)
	} else if savedActiveID != "" {
		// Verify the saved profile still exists
		m.mu.Lock()
		if _, exists := m.profiles[savedActiveID]; exists {
			m.activeID = savedActiveID
		}
		m.mu.Unlock()
	}

	return nil
}

// createDefaultProfile creates the default profile
func (m *Manager) createDefaultProfile() error {
	profile := &Profile{
		ID:          "default",
		Name:        "Default",
		Description: "Default profile",
		Extensions:  []ExtensionRef{},
		Environment: make(map[string]string),
		MCPServers:  make(map[string]ServerConfig),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return m.Save(profile)
}

// LoadProfiles loads all profiles from disk
func (m *Manager) LoadProfiles() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries, err := os.ReadDir(m.basePath)
	if err != nil {
		return fmt.Errorf("reading profiles directory: %w", err)
	}

	m.profiles = make(map[string]*Profile)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		profilePath := filepath.Join(m.basePath, entry.Name())
		profile, err := m.loadProfile(profilePath)
		if err != nil {
			fmt.Printf("Warning: failed to load profile %s: %v\n", entry.Name(), err)
			continue
		}

		m.profiles[profile.ID] = profile
	}

	// Set default as active if no active profile
	if m.activeID == "" && len(m.profiles) > 0 {
		if _, exists := m.profiles["default"]; exists {
			m.activeID = "default"
		}
	}

	return nil
}

// loadProfile loads a single profile from disk
func (m *Manager) loadProfile(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading profile file: %w", err)
	}

	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("parsing profile YAML: %w", err)
	}

	// Validate profile
	if err := m.validator.Validate(&profile); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &profile, nil
}

// Save saves a profile to disk
func (m *Manager) Save(profile *Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.saveInternal(profile)
}

// saveInternal saves a profile without locking (must be called with lock held)
func (m *Manager) saveInternal(profile *Profile) error {
	// Preserve CreatedAt if it exists in cache
	if existing, exists := m.profiles[profile.ID]; exists && !existing.CreatedAt.IsZero() {
		profile.CreatedAt = existing.CreatedAt
	}

	// Update timestamp
	profile.UpdatedAt = time.Now()

	// Validate before saving
	if err := m.validator.Validate(profile); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(profile)
	if err != nil {
		return fmt.Errorf("marshaling profile: %w", err)
	}

	// Write to file atomically
	profilePath := filepath.Join(m.basePath, profile.ID+".yaml")
	tempPath := profilePath + ".tmp"

	// Write to temporary file first
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("writing temporary profile file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, profilePath); err != nil {
		// Clean up temp file on error
		os.Remove(tempPath)
		return fmt.Errorf("renaming profile file: %w", err)
	}

	// Update in-memory cache
	m.profiles[profile.ID] = profile

	return nil
}

// List returns all profiles
func (m *Manager) List() []*Profile {
	m.mu.RLock()
	defer m.mu.RUnlock()

	profiles := make([]*Profile, 0, len(m.profiles))
	for _, p := range m.profiles {
		profiles = append(profiles, p)
	}
	return profiles
}

// Get returns a profile by ID
func (m *Manager) Get(id string) (*Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	profile, exists := m.profiles[id]
	if !exists {
		return nil, fmt.Errorf("profile not found: %s", id)
	}
	return profile, nil
}

// GetActive returns the currently active profile
func (m *Manager) GetActive() (*Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeID == "" {
		return nil, fmt.Errorf("no active profile")
	}

	profile, exists := m.profiles[m.activeID]
	if !exists {
		return nil, fmt.Errorf("active profile not found: %s", m.activeID)
	}
	return profile, nil
}

// SetActive sets the active profile
func (m *Manager) SetActive(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	profile, exists := m.profiles[id]
	if !exists {
		return fmt.Errorf("profile not found: %s", id)
	}

	m.activeID = id

	// Save active profile to persistent state
	if err := m.stateManager.SetActiveProfile(id); err != nil {
		// Log but don't fail the operation
		fmt.Printf("Warning: failed to save active profile state: %v\n", err)
	}

	// Update last used timestamp
	now := time.Now()
	profile.LastUsed = &now
	profile.UsageCount++

	// Save the updated profile (use internal method since we already have the lock)
	return m.saveInternal(profile)
}

// Create creates a new profile
func (m *Manager) Create(profile *Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if profile already exists
	if _, exists := m.profiles[profile.ID]; exists {
		return fmt.Errorf("profile already exists: %s", profile.ID)
	}

	// Set timestamps
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	// Save to disk (use internal method since we already have the lock)
	// saveInternal will also add to cache
	return m.saveInternal(profile)
}

// Delete removes a profile
func (m *Manager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Don't allow deleting the default profile
	if id == "default" {
		return fmt.Errorf("cannot delete default profile")
	}

	// Don't allow deleting the active profile
	if id == m.activeID {
		return fmt.Errorf("cannot delete active profile")
	}

	if _, exists := m.profiles[id]; !exists {
		return fmt.Errorf("profile not found: %s", id)
	}

	// Remove file
	profilePath := filepath.Join(m.basePath, id+".yaml")
	if err := os.Remove(profilePath); err != nil {
		return fmt.Errorf("removing profile file: %w", err)
	}

	// Remove from cache
	delete(m.profiles, id)

	return nil
}

// Clone creates a copy of an existing profile
func (m *Manager) Clone(sourceID, newID, newName string) error {
	source, err := m.Get(sourceID)
	if err != nil {
		return err
	}

	// Create a deep copy
	cloned := &Profile{
		ID:          newID,
		Name:        newName,
		Description: fmt.Sprintf("Cloned from %s", source.Name),
		Icon:        source.Icon,
		Color:       source.Color,
		Extensions:  make([]ExtensionRef, len(source.Extensions)),
		Environment: make(map[string]string),
		MCPServers:  make(map[string]ServerConfig),
		Inherits:    source.Inherits,
		Tags:        append([]string{}, source.Tags...),
	}

	// Copy extensions
	copy(cloned.Extensions, source.Extensions)

	// Copy environment
	for k, v := range source.Environment {
		cloned.Environment[k] = v
	}

	// Copy MCP servers
	for k, v := range source.MCPServers {
		cloned.MCPServers[k] = v
	}

	return m.Create(cloned)
}
