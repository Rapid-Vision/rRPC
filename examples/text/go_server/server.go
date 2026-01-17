package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"examples/text/rpcserver"
)

type textEntry struct {
	ID   int
	Text rpcserver.TextModel
}

type service struct {
	mu     sync.Mutex
	nextID int
	store  map[int]rpcserver.TextModel
}

func newService() *service {
	return &service{
		nextID: 1,
		store:  make(map[int]rpcserver.TextModel),
	}
}

func (s *service) SubmitText(_ context.Context, params rpcserver.SubmitTextParams) (rpcserver.SubmitTextResult, error) {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.store[id] = params.Text
	s.mu.Unlock()
	return rpcserver.SubmitTextResult{Int: id}, nil
}

func (s *service) ComputeStats(_ context.Context, params rpcserver.ComputeStatsParams) (rpcserver.ComputeStatsResult, error) {
	s.mu.Lock()
	text, ok := s.store[params.TextId]
	s.mu.Unlock()
	if !ok {
		return rpcserver.ComputeStatsResult{}, rpcserver.ValidationError{Message: "unknown text id"}
	}

	stats := buildStats(text.Data)
	return rpcserver.ComputeStatsResult{Stats: stats}, nil
}

func buildStats(data string) rpcserver.StatsModel {
	words := strings.Fields(data)
	wordCount := make(map[string]int, len(words))
	for _, word := range words {
		wordCount[word]++
	}

	sentences := splitSentences(data)

	return rpcserver.StatsModel{
		Ascii:      isASCII(data),
		WordCount:  wordCount,
		TotalWords: len(words),
		Sentences:  sentences,
	}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

func splitSentences(s string) []rpcserver.SliceModel {
	var slices []rpcserver.SliceModel
	start := 0
	for i, r := range s {
		if r == '.' || r == '!' || r == '?' {
			if i+1 > start {
				slices = append(slices, rpcserver.SliceModel{
					Begin: start,
					End:   i + 1,
				})
			}
			start = i + 1
		}
	}
	if start < len(s) {
		slices = append(slices, rpcserver.SliceModel{
			Begin: start,
			End:   len(s),
		})
	}
	return slices
}

func main() {
	handler := rpcserver.CreateHTTPHandler(newService())
	log.Fatal(http.ListenAndServe(":8080", handler))
}
