package model

import (
	"encoding/json"
	"time"
)

// Candle holds one OHLCV record
type Candle struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// TimescaleUpdate is the main structure we expect for "timescale_update" messages.
type TimescaleUpdate struct {
	M   string            `json:"m"`
	P   []json.RawMessage `json:"p"`
	T   int64             `json:"t"`
	TMs int64             `json:"t_ms"`
}

type SdsOneContainer struct {
	Sds1 SdsOne `json:"sds_1"`
}

type SdsOne struct {
	S []Bar `json:"s"`
}

type Bar struct {
	I int       `json:"i"`
	V []float64 `json:"v"` // [timestamp, open, high, low, close, volume]
}
