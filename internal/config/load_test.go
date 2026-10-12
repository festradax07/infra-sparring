package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// writeFile creates a file named name with content in a temp dir and returns its path.
func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

func missingFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "does-not-exist.yaml")
}

func TestRF24_UsesDefaultsWhenNoFileExists(t *testing.T) {
	cfg, err := Load(Paths{Global: missingFile(t), Local: missingFile(t)})

	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !reflect.DeepEqual(cfg, Default()) {
		t.Errorf("got %+v, want defaults %+v", cfg, Default())
	}
}

func TestRF24_LocalFileOverridesGlobalFile(t *testing.T) {
	global := writeFile(t, "config.yaml", `
pistas: 1
examen:
  minutos_por_desafio: 10
nodos: [ubuntu, alma]
`)
	local := writeFile(t, "spar.yaml", `
pistas: 0
nodos: [alma]
`)

	cfg, err := Load(Paths{Global: global, Local: local})
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Hints != 0 {
		t.Errorf("pistas = %d, want 0 from spar.yaml", cfg.Hints)
	}
	if cfg.Exam.MinutesPerChallenge != 10 {
		t.Errorf("minutos_por_desafio = %d, want 10 from the global file", cfg.Exam.MinutesPerChallenge)
	}
	if !reflect.DeepEqual(cfg.Nodes, []string{"alma"}) {
		t.Errorf("nodos = %v, want [alma]: a list in spar.yaml replaces the global one", cfg.Nodes)
	}
	if cfg.LLM.Provider != ProviderAnthropic {
		t.Errorf("llm.provider = %q, want the default %q", cfg.LLM.Provider, ProviderAnthropic)
	}
}

func TestRF24_EmptyFileKeepsDefaults(t *testing.T) {
	cfg, err := Load(Paths{Global: writeFile(t, "config.yaml", ""), Local: missingFile(t)})

	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !reflect.DeepEqual(cfg, Default()) {
		t.Errorf("got %+v, want defaults", cfg)
	}
}

func TestRF25_InvalidValueReportsFileLineAndValidValues(t *testing.T) {
	local := writeFile(t, "spar.yaml", "pack: lfcs\nmodo: incidnete\n")

	_, err := Load(Paths{Global: missingFile(t), Local: local})

	assertErrorContains(t, err, local+":2: modo:", `"incidnete"`, "examen, incidente, mixto")
}

func TestRF25_PointsToTheFileThatSetTheValueLast(t *testing.T) {
	global := writeFile(t, "config.yaml", "pistas: 1\n")
	local := writeFile(t, "spar.yaml", "\n\npistas: 5\n")

	_, err := Load(Paths{Global: global, Local: local})

	assertErrorContains(t, err, local+":3: pistas:")
}

func TestRF25_ReportsEveryProblemAtOnce(t *testing.T) {
	local := writeFile(t, "spar.yaml", "modo: x\npistas: 9\n")

	_, err := Load(Paths{Global: missingFile(t), Local: local})

	assertErrorContains(t, err, ":1: modo:", ":2: pistas:")
}

func TestRF25_RejectsUnknownKeysWithLine(t *testing.T) {
	local := writeFile(t, "spar.yaml", "llm:\n  provider: anthropic\n  modle: x\n")

	_, err := Load(Paths{Global: missingFile(t), Local: local})

	assertErrorContains(t, err, local+":3:", `"llm.modle"`, "model")
}

func TestRF25_RejectsWrongValueType(t *testing.T) {
	local := writeFile(t, "spar.yaml", "pistas: dos\n")

	_, err := Load(Paths{Global: missingFile(t), Local: local})

	assertErrorContains(t, err, local, "line 1")
}

func TestRF25_RejectsMalformedYAML(t *testing.T) {
	local := writeFile(t, "spar.yaml", "modo: [examen\n")

	_, err := Load(Paths{Global: missingFile(t), Local: local})

	assertErrorContains(t, err, local, "YAML mal formado")
}

func TestDefaultPaths_UsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg")

	paths, err := DefaultPaths()

	if err != nil {
		t.Fatalf("default paths: %v", err)
	}
	if paths.Global != "/tmp/xdg/spar/config.yaml" || paths.Local != "spar.yaml" {
		t.Errorf("got %+v", paths)
	}
}

func assertErrorContains(t *testing.T, err error, fragments ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("got nil error, want an error")
	}
	for _, fragment := range fragments {
		if !strings.Contains(err.Error(), fragment) {
			t.Errorf("error %q does not contain %q", err, fragment)
		}
	}
}
