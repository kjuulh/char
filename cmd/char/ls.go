package char

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/char/pkg/register"
	"git.front.kjuulh.io/kjuulh/char/pkg/schema"
	"github.com/spf13/cobra"
)

func NewLsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			_, err := schema.ParseFile(ctx, ".char.yml")
			if err != nil {
				return err
			}

			r, err := register.
				NewPluginRegisterBuilder().
				Add("gocli", "plugins/gocli/main.go").
				Add("rust", "plugins/rust/main.go").
				Build(ctx)
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
