package analyzer

import (
	"testing"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

func TestFindDuplicateKeysConflict(t *testing.T) {
	entries := []types.I18nEntry{
		{Key: "save", Value: "Save", File: "common.json", Locale: "en"},
		{Key: "save", Value: "Save it", File: "buttons.json", Locale: "en"},
		{Key: "cancel", Value: "Cancel", File: "common.json", Locale: "en"},
	}

	issues := FindDuplicateKeys(entries)

	if len(issues) != 1 {
		t.Fatalf("expected 1 duplicate issue, got %d", len(issues))
	}
	if issues[0].Key != "save" {
		t.Errorf("expected key 'save', got %q", issues[0].Key)
	}
	if len(issues[0].Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(issues[0].Values))
	}
}

func TestFindDuplicateKeysSameValue(t *testing.T) {
	// Same key, same value in different files — NOT a conflict
	entries := []types.I18nEntry{
		{Key: "save", Value: "Save", File: "common.json", Locale: "en"},
		{Key: "save", Value: "Save", File: "buttons.json", Locale: "en"},
	}

	issues := FindDuplicateKeys(entries)
	if len(issues) != 0 {
		t.Errorf("expected 0 issues for same values, got %d", len(issues))
	}
}

func TestFindDuplicateKeysNone(t *testing.T) {
	entries := []types.I18nEntry{
		{Key: "save", Value: "Save", File: "en.json", Locale: "en"},
		{Key: "cancel", Value: "Cancel", File: "en.json", Locale: "en"},
	}

	issues := FindDuplicateKeys(entries)
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}
