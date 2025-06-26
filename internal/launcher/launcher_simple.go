package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// SimpleLauncher handles launching Gemini CLI with profiles
type SimpleLauncher struct {
	profileManager   *profile.Manager
	extensionManager *extension.Manager
	geminiPath       string
	homeDir          string
}

// NewSimpleLauncher creates a new launcher instance
func NewSimpleLauncher(pm *profile.Manager, em *extension.Manager, geminiPath string) *SimpleLauncher {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "."
	}

	return &SimpleLauncher{
		profileManager:   pm,
		extensionManager: em,
		geminiPath:       geminiPath,
		homeDir:          homeDir,
	}
}

// Launch executes Gemini CLI with the current profile
func (l *SimpleLauncher) Launch(profile *profile.Profile, extensions []*extension.Extension) error {
	if profile == nil {
		return fmt.Errorf("profile cannot be nil")
	}
	
	// Open debug log file
	debugLog, err := os.OpenFile("/tmp/gemini-cli-manager-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer debugLog.Close()
		fmt.Fprintf(debugLog, "\n=== Launch attempt at %s ===\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(debugLog, "Initial geminiPath: %s\n", l.geminiPath)
	}
	
	// Prepare the environment by setting up extension symlinks
	// This is all we need! Each extension has its own gemini-extension.json
	// with mcpServers configuration that Gemini CLI will automatically load
	envPreparer := NewEnvironmentPreparer()
	if err := envPreparer.PrepareExtensions(extensions); err != nil {
		if debugLog != nil {
			fmt.Fprintf(debugLog, "ERROR: Failed to prepare extensions: %v\n", err)
		}
		return fmt.Errorf("preparing extensions: %w", err)
	}
	
	if debugLog != nil {
		managerPath, geminiPath := envPreparer.GetManagedExtensionPaths()
		fmt.Fprintf(debugLog, "Extension paths:\n")
		fmt.Fprintf(debugLog, "  Manager: %s\n", managerPath)
		fmt.Fprintf(debugLog, "  Gemini:  %s\n", geminiPath)
		fmt.Fprintf(debugLog, "Prepared %d enabled extensions\n", len(extensions))
		fmt.Fprintf(debugLog, "Extensions contain their own MCP server configs in gemini-extension.json\n")
	}
	
	// Find the full path to the gemini binary
	geminiPath := l.geminiPath
	
	// If it's not an absolute path, search in PATH
	if !strings.HasPrefix(geminiPath, "/") {
		if debugLog != nil {
			fmt.Fprintf(debugLog, "Searching for gemini in PATH...\n")
		}
		if fullPath, err := exec.LookPath(geminiPath); err == nil {
			geminiPath = fullPath
			if debugLog != nil {
				fmt.Fprintf(debugLog, "Found gemini at: %s\n", fullPath)
			}
		} else {
			if debugLog != nil {
				fmt.Fprintf(debugLog, "ERROR: gemini not found in PATH: %v\n", err)
			}
			return fmt.Errorf("gemini binary not found in PATH: %w", err)
		}
	}
	
	// Verify the binary exists and is executable
	if info, err := os.Stat(geminiPath); err != nil {
		if debugLog != nil {
			fmt.Fprintf(debugLog, "ERROR: stat failed for %s: %v\n", geminiPath, err)
		}
		return fmt.Errorf("gemini binary not found at %s: %w", geminiPath, err)
	} else if debugLog != nil {
		fmt.Fprintf(debugLog, "Binary exists at %s, mode: %v\n", geminiPath, info.Mode())
	}
	
	// Build command
	args := []string{geminiPath}
	
	// Add any profile-specific arguments
	// For now, gemini CLI discovers extensions from ~/.gemini/extensions/
	// and loads enabled ones automatically
	
	// Set up environment
	env := os.Environ()
	
	// Add profile environment variables
	if profile != nil && profile.Environment != nil {
		envMap := make(map[string]string)
		
		// Parse existing environment
		for _, e := range env {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}
		
		// Add profile environment
		for k, v := range profile.Environment {
			envMap[k] = v
		}
		
		// Add Gemini-specific variables
		envMap["GEMINI_PROFILE"] = profile.Name
		envMap["GEMINI_PROFILE_ID"] = profile.ID
		
		// Convert back to slice
		env = make([]string, 0, len(envMap))
		for k, v := range envMap {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	
	// Debug: Log what we're trying to exec
	if debugLog != nil {
		fmt.Fprintf(debugLog, "About to exec:\n")
		fmt.Fprintf(debugLog, "  Path: %s\n", geminiPath)
		fmt.Fprintf(debugLog, "  Args: %v\n", args)
		fmt.Fprintf(debugLog, "  Profile: %s\n", profile.Name)
		fmt.Fprintf(debugLog, "  Env vars added: GEMINI_PROFILE=%s, GEMINI_PROFILE_ID=%s\n", profile.Name, profile.ID)
	}
	
	// Use syscall.Exec to replace our process with Gemini CLI
	// This is the cleanest way to hand over the terminal
	err = syscall.Exec(geminiPath, args, env)
	
	// If we reach here, exec failed
	if debugLog != nil {
		fmt.Fprintf(debugLog, "ERROR: syscall.Exec failed: %v\n", err)
	}
	return fmt.Errorf("syscall.Exec failed for %s: %w", geminiPath, err)
}

// CreateLaunchScript generates a standalone launch script
func (l *SimpleLauncher) CreateLaunchScript(profile *profile.Profile, scriptPath string) error {
	if profile == nil {
		return fmt.Errorf("profile cannot be nil")
	}
	
	script := &strings.Builder{}
	
	// Script header
	fmt.Fprintf(script, "#!/bin/bash\n")
	fmt.Fprintf(script, "# Generated by Gemini CLI Manager\n")
	fmt.Fprintf(script, "# Profile: %s\n", profile.Name)
	fmt.Fprintf(script, "# Generated: %s\n\n", time.Now().Format(time.RFC3339))
	
	// Environment variables
	if profile != nil && len(profile.Environment) > 0 {
		fmt.Fprintf(script, "# Profile Environment\n")
		for k, v := range profile.Environment {
			// Escape single quotes by replacing ' with '\''
			escapedValue := strings.ReplaceAll(v, "'", "'\"'\"'")
			fmt.Fprintf(script, "export %s='%s'\n", k, escapedValue)
		}
		fmt.Fprintf(script, "export GEMINI_PROFILE=\"%s\"\n", profile.Name)
		fmt.Fprintf(script, "\n")
	}
	
	// Note about extensions
	fmt.Fprintf(script, "# Extensions are managed in ~/.gemini/extensions/\n")
	fmt.Fprintf(script, "# Enable/disable them using the Gemini CLI Manager\n\n")
	
	// Launch Gemini
	fmt.Fprintf(script, "# Launch Gemini CLI\n")
	fmt.Fprintf(script, "exec %s \"$@\"\n", l.geminiPath)
	
	// Write script
	if err := os.WriteFile(scriptPath, []byte(script.String()), 0755); err != nil {
		return fmt.Errorf("writing script: %w", err)
	}
	
	return nil
}