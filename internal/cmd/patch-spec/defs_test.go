package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestHoist(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name: "component owner keeps equal and renames different",
			input: `components:
  schemas:
    Granularity: {type: string, enum: [word]}
    Usage: {type: object}
    Request:
      properties:
        g: {$ref: '#/$defs/Granularity'}
        u: {$ref: '#/$defs/Usage'}
      $defs:
        Granularity: {type: string, enum: [word]}
        Usage: {type: object, properties: {n: {type: integer}}}
`,
			want: `components:
  schemas:
    Granularity: {type: string, enum: [word]}
    Usage: {type: object}
    RequestUsage: {type: object, properties: {n: {type: integer}}}
    Request:
      properties:
        g: {$ref: '#/components/schemas/Granularity'}
        u: {$ref: '#/components/schemas/RequestUsage'}
`,
		},
		{
			name: "inline owner becomes a component",
			input: `paths:
  /events:
    get:
      responses:
        '200':
          content:
            text/event-stream:
              schema:
                title: Events
                properties:
                  data:
                    oneOf: [{$ref: '#/$defs/Delta'}]
                    discriminator:
                      propertyName: type
                      mapping: {delta: '#/$defs/Delta'}
                $defs:
                  Delta: {type: object}
components:
  schemas: {}
`,
			want: `paths:
  /events:
    get:
      responses:
        '200':
          content:
            text/event-stream:
              schema: {$ref: '#/components/schemas/Events'}
components:
  schemas:
    Delta: {type: object}
    Events:
      title: Events
      properties:
        data:
          oneOf: [{$ref: '#/components/schemas/Delta'}]
          discriminator:
            propertyName: type
            mapping: {delta: '#/components/schemas/Delta'}
`,
		},
		{
			name: "undefined reference",
			input: `components:
  schemas:
    Request:
      properties:
        g: {$ref: '#/$defs/Missing'}
      $defs:
        Other: {type: string}
`,
			wantErr: errUnknownRef,
		},
		{
			name: "inline owner without title",
			input: `paths:
  /x:
    schema:
      $defs:
        A: {type: string}
components:
  schemas: {}
`,
			wantErr: errMissingTitle,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var spec map[string]any
			require.NoError(t, yaml.Unmarshal([]byte(tc.input), &spec))

			err := hoist(spec)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)

			var want map[string]any
			require.NoError(t, yaml.Unmarshal([]byte(tc.want), &want))
			assert.Equal(t, want, spec)
		})
	}
}
