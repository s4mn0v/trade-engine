package app

import (
	"time"
)

// MockTrade represents a simplified result of a trade
type MockTrade struct {
	ID        int
	Symbol    string
	Type      string // "BUY" or "SELL"
	Price     float64
	Profit    float64
	Timestamp time.Time
}

// BacktestResult contains the overall summary
type BacktestResult struct {
	Trades       []MockTrade
	InitialCap   float64
	FinalBalance float64
	NetProfitPct float64
}

// RunBacktest simulates the execution.
// In a real app, this would be a long-running process.
func RunBacktest(dataFile, strategy, indicator string, investment, commission float64) BacktestResult {
	// Simulated delay is handled by the UI Tick for UX purposes.
	// This function returns the final static data.

	mockTrades := []MockTrade{
		{ID: 1, Symbol: "BTC/USDT", Type: "BUY", Price: 42000.00, Timestamp: time.Now()},
		{ID: 2, Symbol: "BTC/USDT", Type: "SELL", Price: 43500.00, Profit: 150.0, Timestamp: time.Now().Add(time.Hour)},
		{ID: 3, Symbol: "BTC/USDT", Type: "BUY", Price: 41000.00, Timestamp: time.Now().Add(time.Hour * 2)},
		{ID: 4, Symbol: "BTC/USDT", Type: "SELL", Price: 44000.00, Profit: 300.0, Timestamp: time.Now().Add(time.Hour * 3)},
	}

	return BacktestResult{
		Trades:       mockTrades,
		InitialCap:   investment,
		FinalBalance: investment + 450.0,
		NetProfitPct: 4.5,
	}
}
