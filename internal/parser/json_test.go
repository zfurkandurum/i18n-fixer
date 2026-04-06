package parser

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func TestParseJSONFlat(t *testing.T) {
	entries, err := ParseJSON(filepath.Join(testdataDir(), "flat.json"), ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	keys := make(map[string]string)
	for _, e := range entries {
		keys[e.Key] = e.Value
	}

	if keys["save"] != "Save" {
		t.Errorf("expected save=Save, got %q", keys["save"])
	}
	if keys["cancel"] != "Cancel" {
		t.Errorf("expected cancel=Cancel, got %q", keys["cancel"])
	}
}

func TestParseJSONNested(t *testing.T) {
	entries, err := ParseJSON(filepath.Join(testdataDir(), "nested.json"), ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys := make(map[string]string)
	for _, e := range entries {
		keys[e.Key] = e.Value
	}

	if keys["common.save"] != "Save" {
		t.Errorf("expected common.save=Save, got %q", keys["common.save"])
	}
	if keys["errors.network.timeout"] != "Connection timed out" {
		t.Errorf("expected errors.network.timeout=Connection timed out, got %q", keys["errors.network.timeout"])
	}
}

func TestParseJSONEmpty(t *testing.T) {
	entries, err := ParseJSON(filepath.Join(testdataDir(), "empty.json"), ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseJSONNonexistent(t *testing.T) {
	_, err := ParseJSON(filepath.Join(testdataDir(), "nonexistent.json"), ".")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestDetectLocale(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"src/locales/en/common.json", "en"},
		{"src/locales/fr/common.json", "fr"},
		{"res/values-fr/strings.xml", "fr"},
		{"res/values/strings.xml", "default"},
		{"lib/l10n/app_en.arb", "en"},
		{"en.json", "en"},
		{"zh-CN.json", "zh-CN"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := detectLocale(tt.path)
			if got != tt.expected {
				t.Errorf("detectLocale(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func sortedKeys(entries []string) []string {
	sort.Strings(entries)
	return entries
}
