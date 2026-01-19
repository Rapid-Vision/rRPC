package rpcserver

import (
	"context"
	"encoding/json"
	"errors"
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

type rpcError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type InputError struct {
	Message string
}

func (e InputError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e ForbiddenError) Error() string {
	return e.Message
}

type NotImplementedError struct {
	Message string
}

func (e NotImplementedError) Error() string {
	return e.Message
}

const (
	errorTypeCustom     = "custom"
	errorTypeValidation = "validation"
	errorTypeInput      = "input"
	errorTypeAuth       = "unauthorized"
	errorTypeForbidden  = "forbidden"
	errorTypeNotImpl    = "not_implemented"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	errType := errorTypeCustom
	msg := "error"
	if err != nil {
		msg = err.Error()
	}

	var validationErrPtr *ValidationError
	if errors.As(err, &validationErrPtr) {
		status = http.StatusBadRequest
		errType = errorTypeValidation
		msg = validationErrPtr.Error()
	}
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		status = http.StatusBadRequest
		errType = errorTypeValidation
		msg = validationErr.Error()
	}
	var inputErrPtr *InputError
	if errors.As(err, &inputErrPtr) {
		status = http.StatusBadRequest
		errType = errorTypeInput
		msg = inputErrPtr.Error()
	}
	var inputErr InputError
	if errors.As(err, &inputErr) {
		status = http.StatusBadRequest
		errType = errorTypeInput
		msg = inputErr.Error()
	}
	var unauthorizedErrPtr *UnauthorizedError
	if errors.As(err, &unauthorizedErrPtr) {
		status = http.StatusUnauthorized
		errType = errorTypeAuth
		msg = unauthorizedErrPtr.Error()
	}
	var unauthorizedErr UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		status = http.StatusUnauthorized
		errType = errorTypeAuth
		msg = unauthorizedErr.Error()
	}
	var forbiddenErrPtr *ForbiddenError
	if errors.As(err, &forbiddenErrPtr) {
		status = http.StatusForbidden
		errType = errorTypeForbidden
		msg = forbiddenErrPtr.Error()
	}
	var forbiddenErr ForbiddenError
	if errors.As(err, &forbiddenErr) {
		status = http.StatusForbidden
		errType = errorTypeForbidden
		msg = forbiddenErr.Error()
	}
	var notImplErrPtr *NotImplementedError
	if errors.As(err, &notImplErrPtr) {
		status = http.StatusNotImplemented
		errType = errorTypeNotImpl
		msg = notImplErrPtr.Error()
	}
	var notImplErr NotImplementedError
	if errors.As(err, &notImplErr) {
		status = http.StatusNotImplemented
		errType = errorTypeNotImpl
		msg = notImplErr.Error()
	}

	writeJSON(w, status, rpcError{Type: errType, Message: msg})
}

func WriteAuthError(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusUnauthorized, rpcError{Type: errorTypeAuth, Message: message})
}

func WriteUnauthorizedError(w http.ResponseWriter, message string) {
	WriteAuthError(w, message)
}
