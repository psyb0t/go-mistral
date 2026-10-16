# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `mistral`, a REST client generated with oapi-codegen from Mistral AI's OpenAPI document (spec version 1.0.0).
- `make generate` downloads the current spec from Mistral AI's platform-docs-public repository and regenerates the client. Boolean `const` values are dropped before generation, since one of them collides with a component schema name. The spec is not committed.
