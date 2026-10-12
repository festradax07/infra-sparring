// Package lab abstracts the lab VMs where challenges are validated and played.
//
// The core only depends on the Environment interface; the Vagrant adapter and
// the test fake are interchangeable implementations.
package lab

import (
	"context"
	"io"
)

// Environment controls the lab nodes. Implementations must never copy the
// script to the node: it is sent through stdin so no artifact stays on disk
// inside the VM (RF-18).
type Environment interface {
	// Snapshot saves the current state of node under name.
	Snapshot(ctx context.Context, node, name string) error
	// Restore returns node to the state saved under name.
	Restore(ctx context.Context, node, name string) error
	// Exec runs script as root on node and reports how it ended. A non-zero
	// exit code is not an error: err is reserved for failures to run at all.
	Exec(ctx context.Context, node string, script io.Reader) (ExecResult, error)
}

// ExecResult is the outcome of a script that ran on a node.
type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Succeeded reports whether the script exited with code zero.
func (r ExecResult) Succeeded() bool {
	return r.ExitCode == 0
}
