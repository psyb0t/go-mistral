package main

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
)

const (
	defsKey            = "$defs"
	refKey             = "$ref"
	titleKey           = "title"
	localRefPrefix     = "#/$defs/"
	componentRefPrefix = "#/components/schemas/"
)

var (
	errNoSchemas    = errors.New("spec has no components.schemas")
	errBadDefs      = errors.New("$defs is not a map of schemas")
	errNestedDefs   = errors.New("$defs nested inside a $defs schema")
	errUnknownRef   = errors.New("reference to an undefined $defs entry")
	errNameTaken    = errors.New("schema name taken by a different component")
	errMissingTitle = errors.New("inline schema with $defs has no title")
)

// hoist moves schemas out of local $defs blocks into components/schemas. The
// spec references them as "#/$defs/<Name>", which kin-openapi resolves against
// the document root and fails on. Each definition becomes a component under
// its own name, or under
// the owning schema's name plus its own when a different component already
// uses that name. A schema with $defs inlined outside components/schemas is
// hoisted too, under its title, and replaced by a reference.
func hoist(spec map[string]any) error {
	components, _ := spec["components"].(map[string]any)

	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return errNoSchemas
	}

	// Component owners keep their place, so their names are known up front.
	for _, name := range sortedKeys(schemas) {
		owner, ok := schemas[name].(map[string]any)
		if !ok || owner[defsKey] == nil {
			continue
		}

		if err := moveDefs(owner, name, schemas); err != nil {
			return fmt.Errorf("schema %s: %w", name, err)
		}
	}

	for _, key := range sortedKeys(spec) {
		if key == "components" {
			continue
		}

		replaced, err := hoistInline(spec[key], schemas)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}

		spec[key] = replaced
	}

	return nil
}

// hoistInline replaces every schema with a $defs block below node by a
// reference to a component named after the schema's title.
func hoistInline(node any, schemas map[string]any) (any, error) {
	switch value := node.(type) {
	case map[string]any:
		if value[defsKey] != nil {
			return hoistOwner(value, schemas)
		}

		for _, key := range sortedKeys(value) {
			replaced, err := hoistInline(value[key], schemas)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}

			value[key] = replaced
		}
	case []any:
		for i, child := range value {
			replaced, err := hoistInline(child, schemas)
			if err != nil {
				return nil, fmt.Errorf("item %d: %w", i, err)
			}

			value[i] = replaced
		}
	}

	return node, nil
}

func hoistOwner(owner, schemas map[string]any) (any, error) {
	title, _ := owner[titleKey].(string)
	if title == "" {
		return nil, errMissingTitle
	}

	if err := moveDefs(owner, title, schemas); err != nil {
		return nil, fmt.Errorf("schema %s: %w", title, err)
	}

	if err := addSchema(schemas, title, owner); err != nil {
		return nil, err
	}

	return map[string]any{refKey: componentRefPrefix + title}, nil
}

// moveDefs adds owner's $defs entries to schemas and points owner's
// references at them.
func moveDefs(
	owner map[string]any, ownerName string, schemas map[string]any,
) error {
	defs, ok := owner[defsKey].(map[string]any)
	if !ok {
		return errBadDefs
	}

	delete(owner, defsKey)

	names, err := componentNames(defs, ownerName, schemas)
	if err != nil {
		return err
	}

	for _, name := range sortedKeys(defs) {
		def, err := rewriteRefs(defs[name], names)
		if err != nil {
			return fmt.Errorf("$defs %s: %w", name, err)
		}

		if err := addSchema(schemas, names[name], def); err != nil {
			return err
		}
	}

	for _, key := range sortedKeys(owner) {
		rewritten, err := rewriteRefs(owner[key], names)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}

		owner[key] = rewritten
	}

	return nil
}

// componentNames keeps each definition's name unless a different component
// already has it.
func componentNames(
	defs map[string]any, ownerName string, schemas map[string]any,
) (map[string]string, error) {
	same := make(map[string]string, len(defs))
	for name := range defs {
		same[name] = name
	}

	names := make(map[string]string, len(defs))

	for name, def := range defs {
		if hasDefs(def) {
			return nil, fmt.Errorf("$defs %s: %w", name, errNestedDefs)
		}

		rewritten, err := rewriteRefs(def, same)
		if err != nil {
			return nil, fmt.Errorf("$defs %s: %w", name, err)
		}

		names[name] = name

		existing, taken := schemas[name]
		if taken && !reflect.DeepEqual(existing, rewritten) {
			names[name] = ownerName + name
		}
	}

	return names, nil
}

func addSchema(schemas map[string]any, name string, schema any) error {
	if existing, ok := schemas[name]; ok {
		if reflect.DeepEqual(existing, schema) {
			return nil
		}

		return fmt.Errorf("%s: %w", name, errNameTaken)
	}

	schemas[name] = schema

	return nil
}

// rewriteRefs returns a copy of node with "#/$defs/<Name>" strings, in $ref
// values and discriminator mappings alike, pointing at the named components.
func rewriteRefs(node any, names map[string]string) (any, error) {
	switch value := node.(type) {
	case string:
		return rewriteRef(value, names)
	case map[string]any:
		out := make(map[string]any, len(value))

		for key, child := range value {
			rewritten, err := rewriteRefs(child, names)
			if err != nil {
				return nil, err
			}

			out[key] = rewritten
		}

		return out, nil
	case []any:
		out := make([]any, len(value))

		for i, child := range value {
			rewritten, err := rewriteRefs(child, names)
			if err != nil {
				return nil, err
			}

			out[i] = rewritten
		}

		return out, nil
	default:
		return node, nil
	}
}

func rewriteRef(value string, names map[string]string) (string, error) {
	name, ok := strings.CutPrefix(value, localRefPrefix)
	if !ok {
		return value, nil
	}

	target, ok := names[name]
	if !ok {
		return "", fmt.Errorf("%s: %w", value, errUnknownRef)
	}

	return componentRefPrefix + target, nil
}

func hasDefs(node any) bool {
	switch value := node.(type) {
	case map[string]any:
		if value[defsKey] != nil {
			return true
		}

		for _, child := range value {
			if hasDefs(child) {
				return true
			}
		}
	case []any:
		if slices.ContainsFunc(value, hasDefs) {
			return true
		}
	}

	return false
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
