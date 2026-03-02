```
trade-engine/
├── cmd/
│   └── engine/
│       └── main.go              # App entrypoint (wire dependencies here only)
│
├── internal/
│   ├── app/                     # Application orchestration layer
│   │   └── engine.go
│   │
│   ├── ui/                      # Bubble Tea UI layer
│   │   ├── model.go
│   │   ├── view.go
│   │   └── components/
│   │
│   ├── domain/                  # Core business models (NO dependencies)
│   │   ├── candle.go
│   │   ├── signal.go
│   │   └── strategy.go          # Strategy interface
│   │
│   ├── indicators/              # Pure math (stateless, deterministic)
│   │   ├── rsi.go
│   │   ├── mfi.go
│   │   └── moving_avg.go
│   │
│   ├── strategies/              # Implementations of domain strategies
│   │   ├── hybrid_darsi.go
│   │   └── buy_and_hold.go
│   │
│   ├── backtest/                # Backtesting engine
│   │   ├── engine.go
│   │   ├── portfolio.go
│   │   └── metrics.go
│   │
│   └── data/
│       ├── loader.go            # CSV parsing
│       └── repository.go        # Data abstraction
│
├── data/                        # Raw datasets (NOT business logic)
│   ├── btcusdt_1h.csv
│   └── usdtd_1h.csv
│
├── pkg/                         # Optional: reusable public packages
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
