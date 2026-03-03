package main

import (
	"fmt"

	"github.com/s4mn0v/trade-engine/internal/data"
)

// Example verification snippet
func main() {
	candles, err := data.LoadCandlesFromCSV("/home/zam/code/data.csv")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully loaded %d candles.\n", len(candles))
	fmt.Printf("First Candle: %+v\n", candles[0])
}
