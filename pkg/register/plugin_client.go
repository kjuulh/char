package register

import (
	"context"
	"encoding/json"
	"net/rpc"
)

type PluginClient struct {
	client *rpc.Client
}

// Do implements Plugin
func (pc *PluginClient) Do(ctx context.Context, cmd *DoCommand) error {
	err := pc.client.Call("Plugin.Do", cmd, new(string))
	if err != nil {
		return err
	}

	return nil
}

var _ Plugin = &PluginClient{}

func (pc *PluginClient) About(ctx context.Context) (*About, error) {
	var resp string
	err := pc.client.Call("Plugin.About", new(any), &resp)
	if err != nil {
		return nil, err
	}

	var about About
	err = json.Unmarshal([]byte(resp), &about)
	if err != nil {
		return nil, err
	}

	return &about, nil
}
