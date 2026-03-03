package backtesting

import (
	"github.com/s4mn0v/trade-engine/internal/domain" // REPLACE WITH YOUR MODULE PATH
)

type Summary struct {
	InitialBalance float64
	FinalBalance   float64
	TotalNetProfit float64
	ProfitPct      float64
	WinRate        float64
	MaxDrawdown    float64
	TotalTrades    int
}

// CalculateMetrics transforms trades into performance data.
func CalculateMetrics(trades []domain.Trade, initialBalance, commissionRate float64) Summary {
	currentBalance := initialBalance
	maxBalance := initialBalance
	maxDrawdown := 0.0
	wins := 0

	for _, t := range trades {
		// Calculate raw PnL based on side and leverage
		rawPnL := t.Profit() * t.Leverage

		// Entry and Exit Commissions (Total 2 trades per position)
		entryFee := t.EntryPrice * (commissionRate / 100) * t.Leverage
		exitFee := t.ExitPrice * (commissionRate / 100) * t.Leverage

		netPnL := rawPnL - entryFee - exitFee
		currentBalance += netPnL

		if netPnL > 0 {
			wins++
		}

		// Drawdown tracking
		if currentBalance > maxBalance {
			maxBalance = currentBalance
		}
		dd := (maxBalance - currentBalance) / maxBalance
		if dd > maxDrawdown {
			maxDrawdown = dd
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
		MaxDrawdown:    maxDrawdown * 100, // As percentage
		TotalTrades:    len(trades),
	}
}
