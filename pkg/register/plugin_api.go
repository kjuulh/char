package register

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type PluginAPI struct {
	Impl Plugin
}

func (pa *PluginAPI) Server(*plugin.MuxBroker) (any, error) {
	return &PluginServer{Impl: pa.Impl}, nil
}

func (*PluginAPI) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &PluginClient{client: c}, nil
}
