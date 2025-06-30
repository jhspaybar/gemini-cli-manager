package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/76creates/stickers/flexbox"
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

	// Exactly mimic the main app's layout calculation
	tabBarHeight := 2
	statusBarHeight := 3
	
	// Create flexbox with exact terminal dimensions
	fb := flexbox.New(m.width, m.height)
	
	// Tab bar row
	tabRow := fb.NewRow()
	tabCell := flexbox.NewCell(1, 1)
	tabContent := lipgloss.NewStyle().
		Width(m.width).
		Height(tabBarHeight).
		Background(lipgloss.Color("235")).
		Render("[ Extensions ] [ Profiles ]")
	tabCell.SetContent(tabContent)
	tabRow.AddCells(tabCell)
	tabRow.LockHeight(tabBarHeight)
	
	// Content row (fills remaining space)
	contentRow := fb.NewRow()
	contentCell := flexbox.NewCell(1, 1)
	
	// Calculate content height the same way
	calculatedContentHeight := m.height - tabBarHeight - statusBarHeight
	
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(calculatedContentHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39"))
	
	// Content with debug info
	var lines []string
	lines = append(lines, fmt.Sprintf("Terminal Size: %d×%d", m.width, m.height))
	lines = append(lines, fmt.Sprintf("Tab Bar Height: %d", tabBarHeight))
	lines = append(lines, fmt.Sprintf("Status Bar Height: %d", statusBarHeight))
	lines = append(lines, fmt.Sprintf("Calculated Content Height: %d", calculatedContentHeight))
	lines = append(lines, fmt.Sprintf("Expected Status Bar Start: row %d", tabBarHeight+calculatedContentHeight))
	lines = append(lines, "")
	lines = append(lines, "This uses the exact same layout calculation as the main app")
	
	contentCell.SetContent(contentStyle.Render(strings.Join(lines, "\n")))
	contentRow.AddCells(contentCell)
	
	// Status bar row
	statusRow := fb.NewRow()
	statusCell := flexbox.NewCell(1, 1)
	
	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(statusBarHeight).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("252"))
	
	statusContent := fmt.Sprintf(" Status: Should fill rows %d-%d (last row: %d) ",
		m.height-statusBarHeight, m.height-1, m.height-1)
	
	statusCell.SetContent(statusStyle.Render(statusContent))
	statusRow.AddCells(statusCell)
	statusRow.LockHeight(statusBarHeight)
	
	// Add all rows
	fb.AddRows([]*flexbox.Row{tabRow, contentRow, statusRow})
	
	// Render
	rendered := fb.Render()
	
	// Add debug line counter
	lines = strings.Split(rendered, "\n")
	debugInfo := fmt.Sprintf("\n[DEBUG] Rendered %d lines (expected %d)", len(lines), m.height)
	
	return rendered + debugInfo
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