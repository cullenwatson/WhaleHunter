# Roadmap
Goal - auto place option trades based on custom indicators and whale option volume
1. ~~Sign into TradingView and solve the captcha~~
* ~~Save the cookies in case restart of program~~
2. ~~Start websocket connection to TradingView~~
3. Load desired chart template and indicators
4. Load multiple desired stock symbols and timeframe

5. Receive and parse stock feed and custom indicator data
6. Support concurrent connection feeds to multiple stock symbols


## why not other repos

- they dont auto sign into acc
- they dont give custom indicator values that are only available in your acc (private scripts)


usage
```go
func main() {
	LoadEnvVars()

	authToken := GetAuthToken()

	symbol := "AAPL"
	timeframe := "1D"
	candles := 100

	client := CreateTradingViewClient(symbol, timeframe, candles, authToken)

	RunTradingViewClient(client)
}
```
current output
```bash
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #97 => Date=2024-12-31 O=138.03 H=138.07 L=133.83 C=134.29 Vol=155659211
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #98 => Date=2025-01-02 O=136.00 H=138.88 L=134.63 C=138.31 Vol=198247166
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #99 => Date=2025-01-03 O=140.01 H=144.90 L=139.73 C=144.47 Vol=229322478
2025-01-07T06:35:24-06:00 | INFO  | [MAIN] NVDA Candle #100 => Date=2025-01-06 O=148.59 H=152.16 L=147.82 C=149.43 Vol=265377359
```