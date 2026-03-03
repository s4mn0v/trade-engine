package domain

import "time"

// Candle represents a single OHLCV bar from market data.
type Candle struct {
	Index     int       // Position in the original CSV/Data source
	Timestamp time.Time // Opening time of the candle
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}
