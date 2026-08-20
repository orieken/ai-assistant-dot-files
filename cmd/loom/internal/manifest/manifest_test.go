package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteCreatesExpectedJSON(t *testing.T) {
	target := t.TempDir()
	installedAt := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	installed := Manifest{Version: "0.1.0", InstalledAt: installedAt, Platforms: []PlatformRecord{{Name: "claude-code", Paths: []string{".claude/agents"}}}}
	if err := Write(target, installed); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	assertManifestContent(t, filepath.Join(target, Filename))
}

func assertManifestContent(t *testing.T, path string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	want := "\"installedAt\": \"2026-08-16T12:00:00Z\""
	if !strings.Contains(string(content), want) {
		t.Fatalf("manifest %q does not contain %q", content, want)
	}
}

func TestReadIfExistsDistinguishesMissingManifest(t *testing.T) {
	if _, exists, err := ReadIfExists(t.TempDir()); err != nil || exists {
		t.Fatalf("missing manifest: exists=%t err=%v", exists, err)
	}
}

func TestReadReturnsWrittenManifest(t *testing.T) {
	target := t.TempDir()
	want := Manifest{Version: "0.1.0", Platforms: []PlatformRecord{{Name: "cursor"}}}
	if err := Write(target, want); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	got, err := Read(target)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if got.Version != want.Version || len(got.Platforms) != 1 || got.Platforms[0].Name != "cursor" {
		t.Fatalf("manifest = %+v, want %+v", got, want)
	}
}
