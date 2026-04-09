package scanner

import (
	"bufio"
	"os"
	"regexp"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// ScanKeyUsage scans a source file for i18n function calls and extracts keys.
func ScanKeyUsage(filePath string, patterns []*regexp.Regexp) ([]types.UsedKey, []types.DynamicKeyWarning) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	var keys []types.UsedKey
	var warnings []types.DynamicKeyWarning

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

	return keys, warnings
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
