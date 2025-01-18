package core

import (
	"fmt"
	"os"

	capsolver_go "github.com/capsolver/capsolver-go"
	"github.com/rs/zerolog/log"
)

func solveCaptcha() (string, error) {
	apiKey := os.Getenv("CAPSOLVER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("CAPSOLVER_API_KEY environment variable is not set")
	}

	log.Info().Msg("Starting captcha solving process...")
	capSolver := capsolver_go.CapSolver{ApiKey: apiKey}
	solution, err := capSolver.Solve(
		map[string]any{
			"type":       "ReCaptchaV2taskProxyLess",
			"websiteURL": "https://www.tradingview.com",
			"websiteKey": "6Lcqv24UAAAAAIvkElDvwPxD0R8scDnMpizaBcHQ",
		})

	if err != nil {
		log.Error().Err(err).Msg("Failed to solve captcha")
		return "", fmt.Errorf("failed to solve captcha: %v", err)
	}

	gResponse := solution.Solution.GRecaptchaResponse
	log.Info().Msgf("Successfully solved captcha: %s", gResponse[:20])
	return gResponse, nil
}
