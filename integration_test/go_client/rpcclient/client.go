// THIS CODE IS GENERATED

package rpcclient

import (
	"bytes"
	"context"
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

type ErrHTTP struct {
	Status int
	Body   string
}

func (e ErrHTTP) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("rpc error: status %d", e.Status)
	}
	return fmt.Sprintf("rpc error: status %d: %s", e.Status, e.Body)
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
	baseURL     string
	client      *http.Client
	headers     map[string]string
	bearerToken string
}

func NewRPCClient(baseURL string) *RPCClient {
	return &RPCClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  http.DefaultClient,
		headers: map[string]string{},
	}
}

func (c *RPCClient) WithHTTPClient(client *http.Client) *RPCClient {
	if client == nil {
		client = http.DefaultClient
	}
	next := c.clone()
	next.client = client
	return next
}

func (c *RPCClient) WithHeaders(headers map[string]string) *RPCClient {
	next := c.clone()
	if headers == nil {
		return next
	}
	for key, value := range headers {
		next.headers[key] = value
	}
	return next
}

func (c *RPCClient) WithBearerToken(token string) *RPCClient {
	next := c.clone()
	next.bearerToken = token
	return next
}

func (c *RPCClient) clone() *RPCClient {
	copiedHeaders := make(map[string]string, len(c.headers))
	for key, value := range c.headers {
		copiedHeaders[key] = value
	}
	return &RPCClient{
		baseURL:     c.baseURL,
		client:      c.client,
		headers:     copiedHeaders,
		bearerToken: c.bearerToken,
	}
}
func (c *RPCClient) TestEmpty(ctx context.Context) (EmptyModel, error) {
	var zero EmptyModel
	var res TestEmptyResult
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_empty", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}
func (c *RPCClient) TestNoReturn(ctx context.Context) error {
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_no_return", payload, nil); err != nil {
		return err
	}
	return nil
}
func (c *RPCClient) TestBasic(ctx context.Context, params TestBasicParams) (TextModel, error) {
	var zero TextModel
	var res TestBasicResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_basic", payload, &res); err != nil {
		return zero, err
	}
	return res.Text, nil
}
func (c *RPCClient) TestListMap(ctx context.Context, params TestListMapParams) (NestedModel, error) {
	var zero NestedModel
	var res TestListMapResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_list_map", payload, &res); err != nil {
		return zero, err
	}
	return res.Nested, nil
}
func (c *RPCClient) TestOptional(ctx context.Context, params TestOptionalParams) (FlagsModel, error) {
	var zero FlagsModel
	var res TestOptionalResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_optional", payload, &res); err != nil {
		return zero, err
	}
	return res.Flags, nil
}
func (c *RPCClient) TestValidationError(ctx context.Context, params TestValidationErrorParams) (TextModel, error) {
	var zero TextModel
	var res TestValidationErrorResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_validation_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Text, nil
}
func (c *RPCClient) TestUnauthorizedError(ctx context.Context) (EmptyModel, error) {
	var zero EmptyModel
	var res TestUnauthorizedErrorResult
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_unauthorized_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}
func (c *RPCClient) TestForbiddenError(ctx context.Context) (EmptyModel, error) {
	var zero EmptyModel
	var res TestForbiddenErrorResult
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_forbidden_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}
func (c *RPCClient) TestNotImplementedError(ctx context.Context) (EmptyModel, error) {
	var zero EmptyModel
	var res TestNotImplementedErrorResult
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_not_implemented_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}
func (c *RPCClient) TestCustomError(ctx context.Context) (EmptyModel, error) {
	var zero EmptyModel
	var res TestCustomErrorResult
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_custom_error", payload, &res); err != nil {
		return zero, err
	}
	return res.Empty, nil
}
func (c *RPCClient) TestMapReturn(ctx context.Context) (map[string]TextModel, error) {
	var zero map[string]TextModel
	var res TestMapReturnResult
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_map_return", payload, &res); err != nil {
		return zero, err
	}
	return res.Result, nil
}
func (c *RPCClient) TestJson(ctx context.Context, params TestJsonParams) (any, error) {
	var zero any
	var res TestJsonResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_json", payload, &res); err != nil {
		return zero, err
	}
	return res.Json, nil
}
func (c *RPCClient) TestRaw(ctx context.Context, params TestRawParams) (json.RawMessage, error) {
	var zero json.RawMessage
	var res TestRawResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_raw", payload, &res); err != nil {
		return zero, err
	}
	return res.Raw, nil
}
func (c *RPCClient) TestMixedPayload(ctx context.Context, params TestMixedPayloadParams) (PayloadModel, error) {
	var zero PayloadModel
	var res TestMixedPayloadResult
	var payload any
	payload = params
	if err := c.doRequest(ctx, "/rpc/test_mixed_payload", payload, &res); err != nil {
		return zero, err
	}
	return res.Payload, nil
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

func (c *RPCClient) doRequest(ctx context.Context, path string, payload any, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	url := c.baseURL + path
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode payload: %w", err)
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.bearerToken != "" {
		hasAuthHeader := false
		for key := range c.headers {
			if http.CanonicalHeaderKey(key) == "Authorization" {
				hasAuthHeader = true
				break
			}
		}
		if !hasAuthHeader {
			req.Header.Set("Authorization", "Bearer "+c.bearerToken)
		}
	}
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
				return ErrHTTP{Status: resp.StatusCode, Body: strings.TrimSpace(string(raw))}
			}
		}
		return ErrHTTP{Status: resp.StatusCode}
	}
	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
