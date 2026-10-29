// Command patch-spec fixes the parts of Mistral AI's spec that oapi-codegen
// cannot generate from, after scripts/patch-spec.sh has dropped the boolean
// consts. It rewrites api/openapi.patched.yaml in place.
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	specPath   = "api/openapi.patched.yaml"
	outputMode = 0o644
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("read spec: %w", err)
	}

	var spec map[string]any

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("decode spec: %w", err)
	}

	if err := hoist(spec); err != nil {
		return fmt.Errorf("hoist $defs: %w", err)
	}

	nameDuplicateParams(spec)

	data, err = yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("encode spec: %w", err)
	}

	if err := os.WriteFile(specPath, data, outputMode); err != nil {
		return fmt.Errorf("write spec: %w", err)
	}

	return nil
}
