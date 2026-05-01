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

func TestParseJSONPlural(t *testing.T) {
	entries, err := ParseJSON(filepath.Join(testdataDir(), "plural.json"), ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys := make(map[string]string)
	for _, e := range entries {
		keys[e.Key] = e.Value
	}

	// items.count is an ICU plural object with all required categories.
	// It should be collapsed into a single entry whose value is the `other`
	// form, and no sub-keys should be emitted.
	if v, ok := keys["items.count"]; !ok {
		t.Errorf("expected items.count to be present (collapsed plural), got map=%v", keys)
	} else if v != "{} items" {
		t.Errorf("expected items.count value=\"{} items\" (other form), got %q", v)
	}
	for _, sub := range []string{"items.count.zero", "items.count.one", "items.count.other"} {
		if _, ok := keys[sub]; ok {
			t.Errorf("did not expect sub-key %q to be emitted (parent is plural)", sub)
		}
	}

	// items.subtitle sits next to a plural sibling but is itself a flat string;
	// it must still be flattened normally.
	if keys["items.subtitle"] != "List of items" {
		t.Errorf("expected items.subtitle=\"List of items\", got %q", keys["items.subtitle"])
	}

	// messages.unread has one+other (no zero/two/few/many) — still a valid
	// plural per ICU (only `other` is mandatory).
	if v, ok := keys["messages.unread"]; !ok || v != "{} unread messages" {
		t.Errorf("expected messages.unread to be collapsed (other form), got %q (present=%v)", v, ok)
	}

	// size has `other` but also `small` (non-category) — must NOT be treated
	// as plural. Both leaf keys flatten as before.
	if _, ok := keys["size"]; ok {
		t.Errorf("did not expect size to be collapsed (mixed keys)")
	}
	if keys["size.small"] != "Small" || keys["size.other"] != "Other size" {
		t.Errorf("expected size.small and size.other to flatten normally, got %v", keys)
	}
}

func TestIsPluralObject(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  bool
	}{
		{
			name:  "empty",
			input: map[string]interface{}{},
			want:  false,
		},
		{
			name:  "missing other",
			input: map[string]interface{}{"one": "a", "two": "b"},
			want:  false,
		},
		{
			name:  "all categories",
			input: map[string]interface{}{"zero": "z", "one": "o", "two": "t", "few": "f", "many": "m", "other": "x"},
			want:  true,
		},
		{
			name:  "one+other",
			input: map[string]interface{}{"one": "a", "other": "b"},
			want:  true,
		},
		{
			name:  "non-category mixed",
			input: map[string]interface{}{"one": "a", "other": "b", "foo": "c"},
			want:  false,
		},
		{
			name:  "non-string value",
			input: map[string]interface{}{"other": map[string]interface{}{"nested": "x"}},
			want:  false,
		},
		{
			name:  "only other",
			input: map[string]interface{}{"other": "x"},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPluralObject(tt.input); got != tt.want {
				t.Errorf("isPluralObject(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
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
		// Filename-as-locale (namespace pattern): auth/en.json → "en"
		{"src/assets/i18n/auth/en.json", "en"},
		{"src/assets/i18n/auth/fr.json", "fr"},
		{"src/assets/i18n/common/de.json", "de"},
		{"src/assets/i18n/layout/es.json", "es"},
		// Directory-as-locale: locales/en/common.json → "en"
		{"src/locales/en/common.json", "en"},
		{"src/locales/fr/common.json", "fr"},
		// Android: values-fr → "fr"
		{"res/values-fr/strings.xml", "fr"},
		{"res/values/strings.xml", "default"},
		// Flutter ARB: app_en.arb → "en"
		{"lib/l10n/app_en.arb", "en"},
		// Direct filename
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
