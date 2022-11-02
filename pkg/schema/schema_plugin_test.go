package schema_test

import (
	"testing"

	"git.front.kjuulh.io/kjuulh/char/pkg/schema"
	"github.com/stretchr/testify/require"
)

func TestSchemaNameCanParse(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		inputString schema.CharSchemaPluginName
		expected    schema.PluginOps
	}{
		{
			name:        "default string",
			inputString: `kju123K_-ulh/someRepo-._123`,
			expected: schema.PluginOps{
				Org:            "kju123K_-ulh",
				RepositoryName: "someRepo-._123",
			},
		},
		{
			name:        "default string with path",
			inputString: `kju123K_-ulh/someRepo-._123#somepath/sometoherpath/somethridpath`,
			expected: schema.PluginOps{
				Org:            "kju123K_-ulh",
				RepositoryName: "someRepo-._123",
				Path:           "somepath/sometoherpath/somethridpath",
			},
		},
		{
			name:        "default string with version",
			inputString: `kju123K_-ulh/someRepo-._123@12l3.jk1lj`,
			expected: schema.PluginOps{
				Org:            "kju123K_-ulh",
				RepositoryName: "someRepo-._123",
				Version:        "12l3.jk1lj",
			},
		},
		{
			name:        "default string with version and path",
			inputString: `kju123K_-ulh/someRepo-._123@12l3.jk1lj#somepath/sometoherpath/somethridpath`,
			expected: schema.PluginOps{
				Org:            "kju123K_-ulh",
				RepositoryName: "someRepo-._123",
				Version:        "12l3.jk1lj",
				Path:           "somepath/sometoherpath/somethridpath",
			},
		},
		{
			name:        "default string with path and version",
			inputString: `kju123K_-ulh/someRepo-._123#somepath/sometoherpath/somethridpath@12l3.jk1lj`,
			expected: schema.PluginOps{
				Org:            "kju123K_-ulh",
				RepositoryName: "someRepo-._123",
				Version:        "12l3.jk1lj",
				Path:           "somepath/sometoherpath/somethridpath",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.inputString.Get()
			require.Equal(t, tc.expected, *actual)
		})
	}
}
