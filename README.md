# TokenShim

**Don't give AI agents keys, use a Shim.**

TokenShim is a local security tool for managing and auditing API credentials in AI agent environments. It has two core capabilities:

| Feature | Status |
|---------|--------|
| **Doctor** — scan files and env vars for exposed secrets | ✅ Active |
| **Token proxy** — intercept and swap masked tokens in-flight | 🚧 In progress |

---

## Doctor Mode

Doctor scans your local machine for exposed API keys and secrets across files and environment variables, and produces an exportable report.

```sh
# Run all checks (files + environment)
tokenshim doctor check

# Scan only environment variables
tokenshim doctor check --env

# Scan only files
tokenshim doctor check --files

# Use a custom config (supports wildcards: **/.env*, **/env.*, **/*.env)
tokenshim doctor check --config /path/to/doctor.yaml

# Export as HTML
tokenshim doctor check --output html --export report.html

# Supported formats: json (default), yaml, xml, csv, html
tokenshim doctor check --output yaml --export report.yaml
```

Doctor reads `~/.config/tokenshim/doctor.yaml` for the list of files to scan. If the config does not exist, built-in defaults are used. All secret values are **redacted** in reports (`sk-ab****xyz`).

**Default scan paths:**

```yaml
scan_paths:
  - "~/.env"
  - "~/.bashrc"
  - "~/.zshrc"
  - "~/.bash_profile"
  - "~/.profile"
  - "~/.config/fish/config.fish"
  - ".env"
  - ".env.local"
  - ".env.development"
  - ".env.production"
  - ".env.staging"
```

**Detected secret types:** OpenAI, Anthropic, AWS (Key ID + Secret), GitHub, Stripe, HuggingFace, Google API Key, Slack, Twilio, and generic high-entropy API secrets.

---

## Token Proxy *(in progress)*

The proxy mode keeps real API credentials out of AI agent processes entirely. The agent is given a masked token and a local proxy address. When it makes an API call, Shim swaps the token for the real credential in-flight and forwards the request. The real key never touches the agent's environment.

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

Once complete, the proxy will support:

| Service       | Status   |
|---------------|----------|
| OpenAI        | Planned  |
| Anthropic     | Planned  |
| Cohere        | Planned  |
| Mistral       | Planned  |
| AWS Bedrock   | Planned  |
| Google Vertex | Planned  |

---

## Installation

```sh
# From source
git clone https://github.com/gearsec/tokenshim.git
cd tokenshim
make build
```

---

## Directory Structure

```
tokenshim/
├── cmd/
│   └── tokenshim/      # CLI entrypoint
│       └── main.go
├── internal/
│   ├── cli/            # Cobra command definitions
│   │   ├── root.go
│   │   ├── exec.go
│   │   ├── secrets.go
│   │   ├── profile.go
│   │   ├── log.go
│   │   └── doctor.go   # doctor / doctor check commands
│   ├── config/         # Shared config manager (viper-backed)
│   │   ├── item.go
│   │   └── manager.go
│   ├── doctor/         # Secret detection engine
│   │   ├── patterns.go # Compiled regex patterns (OpenAI, AWS, GitHub, etc.)
│   │   ├── config.go   # doctor.yaml loading and path resolution
│   │   ├── scanner.go  # File and environment variable scanning
│   │   └── report.go   # Multi-format report serialization (JSON/YAML/XML/CSV/HTML)
│   ├── injection/      # In-flight credential injection (in progress)
│   └── keyring/        # OS keyring integration (Keychain, Secret Service, WCM)
├── pkg/
│   ├── proxy/          # Local HTTP/HTTPS proxy engine (in progress)
│   └── masking/        # Token masking and alias generation (in progress)
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
- **No plaintext at rest.** Secrets are stored in the OS keyring (Keychain on macOS, Secret Service on Linux, Windows Credential Manager on Windows) — never written to disk.
- **Auditable.** Every scan produces a full report of what was checked and what was found.
- **Minimal surface area.** No features that expand the trust boundary beyond what is necessary.

---

## License

MIT — see [LICENSE](LICENSE).

---

## Security

See [SECURITY.md](SECURITY.md) for the vulnerability disclosure process.
