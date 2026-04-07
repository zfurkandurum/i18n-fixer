package analyzer

import (
	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// AnalyzeCompleteness calculates translation coverage percentage per locale.
func AnalyzeCompleteness(i18nEntries []types.I18nEntry, locales []string) []types.LocaleCoverage {
	// Collect all unique keys across all locales
	allKeys := make(map[string]bool)
	for _, e := range i18nEntries {
		allKeys[e.Key] = true
	}
	totalKeys := len(allKeys)

	if totalKeys == 0 {
		return nil
	}

	// Count keys per locale
	keysPerLocale := make(map[string]map[string]bool)
	for _, locale := range locales {
		keysPerLocale[locale] = make(map[string]bool)
	}
	for _, e := range i18nEntries {
		if m, ok := keysPerLocale[e.Locale]; ok {
			m[e.Key] = true
		}
	}

	var coverage []types.LocaleCoverage
	for _, locale := range locales {
		hasKeys := len(keysPerLocale[locale])
		pct := float64(hasKeys) / float64(totalKeys) * 100
		coverage = append(coverage, types.LocaleCoverage{
			Locale:     locale,
			TotalKeys:  totalKeys,
			HasKeys:    hasKeys,
			Percentage: pct,
		})
	}

	return coverage
}

// OverallCompleteness returns the average completeness across all locales.
func OverallCompleteness(coverage []types.LocaleCoverage) float64 {
	if len(coverage) == 0 {
		return 0
	}
	var sum float64
	for _, c := range coverage {
		sum += c.Percentage
	}
	return sum / float64(len(coverage))
}
