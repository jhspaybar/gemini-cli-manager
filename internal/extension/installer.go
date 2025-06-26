package extension

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// InstallProgress represents the progress of an installation
type InstallProgress struct {
	Stage   string
	Message string
	Percent int
}

// Installer handles extension installation from various sources
type Installer struct {
	extensionsDir string
	progressChan  chan InstallProgress
}

// NewInstaller creates a new extension installer
func NewInstaller(extensionsDir string) *Installer {
	return &Installer{
		extensionsDir: extensionsDir,
	}
}

// SetProgressChannel sets the channel for progress updates
func (i *Installer) SetProgressChannel(ch chan InstallProgress) {
	i.progressChan = ch
}

// sendProgress sends a progress update if channel is set
func (i *Installer) sendProgress(stage, message string, percent int) {
	if i.progressChan != nil {
		select {
		case i.progressChan <- InstallProgress{Stage: stage, Message: message, Percent: percent}:
		default:
			// Don't block if channel is full
		}
	}
}

// Install installs an extension from a source (URL or local path)
func (i *Installer) Install(source string, isPath bool) (*Extension, error) {
	if isPath {
		return i.InstallFromPath(source)
	}
	return i.InstallFromURL(source)
}

// InstallFromURL installs an extension from a remote URL
func (i *Installer) InstallFromURL(url string) (*Extension, error) {
	i.sendProgress("download", "Starting download...", 0)
	
	// Handle different URL types
	if strings.Contains(url, "github.com") {
		return i.installFromGitHub(url)
	}
	
	// Generic URL download
	return i.installFromDirectURL(url)
}

// InstallFromPath installs an extension from a local path
func (i *Installer) InstallFromPath(path string) (*Extension, error) {
	i.sendProgress("validate", "Validating extension...", 0)
	
	// Expand home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("expanding home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	}
	
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}
	
	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("accessing path: %w", err)
	}
	
	// If it's a directory, validate and copy
	if info.IsDir() {
		return i.installFromDirectory(absPath)
	}
	
	// If it's a file, check if it's an archive
	if strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tgz") {
		return i.installFromArchive(absPath)
	}
	
	return nil, fmt.Errorf("unsupported file type")
}

// installFromGitHub handles GitHub URL installation
func (i *Installer) installFromGitHub(url string) (*Extension, error) {
	// Parse GitHub URL
	// Supports:
	// - https://github.com/user/repo
	// - https://github.com/user/repo.git
	// - git@github.com:user/repo.git
	
	var repoURL string
	if strings.HasPrefix(url, "git@github.com:") {
		// Convert SSH URL to HTTPS
		repoURL = strings.Replace(url, "git@github.com:", "https://github.com/", 1)
		repoURL = strings.TrimSuffix(repoURL, ".git")
	} else {
		repoURL = strings.TrimSuffix(url, ".git")
	}
	
	// Extract user and repo
	parts := strings.Split(strings.TrimPrefix(repoURL, "https://github.com/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GitHub URL format")
	}
	
	user := parts[0]
	repo := parts[1]
	
	// Try to download the latest release first
	i.sendProgress("download", "Checking for latest release...", 10)
	ext, err := i.downloadGitHubRelease(user, repo)
	if err == nil {
		return ext, nil
	}
	
	// Fall back to cloning the repository
	i.sendProgress("download", "Cloning repository...", 20)
	return i.cloneGitRepository(url)
}

// downloadGitHubRelease downloads the latest release from GitHub
func (i *Installer) downloadGitHubRelease(user, repo string) (*Extension, error) {
	// GitHub API URL for latest release
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", user, repo)
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("fetching release info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no releases found")
	}
	
	// Parse release info
	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parsing release info: %w", err)
	}
	
	// Look for a suitable asset (zip or tar.gz)
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".zip") || strings.HasSuffix(asset.Name, ".tar.gz") {
			return i.downloadAndInstallAsset(asset.BrowserDownloadURL, asset.Name)
		}
	}
	
	return nil, fmt.Errorf("no suitable release assets found")
}

// downloadAndInstallAsset downloads and installs a release asset
func (i *Installer) downloadAndInstallAsset(url, filename string) (*Extension, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "gemini-ext-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	// Download the asset
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("downloading asset: %w", err)
	}
	defer resp.Body.Close()
	
	// Copy to temp file with progress
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("saving asset: %w", err)
	}
	
	tmpFile.Close()
	
	// Install from the downloaded file
	return i.installFromArchive(tmpFile.Name())
}

// cloneGitRepository clones a git repository
func (i *Installer) cloneGitRepository(url string) (*Extension, error) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "gemini-ext-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Clone the repository
	cmd := exec.Command("git", "clone", "--depth", "1", url, tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cloning repository: %w\nOutput: %s", err, output)
	}
	
	i.sendProgress("install", "Installing extension...", 80)
	
	// Install from the cloned directory
	return i.installFromDirectory(tmpDir)
}

// installFromDirectory installs an extension from a directory
func (i *Installer) installFromDirectory(srcDir string) (*Extension, error) {
	// Load and validate the extension
	manifestPath := filepath.Join(srcDir, "gemini-extension.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}
	
	var ext Extension
	if err := json.Unmarshal(manifestData, &ext); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}
	
	// Validate the extension
	validator := NewValidator()
	if err := validator.Validate(&ext); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	// Create destination directory
	destDir := filepath.Join(i.extensionsDir, ext.ID)
	
	// Check if extension already exists
	if _, err := os.Stat(destDir); err == nil {
		return nil, fmt.Errorf("extension %s already exists", ext.ID)
	}
	
	// Copy the extension
	i.sendProgress("install", "Copying files...", 90)
	if err := i.copyDirectory(srcDir, destDir); err != nil {
		return nil, fmt.Errorf("copying extension: %w", err)
	}
	
	// Set the path
	ext.Path = destDir
	ext.Status = StatusInstalled
	ext.Enabled = false // Start disabled
	
	i.sendProgress("complete", "Installation complete!", 100)
	
	return &ext, nil
}

// installFromArchive installs an extension from an archive file
func (i *Installer) installFromArchive(archivePath string) (*Extension, error) {
	// Create temp directory for extraction
	tmpDir, err := os.MkdirTemp("", "gemini-ext-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	
	i.sendProgress("extract", "Extracting archive...", 50)
	
	// Extract based on file type
	if strings.HasSuffix(archivePath, ".zip") {
		if err := i.extractZip(archivePath, tmpDir); err != nil {
			return nil, fmt.Errorf("extracting zip: %w", err)
		}
	} else if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		if err := i.extractTarGz(archivePath, tmpDir); err != nil {
			return nil, fmt.Errorf("extracting tar.gz: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported archive format")
	}
	
	// Find the extension directory (might be nested)
	extDir, err := i.findExtensionDir(tmpDir)
	if err != nil {
		return nil, err
	}
	
	// Install from the extracted directory
	return i.installFromDirectory(extDir)
}

// installFromDirectURL downloads and installs from a direct URL
func (i *Installer) installFromDirectURL(url string) (*Extension, error) {
	// Determine file type from URL
	filename := filepath.Base(url)
	if !strings.HasSuffix(filename, ".zip") && !strings.HasSuffix(filename, ".tar.gz") {
		return nil, fmt.Errorf("unsupported file type (expected .zip or .tar.gz)")
	}
	
	// Download to temp file
	tmpFile, err := os.CreateTemp("", "gemini-ext-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	i.sendProgress("download", "Downloading extension...", 30)
	
	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("downloading file: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}
	
	// Copy to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("saving file: %w", err)
	}
	
	tmpFile.Close()
	
	// Install from the downloaded file
	return i.installFromArchive(tmpFile.Name())
}

// copyDirectory recursively copies a directory
func (i *Installer) copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)
		
		// Create directories
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		
		// Copy files
		return i.copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func (i *Installer) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	
	// Copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	return os.Chmod(dst, srcInfo.Mode())
}

// extractZip extracts a zip archive
func (i *Installer) extractZip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	
	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)
		
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}
		
		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		
		// Extract file
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		
		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// extractTarGz extracts a tar.gz archive
func (i *Installer) extractTarGz(tarPath, destDir string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()
	
	tr := tar.NewReader(gzr)
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		
		path := filepath.Join(destDir, header.Name)
		
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create directory if needed
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			
			// Extract file
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			
			// Set file permissions
			if err := os.Chmod(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// findExtensionDir finds the directory containing gemini-extension.json
func (i *Installer) findExtensionDir(rootDir string) (string, error) {
	var extensionDir string
	
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			// Check if this directory contains gemini-extension.json
			manifestPath := filepath.Join(path, "gemini-extension.json")
			if _, err := os.Stat(manifestPath); err == nil {
				extensionDir = path
				return filepath.SkipAll // Stop walking
			}
		}
		
		return nil
	})
	
	if err != nil {
		return "", err
	}
	
	if extensionDir == "" {
		return "", fmt.Errorf("no gemini-extension.json found in archive")
	}
	
	return extensionDir, nil
}

// InstallProgressCallback is a function that receives installation progress
type InstallProgressCallback func(stage, message string, percent int)

// InstallWithProgress installs an extension with progress callback
func (i *Installer) InstallWithProgress(source string, isPath bool, callback InstallProgressCallback) (*Extension, error) {
	// Create a progress channel
	progressChan := make(chan InstallProgress, 10)
	i.SetProgressChannel(progressChan)
	defer close(progressChan)
	
	// Start a goroutine to handle progress updates
	done := make(chan bool)
	go func() {
		for {
			select {
			case progress, ok := <-progressChan:
				if !ok {
					return
				}
				callback(progress.Stage, progress.Message, progress.Percent)
			case <-done:
				return
			}
		}
	}()
	
	// Perform installation
	ext, err := i.Install(source, isPath)
	done <- true
	
	return ext, err
}