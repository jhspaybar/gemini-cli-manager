package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width      int
	height     int
	ready      bool
	showRulers bool
}

func initialModel() model {
	return model{showRulers: true}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.showRulers = !m.showRulers
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Define our layout
	statusBarHeight := 3 // Same as main app
	contentHeight := m.height - statusBarHeight

	// Styles
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39"))

	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(statusBarHeight).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252"))

	rulerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	// Build content area
	var contentLines []string
	
	// Add some content
	contentLines = append(contentLines, infoStyle.Render(fmt.Sprintf("Terminal Size: %d×%d", m.width, m.height)))
	contentLines = append(contentLines, fmt.Sprintf("Content Area: %d×%d", m.width, contentHeight))
	contentLines = append(contentLines, fmt.Sprintf("Status Bar: %d×%d", m.width, statusBarHeight))
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, "Press 'r' to toggle rulers, 'q' to quit")
	
	// Show row markers if enabled
	if m.showRulers {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "Content rows:")
		for i := 0; i < contentHeight-len(contentLines)-2; i++ {
			if i%5 == 0 {
				contentLines = append(contentLines, rulerStyle.Render(fmt.Sprintf("  Row %d", i)))
			}
		}
	}

	content := contentStyle.Render(strings.Join(contentLines, "\n"))

	// Build status bar with exact positioning info
	statusContent := fmt.Sprintf(
		" Status Bar | Starts at row %d | Height: %d | Total height: %d ",
		contentHeight, statusBarHeight, m.height,
	)
	
	// Add position markers
	if m.showRulers {
		statusContent += fmt.Sprintf("| Rows %d-%d ", contentHeight, m.height-1)
	}

	statusBar := statusStyle.Render(statusContent)

	// Combine without any spacing
	return content + statusBar
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}