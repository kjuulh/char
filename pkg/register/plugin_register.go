package register

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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
	pr.plugins[name] = PluginAPI{
		path: path,
	}

	return pr
}

func (pr *PluginRegisterBuilder) Build(ctx context.Context) (*PluginRegister, error) {
	clients := make(map[string]*pluginClientWrapper, 0)
	errgroup, _ := errgroup.WithContext(ctx)

	if err := os.MkdirAll(".char/plugins/", 0755); err != nil {
		return nil, err
	}

	for name, p := range pr.plugins {
		name, p := name, p

		errgroup.Go(func() error {
			pluginPath := fmt.Sprintf(".char/plugins/%s/dist/plugin", name)

			_, err := os.Stat(pluginPath)
			if err != nil || os.Getenv("CHAR_DEV_MODE") == "true" {
				log.Printf("building: %s", name)
				cmd := exec.Command(
					"sh",
					"-c",
					fmt.Sprintf(
						"(cd .char/plugins/%s; go build -o dist/plugin %s/main.go)",
						name,
						strings.TrimSuffix(
							strings.TrimSuffix(
								p.path,
								"main.go",
							),
							"/",
						),
					),
				)
				output, err := cmd.CombinedOutput()
				if len(output) > 0 {
					log.Println(string(output))
				}
				if err != nil {
					return fmt.Errorf("could not build plugin: %w", err)
				}
			}

			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "BASIC_PLUGIN",
					MagicCookieValue: "char",
				},
				Logger: hclog.New(&hclog.LoggerOptions{
					Name:   "char",
					Output: os.Stdout,
					Level:  hclog.Error,
				}),
				Cmd: exec.Command(
					fmt.Sprintf(
						".char/plugins/%s/dist/plugin",
						name,
					),
				),
				Plugins: map[string]plugin.Plugin{
					PluginKey: &p,
				},
			})

			rpcClient, err := client.Client()
			if err != nil {
				return err
			}

			raw, err := rpcClient.Dispense("plugin")
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

type CommandAboutItem struct {
	Name     string
	Args     []string
	Required []string
}

type CommandAboutItems []*CommandAboutItem

func FromAboutCommands(commands []*AboutCommand) CommandAboutItems {
	cai := make(CommandAboutItems, 0)
	for _, command := range commands {
		cai = append(cai, &CommandAboutItem{
			Name:     command.Name,
			Args:     command.Args,
			Required: command.Required,
		})
	}
	return cai
}

type AboutItem struct {
	Name       string
	Version    string
	About      string
	Vars       []string
	Commands   CommandAboutItems
	ClientName string
}

func (pr *PluginRegister) About(ctx context.Context) ([]AboutItem, error) {
	list := make([]AboutItem, 0)

	errgroup, ctx := errgroup.WithContext(ctx)

	for name, c := range pr.clients {
		name, c := name, c
		errgroup.Go(func() error {
			about, err := c.plugin.About(ctx)
			if err != nil {
				return err
			}

			list = append(list, AboutItem{
				Name:       about.Name,
				Version:    about.Version,
				About:      about.About,
				Vars:       about.Vars,
				Commands:   FromAboutCommands(about.Commands),
				ClientName: name,
			})
			return nil
		})
	}

	if err := errgroup.Wait(); err != nil {
		return nil, err
	}

	return list, nil
}

func (pr *PluginRegister) Do(ctx context.Context, clientName string, commandName string, args map[string]string) error {
	errgroup, ctx := errgroup.WithContext(ctx)

	client, ok := pr.clients[clientName]
	if !ok {
		return fmt.Errorf("plugin was not found: %s", clientName)
	}

	errgroup.Go(func() error {
		return client.plugin.Do(ctx, commandName, args)
	})

	return nil
}
