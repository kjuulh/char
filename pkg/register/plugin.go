package register

import "context"

type AboutCommand struct {
	Name     string   `json:"name" yaml:"name"`
	Args     []string `json:"args" yaml:"args"`
	Required []string `json:"required" yaml:"required"`
}

type About struct {
	Name     string          `json:"name"`
	Version  string          `json:"version"`
	About    string          `json:"about"`
	Vars     []string        `json:"vars"`
	Commands []*AboutCommand `json:"commands"`
}

type DoCommand struct {
	CommandName string            `json:"commandName"`
	Args        map[string]string `json:"args"`
}

type Plugin interface {
	About(ctx context.Context) (*About, error)
	Do(ctx context.Context, cmd *DoCommand) error
}

const PluginKey = "plugin"
