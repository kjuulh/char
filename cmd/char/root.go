package char

import "github.com/spf13/cobra"

func NewCharCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "char",
	}

	cmd.AddCommand(
		NewLsCommand(),
	)

	return cmd
}
