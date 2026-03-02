package events

import "encoding/json"

type Event struct {
    Type string          `json:"type"`
    Room string          `json:"room"`
    User string          `json:"user"`
    Data json.RawMessage `json:"data"`
}
