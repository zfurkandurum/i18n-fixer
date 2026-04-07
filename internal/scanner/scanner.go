package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
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

// hardcodedSkipDirs lists directories that are always skipped during source scanning.
var hardcodedSkipDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"dist":         true,
	"build":        true,
	".dart_tool":   true,
	"Pods":         true,
	"DerivedData":  true,
	".gradle":      true,
	".next":        true,
	".nuxt":        true,
	".svelte-kit":  true,
	".angular":     true,
	"vendor":       true,
	"__pycache__":  true,
	".build":       true,
	"coverage":     true,
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
			if hardcodedSkipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Check extension
		if !extSet[filepath.Ext(path)] {
			return nil
		}

		// Skip large files (> 1MB)
		if info.Size() > 1024*1024 {
			return nil
		}

		// Check ignore patterns using doublestar for proper ** glob support
		relPath, _ := filepath.Rel(rootDir, path)
		// Normalize to forward slashes for doublestar
		relPath = filepath.ToSlash(relPath)
		for _, pattern := range ignorePatterns {
			pattern = filepath.ToSlash(pattern)
			if matched, _ := doublestar.Match(pattern, relPath); matched {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
