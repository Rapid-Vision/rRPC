package rpc_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

type TestEmptyParams struct {
}

type TestEmptyResult struct {
	Empty EmptyModel `json:"empty"`
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

type RPCErrorType string

const (
	RPCErrorCustom         RPCErrorType = "custom"
	RPCErrorValidation     RPCErrorType = "validation"
	RPCErrorInput          RPCErrorType = "input"
	RPCErrorUnauthorized   RPCErrorType = "unauthorized"
	RPCErrorForbidden      RPCErrorType = "forbidden"
	RPCErrorNotImplemented RPCErrorType = "not_implemented"
)

type RPCError struct {
	Type    RPCErrorType `json:"type"`
	Message string       `json:"message"`
}

type RPCErrorException struct {
	Err RPCError
}

func (e RPCErrorException) Error() string {
	return e.Err.Message
}

type CustomRPCError struct {
	RPCError
}

func (e CustomRPCError) Error() string {
	return e.Message
}

type ValidationRPCError struct {
	RPCError
}

func (e ValidationRPCError) Error() string {
	return e.Message
}

type InputRPCError struct {
	RPCError
}

func (e InputRPCError) Error() string {
	return e.Message
}

type UnauthorizedRPCError struct {
	RPCError
}

func (e UnauthorizedRPCError) Error() string {
	return e.Message
}

type ForbiddenRPCError struct {
	RPCError
}

func (e ForbiddenRPCError) Error() string {
	return e.Message
}

type NotImplementedRPCError struct {
	RPCError
}

func (e NotImplementedRPCError) Error() string {
	return e.Message
}

type RPCClient struct {
	baseURL string
	client  *http.Client
	headers map[string]string
}

func NewRPCClient(baseURL string) *RPCClient {
	return NewRPCClientWithHTTPAndHeaders(baseURL, nil, nil)
}

func NewRPCClientWithHTTP(baseURL string, client *http.Client) *RPCClient {
	return NewRPCClientWithHTTPAndHeaders(baseURL, client, nil)
}

func NewRPCClientWithHeaders(baseURL string, headers map[string]string) *RPCClient {
	return NewRPCClientWithHTTPAndHeaders(baseURL, nil, headers)
}

func NewRPCClientWithHTTPAndHeaders(baseURL string, client *http.Client, headers map[string]string) *RPCClient {
	if client == nil {
		client = http.DefaultClient
	}
	copiedHeaders := make(map[string]string, len(headers))
	for key, value := range headers {
		copiedHeaders[key] = value
	}
	return &RPCClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		headers: copiedHeaders,
	}
}

func (c *RPCClient) TestEmpty() (EmptyModel, error) {
	var zero EmptyModel
	var res TestEmptyResult
	var payload any
	payload = nil
	if err := c.doRequest("/rpc/test_empty", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}

func (c *RPCClient) TestBasic(params TestBasicParams) (TextModel, error) {
	var zero TextModel
	var res TestBasicResult
	var payload any
	payload = params
	if err := c.doRequest("/rpc/test_basic", payload, &res); err != nil {
		return zero, err
	}
	return res.Text, nil
}

func (c *RPCClient) TestListMap(params TestListMapParams) (NestedModel, error) {
	var zero NestedModel
	var res TestListMapResult
	var payload any
	payload = params
	if err := c.doRequest("/rpc/test_list_map", payload, &res); err != nil {
		return zero, err
	}
	return res.Nested, nil
}

func (c *RPCClient) TestOptional(params TestOptionalParams) (FlagsModel, error) {
	var zero FlagsModel
	var res TestOptionalResult
	var payload any
	payload = params
	if err := c.doRequest("/rpc/test_optional", payload, &res); err != nil {
		return zero, err
	}
	return res.Flags, nil
}

func (c *RPCClient) TestValidationError(params TestValidationErrorParams) (TextModel, error) {
	var zero TextModel
	var res TestValidationErrorResult
	var payload any
	payload = params
	if err := c.doRequest("/rpc/test_validation_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Text, nil
}

func (c *RPCClient) TestUnauthorizedError() (EmptyModel, error) {
	var zero EmptyModel
	var res TestUnauthorizedErrorResult
	var payload any
	payload = nil
	if err := c.doRequest("/rpc/test_unauthorized_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}

func (c *RPCClient) TestForbiddenError() (EmptyModel, error) {
	var zero EmptyModel
	var res TestForbiddenErrorResult
	var payload any
	payload = nil
	if err := c.doRequest("/rpc/test_forbidden_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}

func (c *RPCClient) TestNotImplementedError() (EmptyModel, error) {
	var zero EmptyModel
	var res TestNotImplementedErrorResult
	var payload any
	payload = nil
	if err := c.doRequest("/rpc/test_not_implemented_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}

func (c *RPCClient) TestCustomError() (EmptyModel, error) {
	var zero EmptyModel
	var res TestCustomErrorResult
	var payload any
	payload = nil
	if err := c.doRequest("/rpc/test_custom_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}

func (c *RPCClient) TestMapReturn() (map[string]TextModel, error) {
	var zero map[string]TextModel
	var res TestMapReturnResult
	var payload any
	payload = nil
	if err := c.doRequest("/rpc/test_map_return", payload, &res); err != nil {
		return zero, err
	}
	return res.Result, nil
}

func errorFromRPCError(err RPCError) error {
	switch err.Type {
	case RPCErrorCustom:
		return CustomRPCError{RPCError: err}
	case RPCErrorValidation:
		return ValidationRPCError{RPCError: err}
	case RPCErrorInput:
		return InputRPCError{RPCError: err}
	case RPCErrorUnauthorized:
		return UnauthorizedRPCError{RPCError: err}
	case RPCErrorForbidden:
		return ForbiddenRPCError{RPCError: err}
	case RPCErrorNotImplemented:
		return NotImplementedRPCError{RPCError: err}
	default:
		return RPCErrorException{Err: err}
	}
}

func (c *RPCClient) doRequest(path string, payload any, out any) error {
	url := c.baseURL + path
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode payload: %w", err)
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if len(raw) > 0 {
			var rpcErr RPCError
			if err := json.Unmarshal(raw, &rpcErr); err == nil && rpcErr.Type != "" {
				return errorFromRPCError(rpcErr)
			}
			if strings.TrimSpace(string(raw)) != "" {
				return fmt.Errorf("rpc error: %s", strings.TrimSpace(string(raw)))
			}
		}
		return fmt.Errorf("rpc error: status %d", resp.StatusCode)
	}
	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
