package cli

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// closeModalMsg is sent when a modal wants to close
type closeModalMsg struct{}

// profileSavedMsg is sent when a profile is successfully saved
type profileSavedMsg struct {
	profile *profile.Profile
	isNew   bool
}

// profileSwitchMsg is sent when switching to a different profile
type profileSwitchMsg struct {
	profile *profile.Profile
}

// execGeminiMsg signals that we should exec Gemini after quitting
type execGeminiMsg struct {
	profile    *profile.Profile
	extensions []*extension.Extension
}

// installStartMsg is sent to start extension installation
type installStartMsg struct {
	source string
	isPath bool
}

// installProgressMsg is sent for installation progress updates
type installProgressMsg struct {
	stage   string
	message string
	percent int
}

// extensionsLoadedMsg is sent when extensions are loaded
type extensionsLoadedMsg struct {
	extensions []*extension.Extension
	err        error
}

// extensionDeletedMsg is sent when an extension is deleted
type extensionDeletedMsg struct {
	extensionID string
	err         error
}

// profilesLoadedMsg is sent when profiles are loaded
type profilesLoadedMsg struct {
	profiles []*profile.Profile
	err      error
}

// profileDeletedMsg is sent when a profile is deleted
type profileDeletedMsg struct {
	profileID string
	err       error
}

// profileActivatedMsg is sent when a profile is activated
type profileActivatedMsg struct {
	profile *profile.Profile
	err     error
}

// extensionSelectedMsg is sent when an extension is selected for details
type extensionSelectedMsg struct {
	extension *extension.Extension
}

// extensionSavedMsg is sent when an extension is saved
type extensionSavedMsg struct {
	extension *extension.Extension
}

// initCompleteMsg is sent when initialization is complete
type initCompleteMsg struct {
	err error
}

// managersInitializedMsg is sent when managers are initialized
type managersInitializedMsg struct {
	err error
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
	LogMessage("Model.Update", msg)
	
	// Handle modal updates first
	if m.showingModal && m.modal != nil {
		LogDebug("Modal is showing, type: %T", m.modal)
		
		// Allow Ctrl+C to quit even in modals
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
			LogDebug("Ctrl+C pressed in modal, quitting")
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
			case *ExtensionEditForm:
				modal.SetSize(msg.Width, msg.Height)
				m.modal = modal
			}
		}
		
		LogDebug("Calling modal.Update with message: %T", msg)
		updatedModal, cmd := m.modal.Update(msg)
		LogDebug("Modal.Update returned, cmd: %v", cmd)
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
		case installStartMsg:
			// Start the installation process
			return m, m.performInstallation(msg.source, msg.isPath)
		case installProgressMsg:
			// Update the modal with progress if it's an ExtensionInstallForm
			if form, ok := m.modal.(ExtensionInstallForm); ok {
				form.progress = fmt.Sprintf("%s: %s (%d%%)", msg.stage, msg.message, msg.percent)
				m.modal = form
			}
			return m, nil
		case installCompleteMsg:
			if msg.err == nil {
				// Success - close modal and refresh extensions
				m.showingModal = false
				m.modal = nil
				// Reload extensions asynchronously
				return m, m.loadExtensionsCmd()
			} else {
				// If error, keep modal open to show error
				m.err = msg.err
			}
		case closeModalMsg:
			// Close any open modal
			m.showingModal = false
			m.modal = nil
		case extensionSavedMsg:
			// Extension was saved successfully
			m.showingModal = false
			m.modal = nil
			// Update the selected extension
			if msg.extension != nil {
				m.selectedExtension = msg.extension
			}
			// Reload extensions list
			return m, m.loadExtensionsCmd()
		case profileSavedMsg:
			// Profile was saved successfully
			m.showingModal = false
			m.modal = nil
			
			// Store the new profile info for cursor positioning
			if msg.isNew {
				m.newProfileID = msg.profile.ID
			}
			
			// Refresh profiles list asynchronously
			return m, m.loadProfilesCmd()
		case profileSwitchMsg:
			// Close modal first
			m.showingModal = false
			m.modal = nil
			// Switch to the selected profile asynchronously
			return m, m.activateProfileCmd(msg.profile)
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
		
		// Return with the command from modal update
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
		
	case extensionsLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.extensions = msg.extensions
			// Update filtered list too
			if m.searchBar.Value() != "" {
				m.filteredExtensions = filterExtensions(m.extensions, m.searchBar.Value())
			} else {
				m.filteredExtensions = m.extensions
			}
		}
		// Check if initialization is complete
		if m.loading && m.checkInitComplete() {
			m.loading = false
		}
		return m, nil
		
	case extensionDeletedMsg:
		if msg.err != nil {
			m.err = NewFileSystemError("delete", msg.extensionID, msg.err)
		} else {
			// Reload extensions after successful deletion
			return m, m.loadExtensionsCmd()
		}
		
	case profilesLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.profiles = msg.profiles
			// Update filtered list too
			if m.searchBar.Value() != "" {
				m.filteredProfiles = filterProfiles(m.profiles, m.searchBar.Value())
			} else {
				m.filteredProfiles = m.profiles
			}
			
			// If we just created a new profile, position cursor on it
			if m.newProfileID != "" {
				for i, prof := range m.profiles {
					if prof.ID == m.newProfileID {
						m.profilesCursor = i
						break
					}
				}
				m.newProfileID = "" // Clear after use
			}
		}
		// Check if initialization is complete
		if m.loading && m.checkInitComplete() {
			m.loading = false
		}
		return m, nil
		
	case profileDeletedMsg:
		if msg.err != nil {
			m.err = WrapError(msg.err, "profile deletion")
		} else {
			// Reload profiles after successful deletion
			// Adjust cursor if needed
			if m.profilesCursor >= len(m.filteredProfiles)-1 && m.profilesCursor > 0 {
				m.profilesCursor--
			}
			return m, m.loadProfilesCmd()
		}
		return m, nil
		
	case profileActivatedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.currentProfile = msg.profile
			// During initialization, check if we're done loading
			if m.loading {
				// Check if initialization is complete
				if m.checkInitComplete() {
					return m, func() tea.Msg {
						return initCompleteMsg{err: nil}
					}
				}
			} else {
				// Normal profile activation - reload profiles to update active indicator
				return m, m.loadProfilesCmd()
			}
		}
		return m, nil
		
	case managersInitializedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.loading = false
			return m, nil
		}
		// Managers initialized, now load data
		return m, m.loadInitialDataCmd()
		
	case initCompleteMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
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
			// Cycle through tabs (skip if in detail view)
			if m.currentView != ViewExtensionDetail {
				switch m.currentView {
				case ViewExtensions:
					m.currentView = ViewProfiles
				case ViewProfiles:
					m.currentView = ViewSettings
				case ViewSettings:
					m.currentView = ViewHelp
				case ViewHelp:
					m.currentView = ViewExtensions
				}
			}
			return m, nil
		case "left", "h":
			// Previous tab (skip if in detail view)
			if m.currentView != ViewExtensionDetail {
				switch m.currentView {
				case ViewExtensions:
					m.currentView = ViewHelp
				case ViewProfiles:
					m.currentView = ViewExtensions
				case ViewSettings:
					m.currentView = ViewProfiles
				case ViewHelp:
					m.currentView = ViewSettings
				}
			}
			return m, nil
		case "right", "l":
			// Next tab (skip if in detail view)
			if m.currentView != ViewExtensionDetail {
				switch m.currentView {
				case ViewExtensions:
					m.currentView = ViewProfiles
				case ViewProfiles:
					m.currentView = ViewSettings
				case ViewSettings:
					m.currentView = ViewHelp
				case ViewHelp:
					m.currentView = ViewExtensions
				}
			}
			return m, nil
		case "esc":
			// Handle escape based on current view
			if m.currentView == ViewExtensionDetail {
				// In detail view, go back to extensions
				m.currentView = ViewExtensions
				m.selectedExtension = nil
			}
			// Otherwise, escape does nothing in main navigation
			return m, nil
		case "L":
			// Quick launch
			return m.startLaunch()
		case "/":
			// Activate search (not in detail views)
			if m.currentView != ViewExtensionDetail {
				m.searchActive = true
				return m, m.searchBar.Focus()
			}
			return m, nil
		case "ctrl+p":
			// Quick switch profiles
			return m.showProfileQuickSwitch()
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
		case key.Matches(msg, keyLaunch):
			// Quick launch
			return m.startLaunch()
		}

		// Handle view-specific navigation
		return m.updateContent(msg)
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
	case ViewExtensionDetail:
		return m.updateExtensionDetail(msg)
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
	case " ", "enter":
		// View extension details
		if m.extensionsCursor < len(m.filteredExtensions) {
			ext := m.filteredExtensions[m.extensionsCursor]
			LogDebug("Selected extension: %s, switching to detail view", ext.Name)
			m.selectedExtension = ext
			m.currentView = ViewExtensionDetail
			m.err = nil // Clear any errors
			return m, nil
		}
	case "n":
		// Add new extension
		return m.showExtensionInstallForm()
	case "d":
		// Delete extension
		if m.extensionsCursor < len(m.filteredExtensions) {
			ext := m.filteredExtensions[m.extensionsCursor]
			// Delete asynchronously
			return m, m.deleteExtensionCmd(ext.ID)
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
			return m, m.activateProfileCmd(profile)
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
				// Delete the profile asynchronously
				// Adjust cursor before deletion
				if m.profilesCursor >= len(m.filteredProfiles)-1 && m.profilesCursor > 0 {
					m.profilesCursor--
				}
				return m, m.deleteProfileCmd(prof.ID)
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

func (m Model) updateExtensionDetail(msg tea.KeyMsg) (Model, tea.Cmd) {
	LogDebug("updateExtensionDetail: key=%s", msg.String())
	
	switch msg.String() {
	case "esc":
		// Go back to extensions list
		LogDebug("ESC pressed, going back to extensions list")
		m.currentView = ViewExtensions
		m.selectedExtension = nil
		return m, nil
	case "d":
		// Delete the extension
		if m.selectedExtension != nil {
			extID := m.selectedExtension.ID
			LogDebug("Delete extension: %s", extID)
			// Go back to list first
			m.currentView = ViewExtensions
			m.selectedExtension = nil
			// Then delete
			return m, m.deleteExtensionCmd(extID)
		}
	case "e":
		// Edit extension
		if m.selectedExtension != nil {
			LogDebug("Edit extension requested: %s", m.selectedExtension.Name)
			return m.showExtensionEditForm(m.selectedExtension)
		}
		LogDebug("No extension selected for editing")
		return m, nil
	}
	return m, nil
}

// startLaunch initiates the launch process
func (m Model) startLaunch() (Model, tea.Cmd) {
	if m.currentProfile == nil {
		m.err = fmt.Errorf("no profile selected")
		return m, nil
	}
	
	// Get only the extensions that are in the current profile
	profileExtensions := m.getProfileExtensions(m.currentProfile)
	
	// Create launch modal
	modal := NewSimpleLaunchModal(m.currentProfile, profileExtensions, m.launcher)
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
			// Switch profile - return profile switch message
			return func() tea.Msg {
				return profileSwitchMsg{profile: p}
			}
		},
		func() tea.Cmd {
			// Cancel - return close modal message
			return func() tea.Msg {
				return closeModalMsg{}
			}
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
			// Return a message to start the installation
			return func() tea.Msg {
				return installStartMsg{
					source: source,
					isPath: isPath,
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

// performInstallation performs the extension installation without goroutines
func (m Model) performInstallation(source string, isPath bool) tea.Cmd {
	return func() tea.Msg {
		// For now, use the regular Install method without progress
		// We'll refactor the extension manager to support non-blocking progress later
		ext, err := m.extensionManager.Install(source, isPath)
		
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
}

// Command functions for async operations

// loadExtensionsCmd loads the list of extensions
func (m Model) loadExtensionsCmd() tea.Cmd {
	return func() tea.Msg {
		extensions := m.extensionManager.List()
		return extensionsLoadedMsg{
			extensions: extensions,
			err:        nil,
		}
	}
}

// deleteExtensionCmd deletes an extension
func (m Model) deleteExtensionCmd(extID string) tea.Cmd {
	return func() tea.Msg {
		err := m.extensionManager.MoveToTrash(extID)
		return extensionDeletedMsg{
			extensionID: extID,
			err:         err,
		}
	}
}

// loadProfilesCmd loads the list of profiles
func (m Model) loadProfilesCmd() tea.Cmd {
	return func() tea.Msg {
		profiles := m.profileManager.List()
		return profilesLoadedMsg{
			profiles: profiles,
			err:      nil,
		}
	}
}

// deleteProfileCmd deletes a profile
func (m Model) deleteProfileCmd(profID string) tea.Cmd {
	return func() tea.Msg {
		err := m.profileManager.Delete(profID)
		return profileDeletedMsg{
			profileID: profID,
			err:       err,
		}
	}
}

// activateProfileCmd activates a profile
func (m Model) activateProfileCmd(prof *profile.Profile) tea.Cmd {
	return func() tea.Msg {
		err := m.profileManager.SetActive(prof.ID)
		if err != nil {
			return profileActivatedMsg{
				profile: nil,
				err:     err,
			}
		}
		return profileActivatedMsg{
			profile: prof,
			err:     nil,
		}
	}
}

// showExtensionEditForm shows the extension edit form
func (m Model) showExtensionEditForm(ext *extension.Extension) (Model, tea.Cmd) {
	LogDebug("showExtensionEditForm called for extension: %s", ext.Name)
	
	// Create form
	form := NewExtensionEditForm(ext)
	LogDebug("Created form, setting size to %dx%d", m.windowWidth, m.windowHeight)
	form.SetSize(m.windowWidth, m.windowHeight)
	form.SetRenderer(m.markdownRenderer) // Use cached renderer
	form.SetCallbacks(
		func(e *extension.Extension) tea.Cmd {
			// Save extension - return saved message
			return func() tea.Msg {
				return extensionSavedMsg{extension: e}
			}
		},
		func() tea.Cmd {
			// Cancel - return close modal message
			return func() tea.Msg {
				return closeModalMsg{}
			}
		},
	)
	
	m.showingModal = true
	m.modal = &form
	m.err = nil
	
	LogDebug("Modal set, showingModal=%v, calling Init", m.showingModal)
	// Initialize the form
	initCmd := form.Init()
	LogDebug("Init returned cmd: %v", initCmd)
	return m, initCmd
}