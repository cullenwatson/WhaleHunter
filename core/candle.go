package core

import (
	"encoding/json"
	"fmt"
	"github.com/cullenwatson/WhaleHunter/model"
	"strings"
	"time"
)

func ExtractCandles(msg string) ([]model.Candle, error) {
	var update model.TimescaleUpdate
	if err := json.Unmarshal([]byte(msg), &update); err != nil {
		return nil, fmt.Errorf("failed to unmarshal timescale_update: %w", err)
	}

	if update.M != "timescale_update" {
		return nil, fmt.Errorf("message is not timescale_update (got: %s)", update.M)
	}
	if len(update.P) < 2 {
		return nil, fmt.Errorf("unexpected structure in timescale_update")
	}

	var container model.SdsOneContainer
	if err := json.Unmarshal(update.P[1], &container); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sds_1 container: %w", err)
	}

	bars := container.Sds1.S
	var out []model.Candle

	for _, b := range bars {
		if len(b.V) < 6 {
			continue
		}
		epoch := int64(b.V[0])
		openV := b.V[1]
		highV := b.V[2]
		lowV := b.V[3]
		closeV := b.V[4]
		vol := b.V[5]

		t := time.Unix(epoch, 0)
		out = append(out, model.Candle{
			Date:   t,
			Open:   openV,
			High:   highV,
			Low:    lowV,
			Close:  closeV,
			Volume: vol,
		})
	}
	return out, nil
}

func ParseMultipleMessages(raw string) []string {
	var results []string
	// Split on "~m~"
	parts := strings.Split(raw, "~m~")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 0 && (p[0] == '{' || p[0] == '[') {
			results = append(results, p)
		}
	}
	return results
}
