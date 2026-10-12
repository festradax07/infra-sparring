package labtest

import (
	"context"
	"strings"
	"testing"

	"github.com/festradax07/infra-sparring/internal/lab"
)

func TestFakeEnvironment_RestoreRequiresSnapshot(t *testing.T) {
	ctx := context.Background()
	env := &FakeEnvironment{}

	if err := env.Restore(ctx, "ubuntu", "base"); err == nil {
		t.Fatal("restore without snapshot: got nil error, want an error")
	}

	if err := env.Snapshot(ctx, "ubuntu", "base"); err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if err := env.Restore(ctx, "ubuntu", "base"); err != nil {
		t.Fatalf("restore after snapshot: %v", err)
	}
}

func TestFakeEnvironment_ExecRecordsScriptAndUsesOnExec(t *testing.T) {
	env := &FakeEnvironment{
		OnExec: func(node, script string) (lab.ExecResult, error) {
			return lab.ExecResult{ExitCode: 1, Stderr: "check failed"}, nil
		},
	}

	result, err := env.Exec(context.Background(), "alma", strings.NewReader("systemctl is-active nginx"))
	if err != nil {
		t.Fatalf("exec: %v", err)
	}

	if result.Succeeded() {
		t.Error("got success, want the failure returned by OnExec")
	}
	calls := env.Calls()
	if len(calls) != 1 || calls[0].Script != "systemctl is-active nginx" {
		t.Errorf("recorded calls = %+v, want one exec with the script", calls)
	}
}
