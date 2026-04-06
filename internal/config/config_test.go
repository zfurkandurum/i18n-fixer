package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNoConfigFile(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Preset != "" {
		t.Errorf("expected empty preset, got %q", cfg.Preset)
	}
}

func TestLoadValidConfig(t *testing.T) {
	dir := t.TempDir()
	configContent := `{"preset": "react-i18next", "defaultLocale": "en", "verbose": true}`
	err := os.WriteFile(filepath.Join(dir, ".i18n-fixer.json"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Preset != "react-i18next" {
		t.Errorf("expected preset %q, got %q", "react-i18next", cfg.Preset)
	}
	if cfg.DefaultLocale != "en" {
		t.Errorf("expected defaultLocale %q, got %q", "en", cfg.DefaultLocale)
	}
	if !cfg.Verbose {
		t.Error("expected verbose to be true")
	}
}

func TestLoadWalksUpDirectories(t *testing.T) {
	root := t.TempDir()
	subdir := filepath.Join(root, "src", "components")
	err := os.MkdirAll(subdir, 0755)
	if err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	configContent := `{"preset": "vue-i18n"}`
	err = os.WriteFile(filepath.Join(root, ".i18n-fixer.json"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(subdir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Preset != "vue-i18n" {
		t.Errorf("expected preset %q, got %q", "vue-i18n", cfg.Preset)
	}
}

func TestDefaults(t *testing.T) {
	cfg := Defaults()
	if cfg.Format != "console" {
		t.Errorf("expected default format %q, got %q", "console", cfg.Format)
	}
}
