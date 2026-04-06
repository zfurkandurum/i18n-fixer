package scanner

import (
	"bufio"
	"os"
	"regexp"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
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
