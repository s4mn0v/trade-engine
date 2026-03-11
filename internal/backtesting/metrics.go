package backtesting

import (
	"github.com/s4mn0v/trade-engine/internal/domain"
)

type Summary struct {
	InitialBalance float64
	FinalBalance   float64
	TotalNetProfit float64
	ProfitPct      float64
	WinRate        float64
	MaxDrawdown    float64
	TotalTrades    int
	Leverage       float64
	Commission     float64
}

func CalculateMetrics(trades []domain.Trade, initialBalance, commissionPercent, leverage float64) Summary {
	currentBalance := initialBalance
	maxBalance := initialBalance
	maxDrawdown := 0.0
	wins := 0
	commissionRate := commissionPercent / 100

	for _, t := range trades {
		rawPnL := t.Profit() * t.Leverage
		entryFee := t.EntryPrice * commissionRate * t.Leverage
		exitFee := t.ExitPrice * commissionRate * t.Leverage

		netPnL := rawPnL - entryFee - exitFee
		currentBalance += netPnL

		if netPnL > 0 {
			wins++
		}
		if currentBalance > maxBalance {
			maxBalance = currentBalance
		}
		if maxBalance > 0 {
			dd := (maxBalance - currentBalance) / maxBalance
			if dd > maxDrawdown {
				maxDrawdown = dd
			}
		}
		if currentBalance <= 0 {
			currentBalance = 0
			break
		}
	}

	profitPct := 0.0
	if initialBalance > 0 {
		profitPct = ((currentBalance - initialBalance) / initialBalance) * 100
	}

	winRate := 0.0
	if len(trades) > 0 {
		winRate = (float64(wins) / float64(len(trades))) * 100
	}

	return Summary{
		InitialBalance: initialBalance,
		FinalBalance:   currentBalance,
		TotalNetProfit: currentBalance - initialBalance,
		ProfitPct:      profitPct,
		WinRate:        winRate,
		MaxDrawdown:    maxDrawdown * 100,
		TotalTrades:    len(trades),
		Leverage:       leverage,
		Commission:     commissionPercent,
	}
}
