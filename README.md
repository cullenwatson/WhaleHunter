# Roadmap

Goal - auto place option trades based on custom indicators and whale option volume

### Phase 1 - TradingView connection

1. ~~Sign into TradingView and solve the captcha~~

* ~~Save the cookies in case restart of program~~

2. ~~Start websocket connection to TradingView~~
3. ~~Receive and parse stock feed~~
4. ~~Receive custom indicator data~~
5. Support concurrent connection feeds to multiple stock symbols

#### why not other repos for Phase 1

- they dont auto sign into acc
- they dont give custom indicator values that are only available in your acc (private scripts)

### Phase 2 - UnusualWhales

1. Connect to unusualwhales or cheddarflow
2. Allow custom tracking / filters

### Phase 3 - Merge phase 1 & 2 for confluence

1. Connect the TradingView and unusualwhale microservice
2. Alert if both are confluent

### Phase 4 - Auto trade

1. ThinkOrSwim microservice
2. Auto purchase the option
3. Highly configurable for the trades desired
4. set take profit / stop loss and monitor win rate
5. View open positions

#### current usage

```go
func main() {
    cfg := core.SessionConfig{
        Symbol:              "NVDA",
        Timeframe:             "1D",
        CandlesRequested:         5,
        AuthToken:        authToken,
        Indicators:       []string{"SuperTrend", "MMRI"},
    }
    go core.RunTradingViewSession(cfg)
}
```

#### current output
```bash
NVDA Candle:  1 => Date=2025-01-13 O=129.99 H=133.49 L=129.51 C=133.23 Vol=204808914
NVDA Candle:  2 => Date=2025-01-14 O=136.05 H=136.38 L=130.05 C=131.76 Vol=195590485
NVDA Candle:  3 => Date=2025-01-15 O=133.65 H=136.45 L=131.29 C=136.24 Vol=185217338
NVDA Candle:  4 => Date=2025-01-16 O=138.64 H=138.75 L=133.49 C=133.57 Vol=209235583
NVDA Candle:  5 => Date=2025-01-17 O=136.69 H=138.50 L=135.46 C=137.71 Vol=201188760
NVDA Indicator: SuperTrend => bearish
NVDA Indicator: MMRI => 312.19
```