// Package cli defines the spar commands. Each command lives in its own file
// and only translates user input into calls to the internal packages.
package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

// errNotImplemented is returned by commands whose behavior is not built yet.
var errNotImplemented = errors.New("todavía no está implementado")

// NewRootCommand builds the spar command tree. It creates a fresh tree on
// every call so tests never share state.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "spar",
		Short: "Tu sparring para certificaciones de Linux",
		Long: "spar genera fallas nuevas en tu laboratorio de VMs, vos las diagnosticás\n" +
			"desde la terminal y al final recibís un informe de cómo las resolviste.",
		// Errors are printed once by main; usage is noise for runtime errors.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(
		newStartCommand(),
		newHintCommand(),
		newReadyCommand(),
		newGiveUpCommand(),
		newConfigCommand(),
		newAuthCommand(),
	)

	return root
}

// notImplemented is a placeholder RunE for commands still under construction.
func notImplemented(*cobra.Command, []string) error {
	return errNotImplemented
}
