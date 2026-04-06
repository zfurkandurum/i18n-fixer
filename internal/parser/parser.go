package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// Parse reads an i18n file and returns flattened key-value entries.
// The format parameter determines which parser to use.
func Parse(filePath, format, keySeparator string) ([]types.I18nEntry, error) {
	locale := detectLocale(filePath)

	var entries []types.I18nEntry
	var err error

	switch format {
	case "json":
		entries, err = ParseJSON(filePath, keySeparator)
	case "yaml":
		entries, err = ParseYAML(filePath, keySeparator)
	case "xml":
		entries, err = ParseXML(filePath)
	case "strings":
		entries, err = ParseStrings(filePath)
	case "arb":
		entries, err = ParseARB(filePath)
	default:
		return nil, fmt.Errorf("unsupported i18n file format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	// Set locale and file path on all entries
	for i := range entries {
		entries[i].File = filePath
		entries[i].Locale = locale
	}

	return entries, nil
}

// detectLocale extracts the locale from a file path.
// Supports patterns like:
//   - src/locales/en/common.json → "en"
//   - res/values-fr/strings.xml → "fr"
//   - lib/l10n/app_en.arb → "en"
//   - en.json → "en"
func detectLocale(filePath string) string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Check parent directory name: locales/en/... or values-fr/...
	dirName := filepath.Base(dir)

	// Android pattern: values-fr → fr
	if strings.HasPrefix(dirName, "values-") {
		return strings.TrimPrefix(dirName, "values-")
	}
	if dirName == "values" {
		return "default"
	}

	// Common pattern: locales/en, lang/fr, etc.
	// But exclude known non-locale directory names
	nonLocale := map[string]bool{
		"l10n": true, "i18n": true, "locales": true, "locale": true,
		"lang": true, "translations": true, "messages": true, "src": true,
		"assets": true, "public": true, "lib": true, "res": true,
	}
	if len(dirName) >= 2 && len(dirName) <= 5 && !nonLocale[dirName] {
		return dirName
	}

	// Flutter ARB pattern: app_en.arb → en
	if parts := strings.Split(name, "_"); len(parts) >= 2 {
		lastPart := parts[len(parts)-1]
		if len(lastPart) >= 2 && len(lastPart) <= 5 {
			return lastPart
		}
	}

	// Filename is the locale: en.json → en
	if len(name) >= 2 && len(name) <= 5 {
		return name
	}

	return "unknown"
}
