// Package mistral is a Go client for the Mistral AI REST API, generated with
// oapi-codegen from the OpenAPI document Mistral AI publishes.
package mistral

//go:generate bash scripts/fetch-spec.sh openapi.yaml
//go:generate bash scripts/patch-spec.sh
//go:generate go run ./internal/cmd/patch-spec
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.yaml
