package scanner

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// ScanKeyUsage scans a source file for i18n function calls and extracts keys.
//
// Returns three things:
//   - resolved static keys (definite key references)
//   - dynamic-key warnings (keys with runtime interpolation/concat)
//   - inferred dynamic prefixes — when a captured key is dynamic but a usable
//     static prefix can be extracted (e.g. `'a.b.${var}'` → `a.b.`), the
//     prefix is returned so the unused-key analyzer can mark all keys under
//     that prefix as used.
//
// `keySeparator` (typically ".") gates prefix inference: an extracted prefix
// is only emitted if it ends with the separator, to avoid false matches like
// `abc$d` extracting "abc" and shadowing unrelated top-level keys.
func ScanKeyUsage(filePath, keySeparator string, patterns []*regexp.Regexp) (
	[]types.UsedKey,
	[]types.DynamicKeyWarning,
	[]string,
) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil
	}
	defer f.Close()

	var keys []types.UsedKey
	var warnings []types.DynamicKeyWarning
	var inferred []string

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatchIndex(line, -1)
			for _, match := range matches {
				keyIdx := pattern.SubexpIndex("key")
				if keyIdx < 0 || keyIdx*2+1 >= len(match) {
					continue
				}

				start := match[keyIdx*2]
				end := match[keyIdx*2+1]
				if start < 0 || end < 0 {
					continue
				}

				key := line[start:end]
				rawMatch := line[match[0]:match[1]]

				if IsDynamicKey(key) {
					warnings = append(warnings, types.DynamicKeyWarning{
						RawExpression: rawMatch,
						File:          filePath,
						Line:          lineNum,
					})
					if p := extractStaticPrefix(key, keySeparator); p != "" {
						inferred = append(inferred, p)
					}
				} else {
					keys = append(keys, types.UsedKey{
						Key:      key,
						File:     filePath,
						Line:     lineNum,
						Column:   start + 1,
						RawMatch: rawMatch,
					})
				}
			}
		}
	}

	return keys, warnings, inferred
}

// dynamicMarkers lists substrings whose first occurrence in a captured i18n
// key marks the start of an interpolation/concat. They cover the common
// runtime composition forms across frameworks:
//
//   - `${`   JS/TS template literals, Dart, Kotlin, PHP
//   - `$`    Dart bare interpolation (`$var`), PHP `$var`
//   - `#{`   Ruby
//   - `\(`   Swift
//   - `{`    Python f-string (also covers `${` since `$` precedes `{`)
//   - `+`    string concatenation (e.g. `"a.b." + var`)
//   - "`"    stray template literal start (defensive)
//   - " "    whitespace inside a captured key implies concatenation/format
var dynamicMarkers = []string{"${", "$", "#{", "\\(", "{", "+", "`", " "}

// extractStaticPrefix returns the static portion of a captured-but-dynamic
// key up to (and including) the last keySeparator before the first dynamic
// marker. Returns "" when no usable prefix exists.
//
// Examples (separator="."):
//
//	"a.b.${var}"      -> "a.b."
//	"a.${var}.c"      -> "a."
//	"a.\\(var)"       -> "a."     // Swift
//	"a.#{var}"        -> "a."     // Ruby
//	"a.{var}"         -> "a."     // Python f-string
//	"a.b.c"           -> ""       // no dynamic marker
//	"${var}.foo"      -> ""       // no static prefix before marker
//	"abc$d"           -> ""       // doesn't end with separator
func extractStaticPrefix(key, sep string) string {
	if sep == "" {
		return ""
	}
	minIdx := -1
	for _, marker := range dynamicMarkers {
		i := strings.Index(key, marker)
		if i >= 0 && (minIdx == -1 || i < minIdx) {
			minIdx = i
		}
	}
	if minIdx <= 0 {
		return ""
	}
	static := key[:minIdx]
	if !strings.HasSuffix(static, sep) {
		return ""
	}
	return static
}

// ScanDynamicPrefixes scans a source file for dynamic prefix patterns and returns
// the unique prefixes found (named capture group "prefix"). These prefixes indicate
// that all i18n keys starting with the prefix are used dynamically at runtime.
func ScanDynamicPrefixes(filePath string, patterns []*regexp.Regexp) []string {
	if len(patterns) == 0 {
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	seen := make(map[string]bool)
	var prefixes []string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatchIndex(line, -1)
			for _, match := range matches {
				idx := pattern.SubexpIndex("prefix")
				if idx < 0 || idx*2+1 >= len(match) {
					continue
				}
				start := match[idx*2]
				end := match[idx*2+1]
				if start < 0 || end < 0 {
					continue
				}
				prefix := line[start:end]
				if prefix != "" && !seen[prefix] {
					seen[prefix] = true
					prefixes = append(prefixes, prefix)
				}
			}
		}
	}

	return prefixes
}

// CompilePatterns compiles regex pattern strings, skipping invalid ones.
func CompilePatterns(patternStrs []string) []*regexp.Regexp {
	var compiled []*regexp.Regexp
	for _, p := range patternStrs {
		re, err := regexp.Compile(p)
		if err != nil {
			continue
		}
		compiled = append(compiled, re)
	}
	return compiled
}
