package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// ScanResult holds all findings from scanning source files.
type ScanResult struct {
	UsedKeys    []types.UsedKey
	DynamicKeys []types.DynamicKeyWarning
	Hardcoded   []types.HardcodedString
}

// Scan orchestrates parallel scanning of source files.
func Scan(rootDir string, preset types.FrameworkPreset) (*ScanResult, error) {
	files, err := findSourceFiles(rootDir, preset.FileExtensions, preset.IgnorePatterns)
	if err != nil {
		return nil, err
	}

	keyPatterns := CompilePatterns(preset.I18nFunctionPatterns)
	hardcodedPatterns := CompilePatterns(preset.HardcodedStringPatterns)
	exclusionPatterns := CompilePatterns(preset.HardcodedStringExclusions)

	result := &ScanResult{}
	var mu sync.Mutex

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
				keys, dynWarnings := ScanKeyUsage(filePath, keyPatterns)
				hardcoded := ScanHardcoded(filePath, hardcodedPatterns, exclusionPatterns)

				mu.Lock()
				result.UsedKeys = append(result.UsedKeys, keys...)
				result.DynamicKeys = append(result.DynamicKeys, dynWarnings...)
				result.Hardcoded = append(result.Hardcoded, hardcoded...)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return result, nil
}

func findSourceFiles(rootDir string, extensions, ignorePatterns []string) ([]string, error) {
	extSet := make(map[string]bool)
	for _, ext := range extensions {
		extSet[ext] = true
	}

	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			// Quick skip for common directories
			if name == "node_modules" || name == ".git" || name == "dist" ||
				name == "build" || name == ".dart_tool" || name == "Pods" ||
				name == "DerivedData" || name == ".gradle" || name == ".next" ||
				name == ".nuxt" || name == ".svelte-kit" || name == ".angular" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check extension
		ext := filepath.Ext(path)
		if !extSet[ext] {
			return nil
		}

		// Skip large files (> 1MB)
		if info.Size() > 1024*1024 {
			return nil
		}

		// Check ignore patterns
		relPath, _ := filepath.Rel(rootDir, path)
		for _, pattern := range ignorePatterns {
			if matched, _ := filepath.Match(pattern, relPath); matched {
				return nil
			}
			// Also check against just the filename for simple patterns
			if matched, _ := filepath.Match(pattern, info.Name()); matched {
				return nil
			}
			// Check if path contains the ignore pattern directory
			if strings.Contains(pattern, "**") {
				cleanPattern := strings.ReplaceAll(pattern, "**", "")
				cleanPattern = strings.Trim(cleanPattern, "/")
				if cleanPattern != "" && strings.Contains(relPath, cleanPattern) {
					return nil
				}
			}
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
