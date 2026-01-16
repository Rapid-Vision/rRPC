# Errors

All error responses are JSON with this shape:
```json
{
  "type": "validation",
  "message": "field name is required"
}
```

## Error types

### Framework errors
- `input`: JSON decode or request format errors
### User errors
- `custom`: default for user-returned errors
- `validation`: validation failures in handler logic
- `unauthorized`: authentication missing/invalid
- `forbidden`: authenticated but not allowed
- `not_implemented`: endpoint not implemented

## Status code mapping
- `validation` -> `400 Bad Request`
- `input` -> `400 Bad Request`
- `unauthorized` -> `401 Unauthorized`
- `forbidden` -> `403 Forbidden`
- `not_implemented` -> `501 Not Implemented`
- `custom` -> `500 Internal Server Error`

## Returning errors from server code (Go)
Generated servers expose error types. Return them from handlers:
```go
func (s *service) GetUser(params rpcserver.GetUserParams) (rpcserver.GetUserResult, error) {
	if params.Id <= 0 {
		return rpcserver.GetUserResult{}, rpcserver.ValidationError{Message: "id must be positive"}
	}
	if !s.authenticated(params.token) {
		return rpcserver.GetUserResult{}, rpcserver.UnauthorizedError{Message: "missing token"}
	}
	return rpcserver.GetUserResult{}, rpcserver.NotImplementedError{Message: "not implemented"}
}
```

Input errors are produced automatically when JSON decoding fails in generated handlers.
