package config

import (
	"strings"
	"testing"
)

func TestRF25_DefaultsAreValid(t *testing.T) {
	if problems := Default().Validate(); len(problems) > 0 {
		t.Fatalf("defaults have problems: %+v", problems)
	}
}

func TestRF25_ValidateRules(t *testing.T) {
	tests := []struct {
		name     string
		modify   func(*Config)
		wantPath string
	}{
		{"unknown mode", func(c *Config) { c.Mode = "incidnete" }, "modo"},
		{"empty pack", func(c *Config) { c.Pack = " " }, "pack"},
		{"unknown provider", func(c *Config) { c.LLM.Provider = "openai" }, "llm.provider"},
		{"empty model", func(c *Config) { c.LLM.Model = "" }, "llm.model"},
		{"api key instead of variable name", func(c *Config) { c.LLM.APIKeyEnv = "sk-ant-123" }, "llm.api_key_env"},
		{"initial level too high", func(c *Config) { c.Difficulty.Initial = "11" }, "dificultad.inicial"},
		{"initial level not a number", func(c *Config) { c.Difficulty.Initial = "alto" }, "dificultad.inicial"},
		{"unknown ceiling", func(c *Config) { c.Difficulty.Ceiling = "max" }, "dificultad.techo"},
		{"zero exam minutes", func(c *Config) { c.Exam.MinutesPerChallenge = 0 }, "examen.minutos_por_desafio"},
		{"zero mixed tasks", func(c *Config) { c.Mixed.Tasks = 0 }, "mixto.tareas"},
		{"no nodes", func(c *Config) { c.Nodes = nil }, "nodos"},
		{"repeated node", func(c *Config) { c.Nodes = []string{"alma", "alma"} }, "nodos"},
		{"too many hints", func(c *Config) { c.Hints = 5 }, "pistas"},
		{"negative hints", func(c *Config) { c.Hints = -1 }, "pistas"},
		{"unknown language", func(c *Config) { c.Language = "fr" }, "idioma"},
		{"empty reports dir", func(c *Config) { c.Reports.Dir = "" }, "informes.carpeta"},
		{"no report formats", func(c *Config) { c.Reports.Formats = nil }, "informes.formato"},
		{"unknown report format", func(c *Config) { c.Reports.Formats = []string{"pdf"} }, "informes.formato"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			tt.modify(&cfg)

			problems := cfg.Validate()

			if len(problems) != 1 || problems[0].Path != tt.wantPath {
				t.Errorf("got problems %+v, want exactly one at %q", problems, tt.wantPath)
			}
		})
	}
}

func TestRF25_AcceptsNumericInitialLevel(t *testing.T) {
	cfg := Default()
	cfg.Difficulty.Initial = "4"

	if problems := cfg.Validate(); len(problems) > 0 {
		t.Errorf("got problems %+v, want none", problems)
	}
}

func TestADR0003_InvalidAPIKeyEnvNeverEchoesTheValue(t *testing.T) {
	const secret = "sk-ant-api03-secret"
	cfg := Default()
	cfg.LLM.APIKeyEnv = secret

	for _, p := range cfg.Validate() {
		if strings.Contains(p.Message, secret) {
			t.Fatalf("problem message leaks the value: %q", p.Message)
		}
	}
}
