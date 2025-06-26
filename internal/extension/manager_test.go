package extension

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestExtension(t *testing.T, dir, id string) string {
	extDir := filepath.Join(dir, id)
	if err := os.MkdirAll(extDir, 0755); err != nil {
		t.Fatalf("Failed to create extension dir: %v", err)
	}

	manifest := Extension{
		Name:        id,
		Version:     "1.0.0",
		Description: "Test extension " + id,
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}

	manifestPath := filepath.Join(extDir, "gemini-extension.json")
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	return extDir
}

func TestManager_Scan(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ext-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("creates directory if not exists", func(t *testing.T) {
		nonExistentDir := filepath.Join(tmpDir, "new-dir")
		m := NewManager(nonExistentDir)
		if err := m.Scan(); err != nil {
			t.Errorf("Scan() error = %v", err)
		}
		if _, err := os.Stat(nonExistentDir); os.IsNotExist(err) {
			t.Error("Directory was not created")
		}
	})

	t.Run("handles permission errors", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}
		
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		os.MkdirAll(readOnlyDir, 0755)
		os.Chmod(readOnlyDir, 0444)
		defer os.Chmod(readOnlyDir, 0755)

		m := NewManager(filepath.Join(readOnlyDir, "subdir"))
		err := m.Scan()
		if err == nil {
			t.Error("Expected permission error")
		}
	})
}

func TestManager_ScanExtensions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ext-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	// Create test extensions
	createTestExtension(t, tmpDir, "ext1")
	createTestExtension(t, tmpDir, "ext2")
	createTestExtension(t, tmpDir, "ext3")

	// Create invalid extension (no manifest)
	invalidDir := filepath.Join(tmpDir, "invalid")
	os.MkdirAll(invalidDir, 0755)

	// Create non-directory file
	os.WriteFile(filepath.Join(tmpDir, "notadir.txt"), []byte("test"), 0644)

	err = manager.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Verify extensions were loaded
	exts := manager.List()
	if len(exts) != 3 {
		t.Errorf("List() returned %d extensions, want 3", len(exts))
	}

	// Verify IDs were set from directory names
	for _, ext := range exts {
		if !strings.HasPrefix(ext.ID, "ext") {
			t.Errorf("Extension ID %q doesn't match expected pattern", ext.ID)
		}
	}
}

func TestManager_Get(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ext-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)
	createTestExtension(t, tmpDir, "test-ext")
	manager.Scan()

	t.Run("existing extension", func(t *testing.T) {
		ext, err := manager.Get("test-ext")
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if ext == nil {
			t.Error("Get() returned nil extension")
		}
		if ext.ID != "test-ext" {
			t.Errorf("Get() returned extension with ID %q, want %q", ext.ID, "test-ext")
		}
	})

	t.Run("non-existent extension", func(t *testing.T) {
		_, err := manager.Get("non-existent")
		if err == nil {
			t.Error("Expected error for non-existent extension")
		}
	})
}

func TestManager_Remove(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ext-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)
	
	extPath := createTestExtension(t, tmpDir, "removable-ext")
	manager.Scan()

	// Verify extension exists
	if _, err := manager.Get("removable-ext"); err != nil {
		t.Fatalf("Extension not found before removal: %v", err)
	}

	// Remove extension
	if err := manager.Remove("removable-ext"); err != nil {
		t.Errorf("Remove() error = %v", err)
	}

	// Verify extension is gone from manager
	if _, err := manager.Get("removable-ext"); err == nil {
		t.Error("Extension still exists in manager after removal")
	}

	// Verify directory was moved to trash
	trashPath := filepath.Join(tmpDir, ".trash")
	entries, err := os.ReadDir(trashPath)
	if err != nil {
		t.Fatalf("Failed to read trash directory: %v", err)
	}

	found := false
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "removable-ext") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Extension not found in trash")
	}

	// Verify original path no longer exists
	if _, err := os.Stat(extPath); !os.IsNotExist(err) {
		t.Error("Original extension directory still exists")
	}
}

func TestManager_ConcurrentAccess(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ext-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)
	

	// Create multiple extensions
	for i := 0; i < 10; i++ {
		createTestExtension(t, tmpDir, fmt.Sprintf("ext%d", i))
	}
	manager.Scan()

	// Test concurrent reads and writes
	done := make(chan bool)
	errors := make(chan error, 100)

	// Concurrent readers
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < 20; j++ {
				ext, err := manager.Get(fmt.Sprintf("ext%d", j%10))
				if err != nil {
					errors <- err
					return
				}
				if ext == nil {
					errors <- fmt.Errorf("Get returned nil")
					return
				}
				
				// List extensions
				exts := manager.List()
				if len(exts) == 0 {
					errors <- fmt.Errorf("List returned empty")
					return
				}
			}
		}(i)
	}

	// Concurrent state changes
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer func() { done <- true }()
			extID := fmt.Sprintf("ext%d", id)
			
			// Toggle enable/disable
			for j := 0; j < 10; j++ {
				if j%2 == 0 {
					if err := manager.Enable(extID); err != nil {
						errors <- err
						return
					}
				} else {
					if err := manager.Disable(extID); err != nil {
						errors <- err
						return
					}
				}
			}
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 8; i++ {
		<-done
	}

	// Check for errors
	select {
	case err := <-errors:
		t.Errorf("Concurrent access error: %v", err)
	default:
		// No errors
	}
}

func TestManager_EdgeCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ext-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)
	

	t.Run("remove non-existent extension", func(t *testing.T) {
		err := manager.Remove("non-existent")
		if err == nil {
			t.Error("Expected error when removing non-existent extension")
		}
	})

	t.Run("enable/disable non-existent extension", func(t *testing.T) {
		err := manager.Enable("non-existent")
		if err == nil {
			t.Error("Expected error when enabling non-existent extension")
		}

		err = manager.Disable("non-existent")
		if err == nil {
			t.Error("Expected error when disabling non-existent extension")
		}
	})

	t.Run("scan with corrupted manifest", func(t *testing.T) {
		corruptDir := filepath.Join(tmpDir, "corrupt")
		os.MkdirAll(corruptDir, 0755)
		manifestPath := filepath.Join(corruptDir, "gemini-extension.json")
		os.WriteFile(manifestPath, []byte("{invalid json}"), 0644)

		err := manager.Scan()
		// Should not fail completely, just skip the corrupted one
		if err != nil {
			t.Errorf("Scan should not fail with corrupted manifest: %v", err)
		}
		
		// Should have 0 extensions loaded
		if len(manager.List()) != 0 {
			t.Errorf("Should have 0 extensions with corrupted manifest, got %d", len(manager.List()))
		}
	})

	t.Run("extension with very long path", func(t *testing.T) {
		// Create nested directories
		longPath := tmpDir
		for i := 0; i < 20; i++ {
			longPath = filepath.Join(longPath, fmt.Sprintf("very-long-directory-name-%d", i))
		}
		
		// This might fail on some systems due to path length limits
		err := os.MkdirAll(longPath, 0755)
		if err != nil {
			t.Skipf("Skipping long path test: %v", err)
		}

		// Try to use it as extension directory
		m := NewManager(longPath)
		err = m.Scan()
		if err != nil && !strings.Contains(err.Error(), "too long") {
			// We expect it might fail, but for path length reasons
			t.Logf("Long path scan: %v", err)
		}
	})
}