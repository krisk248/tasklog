package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/krisk248/tasklog/internal/app"
)

func main() {
	// Create new model
	m := app.NewModel()

	// Create Bubbletea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running tasklog: %v\n", err)
		os.Exit(1)
	}
}
