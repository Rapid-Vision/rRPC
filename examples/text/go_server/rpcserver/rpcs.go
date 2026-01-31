// THIS CODE IS GENERATED

package rpcserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type SubmitTextParams struct {
	Text TextModel `json:"text"`
}
type SubmitTextResult struct {
	Int int `json:"int"`
}

type ComputeStatsParams struct {
	TextId int `json:"text_id"`
}
type ComputeStatsResult struct {
	Stats StatsModel `json:"stats"`
}
type RPCHandler interface {
	SubmitText(context.Context, SubmitTextParams) (SubmitTextResult, error)
	ComputeStats(context.Context, ComputeStatsParams) (ComputeStatsResult, error)
}

func CreateHTTPHandler(rpc RPCHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /rpc/submit_text", CreateSubmitTextHandler(rpc))
	mux.Handle("POST /rpc/compute_stats", CreateComputeStatsHandler(rpc))
	return mux
}

func CreateSubmitTextHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params SubmitTextParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.SubmitText(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateComputeStatsHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params ComputeStatsParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.ComputeStats(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}
