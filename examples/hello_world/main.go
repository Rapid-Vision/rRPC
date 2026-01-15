package main

import (
	"log"
	"net/http"

	"examples/hello_world/rpc"
)

type service struct{}

func (s *service) HelloWorld(params rpc.HelloWorldParams) (rpc.HelloWorldResult, error) {
	name := params.Name
	if params.Surname != nil && *params.Surname != "" {
		name = name + " " + *params.Surname
	}
	msg := rpc.GreetingMessageModel{
		Message: "Hello, " + name + "!",
	}
	return rpc.HelloWorldResult{GreetingMessage: msg}, nil
}

func main() {
	handler := rpc.CreateHTTPHandler(&service{})
	log.Fatal(http.ListenAndServe(":8080", handler))
}
