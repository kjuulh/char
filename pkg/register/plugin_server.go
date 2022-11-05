package register

import (
	"context"
	"encoding/json"
)

type PluginServer struct {
	Impl Plugin
}

func (ps *PluginServer) Do(args string, resp *string) error {
	var doReq DoRequest
	if err := json.Unmarshal([]byte(args), &doReq); err != nil {
		return err
	}

	if err := ps.Impl.Do(context.Background(), doReq.CommandName, doReq.Args); err != nil {
		return err
	}

	return nil
}

func (ps *PluginServer) About(args any, resp *string) error {
	r, err := ps.Impl.About(context.Background())
	if err != nil {
		return err
	}

	respB, err := json.Marshal(r)
	if err != nil {
		return err
	}

	*resp = string(respB)
	return nil
}
