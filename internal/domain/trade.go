package domain

import "time"

// Trade represents a completed trading position.
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
}

// Profit calculates the raw difference between exit and entry.
// (Note: Logic is kept minimal as this is a domain model property).
func (t Trade) Profit() float64 {
	if t.Side == SideLong {
		return t.ExitPrice - t.EntryPrice
	}
	return t.EntryPrice - t.ExitPrice
}
