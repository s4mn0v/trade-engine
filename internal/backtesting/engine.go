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

func (e *Engine) Run() []domain.Trade {
	signals := e.Strategy.Generate(e.Candles)
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

		if signal.Action == domain.ActionBuy && e.Executor.ActiveTrade == nil {
			e.Executor.OpenPosition(i, candle, domain.SideLong)
		} else if signal.Action == domain.ActionSell && e.Executor.ActiveTrade != nil {
			trade := e.Executor.ClosePosition(i, candle)
			if trade != nil {
				// Record balance before trade resolution
				trade.BalanceBefore = currentBalance

				rawPnL := trade.Profit() * trade.Leverage
				entryFee := trade.EntryPrice * e.Executor.CommissionRate * trade.Leverage
				exitFee := trade.ExitPrice * e.Executor.CommissionRate * trade.Leverage

				netPnL := rawPnL - entryFee - exitFee
				currentBalance += netPnL

				// Record balance after trade resolution
				trade.BalanceAfter = currentBalance
				completedTrades = append(completedTrades, *trade)

				if currentBalance <= 0 {
					break
				}
			}
		}
	}
	return completedTrades
}
