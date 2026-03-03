package strategies

import (
	"github.com/s4mn0v/trade-engine/internal/domain" // REPLACE WITH YOUR MODULE PATH
)

type BuyAndHold struct{}

func (s *BuyAndHold) Name() string {
	return "Buy and Hold"
}

func (s *BuyAndHold) Generate(candles []domain.Candle) []domain.Signal {
	if len(candles) == 0 {
		return nil
	}

	// Just one signal at the very beginning
	return []domain.Signal{
		{
			Index:     0,
			Timestamp: candles[0].Timestamp,
			Action:    domain.ActionBuy,
			Side:      domain.SideLong,
		},
	}
}
