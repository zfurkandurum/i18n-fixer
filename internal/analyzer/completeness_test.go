package analyzer

import (
	"math"
	"testing"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

func TestAnalyzeCompleteness(t *testing.T) {
	entries := []types.I18nEntry{
		{Key: "save", Locale: "en"},
		{Key: "cancel", Locale: "en"},
		{Key: "delete", Locale: "en"},
		{Key: "save", Locale: "fr"},
		{Key: "cancel", Locale: "fr"},
		// fr missing "delete"
	}

	coverage := AnalyzeCompleteness(entries, []string{"en", "fr"})

	if len(coverage) != 2 {
		t.Fatalf("expected 2 locales, got %d", len(coverage))
	}

	for _, c := range coverage {
		if c.Locale == "en" {
			if c.HasKeys != 3 || c.TotalKeys != 3 {
				t.Errorf("en: expected 3/3, got %d/%d", c.HasKeys, c.TotalKeys)
			}
			if math.Abs(c.Percentage-100.0) > 0.1 {
				t.Errorf("en: expected 100%%, got %.1f%%", c.Percentage)
			}
		}
		if c.Locale == "fr" {
			if c.HasKeys != 2 || c.TotalKeys != 3 {
				t.Errorf("fr: expected 2/3, got %d/%d", c.HasKeys, c.TotalKeys)
			}
			if math.Abs(c.Percentage-66.7) > 0.1 {
				t.Errorf("fr: expected ~66.7%%, got %.1f%%", c.Percentage)
			}
		}
	}
}

func TestAnalyzeCompletenessEmpty(t *testing.T) {
	coverage := AnalyzeCompleteness(nil, []string{"en"})
	if coverage != nil {
		t.Errorf("expected nil for empty entries, got %v", coverage)
	}
}

func TestOverallCompleteness(t *testing.T) {
	coverage := []types.LocaleCoverage{
		{Locale: "en", Percentage: 100},
		{Locale: "fr", Percentage: 50},
	}
	overall := OverallCompleteness(coverage)
	if math.Abs(overall-75.0) > 0.1 {
		t.Errorf("expected 75%%, got %.1f%%", overall)
	}
}
