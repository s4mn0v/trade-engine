package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/s4mn0v/trade-engine/internal/domain"
)

// CSV Column Mapping based on your provided data
const (
	ColTimestamp   = 0
	ColOpen        = 1
	ColHigh        = 2
	ColLow         = 3
	ColClose       = 4
	ColBaseVolume  = 5
	ColUSDTVolume  = 6
	ColQuoteVolume = 7
)

// LoadCandlesFromCSV reads a file and converts its rows into domain objects.
func LoadCandlesFromCSV(filePath string) ([]domain.Candle, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read the header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	var candles []domain.Candle
	lineCount := 0

	// Use a loop to read row by row (better for 480,000 lines than ReadAll)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row %d: %w", lineCount+1, err)
		}

		// 1. Parse Unix Milliseconds Timestamp
		ms, err := strconv.ParseInt(record[ColTimestamp], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid timestamp format: %s", lineCount+1, record[ColTimestamp])
		}

		open, _ := strconv.ParseFloat(record[ColOpen], 64)
		high, _ := strconv.ParseFloat(record[ColHigh], 64)
		low, _ := strconv.ParseFloat(record[ColLow], 64)
		close, _ := strconv.ParseFloat(record[ColClose], 64)
		baseVol, _ := strconv.ParseFloat(record[ColBaseVolume], 64)
		usdtVol, _ := strconv.ParseFloat(record[ColUSDTVolume], 64)
		quoteVol, _ := strconv.ParseFloat(record[ColQuoteVolume], 64)

		candles = append(candles, domain.Candle{
			Index:       lineCount,
			Timestamp:   time.UnixMilli(ms),
			Open:        open,
			High:        high,
			Low:         low,
			Close:       close,
			BaseVolume:  baseVol,
			USDTVolume:  usdtVol,
			QuoteVolume: quoteVol,
		})

		lineCount++
	}

	return candles, nil
}
