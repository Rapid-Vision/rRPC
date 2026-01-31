// THIS CODE IS GENERATED

package rpcclient

import (
	"context"
	"encoding/json"
)

type TestEmptyParams struct {
}
type TestEmptyResult struct {
	Empty EmptyModel `json:"empty"`
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

type TestNoReturnParams struct {
}

func (c *RPCClient) TestNoReturn(ctx context.Context) error {
	var payload any
	payload = nil
	if err := c.doRequest(ctx, "/rpc/test_no_return", payload, nil); err != nil {
		return err
	}
	return nil
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

type TestListMapParams struct {
	Texts []TextModel       `json:"texts"`
	Flags map[string]string `json:"flags"`
}
type TestListMapResult struct {
	Nested NestedModel `json:"nested"`
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

type TestOptionalParams struct {
	Text *TextModel `json:"text"`
	Flag *bool      `json:"flag"`
}
type TestOptionalResult struct {
	Flags FlagsModel `json:"flags"`
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

type TestValidationErrorParams struct {
	Text TextModel `json:"text"`
}
type TestValidationErrorResult struct {
	Text TextModel `json:"text"`
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

type TestUnauthorizedErrorParams struct {
}
type TestUnauthorizedErrorResult struct {
	Empty EmptyModel `json:"empty"`
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

type TestForbiddenErrorParams struct {
}
type TestForbiddenErrorResult struct {
	Empty EmptyModel `json:"empty"`
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

type TestNotImplementedErrorParams struct {
}
type TestNotImplementedErrorResult struct {
	Empty EmptyModel `json:"empty"`
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

type TestCustomErrorParams struct {
}
type TestCustomErrorResult struct {
	Empty EmptyModel `json:"empty"`
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

type TestMapReturnParams struct {
}
type TestMapReturnResult struct {
	Result map[string]TextModel `json:"result"`
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

type TestJsonParams struct {
	Data any `json:"data"`
}
type TestJsonResult struct {
	Json any `json:"json"`
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

type TestRawParams struct {
	Payload json.RawMessage `json:"payload"`
}
type TestRawResult struct {
	Raw json.RawMessage `json:"raw"`
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

type TestMixedPayloadParams struct {
	Payload PayloadModel `json:"payload"`
}
type TestMixedPayloadResult struct {
	Payload PayloadModel `json:"payload"`
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
