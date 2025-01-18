package indicator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func parseSuperTrend(rawJSON string) (string, error) {
	var top duMessage
	if err := json.Unmarshal([]byte(rawJSON), &top); err != nil {
		return "", fmt.Errorf("failed to unmarshal top-level du: %v", err)
	}
	if top.M != "du" || len(top.P) < 2 {
		return "", errors.New("not a 'du' message")
	}

	var stMap map[string]st1Payload
	if err := json.Unmarshal(top.P[1], &stMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal du payload: %v", err)
	}

	var stPayload st1Payload
	foundAny := false
	for k, v := range stMap {
		if strings.HasPrefix(k, "st") {
			stPayload = v
			foundAny = true
			break
		}
	}
	if !foundAny {
		return "", errors.New("no stN key found in du payload")
	}
	if len(stPayload.St) == 0 {
		return "", errors.New("empty 'st' array in supertrend payload")
	}

	lastBar := stPayload.St[len(stPayload.St)-1]
	values := lastBar.V
	if len(values) < 9 {
		return "", errors.New("supertrend array doesn't have enough fields")
	}

	// Basic example: 1.0 => bullish, 2.0 => bearish
	hasBearish := false
	hasBullish := false
	for _, val := range values {
		if val == 2.0 {
			hasBearish = true
		} else if val == 1.0 {
			hasBullish = true
		}
	}
	switch {
	case hasBearish:
		return "bearish", nil
	case hasBullish:
		return "bullish", nil
	default:
		return "", fmt.Errorf("no SuperTrend signal found in last bar: %v", values)
	}
}
