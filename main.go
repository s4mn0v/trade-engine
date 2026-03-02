package main

import (
	"fmt"
	"os"
	"time"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type sessionState int

const (
	stateFilePicker sessionState = iota
	stateConfig
	stateExecuting
	stateFinished
)

// Styles
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).MarginBottom(1)
	logStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	docStyle     = lipgloss.NewStyle().Padding(1, 2)
)

type model struct {
	state        sessionState
	filepicker   filepicker.Model
	inputs       []textinput.Model
	focusIndex   int
	logs         []string
	progressPct  int // 0-100 for the v2 ProgressBar
	selectedFile string
	quitting     bool
	width        int
}

func initialModel() model {
	// 1. Setup Filepicker
	fp := filepicker.New()
	fp.AllowedTypes = []string{".csv"}
	fp.CurrentDirectory, _ = os.Getwd()

	// 2. Setup Text Inputs (Amount and Commission)
	amount := textinput.New()
	amount.Placeholder = "Investment Amount (e.g. 1000)"
	amount.Focus()

	comm := textinput.New()
	comm.Placeholder = "Commission % (e.g. 0.1)"

	return model{
		state:       stateFilePicker,
		filepicker:  fp,
		inputs:      []textinput.Model{amount, comm},
		logs:        []string{"[SYSTEM] Engine initialized.", "[SYSTEM] Awaiting data source..."},
		progressPct: 0,
	}
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

// Simulated Tick for UI demonstration
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	switch m.state {
	case stateFilePicker:
		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)

		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			m.state = stateConfig
			return m, nil
		}
		return m, cmd

	case stateConfig:
		if msg, ok := msg.(tea.KeyPressMsg); ok {
			switch msg.String() {
			case "tab", "shift+tab", "up", "down":
				if msg.String() == "up" || msg.String() == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				// Cycle focus between inputs and the "Execute" pseudo-button
				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i < len(m.inputs); i++ {
					if i == m.focusIndex {
						cmds[i] = m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
				return m, tea.Batch(cmds...)

			case "enter":
				// If on last input or on the "pseudo-button" index
				m.state = stateExecuting
				return m, tick()
			}
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)

	case stateExecuting:
		switch msg.(type) {
		case tickMsg:
			if m.progressPct >= 100 {
				m.state = stateFinished
				return m, nil
			}

			// UI Simulation Logic
			m.progressPct += 5
			m.logs = append(m.logs, fmt.Sprintf("Evaluating Strategy Logic at index %d...", m.progressPct*12))
			if len(m.logs) > 8 {
				m.logs = m.logs[1:]
			}
			return m, tick()
		}

	case stateFinished:
		if msg, ok := msg.(tea.KeyPressMsg); ok && msg.String() == "enter" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("Closing strategy engine...")
	}

	var content string

	switch m.state {
	case stateFilePicker:
		content = titleStyle.Render("STEP 1: DATA SELECTION") + "\n" +
			"Select the .csv file for backtesting\n\n" +
			m.filepicker.View()

	case stateConfig:
		content = titleStyle.Render("STEP 2: PARAMETERS") + "\n" +
			fmt.Sprintf("File: %s\n\n", focusedStyle.Render(m.selectedFile))

		for i := range m.inputs {
			content += m.inputs[i].View() + "\n"
		}

		button := "\n[ EXECUTE STRATEGY ]"
		if m.focusIndex == len(m.inputs) {
			button = "\n" + focusedStyle.Render("[ EXECUTE STRATEGY ]")
		}
		content += button

	case stateExecuting:
		content = titleStyle.Render("STEP 3: EXECUTION CANVAS") + "\n" +
			fmt.Sprintf("Processing: %s\n\n", m.selectedFile)

		// The "Canvas" area
		logBox := ""
		for _, l := range m.logs {
			logBox += logStyle.Render("> "+l) + "\n"
		}

		canvas := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Width(60).
			Height(10).
			Render(logBox)

		content += canvas + "\n\n(The progress bar is displayed in the terminal status section)"

	case stateFinished:
		content = titleStyle.Render("STEP 4: RESULTS") + "\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render("STRATEGY COMPLETED SUCCESSFULLY") + "\n\n" +
			fmt.Sprintf("Configured Amount: $%s\n", m.inputs[0].Value()) +
			fmt.Sprintf("Configured Comm:   %s%%\n", m.inputs[1].Value()) +
			"Status:            Logs saved to testing.txt\n\n" +
			"Press Enter to exit."
	}

	// Create the view
	v := tea.NewView(docStyle.Render(content))

	// Handle Progress Bar via tea.View properties (v2 way)
	if m.state == stateExecuting {
		v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.progressPct)
	} else if m.state == stateFinished {
		v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, 100)
	}

	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fatal Error: %v", err)
		os.Exit(1)
	}
}
