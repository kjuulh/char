package char

import (
	"context"
	"log"

	"git.front.kjuulh.io/kjuulh/char/pkg/charcontext"
	"github.com/spf13/cobra"
)

type RequiredArg struct {
	Required bool
	Value    string
}

func NewDoCommand(charctx *charcontext.CharContext) *cobra.Command {
	cmd := &cobra.Command{
		Use: "do",
	}
	about, err := charctx.About(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range about {
		for _, c := range a.Commands {
			requiredArgs := make(map[string]*RequiredArg, len(c.Args))
			for _, arg := range c.Args {
				requiredArgs[arg] = &RequiredArg{
					Required: false,
				}
			}
			for _, required := range c.Required {
				if _, ok := requiredArgs[required]; ok {
					requiredArg := requiredArgs[required]
					requiredArg.Required = true
				}
			}

			doCmd := &cobra.Command{
				Use: c.Name,
				RunE: func(cmd *cobra.Command, args []string) error {
					if err := cmd.ParseFlags(args); err != nil {
						return err
					}

					if err := charctx.Do(cmd.Context(), a.Name, c.Name); err != nil {
						return err
					}

					return nil
				},
			}

			for argName, argValue := range requiredArgs {
				doCmd.PersistentFlags().StringVar(&argValue.Value, argName, "", "")
				if argValue.Required {
					doCmd.MarkPersistentFlagRequired(argName)
				}
			}

			cmd.AddCommand(doCmd)
		}
	}

	return cmd
}
