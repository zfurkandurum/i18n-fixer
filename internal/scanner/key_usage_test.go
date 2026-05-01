package scanner

import (
	"os"
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

	keys, warnings, _ := ScanKeyUsage(filepath.Join(testdataDir(), "sample_react.tsx"), ".", patterns)

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

func TestExtractStaticPrefix(t *testing.T) {
	tests := []struct {
		name string
		key  string
		sep  string
		want string
	}{
		{"js/dart template literal", "a.b.${var}", ".", "a.b."},
		{"middle interpolation", "a.${var}.c", ".", "a."},
		{"swift backslash-paren", "a.\\(var)", ".", "a."},
		{"ruby pound-brace", "a.#{var}", ".", "a."},
		{"python f-string", "a.{var}", ".", "a."},
		{"plain bare-dollar var (Dart)", "a.b.$var", ".", "a.b."},
		{"string concat", "a.b. + var", ".", "a.b."},
		{"no dynamic marker", "a.b.c", ".", ""},
		{"no static prefix before marker", "${var}.foo", ".", ""},
		{"no separator before marker", "abc$d", ".", ""},
		{"empty key", "", ".", ""},
		{"empty separator", "a.b.${var}", "", ""},
		{"slash separator", "a/b/${var}", "/", "a/b/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractStaticPrefix(tt.key, tt.sep)
			if got != tt.want {
				t.Errorf("extractStaticPrefix(%q, %q) = %q, want %q", tt.key, tt.sep, got, tt.want)
			}
		})
	}
}

// TestScanKeyUsageInfersPrefixes verifies that interpolated keys captured by
// loose-regex presets (like i18next) yield a static prefix in the third
// return value, and that the captured key shows up as a dynamic warning
// (not a definite static key).
func TestScanKeyUsageInfersPrefixes(t *testing.T) {
	tmp := t.TempDir()

	// Loose regex matching `t('...')` calls including interpolated content.
	patterns := CompilePatterns([]string{
		`\bt\(['"](?P<key>[^'"]+)['"]`,
	})

	cases := []struct {
		name        string
		content     string
		wantPrefix  string
		wantWarning bool
	}{
		{
			name:        "js template-literal interpolation in single quotes",
			content:     "t('a.b.${var}')",
			wantPrefix:  "a.b.",
			wantWarning: true,
		},
		{
			name:        "ruby pound-brace",
			content:     `t("a.#{var}")`,
			wantPrefix:  "a.",
			wantWarning: true,
		},
		{
			name:        "static call yields no prefix and no warning",
			content:     "t('a.b.c')",
			wantPrefix:  "",
			wantWarning: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmp, "src.txt")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("write fixture: %v", err)
			}

			_, warnings, inferred := ScanKeyUsage(path, ".", patterns)

			if tt.wantPrefix == "" {
				if len(inferred) != 0 {
					t.Errorf("expected no inferred prefix, got %v", inferred)
				}
			} else {
				found := false
				for _, p := range inferred {
					if p == tt.wantPrefix {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected inferred prefix %q, got %v", tt.wantPrefix, inferred)
				}
			}

			gotWarning := len(warnings) > 0
			if gotWarning != tt.wantWarning {
				t.Errorf("warning presence: got %v, want %v (warnings=%v)", gotWarning, tt.wantWarning, warnings)
			}
		})
	}
}
