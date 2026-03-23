# Security Policy

## Scope

TokenShim is a local security tool for managing and auditing API credentials in AI agent environments. Its security guarantee is specific: **real API credentials must never be exposed in files, environment variables, or agent processes.** All vulnerability reports are evaluated against this guarantee first.

### In scope

- Secret detection bypass — patterns that fail to detect exposed credentials
- False negatives in the doctor scanner (missed secrets across supported types)
- Credential exposure introduced by TokenShim itself during scanning or reporting
- Report output leaking unredacted secret values
- OS keyring exposure — secrets readable by unintended processes
- Future proxy/injection logic flaws — cases where real credentials are forwarded to unintended hosts
- TLS downgrade or MITM conditions on the loopback proxy interface (when implemented)
- Dependency-level supply chain compromises affecting credential handling

### Out of scope

- Attacks requiring root/kernel-level access to the host machine
- Social engineering or physical access scenarios
- Issues in upstream AI services (OpenAI, Anthropic, etc.) — report those to the respective vendor
- Theoretical vulnerabilities without a working proof-of-concept

---

## Execution Model

TokenShim runs entirely on the local machine. It has no cloud component, makes no outbound connections to GearSec infrastructure, and collects no telemetry. The attack surface is limited to:

1. Files and environment variables being scanned by doctor mode
2. Report output (all secret values are redacted as `sk-ab****xyz`)
3. The OS keyring (Keychain on macOS, Secret Service on Linux, Windows Credential Manager on Windows)
4. The loopback proxy interface (in progress)

---

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Send a report to: **security@gearsec.io**

Include:
- A clear description of the vulnerability and the affected component
- The version of TokenShim and OS/architecture where you reproduced it
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
| Critical | Unredacted secret exposed in report output or extracted from OS keyring by an unintended process |
| High | Doctor scanner consistently fails to detect a supported secret type |
| Medium | Information disclosure that does not directly expose credentials |
| Low | Minor scanner inaccuracy or edge case with limited real-world impact |

---

## Supported Versions

Only the latest release receives security fixes.

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |
| < Latest | No       |

---

## Acknowledgements

Security researchers who report valid vulnerabilities will be credited in release notes unless they prefer to remain anonymous.
