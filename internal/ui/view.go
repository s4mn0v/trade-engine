package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1).MarginBottom(1)
	logStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	docStyle     = lipgloss.NewStyle().Padding(1, 2)
)

func (m Model) View() tea.View {
	if m.Quitting {
		return tea.NewView("Exiting...")
	}

	var content string

	switch m.State {
	case StateDataPicker:
		content = titleStyle.Render("1. SELECT DATA (.CSV)") + "\n" + m.Filepicker.View()

	case StateStrategyPicker:
		content = titleStyle.Render("2. SELECT STRATEGY (.GO)") + "\n" +
			"Press 's' to skip strategy selection.\n\n" + m.Filepicker.View()

	case StateIndicatorPicker:
		content = titleStyle.Render("3. SELECT INDICATOR (.GO)") + "\n" +
			"Press 's' to skip indicator selection.\n\n" + m.Filepicker.View()

	case StateConfig:
		content = titleStyle.Render("4. CONFIGURATION") + "\n" +
			fmt.Sprintf("Data: %s | Strat: %s | Ind: %s\n\n",
				formatFile(m.DataFile), formatFile(m.StrategyFile), formatFile(m.IndicatorFile))

		content += "Investment: " + m.Inputs[0].View() + "\n"
		content += "Commission: " + m.Inputs[1].View() + " %\n"

		btn := "\n[ START ENGINE ]"
		if m.FocusIndex == 2 {
			btn = "\n" + focusedStyle.Render("[ START ENGINE ]")
		}
		content += btn

	case StateExecuting:
		content = titleStyle.Render("5. ENGINE LOGS") + "\n"
		logBox := ""
		for _, l := range m.Logs {
			logBox += logStyle.Render("> "+l) + "\n"
		}
		content += lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).Width(60).Height(10).Render(logBox)

	case StateFinished:
		content = titleStyle.Render("6. BACKTEST RESULTS") + "\n" +
			lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render("SIMULATION SUCCESSFUL") + "\n\n" +
			fmt.Sprintf("Initial Capital: $%0.2f\n", m.Results.InitialCap) +
			fmt.Sprintf("Final Balance: $%0.2f\n", m.Results.FinalBalance) +
			fmt.Sprintf("Net Profit: %0.2f%%\n", m.Results.NetProfitPct) +
			fmt.Sprintf("Total Trades: %d\n\n", len(m.Results.Trades)) +
			"Full report saved to: " + focusedStyle.Render("results.txt") + "\n\n" +
			"Press Enter to Exit..."
	}

	v := tea.NewView(docStyle.Render(content))
	if m.State == StateExecuting {
		v.ProgressBar = tea.NewProgressBar(tea.ProgressBarDefault, m.ProgressPct)
	}
	return v
}

func formatFile(path string) string {
	if path == "" {
		return focusedStyle.Render("(none)")
	}
	return focusedStyle.Render(path)
}
