// Command spar is the entry point of the CLI. It only wires the root command;
// all behavior lives in internal packages.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/festradax07/infra-sparring/internal/cli"
)

func main() {
	// Cancel the context on Ctrl+C or SIGTERM so long VM or network
	// operations can stop cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cli.NewRootCommand().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		stop()
		os.Exit(1)
	}
}
