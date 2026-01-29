# Go Guide

This page covers server and Go client usage. See [schema_language.md](docs/schema_language.md) for schema syntax.

## Generate code
```bash
rRPC server -o . hello.rrpc
rRPC client --lang go -o . hello.rrpc
```
Default packages are `rpcserver` for servers and `rpcclient` for clients.

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
client := rpcclient.NewRPCClient("http://localhost:8080")
greeting, err := client.HelloWorld(context.Background(), rpcclient.HelloWorldParams{Name: "Ada"})
```
Go clients return the RPC result type directly (not the wrapper struct).

## Context and HTTP client
RPC methods accept `context.Context` for cancellation and deadlines. To customize timeouts, proxies, or TLS, pass your own `*http.Client`:
```go
httpClient := &http.Client{Timeout: 5 * time.Second}
client := rpcclient.NewRPCClientWithHTTP("http://localhost:8080", httpClient)
```

## Headers and auth
Use the headers-capable constructor for middleware or auth:
```go
client := rpcclient.NewRPCClientWithHeaders(
	"http://localhost:8080",
	map[string]string{"Authorization": "Bearer token"},
)
```

## Error handling
Errors are returned as typed Go errors on non-2xx responses:
- `rpcclient.ValidationRPCError`
- `rpcclient.InputRPCError`
- `rpcclient.UnauthorizedRPCError`
- `rpcclient.ForbiddenRPCError`
- `rpcclient.NotImplementedRPCError`
- `rpcclient.CustomRPCError`

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
