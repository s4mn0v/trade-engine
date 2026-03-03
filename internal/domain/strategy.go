package domain

// Strategy defines the behavior required to generate trading signals.
type Strategy interface {
	// Name returns the identifier of the strategy (e.g., "RSI-Crossover")
	Name() string

	// Generate analyzes a slice of candles and returns a slice of signals.
	// This is where the core logic of a strategy will eventually live.
	Generate(candles []Candle) []Signal
}
