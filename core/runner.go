package core

import (
	"fmt"
	"strings"

	"github.com/cullenwatson/WhaleHunter/indicator"
	"github.com/cullenwatson/WhaleHunter/model"
	"github.com/rs/zerolog/log"
)

func RunTradingViewSession(
	symbol string,
	timeframe string,
	candlesRequested int,
	authToken string,
	indicators []string,
	candleChan chan<- []model.Candle,
	indicatorChan chan<- string,
) {
	tvClient := NewTradingViewClient(symbol, timeframe, candlesRequested, authToken, indicators)

	// Connect
	if err := tvClient.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect")
	}
	defer tvClient.Ws.Close()

	// Start the TV client
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

		// Parse out potential JSON chunks
		jsonChunks := ParseMultipleMessages(string(rawMsg))
		for _, chunk := range jsonChunks {
			log.Info().Msgf("Received: %s", chunk)

			// A) Candle data?
			if candles, err := ExtractCandles(chunk); err == nil {
				candleChan <- candles
				continue
			}

			// B) "du" study data => we try *each* indicator
			if strings.Contains(chunk, "\"du\"") {
				for _, indName := range tvClient.Indicators {
					parseFunc := indicator.Indicators[indName].ParseFunc
					result, parseErr := parseFunc(chunk)
					if parseErr == nil {
						// e.g. "MMRI => 12.3456"
						// or "SuperTrend => bullish"
						indicatorChan <- fmt.Sprintf("%s => %s", indName, result)
					} else {
						// It's normal if the chunk doesn't match this particular indicator
						log.Debug().Err(parseErr).
							Str("indicator", indName).
							Msg("Skipping parse result for chunk")
					}
				}
			}
		}
	}

	close(candleChan)
	close(indicatorChan)
}
