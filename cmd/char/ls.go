package char

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/char/pkg/plugins/provider"
	"git.front.kjuulh.io/kjuulh/char/pkg/register"
	"git.front.kjuulh.io/kjuulh/char/pkg/schema"
	"github.com/spf13/cobra"
)

func NewLsCommand() *cobra.Command {
	gpp := provider.NewGitPluginProvider()

	cmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			s, err := schema.ParseFile(ctx, ".char.yml")
			if err != nil {
				return err
			}

			plugins, err := s.GetPlugins(ctx)
			if err != nil {
				return err
			}

			err = gpp.FetchPlugins(ctx, s.Registry, plugins)
			if err != nil {
				return err
			}

			builder := register.NewPluginRegisterBuilder()

			for name, plugin := range plugins {
				builder = builder.Add(name.Hash(), plugin.Opts.Path)
			}

			r, err := builder.Build(ctx)
			if err != nil {
				return err
			}
			defer r.Close()

			about, err := r.About(ctx)
			if err != nil {
				return err
			}

			for _, a := range about {
				fmt.Printf("plugin: %s\n", a.Name)
				fmt.Printf("\tversion: %s\n", a.Version)
				fmt.Printf("\tabout: %s\n", a.About)
			}

			return nil
		},
	}

	return cmd
}
