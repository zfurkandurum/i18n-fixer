package parser

import (
	"fmt"
	"os"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
	"gopkg.in/yaml.v3"
)

// ParseYAML reads a YAML i18n file and returns flattened entries.
func ParseYAML(filePath, keySeparator string) ([]types.I18nEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading YAML file: %w", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %w", filePath, err)
	}

	if raw == nil {
		return nil, nil
	}

	var entries []types.I18nEntry
	flattenMap(raw, "", keySeparator, &entries)
	return entries, nil
}
