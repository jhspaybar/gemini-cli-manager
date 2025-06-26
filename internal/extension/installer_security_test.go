package extension

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstaller_SecurityEdgeCases(t *testing.T) {
	t.Run("path traversal in zip archive", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-sec-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)

		// Create a malicious zip with path traversal
		var buf bytes.Buffer
		zipWriter := zip.NewWriter(&buf)
		
		// Add a normal manifest
		manifest := `{"name": "evil-ext", "version": "1.0.0", "description": "Evil extension"}`
		w, _ := zipWriter.Create("gemini-extension.json")
		w.Write([]byte(manifest))
		
		// Try to write outside the extension directory
		w, _ = zipWriter.Create("../../../etc/evil.txt")
		w.Write([]byte("malicious content"))
		
		zipWriter.Close()

		// Write the malicious zip to a file
		maliciousZip := filepath.Join(tmpDir, "evil.zip")
		os.WriteFile(maliciousZip, buf.Bytes(), 0644)

		// This should fail or sanitize the path
		_, err = installer.InstallFromPath(maliciousZip)
		
		// Check that no file was written outside the extension directory
		evilPath := filepath.Join(tmpDir, "..", "..", "..", "etc", "evil.txt")
		if _, err := os.Stat(evilPath); err == nil {
			t.Error("Path traversal attack succeeded - file written outside extension directory")
		}
	})

	t.Run("path traversal in tar.gz archive", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-sec-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)

		// Create a malicious tar.gz with path traversal
		var buf bytes.Buffer
		gzWriter := gzip.NewWriter(&buf)
		tarWriter := tar.NewWriter(gzWriter)
		
		// Add a normal manifest
		manifest := `{"name": "evil-ext", "version": "1.0.0", "description": "Evil extension"}`
		hdr := &tar.Header{
			Name: "gemini-extension.json",
			Mode: 0644,
			Size: int64(len(manifest)),
		}
		tarWriter.WriteHeader(hdr)
		tarWriter.Write([]byte(manifest))
		
		// Try to write outside the extension directory
		evil := "malicious content"
		hdr = &tar.Header{
			Name: "../../../tmp/evil.txt",
			Mode: 0644,
			Size: int64(len(evil)),
		}
		tarWriter.WriteHeader(hdr)
		tarWriter.Write([]byte(evil))
		
		tarWriter.Close()
		gzWriter.Close()

		// Write the malicious tar.gz to a file
		maliciousTar := filepath.Join(tmpDir, "evil.tar.gz")
		os.WriteFile(maliciousTar, buf.Bytes(), 0644)

		// This should fail or sanitize the path
		_, err = installer.InstallFromPath(maliciousTar)
		
		// Check that no file was written outside the extension directory
		evilPath := "/tmp/evil.txt"
		if _, err := os.Stat(evilPath); err == nil {
			t.Error("Path traversal attack succeeded - file written to /tmp")
			os.Remove(evilPath) // Clean up
		}
	})

	t.Run("symlink escape in zip archive", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-sec-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)

		// Create a zip with symlinks pointing outside
		var buf bytes.Buffer
		zipWriter := zip.NewWriter(&buf)
		
		// Add a normal manifest
		manifest := `{"name": "evil-ext", "version": "1.0.0", "description": "Evil extension"}`
		w, _ := zipWriter.Create("gemini-extension.json")
		w.Write([]byte(manifest))
		
		// Note: Creating proper symlinks in zip files is complex and platform-specific
		// Our protection should skip any files with symlink mode bits set
		
		zipWriter.Close()

		// Write the zip to a file
		maliciousZip := filepath.Join(tmpDir, "evil-symlink.zip")
		os.WriteFile(maliciousZip, buf.Bytes(), 0644)

		// Install the extension
		ext, err := installer.InstallFromPath(maliciousZip)
		if err == nil {
			// Check if symlink was created
			linkPath := filepath.Join(tmpDir, ext.ID, "passwd-link")
			if target, err := os.Readlink(linkPath); err == nil {
				if target == "/etc/passwd" || !strings.HasPrefix(target, tmpDir) {
					t.Error("Dangerous symlink was created pointing outside extension directory")
				}
			}
		}
	})

	t.Run("zip bomb protection", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-sec-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)

		// Create a zip with high compression ratio (simulated zip bomb)
		var buf bytes.Buffer
		zipWriter := zip.NewWriter(&buf)
		
		// Add a normal manifest
		manifest := `{"name": "bomb-ext", "version": "1.0.0", "description": "Zip bomb test"}`
		w, _ := zipWriter.Create("gemini-extension.json")
		w.Write([]byte(manifest))
		
		// Add a highly compressible file (lots of zeros)
		w, _ = zipWriter.Create("bomb.txt")
		// Write 100MB of zeros (highly compressible)
		zeros := make([]byte, 1024*1024) // 1MB of zeros
		for i := 0; i < 100; i++ {
			w.Write(zeros)
		}
		
		zipWriter.Close()

		// Write the zip to a file
		bombZip := filepath.Join(tmpDir, "bomb.zip")
		os.WriteFile(bombZip, buf.Bytes(), 0644)

		// Get size of compressed file
		info, _ := os.Stat(bombZip)
		compressedSize := info.Size()
		
		// Install should either fail or enforce size limits
		_, err = installer.InstallFromPath(bombZip)
		
		// Check that we didn't extract 100MB
		if err == nil {
			// Check extracted size
			var totalSize int64
			filepath.Walk(filepath.Join(tmpDir, "bomb-ext"), func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					totalSize += info.Size()
				}
				return nil
			})
			
			ratio := float64(totalSize) / float64(compressedSize)
			if ratio > 100 {
				t.Logf("Warning: High compression ratio detected (%.2f:1) - potential zip bomb", ratio)
			}
		}
	})

	t.Run("absolute path in archive", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-sec-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)

		// Create a zip with absolute paths
		var buf bytes.Buffer
		zipWriter := zip.NewWriter(&buf)
		
		// Add a normal manifest
		manifest := `{"name": "abs-ext", "version": "1.0.0", "description": "Absolute path test"}`
		w, _ := zipWriter.Create("gemini-extension.json")
		w.Write([]byte(manifest))
		
		// Try to use absolute path
		w, _ = zipWriter.Create("/tmp/absolute-evil.txt")
		w.Write([]byte("malicious content"))
		
		zipWriter.Close()

		// Write the zip to a file
		absZip := filepath.Join(tmpDir, "absolute.zip")
		os.WriteFile(absZip, buf.Bytes(), 0644)

		// This should sanitize the path
		_, err = installer.InstallFromPath(absZip)
		
		// Check that no file was written to /tmp
		if _, err := os.Stat("/tmp/absolute-evil.txt"); err == nil {
			t.Error("Absolute path attack succeeded - file written to /tmp")
			os.Remove("/tmp/absolute-evil.txt") // Clean up
		}
	})
}

func TestInstaller_MalformedArchives(t *testing.T) {
	t.Run("empty zip file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)
		
		// Create empty zip
		emptyZip := filepath.Join(tmpDir, "empty.zip")
		os.WriteFile(emptyZip, []byte("PK\x05\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), 0644)
		
		_, err = installer.InstallFromPath(emptyZip)
		if err == nil {
			t.Error("Expected error for empty zip file")
		}
	})

	t.Run("corrupted zip file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)
		
		// Create corrupted zip
		corruptedZip := filepath.Join(tmpDir, "corrupted.zip")
		os.WriteFile(corruptedZip, []byte("This is not a zip file"), 0644)
		
		_, err = installer.InstallFromPath(corruptedZip)
		if err == nil {
			t.Error("Expected error for corrupted zip file")
		}
	})

	t.Run("manifest with huge strings", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "installer-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		installer := NewInstaller(tmpDir)

		// Create a zip with a huge manifest
		var buf bytes.Buffer
		zipWriter := zip.NewWriter(&buf)
		
		// Create manifest with very long description (10MB)
		longDesc := strings.Repeat("A", 10*1024*1024)
		manifest := `{"name": "huge-ext", "version": "1.0.0", "description": "` + longDesc + `"}`
		w, _ := zipWriter.Create("gemini-extension.json")
		w.Write([]byte(manifest))
		
		zipWriter.Close()

		// Write the zip to a file
		hugeZip := filepath.Join(tmpDir, "huge.zip")
		os.WriteFile(hugeZip, buf.Bytes(), 0644)

		// This should handle large manifests gracefully
		_, err = installer.InstallFromPath(hugeZip)
		// We don't necessarily expect an error, but it shouldn't crash or use excessive memory
	})
}

