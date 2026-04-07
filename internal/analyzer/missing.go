package analyzer

import (
	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// FindMissingKeys finds keys used in code but absent from translation files.
func FindMissingKeys(usedKeys []types.UsedKey, i18nEntries []types.I18nEntry, locales []string) []types.MissingKeyIssue {
	// Build key sets per locale
	keysByLocale := make(map[string]map[string]bool)
	for _, locale := range locales {
		keysByLocale[locale] = make(map[string]bool)
	}
	for _, entry := range i18nEntries {
		if _, ok := keysByLocale[entry.Locale]; ok {
			keysByLocale[entry.Locale][entry.Key] = true
		}
	}

	// Group used keys by key name
	usedKeyLocations := make(map[string][]types.Location)
	for _, uk := range usedKeys {
		usedKeyLocations[uk.Key] = append(usedKeyLocations[uk.Key], types.Location{
			File:   uk.File,
			Line:   uk.Line,
			Column: uk.Column,
		})
	}

	var issues []types.MissingKeyIssue

	for key, locations := range usedKeyLocations {
		var missingFrom []string
		for _, locale := range locales {
			if !keysByLocale[locale][key] {
				missingFrom = append(missingFrom, locale)
			}
		}

		if len(missingFrom) > 0 {
			issues = append(issues, types.MissingKeyIssue{
				Key:                key,
				UsedIn:             locations,
				MissingFromLocales: missingFrom,
			})
		}
	}

	return issues
}
