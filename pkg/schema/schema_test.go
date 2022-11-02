package schema_test

import (
	"context"
	"testing"

	"git.front.kjuulh.io/kjuulh/char/pkg/schema"
	"github.com/stretchr/testify/require"
)

func TestSchemaParse(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		input    string
		expected *schema.CharSchema
	}{
		{
			name: "with plugins",
			input: `
registry: git.front.kjuulh.io
plugins:
  "kjuulh/char#plugins/gocli": {}
  "kjuulh/char#plugins/rust": {}
`,
			expected: &schema.CharSchema{
				Registry: "git.front.kjuulh.io",
				Plugins: map[schema.CharSchemaPluginName]*schema.CharSchemaPlugin{
					"kjuulh/char#plugins/gocli": {},
					"kjuulh/char#plugins/rust":  {},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := schema.Parse([]byte(tc.input))

			require.NoError(t, err)
			require.Equal(t, tc.expected, s)
		})
	}
}

func TestGetPlugins(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		input    string
		expected schema.CharSchemaPlugins
	}{
		{
			name: "with plugins",
			input: `
registry: git.front.kjuulh.io
plugins:
  "kjuulh/char#plugins/gocli@v1.9.0": {}
  "kjuulh/char#plugins/rust": {}
`,
			expected: map[schema.CharSchemaPluginName]*schema.CharSchemaPlugin{
				"kjuulh/char#plugins/gocli@v1.9.0": {
					Opts: &schema.PluginOps{
						Org:            "kjuulh",
						RepositoryName: "char",
						Path:           "plugins/gocli",
						Version:        "v1.9.0",
					},
				},
				"kjuulh/char#plugins/rust": {
					Opts: &schema.PluginOps{
						Org:            "kjuulh",
						RepositoryName: "char",
						Path:           "plugins/rust",
						Version:        "",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := schema.Parse([]byte(tc.input))
			require.NoError(t, err)

			plugins, err := s.GetPlugins(context.Background())
			require.NoError(t, err)

			require.Equal(t, tc.expected, plugins)
		})
	}
}
