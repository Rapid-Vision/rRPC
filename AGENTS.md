# Repository Guidelines

## Project Structure & Module Organization
- `main.go` is the CLI entry point; it wires into Cobra commands in `cmd/`.
- `cmd/` holds the CLI command tree (currently `root.go`).
- `internal/` contains implementation packages such as the lexer (`internal/lexer/`).
- `examples/` stores sample schemas like `examples/example.rrpc`.

## Build, Test, and Development Commands
- `go build ./...` builds the CLI and all packages.
- `go run ./...` runs the CLI from source (handy during development).
- `go test ./...` runs all tests (there are currently no tests, but keep this green).
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
