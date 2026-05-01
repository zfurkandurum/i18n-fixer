package parser

import (
	"path/filepath"
	"testing"
)

func TestParseYAML(t *testing.T) {
	entries, err := ParseYAML(filepath.Join(testdataDir(), "en.yaml"), ".")
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
	if keys["common.cancel"] != "Cancel" {
		t.Errorf("expected common.cancel=Cancel, got %q", keys["common.cancel"])
	}
	if keys["errors.required"] != "This field is required" {
		t.Errorf("expected errors.required=This field is required, got %q", keys["errors.required"])
	}
}

func TestParseYAMLPlural(t *testing.T) {
	entries, err := ParseYAML(filepath.Join(testdataDir(), "plural.yaml"), ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys := make(map[string]string)
	for _, e := range entries {
		keys[e.Key] = e.Value
	}

	// items.count is a plural object — collapse to one entry, value = `other`.
	if v, ok := keys["items.count"]; !ok {
		t.Errorf("expected items.count to be present (collapsed plural), got map=%v", keys)
	} else if v != "{} items" {
		t.Errorf("expected items.count value=\"{} items\", got %q", v)
	}
	for _, sub := range []string{"items.count.zero", "items.count.one", "items.count.other"} {
		if _, ok := keys[sub]; ok {
			t.Errorf("did not expect sub-key %q (parent is plural)", sub)
		}
	}

	// Sibling flat key under same parent must still flatten normally.
	if keys["items.subtitle"] != "List of items" {
		t.Errorf("expected items.subtitle=\"List of items\", got %q", keys["items.subtitle"])
	}

	// one+other (no zero/two/few/many) — still valid ICU plural.
	if v, ok := keys["messages.unread"]; !ok || v != "{} unread messages" {
		t.Errorf("expected messages.unread to collapse to other-form, got %q (present=%v)", v, ok)
	}

	// Mixed keys (non-category sibling) — must NOT collapse.
	if _, ok := keys["size"]; ok {
		t.Errorf("did not expect size to be collapsed (mixed keys)")
	}
	if keys["size.small"] != "Small" || keys["size.other"] != "Other size" {
		t.Errorf("expected size.small and size.other to flatten normally, got %v", keys)
	}
}
