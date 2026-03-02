```
trade-engine/
├── cmd/
│   └── engine/
│       └── main.go                # Entry point (wires everything together)
│
├── internal/
│   ├── domain/                    # Core business concepts (pure)
│   │   ├── candle.go
│   │   ├── signal.go
│   │   ├── trade.go
│   │   └── strategy.go
│   │
│   ├── indicators/                # Pure math (stateless)
│   │   ├── rsi.go
│   │   ├── mfi.go
│   │   └── moving_avg.go
│   │
│   ├── strategies/                # Strategy implementations
│   │   ├── hybrid_darsi.go
│   │   └── buy_and_hold.go
│   │
│   ├── backtest/                  # Execution engine
│   │   ├── engine.go              # Runs simulation
│   │   ├── executor.go            # Handles positions & leverage logic
│   │   └── metrics.go             # (optional) PnL, winrate, drawdown
│   │
│   ├── data/                      # CSV loading
│   │   └── loader.go
│   │
│   ├── report/                    # Output layer (file writing)
│   │   └── writer.go              # Writes results.txt
│   │
│   └── ui/                        # Optional Bubble Tea UI
│       ├── model.go
│       ├── view.go
│       └── components/
│
├── data/                          # Raw datasets
│   ├── btcusdt_1h.csv
│   └── usdtd_1h.csv
│
├── results.txt                    # Generated output file
├── go.mod
├── go.sum
└── README.md
```

## Why This Structure Is Correct

This structure follows clean architecture principles: separation of concerns, single responsibility, and dependency direction toward the domain layer. Each folder exists for a specific reason and prevents logic from leaking across layers.

---

## `cmd/engine/main.go`

**Responsibility:** Application entry point and dependency wiring.

This file should:

* Load configuration (if any)
* Load market data via `internal/data`
* Instantiate a strategy from `internal/strategies`
* Call the backtest engine
* Call the report writer

It must not contain:

* Trading logic
* Indicator calculations
* File formatting logic
* Business rules

It exists only to orchestrate the application.

---

## `internal/domain/`

**Responsibility:** Core trading concepts and business models.

This is the most important layer. It defines the fundamental entities of the system.

### `candle.go`

Defines the market data model:

* Price values (open, high, low, close)
* Volume
* Timestamp
* CSV index reference

No parsing logic belongs here.

### `signal.go`

Defines trading decisions:

* Buy / Sell
* Long / Short (if represented here)
* Time or index reference

A signal is a decision, not an executed trade.

### `trade.go`

Represents a completed position:

* Entry price
* Exit price
* Entry/exit index
* Entry/exit timestamp
* Long or short
* Leverage used

This is the final result produced by the backtest engine.

### `strategy.go`

Defines the Strategy interface:

* Name()
* Generate([]Candle) []Signal

It defines behavior, not implementation.

The domain layer must not depend on:

* CSV
* File writing
* UI
* Backtest execution
* Indicators

All other layers depend on domain — never the reverse.

---

## `internal/indicators/`

**Responsibility:** Pure mathematical functions.

Contains stateless calculations such as:

* RSI
* MFI
* Moving averages

Indicators:

* Accept numeric inputs
* Return computed values
* Have no side effects
* Do not know about trades, portfolios, or files

They are utilities used by strategies.

---

## `internal/strategies/`

**Responsibility:** Trading decision logic.

Each file implements the `domain.Strategy` interface.

A strategy:

* Uses indicators
* Analyzes candles
* Produces signals

A strategy must not:

* Execute trades
* Track balances
* Handle leverage mechanics
* Write files

It only decides when to enter or exit.

---

## `internal/backtest/`

**Responsibility:** Trade execution simulation.

This layer converts signals into real simulated trades.

### `engine.go`

Coordinates the backtest process:

* Iterates over candles
* Feeds signals into execution logic
* Collects completed trades

Returns:

```
[]domain.Trade
```

### `executor.go`

Handles position mechanics:

* Opening positions
* Closing positions
* Long/short logic
* Leverage calculation
* Price selection (entry/exit)

This is where trading mechanics live.

### `metrics.go` (optional)

Calculates:

* PnL
* Win rate
* Drawdown
* Other performance statistics

Metrics are derived from completed trades.

---

## `internal/data/`

**Responsibility:** Market data loading.

### `loader.go`

* Reads CSV files
* Parses timestamps
* Converts rows into `domain.Candle`
* Assigns CSV index numbers

It must not:

* Run strategies
* Execute trades
* Calculate indicators
* Write reports

It only transforms external data into domain objects.

---

## `internal/report/`

**Responsibility:** Output formatting and file writing.

### `writer.go`

* Accepts `[]domain.Trade`
* Formats results according to your standard
* Writes `results.txt`

This layer isolates file I/O from business logic.

If later you want:

* JSON output
* Database storage
* Web API responses

You change only this layer.

---

## `internal/ui/` (Optional)

**Responsibility:** Presentation layer.

Contains Bubble Tea models and view logic.

It may:

* Display backtest results
* Let the user select strategies
* Show progress

It must not:

* Contain trading rules
* Execute trades directly
* Write files directly
* Parse CSV

UI depends on application services — not the other way around.

---

## `data/`

**Responsibility:** Raw datasets.

Contains CSV files only.

No logic belongs here.

---

## `results.txt`

Generated output file.

Not part of business logic.
Not version-controlled in most professional setups.
Represents execution output only.

---

## Architectural Direction

Dependency direction should always be:

```
UI / CLI
    ↓
Backtest Engine
    ↓
Strategy Implementations
    ↓
Domain
```

And separately:

```
CSV Loader → Domain
Report Writer ← Domain
```

Domain is the center. Everything depends on it.

---

## Why This Design Works

This structure:

* Keeps trading logic independent of infrastructure
* Makes the engine testable
* Allows strategies to be swapped easily
* Allows new output formats without touching trading logic
* Allows switching from CSV to API without rewriting strategies
* Prevents tight coupling between UI and execution logic

It scales from a CLI backtester to a production trading engine without structural rewrites.

This is why the structure is correct.
