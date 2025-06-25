package cli

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
)

// Additional key bindings
var (
	keyTab = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch pane"),
	)
	keySpace = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	)
	keyNew = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	)
	keyDelete = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	)
	keyEdit = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	)
	keyLaunch = key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "launch Gemini"),
	)
)

// Update handles all messages and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle modal updates first
	if m.showingModal && m.modal != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// Update main model size
			m.windowWidth = msg.Width
			m.windowHeight = msg.Height
			// Update modal size
			if modal, ok := m.modal.(SimpleLaunchModal); ok {
				modal.SetSize(msg.Width, msg.Height)
				m.modal = modal
			}
		}
		
		updatedModal, cmd := m.modal.Update(msg)
		m.modal = updatedModal
		
		// Check if modal wants to close
		switch msg := msg.(type) {
		case LaunchCompleteMsg:
			m.showingModal = false
			m.modal = nil
			if msg.Error != nil {
				m.err = msg.Error
			} else {
				// Successfully launched - we could quit or show a success message
				return m, tea.Quit
			}
		}
		
		return m, cmd
	}
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.help.Width = msg.Width
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case tea.KeyMsg:
		// Handle vim-style navigation and arrow keys
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			if m.currentView == ViewHelp {
				m.currentView = ViewExtensions
			} else {
				m.currentView = ViewHelp
			}
			return m, nil
		case "tab":
			// Switch between sidebar and content
			if m.focusedPane == PaneSidebar {
				m.focusedPane = PaneContent
			} else {
				m.focusedPane = PaneSidebar
			}
			return m, nil
		case "left", "h":
			// Move focus to sidebar
			if m.focusedPane == PaneContent {
				m.focusedPane = PaneSidebar
			}
			return m, nil
		case "right", "l":
			// Move focus to content
			if m.focusedPane == PaneSidebar {
				m.focusedPane = PaneContent
			}
			return m, nil
		case "esc":
			// Escape returns to sidebar
			if m.focusedPane == PaneContent {
				m.focusedPane = PaneSidebar
			}
			return m, nil
		case "L":
			// Quick launch
			return m.startLaunch()
		}
		
		// Handle pane-specific navigation
		if m.focusedPane == PaneSidebar {
			return m.updateSidebar(msg)
		} else {
			return m.updateContent(msg)
		}
		
		// Fallback to key matching for any keys we missed
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			if m.currentView == ViewHelp {
				m.currentView = ViewExtensions
			} else {
				m.currentView = ViewHelp
			}
			return m, nil
		case key.Matches(msg, keyTab):
			// Switch between sidebar and content
			if m.focusedPane == PaneSidebar {
				m.focusedPane = PaneContent
			} else {
				m.focusedPane = PaneSidebar
			}
			return m, nil
		case key.Matches(msg, keyLaunch):
			// Quick launch
			return m, tea.Println("Launching Gemini CLI with current profile...")
		}

		// Handle pane-specific navigation
		if m.focusedPane == PaneSidebar {
			return m.updateSidebar(msg)
		} else {
			return m.updateContent(msg)
		}
	}

	return m, nil
}

// updateSidebar handles navigation in the sidebar
func (m Model) updateSidebar(msg tea.KeyMsg) (Model, tea.Cmd) {
	totalItems := len(m.sidebarItems) + 1 // +1 for Launch button
	
	switch msg.String() {
	case "up", "k":
		if m.sidebarCursor > 0 {
			m.sidebarCursor--
		}
	case "down", "j":
		if m.sidebarCursor < totalItems-1 {
			m.sidebarCursor++
		}
	case "enter":
		if m.sidebarCursor < len(m.sidebarItems) {
			// Navigate to view but keep focus on sidebar
			m.currentView = m.sidebarItems[m.sidebarCursor].View
			// Don't auto-switch focus - let user decide with arrow keys
		} else {
			// Launch button
			return m.startLaunch()
		}
	}
	
	return m, nil
}

// updateContent handles navigation in the content area
func (m Model) updateContent(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.currentView {
	case ViewExtensions:
		return m.updateExtensions(msg)
	case ViewProfiles:
		return m.updateProfiles(msg)
	case ViewSettings:
		return m.updateSettings(msg)
	case ViewHelp:
		return m.updateHelp(msg)
	}
	return m, nil
}

func (m Model) updateExtensions(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.extensionsCursor > 0 {
			m.extensionsCursor--
		}
	case "down", "j":
		if m.extensionsCursor < len(m.extensions)-1 {
			m.extensionsCursor++
		}
	case " ":
		// Toggle extension
		if m.extensionsCursor < len(m.extensions) {
			ext := m.extensions[m.extensionsCursor]
			if ext.Enabled {
				m.extensionManager.Disable(ext.ID)
			} else {
				m.extensionManager.Enable(ext.ID)
			}
			// Reload extensions
			m.extensions = m.extensionManager.List()
		}
	case "enter":
		// TODO: Show extension details
		return m, tea.Println("Extension details not yet implemented")
	case "n":
		// TODO: Add new extension
		return m, tea.Println("Add extension not yet implemented")
	case "d":
		// TODO: Delete extension
		return m, tea.Println("Delete extension not yet implemented")
	}
	return m, nil
}

func (m Model) updateProfiles(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.profilesCursor > 0 {
			m.profilesCursor--
		}
	case "down", "j":
		if m.profilesCursor < len(m.profiles)-1 {
			m.profilesCursor++
		}
	case "enter":
		// Activate profile
		if m.profilesCursor < len(m.profiles) {
			profile := m.profiles[m.profilesCursor]
			if err := m.profileManager.SetActive(profile.ID); err == nil {
				m.currentProfile = profile
				return m, tea.Println("Profile activated: " + profile.Name)
			}
		}
	case "n":
		// TODO: Create new profile
		return m, tea.Println("Create profile not yet implemented")
	case "e":
		// TODO: Edit profile
		return m, tea.Println("Edit profile not yet implemented")
	case "d":
		// TODO: Delete profile
		return m, tea.Println("Delete profile not yet implemented")
	}
	return m, nil
}

func (m Model) updateSettings(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "e":
		// TODO: Edit settings
		return m, tea.Println("Edit settings not yet implemented")
	}
	return m, nil
}

func (m Model) updateHelp(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Any key returns to previous view
	m.currentView = ViewExtensions
	return m, nil
}

// startLaunch initiates the launch process
func (m Model) startLaunch() (Model, tea.Cmd) {
	if m.currentProfile == nil {
		m.err = fmt.Errorf("no profile selected")
		return m, nil
	}
	
	// Create launch modal
	modal := NewSimpleLaunchModal(m.currentProfile, m.extensions, m.launcher)
	modal.SetSize(m.windowWidth, m.windowHeight)
	modal.SetCallbacks(
		func() tea.Cmd { return tea.Quit }, // On success, quit
		func() tea.Cmd { // On cancel
			m.showingModal = false
			m.modal = nil
			return nil
		},
	)
	
	m.showingModal = true
	m.modal = modal
	
	// Initialize the modal
	return m, modal.Init()
}