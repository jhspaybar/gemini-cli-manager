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

	// Mirror the main app's layout calculation
	horizontalPadding := 1  // Same as main app
	verticalPadding := 0    // Same as main app
	contentWidth := m.width - (horizontalPadding * 2)
	contentHeight := m.height - verticalPadding

	// Tab bar (fixed height of 2)
	tabHeight := 2
	tabBar := lipgloss.NewStyle().
		Width(contentWidth).
		Height(tabHeight).
		Background(lipgloss.Color("235")).
		Render("[ Extensions ] [ Profiles ]")

	// Status bar (fixed height of 3)
	statusHeight := 3
	statusBar := lipgloss.NewStyle().
		Width(contentWidth).
		Height(statusHeight).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252")).
		Render(fmt.Sprintf(" Status Bar | Height: %d | Should reach row %d (last row) ",
			statusHeight, m.height-1))

	// Content area calculation
	totalChrome := tabHeight + statusHeight + 4 // Same calculation as main app
	contentBoxHeight := contentHeight - totalChrome
	if contentBoxHeight < 3 {
		contentBoxHeight = 3
	}

	// Content box
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentBoxHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1)

	// Debug info
	var lines []string
	lines = append(lines, fmt.Sprintf("Terminal Size: %d×%d", m.width, m.height))
	lines = append(lines, fmt.Sprintf("Content Size: %d×%d", contentWidth, contentHeight))
	lines = append(lines, fmt.Sprintf("Tab Height: %d", tabHeight))
	lines = append(lines, fmt.Sprintf("Status Height: %d", statusHeight))
	lines = append(lines, fmt.Sprintf("Total Chrome: %d", totalChrome))
	lines = append(lines, fmt.Sprintf("Content Box Height: %d", contentBoxHeight))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Expected layout:"))
	lines = append(lines, fmt.Sprintf("  Tabs: rows 0-%d", tabHeight-1))
	lines = append(lines, fmt.Sprintf("  Content: rows %d-%d", tabHeight, tabHeight+contentBoxHeight-1))
	lines = append(lines, fmt.Sprintf("  Status: rows %d-%d", tabHeight+contentBoxHeight, m.height-1))
	
	content := contentStyle.Render(strings.Join(lines, "\n"))

	// Join components exactly like main app
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabBar,
		content,
		statusBar,
	)

	// Apply padding like main app
	final := lipgloss.NewStyle().
		PaddingLeft(horizontalPadding).
		PaddingRight(horizontalPadding).
		Render(fullContent)

	// Add line count debug
	lineCount := strings.Count(final, "\n") + 1
	if lineCount < m.height {
		final += fmt.Sprintf("\n\n[DEBUG] Rendered %d lines, terminal has %d lines (%d short)",
			lineCount, m.height, m.height-lineCount)
	}

	return final
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