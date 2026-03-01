package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	commission = 0.00006 // 0.006%
	inputState = iota
	runningState
	resultState
)

type Record struct {
	Timestamp time.Time
	Close     float64
}

type model struct {
	state       int
	fileInput   textinput.Model
	amountInput textinput.Model
	records     []Record
	balance     float64
	position    float64
	initial     float64
	err         error
}

func initialModel() model {
	fi := textinput.New()
	fi.Placeholder = "data.csv"
	fi.Focus()

	ai := textinput.New()
	ai.Placeholder = "1000"

	return model{
		state:       inputState,
		fileInput:   fi,
		amountInput: ai,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == inputState {
				if m.fileInput.Focused() {
					m.fileInput.Blur()
					m.amountInput.Focus()
				} else {
					m.initial, _ = strconv.ParseFloat(m.amountInput.Value(), 64)
					m.balance = m.initial
					m.state = runningState
					return m, m.runSimulation
				}
			}
		}
	case error:
		m.err = msg
		return m, nil
	case []Record:
		m.records = msg
		m.executeStrategy()
		m.state = resultState
		return m, nil
	}

	if m.fileInput.Focused() {
		m.fileInput, cmd = m.fileInput.Update(msg)
	} else {
		m.amountInput, cmd = m.amountInput.Update(msg)
	}
	return m, cmd
}

func (m model) runSimulation() tea.Msg {
	f, err := os.Open(m.fileInput.Value())
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	_, _ = reader.Read() // Skip headers

	var records []Record
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Header: Timestamp,Open,High,Low,Close...
		t, _ := time.Parse("2006-01-02 15:04:05-07:00", line[0])
		c, _ := strconv.ParseFloat(line[4], 64)
		records = append(records, Record{Timestamp: t, Close: c})
	}
	return records
}

func (m *model) executeStrategy() {
	lastTradeDay := -1

	for _, r := range m.records {
		day := int(r.Timestamp.Weekday())
		hour := r.Timestamp.Hour()

		// Buy Monday 00:00
		if day == int(time.Monday) && hour == 0 && m.balance > 0 && lastTradeDay != day {
			cost := m.balance * commission
			m.position = (m.balance - cost) / r.Close
			m.balance = 0
			lastTradeDay = day
		}

		// Sell Sunday 23:00 (or last available hour of Sunday)
		if day == int(time.Sunday) && hour == 22 && m.position > 0 && lastTradeDay != day {
			gross := m.position * r.Close
			m.balance = gross * (1 - commission)
			m.position = 0
			lastTradeDay = day
		}
	}

	// Final liquidation if position held
	if m.position > 0 {
		m.balance = (m.position * m.records[len(m.records)-1].Close) * (1 - commission)
		m.position = 0
	}
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	switch m.state {
	case inputState:
		return fmt.Sprintf(
			"CSV Path: %s\nInitial Investment: %s\n\n(Press Enter to switch/confirm)",
			m.fileInput.View(),
			m.amountInput.View(),
		)
	case runningState:
		return "Loading data and running simulation..."
	case resultState:
		pnl := m.balance - m.initial
		return fmt.Sprintf(
			"Results:\nInitial: %.2f USDT\nFinal:   %.2f USDT\nPnL:     %.2f USDT (%.2f%%)\n\nPress 'q' to quit",
			m.initial, m.balance, pnl, (pnl/m.initial)*100,
		)
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
