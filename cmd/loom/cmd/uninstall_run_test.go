package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/loom/cmd/loom/internal/manifest"
)

func TestExecuteUninstallRemovesManifestOwnedPaths(t *testing.T) {
	target := t.TempDir()
	writeUninstallFixture(t, target)
	var output bytes.Buffer
	request := uninstallRequest{target: target, display: target}
	if err := executeUninstall(request, &output); err != nil {
		t.Fatalf("execute uninstall: %v", err)
	}
	assertPathMissing(t, filepath.Join(target, ".claude", "agents"))
	assertPathMissing(t, filepath.Join(target, ".cursor", "agents"))
	assertPathMissing(t, filepath.Join(target, manifest.Filename))
	if !strings.Contains(output.String(), "2 paths removed. Manifest deleted.") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestExecuteUninstallDryRunPreservesInstall(t *testing.T) {
	target := t.TempDir()
	writeUninstallFixture(t, target)
	var output bytes.Buffer
	request := uninstallRequest{target: target, display: target, isDryRun: true}
	if err := executeUninstall(request, &output); err != nil {
		t.Fatalf("execute dry run: %v", err)
	}
	assertPathExists(t, filepath.Join(target, ".claude", "agents"))
	assertPathExists(t, filepath.Join(target, manifest.Filename))
	if !strings.Contains(output.String(), "2 paths would be removed. Manifest would be deleted.") {
		t.Fatalf("unexpected output: %s", output.String())
	}
}

func TestExecuteUninstallFiltersPlatformAndRetainsManifest(t *testing.T) {
	target := t.TempDir()
	writeUninstallFixture(t, target)
	request := uninstallRequest{target: target, display: target, platform: "cursor"}
	if err := executeUninstall(request, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute filtered uninstall: %v", err)
	}
	assertPathExists(t, filepath.Join(target, ".claude", "agents"))
	assertPathMissing(t, filepath.Join(target, ".cursor", "agents"))
	installed, err := manifest.Read(target)
	if err != nil {
		t.Fatalf("read retained manifest: %v", err)
	}
	if len(installed.Platforms) != 1 || installed.Platforms[0].Name != "claude-code" {
		t.Fatalf("remaining records = %+v", installed.Platforms)
	}
}

func TestExecuteUninstallReportsMissingManifest(t *testing.T) {
	target := t.TempDir()
	request := uninstallRequest{target: target, display: target}
	err := executeUninstall(request, &bytes.Buffer{})
	want := "loom: no manifest found at " + filepath.Join(target, manifest.Filename)
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %v, want %q", err, want)
	}
}

func writeUninstallFixture(t *testing.T, target string) {
	t.Helper()
	records := []manifest.PlatformRecord{
		{Name: "claude-code", Paths: []string{".claude/agents"}},
		{Name: "cursor", Paths: []string{".cursor/agents"}},
	}
	for _, record := range records {
		if err := os.MkdirAll(filepath.Join(target, record.Paths[0]), 0o755); err != nil {
			t.Fatalf("create %s: %v", record.Paths[0], err)
		}
	}
	if err := manifest.Write(target, manifest.Manifest{Version: "v3.3.14", Platforms: records}); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); err != nil {
		t.Fatalf("expected %s: %v", path, err)
	}
}

func assertPathMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be missing: %v", path, err)
	}
}
