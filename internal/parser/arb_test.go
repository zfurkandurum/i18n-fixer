package parser

import (
	"path/filepath"
	"testing"
)

func TestParseARB(t *testing.T) {
	entries, err := ParseARB(filepath.Join(testdataDir(), "app_en.arb"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys := make(map[string]string)
	for _, e := range entries {
		keys[e.Key] = e.Value
	}

	// Should have 3 entries: helloWorld, saveButton, itemCount
	// @-prefixed keys (@@locale, @helloWorld, @saveButton) should be excluded
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries (excluding @ keys), got %d: %v", len(entries), keys)
	}

	if keys["helloWorld"] != "Hello World" {
		t.Errorf("expected helloWorld=Hello World, got %q", keys["helloWorld"])
	}
	if keys["saveButton"] != "Save" {
		t.Errorf("expected saveButton=Save, got %q", keys["saveButton"])
	}

	// Verify @-prefixed keys are excluded
	if _, ok := keys["@@locale"]; ok {
		t.Error("@@locale should be excluded")
	}
	if _, ok := keys["@helloWorld"]; ok {
		t.Error("@helloWorld metadata should be excluded")
	}
}
