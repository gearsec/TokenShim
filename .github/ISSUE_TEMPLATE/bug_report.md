---
name: Bug Report
about: Something isn't working as expected
title: "[BUG] "
labels: bug
assignees: ""
---

## Description

A clear, concise description of the problem.

## Environment

- Shim version (`shim version`):
- OS and architecture:
- Go version (if built from source):
- Relevant service (OpenAI, Anthropic, etc.):

## Steps to Reproduce

1.
2.
3.

## Expected Behavior

What you expected to happen.

## Actual Behavior

What actually happened. Include any error output below.

```
paste error output here
```

## Debug Log

If applicable, attach the output of `shim --log-level debug exec -- <your command>`.

**Before attaching logs:** confirm they contain no real API credentials. Shim should never log real credentials, but verify before posting.

## Additional Context

Any other relevant context (config snippets with secrets redacted, proxy settings, etc.).
