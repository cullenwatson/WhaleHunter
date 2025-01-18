package indicator

import (
	"encoding/json"
	"errors"
	"fmt"
)

func parseMMRI(rawJSON string) (string, error) {
	// 1) Unmarshal the top-level "du" message
	var top duMessage
	if err := json.Unmarshal([]byte(rawJSON), &top); err != nil {
		return "", fmt.Errorf("failed to unmarshal top-level du: %v", err)
	}
	if top.M != "du" {
		return "", errors.New("message is not a 'du' study update")
	}

	// 2) We expect at least two items in P: the "du" key plus the actual st1 data
	if len(top.P) < 2 {
		return "", errors.New("unexpected structure: no st1 payload found in 'du' message")
	}

	// 3) Decode the st1 payload
	var payload studyPayload
	if err := json.Unmarshal(top.P[1], &payload); err != nil {
		return "", fmt.Errorf("failed to unmarshal st1 payload: %v", err)
	}
	stArray := payload.St1.St
	if len(stArray) == 0 {
		return "", errors.New("no st data found in du message (empty 'St' array)")
	}

	// 4) Grab the LAST bar => latest data
	lastBar := stArray[len(stArray)-1]
	values := lastBar.V
	if len(values) < 2 {
		return "", errors.New("not enough fields in the last bar")
	}

	// 5) The first real numeric (non-1e+100) after index 0 is your MMRI
	for i := 1; i < len(values); i++ {
		val := values[i]
		// 1e+100 is effectively a placeholder meaning "no real value",
		// so skip it and any nonsense 0.0 if needed
		if val != 1e+100 {
			mmriValue := val
			// Convert to string with a reasonable format
			return fmt.Sprintf("%.4f", mmriValue), nil
		}
	}

	// If we somehow never found a valid value:
	return "", errors.New("could not find a valid MMRI (non-1e+100) in last bar")
}
