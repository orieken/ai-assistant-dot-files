package platform

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetectFindsMarkerDirectoriesAndFiles(t *testing.T) {
	target := t.TempDir()
	createMarkerDirectory(t, target, ".claude")
	createMarkerDirectory(t, target, ".github")
	createMarkerFile(t, target, "AGENTS.md")
	detected, err := Detect(target)
	if err != nil {
		t.Fatalf("detect platforms: %v", err)
	}
	want := []string{"claude-code", "github-copilot", "openai-codex"}
	if !reflect.DeepEqual(detected, want) {
		t.Fatalf("detected = %v, want %v", detected, want)
	}
}

func TestDetectRequiresDirectoryMarkersToBeDirectories(t *testing.T) {
	target := t.TempDir()
	createMarkerFile(t, target, ".cursor")
	detected, err := Detect(target)
	if err != nil {
		t.Fatalf("detect platforms: %v", err)
	}
	if len(detected) != 0 {
		t.Fatalf("detected file as directory marker: %v", detected)
	}
}

func createMarkerDirectory(t *testing.T, target, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(target, path), 0o755); err != nil {
		t.Fatalf("create marker directory: %v", err)
	}
}

func createMarkerFile(t *testing.T, target, path string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(target, path), []byte("marker"), 0o644); err != nil {
		t.Fatalf("create marker file: %v", err)
	}
}
