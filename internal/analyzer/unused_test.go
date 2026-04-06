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

	issues := FindUnusedKeys(usedKeys, i18nEntries)

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

	issues := FindUnusedKeys(usedKeys, i18nEntries)
	if len(issues) != 0 {
		t.Errorf("expected 0 unused keys, got %d", len(issues))
	}
}
