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

## Language support
| Language | Server | Client |
| --- | --- | --- |
| Go | ✅ | ✅ |
| Python | ❌ | ✅ |

Other languages can be supported via OpenAPI toolkits.

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

## OpenAPI
You can generate an OpenAPI 3.0 schema for a given rRPC schema:
```bash
rRPC openapi -o . hello.rrpc
```
This writes `openapi.json` under the output directory. Use it with tools like Swagger UI or client generators.

## Comparison & Fit
This project focuses on a small, typed, JSON-over-HTTP RPC flow.

### Compared to other tools
- **[gRPC](https://grpc.io/)**: gRPC is a full-featured RPC system with strong tooling, streaming, and HTTP/2. rRPC is lighter and simpler but lacks streaming, interceptors, and a mature ecosystem.
- **[OpenAPI](https://www.openapis.org/)**: OpenAPI is an API description format with broad tooling for REST-style endpoints. rRPC is RPC-oriented and does not target REST semantics or multiple transports.
- **[GraphQL](https://graphql.org/)**: GraphQL offers flexible client queries and a rich type system. rRPC is schema-first but request/response shapes are fixed per method and not queryable.
- **[CUE](https://cuelang.org/)**: CUE is a general configuration and validation language. rRPC is narrowly scoped to RPC schema + codegen rather than validation or policy.
- **[TypeSpec](https://typespec.io/)**: TypeSpec is a rich API modeling language with multiple emitters. rRPC is smaller, has a simpler DSL, and targets a limited set of generators.

### When this may be useful
- You want a small schema language and minimal runtime.
- You want strict typing with simple JSON over HTTP.

### When this is not a good fit
- You need streaming, bidirectional RPC, or advanced middleware.
- You need multi-language support beyond Go/Python.
- You want REST or GraphQL semantics and tooling.
