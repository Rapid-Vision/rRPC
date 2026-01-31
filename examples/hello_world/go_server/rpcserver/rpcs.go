// THIS CODE IS GENERATED

package rpcserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HelloWorldParams struct {
	Name    string  `json:"name"`
	Surname *string `json:"surname"`
}
type HelloWorldResult struct {
	GreetingMessage GreetingMessageModel `json:"greeting_message"`
}
type RPCHandler interface {
	HelloWorld(context.Context, HelloWorldParams) (HelloWorldResult, error)
}

func CreateHTTPHandler(rpc RPCHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /rpc/hello_world", CreateHelloWorldHandler(rpc))
	return mux
}

func CreateHelloWorldHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params HelloWorldParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.HelloWorld(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}
