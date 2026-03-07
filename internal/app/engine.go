package app

import (
	"fmt"

	"github.com/s4mn0v/trade-engine/internal/backtesting"
	"github.com/s4mn0v/trade-engine/internal/data"
	"github.com/s4mn0v/trade-engine/internal/domain"
	"github.com/s4mn0v/trade-engine/internal/report"
	"github.com/s4mn0v/trade-engine/internal/strategy"
)

// RunFullBacktest performs the end-to-end backtest process.
func RunFullBacktest(dataPath, stratPath, indPath string, investment, commission, leverage float64) (backtesting.Summary, error) {
	// 1. Load Data (Infrastructure Layer)
	candles, err := data.LoadCandlesFromCSV(dataPath)
	if err != nil {
		return backtesting.Summary{}, err
	}

	// 2. Instantiate Strategy (Strategy Layer)
	// strat := strategies.NewHybridDARSI()
	var strat domain.Strategy
	if stratPath != "" {
		strat, err = strategy.LoadStrategyScript(stratPath)
		if err != nil {
			return backtesting.Summary{}, err
		}
	} else {
		fmt.Errorf("You need to select an strategy...", err)
	}

	// 3. Prepare Execution (Backtesting Layer)
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
	err = report.ExportResults(trades, summary)
	if err != nil {
		return summary, fmt.Errorf("Failed to save results: %w", err)
	}

	return summary, nil
}
