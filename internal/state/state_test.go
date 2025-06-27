package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateManager(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "state-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	t.Run("Save and load active profile", func(t *testing.T) {
		// Set active profile
		testProfileID := "test-profile-123"
		err := manager.SetActiveProfile(testProfileID)
		if err != nil {
			t.Errorf("Failed to set active profile: %v", err)
		}

		// Load it back
		loadedID, err := manager.GetActiveProfile()
		if err != nil {
			t.Errorf("Failed to get active profile: %v", err)
		}

		if loadedID != testProfileID {
			t.Errorf("Loaded profile ID = %s, want %s", loadedID, testProfileID)
		}

		// Verify file exists
		statePath := filepath.Join(tmpDir, "state.json")
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			t.Errorf("State file was not created")
		}
	})

	t.Run("Handle missing state file", func(t *testing.T) {
		// Create new manager with non-existent file
		newManager := NewManager(filepath.Join(tmpDir, "nonexistent"))

		// Should return empty string, not error
		id, err := newManager.GetActiveProfile()
		if err != nil {
			t.Errorf("Expected no error for missing file, got: %v", err)
		}
		if id != "" {
			t.Errorf("Expected empty profile ID, got: %s", id)
		}
	})
}
