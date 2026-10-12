package cli

import "github.com/spf13/cobra"

func newGiveUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "giveup",
		Short: "Muestra la causa y cuenta el desafío como no resuelto",
		Args:  cobra.NoArgs,
		RunE:  notImplemented,
	}
}
