package scanner

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func TestScanKeyUsageReact(t *testing.T) {
	patterns := CompilePatterns([]string{
		`\bt\(['"](?P<key>[^'"]+)['"]`,
	})

	keys, warnings := ScanKeyUsage(filepath.Join(testdataDir(), "sample_react.tsx"), patterns)

	// Should find: dashboard.title, dashboard.welcome, common.save, static.key.here
	// t(`errors.${errorCode}`) uses backticks, not quotes — regex won't match it
	// That's correct behavior: backtick-based calls need separate handling
	if len(keys) < 4 {
		t.Errorf("expected at least 4 keys, got %d", len(keys))
	}

	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k.Key] = true
	}

	if !keySet["dashboard.title"] {
		t.Error("expected to find key 'dashboard.title'")
	}
	if !keySet["common.save"] {
		t.Error("expected to find key 'common.save'")
	}
	if !keySet["static.key.here"] {
		t.Error("expected to find key 'static.key.here'")
	}

	_ = warnings // dynamic keys from backtick calls are not captured by quote-based patterns
}

func TestCompilePatternsSkipsInvalid(t *testing.T) {
	patterns := CompilePatterns([]string{
		`\bt\(['"](?P<key>[^'"]+)['"]`,
		`[invalid regex`,
		`valid\.pattern`,
	})

	if len(patterns) != 2 {
		t.Errorf("expected 2 valid patterns, got %d", len(patterns))
	}
}
