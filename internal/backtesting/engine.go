package backtesting

import (
	"github.com/s4mn0v/trade-engine/internal/domain" // REPLACE WITH YOUR MODULE PATH
)

// Engine runs the backtest loop.
type Engine struct {
	Strategy   domain.Strategy
	Candles    []domain.Candle
	Executor   *Executor
	Investment float64
}

// Run executes the strategy against the data.
func (e *Engine) Run() []domain.Trade {
	// 1. Generate all signals from the strategy
	signals := e.Strategy.Generate(e.Candles)

	// Create a map for quick lookup: Index -> Signal
	signalMap := make(map[int]domain.Signal)
	for _, s := range signals {
		signalMap[s.Index] = s
	}

	var completedTrades []domain.Trade

	// 2. Iterate through market timeline
	for i, candle := range e.Candles {
		signal, hasSignal := signalMap[i]

		if !hasSignal {
			continue
		}

		// Logic for Long positions
		if signal.Action == domain.ActionBuy && e.Executor.ActiveTrade == nil {
			e.Executor.OpenPosition(i, candle, domain.SideLong)
		} else if signal.Action == domain.ActionSell && e.Executor.ActiveTrade != nil {
			trade := e.Executor.ClosePosition(i, candle)
			if trade != nil {
				completedTrades = append(completedTrades, *trade)
			}
		}
	}

	return completedTrades
}
