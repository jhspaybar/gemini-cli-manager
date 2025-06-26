package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

func TestSettingsGenerator_GenerateSettings(t *testing.T) {
	tests := []struct {
		name       string
		profile    *profile.Profile
		extensions []*extension.Extension
		want       *GeminiSettings
	}{
		{
			name: "basic profile with extensions",
			profile: &profile.Profile{
				ID:   "test-profile",
				Name: "Test Profile",
				Environment: map[string]string{
					"GEMINI_API_KEY": "test-key-123",
					"LOG_LEVEL":      "debug",
				},
			},
			extensions: []*extension.Extension{
				{
					Name:    "test-ext1",
					Version: "1.0.0",
					MCPServers: map[string]extension.MCPServer{
						"server1": {
							Command: "node",
							Args:    []string{"server.js"},
							Env: map[string]string{
								"API_KEY": "$GEMINI_API_KEY",
								"DEBUG":   "true",
							},
						},
					},
				},
				{
					Name:    "test-ext2",
					Version: "2.0.0",
					MCPServers: map[string]extension.MCPServer{
						"server1": { // Naming conflict
							Command: "python",
							Args:    []string{"-m", "server"},
						},
					},
					ContextFileName: "CUSTOM.md",
				},
			},
			want: &GeminiSettings{
				ContextFileName: "CUSTOM.md",
				MCPServers: map[string]extension.MCPServer{
					"server1": {
						Command: "node",
						Args:    []string{"server.js"},
						Env: map[string]string{
							"API_KEY": "$GEMINI_API_KEY",
							"DEBUG":   "true",
						},
					},
					"test-ext2_server1": {
						Command: "python",
						Args:    []string{"-m", "server"},
						Env:     map[string]string{},
					},
				},
			},
		},
		{
			name: "empty profile with no extensions",
			profile: &profile.Profile{
				ID:   "empty",
				Name: "Empty Profile",
			},
			extensions: []*extension.Extension{},
			want: &GeminiSettings{
				MCPServers: map[string]extension.MCPServer{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewSettingsGenerator(tt.profile, tt.extensions)
			got, err := gen.GenerateSettings()
			if err != nil {
				t.Fatalf("GenerateSettings() error = %v", err)
			}

			// Compare as JSON for easier debugging
			gotJSON, _ := json.MarshalIndent(got, "", "  ")
			wantJSON, _ := json.MarshalIndent(tt.want, "", "  ")

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("GenerateSettings() mismatch:\ngot:\n%s\nwant:\n%s", gotJSON, wantJSON)
			}
		})
	}
}

func TestSettingsGenerator_WriteSettings(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "settings-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	profile := &profile.Profile{
		ID:   "test",
		Name: "Test",
	}

	extensions := []*extension.Extension{
		{
			Name:    "test-ext",
			Version: "1.0.0",
			MCPServers: map[string]extension.MCPServer{
				"test-server": {
					Command: "echo",
					Args:    []string{"hello"},
				},
			},
		},
	}

	gen := NewSettingsGenerator(profile, extensions)
	settingsPath := filepath.Join(tmpDir, "settings.json")

	if err := gen.WriteSettings(settingsPath); err != nil {
		t.Fatalf("WriteSettings() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Error("Settings file was not created")
	}

	// Verify content
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read settings file: %v", err)
	}

	var settings GeminiSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("Failed to parse settings file: %v", err)
	}

	// Check that MCPServers were populated

	if len(settings.MCPServers) != 1 {
		t.Errorf("MCPServers count = %d, want 1", len(settings.MCPServers))
	}
}

func TestSettingsGenerator_NoEnvironmentExpansion(t *testing.T) {
	// Verify that environment variables are NOT expanded
	// Gemini CLI will handle this itself

	profile := &profile.Profile{
		ID:   "test",
		Name: "Test",
		Environment: map[string]string{
			"CUSTOM_KEY": "profile-key",
			"LOG_LEVEL":  "info",
		},
	}

	extensions := []*extension.Extension{
		{
			Name:    "test-ext",
			Version: "1.0.0",
			MCPServers: map[string]extension.MCPServer{
				"server": {
					Command: "$HOME/bin/server",
					CWD:     "${HOME}/workspace",
					Env: map[string]string{
						"API_KEY":    "$TEST_API_KEY",    // From system env
						"CUSTOM_KEY": "$CUSTOM_KEY",      // From profile env
						"LOG_LEVEL":  "${LOG_LEVEL}",     // From profile env
						"STATIC":     "static-value",     // No expansion
					},
				},
			},
		},
	}

	gen := NewSettingsGenerator(profile, extensions)
	settings, err := gen.GenerateSettings()
	if err != nil {
		t.Fatalf("GenerateSettings() error = %v", err)
	}

	server := settings.MCPServers["server"]
	
	// Verify environment variables are NOT expanded
	// Gemini CLI will handle all expansion
	if server.Env["API_KEY"] != "$TEST_API_KEY" {
		t.Errorf("API_KEY = %v, want $TEST_API_KEY (unexpanded)", server.Env["API_KEY"])
	}
	if server.Env["CUSTOM_KEY"] != "$CUSTOM_KEY" {
		t.Errorf("CUSTOM_KEY = %v, want $CUSTOM_KEY (unexpanded)", server.Env["CUSTOM_KEY"])
	}
	if server.Env["LOG_LEVEL"] != "${LOG_LEVEL}" {
		t.Errorf("LOG_LEVEL = %v, want ${LOG_LEVEL} (unexpanded)", server.Env["LOG_LEVEL"])
	}
	if server.Env["STATIC"] != "static-value" {
		t.Errorf("STATIC = %v, want static-value", server.Env["STATIC"])
	}

	// Verify paths are NOT expanded
	if server.Command != "$HOME/bin/server" {
		t.Errorf("Command = %v, want $HOME/bin/server (unexpanded)", server.Command)
	}
	if server.CWD != "${HOME}/workspace" {
		t.Errorf("CWD = %v, want ${HOME}/workspace (unexpanded)", server.CWD)
	}
}

