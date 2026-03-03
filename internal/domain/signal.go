package domain

import "time"

type (
	Action string
	Side   string
)

const (
	ActionBuy  Action = "BUY"
	ActionSell Action = "SELL"

	SideLong  Side = "LONG"
	SideShort Side = "SHORT"
)

// Signal represents a decision made by a strategy at a specific point in time.
type Signal struct {
	Index     int       // The candle index where the signal was generated
	Timestamp time.Time // When the signal occurred
	Action    Action    // BUY or SELL
	Side      Side      // LONG or SHORT
}
