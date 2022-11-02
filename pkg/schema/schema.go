package schema

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CharSchema struct {
	Registry string            `json:"registry" yaml:"registry"`
	Plugins  CharSchemaPlugins `json:"plugins" yaml:"plugins"`
}

func ParseFile(ctx context.Context, path string) (*CharSchema, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("could not parse file, as it is not found or permitted: %s", path)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return Parse(file)
}

func Parse(content []byte) (*CharSchema, error) {
	var schema CharSchema
	if err := yaml.Unmarshal(content, &schema); err != nil {
		return nil, fmt.Errorf("could not deserialize yaml into CharSchema: %w", err)
	}

	return &schema, nil
}

func (cs *CharSchema) GetPlugins(ctx context.Context) (CharSchemaPlugins, error) {
	plugins := make(map[CharSchemaPluginName]*CharSchemaPlugin, len(cs.Plugins))
	for n, plugin := range cs.Plugins {
		po, err := n.Get()
		if err != nil {
			return nil, err
		}
		plugin.Opts = po
		plugins[n] = plugin
	}
	return plugins, nil
}
