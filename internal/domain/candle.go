package domain

import "time"

// Candle represents a single OHLCV bar from market data.
type Candle struct {
	Index       int
	Timestamp   time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
	BaseVolume  float64
	USDTVolume  float64
	QuoteVolume float64
}
