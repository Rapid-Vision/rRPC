# rRPC
rRPC is a simple codegen tool for creating RPC APIs from a defined schema. It does not aim to create a new procol. rRPC generates boilerplate code from a defined schema. It's like [sqlc](https://sqlc.dev) for APIs.

## Motivation
The industry standard for communication between services is [gRPC](https://grpc.io/). It may be good for Google-scale services, but has several disadvantages: 
- Official protobuf compiler lacks types for Python
- It is bloated
- HTTP/2 may introduce bugs of its own
- Binary protocol is harder to debug

For small and medium size projects the performance and industrial adoption of the gRPC toolkit may not outweigh these problems.

## Goals
This project aims to provide a simple tool with following properties:
- Generated code has strict typing
- Single binary for everything
- JSON over HTTP protocol
- Generated code does not add dependencies
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

## Docs
- [Getting started](docs/getting_started.md)
- [Schema language description](docs/schema_language.md)
- [Error handling](docs/errors.md)
- [Go guide](docs/go.md)
- [Python guide](docs/python.md)
- [Protocol description](docs/protocol.md)

## Usage examples
See [`examples/`](examples/) directory for server, client and Makefile implemenation examples. Also [`integration_test/`](integration_test/) may be useful as reference too.

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
