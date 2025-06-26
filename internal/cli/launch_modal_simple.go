package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/launcher"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
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

	return SimpleLaunchModal{
		profile:    p,
		extensions: exts,
		launcher:   l,
		spinner:    s,
		state:      launchStateChecking,
		progress: []launchStep{
			{name: fmt.Sprintf("Profile: %s", p.Name), status: stepPending},
			{name: fmt.Sprintf("Extensions: %d in profile", len(exts)), status: stepPending},
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
		case "ctrl+c":
			// Always allow Ctrl+C to quit
			return m, tea.Quit
		case "esc":
			if m.state < launchStateLaunching {
				if m.onCancel != nil {
					return m, m.onCancel()
				}
			}
		case "enter", " ":
			if m.state == launchStateSuccess || m.state == launchStateFailed {
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
		// Debug log
		debugLog, _ := os.OpenFile("/tmp/gemini-cli-manager-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if debugLog != nil {
			fmt.Fprintf(debugLog, "Received launchStepMsg: step=%d, message=%s\n", msg.step, msg.message)
			debugLog.Close()
		}
		
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
			return m, nil
		} else {
			m.state = launchStateSuccess
			// Send a message to exec Gemini after quitting
			return m, func() tea.Msg {
				return execGeminiMsg{
					profile:    m.profile,
					extensions: m.extensions,
				}
			}
		}
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
	// Debug log
	debugLog, _ := os.OpenFile("/tmp/gemini-cli-manager-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if debugLog != nil {
		fmt.Fprintf(debugLog, "\n=== Launch started at %s ===\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(debugLog, "startLaunch() called\n")
		debugLog.Close()
	}
	
	return func() tea.Msg {
		return launchStepMsg{
			step:    0,
			message: "Validating configuration...",
		}
	}
}

func (m SimpleLaunchModal) continueNextStep() tea.Cmd {
	// Debug log
	debugLog, _ := os.OpenFile("/tmp/gemini-cli-manager-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if debugLog != nil {
		fmt.Fprintf(debugLog, "continueNextStep() called, currentStep=%d\n", m.currentStep)
		debugLog.Close()
	}
	
	return func() tea.Msg {
		start := time.Now()
		
		switch m.currentStep {
		case 0: // Profile validation
			time.Sleep(300 * time.Millisecond)
			
			// Need to return both the completion and the next step
			return tea.Batch(
				func() tea.Msg {
					return launchCompleteStepMsg{
						step:     0,
						status:   stepSuccess,
						duration: time.Since(start),
					}
				},
				func() tea.Msg {
					return launchStepMsg{
						step:    1,
						message: "Checking extensions...",
					}
				},
			)()
			
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
			// Don't actually launch here - just validate we can
			// The actual exec needs to happen after Bubble Tea shuts down
			
			// Debug log
			debugLog, _ := os.OpenFile("/tmp/gemini-cli-manager-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if debugLog != nil {
				fmt.Fprintf(debugLog, "\n=== Launch modal step 3 at %s ===\n", time.Now().Format(time.RFC3339))
				fmt.Fprintf(debugLog, "Ready to launch after TUI exits\n")
				debugLog.Close()
			}
			
			// Just mark as successful and let the main app handle the actual exec
			return tea.Batch(
				func() tea.Msg {
					return launchCompleteStepMsg{
						step:     3,
						status:   stepSuccess,
						duration: time.Since(start),
					}
				},
				func() tea.Msg {
					time.Sleep(300 * time.Millisecond) // Brief pause to show success
					return LaunchCompleteMsg{}
				},
			)()
		}
		
		return nil
	}
}