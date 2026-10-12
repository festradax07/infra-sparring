package cli

import "github.com/spf13/cobra"

func newReadyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ready",
		Short: "Verifica tu solución y genera el informe",
		Args:  cobra.NoArgs,
		RunE:  notImplemented,
	}
}
