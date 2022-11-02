package schema

type CharSchema struct {
	Registry string            `json:"registry" yaml:"registry"`
	Plugins  CharSchemaPlugins `json:"plugins" yaml:"plugins"`
}
