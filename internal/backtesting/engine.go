package backtesting

import (
	"github.com/s4mn0v/trade-engine/internal/domain"
)

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
	currentBalance := e.Investment

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
				// --- LIQUIDATION LOGIC START ---
				// 1. Calculate Profit/Loss
				rawPnL := trade.Profit() * trade.Leverage

				// 2. Calculate Fees (Entry + Exit) using the Executor's CommissionRate
				// Note: CommissionRate is already decimal (e.g., 0.0006)
				entryFee := trade.EntryPrice * e.Executor.CommissionRate * trade.Leverage
				exitFee := trade.ExitPrice * e.Executor.CommissionRate * trade.Leverage

				netPnL := rawPnL - entryFee - exitFee
				currentBalance += netPnL

				completedTrades = append(completedTrades, *trade)

				// 3. Check for Bankruptcy
				if currentBalance <= 0 {
					// Stop the backtest immediately
					break
				}
				// --- LIQUIDATION LOGIC END ---
			}
		}
	}

	return completedTrades
}
