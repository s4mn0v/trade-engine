package strategy

import (
	"fmt"
	"reflect"

	"github.com/s4mn0v/trade-engine/internal/domain"
	"github.com/s4mn0v/trade-engine/internal/indicators"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type ScriptedStrategy struct {
	fileName string
	genFunc  func([]domain.Candle) []domain.Signal
}

func (s *ScriptedStrategy) Name() string                               { return s.fileName }
func (s *ScriptedStrategy) Generate(c []domain.Candle) []domain.Signal { return s.genFunc(c) }

func LoadStrategyScript(path string) (domain.Strategy, error) {
	i := interp.New(interp.Options{})

	// 1. Load standard library
	i.Use(stdlib.Symbols)

	// 2. Export our internal domain types to the interpreter
	// This allows the user's script to use "internal/domain"
	i.Use(interp.Exports{
		"github.com/s4mn0v/trade-engine/internal/domain/domain": {
			"Candle":     reflect.ValueOf((*domain.Candle)(nil)),
			"Signal":     reflect.ValueOf((*domain.Signal)(nil)),
			"ActionBuy":  reflect.ValueOf(domain.ActionBuy),
			"ActionSell": reflect.ValueOf(domain.ActionSell),
			"SideLong":   reflect.ValueOf(domain.SideLong),
			"SideShort":  reflect.ValueOf(domain.SideShort),
		},
		"github.com/s4mn0v/trade-engine/internal/indicators/indicators": {
			"CalculateRSI":              reflect.ValueOf(indicators.CalculateRSI),
			"CalculateMFI":              reflect.ValueOf(indicators.CalculateMFI),
			"CalculateSMMA":             reflect.ValueOf(indicators.CalculateSMMA),
			"CalculateHybridOscillator": reflect.ValueOf(indicators.CalculateHybridOscillator),
		},
	})

	_, err := i.EvalPath(path)
	if err != nil {
		return nil, fmt.Errorf("interpreter error: %w", err)
	}

	v, err := i.Eval("main.Generate")
	if err != nil {
		return nil, fmt.Errorf("script missing 'func Generate' in package main: %w", err)
	}

	fn, ok := v.Interface().(func([]domain.Candle) []domain.Signal)
	if !ok {
		return nil, fmt.Errorf("Generate function signature mismatch")
	}

	return &ScriptedStrategy{
		fileName: path,
		genFunc:  fn,
	}, nil
}
