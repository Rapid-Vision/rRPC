# Go Guide

This page covers server and Go client usage. See [schema_language.md](docs/schema_language.md) for schema syntax.

## Generate code
```bash
rRPC server -o . hello.rrpc
rRPC client --lang go -o . hello.rrpc
```
Defaults: server package `rpcserver`, client package `rpc_client`.

## Implement the server
Generated handlers expect a context:
```go
type service struct{}

func (s *service) HelloWorld(ctx context.Context, params rpcserver.HelloWorldParams) (rpcserver.HelloWorldResult, error) {
	_ = ctx
	return rpcserver.HelloWorldResult{
		GreetingMessage: rpcserver.GreetingMessageModel{Message: "Hello, " + params.Name + "!"},
	}, nil
}
```

## Start HTTP server
```go
handler := rpcserver.CreateHTTPHandler(&service{})
http.ListenAndServe(":8080", handler)
```

## Go client usage
```go
client := rpc_client.NewRPCClient("http://localhost:8080")
greeting, err := client.HelloWorld(rpc_client.HelloWorldParams{Name: "Ada"})
```

## Headers and auth
Use the headers-capable constructor for middleware or auth:
```go
client := rpc_client.NewRPCClientWithHeaders(
	"http://localhost:8080",
	map[string]string{"Authorization": "Bearer token"},
)
```

## Error handling
Errors are returned as typed Go errors on non-2xx responses:
- `rpc_client.ValidationRPCError`
- `rpc_client.InputRPCError`
- `rpc_client.UnauthorizedRPCError`
- `rpc_client.ForbiddenRPCError`
- `rpc_client.NotImplementedRPCError`
- `rpc_client.CustomRPCError`

For middleware, generated servers expose helpers like:
```go
rpcserver.WriteUnauthorizedError(w, "missing token")
```

## Prefixes
Routes are prefixed with `/rpc` by default. Override with:
```bash
rRPC server --prefix api hello.rrpc
rRPC client --lang go --prefix api hello.rrpc
```
