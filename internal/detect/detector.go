package detect

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/i18n-fixer/i18n-fixer/internal/preset"
	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// Detect auto-detects frameworks in the given project directory.
// Returns all matching presets (supports monorepos with multiple frameworks).
func Detect(rootDir string) ([]types.FrameworkPreset, error) {
	allPresets := preset.All()
	var matched []types.FrameworkPreset

	for _, p := range allPresets {
		if matchesProject(rootDir, p.ProjectMarkers) {
			matched = append(matched, p)
		}
	}

	// If both a specific and general preset match, prefer the specific one.
	// e.g., react-i18next > react-native (both have package.json markers)
	matched = deduplicatePresets(matched)

	return matched, nil
}

func matchesProject(rootDir string, markers []types.ProjectMarker) bool {
	for _, marker := range markers {
		if matchesMarker(rootDir, marker) {
			return true
		}
	}
	return false
}

func matchesMarker(rootDir string, marker types.ProjectMarker) bool {
	markerPath := filepath.Join(rootDir, marker.File)

	// Check for glob-like markers (e.g., "*.xcodeproj")
	if strings.Contains(marker.File, "*") {
		matches, err := filepath.Glob(markerPath)
		if err != nil || len(matches) == 0 {
			return false
		}
		// Glob markers don't need content checks
		return len(marker.ContainsAny) == 0
	}

	data, err := os.ReadFile(markerPath)
	if err != nil {
		return false
	}

	if len(marker.ContainsAny) == 0 {
		return true
	}

	content := string(data)
	for _, keyword := range marker.ContainsAny {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	return false
}

// deduplicatePresets removes less specific presets when a more specific one matches.
// For example, if both react-i18next and react-native match, keep both
// (they could be in a monorepo), but if react-i18next matches, don't also
// include a generic "i18next" preset if one existed.
func deduplicatePresets(presets []types.FrameworkPreset) []types.FrameworkPreset {
	if len(presets) <= 1 {
		return presets
	}
	// For now, return all matches. Monorepo support means multiple presets are valid.
	return presets
}
