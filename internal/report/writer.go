package report

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/s4mn0v/trade-engine/internal/backtesting"
	"github.com/s4mn0v/trade-engine/internal/domain"
)

func ExportResults(trades []domain.Trade, summary backtesting.Summary) error {
	if err := exportTXT("results.txt", trades, summary); err != nil {
		return err
	}
	return exportCSV("results.csv", trades, summary)
}

func exportTXT(filename string, trades []domain.Trade, summary backtesting.Summary) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "================================================================================")
	fmt.Fprintln(file, "                          STRATEGY BACKTEST REPORT                              ")
	fmt.Fprintln(file, "================================================================================")
	fmt.Fprintf(file, "Initial Balance:  $%10.2f\n", summary.InitialBalance)
	fmt.Fprintf(file, "Final Balance:    $%10.2f\n", summary.FinalBalance)
	fmt.Fprintf(file, "Net Profit/Loss:  $%10.2f (%0.2f%%)\n", summary.TotalNetProfit, summary.ProfitPct)
	fmt.Fprintf(file, "Win Rate:         %0.2f%%\n", summary.WinRate)
	fmt.Fprintf(file, "Max Drawdown:     %0.2f%%\n", summary.MaxDrawdown)
	fmt.Fprintf(file, "Total Trades:     %d\n", summary.TotalTrades)
	fmt.Fprintln(file, "================================================================================")

	w := tabwriter.NewWriter(file, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "\nID\tSide\tEntry Time\tEntry Price\tExit Time\tExit Price\tPnL\tBalance")
	for i, t := range trades {
		pnl := (t.BalanceAfter - t.BalanceBefore)
		fmt.Fprintf(w, "%d\t%s\t%s\t%0.2f\t%s\t%0.2f\t%+0.2f\t%0.2f\n",
			i+1, t.Side, t.EntryTimestamp.Format("01-02 15:04"), t.EntryPrice,
			t.ExitTimestamp.Format("01-02 15:04"), t.ExitPrice, pnl, t.BalanceAfter)
	}
	w.Flush()
	return nil
}

func exportCSV(filename string, trades []domain.Trade, summary backtesting.Summary) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"ID", "Side", "Entry Time", "Exit Time", "Entry Price", "Exit Price", "PnL", "Balance Before", "Balance After"})

	for i, t := range trades {
		pnl := t.BalanceAfter - t.BalanceBefore
		writer.Write([]string{
			strconv.Itoa(i + 1),
			string(t.Side),
			t.EntryTimestamp.Format("2006-01-02 15:04:05"),
			t.ExitTimestamp.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%.2f", t.EntryPrice),
			fmt.Sprintf("%.2f", t.ExitPrice),
			fmt.Sprintf("%.2f", pnl),
			fmt.Sprintf("%.2f", t.BalanceBefore),
			fmt.Sprintf("%.2f", t.BalanceAfter),
		})
	}

	// Summary Footer
	writer.Write([]string{""})
	writer.Write([]string{"SUMMARY STATISTICS"})
	writer.Write([]string{"Initial Balance", fmt.Sprintf("%.2f", summary.InitialBalance)})
	writer.Write([]string{"Final Balance", fmt.Sprintf("%.2f", summary.FinalBalance)})
	writer.Write([]string{"Net Profit %", fmt.Sprintf("%.2f%%", summary.ProfitPct)})
	writer.Write([]string{"Win Rate", fmt.Sprintf("%.2f%%", summary.WinRate)})
	writer.Write([]string{"Max Drawdown", fmt.Sprintf("%.2f%%", summary.MaxDrawdown)})
	writer.Write([]string{"Total Trades", strconv.Itoa(summary.TotalTrades)})

	return nil
}
