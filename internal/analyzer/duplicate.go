package analyzer

import (
	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// FindDuplicateKeys detects keys that appear multiple times in the same locale
// with different values (conflicting translations).
func FindDuplicateKeys(i18nEntries []types.I18nEntry) []types.DuplicateKeyIssue {
	// Group entries by locale+key
	type localeKey struct {
		locale string
		key    string
	}

	seen := make(map[localeKey][]types.DuplicateValue)

	for _, e := range i18nEntries {
		lk := localeKey{locale: e.Locale, key: e.Key}
		seen[lk] = append(seen[lk], types.DuplicateValue{
			Value: e.Value,
			File:  e.File,
		})
	}

	var issues []types.DuplicateKeyIssue

	for lk, values := range seen {
		if len(values) < 2 {
			continue
		}

		// Check if values actually differ
		hasConflict := false
		first := values[0].Value
		for _, v := range values[1:] {
			if v.Value != first {
				hasConflict = true
				break
			}
		}

		if hasConflict {
			issues = append(issues, types.DuplicateKeyIssue{
				Key:    lk.key,
				Locale: lk.locale,
				Values: values,
			})
		}
	}

	return issues
}
