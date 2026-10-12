package cli

import "github.com/spf13/cobra"

func newAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Maneja la API key del provider",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "login",
			Short: "Valida la API key y la guarda en el llavero del sistema",
			Args:  cobra.NoArgs,
			RunE:  notImplemented,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Muestra de dónde se está tomando la API key",
			Args:  cobra.NoArgs,
			RunE:  notImplemented,
		},
		&cobra.Command{
			Use:   "logout",
			Short: "Borra la API key del llavero",
			Args:  cobra.NoArgs,
			RunE:  notImplemented,
		},
	)

	return cmd
}
