# Python Guide

This page covers generating Python clients and servers. See [schema_language.md](docs/schema_language.md) for schema syntax.

## Generate code
```bash
rRPC client -o . hello.rrpc
```
The default output package is `rpcclient`.

## Generate a FastAPI server
```bash
rRPC server --lang py -o . hello.rrpc
```
The default output package is `rpcserver`. The generated server uses FastAPI and Pydantic.

Example `server.py`:
```python
from rpcserver import create_app, RPCHandlers
from rpcserver.models import GreetingModel

class Service(RPCHandlers):
    def hello(self, name: str) -> GreetingModel:
        return GreetingModel(message=f"Hello, {name}!")

app = create_app(Service())
```
Run with:
```bash
uvicorn server:app --host 127.0.0.1 --port 8080
```

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

## Pydantic validation
To enable input validation, generate the client with Pydantic models:
```bash
rRPC client --py-pydantic -o . hello.rrpc
```
The client validates RPC inputs with Pydantic before sending requests.

## Prefixes
Routes are prefixed with `/rpc` by default. Override with:
```bash
rRPC client --prefix api hello.rrpc
```
