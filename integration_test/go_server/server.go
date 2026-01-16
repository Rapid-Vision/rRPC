package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"integration_test/rpcserver"
)

const bearerToken = "test_token"

type service struct{}

func (s *service) TestEmpty(_ context.Context, params rpcserver.TestEmptyParams) (rpcserver.TestEmptyResult, error) {
	_ = params
	return rpcserver.TestEmptyResult{Empty: rpcserver.EmptyModel{}}, nil
}

func (s *service) TestBasic(_ context.Context, params rpcserver.TestBasicParams) (rpcserver.TestBasicResult, error) {
	title := params.Text.Title
	if title == nil && params.Note != nil {
		title = params.Note
	}
	out := rpcserver.TextModel{
		Title: title,
		Body:  strings.TrimSpace(params.Text.Body),
	}
	return rpcserver.TestBasicResult{Text: out}, nil
}

func (s *service) TestListMap(_ context.Context, params rpcserver.TestListMapParams) (rpcserver.TestListMapResult, error) {
	flags := rpcserver.FlagsModel{
		Enabled: true,
		Retries: len(params.Texts),
		Labels:  []string{"ok"},
		Meta:    params.Flags,
	}
	out := rpcserver.NestedModel{
		Text:   params.Texts[0],
		Flags:  &flags,
		Items:  params.Texts,
		Lookup: map[string]rpcserver.TextModel{"first": params.Texts[0]},
	}
	return rpcserver.TestListMapResult{Nested: out}, nil
}

func (s *service) TestOptional(_ context.Context, params rpcserver.TestOptionalParams) (rpcserver.TestOptionalResult, error) {
	enabled := params.Flag != nil && *params.Flag
	result := rpcserver.FlagsModel{
		Enabled: enabled,
		Retries: 0,
		Labels:  []string{},
		Meta:    map[string]string{},
	}
	return rpcserver.TestOptionalResult{Flags: result}, nil
}

func (s *service) TestValidationError(_ context.Context, params rpcserver.TestValidationErrorParams) (rpcserver.TestValidationErrorResult, error) {
	if strings.TrimSpace(params.Text.Body) == "" {
		return rpcserver.TestValidationErrorResult{}, &rpcserver.ValidationError{Message: "body is required"}
	}
	return rpcserver.TestValidationErrorResult{Text: params.Text}, nil
}

func (s *service) TestUnauthorizedError(_ context.Context, params rpcserver.TestUnauthorizedErrorParams) (rpcserver.TestUnauthorizedErrorResult, error) {
	_ = params
	return rpcserver.TestUnauthorizedErrorResult{}, rpcserver.UnauthorizedError{Message: "missing token"}
}

func (s *service) TestForbiddenError(_ context.Context, params rpcserver.TestForbiddenErrorParams) (rpcserver.TestForbiddenErrorResult, error) {
	_ = params
	return rpcserver.TestForbiddenErrorResult{}, rpcserver.ForbiddenError{Message: "not allowed"}
}

func (s *service) TestNotImplementedError(_ context.Context, params rpcserver.TestNotImplementedErrorParams) (rpcserver.TestNotImplementedErrorResult, error) {
	_ = params
	return rpcserver.TestNotImplementedErrorResult{}, rpcserver.NotImplementedError{Message: "not implemented"}
}

func (s *service) TestCustomError(_ context.Context, params rpcserver.TestCustomErrorParams) (rpcserver.TestCustomErrorResult, error) {
	_ = params
	return rpcserver.TestCustomErrorResult{}, errors.New("custom failure")
}

func (s *service) TestMapReturn(_ context.Context, params rpcserver.TestMapReturnParams) (rpcserver.TestMapReturnResult, error) {
	_ = params
	text := rpcserver.TextModel{
		Title: nil,
		Body:  "mapped",
	}
	return rpcserver.TestMapReturnResult{Result: map[string]rpcserver.TextModel{"a": text}}, nil
}

func main() {
	handler := rpcserver.CreateHTTPHandler(&service{})
	log.Fatal(http.ListenAndServe(":8080", authMiddleware(handler)))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+bearerToken {
			writeAuthError(w, "missing or invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeAuthError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = fmt.Fprintf(w, `{"type":"unauthorized","message":%q}`, message)
}
