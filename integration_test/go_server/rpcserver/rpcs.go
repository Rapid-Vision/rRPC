// THIS CODE IS GENERATED

package rpcserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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
