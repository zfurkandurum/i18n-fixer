package parser

import (
	"path/filepath"
	"testing"
)

func TestParseXML(t *testing.T) {
	entries, err := ParseXML(filepath.Join(testdataDir(), "strings.xml"))
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

	if keys["app_name"] != "My App" {
		t.Errorf("expected app_name=My App, got %q", keys["app_name"])
	}
	if keys["hello_world"] != "Hello World" {
		t.Errorf("expected hello_world=Hello World, got %q", keys["hello_world"])
	}
}
