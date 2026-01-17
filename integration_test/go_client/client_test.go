package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	client "integration_test/client/rpc_client"
)

const (
	baseURL     = "http://localhost:8080"
	bearerToken = "test_token"
)

func newClient() *client.RPCClient {
	return client.NewRPCClientWithHeaders(baseURL, map[string]string{
		"Authorization": "Bearer " + bearerToken,
	})
}

func TestEmpty(t *testing.T) {
	rpc := newClient()
	if _, err := rpc.TestEmpty(); err != nil {
		t.Fatalf("TestEmpty failed: %v", err)
	}
}

func TestBasic(t *testing.T) {
	rpc := newClient()
	note := "note"
	res, err := rpc.TestBasic(client.TestBasicParams{
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
	res, err := rpc.TestListMap(client.TestListMapParams{
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
	res, err := rpc.TestOptional(client.TestOptionalParams{
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
	_, err := rpc.TestValidationError(client.TestValidationErrorParams{
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
	_, err := rpc.TestUnauthorizedError()
	var uErr client.UnauthorizedRPCError
	if err == nil || !errors.As(err, &uErr) {
		t.Fatalf("expected UnauthorizedRPCError, got %v", err)
	}
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
	rpc := client.NewRPCClient(baseURL)
	_, err := rpc.TestEmpty()
	var uErr client.UnauthorizedRPCError
	if err == nil || !errors.As(err, &uErr) {
		t.Fatalf("expected UnauthorizedRPCError from middleware, got %v", err)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	rpc := client.NewRPCClientWithHeaders(baseURL, map[string]string{
		"Authorization": "Bearer bad_token",
	})
	_, err := rpc.TestEmpty()
	var uErr client.UnauthorizedRPCError
	if err == nil || !errors.As(err, &uErr) {
		t.Fatalf("expected UnauthorizedRPCError from middleware, got %v", err)
	}
}

func TestForbiddenError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestForbiddenError()
	var fErr client.ForbiddenRPCError
	if err == nil || !errors.As(err, &fErr) {
		t.Fatalf("expected ForbiddenRPCError, got %v", err)
	}
}

func TestNotImplementedError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestNotImplementedError()
	var nErr client.NotImplementedRPCError
	if err == nil || !errors.As(err, &nErr) {
		t.Fatalf("expected NotImplementedRPCError, got %v", err)
	}
}

func TestCustomError(t *testing.T) {
	rpc := newClient()
	_, err := rpc.TestCustomError()
	var cErr client.CustomRPCError
	if err == nil || !errors.As(err, &cErr) {
		t.Fatalf("expected CustomRPCError, got %v", err)
	}
}

func TestMapReturn(t *testing.T) {
	rpc := newClient()
	res, err := rpc.TestMapReturn()
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

func stringPtr(value string) *string {
	return &value
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
