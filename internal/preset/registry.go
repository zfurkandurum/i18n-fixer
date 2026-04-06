package preset

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

//go:embed builtin/*.json
var builtinFS embed.FS

var builtinPresets map[string]types.FrameworkPreset

func init() {
	builtinPresets = make(map[string]types.FrameworkPreset)

	entries, err := builtinFS.ReadDir("builtin")
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := builtinFS.ReadFile(filepath.Join("builtin", entry.Name()))
		if err != nil {
			continue
		}

		var p types.FrameworkPreset
		if err := json.Unmarshal(data, &p); err != nil {
			continue
		}

		builtinPresets[p.Name] = p
	}
}

// Get returns a built-in preset by name.
func Get(name string) (types.FrameworkPreset, error) {
	p, ok := builtinPresets[name]
	if !ok {
		return types.FrameworkPreset{}, fmt.Errorf("unknown preset: %s", name)
	}
	return p, nil
}

// All returns all built-in presets.
func All() map[string]types.FrameworkPreset {
	return builtinPresets
}

// LoadCustom reads a preset from a JSON file path.
func LoadCustom(path string) (types.FrameworkPreset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.FrameworkPreset{}, fmt.Errorf("reading preset file: %w", err)
	}

	var p types.FrameworkPreset
	if err := json.Unmarshal(data, &p); err != nil {
		return types.FrameworkPreset{}, fmt.Errorf("parsing preset file: %w", err)
	}

	if p.Name == "" {
		return types.FrameworkPreset{}, fmt.Errorf("preset must have a name")
	}

	return p, nil
}

// Names returns sorted list of all built-in preset names.
func Names() []string {
	names := make([]string, 0, len(builtinPresets))
	for name := range builtinPresets {
		names = append(names, name)
	}
	return names
}
