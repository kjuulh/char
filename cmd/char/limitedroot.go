package char

import (
	"github.com/spf13/cobra"
)

func NewLimitedCharCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "char",
	}

	return cmd
}
