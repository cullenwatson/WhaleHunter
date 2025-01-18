package indicator

import "encoding/json"

type duMessage struct {
	M string            `json:"m"`
	P []json.RawMessage `json:"p"`
}

type stBar struct {
	I int       `json:"i"`
	V []float64 `json:"v"`
}

type st1Payload struct {
	Node string  `json:"node"`
	St   []stBar `json:"st"`
}

type studyPayload struct {
	St1 st1Payload `json:"st1"`
}
