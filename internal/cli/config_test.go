package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRF25_ConfigCheckAcceptsValidConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Chdir(t.TempDir())
	var out bytes.Buffer
	root := NewRootCommand()
	root.SetOut(&out)
	root.SetArgs([]string{"config", "check"})

	err := root.Execute()

	if err != nil {
		t.Fatalf("config check: %v", err)
	}
	if !strings.Contains(out.String(), "válida") {
		t.Errorf("output %q does not confirm the config is valid", out.String())
	}
}

func TestRF25_ConfigCheckRejectsInvalidLocalFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Chdir(t.TempDir())
	if err := os.WriteFile("spar.yaml", []byte("modo: incidnete\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	root := NewRootCommand()
	root.SetArgs([]string{"config", "check"})

	err := root.Execute()

	if err == nil || !strings.Contains(err.Error(), "spar.yaml:1: modo:") {
		t.Fatalf("got error %v, want one pointing to spar.yaml:1", err)
	}
}
