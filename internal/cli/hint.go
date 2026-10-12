package cli

import "github.com/spf13/cobra"

func newHintCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "hint",
		Short: "Da la siguiente pista (máximo dos)",
		Args:  cobra.NoArgs,
		RunE:  notImplemented,
	}
}
