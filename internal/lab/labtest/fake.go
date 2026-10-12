// Package labtest provides an in-memory lab.Environment for unit tests, so
// they run without VMs.
package labtest

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/festradax07/infra-sparring/internal/lab"
)

// Compile-time check that FakeEnvironment satisfies lab.Environment.
var _ lab.Environment = (*FakeEnvironment)(nil)

// ExecFunc decides the result of a script run on node.
type ExecFunc func(node, script string) (lab.ExecResult, error)

// FakeEnvironment records every call and answers Exec through OnExec.
// Restore fails for snapshots that were never taken, like the real adapter.
type FakeEnvironment struct {
	// OnExec answers Exec calls. When nil, every script succeeds.
	OnExec ExecFunc

	mu        sync.Mutex
	snapshots map[string]bool
	calls     []Call
}

// Call is one recorded operation on the fake.
type Call struct {
	Op     string // "snapshot", "restore" or "exec"
	Node   string
	Name   string // snapshot name, empty for exec
	Script string // script content, empty for snapshot and restore
}

// Snapshot records the snapshot so a later Restore can find it.
func (f *FakeEnvironment) Snapshot(_ context.Context, node, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.snapshots == nil {
		f.snapshots = make(map[string]bool)
	}
	f.snapshots[snapshotKey(node, name)] = true
	f.calls = append(f.calls, Call{Op: "snapshot", Node: node, Name: name})
	return nil
}

// Restore fails if the snapshot does not exist.
func (f *FakeEnvironment) Restore(_ context.Context, node, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls = append(f.calls, Call{Op: "restore", Node: node, Name: name})
	if !f.snapshots[snapshotKey(node, name)] {
		return fmt.Errorf("restore %s: snapshot %q not found", node, name)
	}
	return nil
}

// Exec reads the whole script and delegates the result to OnExec.
func (f *FakeEnvironment) Exec(_ context.Context, node string, script io.Reader) (lab.ExecResult, error) {
	content, err := io.ReadAll(script)
	if err != nil {
		return lab.ExecResult{}, fmt.Errorf("exec %s: read script: %w", node, err)
	}

	f.mu.Lock()
	f.calls = append(f.calls, Call{Op: "exec", Node: node, Script: string(content)})
	onExec := f.OnExec
	f.mu.Unlock()

	if onExec == nil {
		return lab.ExecResult{}, nil
	}
	return onExec(node, string(content))
}

// Calls returns a copy of the operations recorded so far, in order.
func (f *FakeEnvironment) Calls() []Call {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]Call(nil), f.calls...)
}

func snapshotKey(node, name string) string {
	return node + "/" + name
}
