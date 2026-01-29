# Python Guide

This page covers generating a Python client. See [schema_language.md](docs/schema_language.md) for schema syntax.

## Generate code
```bash
rRPC client -o . hello.rrpc
```
The default output package is `rpcclient`.

## Client usage
```python
from rpcclient import RPCClient

rpc = RPCClient("http://localhost:8080")
greeting = rpc.hello_world(name="Ada", surname="Lovelace")
```

## Timeout
Pass a timeout (seconds):
```python
rpc = RPCClient("http://localhost:8080", timeout=5.0)
```

## Headers and auth
Pass custom headers when creating the client:
```python
rpc = RPCClient(
    "http://localhost:8080",
    headers={"Authorization": "Bearer token"},
)
```

## Error handling
Non-2xx responses are parsed into typed exceptions:
- `CustomRPCError`
- `ValidationRPCError`
- `InputRPCError`
- `UnauthorizedRPCError`
- `ForbiddenRPCError`
- `NotImplementedRPCError`

Example:
```python
try:
    rpc.get_user(user_id=0)
except ValidationRPCError as err:
    print(err.error.message)
```

## Data classes
Generated models are `@dataclass` types with `from_dict(...)` helpers, and the client
uses dataclass serialization for payloads while handling nested lists/maps automatically.

## Prefixes
Routes are prefixed with `/rpc` by default. Override with:
```bash
rRPC client --prefix api hello.rrpc
```
