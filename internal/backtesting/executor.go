package backtesting

import (
	"github.com/s4mn0v/trade-engine/internal/domain" // REPLACE WITH YOUR MODULE PATH
)

// Executor handles the state of an active position.
type Executor struct {
	CommissionRate float64
	Leverage       float64
	ActiveTrade    *domain.Trade
}

// NewExecutor initializes a trade simulator.
func NewExecutor(commission, leverage float64) *Executor {
	return &Executor{
		CommissionRate: commission / 100, // Convert percent to decimal
		Leverage:       leverage,
	}
}

// OpenPosition creates a new trade record.
func (e *Executor) OpenPosition(index int, candle domain.Candle, side domain.Side) {
	if e.ActiveTrade != nil {
		return // Already in a trade
	}

	e.ActiveTrade = &domain.Trade{
		ID:             index,
		Side:           side,
		Leverage:       e.Leverage,
		EntryPrice:     candle.Close,
		EntryIndex:     index,
		EntryTimestamp: candle.Timestamp,
	}
}

// ClosePosition completes the trade record and returns it.
func (e *Executor) ClosePosition(index int, candle domain.Candle) *domain.Trade {
	if e.ActiveTrade == nil {
		return nil
	}

	trade := e.ActiveTrade
	trade.ExitPrice = candle.Close
	trade.ExitIndex = index
	trade.ExitTimestamp = candle.Timestamp

	e.ActiveTrade = nil
	return trade
}
