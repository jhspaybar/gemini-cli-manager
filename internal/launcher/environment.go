package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/gemini-cli/manager/internal/extension"
)

// EnvironmentPreparer handles setting up the Gemini environment before launch
type EnvironmentPreparer struct {
	managerExtDir string // Our extension directory (~/.gemini-cli-manager/extensions)
	geminiExtDir  string // Gemini's extension directory (~/.gemini/extensions)
}

// NewEnvironmentPreparer creates a new environment preparer
func NewEnvironmentPreparer() *EnvironmentPreparer {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "."
	}
	
	return &EnvironmentPreparer{
		managerExtDir: filepath.Join(homeDir, ".gemini-cli-manager", "extensions"),
		geminiExtDir:  filepath.Join(homeDir, ".gemini", "extensions"),
	}
}

// PrepareExtensions sets up extensions for Gemini to use
// This creates symlinks from Gemini's extension directory to our managed extensions
func (ep *EnvironmentPreparer) PrepareExtensions(extensions []*extension.Extension) error {
	// Create Gemini's extension directory if it doesn't exist
	if err := os.MkdirAll(ep.geminiExtDir, 0755); err != nil {
		return fmt.Errorf("creating Gemini extension directory: %w", err)
	}
	
	// Clean up any existing symlinks we created previously
	if err := ep.cleanupSymlinks(); err != nil {
		return fmt.Errorf("cleaning up old symlinks: %w", err)
	}
	
	// Create symlinks for all provided extensions (already filtered by profile)
	for _, ext := range extensions {
		srcPath := filepath.Join(ep.managerExtDir, ext.ID)
		dstPath := filepath.Join(ep.geminiExtDir, ext.ID)
		
		// Check if source exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("extension source not found: %s", srcPath)
		}
		
		// Create symlink
		if err := os.Symlink(srcPath, dstPath); err != nil {
			// If it already exists and is not our symlink, skip it
			if os.IsExist(err) {
				continue
			}
			return fmt.Errorf("creating symlink for %s: %w", ext.ID, err)
		}
	}
	
	return nil
}

// cleanupSymlinks removes any symlinks we created in Gemini's extension directory
func (ep *EnvironmentPreparer) cleanupSymlinks() error {
	entries, err := os.ReadDir(ep.geminiExtDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, nothing to clean
		}
		return err
	}
	
	for _, entry := range entries {
		path := filepath.Join(ep.geminiExtDir, entry.Name())
		
		// Check if it's a symlink
		info, err := os.Lstat(path)
		if err != nil {
			continue
		}
		
		if info.Mode()&os.ModeSymlink != 0 {
			// Check if it points to our manager directory
			target, err := os.Readlink(path)
			if err != nil {
				continue
			}
			
			// If it points to our managed extensions, remove it
			if strings.HasPrefix(target, ep.managerExtDir) {
				os.Remove(path)
			}
		}
	}
	
	return nil
}

// GetManagedExtensionPaths returns the paths where our managed extensions are stored
func (ep *EnvironmentPreparer) GetManagedExtensionPaths() (managerPath, geminiPath string) {
	return ep.managerExtDir, ep.geminiExtDir
}