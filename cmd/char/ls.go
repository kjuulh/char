package char

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/char/pkg/charcontext"
	"github.com/spf13/cobra"
)

func NewLsCommand(charctx *charcontext.CharContext) *cobra.Command {

	cmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			about, err := charctx.About(ctx)
			if err != nil {
				return err
			}

			for _, a := range about {
				fmt.Printf("plugin: %s\n", a.Name)
				fmt.Printf("\tversion: %s\n", a.Version)
				fmt.Printf("\tabout: %s\n", a.About)
				if len(a.Vars) > 0 {
					fmt.Println("\tVars:")
					for _, av := range a.Vars {
						fmt.Printf("\t\t%s\n", av)
					}
				}
				if len(a.Commands) > 0 {
					fmt.Println("\tCommands:")
					for _, ac := range a.Commands {
						fmt.Printf("\t\t%s\n", ac.Name)
						if len(ac.Args) == 0 {
							continue
						}
						fmt.Println("\t\tArgs")
						for _, aca := range ac.Args {
							isrequired := false
							for _, acr := range ac.Required {
								if acr == aca {
									isrequired = true
								}
							}
							if isrequired {
								fmt.Printf("\t\t\t%s: required\n", aca)
							} else {
								fmt.Printf("\t\t\t%s\n", aca)
							}
						}
					}
				}

				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}
