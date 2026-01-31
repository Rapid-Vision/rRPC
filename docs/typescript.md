# TypeScript Client

This page covers generating a TypeScript client. See [schema_language.md](docs/schema_language.md) for schema syntax.

## Generate a client
```bash
rRPC client --lang ts -o . hello.rrpc
```
The default output package is `rpcclient`.
For zod input validation:
```bash
rRPC client --lang ts --ts-zod -o . hello.rrpc
```

## Basic usage
```ts
import { RPCClient } from "./rpcclient";

const rpc = new RPCClient("http://localhost:8080");
const greeting = await rpc.hello({ name: "Ada" });
```

## Prefixes
Routes are prefixed with `/rpc` by default. Override with:
```bash
rRPC client --lang ts --prefix api hello.rrpc
```

## Options
```ts
const rpc = new RPCClient("localhost:8080", {
	prefix: "/rpc",
	bearerToken: "token",
	headers: { "X-Trace-Id": "trace" },
	timeoutMs: 2000,
	fetchFn: customFetch,
});
```

- `prefix` configures the RPC path prefix.
- `bearerToken` automatically sets `Authorization: Bearer <token>` unless you override it in `headers`.
- `headers` adds custom headers to every request.
- `timeoutMs` sets an abort timeout in milliseconds.
- `fetchFn` lets you inject a custom `fetch` implementation for testing or instrumentation.

## Zod validation
When generated with `--ts-zod`, the client validates RPC inputs using zod before sending requests.
Install zod in your project:
```bash
npm install zod
```
The generated file exports `*Schema` constants (e.g. `UserModelSchema`, `HelloParamsSchema`) that you can reuse.

## Error handling
RPC errors are thrown as typed exceptions:
- `ValidationRPCError`
- `InputRPCError`
- `UnauthorizedRPCError`
- `ForbiddenRPCError`
- `NotImplementedRPCError`
- `CustomRPCError`

Non-JSON error responses are mapped to `CustomRPCError` with `type = "custom"`.
