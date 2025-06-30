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

	// Test different ways of setting height
	
	// Method 1: Using Height() and filling with content
	method1 := lipgloss.NewStyle().
		Width(m.width).
		Height(3).
		Background(lipgloss.Color("235")).
		Render("Method 1: Height(3)")
	
	// Method 2: Using MaxHeight and content
	method2 := lipgloss.NewStyle().
		Width(m.width).
		MaxHeight(3).
		Background(lipgloss.Color("236")).
		Render("Method 2: MaxHeight(3)\nLine 2\nLine 3")
	
	// Method 3: Manual newlines to exact height
	lines := []string{"Method 3: Manual newlines"}
	for i := 1; i < 3; i++ {
		lines = append(lines, fmt.Sprintf("Line %d", i+1))
	}
	method3 := lipgloss.NewStyle().
		Width(m.width).
		Background(lipgloss.Color("237")).
		Render(strings.Join(lines, "\n"))
	
	// Method 4: Fill remaining space manually
	usedHeight := lipgloss.Height(method1) + lipgloss.Height(method2) + lipgloss.Height(method3)
	remainingHeight := m.height - usedHeight - 3 // -3 for the debug info below
	
	var fillLines []string
	fillLines = append(fillLines, "Method 4: Fill remaining space")
	for i := 1; i < remainingHeight; i++ {
		fillLines = append(fillLines, fmt.Sprintf("Fill line %d/%d", i, remainingHeight-1))
	}
	method4 := lipgloss.NewStyle().
		Width(m.width).
		Background(lipgloss.Color("238")).
		Render(strings.Join(fillLines, "\n"))
	
	// Combine all methods
	combined := lipgloss.JoinVertical(
		lipgloss.Left,
		method1,
		method2,
		method3,
		method4,
	)
	
	// Debug info
	debugInfo := fmt.Sprintf("\nTerminal: %dx%d | Used: %d | Remaining: %d | M1: %d, M2: %d, M3: %d, M4: %d",
		m.width, m.height, usedHeight, remainingHeight,
		lipgloss.Height(method1), lipgloss.Height(method2), 
		lipgloss.Height(method3), lipgloss.Height(method4))
	
	// Count actual lines
	result := combined + debugInfo
	actualLines := strings.Count(result, "\n") + 1
	result += fmt.Sprintf(" | Lines: %d", actualLines)
	
	return result
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