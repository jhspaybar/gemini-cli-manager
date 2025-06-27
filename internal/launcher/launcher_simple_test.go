package launcher

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	profilepkg "github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

func TestSimpleLauncher_Launch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "launcher-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock managers (they won't be used in Launch)
	launcher := NewSimpleLauncher(nil, nil, "")

	t.Run("launch with nil profile", func(t *testing.T) {
		err := launcher.Launch(nil, []*extension.Extension{})
		if err == nil {
			t.Error("Expected error for nil profile")
		}
	})

	t.Run("launch with missing gemini", func(t *testing.T) {
		// Set an invalid path
		launcher.geminiPath = "/nonexistent/gemini"

		profile := &profilepkg.Profile{
			ID:   "test",
			Name: "Test Profile",
		}

		err := launcher.Launch(profile, []*extension.Extension{})
		if err == nil {
			t.Error("Expected error when gemini not found")
		}
	})
}

func TestSimpleLauncher_CreateLaunchScript(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "launcher-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	launcher := NewSimpleLauncher(nil, nil, "gemini")

	profile := &profilepkg.Profile{
		ID:   "test-profile",
		Name: "Test Profile",
		Environment: map[string]string{
			"VAR1": "value1",
			"VAR2": "value with spaces",
			"VAR3": "value'with'quotes",
		},
	}

	scriptPath := filepath.Join(tmpDir, "launch.sh")

	t.Run("create basic script", func(t *testing.T) {
		err := launcher.CreateLaunchScript(profile, scriptPath)
		if err != nil {
			t.Fatalf("CreateLaunchScript() error = %v", err)
		}

		// Verify script was created
		info, err := os.Stat(scriptPath)
		if err != nil {
			t.Fatalf("Script not created: %v", err)
		}

		// Verify it's executable
		if runtime.GOOS != "windows" && info.Mode()&0111 == 0 {
			t.Error("Script is not executable")
		}

		// Read and verify content
		content, err := os.ReadFile(scriptPath)
		if err != nil {
			t.Fatalf("Failed to read script: %v", err)
		}

		scriptContent := string(content)

		// Verify shebang
		if !strings.HasPrefix(scriptContent, "#!/bin/bash") {
			t.Error("Script missing shebang")
		}

		// Verify environment variables
		if !strings.Contains(scriptContent, "export VAR1='value1'") {
			t.Error("VAR1 not properly exported")
		}
		if !strings.Contains(scriptContent, "export VAR2='value with spaces'") {
			t.Error("VAR2 with spaces not properly quoted")
		}
		if !strings.Contains(scriptContent, "export VAR3='value'\"'\"'with'\"'\"'quotes'") {
			t.Error("VAR3 with quotes not properly escaped")
		}

		// Verify exec command
		if !strings.Contains(scriptContent, "exec gemini") {
			t.Error("Script missing exec gemini command")
		}
	})

	t.Run("create script with special characters", func(t *testing.T) {
		specialProfile := &profilepkg.Profile{
			ID:   "test",
			Name: "Test",
			Environment: map[string]string{
				"SPECIAL": `value with $VAR and $(command) and ` + "`backticks`",
				"NEWLINE": "value\nwith\nnewlines",
				"TAB":     "value\twith\ttabs",
			},
		}

		scriptPath := filepath.Join(tmpDir, "special.sh")
		err := launcher.CreateLaunchScript(specialProfile, scriptPath)
		if err != nil {
			t.Fatalf("CreateLaunchScript() error = %v", err)
		}

		// Verify script can be parsed (on Unix systems)
		if runtime.GOOS != "windows" {
			cmd := exec.Command("bash", "-n", scriptPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Script has syntax errors: %v\nOutput: %s", err, output)
			}
		}
	})

	t.Run("nil profile", func(t *testing.T) {
		err := launcher.CreateLaunchScript(nil, scriptPath)
		if err == nil {
			t.Error("Expected error for nil profile")
		}
	})
}
