package core

import (
	"os"

	"github.com/cullenwatson/WhaleHunter/model"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func GetTradinViewAuthToken() string {
	username := os.Getenv("TRADINGVIEW_USERNAME")
	password := os.Getenv("TRADINGVIEW_PASSWORD")

	creds := model.Credentials{
		Username: username,
		Password: password,
	}

	authToken, err := SignIn(creds)
	if err != nil {
		log.Fatal().Err(err).Msg("Error signing in")
	}

	log.Info().Str("token", authToken).Msg("Auth token received")
	return authToken
}

func LoadEnvVarsOrDie() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
	log.Info().Msg(".env loaded successfully")
}

func HandleCandleBatch(symbol string, candles []model.Candle) {
	for i, c := range candles {
		log.Info().Msgf(
			"[MAIN] %s Candle #%d => Date=%s O=%.2f H=%.2f L=%.2f C=%.2f Vol=%.0f",
			symbol, i+1, c.Date.Format("2006-01-02"),
			c.Open, c.High, c.Low, c.Close, c.Volume,
		)
	}
}
