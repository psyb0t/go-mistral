package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNameDuplicateParams(t *testing.T) {
	t.Parallel()

	input := `paths:
  /items:
    get:
      parameters:
        - {name: sort.direction, in: query}
        - {name: sort_by, in: query}
        - {name: sort_direction, in: query}
        - {name: sortDirection, in: query}
`
	want := `paths:
  /items:
    get:
      parameters:
        - {name: sort.direction, in: query}
        - {name: sort_by, in: query}
        - {name: sort_direction, in: query, x-go-name: SortDirection2}
        - {name: sortDirection, in: query, x-go-name: SortDirection3}
`

	var spec, wantSpec map[string]any
	require.NoError(t, yaml.Unmarshal([]byte(input), &spec))
	require.NoError(t, yaml.Unmarshal([]byte(want), &wantSpec))

	nameDuplicateParams(spec)

	assert.Equal(t, wantSpec, spec)
}
