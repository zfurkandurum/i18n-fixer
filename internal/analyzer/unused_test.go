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

	issues := FindUnusedKeys(usedKeys, i18nEntries, nil, nil)

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

	issues := FindUnusedKeys(usedKeys, i18nEntries, nil, nil)
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
	issues := FindUnusedKeys(usedKeys, i18nEntries, []string{"ERRORS.*"}, nil)

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

func TestFindUnusedKeysDynamicPrefixes(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "common.save"},
	}
	i18nEntries := []types.I18nEntry{
		{Key: "common.save", File: "en.json", Locale: "en"},
		{Key: "SEASON.TIP_ALGAE_RISK_HIGH", File: "en.json", Locale: "en"},
		{Key: "SEASON.TIP_LEAF_REMOVAL", File: "en.json", Locale: "en"},
		{Key: "SEASON.TIPS", File: "en.json", Locale: "en"},
		{Key: "common.cancel", File: "en.json", Locale: "en"},
	}

	// Dynamic prefix "SEASON.TIP_" covers SEASON.TIP_* but not SEASON.TIPS
	issues := FindUnusedKeys(usedKeys, i18nEntries, nil, []string{"SEASON.TIP_"})

	unusedKeySet := make(map[string]bool)
	for _, issue := range issues {
		unusedKeySet[issue.Key] = true
	}

	if unusedKeySet["SEASON.TIP_ALGAE_RISK_HIGH"] {
		t.Error("SEASON.TIP_ALGAE_RISK_HIGH should be excluded via dynamic prefix")
	}
	if unusedKeySet["SEASON.TIP_LEAF_REMOVAL"] {
		t.Error("SEASON.TIP_LEAF_REMOVAL should be excluded via dynamic prefix")
	}
	if !unusedKeySet["SEASON.TIPS"] {
		t.Error("SEASON.TIPS should still be reported (prefix SEASON.TIP_ != SEASON.TIPS)")
	}
	if !unusedKeySet["common.cancel"] {
		t.Error("common.cancel should still be reported as unused")
	}
}

// TestFindUnusedKeysDynamicPrefixRegression covers the integration between
// the scanner's inferred-prefix output (Fix 2) and the unused-key analyzer.
// When code interpolates a key like `'a.b.${var}'.tr()`, the scanner emits
// an inferred prefix `a.b.`. Sub-keys defined under that prefix in the
// translation file must NOT be flagged as unused even though no static
// reference points to them.
func TestFindUnusedKeysDynamicPrefixRegression(t *testing.T) {
	// No static UsedKey references — runtime resolves the suffix.
	usedKeys := []types.UsedKey{}

	i18nEntries := []types.I18nEntry{
		{Key: "a.b.x", File: "en.json", Locale: "en"},
		{Key: "a.b.y", File: "en.json", Locale: "en"},
		{Key: "a.c", File: "en.json", Locale: "en"}, // outside prefix → unused
	}

	dynamicPrefixes := []string{"a.b."}
	issues := FindUnusedKeys(usedKeys, i18nEntries, nil, dynamicPrefixes)

	unused := make(map[string]bool)
	for _, issue := range issues {
		unused[issue.Key] = true
	}

	if unused["a.b.x"] {
		t.Error("a.b.x should be excluded via inferred prefix a.b.")
	}
	if unused["a.b.y"] {
		t.Error("a.b.y should be excluded via inferred prefix a.b.")
	}
	if !unused["a.c"] {
		t.Error("a.c should still be reported (outside the inferred prefix)")
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
		// Raw string prefix (trailing underscore)
		{"SEASON.TIP_ALGAE", "SEASON.TIP_*", true},
		{"SEASON.TIPS", "SEASON.TIP_*", false},
		{"ONBOARDING.STEP1", "ONBOARDING.*", true},
	}

	for _, tt := range tests {
		got := keyMatchesPattern(tt.key, tt.pattern)
		if got != tt.want {
			t.Errorf("keyMatchesPattern(%q, %q) = %v, want %v", tt.key, tt.pattern, got, tt.want)
		}
	}
}
