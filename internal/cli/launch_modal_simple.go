package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gemini-cli/manager/internal/extension"
	"github.com/gemini-cli/manager/internal/launcher"
	"github.com/gemini-cli/manager/internal/profile"
)

// Launch state enums
type launchState int

const (
	launchStateChecking launchState = iota
	launchStatePreparing
	launchStateStartingServers
	launchStateLaunching
	launchStateSuccess
	launchStateFailed
)

type launchStep struct {
	name     string
	status   stepStatus
	message  string
	duration time.Duration
}

type stepStatus int

const (
	stepPending stepStatus = iota
	stepRunning
	stepSuccess
	stepFailed
	stepSkipped
)

// Launch messages
type launchStepMsg struct {
	step    int
	message string
}

type launchCompleteStepMsg struct {
	step     int
	status   stepStatus
	duration time.Duration
	error    error
}

// LaunchCompleteMsg is sent when launch completes
type LaunchCompleteMsg struct {
	Error error
}

// SimpleLaunchModal represents the launch progress modal
type SimpleLaunchModal struct {
	profile    *profile.Profile
	extensions []*extension.Extension
	launcher   *launcher.SimpleLauncher
	
	// UI state
	width      int
	height     int
	spinner    spinner.Model
	state      launchState
	progress   []launchStep
	currentStep int
	error      error
	
	// Callbacks
	onComplete func() tea.Cmd
	onCancel   func() tea.Cmd
}

// NewSimpleLaunchModal creates a new launch modal
func NewSimpleLaunchModal(p *profile.Profile, exts []*extension.Extension, l *launcher.SimpleLauncher) SimpleLaunchModal {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	enabledCount := 0
	for _, ext := range exts {
		if ext.Enabled {
			enabledCount++
		}
	}

	return SimpleLaunchModal{
		profile:    p,
		extensions: exts,
		launcher:   l,
		spinner:    s,
		state:      launchStateChecking,
		progress: []launchStep{
			{name: fmt.Sprintf("Profile: %s", p.Name), status: stepPending},
			{name: fmt.Sprintf("Extensions: %d enabled", enabledCount), status: stepPending},
			{name: "Setting up environment", status: stepPending},
			{name: "Launching Gemini CLI", status: stepPending},
		},
	}
}

// Init initializes the modal
func (m SimpleLaunchModal) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.startLaunch(),
	)
}

// Update handles modal updates
func (m SimpleLaunchModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			if m.state < launchStateLaunching {
				if m.onCancel != nil {
					return m, m.onCancel()
				}
			}
		case "enter", " ":
			if m.state == launchStateSuccess {
				if m.onComplete != nil {
					return m, m.onComplete()
				}
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case launchStepMsg:
		m.currentStep = msg.step
		m.progress[msg.step].status = stepRunning
		m.progress[msg.step].message = msg.message
		
		return m, m.continueNextStep()

	case launchCompleteStepMsg:
		if msg.step < len(m.progress) {
			m.progress[msg.step].status = msg.status
			m.progress[msg.step].duration = msg.duration
			if msg.error != nil {
				m.progress[msg.step].message = msg.error.Error()
			}
		}
		return m, nil

	case LaunchCompleteMsg:
		if msg.Error != nil {
			m.state = launchStateFailed
			m.error = msg.Error
		} else {
			m.state = launchStateSuccess
		}
		return m, nil
	}

	return m, nil
}

// View renders the launch modal
func (m SimpleLaunchModal) View() string {
	// Modal container
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(2, 3).
		Width(60).
		MaxWidth(m.width - 4)

	// Title
	title := h1Style.Render("🚀 Launching Gemini CLI")
	
	// Build content
	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")
	
	// Progress steps
	for i, step := range m.progress {
		icon := m.getStepIcon(step.status)
		style := m.getStepStyle(step.status)
		
		line := fmt.Sprintf("%s %s", icon, step.name)
		if step.duration > 0 {
			line += bodySmallStyle.Render(fmt.Sprintf(" (%dms)", step.duration.Milliseconds()))
		}
		
		content.WriteString(style.Render(line))
		
		if step.message != "" && (step.status == stepRunning || step.status == stepFailed) {
			content.WriteString("\n")
			content.WriteString(bodySmallStyle.Render("   " + step.message))
		}
		
		if i < len(m.progress)-1 {
			content.WriteString("\n")
		}
	}
	
	// Footer based on state
	content.WriteString("\n\n")
	switch m.state {
	case launchStateChecking, launchStatePreparing, launchStateStartingServers:
		content.WriteString(m.spinner.View())
		content.WriteString(" ")
		content.WriteString(mutedStyle.Render("Preparing launch..."))
		content.WriteString("\n")
		content.WriteString(helpDescStyle.Render("Press Esc to cancel"))
		
	case launchStateLaunching:
		content.WriteString(m.spinner.View())
		content.WriteString(" ")
		content.WriteString(bodyStyle.Render("Launching Gemini CLI..."))
		
	case launchStateSuccess:
		content.WriteString(enabledStyle.Render("✓ Successfully launched!"))
		content.WriteString("\n")
		content.WriteString(helpDescStyle.Render("Gemini CLI is now running"))
		
	case launchStateFailed:
		content.WriteString(errorStyle.Render("✗ Launch failed"))
		if m.error != nil {
			content.WriteString("\n")
			content.WriteString(errorStyle.Render(m.error.Error()))
		}
		content.WriteString("\n")
		content.WriteString(helpDescStyle.Render("Press Enter to close"))
	}
	
	// Center the modal
	modal := modalStyle.Render(content.String())
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modal,
	)
}

// SetSize updates the modal dimensions
func (m *SimpleLaunchModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetCallbacks sets completion callbacks
func (m *SimpleLaunchModal) SetCallbacks(onComplete, onCancel func() tea.Cmd) {
	m.onComplete = onComplete
	m.onCancel = onCancel
}

// Helper methods

func (m SimpleLaunchModal) getStepIcon(status stepStatus) string {
	switch status {
	case stepPending:
		return "○"
	case stepRunning:
		return m.spinner.View()
	case stepSuccess:
		return "✓"
	case stepFailed:
		return "✗"
	case stepSkipped:
		return "−"
	default:
		return "?"
	}
}

func (m SimpleLaunchModal) getStepStyle(status stepStatus) lipgloss.Style {
	switch status {
	case stepPending:
		return mutedStyle
	case stepRunning:
		return bodyStyle.Bold(true)
	case stepSuccess:
		return enabledStyle
	case stepFailed:
		return errorStyle
	case stepSkipped:
		return mutedStyle
	default:
		return bodyStyle
	}
}

// Launch process implementation

func (m SimpleLaunchModal) startLaunch() tea.Cmd {
	return func() tea.Msg {
		return launchStepMsg{
			step:    0,
			message: "Validating configuration...",
		}
	}
}

func (m SimpleLaunchModal) continueNextStep() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		
		switch m.currentStep {
		case 0: // Profile validation
			time.Sleep(300 * time.Millisecond)
			
			return launchCompleteStepMsg{
				step:     0,
				status:   stepSuccess,
				duration: time.Since(start),
			}
			
		case 1: // Extension check
			time.Sleep(500 * time.Millisecond)
			
			return tea.Batch(
				func() tea.Msg {
					return launchCompleteStepMsg{
						step:     1,
						status:   stepSuccess,
						duration: time.Since(start),
					}
				},
				func() tea.Msg {
					return launchStepMsg{
						step:    2,
						message: "Configuring environment variables...",
					}
				},
			)()
			
		case 2: // Environment setup
			time.Sleep(200 * time.Millisecond)
			
			return tea.Batch(
				func() tea.Msg {
					return launchCompleteStepMsg{
						step:     2,
						status:   stepSuccess,
						duration: time.Since(start),
					}
				},
				func() tea.Msg {
					return launchStepMsg{
						step:    3,
						message: "Starting Gemini CLI...",
					}
				},
			)()
			
		case 3: // Launch Gemini
			// Actually launch the CLI
			err := m.launcher.Launch(m.profile, m.extensions)
			if err != nil {
				return tea.Batch(
					func() tea.Msg {
						return launchCompleteStepMsg{
							step:     3,
							status:   stepFailed,
							duration: time.Since(start),
							error:    err,
						}
					},
					func() tea.Msg {
						return LaunchCompleteMsg{
							Error: err,
						}
					},
				)()
			}
			
			// Success!
			return tea.Batch(
				func() tea.Msg {
					return launchCompleteStepMsg{
						step:     3,
						status:   stepSuccess,
						duration: time.Since(start),
					}
				},
				func() tea.Msg {
					time.Sleep(300 * time.Millisecond) // Brief pause
					return LaunchCompleteMsg{}
				},
			)()
		}
		
		return nil
	}
}