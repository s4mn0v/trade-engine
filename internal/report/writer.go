package report

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/s4mn0v/trade-engine/internal/backtesting"
	"github.com/s4mn0v/trade-engine/internal/domain" // REPLACE WITH YOUR MODULE PATH
)

// ExportResults creates the results.txt file with a detailed breakdown of the backtest.
func ExportResults(filename string, trades []domain.Trade, summary backtesting.Summary) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	// 1. Write Header & Summary
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
	fmt.Fprintln(file, "")

	// 2. Write Detailed Trade Log
	fmt.Fprintln(file, "DETAILED TRADE LOG:")
	fmt.Fprintln(file, "--------------------------------------------------------------------------------")

	// Using tabwriter for clean column alignment
	w := tabwriter.NewWriter(file, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tSide\tEntry Time\tEntry Index\tEntry Price\tExit Time\tExit Index\tExit Price\tPnL")

	for i, t := range trades {
		pnl := t.Profit() * t.Leverage
		pnlStr := fmt.Sprintf("%+0.2f", pnl)

		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%0.2f\t%s\t%d\t%0.2f\t%s\n",
			i+1,
			t.Side,
			t.EntryTimestamp.Format("2006-01-02 15:04"),
			t.EntryIndex,
			t.EntryPrice,
			t.ExitTimestamp.Format("2006-01-02 15:04"),
			t.ExitIndex,
			t.ExitPrice,
			pnlStr,
		)
	}
	w.Flush()

	fmt.Fprintln(file, "--------------------------------------------------------------------------------")
	fmt.Fprintln(file, "End of Report")

	return nil
}
