package domain

import "time"

type (
	Action string
	Side   string
)

const (
	ActionBuy  Action = "BUY"
	ActionSell Action = "SELL"
	SideLong   Side   = "LONG"
	SideShort  Side   = "SHORT"
)

type Signal struct {
	Index     int
	Timestamp time.Time
	Action    Action
	Side      Side
	Reason    string // Dynamic explanation defined by the strategy
}
