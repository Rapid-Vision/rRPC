package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	client "integration_test/client/rpcclient"
)

const (
	baseURL     = "http://localhost:8080"
	bearerToken = "test_token"
)

var backgroundCtx = context.Background()

func newClient() *client.RPCClient {
	return client.NewRPCClient(baseURL).WithBearerToken(bearerToken)
}

func TestEmpty(t *testing.T) {
	rpc := newClient()
	if _, err := rpc.TestEmpty(backgroundCtx); err != nil {
		t.Fatalf("TestEmpty failed: %v", err)
	}
}

func TestNoReturn(t *testing.T) {
	rpc := newClient()
	if err := rpc.TestNoReturn(backgroundCtx); err != nil {
		t.Fatalf("TestNoReturn failed: %v", err)
	}
}

func TestBasic(t *testing.T) {
	rpc := newClient()
	note := "note"
	res, err := rpc.TestBasic(backgroundCtx, client.TestBasicParams{
		Text:  client.TextModel{Title: nil, Body: "  hello  "},
		Flag:  true,
		Count: 3,
		Note:  &note,
	})
	if err != nil {
		t.Fatalf("TestBasic failed: %v", err)
	}
	if res.Body != "hello" {
		t.Fatalf("expected body 'hello', got %q", res.Body)
	}
	if res.Title == nil || *res.Title != "note" {
		t.Fatalf("expected title 'note', got %v", res.Title)
	}
}

func TestListMap(t *testing.T) {
	rpc := newClient()
	res, err := rpc.TestListMap(backgroundCtx, client.TestListMapParams{
		Texts: []client.TextModel{
			{Title: stringPtr("t1"), Body: "b1"},
			{Title: stringPtr("t2"), Body: "b2"},
		},
		Flags: map[string]string{"mode": "fast"},
	})
	if err != nil {
		t.Fatalf("TestListMap failed: %v", err)
	}
	if res.Flags == nil {
		t.Fatalf("expected flags")
	}
	if res.Flags.Retries != 2 {
		t.Fatalf("expected retries 2, got %d", res.Flags.Retries)
	}
	if res.Flags.Meta["mode"] != "fast" {
		t.Fatalf("expected mode fast, got %q", res.Flags.Meta["mode"])
	}
	if _, ok := res.Lookup["first"]; !ok {
		t.Fatalf("expected lookup to include 'first'")
	}
}

func TestOptional(t *testing.T) {
	rpc := newClient()
	res, err := rpc.TestOptional(backgroundCtx, client.TestOptionalParams{
		Text: nil,
		Flag: nil,
	})
	if err != nil {
		t.Fatalf("TestOptional failed: %v", err)
	}
	if res.Enabled {
		t.Fatalf("expected enabled false")
	}
}

func TestValidationError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestValidationError(backgroundCtx, client.TestValidationErrorParams{
		Text: client.TextModel{Title: nil, Body: ""},
	})
	var vErr client.ValidationRPCError
	if err == nil || !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationRPCError, got %v", err)
	}
}

func TestInputError(t *testing.T) {
	if err := sendInvalidPayload(); err != nil {
		t.Fatalf("expected input error, got %v", err)
	}
}

func TestUnauthorizedError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestUnauthorizedError(backgroundCtx)
	var uErr client.UnauthorizedRPCError
	if err == nil || !errors.As(err, &uErr) {
		t.Fatalf("expected UnauthorizedRPCError, got %v", err)
	}
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
	rpc := client.NewRPCClient(baseURL)
	_, err := rpc.TestEmpty(backgroundCtx)
	var uErr client.UnauthorizedRPCError
	if err == nil || !errors.As(err, &uErr) {
		t.Fatalf("expected UnauthorizedRPCError from middleware, got %v", err)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	rpc := client.NewRPCClient(baseURL).WithBearerToken("bad_token")
	_, err := rpc.TestEmpty(backgroundCtx)
	var uErr client.UnauthorizedRPCError
	if err == nil || !errors.As(err, &uErr) {
		t.Fatalf("expected UnauthorizedRPCError from middleware, got %v", err)
	}
}

func TestForbiddenError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestForbiddenError(backgroundCtx)
	var fErr client.ForbiddenRPCError
	if err == nil || !errors.As(err, &fErr) {
		t.Fatalf("expected ForbiddenRPCError, got %v", err)
	}
}

func TestNotImplementedError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestNotImplementedError(backgroundCtx)
	var nErr client.NotImplementedRPCError
	if err == nil || !errors.As(err, &nErr) {
		t.Fatalf("expected NotImplementedRPCError, got %v", err)
	}
}

func TestCustomError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestCustomError(backgroundCtx)
	var cErr client.CustomRPCError
	if err == nil || !errors.As(err, &cErr) {
		t.Fatalf("expected CustomRPCError, got %v", err)
	}
}

func TestMapReturn(t *testing.T) {
	rpc := newClient()
	res, err := rpc.TestMapReturn(backgroundCtx)
	if err != nil {
		t.Fatalf("TestMapReturn failed: %v", err)
	}
	text, ok := res["a"]
	if !ok {
		t.Fatalf("expected key 'a'")
	}
	if text.Body != "mapped" {
		t.Fatalf("expected body 'mapped', got %q", text.Body)
	}
}

func TestJson(t *testing.T) {
	rpc := newClient()
	payload := map[string]any{"count": 2, "tags": []any{"a", "b"}}
	res, err := rpc.TestJson(backgroundCtx, client.TestJsonParams{Data: payload})
	if err != nil {
		t.Fatalf("TestJson failed: %v", err)
	}
	result, ok := res.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", res)
	}
	if result["count"] != float64(2) {
		t.Fatalf("expected count 2, got %v", result["count"])
	}
}

func TestRaw(t *testing.T) {
	rpc := newClient()
	raw := json.RawMessage(`{"ok":true}`)
	res, err := rpc.TestRaw(backgroundCtx, client.TestRawParams{Payload: raw})
	if err != nil {
		t.Fatalf("TestRaw failed: %v", err)
	}
	if string(res) != string(raw) {
		t.Fatalf("expected raw %s, got %s", string(raw), string(res))
	}
}

func TestMixedPayload(t *testing.T) {
	rpc := newClient()
	payload := client.PayloadModel{
		Data:    map[string]any{"value": "x"},
		RawData: json.RawMessage(`{"id":1}`),
	}
	res, err := rpc.TestMixedPayload(backgroundCtx, client.TestMixedPayloadParams{Payload: payload})
	if err != nil {
		t.Fatalf("TestMixedPayload failed: %v", err)
	}
	result, ok := res.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", res.Data)
	}
	if result["value"] != "x" {
		t.Fatalf("expected value x, got %v", result["value"])
	}
	if string(res.RawData) != `{"id":1}` {
		t.Fatalf("expected raw data %s, got %s", `{"id":1}`, string(res.RawData))
	}
}

func TestContextCancelled(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, req.Context().Err()
		}),
	}
	rpc := client.NewRPCClient(baseURL).WithHTTPClient(httpClient)
	ctx, cancel := context.WithCancel(backgroundCtx)
	cancel()
	_, err := rpc.TestEmpty(ctx)
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestHTTPError(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("boom")),
				Header:     make(http.Header),
			}, nil
		}),
	}
	rpc := client.NewRPCClient(baseURL).WithHTTPClient(httpClient)
	_, err := rpc.TestEmpty(backgroundCtx)
	var httpErr client.ErrHTTP
	if err == nil || !errors.As(err, &httpErr) {
		t.Fatalf("expected ErrHTTP, got %v", err)
	}
	if httpErr.Status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", httpErr.Status)
	}
	if httpErr.Body != "boom" {
		t.Fatalf("expected body 'boom', got %q", httpErr.Body)
	}
}

func TestWithBearerTokenHeader(t *testing.T) {
	var authHeader string
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			authHeader = req.Header.Get("Authorization")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}
	rpc := client.NewRPCClient(baseURL).WithHTTPClient(httpClient).WithBearerToken(bearerToken)
	if err := rpc.TestNoReturn(backgroundCtx); err != nil {
		t.Fatalf("TestNoReturn failed: %v", err)
	}
	if authHeader != "Bearer "+bearerToken {
		t.Fatalf("expected Authorization header set, got %q", authHeader)
	}
}

func stringPtr(value string) *string {
	return &value
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func sendInvalidPayload() error {
	url := baseURL + "/rpc/test_basic"
	body := strings.NewReader(`{"text":"bad","flag":true,"count":1}`)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		return errors.New("expected 400 status")
	}
	var payload client.RPCError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	if payload.Type != client.RPCErrorInput {
		return errors.New("expected input error type")
	}
	return nil
}
