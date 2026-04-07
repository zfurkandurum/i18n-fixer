package scanner

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// ScanKeyLiterals scans source files for i18n key strings using word-boundary-aware
// literal search. This supplements pattern-based scanning and catches any key that
// appears as a literal string in source, regardless of the surrounding syntax or
// function name used to reference it.
//
// Keys already found by pattern scanning (patternFound) are skipped to avoid duplicates.
// i18nFilePaths are excluded so keys don't match inside their own definition files.
func ScanKeyLiterals(
	rootDir string,
	preset types.FrameworkPreset,
	allKeys []string,
	i18nFilePaths map[string]bool,
	patternFound map[string]bool,
) []types.UsedKey {
	if len(allKeys) == 0 {
		return nil
	}

	// Only search for keys not already found via pattern scanning
	keysToSearch := make([]string, 0, len(allKeys))
	for _, k := range allKeys {
		if !patternFound[k] {
			keysToSearch = append(keysToSearch, k)
		}
	}
	if len(keysToSearch) == 0 {
		return nil
	}

	files, err := findSourceFiles(rootDir, preset.FileExtensions, preset.IgnorePatterns)
	if err != nil || len(files) == 0 {
		return nil
	}

	var mu sync.Mutex
	foundKeys := make(map[string]string) // key → first file found in

	workers := runtime.NumCPU()
	if workers < 1 {
		workers = 1
	}

	fileCh := make(chan string, len(files))
	for _, f := range files {
		fileCh <- f
	}
	close(fileCh)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileCh {
				// Skip i18n definition files — every key exists there by definition
				if i18nFilePaths[filePath] {
					continue
				}

				content, err := os.ReadFile(filePath)
				if err != nil {
					continue
				}
				text := string(content)

				for _, key := range keysToSearch {
					mu.Lock()
					alreadyFound := foundKeys[key] != ""
					mu.Unlock()
					if alreadyFound {
						continue
					}
					if keyExistsInText(text, key) {
						mu.Lock()
						foundKeys[key] = filePath
						mu.Unlock()
					}
				}
			}
		}()
	}
	wg.Wait()

	var result []types.UsedKey
	for key, file := range foundKeys {
		result = append(result, types.UsedKey{
			Key:      key,
			File:     file,
			RawMatch: "(literal)",
		})
	}
	return result
}

// keyExistsInText returns true if key appears as a complete token in text.
// Uses word-boundary check: the byte immediately before and after the match
// must not be a key character ([A-Za-z0-9_.]) to avoid partial matches like
// finding COMMON.SAVE inside COMMON.SAVE_BUTTON.
func keyExistsInText(text, key string) bool {
	idx := 0
	keyLen := len(key)
	for {
		pos := strings.Index(text[idx:], key)
		if pos < 0 {
			return false
		}
		abs := idx + pos
		end := abs + keyLen

		beforeOK := abs == 0 || !isKeyChar(text[abs-1])
		afterOK := end >= len(text) || !isKeyChar(text[end])

		if beforeOK && afterOK {
			return true
		}
		idx = abs + 1
	}
}

// isKeyChar returns true for characters that can appear in an i18n key.
func isKeyChar(c byte) bool {
	return c == '.' || c == '_' ||
		(c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9')
}
