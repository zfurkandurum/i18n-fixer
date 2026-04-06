package scanner

import (
	"path/filepath"
	"testing"
)

func TestScanHardcodedReact(t *testing.T) {
	patterns := CompilePatterns([]string{
		`>[\s]*(?P<str>[A-Z][a-zA-Z0-9 ,.!?'\-]{2,})[\s]*<`,
		`(?:placeholder|title|alt|aria-label|label)=['"](?P<str>[^'"]{2,})['"]`,
	})
	exclusions := CompilePatterns([]string{
		`^https?://`,
		`^[/#.]`,
		`className=`,
	})

	results := ScanHardcoded(filepath.Join(testdataDir(), "sample_react.tsx"), patterns, exclusions)

	if len(results) == 0 {
		t.Error("expected to find hardcoded strings")
	}

	// Should find strings like "Please enter your email address", "Loading data, please wait..."
	found := false
	for _, r := range results {
		if r.Value == "Type here..." {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find hardcoded placeholder 'Type here...'")
	}
}

func TestShouldExclude(t *testing.T) {
	exclusions := CompilePatterns([]string{
		`^https?://`,
		`^[0-9.,]+$`,
	})

	tests := []struct {
		str      string
		expected bool
	}{
		{"https://example.com", true},
		{"123.45", true},
		{"Hello World", false},
		{"a", true},  // too short
		{"42", true},  // no letters
		{"", true},    // empty
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := shouldExclude(tt.str, exclusions)
			if got != tt.expected {
				t.Errorf("shouldExclude(%q) = %v, want %v", tt.str, got, tt.expected)
			}
		})
	}
}
