# Contributing to WardenGrid

## Code Standards

All source code, comments in commit messages, pull request descriptions, and documentation must be written in English. Conversational discussion in issues may be in Turkish or English.

Code must be self-documenting. Function and variable names must clearly describe their purpose without relying on inline comments. Inline comments inside function bodies are not accepted in pull requests.

Unsafe functions that bypass memory or bounds safety are not permitted. Every new feature must include corresponding unit tests under the matching `_test.go` file.

## Architecture Rules

WardenGrid follows a strict layered architecture, documented in `docs/ARCHITECTURE.md`. Contributions must respect the following constraints.

A layer must not depend on a layer above it in the stack.
The detection engine must not depend on whether input originates from live capture or from the simulator.
Every protocol parser must implement the same internal interface used by existing parsers.
The detection and protocol layers must not make external network calls.

## Commit Message Format

```text
<type>: <subject>

[optional body]
```

Accepted types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`.

Example:

```text
feat: add DNP3 frame parser

Implements basic DNP3 application layer parsing following the
same interface used by the Modbus parser.
```

## Pull Request Process

Fork the repository and create a feature branch from `main`.
Ensure `go build ./...` and `go test ./...` pass locally before opening a pull request.
Ensure `go vet ./...` reports no issues.
Open a pull request against `main` with a clear description of the change and its motivation.
All pull requests are validated automatically by the CI workflow across Linux, macOS, and Windows.

## Reporting Security Issues

Do not open a public issue for security vulnerabilities. See `SECURITY.md` for the disclosure process.

## Author

Batuhan ALGUL
