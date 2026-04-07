package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// ParseARB reads a Flutter .arb file (JSON with @-prefixed metadata) and returns entries.
func ParseARB(filePath string) ([]types.I18nEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading ARB file: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing ARB file %s: %w", filePath, err)
	}

	var entries []types.I18nEntry

	for key, value := range raw {
		// Skip @-prefixed metadata keys (@@locale, @keyDescription, etc.)
		if strings.HasPrefix(key, "@") {
			continue
		}

		strVal, ok := value.(string)
		if !ok {
			strVal = fmt.Sprintf("%v", value)
		}

		entries = append(entries, types.I18nEntry{Key: key, Value: strVal})
	}

	return entries, nil
}
