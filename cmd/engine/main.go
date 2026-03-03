package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/s4mn0v/trade-engine/internal/ui"
)

func main() {
	// 1. Create the UI Model
	// The internal/ui package handles the step-by-step flow
	m := ui.New()

	// 2. Initialize the Bubble Tea Program
	p := tea.NewProgram(m)

	// 3. Run the Program
	// Run() blocks until the user quits or the backtest finishes and user hits Enter
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, the engine has encountered an error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Engine shut down gracefully. Check results.txt for details.")
}
