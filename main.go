package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jhspaybar/gemini-cli-manager/internal/cli"
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
	model := cli.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	// Check if we need to exec into Gemini
	if m, ok := finalModel.(cli.Model); ok && m.ShouldExecGemini() {
		profile, extensions, launcher := m.GetExecInfo()
		if launcher != nil && profile != nil {
			// Now we can safely exec from the main thread
			if err := launcher.Launch(profile, extensions); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to launch Gemini: %v\n", err)
				os.Exit(1)
			}
		}
	}
}