# Security Policy

## Scope

Shim is a local credential separation layer. Its security guarantee is narrow and specific: **real API credentials must never appear in the environment, memory space, or logs of the agent process it wraps.** All vulnerability reports are evaluated against this guarantee first.

### In scope

- Credential leak from the Shim proxy to the agent process (environment, stdin/stdout, IPC)
- Masked token bypass — any path by which an agent could use a masked token to recover the real credential
- OS keyring exposure — secrets readable by unintended processes due to misconfigured keyring access controls
- Injection logic flaws — cases where the real credential is forwarded to an unintended host
- TLS downgrade or MITM conditions on the loopback proxy interface
- Privilege escalation via the `shim exec` subprocess model
- Dependency-level supply chain compromises affecting credential handling

### Out of scope

- Attacks requiring root/kernel-level access to the host machine (outside the threat model)
- Social engineering or physical access scenarios
- Issues in upstream AI services (OpenAI, Anthropic, etc.) — report those to the respective vendor
- Theoretical vulnerabilities without a working proof-of-concept

---

## Execution Model

Shim runs entirely on the local machine. It has no cloud component, makes no outbound connections to GearSec infrastructure, and collects no telemetry. The attack surface is limited to:

1. The loopback proxy (`127.0.0.1`, configurable port)
2. The OS keyring (Keychain on macOS, Secret Service on Linux, Windows Credential Manager on Windows)
3. The subprocess environment constructed by `shim exec`

There is no Shim server, relay, or SaaS endpoint that handles credentials.

---

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Send a report to: **security@gearsec.io**

Include:
- A clear description of the vulnerability and the affected component
- The version of Shim and OS/architecture where you reproduced it
- Step-by-step reproduction instructions
- A proof-of-concept demonstrating credential exposure, if applicable
- Your assessment of severity and exploitability

### What to expect

| Stage | Timeline |
|-------|----------|
| Acknowledgement | Within 2 business days |
| Initial triage and severity assessment | Within 5 business days |
| Fix or mitigation plan communicated to reporter | Within 14 days of triage |
| Public disclosure | Coordinated with reporter, typically 90 days after triage |

We follow coordinated disclosure. We will not take legal action against researchers who report in good faith and do not disclose publicly before the agreed date.

---

## Severity Classification

| Severity | Definition |
|----------|------------|
| Critical | Real credential exposed to agent process or extracted from the OS keyring by an unintended process |
| High | Masked token can be used to reconstruct or exchange for real credential |
| Medium | Keyring access or proxy accessible beyond intended scope without additional privileges |
| Low | Information disclosure that does not directly expose credentials |

---

## Supported Versions

Only the latest release receives security fixes. If you are running an older version, upgrade before reporting.

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |
| < Latest | No       |

---

## Acknowledgements

Security researchers who report valid vulnerabilities will be credited in release notes unless they prefer to remain anonymous.
