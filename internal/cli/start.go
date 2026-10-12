package cli

import "github.com/spf13/cobra"

func newStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Genera un desafío nuevo y te entrega el síntoma",
		Args:  cobra.NoArgs,
		RunE:  notImplemented,
	}
}
