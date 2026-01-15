# rRPC
rRPC is a simple codegen tool for creating RPC API's from a defined schema.

## Motivation
Industry standard for communication between services is [gRPC](https://grpc.io/). It may be good for google-scale services, but has several disadvantages: 
- Official protobuf compiler lacks types for python
- It is bloated
- http2 may introduce bugs of it's own
- Binary protocol is harder to debug

For small and medium size projects the perfomance and industrial adoption of the gRPC toolkit may not outweight these problems.

## Goals
This project aims to provide a simple toolkit with following properties:
- Generated code has strict typing
- Single binary for everything
- JSON over HTTP protocol
- Capacity to generate an [OpenAPI](https://www.openapis.org/) schema

At this moment this project aims to generate only a golang server and a python client.

## Usage
Define your schema in a `.rrpc` file (see `examples/example.rrpc`), then generate code from it.

Generate Go server code:
```bash
rrpc server examples/example.rrpc
```

Generate Python client code:
```bash
rrpc client examples/example.rrpc
```
