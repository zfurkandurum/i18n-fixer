package analyzer

import (
	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// FindUnusedKeys finds keys in translation files that are never used in code.
func FindUnusedKeys(usedKeys []types.UsedKey, i18nEntries []types.I18nEntry) []types.UnusedKeyIssue {
	usedKeySet := make(map[string]bool)
	for _, uk := range usedKeys {
		usedKeySet[uk.Key] = true
	}

	// Group i18n entries by key
	entryLocations := make(map[string][]types.LocaleLocation)
	for _, entry := range i18nEntries {
		entryLocations[entry.Key] = append(entryLocations[entry.Key], types.LocaleLocation{
			File:   entry.File,
			Locale: entry.Locale,
		})
	}

	var issues []types.UnusedKeyIssue

	for key, locations := range entryLocations {
		if !usedKeySet[key] {
			issues = append(issues, types.UnusedKeyIssue{
				Key:       key,
				DefinedIn: locations,
			})
		}
	}

	return issues
}
