package main

import (
	"github.com/cullenwatson/WhaleHunter/internal/core"
	"github.com/cullenwatson/WhaleHunter/internal/model"

	mylog "github.com/cullenwatson/WhaleHunter/internal/log"
	"github.com/rs/zerolog/log"
)

func main() {
	mylog.InitLogger()
	core.LoadEnvVarsOrDie()

	authToken := core.GetTradinViewAuthToken()

	cfg := core.SessionConfig{
		Symbol:           "NVDA",
		Timeframe:        "1D",
		CandlesRequested: 5,
		AuthToken:        authToken,
		Indicators:       []string{"SuperTrend", "MMRI"},

		CandleChan:    make(chan []model.Candle),
		IndicatorChan: make(chan string),
	}

	go core.RunTradingViewSession(cfg)

	for {
		select {
		case candleBatch, ok := <-cfg.CandleChan:
			if !ok {
				log.Info().Msg("Candle channel closed. Exiting.")
				return
			}
			core.HandleCandleBatch(cfg.Symbol, candleBatch)

		case studyResult, ok := <-cfg.IndicatorChan:
			if !ok {
				log.Info().Msg("Indicator channel closed. Exiting.")
				return
			}
			log.Info().Msgf("%s %s: %s", cfg.Symbol, "Indicator", studyResult)
		}
	}
}
