package parser

import (
	"encoding/json"
	"os"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// xcstringsFile represents the top-level structure of an Apple String Catalog (.xcstrings).
type xcstringsFile struct {
	SourceLanguage string                    `json:"sourceLanguage"`
	Strings        map[string]xcstringsEntry `json:"strings"`
	Version        string                    `json:"version"`
}

type xcstringsEntry struct {
	Comment       string                           `json:"comment"`
	Localizations map[string]xcstringsLocalization `json:"localizations"`
}

type xcstringsLocalization struct {
	StringUnit *xcstringsStringUnit `json:"stringUnit"`
	Variations *xcstringsVariations `json:"variations"`
}

type xcstringsStringUnit struct {
	State string `json:"state"`
	Value string `json:"value"`
}

type xcstringsVariations struct {
	Plural map[string]xcstringsPluralVariation `json:"plural"`
}

type xcstringsPluralVariation struct {
	StringUnit xcstringsStringUnit `json:"stringUnit"`
}

// ParseXCStrings parses an Apple String Catalog file (Xcode 15+, .xcstrings format).
// Each key can have per-locale translations and optional plural variations.
func ParseXCStrings(filePath string) ([]types.I18nEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var catalog xcstringsFile
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, err
	}

	var entries []types.I18nEntry

	for key, entry := range catalog.Strings {
		for locale, loc := range entry.Localizations {
			if loc.StringUnit != nil && loc.StringUnit.Value != "" {
				entries = append(entries, types.I18nEntry{
					Key:    key,
					Value:  loc.StringUnit.Value,
					File:   filePath,
					Locale: locale,
				})
			}

			// Plural variations: store each form as key.pluralForm
			if loc.Variations != nil && loc.Variations.Plural != nil {
				for pluralForm, variation := range loc.Variations.Plural {
					if variation.StringUnit.Value != "" {
						entries = append(entries, types.I18nEntry{
							Key:    key + "." + pluralForm,
							Value:  variation.StringUnit.Value,
							File:   filePath,
							Locale: locale,
						})
					}
				}
			}
		}

		// If a key has no localizations at all, still register it so missing-key
		// analysis can detect it as untranslated.
		if len(entry.Localizations) == 0 {
			entries = append(entries, types.I18nEntry{
				Key:    key,
				Value:  "",
				File:   filePath,
				Locale: catalog.SourceLanguage,
			})
		}
	}

	return entries, nil
}
