package main

import (
	"fmt"
	"github.com/cullenwatson/WhaleHunter/core"
	"github.com/cullenwatson/WhaleHunter/model"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

func LoadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
}

func GetAuthToken() string {
	username := os.Getenv("TRADINGVIEW_USERNAME")
	password := os.Getenv("TRADINGVIEW_PASSWORD")

	creds := model.Credentials{
		Username: username,
		Password: password,
	}

	authToken, err := core.SignIn(creds)
	if err != nil {
		log.Fatal().Err(err).Msg("Error signing in")
	}

	log.Info().Str("token", authToken).Msg("Auth token received")
	return authToken
}

func CreateTradingViewClient(symbol, timeframe string, candles int, authToken string) *core.TradingViewClient {
	client := core.NewTradingViewClient(symbol, timeframe, candles, authToken)
	return client
}

// RunTradingViewClient now accepts a channel, so it can send candle slices back to main.
func RunTradingViewClient(client *core.TradingViewClient, candleChan chan<- []model.Candle) {
	// Connect
	if err := client.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect")
	}
	defer client.Ws.Close()
	defer close(candleChan)

	if err := client.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}

	// The main receive loop
	for {
		_, rawMsg, err := client.Ws.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("Websocket read error")
			break
		}
		log.Info().Msgf("Received: %s", string(rawMsg))

		// 1) Split out all possible JSON chunks
		jsonChunks := core.ParseMultipleMessages(string(rawMsg))
		for _, chunk := range jsonChunks {
			// 2) Attempt to parse as timescale_update
			candles, err := core.ExtractCandles(chunk)
			if err != nil {
				log.Debug().Err(err).Str("chunk", chunk).Msg("Skipping non-timescale or partial chunk")
				continue
			}
			candleChan <- candles
		}
	}
}

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

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
			log.Info().Msgf("[MAIN] %s Candle #%d => Date=%s O=%.2f H=%.2f L=%.2f C=%.2f Vol=%.0f",
				symbol,
				i+1,
				c.Date.Format("2006-01-02"),
				c.Open, c.High, c.Low, c.Close, c.Volume,
			)
		}
	}

	log.Info().Msg("Candle channel closed. Exiting.")
}
