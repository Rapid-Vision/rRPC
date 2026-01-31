// THIS CODE IS GENERATED

package rpcserver

type TextModel struct {
	Title *string `json:"title"`
	Data  string  `json:"data"`
}
type SliceModel struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}
type StatsModel struct {
	Ascii      bool           `json:"ascii"`
	WordCount  map[string]int `json:"word_count"`
	TotalWords int            `json:"total_words"`
	Sentences  []SliceModel   `json:"sentences"`
}
