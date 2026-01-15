# rRPC
rRPC is a simple codegen tool for creating RPC APIs from a defined schema.

## Motivation
The industry standard for communication between services is [gRPC](https://grpc.io/). It may be good for Google-scale services, but has several disadvantages: 
- Official protobuf compiler lacks types for Python
- It is bloated
- HTTP/2 may introduce bugs of its own
- Binary protocol is harder to debug

For small and medium size projects the performance and industrial adoption of the gRPC toolkit may not outweigh these problems.

## Goals
This project aims to provide a simple toolkit with following properties:
- Generated code has strict typing
- Single binary for everything
- JSON over HTTP protocol
- Capacity to generate an [OpenAPI](https://www.openapis.org/) schema

At this moment this project aims to generate only a Go server and a python client.

## Installation
```bash
go install github.com/Rapid-Vision/rRPC
```

## Tutorial (Hello World)
Use the `examples/hello_world` schema as a minimal end-to-end example.

Schema (`hello.rrpc`):
```rrpc
model GreetingMessage {
    message: string
}

rpc HelloWorld(
    name: string,
    surname: string?,
) GreetingMessage
```

Generate server and client:
```bash
rRPC server -o . hello.rrpc
rRPC client -o . hello.rrpc
```

Implement the server (Go):
```go
package main

import (
	"log"
	"net/http"

	"examples/hello_world/rpcserver"
)

type service struct{}

func (s *service) HelloWorld(params rpcserver.HelloWorldParams) (rpcserver.HelloWorldResult, error) {
	msg := rpcserver.GreetingMessageModel{
		Message: "Hello, " + params.Name + "!",
	}
	return rpcserver.HelloWorldResult{GreetingMessage: msg}, nil
}

func main() {
	handler := rpcserver.CreateHTTPHandler(&service{})
	log.Fatal(http.ListenAndServe(":8080", handler))
}
```

Call it from Python:
```python
from rpc_client import RPCClient

rpc = RPCClient("http://localhost:8080")
greeting = rpc.hello_world(name="Ada", surname="Lovelace")
print("greeting:", greeting)
```

