package analyzer

import (
	"github.com/i18n-fixer/i18n-fixer/internal/scanner"
	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// Analyze runs all analyzers and produces the final audit result.
func Analyze(scanResult *scanner.ScanResult, i18nEntries []types.I18nEntry, keySeparator string, opts Options) *types.AuditResult {
	locales := extractLocales(i18nEntries)

	var missingKeys []types.MissingKeyIssue
	var unusedKeys []types.UnusedKeyIssue
	var hardcodedStrings []types.HardcodedStringIssue

	if !opts.NoMissing {
		missingKeys = FindMissingKeys(scanResult.UsedKeys, i18nEntries, locales)
	}

	if !opts.NoUnused {
		unusedKeys = FindUnusedKeys(scanResult.UsedKeys, i18nEntries)
	}

	if !opts.NoHardcoded {
		hardcodedStrings = GroupHardcodedStrings(scanResult.Hardcoded, keySeparator)
	}

	return &types.AuditResult{
		Summary: types.AuditSummary{
			Locales:              locales,
			TotalKeys:            countUniqueKeys(i18nEntries),
			MissingKeyCount:      len(missingKeys),
			UnusedKeyCount:       len(unusedKeys),
			HardcodedStringCount: len(hardcodedStrings),
			DynamicKeyCount:      len(scanResult.DynamicKeys),
		},
		MissingKeys:      missingKeys,
		UnusedKeys:       unusedKeys,
		HardcodedStrings: hardcodedStrings,
		DynamicKeys:      scanResult.DynamicKeys,
	}
}

// Options controls which analyzers to run.
type Options struct {
	NoMissing   bool
	NoUnused    bool
	NoHardcoded bool
}

func extractLocales(entries []types.I18nEntry) []string {
	seen := make(map[string]bool)
	var locales []string
	for _, e := range entries {
		if !seen[e.Locale] {
			seen[e.Locale] = true
			locales = append(locales, e.Locale)
		}
	}
	return locales
}

func countUniqueKeys(entries []types.I18nEntry) int {
	seen := make(map[string]bool)
	for _, e := range entries {
		seen[e.Key] = true
	}
	return len(seen)
}
