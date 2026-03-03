package strategies

import (
	"github.com/s4mn0v/trade-engine/internal/domain"     // REPLACE WITH YOUR MODULE PATH
	"github.com/s4mn0v/trade-engine/internal/indicators" // REPLACE WITH YOUR MODULE PATH
)

type HybridDARSI struct {
	RSILen    int
	MFILen    int
	SMALen    int
	RSIWeight float64
}

// NewHybridDARSI creates a strategy with default values from the Pine Script
func NewHybridDARSI() *HybridDARSI {
	return &HybridDARSI{
		RSILen:    14,
		MFILen:    14,
		SMALen:    50,
		RSIWeight: 0.5,
	}
}

func (s *HybridDARSI) Name() string {
	return "Hybrid RSI+MFI + TMA"
}

func (s *HybridDARSI) Generate(candles []domain.Candle) []domain.Signal {
	if len(candles) < s.SMALen {
		return nil
	}

	// 1. Calculate the indicators
	hybridOsc, tma := indicators.CalculateHybridOscillator(
		candles,
		s.RSILen,
		s.MFILen,
		s.SMALen,
		s.RSIWeight,
	)

	var signals []domain.Signal

	// 2. Iterate through candles (starting after the SMA is warm)
	for i := s.SMALen; i < len(candles); i++ {
		currentOsc := hybridOsc[i]
		prevOsc := hybridOsc[i-1]
		currentTma := tma[i]
		prevTma := tma[i-1]

		// Logic: Crossover detection

		// BUY: Hybrid crosses ABOVE TMA
		if prevOsc <= prevTma && currentOsc > currentTma {
			signals = append(signals, domain.Signal{
				Index:     i,
				Timestamp: candles[i].Timestamp,
				Action:    domain.ActionBuy,
				Side:      domain.SideLong,
			})
		}

		// SELL: Hybrid crosses BELOW TMA
		if prevOsc >= prevTma && currentOsc < currentTma {
			signals = append(signals, domain.Signal{
				Index:     i,
				Timestamp: candles[i].Timestamp,
				Action:    domain.ActionSell,
				Side:      domain.SideLong, // Closing a Long
			})
		}
	}

	return signals
}
