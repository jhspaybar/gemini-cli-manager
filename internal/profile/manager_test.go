package profile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProfileManager_Create(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "profile-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create manager
	manager := NewManager(tmpDir)
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	// Test creating a new profile
	t.Run("Create new profile", func(t *testing.T) {
		profile := &Profile{
			ID:          "test-profile",
			Name:        "Test Profile",
			Description: "A test profile",
			Extensions: []ExtensionRef{
				{ID: "ext1", Enabled: true},
				{ID: "ext2", Enabled: true},
			},
			Environment: map[string]string{
				"TEST_VAR": "test_value",
			},
		}

		err := manager.Create(profile)
		if err != nil {
			t.Errorf("Failed to create profile: %v", err)
		}

		// Verify the profile was saved to disk
		profilePath := filepath.Join(tmpDir, "test-profile.yaml")
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			t.Errorf("Profile file was not created at %s", profilePath)
		}

		// Verify we can retrieve it
		retrieved, err := manager.Get("test-profile")
		if err != nil {
			t.Errorf("Failed to get profile: %v", err)
		}

		if retrieved.Name != profile.Name {
			t.Errorf("Retrieved profile name = %s, want %s", retrieved.Name, profile.Name)
		}

		if len(retrieved.Extensions) != 2 {
			t.Errorf("Retrieved profile has %d extensions, want 2", len(retrieved.Extensions))
		}
	})
}

func TestProfileManager_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "profile-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create manager
	manager := NewManager(tmpDir)
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	t.Run("Save updates existing profile", func(t *testing.T) {
		// First create a profile
		profile := &Profile{
			ID:          "update-test",
			Name:        "Original Name",
			Description: "Original description",
			Extensions:  []ExtensionRef{},
		}

		if err := manager.Create(profile); err != nil {
			t.Fatalf("Failed to create initial profile: %v", err)
		}

		// Now update it
		profile.Name = "Updated Name"
		profile.Description = "Updated description"
		profile.Extensions = append(profile.Extensions, ExtensionRef{ID: "new-ext", Enabled: true})

		if err := manager.Save(profile); err != nil {
			t.Errorf("Failed to save profile: %v", err)
		}

		// Reload from disk to verify
		manager.profiles = make(map[string]*Profile) // Clear cache
		if err := manager.LoadProfiles(); err != nil {
			t.Fatalf("Failed to reload profiles: %v", err)
		}

		retrieved, err := manager.Get("update-test")
		if err != nil {
			t.Errorf("Failed to get updated profile: %v", err)
		}

		if retrieved.Name != "Updated Name" {
			t.Errorf("Profile name = %s, want %s", retrieved.Name, "Updated Name")
		}

		if retrieved.Description != "Updated description" {
			t.Errorf("Profile description = %s, want %s", retrieved.Description, "Updated description")
		}

		if len(retrieved.Extensions) != 1 {
			t.Errorf("Profile has %d extensions, want 1", len(retrieved.Extensions))
		}
	})
}

func TestProfileManager_CreateWithExtensions(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "profile-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create manager
	manager := NewManager(tmpDir)
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	t.Run("Create profile with multiple extensions", func(t *testing.T) {
		profile := &Profile{
			ID:          "multi-ext",
			Name:        "Multi Extension Profile",
			Description: "Profile with multiple extensions",
			Extensions: []ExtensionRef{
				{ID: "simple-extension", Enabled: true},
				{ID: "mcp-extension", Enabled: true},
				{ID: "disabled-ext", Enabled: false},
			},
			Environment: map[string]string{
				"NODE_ENV": "development",
				"DEBUG":    "true",
			},
		}

		if err := manager.Create(profile); err != nil {
			t.Errorf("Failed to create profile: %v", err)
		}

		// Reload and verify
		manager.profiles = make(map[string]*Profile)
		if err := manager.LoadProfiles(); err != nil {
			t.Fatalf("Failed to reload profiles: %v", err)
		}

		retrieved, err := manager.Get("multi-ext")
		if err != nil {
			t.Errorf("Failed to get profile: %v", err)
		}

		// Verify extensions
		if len(retrieved.Extensions) != 3 {
			t.Errorf("Profile has %d extensions, want 3", len(retrieved.Extensions))
		}

		// Check specific extensions
		foundSimple := false
		foundMCP := false
		enabledCount := 0

		for _, ext := range retrieved.Extensions {
			if ext.ID == "simple-extension" {
				foundSimple = true
			}
			if ext.ID == "mcp-extension" {
				foundMCP = true
			}
			if ext.Enabled {
				enabledCount++
			}
		}

		if !foundSimple {
			t.Error("simple-extension not found in profile")
		}
		if !foundMCP {
			t.Error("mcp-extension not found in profile")
		}
		if enabledCount != 2 {
			t.Errorf("Profile has %d enabled extensions, want 2", enabledCount)
		}

		// Verify environment
		if len(retrieved.Environment) != 2 {
			t.Errorf("Profile has %d environment vars, want 2", len(retrieved.Environment))
		}
		if retrieved.Environment["NODE_ENV"] != "development" {
			t.Errorf("NODE_ENV = %s, want development", retrieved.Environment["NODE_ENV"])
		}
	})
}

func TestProfileManager_Timestamps(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "profile-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create manager
	manager := NewManager(tmpDir)
	if err := manager.Initialize(); err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	t.Run("Timestamps are set correctly", func(t *testing.T) {
		beforeCreate := time.Now()

		profile := &Profile{
			ID:          "timestamp-test",
			Name:        "Timestamp Test",
			Description: "Testing timestamp handling",
		}

		if err := manager.Create(profile); err != nil {
			t.Fatalf("Failed to create profile: %v", err)
		}

		afterCreate := time.Now()

		// Get the profile
		retrieved, err := manager.Get("timestamp-test")
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}

		// Check CreatedAt
		if retrieved.CreatedAt.Before(beforeCreate) || retrieved.CreatedAt.After(afterCreate) {
			t.Errorf("CreatedAt %v is not between %v and %v", retrieved.CreatedAt, beforeCreate, afterCreate)
		}

		// Check UpdatedAt
		if retrieved.UpdatedAt.Before(beforeCreate) || retrieved.UpdatedAt.After(afterCreate) {
			t.Errorf("UpdatedAt %v is not between %v and %v", retrieved.UpdatedAt, beforeCreate, afterCreate)
		}

		// Save the original UpdatedAt for comparison
		originalUpdatedAt := retrieved.UpdatedAt

		// Update the profile
		time.Sleep(10 * time.Millisecond) // Ensure time difference
		beforeUpdate := time.Now()

		retrieved.Description = "Updated description"
		if err := manager.Save(retrieved); err != nil {
			t.Fatalf("Failed to save profile: %v", err)
		}

		afterUpdate := time.Now()

		// Reload profiles from disk to ensure we get fresh data
		manager.profiles = make(map[string]*Profile)
		if err := manager.LoadProfiles(); err != nil {
			t.Fatalf("Failed to reload profiles: %v", err)
		}

		// Get updated profile
		updated, err := manager.Get("timestamp-test")
		if err != nil {
			t.Fatalf("Failed to get updated profile: %v", err)
		}

		// CreatedAt should not change
		if !updated.CreatedAt.Equal(retrieved.CreatedAt) {
			t.Errorf("CreatedAt changed from %v to %v", retrieved.CreatedAt, updated.CreatedAt)
		}

		// UpdatedAt should be updated
		if updated.UpdatedAt.Before(beforeUpdate) || updated.UpdatedAt.After(afterUpdate) {
			t.Errorf("UpdatedAt %v is not between %v and %v", updated.UpdatedAt, beforeUpdate, afterUpdate)
		}

		if !updated.UpdatedAt.After(originalUpdatedAt) {
			t.Errorf("UpdatedAt %v is not after original %v", updated.UpdatedAt, originalUpdatedAt)
		}
	})
}
