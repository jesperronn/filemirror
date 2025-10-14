package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Get initial query from command line args
	initialQuery := ""
	if len(os.Args) > 1 {
		initialQuery = os.Args[1]
	}

	// Create the model
	m := initialModel(initialQuery)

	// Start the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
