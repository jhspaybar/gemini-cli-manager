package cli

import (
	"fmt"
	"strings"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// ProfileQuickSwitchModal represents a quick profile switcher
type ProfileQuickSwitchModal struct {
	profiles         []*profile.Profile
	filteredProfiles []*profile.Profile
	currentProfileID string
	cursor           int
	searchInput      textinput.Model
	width            int
	height           int

	// Callbacks
	onSelect func(p *profile.Profile) tea.Cmd
	onCancel func() tea.Cmd
}

// NewProfileQuickSwitchModal creates a new quick switch modal
func NewProfileQuickSwitchModal(profiles []*profile.Profile, currentID string) ProfileQuickSwitchModal {
	ti := textinput.New()
	ti.Placeholder = "Type to filter profiles..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40
	ti.Prompt = "🔍 "

	modal := ProfileQuickSwitchModal{
		profiles:         profiles,
		filteredProfiles: profiles,
		currentProfileID: currentID,
		searchInput:      ti,
	}

	// Find current profile position
	for i, p := range profiles {
		if p.ID == currentID {
			modal.cursor = i
			break
		}
	}

	return modal
}

// Init initializes the modal
func (m ProfileQuickSwitchModal) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles modal updates
func (m ProfileQuickSwitchModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			if m.onCancel != nil {
				return m, m.onCancel()
			}
			return m, func() tea.Msg { return closeModalMsg{} }

		case "enter":
			if len(m.filteredProfiles) > 0 && m.cursor < len(m.filteredProfiles) {
				selected := m.filteredProfiles[m.cursor]
				if m.onSelect != nil {
					return m, m.onSelect(selected)
				}
			}
			return m, nil

		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "ctrl+n":
			if m.cursor < len(m.filteredProfiles)-1 {
				m.cursor++
			}

		case "home":
			m.cursor = 0

		case "end":
			if len(m.filteredProfiles) > 0 {
				m.cursor = len(m.filteredProfiles) - 1
			}

		default:
			// Update search input
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)

			// Filter profiles
			query := m.searchInput.Value()
			m.filteredProfiles = filterProfilesForQuickSwitch(m.profiles, query)

			// Reset cursor if out of bounds
			if m.cursor >= len(m.filteredProfiles) {
				m.cursor = 0
			}

			return m, cmd
		}
	}

	return m, nil
}

// View renders the modal
func (m ProfileQuickSwitchModal) View() string {
	// Modal dimensions
	modalWidth := min(50, m.width-4)
	modalHeight := min(20, m.height-4)

	// Create flexbox for modal
	fb := flexbox.New(modalWidth, modalHeight)

	// Title row
	titleRow := fb.NewRow()
	titleCell := flexbox.NewCell(1, 1)
	titleCell.SetContent(h2Style.Render("Switch Profile"))
	titleRow.AddCells(titleCell)

	// Search row
	searchRow := fb.NewRow()
	searchCell := flexbox.NewCell(1, 1)
	searchCell.SetContent(m.searchInput.View())
	searchRow.AddCells(searchCell)

	// Profile list row (takes most space)
	listRow := fb.NewRow()
	listCell := flexbox.NewCell(1, 6) // 6:1 ratio
	listCell.SetContent(m.renderProfileList())
	listRow.AddCells(listCell)

	// Help row
	helpRow := fb.NewRow()
	helpCell := flexbox.NewCell(1, 1)
	helpCell.SetContent(keyDescStyle.Render("Enter: Select • Esc: Cancel"))
	helpRow.AddCells(helpCell)

	// Add all rows
	fb.AddRows([]*flexbox.Row{titleRow, searchRow, listRow, helpRow})

	// Render with border and center
	modalContent := fb.Render()
	styledModal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2).
		Render(modalContent)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styledModal,
	)
}

// renderProfileList renders the scrollable profile list
func (m ProfileQuickSwitchModal) renderProfileList() string {
	if len(m.filteredProfiles) == 0 {
		return textMutedStyle.Render("No profiles match your search")
	}

	var content strings.Builder

	// Show up to 10 profiles
	visibleCount := len(m.filteredProfiles)
	if visibleCount > 10 {
		visibleCount = 10
	}

	// Calculate visible range
	start := 0
	if m.cursor >= visibleCount {
		start = m.cursor - visibleCount + 1
	}
	end := start + visibleCount
	if end > len(m.filteredProfiles) {
		end = len(m.filteredProfiles)
		start = end - visibleCount
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		p := m.filteredProfiles[i]

		prefix := "  "
		style := textStyle

		if i == m.cursor {
			prefix = "▶ "
			style = style.Bold(true).Foreground(colorAccent)
		}

		// Current profile indicator
		indicator := ""
		if p.ID == m.currentProfileID {
			indicator = " ●"
			style = style.Foreground(colorSuccess)
		}

		line := fmt.Sprintf("%s%s%s", prefix, p.Name, indicator)
		content.WriteString(style.Render(line))

		// Show description for selected item
		if i == m.cursor && p.Description != "" {
			content.WriteString("\n")
			content.WriteString(textDimStyle.Render("  " + p.Description))
		}

		if i < end-1 {
			content.WriteString("\n")
		}
	}

	// Scroll indicator
	if len(m.filteredProfiles) > visibleCount {
		content.WriteString("\n")
		content.WriteString(textMutedStyle.Render(fmt.Sprintf("  (%d/%d profiles)", m.cursor+1, len(m.filteredProfiles))))
	}

	return content.String()
}

// SetSize updates modal dimensions
func (m *ProfileQuickSwitchModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetCallbacks sets modal callbacks
func (m *ProfileQuickSwitchModal) SetCallbacks(onSelect func(*profile.Profile) tea.Cmd, onCancel func() tea.Cmd) {
	m.onSelect = onSelect
	m.onCancel = onCancel
}

// filterProfilesForQuickSwitch filters profiles for quick switch
func filterProfilesForQuickSwitch(profiles []*profile.Profile, query string) []*profile.Profile {
	if query == "" {
		return profiles
	}

	query = strings.ToLower(query)
	filtered := make([]*profile.Profile, 0)

	// First pass: exact prefix matches
	for _, p := range profiles {
		if strings.HasPrefix(strings.ToLower(p.Name), query) {
			filtered = append(filtered, p)
		}
	}

	// Second pass: contains matches (if not already added)
	seen := make(map[string]bool)
	for _, p := range filtered {
		seen[p.ID] = true
	}

	for _, p := range profiles {
		if !seen[p.ID] && strings.Contains(strings.ToLower(p.Name), query) {
			filtered = append(filtered, p)
		}
	}

	return filtered
}
