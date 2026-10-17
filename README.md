# go-mistral

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-mistral.svg)](https://pkg.go.dev/github.com/psyb0t/go-mistral)
[![CI](https://github.com/psyb0t/go-mistral/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-mistral/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-mistral/badges/version.svg)](https://github.com/psyb0t/go-mistral/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-mistral/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-mistral/badges/importers.svg)](https://github.com/psyb0t/go-mistral/blob/badges/importers.md)

Unofficial Go client for Mistral AI. It covers the [Mistral AI REST API](https://docs.mistral.ai/api/), `https://api.mistral.ai`, and is generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Mistral AI publishes.

The package lives in the repo root and is named `mistral`. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-mistral
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Method names come from the spec's operation IDs, so listing models is `ListModelsV1ModelsGetWithResponse` and chat completions is `ChatCompletionV1ChatCompletionsPostWithResponse`.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-mistral"
)

func main() {
	auth := mistral.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("MISTRAL_API_KEY"))

		return nil
	})

	client, err := mistral.NewClientWithResponses("https://api.mistral.ai", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ListModelsV1ModelsGetWithResponse(
		context.Background(),
		&mistral.ListModelsV1ModelsGetParams{},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is compiled as `Example` in `example_test.go`, so it stays in step with the generated client.

Fields the spec types as one of several shapes, such as a message's `content`, are union types with `From...` and `As...` methods. A user message with plain text content:

```go
var content mistral.UserMessage_Content
if err := content.FromUserMessageContent0("hello"); err != nil {
	panic(err)
}

role := mistral.UserMessageRoleUser

var message mistral.ChatCompletionRequest_Messages_Item
if err := message.FromUserMessage(mistral.UserMessage{Content: &content, Role: &role}); err != nil {
	panic(err)
}
```

## Auth

Mistral AI authenticates with a bearer API key, created in the API keys section of [La Plateforme](https://console.mistral.ai). Pass it on every request with a `RequestEditorFn`, as in the example above. A `RequestEditorFn` can also be passed to a single call as the last argument to scope it to that request.

## Agent integrations

The [skill](.agents/skills/go-mistral) teaches a coding agent how to use go-mistral when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-mistral@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-mistral@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-mistral
```

## Regenerating

Mistral AI's spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `openapi.yaml` from Mistral AI's [platform-docs-public](https://github.com/mistralai/platform-docs-public) repository.
- `scripts/patch-spec.sh` writes `api/openapi.patched.yaml` with every `const: true` and `const: false` line dropped. oapi-codegen turns a property `const` into a named type, and the type it picks for the boolean `stream` flag on the transcription request is already the name of the streaming request schema, so generation fails on a duplicate type name. Those properties become plain booleans.
- oapi-codegen runs on the patched document with `config.yaml`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Mistral AI's published API specification. The specification itself is not redistributed here. This project is not affiliated with Mistral AI.
