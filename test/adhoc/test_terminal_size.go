package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width  int
	height int
	ready  bool
}

func initialModel() model {
	return model{}
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
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Create styles
	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))  // Bright blue for visibility

	cornerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Bright red for corners
		Bold(true)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Bright yellow
		Bold(true)

	// Build the view with exact dimensions
	var lines []string

	// Top border with corner markers
	topLine := cornerStyle.Render("┌") + strings.Repeat("─", m.width-2) + cornerStyle.Render("┐")
	lines = append(lines, topLine)

	// Fill middle with dimension info
	middleHeight := m.height - 2 // Account for top and bottom borders
	for i := 0; i < middleHeight; i++ {
		if i == middleHeight/2-2 {
			// Show terminal dimensions
			info := fmt.Sprintf("Terminal Size: %d×%d", m.width, m.height)
			padding := (m.width - len(info) - 2) / 2
			line := "│" + strings.Repeat(" ", padding) + infoStyle.Render(info) + strings.Repeat(" ", m.width-padding-len(info)-2) + "│"
			lines = append(lines, borderStyle.Render(line))
		} else if i == middleHeight/2-1 {
			// Show position info
			info := "Press 'q' to quit"
			padding := (m.width - len(info) - 2) / 2
			line := "│" + strings.Repeat(" ", padding) + info + strings.Repeat(" ", m.width-padding-len(info)-2) + "│"
			lines = append(lines, borderStyle.Render(line))
		} else if i == middleHeight/2+1 {
			// Show line numbers at edges
			leftInfo := fmt.Sprintf("Line %d", i+2)
			rightInfo := fmt.Sprintf("Col %d", m.width)
			middleSpace := m.width - len(leftInfo) - len(rightInfo) - 2
			line := "│" + leftInfo + strings.Repeat(" ", middleSpace) + rightInfo + "│"
			lines = append(lines, borderStyle.Render(line))
		} else {
			// Empty line with borders
			line := "│" + strings.Repeat(" ", m.width-2) + "│"
			lines = append(lines, borderStyle.Render(line))
		}
	}

	// Bottom border with corner markers
	bottomLine := cornerStyle.Render("└") + strings.Repeat("─", m.width-2) + cornerStyle.Render("┘")
	lines = append(lines, bottomLine)

	// Join all lines
	return strings.Join(lines, "\n")
}

func main() {
	// Create program with alt screen
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // Use alternate screen like main app
		tea.WithMouseCellMotion(), // Enable mouse to match main app
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}