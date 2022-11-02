package schema_test

import (
	"context"
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
			actual, _ := tc.inputString.Get()
			require.Equal(t, tc.expected, *actual)
		})
	}
}

func TestPluginOpt(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		pluginOpt   schema.PluginOps
		cloneUrlOpt schema.CloneUrlOpt
		registry    string
		expected    string
	}{
		{
			name: "ssh values",
			pluginOpt: schema.PluginOps{
				Org:            "kjuulh",
				RepositoryName: "char",
				Path:           "",
				Version:        "",
			},
			cloneUrlOpt: schema.CloneUrlOpt{
				Protocol: schema.GitProtocolSsh,
				SshUser:  "git",
			},
			registry: "git.front.kjuulh.io",
			expected: "git@git.front.kjuulh.io:kjuulh/char.git",
		},
		{
			name: "https values",
			pluginOpt: schema.PluginOps{
				Org:            "kjuulh",
				RepositoryName: "char",
				Path:           "",
				Version:        "",
			},
			cloneUrlOpt: schema.CloneUrlOpt{
				Protocol: schema.GitProtocolHttps,
			},
			registry: "git.front.kjuulh.io",
			expected: "https://git.front.kjuulh.io/kjuulh/char.git",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			url, err := tc.pluginOpt.GetCloneUrl(context.Background(), tc.registry, &tc.cloneUrlOpt)

			require.NoError(t, err)
			require.Equal(t, tc.expected, url)

		})
	}

}
