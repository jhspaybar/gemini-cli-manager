package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jhspaybar/gemini-cli-manager/internal/cli"
)

func main() {
	// Define flags
	var (
		debug    = flag.Bool("debug", false, "Enable debug logging")
		stateDir = flag.String("state-dir", "", "Directory for storing application state (default: ~/.gemini-cli-manager)")
		help     = flag.Bool("help", false, "Show help")
		helpShort = flag.Bool("h", false, "Show help")
	)

	flag.Parse()

	// Show help if requested
	if *help || *helpShort {
		fmt.Println("Gemini CLI Manager - Manage extensions, prompts, and tools for Gemini")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Printf("  %s [flags]\n", os.Args[0])
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --state-dir string   Directory for storing application state (default: ~/.gemini-cli-manager)")
		fmt.Println("  --debug              Enable debug logging")
		fmt.Println("  -h, --help           Show help")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  # Use default state directory")
		fmt.Printf("  %s\n", os.Args[0])
		fmt.Println()
		fmt.Println("  # Use custom state directory for testing")
		fmt.Printf("  %s --state-dir /tmp/gemini-test\n", os.Args[0])
		fmt.Println()
		fmt.Println("  # Run multiple independent setups")
		fmt.Printf("  %s --state-dir ~/gemini-work\n", os.Args[0])
		fmt.Printf("  %s --state-dir ~/gemini-personal\n", os.Args[0])
		os.Exit(0)
	}

	// Set default state directory if not provided
	if *stateDir == "" {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			homePath = "."
		}
		*stateDir = filepath.Join(homePath, ".gemini-cli-manager")
	}

	// Expand home directory if used
	if len(*stateDir) > 0 && (*stateDir)[0] == '~' {
		homePath := os.Getenv("HOME")
		if homePath != "" {
			*stateDir = filepath.Join(homePath, (*stateDir)[1:])
		}
	}

	// Make path absolute
	absStateDir, err := filepath.Abs(*stateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving state directory path: %v\n", err)
		os.Exit(1)
	}
	*stateDir = absStateDir

	// Enable debug logging if requested
	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Create and run the program with state directory
	model := cli.NewModel(*stateDir)
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
