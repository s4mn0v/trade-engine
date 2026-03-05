package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/s4mn0v/trade-engine/internal/app"
	"github.com/s4mn0v/trade-engine/internal/backtesting"
)

type SessionState int

const (
	StateDataPicker SessionState = iota
	StateStrategyPicker
	StateIndicatorPicker
	StateConfig
	StateExecuting
	StateFinished
)

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type Model struct {
	State         SessionState
	Filepicker    filepicker.Model
	Inputs        []textinput.Model
	FocusIndex    int
	Logs          []string
	ProgressPct   int
	DataFile      string
	StrategyFile  string
	IndicatorFile string
	Results       backtesting.Summary
	Quitting      bool

	// We keep these only as starting-point references
	DataRoot    string
	ScriptsRoot string
}

func New() Model {
	// Initialize Inputs
	amount := textinput.New()
	amount.Placeholder = "Investment Amount (Numbers Only)"
	amount.Focus()

	comm := textinput.New()
	comm.Placeholder = "Commission % (Default 0.06)"

	lev := textinput.New()
	lev.Placeholder = "Leverage (Default 1.0)"

	// Setup Starting Paths
	cwd, _ := os.Getwd()
	dataRoot := filepath.Join(cwd, "data")
	scriptsRoot := filepath.Join(cwd, "scripts")

	// Ensure folders exist for convenience, but we won't lock the user here
	os.MkdirAll(dataRoot, 0o755)
	os.MkdirAll(scriptsRoot, 0o755)

	// Initialize Filepicker
	fp := filepicker.New()
	fp.CurrentDirectory = dataRoot
	fp.AllowedTypes = []string{".csv"}

	return Model{
		State:       StateDataPicker,
		Filepicker:  fp,
		DataRoot:    dataRoot,
		ScriptsRoot: scriptsRoot,
		Inputs:      []textinput.Model{amount, comm, lev},
		Logs:        []string{"[SYSTEM] Waiting for data selection..."},
	}
}

func (m Model) Init() tea.Cmd {
	return m.Filepicker.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 1. GLOBAL KEY HANDLING
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		key := keyMsg.String()
		if key == "ctrl+c" || key == "q" {
			m.Quitting = true
			return m, tea.Quit
		}

		if key == "s" {
			if m.State == StateStrategyPicker {
				m.State = StateIndicatorPicker
				return m, m.Filepicker.Init()
			}
			if m.State == StateIndicatorPicker {
				m.State = StateConfig
				return m, nil
			}
		}

		if m.State == StateFinished && key == "enter" {
			return m, tea.Quit
		}
	}

	// 2. STATE-SPECIFIC UPDATES
	switch m.State {
	case StateDataPicker, StateStrategyPicker, StateIndicatorPicker:
		var cmd tea.Cmd
		m.Filepicker, cmd = m.Filepicker.Update(msg)

		// Handle Selection
		if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
			switch m.State {
			case StateDataPicker:
				m.DataFile = path
				m.State = StateStrategyPicker

				// Relocate to scripts folder for convenience, but user can still leave
				m.Filepicker.CurrentDirectory = m.ScriptsRoot
				m.Filepicker.AllowedTypes = []string{".go"}
				return m, m.Filepicker.Init()

			case StateStrategyPicker:
				m.StrategyFile = path
				m.State = StateIndicatorPicker
				return m, m.Filepicker.Init()

			case StateIndicatorPicker:
				m.IndicatorFile = path
				m.State = StateConfig
				return m, nil
			}
		}
		return m, cmd

	case StateConfig:
		if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
			key := keyMsg.String()

			if (m.FocusIndex >= 0 && m.FocusIndex <= 2) && len(key) == 1 {
				r := rune(key[0])
				if !unicode.IsDigit(r) && r != '.' {
					return m, nil
				}
			}

			switch key {
			case "up", "shift+tab":
				m.FocusIndex--
				if m.FocusIndex < 0 {
					m.FocusIndex = 3
				}
			case "down", "tab":
				m.FocusIndex++
				if m.FocusIndex > 3 {
					m.FocusIndex = 0
				}
			case "enter":
				if m.FocusIndex < 2 {
					m.FocusIndex++
				} else {
					if m.Inputs[0].Value() == "" {
						return m, nil
					}
					if m.Inputs[1].Value() == "" {
						m.Inputs[1].SetValue("0.06")
					}
					if m.Inputs[2].Value() == "" {
						m.Inputs[2].SetValue("1.0")
					}
					m.State = StateExecuting
					return m, Tick()
				}
			}

			var cmd tea.Cmd
			cmds := make([]tea.Cmd, 0)
			for i := range m.Inputs {
				if i == m.FocusIndex {
					cmds = append(cmds, m.Inputs[i].Focus())
				} else {
					m.Inputs[i].Blur()
				}
				m.Inputs[i], cmd = m.Inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

	case StateExecuting:
		if _, ok := msg.(TickMsg); ok {
			if m.ProgressPct >= 100 {
				inv, _ := strconv.ParseFloat(m.Inputs[0].Value(), 64)
				comm, _ := strconv.ParseFloat(m.Inputs[1].Value(), 64)
				lev, _ := strconv.ParseFloat(m.Inputs[2].Value(), 64)

				summary, err := app.RunFullBacktest(m.DataFile, m.StrategyFile, m.IndicatorFile, inv, comm, lev)
				if err != nil {
					m.Logs = append(m.Logs, fmt.Sprintf("[ERROR] %v", err))
					return m, nil
				}

				m.Results = summary
				m.State = StateFinished
				return m, nil
			}
			m.ProgressPct += 10
			m.Logs = append(m.Logs, fmt.Sprintf("[%s] Analyzing market data...", time.Now().Format("15:04:05")))
			if len(m.Logs) > 8 {
				m.Logs = m.Logs[1:]
			}
			return m, Tick()
		}
	}

	return m, nil
}
