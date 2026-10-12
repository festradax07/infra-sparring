package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Problem is a setting that breaks a rule, identified by its YAML path
// (for example "llm.provider") so the loader can point to the file and line.
type Problem struct {
	Path    string
	Message string
}

const (
	minLevel = 1
	maxLevel = 10
)

// envVarName matches a conventional environment variable name. Requiring it
// also stops players from pasting the key itself into the file.
var envVarName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks every rule and returns all problems found, so the player
// can fix them in one pass. An empty result means the config is valid.
func (c Config) Validate() []Problem {
	v := &validator{}

	v.required("pack", c.Pack)
	v.oneOf("modo", c.Mode, ModeExam, ModeIncident, ModeMixed)

	v.oneOf("llm.provider", c.LLM.Provider, ProviderAnthropic, ProviderOpenAICompatible)
	v.required("llm.model", c.LLM.Model)
	// The value is never echoed: if it is wrong, it may be the key itself.
	v.check(envVarName.MatchString(c.LLM.APIKeyEnv), "llm.api_key_env",
		"tiene que ser el nombre de una variable de entorno (por ejemplo ANTHROPIC_API_KEY), nunca la key")

	v.initialLevel("dificultad.inicial", c.Difficulty.Initial)
	v.oneOf("dificultad.techo", c.Difficulty.Ceiling, CeilingExam, CeilingFree)

	v.positive("examen.minutos_por_desafio", c.Exam.MinutesPerChallenge)
	v.positive("mixto.tareas", c.Mixed.Tasks)

	v.nonEmptyUnique("nodos", c.Nodes)
	v.check(c.Hints >= 0 && c.Hints <= MaxHints, "pistas",
		fmt.Sprintf("valor inválido %d; tiene que estar entre 0 y %d", c.Hints, MaxHints))
	v.oneOf("idioma", c.Language, "es")

	v.required("informes.carpeta", c.Reports.Dir)
	v.check(len(c.Reports.Formats) > 0, "informes.formato", "tiene que incluir al menos un formato")
	for _, format := range c.Reports.Formats {
		v.oneOf("informes.formato", format, ReportMarkdown, ReportJSON)
	}

	return v.problems
}

// validator accumulates problems instead of stopping at the first one.
type validator struct {
	problems []Problem
}

func (v *validator) check(ok bool, path, message string) {
	if !ok {
		v.problems = append(v.problems, Problem{Path: path, Message: message})
	}
}

func (v *validator) required(path, value string) {
	v.check(strings.TrimSpace(value) != "", path, "no puede estar vacío")
}

func (v *validator) oneOf(path, value string, valid ...string) {
	for _, candidate := range valid {
		if value == candidate {
			return
		}
	}
	v.check(false, path, fmt.Sprintf("valor inválido %q; valores válidos: %s", value, strings.Join(valid, ", ")))
}

func (v *validator) positive(path string, value int) {
	v.check(value > 0, path, fmt.Sprintf("valor inválido %d; tiene que ser mayor que 0", value))
}

func (v *validator) initialLevel(path, value string) {
	if value == InitialLevelAuto {
		return
	}
	level, err := strconv.Atoi(value)
	v.check(err == nil && level >= minLevel && level <= maxLevel, path,
		fmt.Sprintf("valor inválido %q; valores válidos: %s o un número de %d a %d",
			value, InitialLevelAuto, minLevel, maxLevel))
}

func (v *validator) nonEmptyUnique(path string, values []string) {
	v.check(len(values) > 0, path, "tiene que tener al menos un elemento")

	seen := make(map[string]bool, len(values))
	for _, value := range values {
		v.check(strings.TrimSpace(value) != "", path, "no puede tener elementos vacíos")
		v.check(!seen[value], path, fmt.Sprintf("%q está repetido", value))
		seen[value] = true
	}
}
