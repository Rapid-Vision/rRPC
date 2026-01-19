// THIS CODE IS GENERATED

package rpcserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

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

type TestEmptyParams struct {
}
type TestEmptyResult struct {
	Empty EmptyModel `json:"empty"`
}

type TestNoReturnParams struct {
}

type TestBasicParams struct {
	Text  TextModel `json:"text"`
	Flag  bool      `json:"flag"`
	Count int       `json:"count"`
	Note  *string   `json:"note"`
}
type TestBasicResult struct {
	Text TextModel `json:"text"`
}

type TestListMapParams struct {
	Texts []TextModel       `json:"texts"`
	Flags map[string]string `json:"flags"`
}
type TestListMapResult struct {
	Nested NestedModel `json:"nested"`
}

type TestOptionalParams struct {
	Text *TextModel `json:"text"`
	Flag *bool      `json:"flag"`
}
type TestOptionalResult struct {
	Flags FlagsModel `json:"flags"`
}

type TestValidationErrorParams struct {
	Text TextModel `json:"text"`
}
type TestValidationErrorResult struct {
	Text TextModel `json:"text"`
}

type TestUnauthorizedErrorParams struct {
}
type TestUnauthorizedErrorResult struct {
	Empty EmptyModel `json:"empty"`
}

type TestForbiddenErrorParams struct {
}
type TestForbiddenErrorResult struct {
	Empty EmptyModel `json:"empty"`
}

type TestNotImplementedErrorParams struct {
}
type TestNotImplementedErrorResult struct {
	Empty EmptyModel `json:"empty"`
}

type TestCustomErrorParams struct {
}
type TestCustomErrorResult struct {
	Empty EmptyModel `json:"empty"`
}

type TestMapReturnParams struct {
}
type TestMapReturnResult struct {
	Result map[string]TextModel `json:"result"`
}

type TestJsonParams struct {
	Data any `json:"data"`
}
type TestJsonResult struct {
	Json any `json:"json"`
}

type TestRawParams struct {
	Payload json.RawMessage `json:"payload"`
}
type TestRawResult struct {
	Raw json.RawMessage `json:"raw"`
}

type TestMixedPayloadParams struct {
	Payload PayloadModel `json:"payload"`
}
type TestMixedPayloadResult struct {
	Payload PayloadModel `json:"payload"`
}

type RPCHandler interface {
	TestEmpty(context.Context, TestEmptyParams) (TestEmptyResult, error)
	TestNoReturn(context.Context, TestNoReturnParams) error
	TestBasic(context.Context, TestBasicParams) (TestBasicResult, error)
	TestListMap(context.Context, TestListMapParams) (TestListMapResult, error)
	TestOptional(context.Context, TestOptionalParams) (TestOptionalResult, error)
	TestValidationError(context.Context, TestValidationErrorParams) (TestValidationErrorResult, error)
	TestUnauthorizedError(context.Context, TestUnauthorizedErrorParams) (TestUnauthorizedErrorResult, error)
	TestForbiddenError(context.Context, TestForbiddenErrorParams) (TestForbiddenErrorResult, error)
	TestNotImplementedError(context.Context, TestNotImplementedErrorParams) (TestNotImplementedErrorResult, error)
	TestCustomError(context.Context, TestCustomErrorParams) (TestCustomErrorResult, error)
	TestMapReturn(context.Context, TestMapReturnParams) (TestMapReturnResult, error)
	TestJson(context.Context, TestJsonParams) (TestJsonResult, error)
	TestRaw(context.Context, TestRawParams) (TestRawResult, error)
	TestMixedPayload(context.Context, TestMixedPayloadParams) (TestMixedPayloadResult, error)
}

func CreateHTTPHandler(rpc RPCHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /rpc/test_empty", CreateTestEmptyHandler(rpc))
	mux.Handle("POST /rpc/test_no_return", CreateTestNoReturnHandler(rpc))
	mux.Handle("POST /rpc/test_basic", CreateTestBasicHandler(rpc))
	mux.Handle("POST /rpc/test_list_map", CreateTestListMapHandler(rpc))
	mux.Handle("POST /rpc/test_optional", CreateTestOptionalHandler(rpc))
	mux.Handle("POST /rpc/test_validation_error", CreateTestValidationErrorHandler(rpc))
	mux.Handle("POST /rpc/test_unauthorized_error", CreateTestUnauthorizedErrorHandler(rpc))
	mux.Handle("POST /rpc/test_forbidden_error", CreateTestForbiddenErrorHandler(rpc))
	mux.Handle("POST /rpc/test_not_implemented_error", CreateTestNotImplementedErrorHandler(rpc))
	mux.Handle("POST /rpc/test_custom_error", CreateTestCustomErrorHandler(rpc))
	mux.Handle("POST /rpc/test_map_return", CreateTestMapReturnHandler(rpc))
	mux.Handle("POST /rpc/test_json", CreateTestJsonHandler(rpc))
	mux.Handle("POST /rpc/test_raw", CreateTestRawHandler(rpc))
	mux.Handle("POST /rpc/test_mixed_payload", CreateTestMixedPayloadHandler(rpc))
	return mux
}

func CreateTestEmptyHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestEmptyParams
		res, err := rpc.TestEmpty(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestNoReturnHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestNoReturnParams
		if err := rpc.TestNoReturn(r.Context(), params); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, struct{}{})
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
		res, err := rpc.TestBasic(r.Context(), params)
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
		res, err := rpc.TestListMap(r.Context(), params)
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
		res, err := rpc.TestOptional(r.Context(), params)
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
		res, err := rpc.TestValidationError(r.Context(), params)
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
		res, err := rpc.TestUnauthorizedError(r.Context(), params)
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
		res, err := rpc.TestForbiddenError(r.Context(), params)
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
		res, err := rpc.TestNotImplementedError(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestCustomErrorHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestCustomErrorParams
		res, err := rpc.TestCustomError(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestMapReturnHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestMapReturnParams
		res, err := rpc.TestMapReturn(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestJsonHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestJsonParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestJson(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestRawHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestRawParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestRaw(r.Context(), params)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, res)
	})
}

func CreateTestMixedPayloadHandler(rpc RPCHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params TestMixedPayloadParams
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&params); err != nil && err != io.EOF {
			writeError(w, InputError{Message: err.Error()})
			return
		}
		res, err := rpc.TestMixedPayload(r.Context(), params)
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
