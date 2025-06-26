package extension

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name     string
		ext      *Extension
		wantErr  bool
		errField string
	}{
		{
			name: "valid extension",
			ext: &Extension{
				ID:          "test-ext",
				Name:        "test-ext",
				Version:     "1.0.0",
				Description: "A test extension",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			ext: &Extension{
				Name:        "test-extension",
				Version:     "1.0.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "id",
		},
		{
			name: "invalid ID with spaces",
			ext: &Extension{
				ID:          "test ext",
				Name:        "test ext",
				Version:     "1.0.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "id",
		},
		{
			name: "invalid ID with special chars",
			ext: &Extension{
				ID:          "test@ext!",
				Name:        "test@ext!",
				Version:     "1.0.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "id",
		},
		{
			name: "missing name",
			ext: &Extension{
				ID:          "test-ext",
				Version:     "1.0.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "name",
		},
		{
			name: "name mismatch with directory",
			ext: &Extension{
				ID:          "test-ext",
				Name:        "different-name",
				Version:     "1.0.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "name",
		},
		{
			name: "missing version",
			ext: &Extension{
				ID:          "test-ext",
				Name:        "test-ext",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "version",
		},
		{
			name: "invalid version format",
			ext: &Extension{
				ID:          "test-ext",
				Name:        "test-ext",
				Version:     "1.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "version",
		},
		{
			name: "invalid version with letters",
			ext: &Extension{
				ID:          "test-ext",
				Name:        "test-ext",
				Version:     "v1.0.0",
				Description: "A test extension",
			},
			wantErr:  true,
			errField: "version",
		},
		{
			name: "extension with valid MCP config",
			ext: &Extension{
				ID:          "mcp-ext",
				Name:        "mcp-ext",
				Version:     "1.0.0",
				Description: "Extension with MCP server",
				MCPServers: map[string]MCPServer{
					"test-server": {
						Command: "node",
						Args:    []string{"server.js"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "extension with invalid MCP config - missing command",
			ext: &Extension{
				ID:          "mcp-ext",
				Name:        "mcp-ext",
				Version:     "1.0.0",
				Description: "Extension with invalid MCP server",
				MCPServers: map[string]MCPServer{
					"test-server": {
						Args: []string{"server.js"},
					},
				},
			},
			wantErr:  true,
			errField: "mcpServers.test-server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errField != "" {
				if !strings.Contains(err.Error(), tt.errField) {
					t.Errorf("Expected error to contain field %q, got %v", tt.errField, err)
				}
			}
		})
	}
}

func TestValidator_ValidateFileStructure(t *testing.T) {
	validator := NewValidator()

	// Create temp directory for testing
	tmpDir, err := os.MkdirTemp("", "validator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "valid structure with manifest",
			setup: func() string {
				extDir := filepath.Join(tmpDir, "valid-ext")
				os.MkdirAll(extDir, 0755)
				manifestPath := filepath.Join(extDir, "gemini-extension.json")
				os.WriteFile(manifestPath, []byte(`{
					"name": "test",
					"displayName": "Test",
					"version": "1.0.0",
					"description": "Test extension"
				}`), 0644)
				return extDir
			},
			wantErr: false,
		},
		{
			name: "missing manifest file",
			setup: func() string {
				extDir := filepath.Join(tmpDir, "no-manifest")
				os.MkdirAll(extDir, 0755)
				return extDir
			},
			wantErr: true,
		},
		{
			name: "empty directory",
			setup: func() string {
				extDir := filepath.Join(tmpDir, "empty")
				os.MkdirAll(extDir, 0755)
				return extDir
			},
			wantErr: true,
		},
		{
			name: "manifest is directory not file",
			setup: func() string {
				extDir := filepath.Join(tmpDir, "bad-manifest")
				os.MkdirAll(extDir, 0755)
				manifestPath := filepath.Join(extDir, "gemini-extension.json")
				os.MkdirAll(manifestPath, 0755)
				return extDir
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			// Create a dummy extension to test file structure
			ext := &Extension{Path: path}
			err := validator.validateFileStructure(ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFileStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_EdgeCases(t *testing.T) {
	validator := NewValidator()

	t.Run("nil extension", func(t *testing.T) {
		err := validator.Validate(nil)
		if err == nil {
			t.Error("Expected error for nil extension")
		}
	})

	t.Run("extension with non-ASCII ID", func(t *testing.T) {
		ext := &Extension{
			ID:          "测试-extension",
			Name:        "test-extension",
			Version:     "1.0.0",
			Description: "A test extension",
		}
		err := validator.Validate(ext)
		if err == nil {
			t.Error("Expected error for non-ASCII ID")
		}
	})

	t.Run("very long ID", func(t *testing.T) {
		ext := &Extension{
			ID:          strings.Repeat("a", 256),
			Name:        "test-extension",
			Version:     "1.0.0",
			Description: "A test extension",
		}
		err := validator.Validate(ext)
		if err == nil {
			t.Error("Expected error for very long ID")
		}
	})

	t.Run("version with pre-release", func(t *testing.T) {
		ext := &Extension{
			ID:          "test-ext",
			Name:        "test-ext",
			Version:     "1.0.0-beta.1",
			Description: "A test extension",
		}
		err := validator.Validate(ext)
		if err != nil {
			t.Errorf("Should accept semantic version with pre-release: %v", err)
		}
	})

	t.Run("version with build metadata", func(t *testing.T) {
		ext := &Extension{
			ID:          "test-ext",
			Name:        "test-ext",
			Version:     "1.0.0+build.123",
			Description: "A test extension",
		}
		err := validator.Validate(ext)
		if err != nil {
			t.Errorf("Should accept semantic version with build metadata: %v", err)
		}
	})
}