package analyzer

import (
	"strings"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// FindUnusedKeys finds keys in translation files that are never used in code.
// ignorePatterns supports "PREFIX.*" (prefix), "*.SUFFIX" (suffix), or exact key names.
func FindUnusedKeys(usedKeys []types.UsedKey, i18nEntries []types.I18nEntry, ignorePatterns []string) []types.UnusedKeyIssue {
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
		if !usedKeySet[key] && !keyMatchesAnyPattern(key, ignorePatterns) {
			issues = append(issues, types.UnusedKeyIssue{
				Key:       key,
				DefinedIn: locations,
			})
		}
	}

	return issues
}

// keyMatchesAnyPattern returns true if key matches any of the ignore patterns.
// Supported forms:
//   - "PREFIX.*"  → key == PREFIX or key starts with PREFIX.
//   - "*.SUFFIX"  → key == SUFFIX or key ends with .SUFFIX
//   - "exact.key" → exact match
func keyMatchesAnyPattern(key string, patterns []string) bool {
	for _, p := range patterns {
		if keyMatchesPattern(key, p) {
			return true
		}
	}
	return false
}

func keyMatchesPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return key == prefix || strings.HasPrefix(key, prefix+".")
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := strings.TrimPrefix(pattern, "*.")
		return key == suffix || strings.HasSuffix(key, "."+suffix)
	}
	return key == pattern
}
