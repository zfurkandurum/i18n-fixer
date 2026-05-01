package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
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
			if isPluralObject(v) {
				// ICU plural form — treat the parent as a single key. The
				// runtime resolves a sub-category at call time; static
				// analysis only needs to know the parent exists so callers
				// like `t('a.b', count)` or `'a.b'.plural(count)` aren't
				// reported as missing.
				*entries = append(*entries, types.I18nEntry{
					Key:   fullKey,
					Value: v["other"].(string),
				})
			} else {
				flattenMap(v, fullKey, sep, entries)
			}
		case string:
			*entries = append(*entries, types.I18nEntry{Key: fullKey, Value: v})
		default:
			// Numbers, booleans, etc. — convert to string
			*entries = append(*entries, types.I18nEntry{Key: fullKey, Value: fmt.Sprintf("%v", v)})
		}
	}
}

// pluralCategories enumerates ICU MessageFormat plural keywords. ICU mandates
// that any plural form include `other`; the rest are optional per locale.
var pluralCategories = map[string]struct{}{
	"zero":  {},
	"one":   {},
	"two":   {},
	"few":   {},
	"many":  {},
	"other": {},
}

// isPluralObject reports whether m represents an ICU plural form (typically
// the payload of a framework's plural-aware translate function: parent key
// resolved at call site, sub-keys selected at runtime per the count value).
//
// All leaf keys must be plural categories AND `other` must be present
// (ICU mandates it). All values must be strings — a nested object under
// any plural key is ambiguous, so we fall back to regular flattening.
func isPluralObject(m map[string]interface{}) bool {
	if len(m) == 0 {
		return false
	}
	hasOther := false
	for k, v := range m {
		if _, ok := pluralCategories[k]; !ok {
			return false
		}
		if _, ok := v.(string); !ok {
			return false
		}
		if k == "other" {
			hasOther = true
		}
	}
	return hasOther
}
