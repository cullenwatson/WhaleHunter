package main

import (
	"os"
	"time"

	"github.com/cullenwatson/WhaleHunter/core"
	"github.com/cullenwatson/WhaleHunter/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

func main() {
	core.LoadEnvVarsOrDie()

	authToken := core.GetTradinViewAuthToken()

	symbol := "NKE"
	timeframe := "1D"
	candlesRequested := 100

	indicators := []string{"SuperTrend", "MMRI"}

	candleChan := make(chan []model.Candle)
	indicatorChan := make(chan string)

	go core.RunTradingViewSession(symbol, timeframe, candlesRequested, authToken, indicators, candleChan, indicatorChan)

	for {
		select {
		case candleBatch, ok := <-candleChan:
			if !ok {
				log.Info().Msg("Candle channel closed. Exiting.")
				return
			}
			core.HandleCandleBatch(symbol, candleBatch)

		case studyResult, ok := <-indicatorChan:
			if !ok {
				log.Info().Msg("Indicator channel closed. Exiting.")
				return
			}
			log.Info().Msgf("[MAIN] %s Study Update: %s", symbol, studyResult)
		}
	}
}
