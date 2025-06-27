package launcher

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
)

func TestEnvironmentPreparer_PrepareExtensions(t *testing.T) {
	// Create temp directories
	tmpDir, err := os.MkdirTemp("", "env-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	managerExtDir := filepath.Join(tmpDir, "manager", "extensions")
	geminiExtDir := filepath.Join(tmpDir, "gemini", "extensions")

	os.MkdirAll(managerExtDir, 0755)
	os.MkdirAll(geminiExtDir, 0755)

	preparer := &EnvironmentPreparer{
		managerExtDir: managerExtDir,
		geminiExtDir:  geminiExtDir,
	}

	// Create test extensions in manager directory
	ext1Dir := filepath.Join(managerExtDir, "ext1")
	ext2Dir := filepath.Join(managerExtDir, "ext2")
	os.MkdirAll(ext1Dir, 0755)
	os.MkdirAll(ext2Dir, 0755)
	os.WriteFile(filepath.Join(ext1Dir, "test.txt"), []byte("ext1"), 0644)
	os.WriteFile(filepath.Join(ext2Dir, "test.txt"), []byte("ext2"), 0644)

	extensions := []*extension.Extension{
		{ID: "ext1", Path: ext1Dir},
		{ID: "ext2", Path: ext2Dir},
	}

	t.Run("prepare extensions creates symlinks", func(t *testing.T) {
		err := preparer.PrepareExtensions(extensions)
		if err != nil {
			t.Fatalf("PrepareExtensions() error = %v", err)
		}

		// Verify symlinks were created
		for _, ext := range extensions {
			linkPath := filepath.Join(geminiExtDir, ext.ID)
			info, err := os.Lstat(linkPath)
			if err != nil {
				t.Errorf("Symlink not created for %s: %v", ext.ID, err)
				continue
			}

			if info.Mode()&os.ModeSymlink == 0 {
				t.Errorf("Path %s is not a symlink", linkPath)
			}

			// Verify symlink points to correct location
			target, err := os.Readlink(linkPath)
			if err != nil {
				t.Errorf("Cannot read symlink %s: %v", linkPath, err)
				continue
			}

			if target != ext.Path {
				t.Errorf("Symlink target = %q, want %q", target, ext.Path)
			}
		}
	})

	t.Run("prepare extensions cleans old symlinks", func(t *testing.T) {
		// Create an old symlink that should be removed
		oldLink := filepath.Join(geminiExtDir, "old-ext")
		os.Symlink("/nonexistent/path", oldLink)

		// Prepare with only ext1
		err := preparer.PrepareExtensions([]*extension.Extension{
			{ID: "ext1", Path: ext1Dir},
		})
		if err != nil {
			t.Fatalf("PrepareExtensions() error = %v", err)
		}

		// Verify old symlink was removed
		if _, err := os.Lstat(oldLink); !os.IsNotExist(err) {
			t.Error("Old symlink was not removed")
		}

		// Verify ext2 symlink was also removed
		ext2Link := filepath.Join(geminiExtDir, "ext2")
		if _, err := os.Lstat(ext2Link); !os.IsNotExist(err) {
			t.Error("Ext2 symlink was not removed")
		}
	})

	t.Run("prepare with no extensions", func(t *testing.T) {
		// First create some symlinks
		preparer.PrepareExtensions(extensions)

		// Now prepare with empty list
		err := preparer.PrepareExtensions([]*extension.Extension{})
		if err != nil {
			t.Fatalf("PrepareExtensions() error = %v", err)
		}

		// Verify all symlinks were removed
		entries, err := os.ReadDir(geminiExtDir)
		if err != nil {
			t.Fatalf("Failed to read directory: %v", err)
		}

		for _, entry := range entries {
			if entry.Type()&os.ModeSymlink != 0 {
				t.Errorf("Symlink %s was not removed", entry.Name())
			}
		}
	})
}

func TestEnvironmentPreparer_EdgeCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "env-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	managerExtDir := filepath.Join(tmpDir, "manager", "extensions")
	geminiExtDir := filepath.Join(tmpDir, "gemini", "extensions")

	os.MkdirAll(managerExtDir, 0755)

	t.Run("gemini directory doesn't exist", func(t *testing.T) {
		preparer := &EnvironmentPreparer{
			managerExtDir: managerExtDir,
			geminiExtDir:  geminiExtDir,
		}

		extDir := filepath.Join(managerExtDir, "test-ext")
		os.MkdirAll(extDir, 0755)

		extensions := []*extension.Extension{
			{ID: "test-ext", Path: extDir},
		}

		// Should create directory
		err := preparer.PrepareExtensions(extensions)
		if err != nil {
			t.Fatalf("PrepareExtensions() error = %v", err)
		}

		if _, err := os.Stat(geminiExtDir); os.IsNotExist(err) {
			t.Error("Gemini extensions directory was not created")
		}
	})

	t.Run("existing non-symlink file", func(t *testing.T) {
		os.MkdirAll(geminiExtDir, 0755)
		preparer := &EnvironmentPreparer{
			managerExtDir: managerExtDir,
			geminiExtDir:  geminiExtDir,
		}

		// Create a regular file where symlink should go
		conflictPath := filepath.Join(geminiExtDir, "conflict-ext")
		os.WriteFile(conflictPath, []byte("existing file"), 0644)

		extDir := filepath.Join(managerExtDir, "conflict-ext")
		os.MkdirAll(extDir, 0755)

		extensions := []*extension.Extension{
			{ID: "conflict-ext", Path: extDir},
		}

		err := preparer.PrepareExtensions(extensions)
		if err == nil {
			t.Error("Expected error when non-symlink file exists")
		}
	})

	t.Run("broken symlink in gemini dir", func(t *testing.T) {
		os.MkdirAll(geminiExtDir, 0755)
		preparer := &EnvironmentPreparer{
			managerExtDir: managerExtDir,
			geminiExtDir:  geminiExtDir,
		}

		// Create broken symlink
		brokenLink := filepath.Join(geminiExtDir, "broken")
		os.Symlink("/nonexistent/path", brokenLink)

		// Should remove broken symlink
		err := preparer.PrepareExtensions([]*extension.Extension{})
		if err != nil {
			t.Fatalf("PrepareExtensions() error = %v", err)
		}

		if _, err := os.Lstat(brokenLink); !os.IsNotExist(err) {
			t.Error("Broken symlink was not removed")
		}
	})

	t.Run("permission denied", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping permission test on Windows")
		}
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		// Create read-only gemini directory
		roGeminiDir := filepath.Join(tmpDir, "readonly-gemini")
		os.MkdirAll(roGeminiDir, 0755)
		os.Chmod(roGeminiDir, 0444)
		defer os.Chmod(roGeminiDir, 0755)

		preparer := &EnvironmentPreparer{
			managerExtDir: managerExtDir,
			geminiExtDir:  roGeminiDir,
		}

		extDir := filepath.Join(managerExtDir, "test-ext")
		os.MkdirAll(extDir, 0755)

		extensions := []*extension.Extension{
			{ID: "test-ext", Path: extDir},
		}

		err := preparer.PrepareExtensions(extensions)
		if err == nil {
			t.Error("Expected permission error")
		}
	})

	t.Run("extension path doesn't exist", func(t *testing.T) {
		os.MkdirAll(geminiExtDir, 0755)
		preparer := &EnvironmentPreparer{
			managerExtDir: managerExtDir,
			geminiExtDir:  geminiExtDir,
		}

		extensions := []*extension.Extension{
			{ID: "nonexistent", Path: "/nonexistent/path"},
		}

		// Should handle gracefully
		err := preparer.PrepareExtensions(extensions)
		if err != nil {
			t.Logf("PrepareExtensions() with nonexistent path: %v", err)
		}

		// Symlink should still be created even if target doesn't exist
		linkPath := filepath.Join(geminiExtDir, "nonexistent")
		if _, err := os.Lstat(linkPath); os.IsNotExist(err) {
			t.Error("Symlink was not created for nonexistent path")
		}
	})

	t.Run("circular symlink", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping symlink test on Windows")
		}

		os.MkdirAll(geminiExtDir, 0755)
		preparer := &EnvironmentPreparer{
			managerExtDir: managerExtDir,
			geminiExtDir:  geminiExtDir,
		}

		// Create circular symlinks
		link1 := filepath.Join(geminiExtDir, "circular1")
		link2 := filepath.Join(geminiExtDir, "circular2")
		os.Symlink(link2, link1)
		os.Symlink(link1, link2)

		// Should handle cleanup without infinite loop
		err := preparer.PrepareExtensions([]*extension.Extension{})
		if err != nil {
			t.Fatalf("Failed to handle circular symlinks: %v", err)
		}

		// Both should be removed
		if _, err := os.Lstat(link1); !os.IsNotExist(err) {
			t.Error("Circular symlink 1 was not removed")
		}
		if _, err := os.Lstat(link2); !os.IsNotExist(err) {
			t.Error("Circular symlink 2 was not removed")
		}
	})
}
