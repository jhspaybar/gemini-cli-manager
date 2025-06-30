package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width    int
	height   int
	ready    bool
	showGrid bool
}

func initialModel() model {
	return model{showGrid: true}
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
		case "g":
			m.showGrid = !m.showGrid
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Create styles
	gridStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")) // Dark gray for grid

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")) // Bright blue

	cornerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Bright red
		Bold(true)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Bright yellow
		Bold(true)

	rulerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("45")) // Cyan for rulers

	var lines []string

	// Top ruler (column numbers)
	if m.showGrid {
		ruler := " "
		for i := 0; i < m.width-1; i++ {
			if i%10 == 0 {
				ruler += fmt.Sprintf("%d", i/10%10)
			} else {
				ruler += " "
			}
		}
		lines = append(lines, rulerStyle.Render(ruler))

		ruler2 := " "
		for i := 0; i < m.width-1; i++ {
			ruler2 += fmt.Sprintf("%d", i%10)
		}
		lines = append(lines, rulerStyle.Render(ruler2))
	}

	// Calculate content area
	rulerHeight := 0
	if m.showGrid {
		rulerHeight = 2
	}
	contentHeight := m.height - rulerHeight

	// Build content with row numbers
	for row := 0; row < contentHeight; row++ {
		var line string

		// Add row number if grid is shown
		if m.showGrid {
			line = rulerStyle.Render(fmt.Sprintf("%2d", row)) + " "
		}

		// First row - top border
		if row == 0 {
			line += cornerStyle.Render("┌") + strings.Repeat("─", m.width-2) + cornerStyle.Render("┐")
		} else if row == contentHeight-1 {
			// Last row - bottom border
			line += cornerStyle.Render("└") + strings.Repeat("─", m.width-2) + cornerStyle.Render("┘")
		} else {
			// Middle rows
			content := ""
			
			if row == contentHeight/2-2 {
				// Terminal info
				info := fmt.Sprintf("Terminal: %d×%d (Content: %d×%d)", 
					m.width, m.height, m.width, contentHeight)
				padding := (m.width - len(info) - 2) / 2
				content = strings.Repeat(" ", padding) + infoStyle.Render(info)
				content += strings.Repeat(" ", m.width-2-len(strings.Repeat(" ", padding)+info))
			} else if row == contentHeight/2 {
				// Instructions
				info := "Press 'g' to toggle grid, 'q' to quit"
				padding := (m.width - len(info) - 2) / 2
				content = strings.Repeat(" ", padding) + info
				content += strings.Repeat(" ", m.width-2-len(strings.Repeat(" ", padding)+info))
			} else if row == contentHeight/2+2 {
				// Show exact position
				info := fmt.Sprintf("Bottom is at row %d (0-indexed)", contentHeight-1)
				padding := (m.width - len(info) - 2) / 2
				content = strings.Repeat(" ", padding) + info
				content += strings.Repeat(" ", m.width-2-len(strings.Repeat(" ", padding)+info))
			} else {
				// Empty space with dots every 5 characters for visualization
				if m.showGrid {
					for i := 0; i < m.width-2; i++ {
						if i%5 == 0 {
							content += gridStyle.Render("·")
						} else {
							content += " "
						}
					}
				} else {
					content = strings.Repeat(" ", m.width-2)
				}
			}

			line += borderStyle.Render("│") + content + borderStyle.Render("│")
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
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