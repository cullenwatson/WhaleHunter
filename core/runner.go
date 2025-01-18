package core

import (
	"fmt"
	"strings"

	"github.com/cullenwatson/WhaleHunter/indicator"
	"github.com/cullenwatson/WhaleHunter/model"
	"github.com/rs/zerolog/log"
)

func RunTradingViewSession(
	symbol, timeframe string,
	candlesRequested int,
	authToken string,
	indicatorName string,
	candleChan chan<- []model.Candle,
	indicatorChan chan<- string,
) {
	tvClient := NewTradingViewClient(symbol, timeframe, candlesRequested, authToken, indicatorName)

	// Connect
	if err := tvClient.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect")
	}
	defer tvClient.Ws.Close()

	// Start the TV client (create_series, create_study, etc.)
	if err := tvClient.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}

	// Read from the WebSocket
	for {
		_, rawMsg, err := tvClient.Ws.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("Websocket read error")
			break
		}
		log.Info().Msgf("Received: %s", string(rawMsg))

		// Split out all possible JSON chunks
		jsonChunks := ParseMultipleMessages(string(rawMsg))
		for _, chunk := range jsonChunks {
			// A) Try timescale_update → handle as candle data
			candles, err := ExtractCandles(chunk)
			if err == nil {
				candleChan <- candles
				continue
			}

			// B) Otherwise, if it’s a "du" message, parse with the chosen indicator
			if strings.Contains(chunk, "\"du\"") {
				parseFunc := indicator.Indicators[indicatorName].ParseFunc
				result, parseErr := parseFunc(chunk)
				if parseErr != nil {
					log.Debug().Err(parseErr).Str("chunk", chunk).
						Msg("Skipping or error in parse for du chunk")
					continue
				}
				// e.g. "supertrend => bullish"
				indicatorChan <- fmt.Sprintf("%s => %s", indicatorName, result)
			}
		}
	}

	close(candleChan)
	close(indicatorChan)
}
