package rpc

import (
	"encoding/json"
	"io"
	"net/http"
)

type GreetingMessageModel struct {
	Message string `json:"message"`
}

type HelloWorldParams struct {
	Name    string  `json:"name"`
	Surname *string `json:"surname"`
}

type HelloWorldResult struct {
	GreetingMessage GreetingMessageModel `json:"greeting_message"`
}

type RPCHandler interface {
	HelloWorld(HelloWorldParams) (HelloWorldResult, error)
}

func CreateHTTPHandler(rpc RPCHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/hello_world", CreateHelloWorldHandler(rpc))
	return mux
}

func CreateHelloWorldHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params HelloWorldParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		res, err := rpc.HelloWorld(params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

type rpcError struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	msg := "error"
	if err != nil {
		msg = err.Error()
	}
	writeJSON(w, status, rpcError{Message: msg})
}
