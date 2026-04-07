package detect

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func TestDetectReactProject(t *testing.T) {
	dir := filepath.Join(testdataDir(), "react-project")
	presets, err := Detect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(presets) == 0 {
		t.Fatal("expected at least one preset, got none")
	}

	found := false
	for _, p := range presets {
		if p.Name == "react-i18next" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected react-i18next preset to be detected")
	}
}

func TestDetectFlutterProject(t *testing.T) {
	dir := filepath.Join(testdataDir(), "flutter-project")
	presets, err := Detect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(presets) == 0 {
		t.Fatal("expected at least one preset, got none")
	}

	found := false
	for _, p := range presets {
		if p.Name == "flutter" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected flutter preset to be detected")
	}
}

func TestDetectAndroidProject(t *testing.T) {
	dir := filepath.Join(testdataDir(), "android-project")
	presets, err := Detect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(presets) == 0 {
		t.Fatal("expected at least one preset, got none")
	}

	found := false
	for _, p := range presets {
		if p.Name == "android-kotlin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected android-kotlin preset to be detected")
	}
}

func TestDetectEmptyProject(t *testing.T) {
	dir := filepath.Join(testdataDir(), "empty-project")
	presets, err := Detect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(presets) != 0 {
		t.Errorf("expected no presets, got %d", len(presets))
	}
}
