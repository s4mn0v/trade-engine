package domain

import "time"

type Trade struct {
	ID       int
	Side     Side
	Leverage float64

	EntryPrice float64
	ExitPrice  float64

	EntryIndex int
	ExitIndex  int

	EntryTimestamp time.Time
	ExitTimestamp  time.Time

	// Snapshots for reporting
	BalanceBefore float64
	BalanceAfter  float64
}

func (t Trade) Profit() float64 {
	if t.Side == SideLong {
		return t.ExitPrice - t.EntryPrice
	}
	return t.EntryPrice - t.ExitPrice
}
