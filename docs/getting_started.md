# Getting Started

rRPC is a small schema-first RPC generator for a Go server and Python/Go/TypeScript clients.

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

An `.rrpc` file can be formatted using the `rRPC format` command.

## Generate code
```bash
rRPC server -o . hello.rrpc
rRPC server --lang py -o . hello.rrpc
rRPC client -o . hello.rrpc
rRPC client --lang go -o . hello.rrpc
rRPC client --lang ts -o . hello.rrpc
```
Generated code is written to `./<pkg>/` (default packages: `rpcserver` and `rpcclient`).

## Creating a go server

Implement RPCHandler interface
```go
type service struct{}

func (s *service) HelloWorld(_ context.Context, params rpcserver.HelloWorldParams) (rpcserver.HelloWorldResult, error) {
	msg := rpcserver.GreetingMessageModel{
		Message: "Hello, " + params.Name + "!",
	}
	return rpcserver.HelloWorldResult{GreetingMessage: msg}, nil
}
```

Run the Go server
```go
handler := rpcserver.CreateHTTPHandler(&service{})
http.ListenAndServe(":8080", handler)
```

## Creating a python server
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

## Call from Python
```python
rpc = RPCClient("http://localhost:8080")
greeting = rpc.hello(name="Ada")
```

## Call from TypeScript
```ts
import { RPCClient } from "./rpcclient";

const rpc = new RPCClient("http://localhost:8080");
const greeting = await rpc.hello({ name: "Ada" });
```

## Prefixes
Routes are prefixed with `/rpc` by default. Override with `--prefix` flag
```bash
rRPC server --prefix api hello.rrpc
rRPC client --prefix api hello.rrpc
```
