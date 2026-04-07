package analyzer

import (
	"testing"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

func TestGroupHardcodedStrings(t *testing.T) {
	strings := []types.HardcodedString{
		{Value: "Save Changes", File: "src/A.tsx", Line: 10},
		{Value: "Save Changes", File: "src/B.tsx", Line: 20},
		{Value: "Please enter email", File: "src/C.tsx", Line: 30},
	}

	issues := GroupHardcodedStrings(strings, ".")

	if len(issues) != 2 {
		t.Fatalf("expected 2 grouped issues, got %d", len(issues))
	}

	issueMap := make(map[string]types.HardcodedStringIssue)
	for _, issue := range issues {
		issueMap[issue.Value] = issue
	}

	if save, ok := issueMap["Save Changes"]; ok {
		if len(save.Occurrences) != 2 {
			t.Errorf("expected 2 occurrences of 'Save Changes', got %d", len(save.Occurrences))
		}
		if save.SuggestedKey != "save.changes" {
			t.Errorf("expected suggested key 'save.changes', got %q", save.SuggestedKey)
		}
	} else {
		t.Error("expected 'Save Changes' in issues")
	}
}

func TestSuggestKey(t *testing.T) {
	tests := []struct {
		value     string
		separator string
		expected  string
	}{
		{"Save Changes", ".", "save.changes"},
		{"Please enter your email address", ".", "please.enter.your.email"},
		{"Hello!", ".", "hello"},
		{"OK", "_", "ok"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got := suggestKey(tt.value, tt.separator)
			if got != tt.expected {
				t.Errorf("suggestKey(%q, %q) = %q, want %q", tt.value, tt.separator, got, tt.expected)
			}
		})
	}
}
