package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"integration_test/rpcserver"
)

type service struct{}

func (s *service) TestEmpty(params rpcserver.TestEmptyParams) (rpcserver.TestEmptyResult, error) {
	_ = params
	return rpcserver.TestEmptyResult{Empty: rpcserver.EmptyModel{}}, nil
}

func (s *service) TestBasic(params rpcserver.TestBasicParams) (rpcserver.TestBasicResult, error) {
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

func (s *service) TestListMap(params rpcserver.TestListMapParams) (rpcserver.TestListMapResult, error) {
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

func (s *service) TestOptional(params rpcserver.TestOptionalParams) (rpcserver.TestOptionalResult, error) {
	enabled := params.Flag != nil && *params.Flag
	result := rpcserver.FlagsModel{
		Enabled: enabled,
		Retries: 0,
		Labels:  []string{},
		Meta:    map[string]string{},
	}
	return rpcserver.TestOptionalResult{Flags: result}, nil
}

func (s *service) TestValidationError(params rpcserver.TestValidationErrorParams) (rpcserver.TestValidationErrorResult, error) {
	if strings.TrimSpace(params.Text.Body) == "" {
		return rpcserver.TestValidationErrorResult{}, &rpcserver.ValidationError{Message: "body is required"}
	}
	return rpcserver.TestValidationErrorResult{Text: params.Text}, nil
}

func (s *service) TestUnauthorizedError(params rpcserver.TestUnauthorizedErrorParams) (rpcserver.TestUnauthorizedErrorResult, error) {
	_ = params
	return rpcserver.TestUnauthorizedErrorResult{}, rpcserver.UnauthorizedError{Message: "missing token"}
}

func (s *service) TestForbiddenError(params rpcserver.TestForbiddenErrorParams) (rpcserver.TestForbiddenErrorResult, error) {
	_ = params
	return rpcserver.TestForbiddenErrorResult{}, rpcserver.ForbiddenError{Message: "not allowed"}
}

func (s *service) TestNotImplementedError(params rpcserver.TestNotImplementedErrorParams) (rpcserver.TestNotImplementedErrorResult, error) {
	_ = params
	return rpcserver.TestNotImplementedErrorResult{}, rpcserver.NotImplementedError{Message: "not implemented"}
}

func (s *service) TestCustomError(params rpcserver.TestCustomErrorParams) (rpcserver.TestCustomErrorResult, error) {
	_ = params
	return rpcserver.TestCustomErrorResult{}, errors.New("custom failure")
}

func (s *service) TestMapReturn(params rpcserver.TestMapReturnParams) (rpcserver.TestMapReturnResult, error) {
	_ = params
	text := rpcserver.TextModel{
		Title: nil,
		Body:  "mapped",
	}
	return rpcserver.TestMapReturnResult{Result: map[string]rpcserver.TextModel{"a": text}}, nil
}

func main() {
	handler := rpcserver.CreateHTTPHandler(&service{})
	log.Fatal(http.ListenAndServe(":8080", handler))
}
