package indicator

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

//go:embed scripts/supertrend.json
var superTrendScript string

//go:embed scripts/mmri.json
var mmriScript string

type IndicatorDefinition struct {
	Name      string
	Script    map[string]interface{}
	ParseFunc func(rawJSON string) (string, error)
}

// Indicators is a global map of all possible indicator you want to register.
var Indicators = map[string]*IndicatorDefinition{}

func init() {
	var indicator = "SuperTrend"
	Indicators[indicator] = &IndicatorDefinition{
		Name:      indicator,
		Script:    mustUnmarshal(superTrendScript),
		ParseFunc: parseSuperTrend,
	}

	indicator = "MMRI"
	Indicators[indicator] = &IndicatorDefinition{
		Name:      indicator,
		Script:    mustUnmarshal(mmriScript),
		ParseFunc: parseMMRI,
	}
}

func GetIndicatorScript(indicatorName string) (map[string]interface{}, error) {
	ind, ok := Indicators[indicatorName]
	if !ok {
		return nil, fmt.Errorf("no such indicator: %s", indicatorName)
	}
	return ind.Script, nil
}

func mustUnmarshal(raw string) map[string]interface{} {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal embedded script: %v", err))
	}
	return obj
}

func ExtractStudyLabel(rawJSON string) (string, error) {
	var top duMessage
	if err := json.Unmarshal([]byte(rawJSON), &top); err != nil {
		return "", err
	}
	if top.M != "du" || len(top.P) < 2 {
		return "", errors.New("not a valid 'du' message with 2+ params")
	}

	// The second item in P is the big object: e.g. {"st1":{...}} or {"st2":{...}}
	var stMap map[string]json.RawMessage
	if err := json.Unmarshal(top.P[1], &stMap); err != nil {
		return "", err
	}

	// Return whichever key looks like st1, st2, etc.
	for k := range stMap {
		if strings.HasPrefix(k, "st") {
			return k, nil
		}
	}
	return "", errors.New("no stN label found in du message")
}
