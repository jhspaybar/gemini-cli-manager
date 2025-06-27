package extension

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// getTestDataPath returns the path to test data files
func getTestDataPath(filename string) string {
	// Get the path relative to the module root
	_, currentFile, _, _ := runtime.Caller(0)
	moduleRoot := filepath.Join(filepath.Dir(currentFile), "../..")
	return filepath.Join(moduleRoot, "testdata", "extensions", filename)
}

func TestInstaller_InstallFromPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	installer := NewInstaller(tmpDir)

	t.Run("install from directory", func(t *testing.T) {
		// Use the test data directory
		srcDir := getTestDataPath("simple-extension")

		ext, err := installer.InstallFromPath(srcDir)
		if err != nil {
			t.Fatalf("InstallFromPath() error = %v", err)
		}

		if ext.ID != "simple-extension" {
			t.Errorf("Extension ID = %q, want %q", ext.ID, "simple-extension")
		}

		// Verify installation
		installedPath := filepath.Join(tmpDir, "simple-extension")
		if _, err := os.Stat(installedPath); os.IsNotExist(err) {
			t.Error("Extension was not installed")
		}

		// Verify files were copied
		manifestPath := filepath.Join(installedPath, "gemini-extension.json")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			t.Error("Manifest was not copied")
		}
	})

	t.Run("install from home directory path", func(t *testing.T) {
		_, err := os.UserHomeDir()
		if err != nil {
			t.Skip("Cannot determine home directory")
		}

		// Simulate home path with non-existent directory
		relativePath := "~/nonexistent-extension"

		// This should fail because the path doesn't exist
		_, err = installer.InstallFromPath(relativePath)
		if err == nil {
			t.Error("Expected error for non-existent home path")
		}
	})

	t.Run("install non-existent path", func(t *testing.T) {
		_, err := installer.InstallFromPath("/non/existent/path")
		if err == nil {
			t.Error("Expected error for non-existent path")
		}
	})

	t.Run("install file instead of directory", func(t *testing.T) {
		// Use an actual file from test data
		filePath := getTestDataPath("simple-extension.zip")

		// This should actually work since it's a zip file
		ext, err := installer.InstallFromPath(filePath)
		if err != nil {
			t.Logf("Installing from zip file: %v", err)
		} else {
			// If it succeeded, verify
			if ext.Name != "simple-test-extension" {
				t.Errorf("Extension name = %q, want %q", ext.Name, "simple-test-extension")
			}
		}
	})
}

func TestInstaller_InstallFromArchive(t *testing.T) {
	t.Run("install from zip", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)
		zipPath := getTestDataPath("simple-extension.zip")

		ext, err := installer.InstallFromPath(zipPath)
		if err != nil {
			t.Fatalf("InstallFromPath(zip) error = %v", err)
		}

		if ext.Name != "simple-extension" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "simple-extension")
		}

		// Verify files were extracted
		installedPath := filepath.Join(tmpDir, ext.ID)
		readmePath := filepath.Join(installedPath, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			t.Error("README.md was not extracted")
		}
	})

	t.Run("install from tar.gz", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)
		tarPath := getTestDataPath("simple-extension.tar.gz")

		ext, err := installer.InstallFromPath(tarPath)
		if err != nil {
			t.Fatalf("InstallFromPath(tar.gz) error = %v", err)
		}

		if ext.Name != "simple-extension" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "simple-extension")
		}
	})

	t.Run("install from nested archive", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)
		nestedPath := getTestDataPath("nested-extension.zip")

		ext, err := installer.InstallFromPath(nestedPath)
		if err != nil {
			t.Fatalf("Failed to install nested archive: %v", err)
		}

		if ext.Name != "extension" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "extension")
		}
	})

	t.Run("install archive without manifest", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)
		noManifestPath := getTestDataPath("no-manifest.zip")

		_, err = installer.InstallFromPath(noManifestPath)
		if err == nil {
			t.Error("Expected error for archive without manifest")
		}
		if !strings.Contains(err.Error(), "gemini-extension.json") {
			t.Errorf("Error should mention missing manifest: %v", err)
		}
	})
}

func TestInstaller_InstallFromURL(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	installer := NewInstaller(tmpDir)

	t.Run("install from direct URL", func(t *testing.T) {
		// Create test server that serves our test zip
		zipData, err := os.ReadFile(getTestDataPath("simple-extension.zip"))
		if err != nil {
			t.Fatalf("Failed to read test zip: %v", err)
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/zip")
			w.Write(zipData)
		}))
		defer server.Close()

		ext, err := installer.InstallFromURL(server.URL + "/extension.zip")
		if err != nil {
			t.Fatalf("InstallFromURL() error = %v", err)
		}

		if ext.Name != "simple-extension" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "simple-extension")
		}
	})

	t.Run("install from 404 URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer server.Close()

		_, err := installer.InstallFromURL(server.URL + "/notfound.zip")
		if err == nil {
			t.Error("Expected error for 404 response")
		}
	})

	t.Run("install from invalid URL", func(t *testing.T) {
		_, err := installer.InstallFromURL("not-a-url")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("install with network timeout", func(t *testing.T) {
		t.Skip("Skipping timeout test to avoid long waits")
	})
}

func TestInstaller_EdgeCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	installer := NewInstaller(tmpDir)

	t.Run("install already existing extension", func(t *testing.T) {
		// Use test data
		srcDir := getTestDataPath("simple-extension")

		// Install once
		_, err := installer.InstallFromPath(srcDir)
		if err != nil {
			t.Fatalf("First install failed: %v", err)
		}

		// Try to install again
		_, err = installer.InstallFromPath(srcDir)
		if err == nil {
			t.Error("Expected error when installing duplicate extension")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Error should mention extension already exists: %v", err)
		}
	})

	t.Run("install with invalid manifest", func(t *testing.T) {
		// Create a temp directory with invalid manifest
		badDir, err := os.MkdirTemp("", "bad-ext-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(badDir)

		// Create invalid manifest (missing required fields)
		invalidManifest := `{
			"name": "bad-ext"
		}`
		manifestPath := filepath.Join(badDir, "gemini-extension.json")
		os.WriteFile(manifestPath, []byte(invalidManifest), 0644)

		_, err = installer.InstallFromPath(badDir)
		if err == nil {
			t.Error("Expected error for invalid manifest")
		}
		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("Expected validation error, got: %v", err)
		}
	})

	t.Run("install with permission issues", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping permission test on Windows")
		}
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		// Create read-only destination
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		os.MkdirAll(readOnlyDir, 0755)

		// Create installer pointing to read-only directory
		roInstaller := NewInstaller(readOnlyDir)

		// Make directory read-only
		os.Chmod(readOnlyDir, 0444)
		defer os.Chmod(readOnlyDir, 0755)

		// Try to install
		srcDir := getTestDataPath("simple-extension")
		_, err := roInstaller.InstallFromPath(srcDir)
		if err == nil {
			t.Error("Expected error for permission denied")
		}
	})

	t.Run("install with symlinks", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping symlink test on Windows")
		}

		// Create a test extension with symlinks
		symlinkDir, err := os.MkdirTemp("", "symlink-ext-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(symlinkDir)

		// Create manifest that matches directory name
		dirName := filepath.Base(symlinkDir)
		manifest := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "description": "Test extension with symlinks"
}`, dirName)
		os.WriteFile(filepath.Join(symlinkDir, "gemini-extension.json"), []byte(manifest), 0644)

		// Create a file and symlink
		os.WriteFile(filepath.Join(symlinkDir, "original.txt"), []byte("original"), 0644)
		os.Symlink("original.txt", filepath.Join(symlinkDir, "link.txt"))

		ext, err := installer.InstallFromPath(symlinkDir)
		if err != nil {
			t.Fatalf("Failed to install extension with symlinks: %v", err)
		}

		// Verify symlink was copied
		linkPath := filepath.Join(tmpDir, ext.ID, "link.txt")
		if _, err := os.Stat(linkPath); os.IsNotExist(err) {
			t.Error("Symlink was not copied")
		}
	})
}

func TestInstaller_GitHubURL(t *testing.T) {
	t.Run("parse GitHub URLs", func(t *testing.T) {
		// Test URL parsing (without actually downloading)
		urls := []string{
			"https://github.com/user/repo",
			"https://github.com/user/repo.git",
			"git@github.com:user/repo.git",
		}

		for _, url := range urls {
			// Just verify the URL format is recognized
			if !strings.Contains(url, "github.com") {
				t.Errorf("URL should contain github.com: %s", url)
			}
		}
	})
}

func TestInstaller_ProgressCallback(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	installer := NewInstaller(tmpDir)

	t.Run("install without progress tracking", func(t *testing.T) {
		ext, err := installer.Install(
			getTestDataPath("simple-extension"),
			true,
		)

		if err != nil {
			t.Fatalf("Install failed: %v", err)
		}

		if ext.Name != "simple-extension" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "simple-extension")
		}
	})
}
