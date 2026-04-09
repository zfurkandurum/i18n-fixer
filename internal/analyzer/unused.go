package analyzer

import (
	"strings"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// FindUnusedKeys finds keys in translation files that are never used in code.
// ignorePatterns supports "PREFIX.*" (dot-separated prefix), "PREFIX_*" (raw string prefix),
// "*.SUFFIX" (suffix), or exact key names.
// dynamicPrefixes holds raw string prefixes detected from dynamic key patterns
// (e.g. "SEASON.TIP_" from 'SEASON.TIP_' + variable | translate).
func FindUnusedKeys(usedKeys []types.UsedKey, i18nEntries []types.I18nEntry, ignorePatterns []string, dynamicPrefixes []string) []types.UnusedKeyIssue {
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
		if !usedKeySet[key] && !keyMatchesAnyPattern(key, ignorePatterns) && !keyHasDynamicPrefix(key, dynamicPrefixes) {
			issues = append(issues, types.UnusedKeyIssue{
				Key:       key,
				DefinedIn: locations,
			})
		}
	}

	return issues
}

// keyHasDynamicPrefix returns true if the key starts with any of the dynamic prefixes.
func keyHasDynamicPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if p != "" && strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// dedupPrefixes removes duplicate entries from a slice of prefixes.
func dedupPrefixes(prefixes []string) []string {
	seen := make(map[string]bool, len(prefixes))
	result := prefixes[:0]
	for _, p := range prefixes {
		if !seen[p] {
			seen[p] = true
			result = append(result, p)
		}
	}
	return result
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
	// Raw string prefix: "SEASON.TIP_*" matches any key starting with "SEASON.TIP_"
	if strings.HasSuffix(pattern, "*") && !strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(key, prefix)
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := strings.TrimPrefix(pattern, "*.")
		return key == suffix || strings.HasSuffix(key, "."+suffix)
	}
	return key == pattern
}
