package main

import (
	"context"
	"log"

	"git.front.kjuulh.io/kjuulh/char/pkg/register"
)

type GoCliPlugin struct {
}

func (*GoCliPlugin) About(ctx context.Context) (*register.About, error) {
	return &register.About{
		Name:    "rust",
		Version: "v0.0.1",
		About:   "golang cli provides a set of actions and presets supporting golang development",
	}, nil
}

var _ register.Plugin = &GoCliPlugin{}

func main() {
	if err := register.
		NewPluginBuilder(
			"rust",
			&GoCliPlugin{},
		).
		Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
