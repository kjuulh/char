package register

import (
	"context"
	"encoding/json"
)

type PluginServer struct {
	Impl Plugin
}

func (ps *PluginServer) Do(args *DoRequest, resp *string) error {
	//rawReq, ok := args.(string)
	//if !ok {
	//	return errors.New("args is not a string")
	//}

	//var doReq DoRequest
	//if err := json.Unmarshal([]byte(rawReq), &doReq); err != nil {
	//	return err
	//}

	if err := ps.Impl.Do(context.Background(), args.CommandName, args.Args); err != nil {
		return err
	}
	*resp = ""

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
