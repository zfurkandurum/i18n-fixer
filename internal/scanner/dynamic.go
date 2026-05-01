package scanner

import "strings"

// IsDynamicKey checks if a key contains dynamic/computed parts
// that cannot be statically resolved.
func IsDynamicKey(key string) bool {
	// JS/TS/Dart/Kotlin/PHP template literal: ${variable}
	if strings.Contains(key, "${") {
		return true
	}

	// Dart/PHP bare-variable interpolation: $variable (without braces)
	if strings.Contains(key, "$") {
		return true
	}

	// String concatenation indicator
	if strings.Contains(key, "+") {
		return true
	}

	// Backtick (template literal wrapper)
	if strings.Contains(key, "`") {
		return true
	}

	// Ruby string interpolation: #{variable}
	if strings.Contains(key, "#{") {
		return true
	}

	// Swift string interpolation: \(variable)
	if strings.Contains(key, "\\(") {
		return true
	}

	// Python f-string interpolation or stray brace: {variable}
	// (Static i18n keys never legitimately contain `{`.)
	if strings.Contains(key, "{") {
		return true
	}

	// Key is a bare identifier (no dots, no quotes — just a variable name)
	// e.g., t(keyVariable) where key captured is "keyVariable"
	// Real keys typically contain dots or are all-lowercase with underscores
	if !strings.Contains(key, ".") && !strings.Contains(key, "_") && !strings.Contains(key, " ") {
		// Could be a variable reference, but also could be a simple flat key
		// Only flag if it looks like a camelCase variable
		if len(key) > 1 && key[0] >= 'a' && key[0] <= 'z' {
			for _, c := range key[1:] {
				if c >= 'A' && c <= 'Z' {
					return true // camelCase = likely variable
				}
			}
		}
	}

	return false
}
