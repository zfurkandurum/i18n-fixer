package scanner

import "testing"

func TestIsDynamicKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"common.save", false},
		{"errors.network.timeout", false},
		{"simple", false},
		{"save_button", false},
		{"errors.${code}", true},
		{"prefix" + "+" + "suffix", true},
		{"`template`", true},
		{"errorCode", true}, // camelCase = likely variable
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := IsDynamicKey(tt.key)
			if got != tt.expected {
				t.Errorf("IsDynamicKey(%q) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}
