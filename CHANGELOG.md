# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.3, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.2, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.1, 2026-10-16

- Added the `go-mistral` agent skill. It installs with `claude plugin install go-mistral@psyb0t` in Claude Code, `codex plugin add go-mistral@psyb0t` in Codex and `openclaw skills install @psyb0t/go-mistral` from ClawHub.

## v0.1.0, 2026-10-16

- Added `mistral`, a REST client generated with oapi-codegen from Mistral AI's OpenAPI document (spec version 1.0.0).
- `make generate` downloads the current spec from Mistral AI's platform-docs-public repository and regenerates the client. Boolean `const` values are dropped before generation, since one of them collides with a component schema name. The spec is not committed.
