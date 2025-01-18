package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cullenwatson/WhaleHunter/internal/model"
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
			"%s Candle: %4d => Date=%s O=%.2f H=%.2f L=%.2f C=%.2f Vol=%.0f",
			symbol, i+1, c.Date.Format("2006-01-02"),
			c.Open, c.High, c.Low, c.Close, c.Volume,
		)
	}
}

// createMessage replicates "~m~<len>~m~<json>" format
func createMessage(functionName string, paramList interface{}) string {
	payload := map[string]interface{}{
		"m": functionName,
		"p": paramList,
	}
	jsonBytes, _ := json.Marshal(payload)
	jsonStr := string(jsonBytes)
	return fmt.Sprintf("~m~%d~m~%s", len(jsonStr), jsonStr)
}

// generateRandomString helper for session generation
func generateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generateChartSession returns something like "cs_abcdef..."
func generateChartSession() string {
	return "cs_" + generateRandomString(12)
}
