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
