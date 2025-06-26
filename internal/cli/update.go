package cli

import (
	"fmt"
	"time"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/gemini-cli/manager/internal/extension"
	"github.com/gemini-cli/manager/internal/profile"
)

// closeModalMsg is sent when a modal wants to close
type closeModalMsg struct{}

// profileSavedMsg is sent when a profile is successfully saved
type profileSavedMsg struct {
	profile *profile.Profile
	isNew   bool
}

// execGeminiMsg signals that we should exec Gemini after quitting
type execGeminiMsg struct {
	profile    *profile.Profile
	extensions []*extension.Extension
}

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
		// Allow Ctrl+C to quit even in modals
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// Update main model size
			m.windowWidth = msg.Width
			m.windowHeight = msg.Height
			// Update modal size
			switch modal := m.modal.(type) {
			case SimpleLaunchModal:
				modal.SetSize(msg.Width, msg.Height)
				m.modal = modal
			case ProfileForm:
				modal.SetSize(msg.Width, msg.Height)
				m.modal = modal
			case ExtensionInstallForm:
				modal.SetSize(msg.Width, msg.Height)
				m.modal = modal
			case ProfileQuickSwitchModal:
				modal.SetSize(msg.Width, msg.Height)
				m.modal = modal
			}
		}
		
		updatedModal, cmd := m.modal.Update(msg)
		m.modal = updatedModal
		
		// Check if modal wants to close
		switch msg := msg.(type) {
		case LaunchCompleteMsg:
			if msg.Error != nil {
				// Keep modal open to show error
				m.err = msg.Error
			} else {
				// Successfully launched - the Gemini CLI is now running
				// Don't quit immediately, let the modal handle it
			}
		case installCompleteMsg:
			if msg.err == nil {
				// Success - close modal and refresh extensions
				m.showingModal = false
				m.modal = nil
				// Reload extensions
				m.extensions = m.extensionManager.List()
			}
			// If error, keep modal open to show error
		case closeModalMsg:
			// Close any open modal
			m.showingModal = false
			m.modal = nil
		case profileSavedMsg:
			// Profile was saved successfully
			m.showingModal = false
			m.modal = nil
			
			// Refresh profiles list
			m.profiles = m.profileManager.List()
			if m.searchBar.Value() != "" {
				m.filteredProfiles = filterProfiles(m.profiles, m.searchBar.Value())
			} else {
				m.filteredProfiles = m.profiles
			}
			
			// If we created a new profile, select it
			if msg.isNew {
				for i, prof := range m.profiles {
					if prof.ID == msg.profile.ID {
						m.profilesCursor = i
						break
					}
				}
			}
		case UIError:
			// Set error
			m.err = msg
		case execGeminiMsg:
			// Store the launch info and quit
			m.shouldExecGemini = true
			m.execProfile = msg.profile
			m.execExtensions = msg.extensions
			m.showingModal = false
			m.modal = nil
			return m, tea.Quit
		}
		
		return m, cmd
	}
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.help.Width = msg.Width
		// Update search bar width
		m.searchBar.SetWidth(msg.Width / 2)
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case tea.KeyMsg:
		// Handle search mode first
		if m.searchActive {
			switch msg.String() {
			case "esc":
				// Exit search
				m.searchActive = false
				m.searchBar.Blur()
				m.searchBar.Clear()
				// Reset filtered lists
				m.filteredExtensions = m.extensions
				m.filteredProfiles = m.profiles
				return m, nil
			case "enter":
				// Apply search and exit search mode
				m.searchActive = false
				m.searchBar.Blur()
				return m, nil
			default:
				// Update search bar
				var cmd tea.Cmd
				m.searchBar, cmd = m.searchBar.Update(msg)
				
				// Apply filters
				query := m.searchBar.Value()
				m.filteredExtensions = filterExtensions(m.extensions, query)
				m.filteredProfiles = filterProfiles(m.profiles, query)
				
				// Reset cursors if they're out of bounds
				if m.extensionsCursor >= len(m.filteredExtensions) {
					m.extensionsCursor = 0
				}
				if m.profilesCursor >= len(m.filteredProfiles) {
					m.profilesCursor = 0
				}
				
				return m, cmd
			}
		}
		
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
		case "/":
			// Activate search
			m.searchActive = true
			return m, m.searchBar.Focus()
		case "ctrl+p":
			// Quick switch profiles
			return m.showProfileQuickSwitch()
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
		if m.extensionsCursor < len(m.filteredExtensions)-1 {
			m.extensionsCursor++
		}
	case " ":
		// Toggle extension
		if m.extensionsCursor < len(m.filteredExtensions) {
			ext := m.filteredExtensions[m.extensionsCursor]
			if ext.Enabled {
				m.extensionManager.Disable(ext.ID)
			} else {
				m.extensionManager.Enable(ext.ID)
			}
			// Reload extensions
			m.extensions = m.extensionManager.List()
			// Reapply filters if search is active
			if m.searchBar.Value() != "" {
				m.filteredExtensions = filterExtensions(m.extensions, m.searchBar.Value())
			} else {
				m.filteredExtensions = m.extensions
			}
		}
	case "enter":
		// TODO: Show extension details
		return m, tea.Println("Extension details not yet implemented")
	case "n":
		// Add new extension
		return m.showExtensionInstallForm()
	case "d":
		// Delete extension
		if m.extensionsCursor < len(m.filteredExtensions) {
			ext := m.filteredExtensions[m.extensionsCursor]
			// Move to trash instead of permanent delete
			if err := m.extensionManager.MoveToTrash(ext.ID); err != nil {
				m.err = NewFileSystemError("delete", ext.Name, err)
			} else {
				// Reload extensions
				m.extensions = m.extensionManager.List()
				// Reapply filters
				if m.searchBar.Value() != "" {
					m.filteredExtensions = filterExtensions(m.extensions, m.searchBar.Value())
				} else {
					m.filteredExtensions = m.extensions
				}
				// Adjust cursor if needed
				if m.extensionsCursor >= len(m.filteredExtensions) && m.extensionsCursor > 0 {
					m.extensionsCursor--
				}
			}
		}
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
		if m.profilesCursor < len(m.filteredProfiles)-1 {
			m.profilesCursor++
		}
	case "enter":
		// Activate profile
		if m.profilesCursor < len(m.filteredProfiles) {
			profile := m.filteredProfiles[m.profilesCursor]
			if err := m.profileManager.SetActive(profile.ID); err == nil {
				m.currentProfile = profile
				return m, tea.Println("Profile activated: " + profile.Name)
			}
		}
	case "n":
		// Create new profile
		return m.showProfileForm(nil, false)
	case "e":
		// Edit profile
		if m.profilesCursor < len(m.filteredProfiles) {
			return m.showProfileForm(m.filteredProfiles[m.profilesCursor], true)
		}
		return m, nil
	case "d":
		// Delete profile
		if m.profilesCursor < len(m.filteredProfiles) {
			prof := m.filteredProfiles[m.profilesCursor]
			// Don't allow deleting default or active profile
			if prof.ID == "default" {
				m.err = NewValidationError("Cannot delete the default profile", "The default profile is protected")
			} else if m.currentProfile != nil && prof.ID == m.currentProfile.ID {
				m.err = NewValidationError("Cannot delete the active profile", "Switch to another profile first")
			} else {
				// Delete the profile
				if err := m.profileManager.Delete(prof.ID); err != nil {
					m.err = WrapError(err, "profile deletion")
				} else {
					// Reload profiles
					m.profiles = m.profileManager.List()
					// Reapply filters
					if m.searchBar.Value() != "" {
						m.filteredProfiles = filterProfiles(m.profiles, m.searchBar.Value())
					} else {
						m.filteredProfiles = m.profiles
					}
					// Adjust cursor if needed
					if m.profilesCursor >= len(m.filteredProfiles) && m.profilesCursor > 0 {
						m.profilesCursor--
					}
				}
			}
		}
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

// showProfileForm shows the profile creation/edit form
func (m Model) showProfileForm(prof *profile.Profile, isEdit bool) (Model, tea.Cmd) {
	// Get list of extension IDs
	extIDs := make([]string, 0, len(m.extensions))
	for _, ext := range m.extensions {
		extIDs = append(extIDs, ext.ID)
	}
	
	// Create form
	form := NewProfileForm(prof, extIDs, isEdit)
	form.SetSize(m.windowWidth, m.windowHeight)
	form.SetCallbacks(
		func(p *profile.Profile) tea.Cmd {
			// Save profile and close modal
			return func() tea.Msg {
				var err error
				if isEdit {
					err = m.profileManager.Save(p)
				} else {
					err = m.profileManager.Create(p)
				}
				
				if err != nil {
					// Return error message
					return UIError{
						Type:    ErrorTypeFileSystem,
						Message: fmt.Sprintf("Failed to save profile"),
						Details: err.Error(),
					}
				}
				
				// Return success message to trigger refresh
				return profileSavedMsg{profile: p, isNew: !isEdit}
			}
		},
		func() tea.Cmd {
			// Cancel - return a command to close modal
			return func() tea.Msg {
				return closeModalMsg{}
			}
		},
	)
	
	m.showingModal = true
	m.modal = form
	m.err = nil // Clear any previous errors
	
	// Initialize the form
	return m, form.Init()
}

// showProfileQuickSwitch shows the profile quick switch modal
func (m Model) showProfileQuickSwitch() (Model, tea.Cmd) {
	currentID := ""
	if m.currentProfile != nil {
		currentID = m.currentProfile.ID
	}
	
	modal := NewProfileQuickSwitchModal(m.profiles, currentID)
	modal.SetSize(m.windowWidth, m.windowHeight)
	modal.SetCallbacks(
		func(p *profile.Profile) tea.Cmd {
			// Switch profile
			if err := m.profileManager.SetActive(p.ID); err == nil {
				m.currentProfile = p
				// Refresh profile list to update active indicator
				m.profiles = m.profileManager.List()
				if m.searchBar.Value() != "" {
					m.filteredProfiles = filterProfiles(m.profiles, m.searchBar.Value())
				} else {
					m.filteredProfiles = m.profiles
				}
			}
			// Close modal
			m.showingModal = false
			m.modal = nil
			return nil
		},
		func() tea.Cmd {
			// Cancel
			m.showingModal = false
			m.modal = nil
			return nil
		},
	)
	
	m.showingModal = true
	m.modal = modal
	
	return m, modal.Init()
}

// showExtensionInstallForm shows the extension installation form
func (m Model) showExtensionInstallForm() (Model, tea.Cmd) {
	// Create form
	form := NewExtensionInstallForm()
	form.SetSize(m.windowWidth, m.windowHeight)
	form.SetCallbacks(
		func(source string, isPath bool) tea.Cmd {
			// Install extension
			return func() tea.Msg {
				// Create a progress channel for updates
				progressChan := make(chan extension.InstallProgress, 10)
				
				// Start progress monitoring in a goroutine
				go func() {
					for progress := range progressChan {
						// Would send progress updates to the UI here
						// For now, just log them
						fmt.Printf("Install progress: %s - %s (%d%%)\n", 
							progress.Stage, progress.Message, progress.Percent)
					}
				}()
				
				// Install the extension
				ext, err := m.extensionManager.InstallWithProgress(source, isPath, 
					func(stage, message string, percent int) {
						select {
						case progressChan <- extension.InstallProgress{
							Stage:   stage,
							Message: message,
							Percent: percent,
						}:
						default:
							// Don't block if channel is full
						}
					})
				
				close(progressChan)
				
				if err != nil {
					return installCompleteMsg{
						err: err,
					}
				}
				
				return installCompleteMsg{
					extension: ext,
					err:       nil,
				}
			}
		},
		func() tea.Cmd {
			// Cancel - return a command to close modal
			return func() tea.Msg {
				return closeModalMsg{}
			}
		},
	)
	
	m.showingModal = true
	m.modal = form
	m.err = nil
	
	// Initialize the form
	return m, form.Init()
}