package indicator

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed scripts/super_trend.json
var superTrendRaw string

//go:embed scripts/rsi.json
var rsiRaw string

type IndicatorDefinition struct {
	Name      string
	Script    map[string]interface{}
	ParseFunc func(rawJSON string) (string, error)
}

// Indicators is a global map of all possible indicator you want to register.
var Indicators = map[string]*IndicatorDefinition{}

func init() {
	Indicators["supertrend"] = &IndicatorDefinition{
		Name:      "supertrend",
		Script:    mustUnmarshal(superTrendRaw),
		ParseFunc: parseSuperTrend,
	}

	//Indicators["rsi"] = &IndicatorDefinition{
	//	Name:      "rsi",
	//	Script:    mustUnmarshal(rsiRaw),
	//	ParseFunc: parseRsi,
	//}
}

func GetIndicatorScript(indicatorName string) (map[string]interface{}, error) {
	ind, ok := Indicators[indicatorName]
	if !ok {
		return nil, fmt.Errorf("no such indicator: %s", indicatorName)
	}
	return ind.Script, nil
}

// mustUnmarshal is a small helper that unmarshals embedded JSON or panics
func mustUnmarshal(raw string) map[string]interface{} {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal embedded script: %v", err))
	}
	return obj
}
