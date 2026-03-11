package main

import (
	"fmt"

	"github.com/s4mn0v/trade-engine/internal/domain"
)

func Generate(candles []domain.Candle) []domain.Signal {
	var signals []domain.Signal

	for i := 1; i < len(candles); i++ {
		curr, prev := candles[i], candles[i-1]

		if curr.Close > prev.High && curr.USDTVolume > prev.USDTVolume {
			// Define reason dynamically based on market data
			reason := fmt.Sprintf("Price (%0.2f) broke High (%0.2f) with Vol: %0.0f",
				curr.Close, prev.High, curr.USDTVolume)

			signals = append(signals, domain.Signal{
				Index:     curr.Index,
				Timestamp: curr.Timestamp,
				Action:    domain.ActionBuy,
				Side:      domain.SideLong,
				Reason:    reason,
			})
		}

		if curr.Close < prev.Low {
			signals = append(signals, domain.Signal{
				Index:     curr.Index,
				Timestamp: curr.Timestamp,
				Action:    domain.ActionSell,
				Side:      domain.SideLong,
			})
		}
	}
	return signals
}
