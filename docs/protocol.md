# Protocol

rRPC uses JSON over HTTP with a fixed request/response shape. It is recommended to not interact with the protocol directly and stick to using generated methods. This information is helpful for learning and debugging purposes.

## HTTP method and route
All RPCs are called with `POST`.
The default route prefix is `/rpc`.
Example RPC `GetUser` maps to:
```
POST /rpc/get_user
```
Override the prefix with `--prefix` during code generation.

## Requests
The request body is a JSON object with parameter names in snake_case as keys. For example:
```json
{ "user_id": 42 }
```
RPCs with no parameters send an empty body.

## Responses
On success, the server returns `200` with a JSON object that wraps the result:
```json
{ "user": { "id": 1, "name": "Ada" } }
```
The wrapper key is derived from the return type (model name in snake_case or `result` for collections).

## Errors
Non-2xx responses return:
```json
{ "type": "validation", "message": "field is required" }
```
Error `type` values:
`custom`, `validation`, `input`, `unauthorized`, `forbidden`, `not_implemented`.
See `docs/errors.md` for status code mapping and server-side helpers.
