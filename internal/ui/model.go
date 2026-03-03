package ui

import (
	"fmt"
	"os"
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
	Results       backtesting.Summary // FIXED: Use backtesting.Summary instead of app.BacktestResult
	Quitting      bool
}

func New() Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".csv"}
	fp.CurrentDirectory, _ = os.Getwd()

	amount := textinput.New()
	amount.Placeholder = "Investment Amount (Numbers Only)"
	amount.Focus()

	comm := textinput.New()
	comm.Placeholder = "Commission % (Default 0.06)"

	return Model{
		State:      StateDataPicker,
		Filepicker: fp,
		Inputs:     []textinput.Model{amount, comm},
		Logs:       []string{"[SYSTEM] Waiting for data selection..."},
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
				return m, nil
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

	switch m.State {
	case StateDataPicker:
		var cmd tea.Cmd
		m.Filepicker, cmd = m.Filepicker.Update(msg)
		if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
			m.DataFile = path
			m.State = StateStrategyPicker
			m.Filepicker.AllowedTypes = []string{".go"}
			return m, nil
		}
		return m, cmd

	case StateStrategyPicker:
		var cmd tea.Cmd
		m.Filepicker, cmd = m.Filepicker.Update(msg)
		if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
			m.StrategyFile = path
			m.State = StateIndicatorPicker
			m.Filepicker.AllowedTypes = []string{".go"}
			return m, nil
		}
		return m, cmd

	case StateIndicatorPicker:
		var cmd tea.Cmd
		m.Filepicker, cmd = m.Filepicker.Update(msg)
		if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
			m.IndicatorFile = path
			m.State = StateConfig
			return m, nil
		}
		return m, cmd

	case StateConfig:
		if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
			key := keyMsg.String()

			if m.FocusIndex == 0 && len(key) == 1 && !unicode.IsDigit(rune(key[0])) && rune(key[0]) != '.' {
				return m, nil
			}

			switch key {
			case "up", "shift+tab":
				m.FocusIndex--
				if m.FocusIndex < 0 {
					m.FocusIndex = 2
				}
			case "down", "tab":
				m.FocusIndex++
				if m.FocusIndex > 2 {
					m.FocusIndex = 0
				}
			case "enter":
				if m.FocusIndex < 1 {
					m.FocusIndex++
				} else {
					if m.Inputs[0].Value() == "" {
						return m, nil
					}
					if m.Inputs[1].Value() == "" {
						m.Inputs[1].SetValue("0.06")
					}
					m.State = StateExecuting
					return m, Tick()
				}
			}

			var cmd tea.Cmd
			cmds := make([]tea.Cmd, 0) // Fixed initialization
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

				// FIXED: Call the real RunFullBacktest function from the app package
				summary, err := app.RunFullBacktest(m.DataFile, m.StrategyFile, m.IndicatorFile, inv, comm)
				if err != nil {
					m.Logs = append(m.Logs, fmt.Sprintf("[ERROR] %v", err))
					// In a real app, you might transition to an error state here
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
