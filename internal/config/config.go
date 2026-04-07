package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the .i18n-fixer.json configuration file.
type Config struct {
	Preset                  string   `json:"preset,omitempty"`
	Ignore                  []string `json:"ignore,omitempty"`
	DefaultLocale           string   `json:"defaultLocale,omitempty"`
	Format                  string   `json:"format,omitempty"`
	Output                  string   `json:"output,omitempty"`
	StrictUnused            bool     `json:"strictUnused,omitempty"`
	NoHardcoded             bool     `json:"noHardcoded,omitempty"`
	NoMissing               bool     `json:"noMissing,omitempty"`
	NoUnused                bool     `json:"noUnused,omitempty"`
	Verbose                 bool     `json:"verbose,omitempty"`
	// UnusedKeyIgnorePatterns lists key glob patterns to exclude from unused-key reporting.
	// Supports "PREFIX.*" (prefix match), "*.SUFFIX" (suffix match), or exact key names.
	// Example: ["ERRORS.*", "DEPRECATED.*"]
	UnusedKeyIgnorePatterns []string `json:"unusedKeyIgnorePatterns,omitempty"`
	// I18nFunctionPatterns adds extra key-detection patterns on top of the preset's built-in ones.
	I18nFunctionPatterns    []string `json:"i18nFunctionPatterns,omitempty"`
}

const configFileName = ".i18n-fixer.json"

// Load searches for .i18n-fixer.json starting from dir and walking up.
// Returns a zero Config if no file is found (not an error).
func Load(dir string) (Config, error) {
	path := findConfigFile(dir)
	if path == "" {
		return Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func findConfigFile(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}

	for {
		candidate := filepath.Join(abs, configFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(abs)
		if parent == abs {
			break
		}
		abs = parent
	}

	return ""
}

// Defaults returns the default configuration.
func Defaults() Config {
	return Config{
		Format: "console",
	}
}
