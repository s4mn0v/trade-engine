package indicators

import (
	"github.com/s4mn0v/trade-engine/internal/domain"
)

func CalculateHybridOscillator(candles []domain.Candle, rsiLen, mfiLen, smmaLen int, rsiWeight float64) ([]float64, []float64) {
	if len(candles) == 0 {
		return nil, nil
	}

	mfiWeight := 1.0 - rsiWeight

	// 1. Extract raw data
	closes := make([]float64, len(candles))
	highs := make([]float64, len(candles))
	lows := make([]float64, len(candles))
	volumes := make([]float64, len(candles))
	for i, c := range candles {
		closes[i] = c.Close
		highs[i] = c.High
		lows[i] = c.Low
		volumes[i] = c.Volume
	}

	// 2. Calculate components
	rsi := CalculateRSI(closes, rsiLen)
	mfi := CalculateMFI(highs, lows, closes, volumes, mfiLen)

	// 3. Combine into Hybrid Oscillator
	hybridOsc := make([]float64, len(candles))
	for i := range candles {
		hybridOsc[i] = (rsi[i] * rsiWeight) + (mfi[i] * mfiWeight)
	}

	// 4. Calculate TMA (Smoothed Moving Average) on the Hybrid Oscillator
	tma := CalculateSMMA(hybridOsc, smmaLen)

	return hybridOsc, tma
}

// CalculateRSI calculates the Relative Strength Index.
func CalculateRSI(data []float64, period int) []float64 {
	rsi := make([]float64, len(data))
	if len(data) <= period {
		return rsi
	}

	var gain, loss float64
	for i := 1; i <= period; i++ {
		diff := data[i] - data[i-1]
		if diff > 0 {
			gain += diff
		} else {
			loss -= diff
		}
	}

	avgGain := gain / float64(period)
	avgLoss := loss / float64(period)

	for i := period + 1; i < len(data); i++ {
		diff := data[i] - data[i-1]
		var currentGain, currentLoss float64
		if diff > 0 {
			currentGain = diff
		} else {
			currentLoss = -diff
		}

		// Wilder's Smoothing
		avgGain = (avgGain*float64(period-1) + currentGain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + currentLoss) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}
	return rsi
}

// CalculateMFI calculates the Money Flow Index.
func CalculateMFI(highs, lows, closes, volumes []float64, period int) []float64 {
	mfi := make([]float64, len(closes))
	if len(closes) <= period {
		return mfi
	}

	typicalPrices := make([]float64, len(closes))
	for i := range closes {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3
	}

	for i := period; i < len(closes); i++ {
		posFlow := 0.0
		negFlow := 0.0

		for j := i - period + 1; j <= i; j++ {
			rawMoneyFlow := typicalPrices[j] * volumes[j]
			if typicalPrices[j] > typicalPrices[j-1] {
				posFlow += rawMoneyFlow
			} else if typicalPrices[j] < typicalPrices[j-1] {
				negFlow += rawMoneyFlow
			}
		}

		if negFlow == 0 {
			mfi[i] = 100
		} else {
			moneyRatio := posFlow / negFlow
			mfi[i] = 100 - (100 / (1 + moneyRatio))
		}
	}
	return mfi
}

// CalculateSMMA calculates the Smoothed Moving Average (Wilder's MA).
// This matches the Pine Script: smma := na(smma[1]) ? ta.sma(src, len) : (smma[1] * (len - 1) + src) / len
func CalculateSMMA(data []float64, period int) []float64 {
	smma := make([]float64, len(data))
	if len(data) <= period {
		return smma
	}

	// Initial SMA for the first valid point
	var sum float64
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	smma[period-1] = sum / float64(period)

	for i := period; i < len(data); i++ {
		smma[i] = (smma[i-1]*float64(period-1) + data[i]) / float64(period)
	}

	return smma
}
