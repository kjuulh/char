package main

import (
	"context"
	"log"

	"git.front.kjuulh.io/kjuulh/char/pkg/register"
)

type GoCliPlugin struct {
}

// About implements register.Plugin
func (*GoCliPlugin) About() string {
	return "gocli"
}

var _ register.Plugin = &GoCliPlugin{}

func main() {
	if err := register.
		NewPluginBuilder(
			"gocli",
			&GoCliPlugin{},
		).
		Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
