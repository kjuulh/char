package register

import (
	"context"
	"encoding/json"
)

type PluginServer struct {
	Impl Plugin
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
