// THIS CODE IS GENERATED

package rpcserver

import "encoding/json"

type EmptyModel struct {
}
type TextModel struct {
	Title *string `json:"title"`
	Body  string  `json:"body"`
}
type FlagsModel struct {
	Enabled bool              `json:"enabled"`
	Retries int               `json:"retries"`
	Labels  []string          `json:"labels"`
	Meta    map[string]string `json:"meta"`
}
type NestedModel struct {
	Text   TextModel            `json:"text"`
	Flags  *FlagsModel          `json:"flags"`
	Items  []TextModel          `json:"items"`
	Lookup map[string]TextModel `json:"lookup"`
}
type PayloadModel struct {
	Data    any             `json:"data"`
	RawData json.RawMessage `json:"raw_data"`
}
