# Contributing to Shim

## What Belongs Here

Shim has a narrow scope by design. Contributions that expand the core separation guarantee — broader service support, tighter injection logic, better keyring integration — are welcome. Contributions that add features outside that scope (telemetry, cloud sync, AI model management) will be declined.

Before starting significant work, open an issue to describe what you intend to build and why. This avoids wasted effort on contributions that don't fit the project's direction.

---

## Adding a New Service Shim

The most common contribution is adding proxy support for a new AI service. Each service shim must implement the `ServiceShim` interface in `pkg/proxy`:

```go
type ServiceShim interface {
    // Name returns the canonical identifier for this service (e.g., "openai").
    Name() string

    // InjectCredential replaces the masked token in the outbound request
    // with the real credential. It must not log or store the real credential.
    InjectCredential(req *http.Request, real string) error

    // BaseURL returns the upstream host this shim forwards to.
    BaseURL() *url.URL
}
```

Requirements for a service shim to be merged:

- Handles credential injection entirely in `InjectCredential` — no credential handling in any other method
- Includes unit tests that verify the masked token is absent from the forwarded request
- Includes a test fixture that confirms the real credential does not appear in logs
- Documents the exact header or field where the credential is injected
- Has no external dependencies beyond the Go standard library unless unavoidable

---

## Development Setup

```sh
git clone https://github.com/gearsec/shim.git
cd shim
go mod download
make test
```

### Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Compile the `shim` binary |
| `make test` | Run unit and integration tests |
| `make lint` | Run golangci-lint |
| `make vet` | Run go vet |
| `make security` | Run gosec static analysis |

---

## Code Standards

- All credential-handling code must be in `internal/injection` or `internal/keyring`. Nothing outside those packages should ever hold a real credential value.
- No global state.
- No `init()` functions that touch secrets or configuration.
- Errors must be returned, not swallowed. Use `fmt.Errorf("context: %w", err)`.
- No external logging frameworks. Use the stdlib `log/slog`.
- Tests must not make real outbound network calls. Use `httptest.NewServer` for service fixtures.

---

## Pull Request Process

1. Fork the repository and create a branch from `main`.
2. Make your changes. Keep commits focused — one logical change per commit.
3. Run `make test lint vet security` locally. All must pass.
4. Open a pull request with a description of what changed and why.
5. A maintainer will review within 5 business days. Address feedback and push to the same branch.

Pull requests that modify injection logic or keyring integration require review from at least one maintainer before merge.

---

## Security-Sensitive Changes

If your contribution touches credential injection, keyring access, or the subprocess environment construction in `shim exec`, note this explicitly in the PR description. These areas receive heightened review. Do not rush them.

If you discover a security issue while contributing, follow the process in [SECURITY.md](SECURITY.md) rather than including the fix in a normal PR.

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
