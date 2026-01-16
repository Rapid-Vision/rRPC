package rpcserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type EmptyModelModel struct {
}
type TextModelModel struct {
	Title *string `json:"title"`
	Body  string  `json:"body"`
}
type FlagsModelModel struct {
	Enabled bool              `json:"enabled"`
	Retries int               `json:"retries"`
	Labels  []string          `json:"labels"`
	Meta    map[string]string `json:"meta"`
}
type NestedModelModel struct {
	Text   TextModelModel            `json:"text"`
	Flags  *FlagsModelModel          `json:"flags"`
	Items  []TextModelModel          `json:"items"`
	Lookup map[string]TextModelModel `json:"lookup"`
}

type TestEmptyParams struct {
}

type TestEmptyResult struct {
	EmptyModel EmptyModelModel `json:"empty_model"`
}

type TestBasicParams struct {
	Text  TextModelModel `json:"text"`
	Flag  bool           `json:"flag"`
	Count int            `json:"count"`
	Note  *string        `json:"note"`
}

type TestBasicResult struct {
	TextModel TextModelModel `json:"text_model"`
}

type TestListMapParams struct {
	Texts []TextModelModel  `json:"texts"`
	Flags map[string]string `json:"flags"`
}

type TestListMapResult struct {
	NestedModel NestedModelModel `json:"nested_model"`
}

type TestOptionalParams struct {
	Text *TextModelModel `json:"text"`
	Flag *bool           `json:"flag"`
}

type TestOptionalResult struct {
	FlagsModel FlagsModelModel `json:"flags_model"`
}

type TestValidationErrorParams struct {
	Text TextModelModel `json:"text"`
}

type TestValidationErrorResult struct {
	TextModel TextModelModel `json:"text_model"`
}

type TestUnauthorizedErrorParams struct {
}

type TestUnauthorizedErrorResult struct {
	EmptyModel EmptyModelModel `json:"empty_model"`
}

type TestForbiddenErrorParams struct {
}

type TestForbiddenErrorResult struct {
	EmptyModel EmptyModelModel `json:"empty_model"`
}

type TestNotImplementedErrorParams struct {
}

type TestNotImplementedErrorResult struct {
	EmptyModel EmptyModelModel `json:"empty_model"`
}

type RPCHandler interface {
	TestEmpty(TestEmptyParams) (TestEmptyResult, error)
	TestBasic(TestBasicParams) (TestBasicResult, error)
	TestListMap(TestListMapParams) (TestListMapResult, error)
	TestOptional(TestOptionalParams) (TestOptionalResult, error)
	TestValidationError(TestValidationErrorParams) (TestValidationErrorResult, error)
	TestUnauthorizedError(TestUnauthorizedErrorParams) (TestUnauthorizedErrorResult, error)
	TestForbiddenError(TestForbiddenErrorParams) (TestForbiddenErrorResult, error)
	TestNotImplementedError(TestNotImplementedErrorParams) (TestNotImplementedErrorResult, error)
}

func CreateHTTPHandler(rpc RPCHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /rpc/test_empty", CreateTestEmptyHandler(rpc))
	mux.Handle("POST /rpc/test_basic", CreateTestBasicHandler(rpc))
	mux.Handle("POST /rpc/test_list_map", CreateTestListMapHandler(rpc))
	mux.Handle("POST /rpc/test_optional", CreateTestOptionalHandler(rpc))
	mux.Handle("POST /rpc/test_validation_error", CreateTestValidationErrorHandler(rpc))
	mux.Handle("POST /rpc/test_unauthorized_error", CreateTestUnauthorizedErrorHandler(rpc))
	mux.Handle("POST /rpc/test_forbidden_error", CreateTestForbiddenErrorHandler(rpc))
	mux.Handle("POST /rpc/test_not_implemented_error", CreateTestNotImplementedErrorHandler(rpc))
	return mux
}

func CreateTestEmptyHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestEmptyParams
		res, err := rpc.TestEmpty(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestBasicHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestBasicParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestBasic(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestListMapHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestListMapParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestListMap(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestOptionalHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestOptionalParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestOptional(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestValidationErrorHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestValidationErrorParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestValidationError(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestUnauthorizedErrorHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestUnauthorizedErrorParams
		res, err := rpc.TestUnauthorizedError(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestForbiddenErrorHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestForbiddenErrorParams
		res, err := rpc.TestForbiddenError(params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestNotImplementedErrorHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestNotImplementedErrorParams
		res, err := rpc.TestNotImplementedError(params)
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

	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		status = http.StatusBadRequest
		errType = errorTypeValidation
		msg = validationErr.Error()
	}
	var inputErr InputError
	if errors.As(err, &inputErr) {
		status = http.StatusBadRequest
		errType = errorTypeInput
		msg = inputErr.Error()
	}
	var unauthorizedErr UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		status = http.StatusUnauthorized
		errType = errorTypeAuth
		msg = unauthorizedErr.Error()
	}
	var forbiddenErr ForbiddenError
	if errors.As(err, &forbiddenErr) {
		status = http.StatusForbidden
		errType = errorTypeForbidden
		msg = forbiddenErr.Error()
	}
	var notImplErr NotImplementedError
	if errors.As(err, &notImplErr) {
		status = http.StatusNotImplemented
		errType = errorTypeNotImpl
		msg = notImplErr.Error()
	}

	writeJSON(w, status, rpcError{Type: errType, Message: msg})
}
