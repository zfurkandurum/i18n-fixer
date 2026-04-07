package parser

import (
	"path/filepath"
	"testing"
)

func TestParseStrings(t *testing.T) {
	entries, err := ParseStrings(filepath.Join(testdataDir(), "Localizable.strings"))
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

	if keys["greeting"] != "Hello" {
		t.Errorf("expected greeting=Hello, got %q", keys["greeting"])
	}
	if keys["save_button"] != "Save Changes" {
		t.Errorf("expected save_button=Save Changes, got %q", keys["save_button"])
	}
}
