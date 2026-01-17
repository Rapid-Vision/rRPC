# Repository Guidelines

## Project Structure & Module Organization
- `main.go` is the CLI entry point; it wires into Cobra commands in `cmd/`.
- `cmd/` holds CLI commands like `server` and `client`.
- `internal/` contains the lexer, parser, utils, and codegen (`internal/gen/go/`, `internal/gen/python/`).
- `examples/` stores sample schemas and runnable demos (see `examples/hello_world/`).

## Build, Test, and Development Commands
- `go build ./...` builds the CLI and all packages.
- `go run ./...` runs the CLI from source (handy during development).
- `go test ./...` runs unit tests.
- `gofmt -w .` formats all Go code in the repo.

## Coding Style & Naming Conventions
- Use standard Go formatting: tabs for indentation and `gofmt` for alignment.
- Follow Go naming: `CamelCase` for exported types/functions, `lowerCamel` for locals.
- Keep package names short and lowercase (e.g., `lexer`).
- Prefer small, focused functions with clear names; avoid long parameter lists.

## Testing Guidelines
- No test framework is set up yet; when adding tests, use Go’s `testing` package.
- Name test files `*_test.go` and test functions `TestXxx`.
- Run `go test ./...` before submitting changes.

## Commit & Pull Request Guidelines
- Recent commits use short, descriptive phrases (e.g., “README draft”, “Cobra init”).
- Keep commit messages brief and imperative when possible.
- PRs should include a concise description, motivation, and any relevant examples.
- If behavior changes, include a before/after note or sample CLI usage.

## Configuration & Examples
- Sample schema files live in `examples/`; update or add examples when changing parsing or codegen behavior.

## Additional notes
- go version in the go.mod file is correct and real
