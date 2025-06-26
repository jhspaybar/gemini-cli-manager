package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State represents the persistent application state
type State struct {
	ActiveProfileID string    `json:"activeProfileID"`
	LastUpdated     time.Time `json:"lastUpdated"`
}

// Manager handles persistent state operations
type Manager struct {
	statePath string
}

// NewManager creates a new state manager
func NewManager(basePath string) *Manager {
	return &Manager{
		statePath: filepath.Join(basePath, "state.json"),
	}
}

// Load reads the state from disk
func (m *Manager) Load() (*State, error) {
	data, err := os.ReadFile(m.statePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty state if file doesn't exist
			return &State{}, nil
		}
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}

	return &state, nil
}

// Save writes the state to disk
func (m *Manager) Save(state *State) error {
	state.LastUpdated = time.Now()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(m.statePath), 0755); err != nil {
		return fmt.Errorf("creating state directory: %w", err)
	}

	// Write to temporary file first for atomic update
	tempPath := m.statePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("writing temporary state file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, m.statePath); err != nil {
		// Clean up temp file on error
		os.Remove(tempPath)
		return fmt.Errorf("renaming state file: %w", err)
	}

	return nil
}

// SetActiveProfile updates and saves the active profile ID
func (m *Manager) SetActiveProfile(profileID string) error {
	state := &State{
		ActiveProfileID: profileID,
	}
	return m.Save(state)
}

// GetActiveProfile returns the active profile ID from saved state
func (m *Manager) GetActiveProfile() (string, error) {
	state, err := m.Load()
	if err != nil {
		return "", err
	}
	return state.ActiveProfileID, nil
}