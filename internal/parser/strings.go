package parser

import (
	"fmt"
	"os"
	"regexp"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

var stringsRegex = regexp.MustCompile(`"([^"\\]*(?:\\.[^"\\]*)*)"[\s]*=[\s]*"([^"\\]*(?:\\.[^"\\]*)*)";`)

// ParseStrings reads an iOS .strings file and returns entries.
func ParseStrings(filePath string) ([]types.I18nEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading .strings file: %w", err)
	}

	matches := stringsRegex.FindAllSubmatch(data, -1)
	entries := make([]types.I18nEntry, 0, len(matches))

	for _, match := range matches {
		if len(match) >= 3 {
			entries = append(entries, types.I18nEntry{
				Key:   string(match[1]),
				Value: string(match[2]),
			})
		}
	}

	return entries, nil
}
