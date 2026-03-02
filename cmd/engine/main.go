package main

import (
	"fmt"
	"os"

	"github.com/s4mn0v/trade-engine/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	// Initialize the UI Model from the internal package
	m := ui.New()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running engine TUI: %v\n", err)
		os.Exit(1)
	}
}
