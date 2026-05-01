package analyzer

import (
	"testing"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

func TestFindMissingKeys(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "common.save", File: "src/App.tsx", Line: 10},
		{Key: "common.cancel", File: "src/App.tsx", Line: 11},
		{Key: "errors.timeout", File: "src/api.ts", Line: 5},
	}

	i18nEntries := []types.I18nEntry{
		{Key: "common.save", Locale: "en"},
		{Key: "common.save", Locale: "fr"},
		{Key: "common.cancel", Locale: "en"},
		// common.cancel missing from fr
		// errors.timeout missing from both
	}

	locales := []string{"en", "fr"}
	issues := FindMissingKeys(usedKeys, i18nEntries, locales)

	if len(issues) < 2 {
		t.Fatalf("expected at least 2 missing key issues, got %d", len(issues))
	}

	issueMap := make(map[string]types.MissingKeyIssue)
	for _, issue := range issues {
		issueMap[issue.Key] = issue
	}

	if issue, ok := issueMap["common.cancel"]; ok {
		if len(issue.MissingFromLocales) != 1 || issue.MissingFromLocales[0] != "fr" {
			t.Errorf("common.cancel should be missing from fr, got %v", issue.MissingFromLocales)
		}
	} else {
		t.Error("expected common.cancel to be in missing keys")
	}

	if issue, ok := issueMap["errors.timeout"]; ok {
		if len(issue.MissingFromLocales) != 2 {
			t.Errorf("errors.timeout should be missing from 2 locales, got %d", len(issue.MissingFromLocales))
		}
	} else {
		t.Error("expected errors.timeout to be in missing keys")
	}
}

// TestFindMissingKeysPluralRegression covers the case where a translation
// file defines a plural object (e.g. {"a":{"b":{"zero","one","other"}}}).
// After Fix 1, the parser collapses such objects into a single entry under
// the parent key. Code referring to that parent key (`'a.b'.plural(count)`,
// `t('a.b', {count})`, etc.) must NOT be flagged as missing.
func TestFindMissingKeysPluralRegression(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "a.b", File: "src/foo.ts", Line: 10},
	}

	// Parser output after Fix 1: parent collapsed into one entry per locale.
	i18nEntries := []types.I18nEntry{
		{Key: "a.b", Locale: "en", Value: "{} items"},
		{Key: "a.b", Locale: "fr", Value: "{} éléments"},
	}

	issues := FindMissingKeys(usedKeys, i18nEntries, []string{"en", "fr"})
	if len(issues) != 0 {
		t.Errorf("expected 0 missing keys for plural-collapsed parent, got %d: %v", len(issues), issues)
	}
}

func TestFindMissingKeysAllPresent(t *testing.T) {
	usedKeys := []types.UsedKey{
		{Key: "common.save", File: "src/App.tsx", Line: 10},
	}

	i18nEntries := []types.I18nEntry{
		{Key: "common.save", Locale: "en"},
		{Key: "common.save", Locale: "fr"},
	}

	issues := FindMissingKeys(usedKeys, i18nEntries, []string{"en", "fr"})
	if len(issues) != 0 {
		t.Errorf("expected 0 missing keys, got %d", len(issues))
	}
}
