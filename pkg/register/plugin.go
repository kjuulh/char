package register

import "context"

type About struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	About   string `json:"about"`
}

type Plugin interface {
	About(ctx context.Context) (*About, error)
}

const PluginKey = "plugin"
