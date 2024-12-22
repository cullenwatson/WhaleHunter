package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"tw-scanner/tradingview"
)

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
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	username := os.Getenv("TRADINGVIEW_USERNAME")
	password := os.Getenv("TRADINGVIEW_PASSWORD")

	if username == "" || password == "" {
		log.Fatal().Msg("Missing required environment variables TRADINGVIEW_USERNAME and/or TRADINGVIEW_PASSWORD")
	}

	creds := tradingview.Credentials{
		Username: username,
		Password: password,
	}

	authToken, err := tradingview.SignIn(creds)
	if err != nil {
		log.Fatal().Err(err).Msg("Error signing in")
	}

	log.Info().Msg("Successfully signed in to TradingView!")
	log.Info().Str("token", authToken).Msg("Auth token received")
}
