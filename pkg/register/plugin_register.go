package register

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/sync/errgroup"
)

type PluginRegisterBuilder struct {
	plugins map[string]PluginAPI
}

func NewPluginRegisterBuilder() *PluginRegisterBuilder {
	return &PluginRegisterBuilder{
		plugins: make(map[string]PluginAPI),
	}
}

func (pr *PluginRegisterBuilder) Add(name, path string) *PluginRegisterBuilder {
	pr.plugins[name] = PluginAPI{}

	return pr
}

func (pr *PluginRegisterBuilder) Build(ctx context.Context) (*PluginRegister, error) {
	clients := make(map[string]*pluginClientWrapper, 0)
	errgroup, _ := errgroup.WithContext(ctx)

	for name, p := range pr.plugins {
		name, p := name, p

		errgroup.Go(func() error {
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "BASIC_PLUGIN",
					MagicCookieValue: "char",
				},
				Logger: hclog.New(&hclog.LoggerOptions{
					Name:   "char",
					Output: os.Stdout,
					Level:  hclog.Debug,
				}),
				Cmd: exec.Command("sh", "-c", fmt.Sprintf(
					"(cd ./.char/plugins/%s; go run plugins/%s/main.go)",
					name,
					name,
				)),
				Plugins: map[string]plugin.Plugin{
					name: &p,
				},
			})

			rpcClient, err := client.Client()
			if err != nil {
				return err
			}

			raw, err := rpcClient.Dispense(name)
			if err != nil {
				return err
			}

			pluginApi, ok := raw.(Plugin)
			if !ok {
				return errors.New("could not cast as plugin")
			}

			clients[name] = &pluginClientWrapper{
				plugin: pluginApi,
				client: client,
			}

			return nil
		})
	}

	err := errgroup.Wait()
	if err != nil {
		return nil, err
	}

	return &PluginRegister{
		clients: clients,
	}, nil
}

// ---

type pluginClientWrapper struct {
	plugin Plugin
	client *plugin.Client
}

func (pcw *pluginClientWrapper) Close() {
	pcw.client.Kill()
}

// ---

type PluginRegister struct {
	clients map[string]*pluginClientWrapper
}

func (pr *PluginRegister) Close() error {
	errgroup, _ := errgroup.WithContext(context.Background())

	for _, c := range pr.clients {
		c := c

		errgroup.Go(func() error {
			c.Close()
			return nil
		})
	}

	if err := errgroup.Wait(); err != nil {
		return err
	}
	return nil
}

func (pr *PluginRegister) About(ctx context.Context) (map[string]string, error) {
	list := make(map[string]string, len(pr.clients))

	errgroup, ctx := errgroup.WithContext(ctx)

	for n, c := range pr.clients {
		n, c := n, c
		errgroup.Go(func() error {
			about := c.plugin.About()

			list[n] = about
			return nil
		})
	}

	if err := errgroup.Wait(); err != nil {
		return nil, err
	}

	return list, nil
}
