package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"examples/hello_world/rpcserver"
)

type service struct{}

func (s *service) HelloWorld(_ context.Context, params rpcserver.HelloWorldParams) (rpcserver.HelloWorldResult, error) {
	surname := "-"
	if params.Surname != nil {
		surname = *params.Surname
	}

	fmt.Printf("hello_world (name=%s, surname=%s)\n", params.Name, surname)

	name := params.Name
	if params.Surname != nil && *params.Surname != "" {
		name = name + " " + *params.Surname
	}
	msg := rpcserver.GreetingMessageModel{
		Message: "Hello, " + name + "!",
	}
	return rpcserver.HelloWorldResult{GreetingMessage: msg}, nil
}

func main() {
	handler := rpcserver.CreateHTTPHandler(&service{})
	log.Fatal(http.ListenAndServe(":8080", handler))
}
