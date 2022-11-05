package main

import (
	"context"
	"log"

	"git.front.kjuulh.io/kjuulh/char/pkg/register"
	"github.com/hashicorp/go-hclog"
)

type GoCliPlugin struct{}

// Do implements register.Plugin
func (*GoCliPlugin) Do(ctx context.Context, cmd *register.DoCommand) error {
	hclog.L().Info("received command: %s", cmd.CommandName)
	return nil
}

func (*GoCliPlugin) About(ctx context.Context) (*register.About, error) {
	return &register.About{
		Name:    "rust",
		Version: "v0.0.1",
		About:   "rust cli provides a set of actions and presets supporting rust development",
	}, nil
}

var _ register.Plugin = &GoCliPlugin{}

func main() {
	if err := register.
		NewPluginBuilder(
			&GoCliPlugin{},
		).
		Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
