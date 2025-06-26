package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// SearchBar represents a search input component
type SearchBar struct {
	textInput    textinput.Model
	active       bool
	placeholder  string
	width        int
}

// NewSearchBar creates a new search bar
func NewSearchBar(placeholder string) SearchBar {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 50
	ti.Width = 30
	ti.Prompt = "🔍 "
	
	return SearchBar{
		textInput:   ti,
		placeholder: placeholder,
	}
}

// Init initializes the search bar
func (s SearchBar) Init() tea.Cmd {
	return nil
}

// Update handles search bar updates
func (s SearchBar) Update(msg tea.Msg) (SearchBar, tea.Cmd) {
	if !s.active {
		return s, nil
	}
	
	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

// View renders the search bar
func (s SearchBar) View() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1)
		
	if s.active {
		style = style.BorderForeground(colorBorderFocus)
	}
	
	return style.Render(s.textInput.View())
}

// Focus activates the search bar
func (s *SearchBar) Focus() tea.Cmd {
	s.active = true
	return s.textInput.Focus()
}

// Blur deactivates the search bar
func (s *SearchBar) Blur() {
	s.active = false
	s.textInput.Blur()
}

// IsActive returns whether the search bar is active
func (s SearchBar) IsActive() bool {
	return s.active
}

// Value returns the current search query
func (s SearchBar) Value() string {
	return s.textInput.Value()
}

// SetValue sets the search query
func (s *SearchBar) SetValue(value string) {
	s.textInput.SetValue(value)
}

// Clear clears the search query
func (s *SearchBar) Clear() {
	s.textInput.SetValue("")
}

// SetWidth sets the width of the search bar
func (s *SearchBar) SetWidth(width int) {
	s.width = width
	s.textInput.Width = width - 6 // Account for border and padding
}

// filterExtensions filters extensions based on the search query
func filterExtensions(extensions []*extension.Extension, query string) []*extension.Extension {
	if query == "" {
		return extensions
	}
	
	query = strings.ToLower(query)
	filtered := make([]*extension.Extension, 0)
	
	for _, ext := range extensions {
		// Search in name and description
		if strings.Contains(strings.ToLower(ext.Name), query) ||
			strings.Contains(strings.ToLower(ext.Description), query) {
			filtered = append(filtered, ext)
		}
	}
	
	return filtered
}

// filterProfiles filters profiles based on the search query
func filterProfiles(profiles []*profile.Profile, query string) []*profile.Profile {
	if query == "" {
		return profiles
	}
	
	query = strings.ToLower(query)
	filtered := make([]*profile.Profile, 0)
	
	for _, prof := range profiles {
		// Search in name, description, and tags
		if strings.Contains(strings.ToLower(prof.Name), query) ||
			strings.Contains(strings.ToLower(prof.Description), query) {
			filtered = append(filtered, prof)
			continue
		}
		
		// Search in tags
		for _, tag := range prof.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				filtered = append(filtered, prof)
				break
			}
		}
	}
	
	return filtered
}