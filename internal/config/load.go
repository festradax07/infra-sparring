package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Paths lists the config files to read, from lowest to highest priority.
// A missing file is not an error: its settings keep the previous value.
type Paths struct {
	Global string
	Local  string
}

// DefaultPaths returns $XDG_CONFIG_HOME/spar/config.yaml (or
// ~/.config/spar/config.yaml) and spar.yaml in the working directory.
func DefaultPaths() (Paths, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return Paths{}, fmt.Errorf("find home directory: %w", err)
		}
		configDir = filepath.Join(home, ".config")
	}

	return Paths{
		Global: filepath.Join(configDir, "spar", "config.yaml"),
		Local:  "spar.yaml",
	}, nil
}

// Load starts from Default, applies each existing file in priority order and
// validates the result. Errors point to the file and line of the bad setting.
func Load(paths Paths) (Config, error) {
	cfg := Default()
	origins := make(origins)

	for _, path := range []string{paths.Global, paths.Local} {
		data, err := os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return Config{}, fmt.Errorf("leer %s: %w", path, err)
		}
		if err := apply(path, data, &cfg, origins); err != nil {
			return Config{}, err
		}
	}

	if problems := cfg.Validate(); len(problems) > 0 {
		return Config{}, origins.explain(problems)
	}
	return cfg, nil
}

// apply overlays the settings present in data onto cfg. Settings absent from
// data keep their current value; lists present in data replace the old ones.
func apply(file string, data []byte, cfg *Config, origins origins) error {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("%s: YAML mal formado: %w", file, err)
	}
	if len(doc.Content) == 0 {
		return nil // empty file
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return fmt.Errorf("%s:%d: se esperaban claves y valores (clave: valor)", file, root.Line)
	}

	if unknown := recordKeys(root, reflect.TypeOf(*cfg), "", file, origins); len(unknown) > 0 {
		return fmt.Errorf("configuración inválida:\n%w", errors.Join(unknown...))
	}

	if err := root.Decode(cfg); err != nil {
		return fmt.Errorf("%s: tipo de valor inválido: %w", file, err)
	}
	return nil
}

// recordKeys walks a YAML mapping alongside the struct type it decodes into.
// It remembers where each known key was written and reports unknown keys,
// which are usually typos that would otherwise be silently ignored.
func recordKeys(node *yaml.Node, structType reflect.Type, prefix, file string, origins origins) []error {
	fields := yamlFields(structType)
	var unknown []error

	for i := 0; i+1 < len(node.Content); i += 2 {
		key, value := node.Content[i], node.Content[i+1]
		path := joinPath(prefix, key.Value)

		fieldType, ok := fields[key.Value]
		if !ok {
			unknown = append(unknown, fmt.Errorf("%s:%d: clave desconocida %q; claves válidas acá: %s",
				file, key.Line, path, strings.Join(sortedKeys(fields), ", ")))
			continue
		}

		origins[path] = position{file: file, line: key.Line}
		if fieldType.Kind() == reflect.Struct && value.Kind == yaml.MappingNode {
			unknown = append(unknown, recordKeys(value, fieldType, path, file, origins)...)
		}
	}
	return unknown
}

// yamlFields maps each YAML key of a struct to the type of its field.
func yamlFields(structType reflect.Type) map[string]reflect.Type {
	fields := make(map[string]reflect.Type, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		name, _, _ := strings.Cut(field.Tag.Get("yaml"), ",")
		if name != "" && name != "-" {
			fields[name] = field.Type
		}
	}
	return fields
}

// position is where a setting was last written.
type position struct {
	file string
	line int
}

// origins maps a YAML path to the position of the file that set it last.
type origins map[string]position

// explain turns validation problems into one error that names, for each
// problem, the file and line to fix, in file order. Settings never written in
// a file come from the defaults and are reported first, without a position.
func (o origins) explain(problems []Problem) error {
	sorted := append([]Problem(nil), problems...)
	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := o[sorted[i].Path], o[sorted[j].Path]
		if a.file != b.file {
			return a.file < b.file
		}
		return a.line < b.line
	})

	lines := make([]error, 0, len(sorted))
	for _, p := range sorted {
		if pos, ok := o[p.Path]; ok {
			lines = append(lines, fmt.Errorf("%s:%d: %s: %s", pos.file, pos.line, p.Path, p.Message))
		} else {
			lines = append(lines, fmt.Errorf("%s: %s", p.Path, p.Message))
		}
	}
	return fmt.Errorf("configuración inválida:\n%w", errors.Join(lines...))
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func sortedKeys(m map[string]reflect.Type) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
