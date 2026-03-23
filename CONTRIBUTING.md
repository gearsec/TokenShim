# Contributing to TokenShim

## What Belongs Here

TokenShim has a narrow scope by design: local credential auditing and (in progress) agent-safe token proxying. Contributions that strengthen secret detection, improve report quality, or extend proxy support are welcome. Contributions outside that scope will be declined.

Before starting significant work, open an issue to describe what you intend to build and why.

---

## Adding a New Secret Pattern

The most common contribution is adding detection support for a new credential type. Patterns live in `internal/doctor/patterns.go`.

Requirements for a new pattern to be merged:

- Regex must match real credentials issued by the service (provide a sample format, not a real key)
- Must not produce false positives on common non-secret strings
- Includes a unit test with both a matching and a non-matching case
- Secret value must be redacted in all report formats — verify with `tokenshim doctor check --output json`
- Pattern name follows the existing naming convention (e.g. `OpenAI`, `AWSKeyID`, `GitHub`)

---

## Development Setup

```sh
git clone https://github.com/gearsec/TokenShim.git
cd TokenShim
go mod download
go test ./...
```

### Useful commands

```sh
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Lint
golangci-lint run

# Vet
go vet ./...

# Build
go build -o tokenshim ./cmd/tokenshim
```

---

## Code Standards

- All credential-handling code must stay in `internal/doctor` or `internal/keyring`. Nothing outside those packages should hold a real credential value.
- No global state.
- Errors must be returned, not swallowed. Use `fmt.Errorf("context: %w", err)`.
- No external logging frameworks — use stdlib `log/slog`.
- Tests must not make real outbound network calls.

---

## Pull Request Process

1. Fork the repository and create a branch from `main`.
2. Make your changes. Keep commits focused — one logical change per commit.
3. Run `go test -race ./...` and `go vet ./...` locally. Both must pass.
4. Open a pull request with a description of what changed and why.
5. A maintainer will review within 5 business days.

---

## Security-Sensitive Changes

If your contribution touches secret pattern matching, report redaction, or keyring integration, note this explicitly in the PR description. These areas receive heightened review.

If you discover a security issue while contributing, follow the process in [SECURITY.md](SECURITY.md) rather than including the fix in a normal PR.

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
