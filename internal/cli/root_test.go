package cli

import (
	"errors"
	"testing"
)

func TestRootCommand_RegistersPlayerCommands(t *testing.T) {
	want := []string{"start", "hint", "ready", "giveup", "config", "auth"}

	root := NewRootCommand()

	for _, name := range want {
		if _, _, err := root.Find([]string{name}); err != nil {
			t.Errorf("command %q not registered: %v", name, err)
		}
	}
}

func TestRootCommand_UnimplementedCommandsReturnError(t *testing.T) {
	root := NewRootCommand()
	root.SetArgs([]string{"start"})

	err := root.Execute()

	if !errors.Is(err, errNotImplemented) {
		t.Fatalf("got error %v, want %v", err, errNotImplemented)
	}
}
