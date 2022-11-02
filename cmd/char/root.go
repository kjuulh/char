package char

import (
	"git.front.kjuulh.io/kjuulh/char/pkg/charcontext"
	"github.com/spf13/cobra"
)

func NewCharCmd(charctx *charcontext.CharContext) *cobra.Command {
	cmd := &cobra.Command{
		Use: "char",
	}

	cmd.AddCommand(
		NewLsCommand(charctx),
	)

	return cmd
}
