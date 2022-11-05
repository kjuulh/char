package register

import (
	"context"
	"encoding/json"
	"net/rpc"
)

type PluginClient struct {
	client *rpc.Client
}

type DoRequest struct {
	CommandName string            `json:"commandName"`
	Args        map[string]string `json:"args"`
}

// Do implements Plugin
func (pc *PluginClient) Do(ctx context.Context, commandName string, args map[string]string) error {
	req := &DoRequest{
		CommandName: commandName,
		Args:        args,
	}
	//doReq, err := json.Marshal(req)
	//if err != nil {
	//	return err
	//}
	err := pc.client.Call("Plugin.Do", req, new(any))
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
