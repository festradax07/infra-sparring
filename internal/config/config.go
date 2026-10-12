// Package config loads and validates the player's preferences (RF-24, RF-25).
//
// Preferences come from, lowest to highest priority: built-in defaults, the
// global file ~/.config/spar/config.yaml and spar.yaml in the working
// directory. Pack rules (allowlist, protected resources) never live here.
package config

// Valid values for the enumerated settings.
const (
	ModeExam     = "examen"
	ModeIncident = "incidente"
	ModeMixed    = "mixto"

	ProviderAnthropic        = "anthropic"
	ProviderOpenAICompatible = "openai-compatible"

	InitialLevelAuto = "auto"

	CeilingExam = "examen"
	CeilingFree = "libre"

	ReportMarkdown = "md"
	ReportJSON     = "json"
)

// MaxHints is the most hints a challenge can offer (RF-16).
const MaxHints = 2

// Config holds the player's preferences. Identifiers are in English; YAML keys
// are in Spanish because players edit the file by hand.
type Config struct {
	Pack       string     `yaml:"pack"`
	Mode       string     `yaml:"modo"`
	LLM        LLM        `yaml:"llm"`
	Difficulty Difficulty `yaml:"dificultad"`
	Exam       Exam       `yaml:"examen"`
	Incident   Incident   `yaml:"incidente"`
	Mixed      Mixed      `yaml:"mixto"`
	Topics     Topics     `yaml:"temas"`
	Nodes      []string   `yaml:"nodos"`
	Hints      int        `yaml:"pistas"`
	Language   string     `yaml:"idioma"`
	Reports    Reports    `yaml:"informes"`
}

// LLM selects the model provider. It names the environment variable that
// holds the API key, never the key itself (ADR 0003).
type LLM struct {
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env"`
}

// Difficulty sets the starting level and its ceiling.
type Difficulty struct {
	// Initial is "auto" or a level from 1 to 10, kept as text to accept both.
	Initial string `yaml:"inicial"`
	Ceiling string `yaml:"techo"`
}

// Exam configures exam mode (RF-01).
type Exam struct {
	MinutesPerChallenge int `yaml:"minutos_por_desafio"`
}

// Incident configures incident mode (RF-02, phase 2).
type Incident struct {
	Clock         bool `yaml:"reloj"`
	ChainedFaults bool `yaml:"fallas_encadenadas"`
}

// Mixed configures mixed mode (RF-03, phase 2).
type Mixed struct {
	Tasks   int      `yaml:"tareas"`
	Domains []string `yaml:"dominios"`
}

// Topics narrows the competencies a challenge can come from.
type Topics struct {
	Include []string `yaml:"incluir"`
	Exclude []string `yaml:"excluir"`
}

// Reports sets where and in which formats reports are written (RF-22).
type Reports struct {
	Dir     string   `yaml:"carpeta"`
	Formats []string `yaml:"formato"`
}

// Default returns the configuration used when no file overrides a setting.
func Default() Config {
	return Config{
		Pack: "lfcs",
		Mode: ModeExam,
		LLM: LLM{
			Provider:  ProviderAnthropic,
			Model:     "claude-sonnet-5-5",
			APIKeyEnv: "ANTHROPIC_API_KEY",
		},
		Difficulty: Difficulty{Initial: InitialLevelAuto, Ceiling: CeilingExam},
		Exam:       Exam{MinutesPerChallenge: 7},
		Incident:   Incident{Clock: false, ChainedFaults: true},
		Mixed:      Mixed{Tasks: 3, Domains: []string{"storage", "networking", "servicios"}},
		Nodes:      []string{"ubuntu", "alma"},
		Hints:      MaxHints,
		Language:   "es",
		Reports:    Reports{Dir: "reports", Formats: []string{ReportMarkdown, ReportJSON}},
	}
}
