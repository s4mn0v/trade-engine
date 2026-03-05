package main

// The import path here must match the key we defined in interpreter.go
// Note the extra /domain at the end to match the Export key
import (
	"github.com/s4mn0v/trade-engine/internal/domain"
)

func Generate(candles []domain.Candle) []domain.Signal {
	var signals []domain.Signal

	for i := 1; i < len(candles); i++ {
		current := candles[i]
		previous := candles[i-1]

		// Logic: Price breaks high with volume confirmation
		if current.Close > previous.High && current.USDTVolume > previous.USDTVolume {
			signals = append(signals, domain.Signal{
				Index:     current.Index,
				Timestamp: current.Timestamp,
				Action:    domain.ActionBuy,
				Side:      domain.SideLong,
			})
		}

		// Exit Logic: Price drops below previous low
		if current.Close < previous.Low {
			signals = append(signals, domain.Signal{
				Index:     current.Index,
				Timestamp: current.Timestamp,
				Action:    domain.ActionSell,
				Side:      domain.SideLong,
			})
		}
	}

	return signals
}
