package indicator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type mmriPayload struct {
	Node string `json:"node"`
	Ns   struct {
		D string `json:"d"`
	} `json:"ns"`
}

type mmriDWrapper struct {
	GraphicsCmds struct {
		Create struct {
			DwgTableCells []dwgTableCells `json:"dwgtablecells"`
		} `json:"create"`
	} `json:"graphicsCmds"`
}

type dwgTableCells struct {
	Data []cellData `json:"data"`
}

type cellData struct {
	ID  int    `json:"id"`
	Col int    `json:"col"`
	Row int    `json:"row"`
	T   string `json:"t"`
}

func parseMMRI(rawJSON string) (string, error) {
	// 1) Unmarshal the top-level "du" message
	var top duMessage
	if err := json.Unmarshal([]byte(rawJSON), &top); err != nil {
		return "", fmt.Errorf("failed to unmarshal top-level du: %v", err)
	}
	if top.M != "du" || len(top.P) < 2 {
		return "", errors.New("not a 'du' message or missing payload")
	}

	// 2) Inside top.P[1], we expect something like {"st1": { "node": "...", "ns": {...} } }
	var stMap map[string]mmriPayload
	if err := json.Unmarshal(top.P[1], &stMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal du payload for mmri: %v", err)
	}

	// 3) Find whichever key starts with "st" (st1, st2, ...)
	for k, v := range stMap {
		if strings.HasPrefix(k, "st") {
			// 4) Now parse the v.Ns.D field, which is another layer of JSON
			var wrap mmriDWrapper
			if err := json.Unmarshal([]byte(v.Ns.D), &wrap); err != nil {
				return "", fmt.Errorf("failed to unmarshal 'ns.d' JSON: %v", err)
			}

			// 5) Dig into wrap.GraphicsCmds.Create.DwgTableCells
			cells := wrap.GraphicsCmds.Create.DwgTableCells
			if len(cells) == 0 {
				return "", errors.New("no table cells found in mmri 'd' JSON")
			}

			// Need first entry in dwgtablecells
			firstTable := cells[0]
			if len(firstTable.Data) < 2 {
				return "", errors.New("table does not have enough cells for MMRI")
			}

			// Take the second data cell
			valStr := strings.TrimSpace(firstTable.Data[1].T)
			if valStr == "" {
				return "", errors.New("MMRI value cell is empty")
			}
			return valStr, nil
		}
	}
	return "", errors.New("no stN key found for mmri in 'du' payload")
}
