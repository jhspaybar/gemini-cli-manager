package cli

import (
	"testing"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// TestEnterKeyBehavior tests that Enter key works correctly on different views
func TestEnterKeyBehavior(t *testing.T) {
	tests := []struct {
		name           string
		setupModel     func() Model
		keyMsg         tea.KeyMsg
		checkResult    func(t *testing.T, m Model, cmd tea.Cmd)
	}{
		{
			name: "Enter key on extension shows info",
			setupModel: func() Model {
				m := Model{
					currentView: ViewExtensions,
					filteredExtensions: []*extension.Extension{
						{
							ID:          "test-ext",
							Name:        "Test Extension",
							Version:     "1.0.0",
							Description: "A test extension",
						},
					},
					extensionsCursor: 0,
				}
				return m
			},
			keyMsg: tea.KeyMsg{Type: tea.KeyEnter},
			checkResult: func(t *testing.T, m Model, cmd tea.Cmd) {
				// Check that we switched to detail view
				if m.currentView != ViewExtensionDetail {
					t.Errorf("Expected ViewExtensionDetail, got %v", m.currentView)
				}
				
				// Check that selected extension is set
				if m.selectedExtension == nil {
					t.Error("Expected selectedExtension to be set")
				} else if m.selectedExtension.ID != "test-ext" {
					t.Errorf("Expected extension ID 'test-ext', got %s", m.selectedExtension.ID)
				}
			},
		},
		{
			name: "Space key on extension also shows details",
			setupModel: func() Model {
				m := Model{
					currentView: ViewExtensions,
					filteredExtensions: []*extension.Extension{
						{
							ID:          "test-ext",
							Name:        "Test Extension",
							Version:     "1.0.0",
							Description: "A test extension",
						},
					},
					extensionsCursor: 0,
				}
				return m
			},
			keyMsg: tea.KeyMsg{Type: tea.KeySpace},
			checkResult: func(t *testing.T, m Model, cmd tea.Cmd) {
				// Check that we switched to detail view
				if m.currentView != ViewExtensionDetail {
					t.Errorf("Expected ViewExtensionDetail, got %v", m.currentView)
				}
				
				// Check that selected extension is set
				if m.selectedExtension == nil {
					t.Error("Expected selectedExtension to be set")
				}
			},
		},
		{
			name: "Enter key on profile returns activation command",
			setupModel: func() Model {
				m := Model{
					currentView: ViewProfiles,
					filteredProfiles: []*profile.Profile{
						{ID: "test", Name: "Test Profile"},
					},
					profilesCursor: 0,
				}
				return m
			},
			keyMsg: tea.KeyMsg{Type: tea.KeyEnter},
			checkResult: func(t *testing.T, m Model, cmd tea.Cmd) {
				// Should return a command
				if cmd == nil {
					t.Fatal("Expected command for profile activation, got nil")
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupModel()
			
			var updatedModel Model
			var cmd tea.Cmd
			
			// Update based on current view
			switch m.currentView {
			case ViewExtensions:
				updatedModel, cmd = m.updateExtensions(tt.keyMsg)
			case ViewProfiles:
				updatedModel, cmd = m.updateProfiles(tt.keyMsg)
			}
			
			tt.checkResult(t, updatedModel, cmd)
		})
	}
}

// TestRegressionEnterKeyNotBroken verifies Enter key wasn't broken during refactoring
func TestRegressionEnterKeyNotBroken(t *testing.T) {
	t.Run("Extension Enter shows info not nothing", func(t *testing.T) {
		m := Model{
			currentView: ViewExtensions,
			filteredExtensions: []*extension.Extension{
				{ID: "ext1", Name: "Extension 1", Version: "1.0.0"},
			},
			extensionsCursor: 0,
		}
		
		// Before: Enter key did nothing after async refactor
		// After: Enter key should show extension details
		updated, _ := m.updateExtensions(tea.KeyMsg{Type: tea.KeyEnter})
		
		if updated.currentView != ViewExtensionDetail {
			t.Error("Regression: Enter key on extension does not show details")
		}
	})
	
	t.Run("Profile Enter triggers activation", func(t *testing.T) {
		m := Model{
			currentView: ViewProfiles,
			filteredProfiles: []*profile.Profile{
				{ID: "prof1", Name: "Profile 1"},
			},
			profilesCursor: 0,
		}
		
		// Before: Enter key did nothing after async refactor  
		// After: Enter key should return activation command
		_, cmd := m.updateProfiles(tea.KeyMsg{Type: tea.KeyEnter})
		
		if cmd == nil {
			t.Error("Regression: Enter key on profile does not activate")
		}
	})
}