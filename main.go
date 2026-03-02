package main

import (
	"fmt"
	"os"
	"time"
	"unicode"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type sessionState int

const (
	stateDataPicker sessionState = iota
	stateStrategyPicker
	stateIndicatorPicker
	stateConfig
	stateExecuting
	stateFinished
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).MarginBottom(1)
	logStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	docStyle     = lipgloss.NewStyle().Padding(1, 2)
)

type model struct {
	state         sessionState
	filepicker    filepicker.Model
	inputs        []textinput.Model
	focusIndex    int // 0: Amount, 1: Comm, 2: Execute Button
	logs          []string
	progressPct   int
	dataFile      string
	strategyFile  string
	indicatorFile string
	quitting      bool
}

func initialModel() model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".csv"} // Start with Data selection
	fp.CurrentDirectory, _ = os.Getwd()

	amount := textinput.New()
	amount.Placeholder = "Investment Amount (Numbers Only)"
	amount.Focus()

	comm := textinput.New()
	comm.Placeholder = "Commission % (Default 0.06)"

	return model{
		state:      stateDataPicker,
		filepicker: fp,
		inputs:     []textinput.Model{amount, comm},
		logs:       []string{"[SYSTEM] Ready for data selection..."},
	}
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		// Handle Skip 's' for Strategy and Indicator
		if key == "s" {
			if m.state == stateStrategyPicker {
				m.state = stateIndicatorPicker
				return m, nil
			}
			if m.state == stateIndicatorPicker {
				m.state = stateConfig
				return m, nil
			}
		}
	}

	switch m.state {
	case stateDataPicker:
		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.dataFile = path
			m.state = stateStrategyPicker
			m.filepicker.AllowedTypes = []string{".go"} // LOCK to Golang only
			return m, nil
		}
		return m, cmd

	case stateStrategyPicker:
		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.strategyFile = path
			m.state = stateIndicatorPicker
			m.filepicker.AllowedTypes = []string{".go"} // LOCK to Golang only
			return m, nil
		}
		return m, cmd

	case stateIndicatorPicker:
		var cmd tea.Cmd
		m.filepicker, cmd = m.filepicker.Update(msg)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.indicatorFile = path
			m.state = stateConfig
			return m, nil
		}
		return m, cmd

	case stateConfig:
		if msg, ok := msg.(tea.KeyPressMsg); ok {
			key := msg.String()

			// Block non-numeric input for Investment Amount
			if m.focusIndex == 0 && len(key) == 1 {
				r := rune(key[0])
				if !unicode.IsDigit(r) && r != '.' {
					return m, nil
				}
			}

			switch key {
			case "up", "shift+tab":
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = 2
				} // Wrap to button
			case "down", "tab":
				m.focusIndex++
				if m.focusIndex > 2 {
					m.focusIndex = 0
				} // Wrap to start
			case "enter":
				// If on inputs, move to next. If on button or last input, EXECUTE.
				if m.focusIndex < 1 {
					m.focusIndex++
				} else {
					// Validate & Default
					if m.inputs[0].Value() == "" {
						return m, nil
					}
					if m.inputs[1].Value() == "" {
						m.inputs[1].SetValue("0.06")
					}

					m.state = stateExecuting
					return m, tick()
				}
			}

			// Apply focus
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			// Update inputs
			var cmd tea.Cmd
			for i := range m.inputs {
				m.inputs[i], cmd = m.inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

	case stateExecuting:
		if _, ok := msg.(tickMsg); ok {
			if m.progressPct >= 100 {
				m.state = stateFinished
				return m, nil
			}
			m.progressPct += 5
			m.logs = append(m.logs, fmt.Sprintf("Processing block %d...", m.progressPct*10))
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
		return tea.NewView("Exiting...")
	}

	var content string

	switch m.state {
	case stateDataPicker:
		content = titleStyle.Render("STEP 1: SELECT DATA (.CSV)") + "\n" +
			m.filepicker.View()

	case stateStrategyPicker:
		content = titleStyle.Render("STEP 2: SELECT STRATEGY (.GO)") + "\n" +
			"Choose logic or press " + focusedStyle.Render("'s'") + " to skip.\n\n" +
			m.filepicker.View()

	case stateIndicatorPicker:
		content = titleStyle.Render("STEP 3: SELECT INDICATOR (.GO)") + "\n" +
			"Choose indicator or press " + focusedStyle.Render("'s'") + " to skip.\n\n" +
			m.filepicker.View()

	case stateConfig:
		content = titleStyle.Render("STEP 4: BACKTEST CONFIG") + "\n" +
			fmt.Sprintf("Data:      %s\nStrategy:  %s\nIndicator: %s\n\n",
				displayFile(m.dataFile), displayFile(m.strategyFile), displayFile(m.indicatorFile))

		content += "Investment: " + m.inputs[0].View() + "\n"
		content += "Commission: " + m.inputs[1].View() + " %\n"

		btn := "\n[ EXECUTE STRATEGY ]"
		if m.focusIndex == 2 {
			btn = "\n" + focusedStyle.Render("[ EXECUTE STRATEGY ]")
		}
		content += btn

	case stateExecuting:
		content = titleStyle.Render("STEP 5: ANALYSIS CANVAS") + "\n"
		logBox := ""
		for _, l := range m.logs {
			logBox += logStyle.Render("> "+l) + "\n"
		}
		content += lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).Width(60).Height(10).Render(logBox)

	case stateFinished:
		content = titleStyle.Render("STEP 6: COMPLETE") + "\n" +
			focusedStyle.Render("SUCCESS: BACKTEST FINISHED") + "\n\n" +
			fmt.Sprintf("Logs saved to testing.txt\n\nPress Enter to Exit.")
	}

	v := tea.NewView(docStyle.Render(content))
	if m.state == stateExecuting {
		v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.progressPct)
	}
	return v
}

func displayFile(p string) string {
	if p == "" {
		return focusedStyle.Render("(none)")
	}
	return focusedStyle.Render(p)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
