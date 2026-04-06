package analyzer

import (
	"regexp"
	"strings"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

var conventions = map[string]*regexp.Regexp{
	"UPPER_SNAKE": regexp.MustCompile(`^[A-Z][A-Z0-9_]*(\.[A-Z][A-Z0-9_]*)*$`),
	"lower.dot":   regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`),
	"camelCase":   regexp.MustCompile(`^[a-z][a-zA-Z0-9]*(\.[a-z][a-zA-Z0-9]*)*$`),
	"kebab-case":  regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z][a-z0-9]*)*(\.([a-z][a-z0-9]*(-[a-z][a-z0-9]*)*))*$`),
}

// LintKeyNaming checks i18n keys against a naming convention.
// If convention is empty, it auto-detects from the first keys.
func LintKeyNaming(entries []types.I18nEntry, convention string) []types.KeyNamingIssue {
	if convention == "" {
		convention = detectConvention(entries)
	}
	if convention == "" {
		return nil
	}

	re, ok := conventions[convention]
	if !ok {
		return nil
	}

	// Deduplicate: check each unique key only once
	seen := make(map[string]bool)
	var issues []types.KeyNamingIssue

	for _, e := range entries {
		if seen[e.Key] {
			continue
		}
		seen[e.Key] = true

		// Check each segment of the key separately for dot-separated keys
		if !re.MatchString(e.Key) {
			issues = append(issues, types.KeyNamingIssue{
				Key:      e.Key,
				Expected: convention,
				File:     e.File,
				Locale:   e.Locale,
			})
		}
	}

	return issues
}

// detectConvention guesses the naming convention from the first keys.
func detectConvention(entries []types.I18nEntry) string {
	if len(entries) == 0 {
		return ""
	}

	counts := make(map[string]int)
	checked := 0

	seen := make(map[string]bool)
	for _, e := range entries {
		if seen[e.Key] {
			continue
		}
		seen[e.Key] = true

		for name, re := range conventions {
			if re.MatchString(e.Key) {
				counts[name]++
			}
		}
		checked++
		if checked >= 20 {
			break
		}
	}

	// Pick convention with most matches
	best := ""
	bestCount := 0
	for name, count := range counts {
		if count > bestCount {
			best = name
			bestCount = count
		}
	}

	// Only return if at least 60% of keys match
	if bestCount > 0 && float64(bestCount)/float64(checked) >= 0.6 {
		return best
	}

	return ""
}

// ConventionNames returns the list of supported convention names.
func ConventionNames() []string {
	return []string{"UPPER_SNAKE", "lower.dot", "camelCase", "kebab-case"}
}

// IsValidConvention checks if a convention name is supported.
func IsValidConvention(name string) bool {
	_, ok := conventions[strings.TrimSpace(name)]
	return ok
}
