package register

import (
	"log"
	"net/rpc"
)

type PluginClient struct {
	client *rpc.Client
}

var _ Plugin = &PluginClient{}

func (pc *PluginClient) About() string {
	var resp string
	err := pc.client.Call("Plugin.About", new(any), &resp)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
