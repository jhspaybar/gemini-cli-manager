package extension

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func createTestArchive(t *testing.T, format string) string {
	tmpFile, err := os.CreateTemp("", "test-ext-*."+format)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	manifest := Extension{
		Name:        "archive-test",
		DisplayName: "Archive Test",
		Version:     "1.0.0",
		Description: "Test extension from archive",
	}
	manifestData, _ := json.Marshal(manifest)

	switch format {
	case "zip":
		// Reopen file for writing
		file, err := os.Create(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to reopen file: %v", err)
		}
		defer file.Close()
		
		w := zip.NewWriter(file)
		
		// Create directory
		_, err = w.Create("archive-test/")
		if err != nil {
			t.Fatalf("Failed to create zip directory: %v", err)
		}

		// Add manifest
		f, err := w.Create("archive-test/gemini-extension.json")
		if err != nil {
			t.Fatalf("Failed to create zip file: %v", err)
		}
		f.Write(manifestData)

		// Add another file
		f2, _ := w.Create("archive-test/README.md")
		f2.Write([]byte("# Test Extension"))

		w.Close()

	case "tar.gz":
		file, _ := os.Create(tmpFile.Name())
		defer file.Close()

		gw := gzip.NewWriter(file)
		defer gw.Close()

		tw := tar.NewWriter(gw)
		defer tw.Close()

		// Add directory
		hdr := &tar.Header{
			Name:     "archive-test/",
			Mode:     0755,
			Typeflag: tar.TypeDir,
		}
		tw.WriteHeader(hdr)

		// Add manifest
		hdr = &tar.Header{
			Name: "archive-test/gemini-extension.json",
			Mode: 0644,
			Size: int64(len(manifestData)),
		}
		tw.WriteHeader(hdr)
		tw.Write(manifestData)

		// Add another file
		readme := []byte("# Test Extension")
		hdr = &tar.Header{
			Name: "archive-test/README.md",
			Mode: 0644,
			Size: int64(len(readme)),
		}
		tw.WriteHeader(hdr)
		tw.Write(readme)
	}

	return tmpFile.Name()
}

func TestInstaller_InstallFromPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	installer := NewInstaller(tmpDir)

	t.Run("install from directory", func(t *testing.T) {
		// Create source extension
		srcDir := filepath.Join(tmpDir, "source")
		createTestExtension(t, srcDir, "test-ext")

		ext, err := installer.InstallFromPath(filepath.Join(srcDir, "test-ext"))
		if err != nil {
			t.Fatalf("InstallFromPath() error = %v", err)
		}

		if ext.ID != "test-ext" {
			t.Errorf("Extension ID = %q, want %q", ext.ID, "test-ext")
		}

		// Verify installation
		installedPath := filepath.Join(tmpDir, "test-ext")
		if _, err := os.Stat(installedPath); os.IsNotExist(err) {
			t.Error("Extension was not installed")
		}
	})

	t.Run("install from home directory path", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("Cannot determine home directory")
		}

		// Create extension in temp location
		srcDir := filepath.Join(tmpDir, "home-test")
		createTestExtension(t, srcDir, "home-ext")

		// Simulate home path
		relativePath := "~/" + filepath.Base(srcDir) + "/home-ext"
		
		// Temporarily change to a different directory to test path expansion
		oldWd, _ := os.Getwd()
		os.Chdir(home)
		defer os.Chdir(oldWd)

		// This should fail because the path doesn't actually exist in home
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
		// Create a regular file
		tmpFile := filepath.Join(tmpDir, "notadir.txt")
		os.WriteFile(tmpFile, []byte("test"), 0644)

		_, err := installer.InstallFromPath(tmpFile)
		if err == nil {
			t.Error("Expected error when installing from non-directory file")
		}
	})
}

func TestInstaller_InstallFromArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "installer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	installer := NewInstaller(tmpDir)

	t.Run("install from zip", func(t *testing.T) {
		zipPath := createTestArchive(t, "zip")
		defer os.Remove(zipPath)

		ext, err := installer.InstallFromPath(zipPath)
		if err != nil {
			t.Fatalf("InstallFromPath(zip) error = %v", err)
		}

		if ext.Name != "archive-test" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "archive-test")
		}

		// Verify files were extracted
		readmePath := filepath.Join(tmpDir, ext.ID, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			t.Error("README.md was not extracted")
		}
	})

	t.Run("install from tar.gz", func(t *testing.T) {
		tarPath := createTestArchive(t, "tar.gz")
		defer os.Remove(tarPath)

		ext, err := installer.InstallFromPath(tarPath)
		if err != nil {
			t.Fatalf("InstallFromPath(tar.gz) error = %v", err)
		}

		if ext.Name != "archive-test" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "archive-test")
		}
	})

	t.Run("install from nested archive", func(t *testing.T) {
		// Create archive with nested structure
		tmpFile, _ := os.CreateTemp("", "nested-*.zip")
		defer os.Remove(tmpFile.Name())

		w := zip.NewWriter(tmpFile)
		
		// Create nested structure
		w.Create("wrapper/")
		w.Create("wrapper/inner/")
		w.Create("wrapper/inner/extension/")
		
		// Put manifest in nested directory
		manifest := Extension{
			Name:        "nested-ext",
			DisplayName: "Nested Extension",
			Version:     "1.0.0",
			Description: "Nested extension",
		}
		manifestData, _ := json.Marshal(manifest)
		
		f, _ := w.Create("wrapper/inner/extension/gemini-extension.json")
		f.Write(manifestData)
		w.Close()
		tmpFile.Close()

		ext, err := installer.InstallFromPath(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to install nested archive: %v", err)
		}

		if ext.Name != "nested-ext" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "nested-ext")
		}
	})

	t.Run("install archive without manifest", func(t *testing.T) {
		// Create archive without manifest
		tmpFile, _ := os.CreateTemp("", "no-manifest-*.zip")
		defer os.Remove(tmpFile.Name())

		w := zip.NewWriter(tmpFile)
		f, _ := w.Create("test.txt")
		f.Write([]byte("no manifest here"))
		w.Close()
		tmpFile.Close()

		_, err := installer.InstallFromPath(tmpFile.Name())
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
		// Create test server
		zipPath := createTestArchive(t, "zip")
		defer os.Remove(zipPath)

		zipData, _ := os.ReadFile(zipPath)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/zip")
			w.Write(zipData)
		}))
		defer server.Close()

		ext, err := installer.InstallFromURL(server.URL + "/extension.zip")
		if err != nil {
			t.Fatalf("InstallFromURL() error = %v", err)
		}

		if ext.Name != "archive-test" {
			t.Errorf("Extension name = %q, want %q", ext.Name, "archive-test")
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
		// Create source
		srcDir := filepath.Join(tmpDir, "source")
		createTestExtension(t, srcDir, "duplicate")

		// Install once
		_, err := installer.InstallFromPath(filepath.Join(srcDir, "duplicate"))
		if err != nil {
			t.Fatalf("First install failed: %v", err)
		}

		// Try to install again
		_, err = installer.InstallFromPath(filepath.Join(srcDir, "duplicate"))
		if err == nil {
			t.Error("Expected error when installing duplicate extension")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Error should mention extension already exists: %v", err)
		}
	})

	t.Run("install with invalid manifest", func(t *testing.T) {
		// Create extension with invalid manifest
		badDir := filepath.Join(tmpDir, "bad-ext")
		os.MkdirAll(badDir, 0755)
		
		// Invalid manifest (missing required fields)
		manifest := map[string]string{
			"name": "bad-ext",
			// Missing version, displayName, description
		}
		data, _ := json.Marshal(manifest)
		os.WriteFile(filepath.Join(badDir, "gemini-extension.json"), data, 0644)

		_, err := installer.InstallFromPath(badDir)
		if err == nil {
			t.Error("Expected error for invalid manifest")
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
		srcDir := filepath.Join(tmpDir, "source2")
		createTestExtension(t, srcDir, "perm-test")

		_, err := roInstaller.InstallFromPath(filepath.Join(srcDir, "perm-test"))
		if err == nil {
			t.Error("Expected error for permission denied")
		}
	})

	t.Run("install with symlinks", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping symlink test on Windows")
		}

		// Create extension with symlink
		srcDir := filepath.Join(tmpDir, "symlink-ext")
		os.MkdirAll(srcDir, 0755)
		
		// Create manifest
		manifest := Extension{
			Name:        "symlink-ext",
			DisplayName: "Symlink Extension",
			Version:     "1.0.0",
			Description: "Extension with symlinks",
		}
		data, _ := json.Marshal(manifest)
		os.WriteFile(filepath.Join(srcDir, "gemini-extension.json"), data, 0644)

		// Create a file and symlink to it
		os.WriteFile(filepath.Join(srcDir, "original.txt"), []byte("original"), 0644)
		os.Symlink("original.txt", filepath.Join(srcDir, "link.txt"))

		ext, err := installer.InstallFromPath(srcDir)
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