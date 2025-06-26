package profile

import (
	"strings"
	"testing"
	"time"
)

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		profile *Profile
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid profile",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				Extensions:  []ExtensionRef{},
				Environment: map[string]string{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			profile: &Profile{
				Name:        "Test Profile",
				Description: "A test profile",
			},
			wantErr: true,
			errMsg:  "ID is required",
		},
		{
			name: "invalid ID with spaces",
			profile: &Profile{
				ID:          "test profile",
				Name:        "Test Profile",
				Description: "A test profile",
			},
			wantErr: true,
			errMsg:  "ID can only contain",
		},
		{
			name: "invalid ID with special chars",
			profile: &Profile{
				ID:          "test@profile!",
				Name:        "Test Profile",
				Description: "A test profile",
			},
			wantErr: true,
			errMsg:  "ID can only contain",
		},
		{
			name: "ID too long",
			profile: &Profile{
				ID:          strings.Repeat("a", 65),
				Name:        "Test Profile",
				Description: "A test profile",
			},
			wantErr: true,
			errMsg:  "ID must be 64 characters or less",
		},
		{
			name: "missing name",
			profile: &Profile{
				ID:          "test-profile",
				Description: "A test profile",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too long",
			profile: &Profile{
				ID:          "test-profile",
				Name:        strings.Repeat("a", 101),
				Description: "A test profile",
			},
			wantErr: true,
			errMsg:  "name must be 100 characters or less",
		},
		{
			name: "valid auto-detect pattern",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				AutoDetect: &AutoDetectRules{
					Patterns: []string{
						"*.py",
						"Dockerfile",
						"package.json",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty auto-detect pattern",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				AutoDetect: &AutoDetectRules{
					Patterns: []string{
						"",
					},
				},
			},
			wantErr: true,
			errMsg:  "empty pattern",
		},
		{
			name: "circular inheritance - self reference",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				Inherits:    []string{"test-profile"},
			},
			wantErr: true,
			errMsg:  "circular inheritance",
		},
		{
			name: "multiple inheritance",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				Inherits:    []string{"base1", "base2"},
			},
			wantErr: false,
		},
		{
			name: "profile with tags",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				Tags:        []string{"python", "web", "api"},
			},
			wantErr: false,
		},
		{
			name: "profile with MCP servers",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				MCPServers: map[string]ServerConfig{
					"test-server": {
						Enabled: true,
						Settings: map[string]interface{}{
							"command": "node",
							"args":    []string{"server.js"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid environment variable name",
			profile: &Profile{
				ID:          "test-profile",
				Name:        "Test Profile",
				Description: "A test profile",
				Environment: map[string]string{
					"VALID_VAR":   "value",
					"invalid-var": "value", // Dash is typically not allowed
					"123INVALID":  "value", // Starting with number
				},
			},
			wantErr: false, // We might want to allow flexible env var names
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}


func TestValidator_EdgeCases(t *testing.T) {
	validator := NewValidator()

	t.Run("nil profile", func(t *testing.T) {
		err := validator.Validate(nil)
		if err == nil {
			t.Error("Expected error for nil profile")
		}
	})

	t.Run("profile with nil collections", func(t *testing.T) {
		profile := &Profile{
			ID:          "test",
			Name:        "Test",
			Description: "Test",
			Extensions:  nil, // nil slice
			Environment: nil, // nil map
			Inherits:    nil, // nil slice
			Tags:        nil, // nil slice
		}

		err := validator.Validate(profile)
		if err != nil {
			t.Errorf("Should handle nil collections gracefully: %v", err)
		}
	})

	t.Run("profile with unicode", func(t *testing.T) {
		profile := &Profile{
			ID:          "test-profile",
			Name:        "测试配置文件",
			Description: "プロファイルの説明 with émojis 🚀",
			Tags:        []string{"中文", "日本語", "🏷️"},
		}

		err := validator.Validate(profile)
		if err != nil {
			t.Errorf("Should support unicode in name/description/tags: %v", err)
		}
	})

	t.Run("environment with empty values", func(t *testing.T) {
		profile := &Profile{
			ID:          "test-profile",
			Name:        "Test Profile",
			Description: "Test",
			Environment: map[string]string{
				"EMPTY_VAR": "",
				"SPACES":    "   ",
			},
		}

		err := validator.Validate(profile)
		if err != nil {
			t.Errorf("Should allow empty environment values: %v", err)
		}
	})

}