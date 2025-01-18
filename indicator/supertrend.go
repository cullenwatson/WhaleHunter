package indicator

import (
	"encoding/json"
	"errors"
	"fmt"
)

func parseSuperTrend(rawJSON string) (string, error) {
	var top duMessage
	if err := json.Unmarshal([]byte(rawJSON), &top); err != nil {
		return "", fmt.Errorf("failed to unmarshal top-level du: %v", err)
	}
	if top.M != "du" {
		return "", errors.New("message is not a study update (du)")
	}

	// expect at least two items in P: "du" + "st1" data
	if len(top.P) < 2 {
		return "", errors.New("unexpected structure in du: no st1 payload found")
	}

	var payload studyPayload
	if err := json.Unmarshal(top.P[1], &payload); err != nil {
		return "", fmt.Errorf("failed to unmarshal st1 payload: %v", err)
	}
	stArray := payload.St1.St
	if len(stArray) == 0 {
		return "", errors.New("no st data found in du message")
	}

	// Grab the *latest* bar in stArray
	lastBar := stArray[len(stArray)-1]
	values := lastBar.V
	if len(values) < 9 {
		return "", errors.New("supertrend array doesn't have enough fields")
	}

	// Check if ANY 2.0 appears (bearish overrides bullish)
	hasBearish := false
	hasBullish := false
	for _, val := range values {
		switch val {
		case 2.0:
			hasBearish = true
		case 1.0:
			hasBullish = true
		}
	}

	if hasBearish {
		return "bearish", nil
	}
	if hasBullish {
		return "bullish", nil
	}
	return "", fmt.Errorf("no SuperTrend (1.0 or 2.0) found in last bar: %v", values)
}
