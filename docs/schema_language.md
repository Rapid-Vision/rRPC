# Schema Language

rRPC schemas define models and RPCs.

## Models
```rrpc
model User {
    id: int
    name: string
    email: string?
    flags: map[string]
    tags: list[string]
}
```

## RPCs
```rrpc
rpc GetUser(
    id: int,
) User
```
Parameters are named fields. The return type is a single type.

## Types
- Builtins: `string`, `int`, `bool`, `json`, `raw`
- Optional: `string?`, `User?`
- Lists: `list[Type]`
- Maps: `map[Type]` (JSON keys are strings)

## Nesting
Types can be nested:
```rrpc
model Report {
    items: list[User]
    summary: map[string]
}
```

## Comments
Lines starting with `#` are ignored:
```rrpc
# This is a comment
rpc Ping() bool
```
