package parser

import (
	"fmt"
	"os"
	"regexp"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

var xmlStringRegex = regexp.MustCompile(`<string\s+name="([^"]+)"[^>]*>([^<]*)</string>`)

// ParseXML reads an Android strings.xml file and returns entries.
func ParseXML(filePath string) ([]types.I18nEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading XML file: %w", err)
	}

	matches := xmlStringRegex.FindAllSubmatch(data, -1)
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
