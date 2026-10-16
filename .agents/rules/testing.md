# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in package `mistral_test`, and use only the exported API.
- **TST2** Tests never call Mistral AI. They point the client at an `httptest` server and answer with a canned status and body.
- **TST3** Every test calls `t.Parallel()`. Assert the request method, path, auth header and body, then the decoded `JSONnnn` field.
- **TST4** `example_test.go` compiles the README usage example as `Example`. Keep the two in sync.
- **TST5** `make test` runs `go test -race` and must pass.
