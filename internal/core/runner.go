package core

import (
	"strings"

	"github.com/cullenwatson/WhaleHunter/internal/indicator"
	"github.com/cullenwatson/WhaleHunter/internal/model"

	"github.com/rs/zerolog/log"
)

type SessionConfig struct {
	Symbol           string
	Timeframe        string
	CandlesRequested int
	AuthToken        string
	Indicators       []string

	CandleChan    chan []model.Candle
	IndicatorChan chan string
}

func RunTradingViewSession(cfg SessionConfig) {
	tvClient := NewTradingViewClient(
		cfg.Symbol,
		cfg.Timeframe,
		cfg.CandlesRequested,
		cfg.AuthToken,
		cfg.Indicators,
	)

	if err := tvClient.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect")
	}
	defer tvClient.Ws.Close()

	if err := tvClient.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}

	// Read from the WebSocket
	for {
		_, rawMsg, err := tvClient.Ws.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("WebSocket read error")
			break
		}

		jsonChunks := ParseMultipleMessages(string(rawMsg))
		for _, chunk := range jsonChunks {
			log.Debug().Msgf("Received: %s", chunk)

			// A) Candle data?
			if candles, err := ExtractCandles(chunk); err == nil {
				cfg.CandleChan <- candles
				continue
			}

			// B) "du" study data => figure out which stN label is present
			if strings.Contains(chunk, "\"du\"") {
				stLabel, err := indicator.ExtractStudyLabel(chunk)
				if err != nil {
					log.Debug().Err(err).Msg("Skipping this 'du' chunk")
					continue
				}

				// Look up which indicator name belongs to that st label
				indName, ok := tvClient.IndicatorMap[stLabel]
				if !ok {
					log.Debug().Msgf("Received an unknown st label: %s", stLabel)
					continue
				}

				parseFunc := indicator.Indicators[indName].ParseFunc
				result, parseErr := parseFunc(chunk)
				if parseErr == nil {
					cfg.IndicatorChan <- indName + " => " + result
				} else {
					log.Debug().Err(parseErr).
						Str("indicator", indName).
						Msg("Skipping parse result for chunk")
				}
			}
		}
	}

	close(cfg.CandleChan)
	close(cfg.IndicatorChan)
}
