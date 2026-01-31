# Integration tests

This directory contains a full end-to-end test setup for rRPC. It generates a server and clients for all supported languages.

## How it works
- `integration_test/run_tests.py` builds the CLI (`./rRPC`) and regenerates test artifacts.
- It generates:
  - Go server into `integration_test/go_server`
  - Go client into `integration_test/go_client`
  - Python client into `integration_test/py_client`
  - TypeScript client into `integration_test/ts_client`
  - OpenAPI spec into `integration_test/openapi.json`
- It starts the generated Go server on `http://localhost:8080`.
- It runs tests for:
  - Go client (`go test .`)
  - Python client (`python -m unittest test_client.py`)
  - TypeScript client (`bun test test_client.ts`)

## Run automatically

From the repo root:
```bash
python integration_test/run_tests.py
```

## Run manually

Build rRPC tool
```bash
go build -o rRPC .
```

Generate code
```bash
cd integration_test
make
```

Start test server
```bash
cd integration_test/go_server
go run .
```

Run go client tests
```bash
cd integration_test/go_client
go test .
```

Run python client tests
```bash
cd integration_test/py_client
python -m unittest test_client.py
```

Run typescript client tests
```bash
cd integration_test/ts_client
bun test test_client.ts
```
