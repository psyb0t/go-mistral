package main

import (
	"strconv"
	"strings"
	"unicode"
)

const goNameKey = "x-go-name"

// nameDuplicateParams gives a parameter its own Go field name when an earlier
// parameter of the same operation maps to the same one. The list endpoints
// take both `sort.direction` and its alias `sort_direction`, which oapi-codegen
// would both name SortDirection. A later duplicate gets a numeric suffix,
// starting at 2.
func nameDuplicateParams(spec map[string]any) {
	paths, _ := spec["paths"].(map[string]any)

	for _, path := range sortedKeys(paths) {
		item, _ := paths[path].(map[string]any)

		for _, method := range sortedKeys(item) {
			operation, _ := item[method].(map[string]any)
			params, _ := operation["parameters"].([]any)
			nameParams(params)
		}
	}
}

func nameParams(params []any) {
	seen := make(map[string]int, len(params))

	for _, raw := range params {
		param, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		name, _ := param["name"].(string)
		if name == "" || param[goNameKey] != nil {
			continue
		}

		goName := goFieldName(name)

		seen[goName]++
		if count := seen[goName]; count > 1 {
			param[goNameKey] = goName + strconv.Itoa(count)
		}
	}
}

// goFieldName approximates oapi-codegen's name for a parameter: the parts
// between non-alphanumeric characters, each with an upper-case first letter.
func goFieldName(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	var b strings.Builder

	for _, part := range parts {
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}

	return b.String()
}
