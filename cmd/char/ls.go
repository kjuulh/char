package char

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/char/pkg/register"
	"github.com/spf13/cobra"
)

func NewLsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			r, err := register.
				NewPluginRegisterBuilder().
				Add("gocli", "").
				Build(ctx)
			if err != nil {
				return err
			}

			about, err := r.About(ctx)
			if err != nil {
				return err
			}

			for plugin, aboutText := range about {
				fmt.Printf("plugin: %s\n", plugin)
				fmt.Printf("\tabout: %s\n", aboutText)
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}
