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

func TestScanHardcodedDartEnumsAndReturns(t *testing.T) {
	// Patterns mirroring the flutter-easy-localization preset additions:
	// arrow-syntax switch arms, return statements, throw/Exception, and the
	// extended named-parameter list (label, subtitle, helperText, etc.).
	patterns := CompilePatterns([]string{
		`Text\([\s]*['"](?P<str>[^'"]{2,})['"]`,
		`label:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`subtitle:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`helperText:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`errorText:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`placeholder:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`actionText:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`description:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`text:[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`=>[\s]*['"](?P<str>[^'"]{2,})['"]`,
		`return[\s]+['"](?P<str>[^'"]{2,})['"]`,
		`throw[\s]+['"](?P<str>[^'"]{2,})['"]`,
		`Exception\([\s]*['"](?P<str>[^'"]{2,})['"]`,
	})
	exclusions := CompilePatterns([]string{
		`^https?://`,
		`^assets/`,
		`^[0-9.,]+$`,
		`^\s*$`,
		`^[a-z_]+$`,
		`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z0-9_]+)+$`,
	})

	results := ScanHardcoded(
		filepath.Join(testdataDir(), "sample_dart_enums.dart"),
		patterns, exclusions,
	)

	// Every flagged string must start with FOUND_; nothing OK_ should escape.
	for _, r := range results {
		if !startsWith(r.Value, "FOUND_") {
			t.Errorf("unexpected hardcoded match %q at line %d (must start with FOUND_)", r.Value, r.Line)
		}
	}

	// Ensure each expected hit IS present.
	want := []string{
		"FOUND_To Do",
		"FOUND_In Progress",
		"FOUND_Done",
		"FOUND_Information",
		"FOUND_Warning",
		"FOUND_Something failed",
		"FOUND_Bad input",
		"FOUND_Hello world",
		"FOUND_Enter your name",
		"FOUND_Invalid email",
		"FOUND_e.g. Jane Doe",
		"FOUND_All Items",
		"FOUND_Last seen today",
		"FOUND_Add Entry",
		"FOUND_Brief explanation",
		"FOUND_Submit",
	}
	for _, w := range want {
		hit := false
		for _, r := range results {
			if r.Value == w {
				hit = true
				break
			}
		}
		if !hit {
			t.Errorf("expected to detect %q but did not (results: %d)", w, len(results))
		}
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
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
		{"42", true}, // no letters
		{"", true},   // empty
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
