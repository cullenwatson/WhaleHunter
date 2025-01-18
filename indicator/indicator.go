package indicator

import (
	_ "embed"
	"encoding/json"
	"fmt"
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
