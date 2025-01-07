# Roadmap
Goal - auto place option trades based on custom indicators and whale option volume
### Phase 1 - TradingView connection
1. ~~Sign into TradingView and solve the captcha~~
* ~~Save the cookies in case restart of program~~
2. ~~Start websocket connection to TradingView~~
3. Receive and parse stock feed and custom indicator data
4. Support concurrent connection feeds to multiple stock symbols

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
    LoadEnvVars()

    authToken := GetAuthToken()

    symbol := "NVDA"
    timeframe := "1D"
    candlesRequested := 100

    client := CreateTradingViewClient(symbol, timeframe, candlesRequested, authToken)

    candleChan := make(chan []model.Candle)
    go RunTradingViewClient(client, candleChan)

    // Listen for new candles
    for candleBatch := range candleChan {
        for i, c := range candleBatch {
            log.Info().Msgf(
                "[MAIN] %s Candle #%d => Date=%s O=%.2f H=%.2f L=%.2f C=%.2f Vol=%.0f",
                symbol,
                i+1,
                c.Date.Format("2006-01-02"),
                c.Open, c.High, c.Low, c.Close, c.Volume,
            )
        }
    }

    log.Info().Msg("Candle channel closed. Exiting.")
}
```
#### current output
```bash
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #97 => Date=2024-12-31 O=138.03 H=138.07 L=133.83 C=134.29 Vol=155659211
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #98 => Date=2025-01-02 O=136.00 H=138.88 L=134.63 C=138.31 Vol=198247166
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #99 => Date=2025-01-03 O=140.01 H=144.90 L=139.73 C=144.47 Vol=229322478
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #100 => Date=2025-01-06 O=148.59 H=152.16 L=147.82 C=149.43 Vol=265377359
```