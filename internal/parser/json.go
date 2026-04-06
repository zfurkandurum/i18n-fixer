package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// ParseJSON reads a JSON i18n file and returns flattened entries.
// Supports both flat ("common.save": "Save") and nested ({"common": {"save": "Save"}}) keys.
func ParseJSON(filePath, keySeparator string) ([]types.I18nEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading JSON file: %w", err)
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON file %s: %w", filePath, err)
	}

	var entries []types.I18nEntry
	flattenMap(raw, "", keySeparator, &entries)
	return entries, nil
}

func flattenMap(m map[string]interface{}, prefix, sep string, entries *[]types.I18nEntry) {
	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + sep + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			flattenMap(v, fullKey, sep, entries)
		case string:
			*entries = append(*entries, types.I18nEntry{Key: fullKey, Value: v})
		default:
			// Numbers, booleans, etc. — convert to string
			*entries = append(*entries, types.I18nEntry{Key: fullKey, Value: fmt.Sprintf("%v", v)})
		}
	}
}
