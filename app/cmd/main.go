package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
	"tw-scanner/core"
	"tw-scanner/models"
)

func LoadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
}

func GetAuthToken() string {
	username := os.Getenv("TRADINGVIEW_USERNAME")
	password := os.Getenv("TRADINGVIEW_PASSWORD")

	creds := models.Credentials{
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

func RunTradingViewClient(client *core.TradingViewClient) {
	if err := client.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect")
	}
	defer client.Ws.Close()

	if err := client.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}

	for {
		_, message, err := client.Ws.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("Websocket read error")
			break
		}
		log.Info().Msgf("Received: %s", message)
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

	symbol := "AAPL"
	timeframe := "1D"
	candles := 100

	client := CreateTradingViewClient(symbol, timeframe, candles, authToken)

	RunTradingViewClient(client)
}
