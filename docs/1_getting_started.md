# Getting Started

rRPC is a small schema-first RPC generator for a Go server and a Python/Go client.

## Install
```bash
go install github.com/Rapid-Vision/rRPC
```

## Define a schema
Create a `.rrpc` file:
```rrpc
model Greeting {
    message: string
}

rpc Hello(
    name: string,
) Greeting
```

## Generate code
```bash
rRPC server -o . hello.rrpc
rRPC client -o . hello.rrpc
```
Generated code is written to `./<pkg>/` (default packages: `rpcserver` and `rpc_client`).

## Implement RPCHandler interface

```go
type service struct{}

func (s *service) HelloWorld(_ context.Context, params rpcserver.HelloWorldParams) (rpcserver.HelloWorldResult, error) {
	msg := rpcserver.GreetingMessageModel{
		Message: "Hello, " + params.Name + "!",
	}
	return rpcserver.HelloWorldResult{GreetingMessage: msg}, nil
}
```

## Run the Go server
```go
handler := rpcserver.CreateHTTPHandler(&service{})
http.ListenAndServe(":8080", handler)
```

## Call from Python
```python
rpc = RPCClient("http://localhost:8080")
greeting = rpc.hello(name="Ada")
```

## Prefixes
Routes are prefixed with `/rpc` by default. Override with:
```bash
rRPC server --prefix api hello.rrpc
rRPC client --prefix api hello.rrpc
```
