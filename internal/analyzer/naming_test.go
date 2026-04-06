package analyzer

import (
	"testing"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

func TestLintKeyNamingUpperSnake(t *testing.T) {
	entries := []types.I18nEntry{
		{Key: "AUTH.LOGIN_TITLE", File: "en.json", Locale: "en"},
		{Key: "COMMON.SAVE", File: "en.json", Locale: "en"},
		{Key: "auth.loginTitle", File: "en.json", Locale: "en"}, // violation
	}

	issues := LintKeyNaming(entries, "UPPER_SNAKE")

	if len(issues) != 1 {
		t.Fatalf("expected 1 naming issue, got %d", len(issues))
	}
	if issues[0].Key != "auth.loginTitle" {
		t.Errorf("expected key 'auth.loginTitle', got %q", issues[0].Key)
	}
}

func TestLintKeyNamingLowerDot(t *testing.T) {
	entries := []types.I18nEntry{
		{Key: "auth.login.title", File: "en.json", Locale: "en"},
		{Key: "common.save", File: "en.json", Locale: "en"},
		{Key: "AUTH.SAVE", File: "en.json", Locale: "en"}, // violation
	}

	issues := LintKeyNaming(entries, "lower.dot")

	if len(issues) != 1 {
		t.Fatalf("expected 1 naming issue, got %d", len(issues))
	}
	if issues[0].Key != "AUTH.SAVE" {
		t.Errorf("expected key 'AUTH.SAVE', got %q", issues[0].Key)
	}
}

func TestLintKeyNamingAutoDetect(t *testing.T) {
	// Mostly UPPER_SNAKE keys
	entries := []types.I18nEntry{
		{Key: "AUTH.LOGIN", File: "en.json", Locale: "en"},
		{Key: "AUTH.REGISTER", File: "en.json", Locale: "en"},
		{Key: "COMMON.SAVE", File: "en.json", Locale: "en"},
		{Key: "COMMON.DELETE", File: "en.json", Locale: "en"},
		{Key: "common.lowercase", File: "en.json", Locale: "en"}, // violation
	}

	issues := LintKeyNaming(entries, "") // auto-detect

	if len(issues) != 1 {
		t.Fatalf("expected 1 naming issue (auto-detected UPPER_SNAKE), got %d", len(issues))
	}
}

func TestLintKeyNamingNoEntries(t *testing.T) {
	issues := LintKeyNaming(nil, "")
	if issues != nil {
		t.Errorf("expected nil for empty entries, got %v", issues)
	}
}

func TestDetectConvention(t *testing.T) {
	entries := []types.I18nEntry{
		{Key: "AUTH.LOGIN"},
		{Key: "AUTH.REGISTER"},
		{Key: "COMMON.SAVE"},
		{Key: "COMMON.DELETE"},
		{Key: "COMMON.CANCEL"},
	}
	got := detectConvention(entries)
	if got != "UPPER_SNAKE" {
		t.Errorf("expected UPPER_SNAKE, got %q", got)
	}
}
