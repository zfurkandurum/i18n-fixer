package preset

import (
	"testing"
)

func TestAllPresetsLoaded(t *testing.T) {
	all := All()

	expectedPresets := []string{
		"react-i18next", "react-intl", "vue-i18n",
		"angular", "ngx-translate", "svelte-i18n",
		"next-intl", "nuxt-i18n", "ember-intl",
		"flutter", "ios-swift", "android-kotlin", "react-native",
	}

	for _, name := range expectedPresets {
		if _, ok := all[name]; !ok {
			t.Errorf("missing built-in preset: %s", name)
		}
	}

	if len(all) != len(expectedPresets) {
		t.Errorf("expected %d presets, got %d", len(expectedPresets), len(all))
	}
}

func TestGetPreset(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"react-i18next", false},
		{"flutter", false},
		{"ios-swift", false},
		{"nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Get(tt.name)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if p.Name != tt.name {
				t.Errorf("expected name %q, got %q", tt.name, p.Name)
			}
			if p.DisplayName == "" {
				t.Error("displayName should not be empty")
			}
			if len(p.FileExtensions) == 0 {
				t.Error("fileExtensions should not be empty")
			}
			if len(p.I18nFunctionPatterns) == 0 {
				t.Error("i18nFunctionPatterns should not be empty")
			}
			if len(p.ProjectMarkers) == 0 {
				t.Error("projectMarkers should not be empty")
			}
		})
	}
}

func TestPresetFieldsValid(t *testing.T) {
	for name, p := range All() {
		t.Run(name, func(t *testing.T) {
			if p.I18nFileFormat == "" {
				t.Error("i18nFileFormat should not be empty")
			}
			validFormats := map[string]bool{"json": true, "yaml": true, "xml": true, "strings": true, "arb": true}
			if !validFormats[p.I18nFileFormat] {
				t.Errorf("invalid i18nFileFormat: %s", p.I18nFileFormat)
			}
			if p.KeySeparator == "" {
				t.Error("keySeparator should not be empty")
			}
		})
	}
}
