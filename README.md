# Shim

**Don't give AI agents keys, use a Shim.**

Shim is a local proxy that keeps real API credentials out of AI agent processes. The agent is given a masked token and a local proxy address. When it makes an API call, Shim swaps the token for the real credential in-flight and forwards the request. The real key never touches the agent's environment.

---

## How It Works

```
┌──────────────────┐          ┌──────────────────┐
│    AI Agent      │  masked  │                  │
│                  │─────────►│       SHIM       │
│  OPENAI_KEY=     │  token   │     (proxy)      │
│  sk-shim-***     │          │                  │
└──────────────────┘          └────────┬─────────┘
                                       │
                               real key injected
                               on tool calls only
                                       │
                                       ▼
                              ┌──────────────────┐
                              │   API Service    │
                              │  OpenAI, etc.    │
                              └──────────────────┘
```

Shim runs a local HTTP proxy on loopback. The agent starts with a masked, non-functional token (`sk-shim-<id>`) and routes outbound requests through that proxy. The real credential is only injected at the point of a tool call to an external service - it is never added to the prompt or message payload sent to the LLM. Prompt contents are logged, cached, and visible to the model. The real key stays out of all of that.

---

## Usage

### Basic: wrap any command with shim exec

```sh
# Store your real key once
shim secrets set OPENAI_API_KEY sk-real-...

# Run your agent - it receives a masked token and a local proxy address
shim exec -- python agent.py

# Inside agent.py, the environment looks like:
# OPENAI_API_KEY=sk-shim-a1b2c3d4
# OPENAI_BASE_URL=http://127.0.0.1:8742/openai
# Shim handles injection transparently
```

### Named profiles

```sh
# Create a profile for a specific agent workload
shim profile create research --service openai --service anthropic

# Execute with a named profile
shim exec --profile research -- claude-agent --task "summarize docs"
```

### Inspect what the agent sees

```sh
# Print the masked environment that will be injected
shim exec --dry-run -- python agent.py
```

### Audit log

```sh
# View all requests proxied in the last session
shim log show --session latest

# Shim logs: timestamp, masked token used, target host, HTTP method, status
# Shim does NOT log request/response bodies
```

---

## Installation

```sh
# Homebrew (macOS/Linux)
brew install gearsec/tap/shim

# Go install
go install github.com/gearsec/shim/cmd/shim@latest

# From source
git clone https://github.com/gearsec/shim.git
cd shim
make build
```

---

## Directory Structure

```
shim/
├── cmd/
│   └── shim/           # CLI entrypoint
│       └── main.go
├── pkg/
│   ├── proxy/          # Local HTTP/HTTPS proxy engine
│   ├── masking/        # Token masking and alias generation
│   └── config/         # Profile and secrets configuration
├── internal/
│   ├── injection/      # In-flight credential injection
│   └── keyring/        # OS keyring integration (Keychain, Secret Service, WCM)
├── .github/
│   ├── workflows/
│   └── ISSUE_TEMPLATE/
├── SECURITY.md
├── CONTRIBUTING.md
└── CODE_OF_CONDUCT.md
```

---

## Design Principles

- **Local-first.** No cloud component, no telemetry, no account required. All secrets stay on your machine.
- **Process separation.** The agent process cannot access real credentials by design.
- **Minimal surface area.** Shim does one thing: intercept, swap, forward. No features that expand the trust boundary.
- **Auditable.** The proxy access log gives you a complete record of what the agent called and when.
- **No plaintext at rest.** Secrets are stored in the OS keyring (Keychain on macOS, Secret Service on Linux, Windows Credential Manager on Windows) - never written to disk by Shim.

---

## Supported Services

| Service    | Status |
|------------|--------|
| OpenAI     | Stable |
| Anthropic  | Stable |
| Cohere     | Beta   |
| Mistral    | Beta   |
| AWS Bedrock | Planned |
| Google Vertex | Planned |

Adding a new service shim requires implementing a small interface. See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

MIT - see [LICENSE](LICENSE).

---

## Security

See [SECURITY.md](SECURITY.md) for the vulnerability disclosure process.
