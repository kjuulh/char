package charcontext

import (
	"context"
	"log"

	"git.front.kjuulh.io/kjuulh/char/pkg/plugins/provider"
	"git.front.kjuulh.io/kjuulh/char/pkg/register"
	"git.front.kjuulh.io/kjuulh/char/pkg/schema"
)

type CharContext struct {
	contextPath    string
	pluginRegister *register.PluginRegister
	schema         *schema.CharSchema
}

func NewCharContext(ctx context.Context) (*CharContext, error) {
	localPath, err := FindLocalRoot(ctx)
	if err != nil {
		return nil, err
	}
	gpp := provider.NewGitPluginProvider()

	s, err := schema.ParseFile(ctx, ".char.yml")
	if err != nil {
		return nil, err
	}

	plugins, err := s.GetPlugins(ctx)
	if err != nil {
		return nil, err
	}

	err = gpp.FetchPlugins(ctx, s.Registry, plugins)
	if err != nil {
		return nil, err
	}

	builder := register.NewPluginRegisterBuilder()

	for name, plugin := range plugins {
		builder = builder.Add(name.Hash(), plugin.Opts.Path)
	}

	r, err := builder.Build(ctx)
	if err != nil {
		return nil, err
	}

	return &CharContext{
		contextPath:    localPath,
		pluginRegister: r,
		schema:         s,
	}, nil
}

func (cc *CharContext) Close() {
	if err := cc.pluginRegister.Close(); err != nil {
		log.Fatal(err)
	}
}

func (cc *CharContext) About(ctx context.Context) ([]register.AboutItem, error) {
	return cc.pluginRegister.About(ctx)
}
