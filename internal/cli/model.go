package cli

import (
	"os"
	"path/filepath"
	
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gemini-cli/manager/internal/extension"
	"github.com/gemini-cli/manager/internal/launcher"
	"github.com/gemini-cli/manager/internal/profile"
)

// ViewType represents different views in the application
type ViewType int

const (
	ViewExtensions ViewType = iota
	ViewProfiles
	ViewSettings
	ViewHelp
)

// Model represents the application state
type Model struct {
	// View state
	currentView  ViewType
	windowWidth  int
	windowHeight int

	// Navigation state
	focusedPane  PaneType // Which pane has focus
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
	searchBar        SearchBar
	searchActive     bool
	filteredExtensions []*extension.Extension
	filteredProfiles   []*profile.Profile
	
	// Launch state
	shouldExecGemini bool
	execProfile      *profile.Profile
	execExtensions   []*extension.Extension
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

// NewModel creates a new application model
func NewModel() Model {
	// Initialize managers
	homePath := os.Getenv("HOME")
	if homePath == "" {
		homePath = "."
	}
	// Use our own state directory, separate from Gemini's
	managerPath := filepath.Join(homePath, ".gemini-cli-manager")
	
	extManager := extension.NewManager(filepath.Join(managerPath, "extensions"))
	profManager := profile.NewManager(filepath.Join(managerPath, "profiles"))
	
	// Create launcher
	geminiCLIPath := os.Getenv("GEMINI_CLI_PATH")
	if geminiCLIPath == "" {
		geminiCLIPath = "gemini" // Assume it's in PATH
	}
	launcherInstance := launcher.NewSimpleLauncher(profManager, extManager, geminiCLIPath)
	
	m := Model{
		currentView:      ViewExtensions,
		focusedPane:      PaneContent, // Start with content focused
		sidebarCursor:    0,           // Initialize cursor position
		extensionsCursor: 0,           // Initialize extension cursor
		profilesCursor:   0,           // Initialize profile cursor
		sidebarItems: []SidebarItem{
			{Title: "Extensions", Icon: "◆", View: ViewExtensions},
			{Title: "Profiles", Icon: "▣", View: ViewProfiles},
			{Title: "Settings", Icon: "⚙", View: ViewSettings},
			{Title: "Help", Icon: "?", View: ViewHelp},
		},
		extensionManager: extManager,
		profileManager:   profManager,
		launcher:         launcherInstance,
		help:            help.New(),
		keys:            keys,
		searchBar:        NewSearchBar("Search extensions, profiles..."),
	}
	
	// Load initial data
	m.loadData()
	
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
	return nil
}

// Getter methods for testing
func (m Model) GetCurrentView() ViewType { return m.currentView }
func (m Model) GetFocusedPane() PaneType { return m.focusedPane }
func (m Model) GetSidebarCursor() int { return m.sidebarCursor }
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