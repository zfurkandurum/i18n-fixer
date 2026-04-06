package analyzer

import (
	"strings"
	"unicode"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// GroupHardcodedStrings deduplicates and groups hardcoded strings by value,
// generating a suggested i18n key for each.
func GroupHardcodedStrings(strings_ []types.HardcodedString, keySeparator string) []types.HardcodedStringIssue {
	grouped := make(map[string]*types.HardcodedStringIssue)

	for _, s := range strings_ {
		if existing, ok := grouped[s.Value]; ok {
			existing.Occurrences = append(existing.Occurrences, types.Location{
				File:   s.File,
				Line:   s.Line,
				Column: s.Column,
			})
		} else {
			grouped[s.Value] = &types.HardcodedStringIssue{
				Value:        s.Value,
				SuggestedKey: suggestKey(s.Value, keySeparator),
				Occurrences: []types.Location{
					{File: s.File, Line: s.Line, Column: s.Column},
				},
			}
		}
	}

	issues := make([]types.HardcodedStringIssue, 0, len(grouped))
	for _, issue := range grouped {
		issues = append(issues, *issue)
	}

	return issues
}

// suggestKey generates a suggested i18n key from a string value.
// "Please enter your email" → "please.enter.your.email"
func suggestKey(value, separator string) string {
	// Lowercase
	lower := strings.ToLower(value)

	// Remove non-alphanumeric characters except spaces
	var cleaned []rune
	for _, r := range lower {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			cleaned = append(cleaned, r)
		}
	}

	// Split into words, take first 4
	words := strings.Fields(string(cleaned))
	if len(words) > 4 {
		words = words[:4]
	}

	return strings.Join(words, separator)
}
