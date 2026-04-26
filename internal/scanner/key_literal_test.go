package scanner

import (
	"testing"
)

func TestKeyExistsInText(t *testing.T) {
	tests := []struct {
		name string
		text string
		key  string
		want bool
	}{
		{
			name: "single-quoted literal",
			text: `this.translate.instant('COMMON.SAVE')`,
			key:  "COMMON.SAVE",
			want: true,
		},
		{
			name: "double-quoted pipe",
			text: `{{ "COMMON.SAVE" | translate }}`,
			key:  "COMMON.SAVE",
			want: true,
		},
		{
			name: "no partial match — key is prefix of longer key",
			text: `this.t('COMMON.SAVE_BUTTON')`,
			key:  "COMMON.SAVE",
			want: false,
		},
		{
			name: "no partial match — key appears inside word",
			text: `xCOMMON.SAVEx`,
			key:  "COMMON.SAVE",
			want: false,
		},
		{
			name: "object property value",
			text: `label: 'MENU.ITEM'`,
			key:  "MENU.ITEM",
			want: true,
		},
		{
			name: "key at start of string",
			text: `COMMON.SAVE is the key`,
			key:  "COMMON.SAVE",
			want: true,
		},
		{
			name: "key at end of string",
			text: `the key is COMMON.SAVE`,
			key:  "COMMON.SAVE",
			want: true,
		},
		{
			name: "key followed by dot (longer sub-key)",
			text: `COMMON.SAVE.OK`,
			key:  "COMMON.SAVE",
			want: false,
		},
		{
			name: "not in text at all",
			text: `hello world`,
			key:  "COMMON.SAVE",
			want: false,
		},
		{
			name: "lowercase dot key",
			text: `t('common.save')`,
			key:  "common.save",
			want: true,
		},
		{
			name: "lowercase key partial no match",
			text: `t('common.save_button')`,
			key:  "common.save",
			want: false,
		},
		{
			name: "flutter tr() style",
			text: `'menu.title'.tr()`,
			key:  "menu.title",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyExistsInText(tt.text, tt.key)
			if got != tt.want {
				t.Errorf("keyExistsInText(%q, %q) = %v, want %v", tt.text, tt.key, got, tt.want)
			}
		})
	}
}
