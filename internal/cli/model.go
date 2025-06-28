package cli

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/launcher"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

// ViewType represents different views in the application
type ViewType int

const (
	ViewExtensions ViewType = iota
	ViewProfiles
	ViewSettings
	ViewHelp
	ViewExtensionDetail
)

// Model represents the application state
type Model struct {
	// View state
	currentView  ViewType
	windowWidth  int
	windowHeight int

	// Navigation state
	focusedPane   PaneType // Which pane has focus
	sidebarCursor int
	sidebarItems  []SidebarItem

	// Content area state
	extensionsCursor int
	profilesCursor   int

	// Managers
	extensionManager *extension.Manager
	profileManager   *profile.Manager
	launcher         *launcher.SimpleLauncher

	// Data
	extensions     []*extension.Extension
	profiles       []*profile.Profile
	currentProfile *profile.Profile

	// UI components
	help help.Model
	keys keyMap

	// State
	ready        bool
	loading      bool
	err          error
	showingModal bool
	modal        tea.Model

	// Search
	searchBar          *components.SearchBar
	searchActive       bool
	filteredExtensions []*extension.Extension
	filteredProfiles   []*profile.Profile

	// Launch state
	shouldExecGemini bool
	execProfile      *profile.Profile
	execExtensions   []*extension.Extension

	// Temporary state for async operations
	newProfileID string // ID of newly created profile for cursor positioning

	// Detail view state
	selectedExtension *extension.Extension

	// Cached glamour renderer
	markdownRenderer *glamour.TermRenderer

	// Theme state
	currentThemeIndex int
	settingsCursor    int // Cursor for settings list
	
	// Configuration
	stateDir string // State directory path
}

// PaneType represents which pane is focused
type PaneType int

const (
	PaneSidebar PaneType = iota
	PaneContent
)

// SidebarItem represents a sidebar navigation item
type SidebarItem struct {
	Title string
	Icon  string
	View  ViewType
}

// keyMap defines our key bindings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to show in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// NewModel creates a new application model with a custom state directory
func NewModel(stateDir string) Model {
	// Use provided state directory
	extManager := extension.NewManager(filepath.Join(stateDir, "extensions"))
	profManager := profile.NewManager(filepath.Join(stateDir, "profiles"))

	// Create launcher with state directory
	geminiCLIPath := os.Getenv("GEMINI_CLI_PATH")
	if geminiCLIPath == "" {
		geminiCLIPath = "gemini" // Assume it's in PATH
	}
	launcherInstance := launcher.NewSimpleLauncher(profManager, extManager, geminiCLIPath, stateDir)

	m := Model{
		currentView:       ViewExtensions,
		focusedPane:       PaneContent,                  // Start with content focused
		sidebarCursor:     0,                            // Initialize cursor position
		extensionsCursor:  0,                            // Initialize extension cursor
		profilesCursor:    0,                            // Initialize profile cursor
		settingsCursor:    theme.GetCurrentThemeIndex(), // Initialize to current theme
		currentThemeIndex: theme.GetCurrentThemeIndex(),
		sidebarItems: []SidebarItem{
			{Title: "Extensions", Icon: "◆", View: ViewExtensions},
			{Title: "Profiles", Icon: "▣", View: ViewProfiles},
			{Title: "Settings", Icon: "⚙", View: ViewSettings},
			{Title: "Help", Icon: "?", View: ViewHelp},
		},
		extensionManager: extManager,
		profileManager:   profManager,
		launcher:         launcherInstance,
		help:             help.New(),
		keys:             keys,
		searchBar:        components.NewSearchBar(80).SetPlaceholder("Search extensions, profiles..."),
		stateDir:         stateDir,
	}

	// Data will be loaded asynchronously in Init()
	m.loading = true

	// Create markdown renderer once
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err == nil {
		m.markdownRenderer = renderer
	}

	return m
}

// loadData loads extensions and profiles
func (m *Model) loadData() {
	// Initialize profile manager
	if err := m.profileManager.Initialize(); err != nil {
		m.err = err
		return
	}

	// Load profiles
	m.profiles = m.profileManager.List()

	// Get current profile
	if current, err := m.profileManager.GetActive(); err == nil {
		m.currentProfile = current
	} else {
		// No active profile, try to activate default
		if len(m.profiles) > 0 {
			// Look for a profile named "default" or just use the first one
			var defaultProfile *profile.Profile
			for _, p := range m.profiles {
				if p.ID == "default" || p.Name == "Default" {
					defaultProfile = p
					break
				}
			}
			if defaultProfile == nil {
				// Just use the first profile
				defaultProfile = m.profiles[0]
			}

			// Set it as active
			if err := m.profileManager.SetActive(defaultProfile.ID); err == nil {
				m.currentProfile = defaultProfile
			}
		}
	}

	// Scan for extensions
	if err := m.extensionManager.Scan(); err != nil {
		m.err = err
		return
	}

	m.extensions = m.extensionManager.List()

	// Initialize filtered lists
	m.filteredExtensions = m.extensions
	m.filteredProfiles = m.profiles
}

// Init is the first function that will be called
func (m Model) Init() tea.Cmd {
	// Set loading state
	m.loading = true

	// Return batch of initialization commands
	return tea.Batch(
		m.initializeManagersCmd(),
		tea.EnterAltScreen,
	)
}

// Getter methods for testing
func (m Model) GetCurrentView() ViewType { return m.currentView }
func (m Model) GetFocusedPane() PaneType { return m.focusedPane }
func (m Model) GetSidebarCursor() int    { return m.sidebarCursor }
func (m Model) GetExtensionsCursor() int { return m.extensionsCursor }

// ShouldExecGemini returns true if we should exec Gemini after quitting
func (m Model) ShouldExecGemini() bool {
	return m.shouldExecGemini
}

// GetExecInfo returns the information needed to exec Gemini
func (m Model) GetExecInfo() (*profile.Profile, []*extension.Extension, *launcher.SimpleLauncher) {
	return m.execProfile, m.execExtensions, m.launcher
}

// getProfileExtensions returns the extensions that are in the given profile
func (m Model) getProfileExtensions(prof *profile.Profile) []*extension.Extension {
	if prof == nil {
		return nil
	}

	var profileExts []*extension.Extension

	// Build a map of all extensions by ID for quick lookup
	extMap := make(map[string]*extension.Extension)
	for _, ext := range m.extensions {
		extMap[ext.ID] = ext
	}

	// Get extensions that are in the profile
	for _, extRef := range prof.Extensions {
		if ext, exists := extMap[extRef.ID]; exists && extRef.Enabled {
			profileExts = append(profileExts, ext)
		}
	}

	return profileExts
}

// Initialization commands

// initializeManagersCmd initializes the profile and extension managers
func (m Model) initializeManagersCmd() tea.Cmd {
	return func() tea.Msg {
		// Initialize profile manager
		if err := m.profileManager.Initialize(); err != nil {
			return managersInitializedMsg{err: err}
		}

		// Scan for extensions
		if err := m.extensionManager.Scan(); err != nil {
			return managersInitializedMsg{err: err}
		}

		return managersInitializedMsg{err: nil}
	}
}

// loadInitialDataCmd loads profiles and extensions after managers are initialized
func (m Model) loadInitialDataCmd() tea.Cmd {
	return tea.Batch(
		m.loadProfilesCmd(),
		m.loadExtensionsCmd(),
		m.checkActiveProfileCmd(),
	)
}

// checkActiveProfileCmd checks and sets the active profile
func (m Model) checkActiveProfileCmd() tea.Cmd {
	return func() tea.Msg {
		// Get current profile
		if current, err := m.profileManager.GetActive(); err == nil {
			return profileActivatedMsg{profile: current, err: nil}
		}

		// No active profile, try to activate default
		profiles := m.profileManager.List()
		if len(profiles) > 0 {
			// Look for a profile named "default" or just use the first one
			var defaultProfile *profile.Profile
			for _, p := range profiles {
				if p.ID == "default" || p.Name == "Default" {
					defaultProfile = p
					break
				}
			}
			if defaultProfile == nil {
				// Just use the first profile
				defaultProfile = profiles[0]
			}

			// Set it as active
			if err := m.profileManager.SetActive(defaultProfile.ID); err == nil {
				return profileActivatedMsg{profile: defaultProfile, err: nil}
			}
		}

		// No profiles available
		return profileActivatedMsg{profile: nil, err: nil}
	}
}

// checkInitComplete checks if all initialization tasks are complete
func (m Model) checkInitComplete() bool {
	// Check that all data has been loaded
	return m.extensions != nil && m.profiles != nil && (m.currentProfile != nil || len(m.profiles) == 0)
}
