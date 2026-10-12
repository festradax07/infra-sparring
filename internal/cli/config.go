package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/festradax07/infra-sparring/internal/config"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Edita y valida la configuración",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "edit",
			Short: "Abre la configuración con $EDITOR y la valida al guardar",
			Args:  cobra.NoArgs,
			RunE:  notImplemented,
		},
		newConfigCheckCommand(),
	)

	return cmd
}

// newConfigCheckCommand validates the merged configuration (RF-25).
func newConfigCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Valida la configuración sin arrancar un desafío",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			paths, err := config.DefaultPaths()
			if err != nil {
				return err
			}
			if _, err := config.Load(paths); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "La configuración es válida.")
			return nil
		},
	}
}
