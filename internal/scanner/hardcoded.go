package scanner

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// ScanHardcoded scans a source file for hardcoded user-facing strings.
func ScanHardcoded(filePath string, patterns []*regexp.Regexp, exclusions []*regexp.Regexp) []types.HardcodedString {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var results []types.HardcodedString

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip lines annotated with i18n-ignore
		if strings.Contains(line, "i18n-ignore") {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatchIndex(line, -1)
			for _, match := range matches {
				strIdx := pattern.SubexpIndex("str")
				if strIdx < 0 || strIdx*2+1 >= len(match) {
					continue
				}

				start := match[strIdx*2]
				end := match[strIdx*2+1]
				if start < 0 || end < 0 {
					continue
				}

				str := line[start:end]

				if shouldExclude(str, exclusions) {
					continue
				}

				results = append(results, types.HardcodedString{
					Value:   str,
					File:    filePath,
					Line:    lineNum,
					Column:  start + 1,
					Context: strings.TrimSpace(line),
				})
			}
		}
	}

	return results
}

func shouldExclude(str string, exclusions []*regexp.Regexp) bool {
	// Built-in filters
	str = strings.TrimSpace(str)

	// Too short
	if len(str) < 2 {
		return true
	}

	// No letters
	hasLetter := false
	for _, r := range str {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return true
	}

	// Dart/JS string interpolation: ${expr} or $variable
	if strings.Contains(str, "${") {
		return true
	}
	if idx := strings.Index(str, "$"); idx >= 0 {
		rest := []rune(str[idx+1:])
		if len(rest) > 0 && unicode.IsLetter(rest[0]) {
			return true
		}
	}

	// Check user-defined exclusion patterns
	for _, ex := range exclusions {
		if ex.MatchString(str) {
			return true
		}
	}

	return false
}
