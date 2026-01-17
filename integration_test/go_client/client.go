package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	client "integration_test/client/rpc_client"
)

const (
	baseURL     = "http://localhost:8080"
	bearerToken = "test_token"
)

type testCase struct {
	name string
	fn   func() error
}

func runTest(tc testCase) bool {
	if err := tc.fn(); err != nil {
		fmt.Printf("FAIL %s: %v\n", tc.name, err)
		return false
	}
	fmt.Printf("OK %s\n", tc.name)
	return true
}

func main() {
	rpc := client.NewRPCClientWithHeaders(baseURL, map[string]string{
		"Authorization": "Bearer " + bearerToken,
	})

	tests := []testCase{
		{name: "empty", fn: func() error {
			_, err := rpc.TestEmpty()
			return err
		}},
		{name: "basic", fn: func() error {
			note := "note"
			text := client.TextModel{Title: nil, Body: "  hello  "}
			res, err := rpc.TestBasic(client.TestBasicParams{
				Text:  text,
				Flag:  true,
				Count: 3,
				Note:  &note,
			})
			if err != nil {
				return err
			}
			if res.Body != "hello" {
				return fmt.Errorf("expected body 'hello', got %q", res.Body)
			}
			if res.Title == nil || *res.Title != "note" {
				return fmt.Errorf("expected title 'note', got %v", res.Title)
			}
			return nil
		}},
		{name: "list_map", fn: func() error {
			res, err := rpc.TestListMap(client.TestListMapParams{
				Texts: []client.TextModel{
					{Title: stringPtr("t1"), Body: "b1"},
					{Title: stringPtr("t2"), Body: "b2"},
				},
				Flags: map[string]string{"mode": "fast"},
			})
			if err != nil {
				return err
			}
			if res.Flags == nil {
				return fmt.Errorf("expected flags")
			}
			if res.Flags.Retries != 2 {
				return fmt.Errorf("expected retries 2, got %d", res.Flags.Retries)
			}
			if res.Flags.Meta["mode"] != "fast" {
				return fmt.Errorf("expected mode fast, got %q", res.Flags.Meta["mode"])
			}
			if _, ok := res.Lookup["first"]; !ok {
				return fmt.Errorf("expected lookup to include 'first'")
			}
			return nil
		}},
		{name: "optional", fn: func() error {
			res, err := rpc.TestOptional(client.TestOptionalParams{
				Text: nil,
				Flag: nil,
			})
			if err != nil {
				return err
			}
			if res.Enabled {
				return fmt.Errorf("expected enabled false")
			}
			return nil
		}},
		{name: "validation_error", fn: func() error {
			_, err := rpc.TestValidationError(client.TestValidationErrorParams{
				Text: client.TextModel{Title: nil, Body: ""},
			})
			var vErr client.ValidationRPCError
			if err == nil || !errors.As(err, &vErr) {
				return fmt.Errorf("expected ValidationRPCError, got %v", err)
			}
			return nil
		}},
		{name: "input_error", fn: func() error {
			return sendInvalidPayload()
		}},
		{name: "unauthorized_error", fn: func() error {
			_, err := rpc.TestUnauthorizedError()
			var uErr client.UnauthorizedRPCError
			if err == nil || !errors.As(err, &uErr) {
				return fmt.Errorf("expected UnauthorizedRPCError, got %v", err)
			}
			return nil
		}},
		{name: "forbidden_error", fn: func() error {
			_, err := rpc.TestForbiddenError()
			var fErr client.ForbiddenRPCError
			if err == nil || !errors.As(err, &fErr) {
				return fmt.Errorf("expected ForbiddenRPCError, got %v", err)
			}
			return nil
		}},
		{name: "not_implemented_error", fn: func() error {
			_, err := rpc.TestNotImplementedError()
			var nErr client.NotImplementedRPCError
			if err == nil || !errors.As(err, &nErr) {
				return fmt.Errorf("expected NotImplementedRPCError, got %v", err)
			}
			return nil
		}},
		{name: "custom_error", fn: func() error {
			_, err := rpc.TestCustomError()
			var cErr client.CustomRPCError
			if err == nil || !errors.As(err, &cErr) {
				return fmt.Errorf("expected CustomRPCError, got %v", err)
			}
			return nil
		}},
		{name: "map_return", fn: func() error {
			res, err := rpc.TestMapReturn()
			if err != nil {
				return err
			}
			text, ok := res["a"]
			if !ok {
				return fmt.Errorf("expected key 'a'")
			}
			if text.Body != "mapped" {
				return fmt.Errorf("expected body 'mapped', got %q", text.Body)
			}
			return nil
		}},
	}

	passed := 0
	for _, tc := range tests {
		if runTest(tc) {
			passed++
		}
	}
	fmt.Printf("passed %d/%d\n", passed, len(tests))
	if passed != len(tests) {
		os.Exit(1)
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
		return fmt.Errorf("expected 400, got %d", resp.StatusCode)
	}
	var payload client.RPCError
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	if payload.Type != client.RPCErrorInput {
		return fmt.Errorf("expected input error type, got %q", payload.Type)
	}
	return nil
}
