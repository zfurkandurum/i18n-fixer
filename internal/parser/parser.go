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

// Known ISO 639-1 language codes (2-letter) — the authoritative list.
// No hardcoded namespace exclusions — we only match real locale codes.
var knownLocaleCodes = map[string]bool{
	"aa": true, "ab": true, "af": true, "ak": true, "am": true, "an": true,
	"ar": true, "as": true, "av": true, "ay": true, "az": true, "ba": true,
	"be": true, "bg": true, "bh": true, "bi": true, "bm": true, "bn": true,
	"bo": true, "br": true, "bs": true, "ca": true, "ce": true, "ch": true,
	"co": true, "cr": true, "cs": true, "cu": true, "cv": true, "cy": true,
	"da": true, "de": true, "dv": true, "dz": true, "ee": true, "el": true,
	"en": true, "eo": true, "es": true, "et": true, "eu": true, "fa": true,
	"ff": true, "fi": true, "fj": true, "fo": true, "fr": true, "fy": true,
	"ga": true, "gd": true, "gl": true, "gn": true, "gu": true, "gv": true,
	"ha": true, "he": true, "hi": true, "ho": true, "hr": true, "ht": true,
	"hu": true, "hy": true, "hz": true, "ia": true, "id": true, "ie": true,
	"ig": true, "ii": true, "ik": true, "io": true, "is": true, "it": true,
	"iu": true, "ja": true, "jv": true, "ka": true, "kg": true, "ki": true,
	"kj": true, "kk": true, "kl": true, "km": true, "kn": true, "ko": true,
	"kr": true, "ks": true, "ku": true, "kv": true, "kw": true, "ky": true,
	"la": true, "lb": true, "lg": true, "li": true, "ln": true, "lo": true,
	"lt": true, "lu": true, "lv": true, "mg": true, "mh": true, "mi": true,
	"mk": true, "ml": true, "mn": true, "mr": true, "ms": true, "mt": true,
	"my": true, "na": true, "nb": true, "nd": true, "ne": true, "ng": true,
	"nl": true, "nn": true, "no": true, "nr": true, "nv": true, "ny": true,
	"oc": true, "oj": true, "om": true, "or": true, "os": true, "pa": true,
	"pi": true, "pl": true, "ps": true, "pt": true, "qu": true, "rm": true,
	"rn": true, "ro": true, "ru": true, "rw": true, "sa": true, "sc": true,
	"sd": true, "se": true, "sg": true, "si": true, "sk": true, "sl": true,
	"sm": true, "sn": true, "so": true, "sq": true, "sr": true, "ss": true,
	"st": true, "su": true, "sv": true, "sw": true, "ta": true, "te": true,
	"tg": true, "th": true, "ti": true, "tk": true, "tl": true, "tn": true,
	"to": true, "tr": true, "ts": true, "tt": true, "tw": true, "ty": true,
	"ug": true, "uk": true, "ur": true, "uz": true, "ve": true, "vi": true,
	"vo": true, "wa": true, "wo": true, "xh": true, "yi": true, "yo": true,
	"za": true, "zh": true, "zu": true,
}

// isLocaleCode checks if a string is a valid BCP 47 locale code.
// Matches: "en", "fr", "zh-CN", "pt-BR" etc.
// Does NOT match arbitrary short strings like "auth", "common", "app".
func isLocaleCode(s string) bool {
	// Direct 2-letter ISO 639-1 match
	if knownLocaleCodes[strings.ToLower(s)] {
		return true
	}

	// Regional variant: en-US, zh-CN, pt-BR (lang-REGION)
	if parts := strings.SplitN(s, "-", 2); len(parts) == 2 {
		lang := strings.ToLower(parts[0])
		region := parts[1]
		if knownLocaleCodes[lang] && len(region) >= 2 && len(region) <= 3 {
			return true
		}
	}

	return false
}

// detectLocale extracts the locale from a file path.
// Uses only real ISO 639-1 codes — no hardcoded namespace exclusions.
func detectLocale(filePath string) string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	dirName := filepath.Base(dir)

	// Android pattern: values-fr → fr
	if strings.HasPrefix(dirName, "values-") {
		return strings.TrimPrefix(dirName, "values-")
	}
	if dirName == "values" {
		return "default"
	}

	// Priority 1: Filename is a locale code (en.json, fr.json, zh-CN.json)
	if isLocaleCode(name) {
		return name
	}

	// Priority 2: Suffix after underscore (app_en.arb → en)
	if parts := strings.Split(name, "_"); len(parts) >= 2 {
		lastPart := parts[len(parts)-1]
		if isLocaleCode(lastPart) {
			return lastPart
		}
	}

	// Priority 3: Parent directory is a locale code (locales/en/common.json)
	if isLocaleCode(dirName) {
		return dirName
	}

	return "unknown"
}
