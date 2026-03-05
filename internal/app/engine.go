package app

import (
	"github.com/s4mn0v/trade-engine/internal/backtesting"
	"github.com/s4mn0v/trade-engine/internal/data"
	"github.com/s4mn0v/trade-engine/internal/report"
	"github.com/s4mn0v/trade-engine/internal/strategies"
)

// RunFullBacktest performs the end-to-end backtest process.
func RunFullBacktest(dataPath, stratPath, indPath string, investment, commission, leverage float64) (backtesting.Summary, error) {
	// 1. Load Data (Infrastructure Layer)
	candles, err := data.LoadCandlesFromCSV(dataPath)
	if err != nil {
		return backtesting.Summary{}, err
	}

	// 2. Instantiate Strategy (Strategy Layer)
	// For now, we use the HybridDARSI. In the future, stratPath could determine which one to load.
	strat := strategies.NewHybridDARSI()

	// 3. Prepare Execution (Backtesting Layer)
	// We use 1.0 as default leverage for now.
	executor := backtesting.NewExecutor(commission, leverage)
	engine := backtesting.Engine{
		Strategy:   strat,
		Candles:    candles,
		Executor:   executor,
		Investment: investment,
	}

	// 4. Run Simulation
	trades := engine.Run()

	// 5. Calculate Performance (Metrics)
	summary := backtesting.CalculateMetrics(trades, investment, commission)

	// 6. Persist Results (Report Layer)
	err = report.ExportResults("results.txt", trades, summary)
	if err != nil {
		return summary, err
	}

	return summary, nil
}
