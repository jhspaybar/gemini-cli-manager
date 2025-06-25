package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gemini-cli/manager/internal/cli"
)

func main() {
	// Check for debug flag
	debug := false
	for _, arg := range os.Args[1:] {
		if arg == "--debug" {
			debug = true
			break
		}
	}

	// Enable debug logging if requested
	if debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Create and run the program
	p := tea.NewProgram(cli.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}