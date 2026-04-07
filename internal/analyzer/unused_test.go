package analyzer

import (
	"testing"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

func TestFindUnusedKeys(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "common.save"},
	}

	i18nEntries := []types.I18nEntry{
		{Key: "common.save", File: "en.json", Locale: "en"},
		{Key: "common.cancel", File: "en.json", Locale: "en"},
		{Key: "old.feature", File: "en.json", Locale: "en"},
	}

	issues := FindUnusedKeys(usedKeys, i18nEntries, nil)

	if len(issues) != 2 {
		t.Fatalf("expected 2 unused keys, got %d", len(issues))
	}

	unusedKeySet := make(map[string]bool)
	for _, issue := range issues {
		unusedKeySet[issue.Key] = true
	}

	if !unusedKeySet["common.cancel"] {
		t.Error("expected common.cancel to be unused")
	}
	if !unusedKeySet["old.feature"] {
		t.Error("expected old.feature to be unused")
	}
}

func TestFindUnusedKeysNone(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "common.save"},
	}

	i18nEntries := []types.I18nEntry{
		{Key: "common.save", File: "en.json", Locale: "en"},
	}

	issues := FindUnusedKeys(usedKeys, i18nEntries, nil)
	if len(issues) != 0 {
		t.Errorf("expected 0 unused keys, got %d", len(issues))
	}
}

func TestFindUnusedKeysIgnorePatterns(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "common.save"},
	}

	i18nEntries := []types.I18nEntry{
		{Key: "common.save", File: "en.json", Locale: "en"},
		{Key: "common.cancel", File: "en.json", Locale: "en"},
		{Key: "ERRORS.AUTH.invalid", File: "en.json", Locale: "en"},
		{Key: "ERRORS.USER.not_found", File: "en.json", Locale: "en"},
		{Key: "old.feature", File: "en.json", Locale: "en"},
	}

	// Ignore ERRORS.* — these are used dynamically via API responses
	issues := FindUnusedKeys(usedKeys, i18nEntries, []string{"ERRORS.*"})

	unusedKeySet := make(map[string]bool)
	for _, issue := range issues {
		unusedKeySet[issue.Key] = true
	}

	if unusedKeySet["ERRORS.AUTH.invalid"] {
		t.Error("ERRORS.AUTH.invalid should be ignored")
	}
	if unusedKeySet["ERRORS.USER.not_found"] {
		t.Error("ERRORS.USER.not_found should be ignored")
	}
	if !unusedKeySet["common.cancel"] {
		t.Error("common.cancel should still be reported as unused")
	}
}

func TestKeyMatchesPattern(t *testing.T) {
	tests := []struct {
		key     string
		pattern string
		want    bool
	}{
		{"ERRORS.AUTH.invalid", "ERRORS.*", true},
		{"ERRORS", "ERRORS.*", true},
		{"OTHER.KEY", "ERRORS.*", false},
		{"common.cancel", "*.cancel", true},
		{"common.save", "*.cancel", false},
		{"exact.key", "exact.key", true},
		{"exact.key.extra", "exact.key", false},
		{"anything", "*", true},
	}

	for _, tt := range tests {
		got := keyMatchesPattern(tt.key, tt.pattern)
		if got != tt.want {
			t.Errorf("keyMatchesPattern(%q, %q) = %v, want %v", tt.key, tt.pattern, got, tt.want)
		}
	}
}
