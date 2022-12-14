package register

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type PluginBuilder struct {
	serveConfig *plugin.ServeConfig
}

func NewPluginBuilder(p Plugin) *PluginBuilder {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stderr,
		JSONFormat: false,
	})

	var pluginMap = map[string]plugin.Plugin{
		PluginKey: &PluginAPI{
			Impl: p,
		},
	}

	serveConfig := &plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "char",
		},
		Plugins: pluginMap,
		Logger:  logger,
	}

	return &PluginBuilder{
		serveConfig: serveConfig,
	}
}

func (pr *PluginBuilder) Serve(ctx context.Context) error {
	plugin.Serve(
		pr.serveConfig,
	)
	return nil
}
