package frameworkfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveOwnedPathRemovesOnlyTargetPath(t *testing.T) {
	target := t.TempDir()
	owned := filepath.Join(target, ".claude", "agents")
	if err := os.MkdirAll(owned, 0o755); err != nil {
		t.Fatalf("create owned path: %v", err)
	}
	if removed, err := RemoveOwnedPath(target, ".claude/agents", false); err != nil || !removed {
		t.Fatalf("remove owned path: removed=%t err=%v", removed, err)
	}
	if _, err := os.Lstat(owned); !os.IsNotExist(err) {
		t.Fatalf("owned path remains: %v", err)
	}
}

func TestRemoveOwnedPathRejectsUnsafePaths(t *testing.T) {
	for _, path := range []string{".", "..", ".loom-manifest.json"} {
		if _, err := RemoveOwnedPath(t.TempDir(), path, false); err == nil {
			t.Fatalf("expected unsafe path %q to fail", path)
		}
	}
}

func TestRemoveOwnedPathDryRunPreservesPath(t *testing.T) {
	target := t.TempDir()
	owned := filepath.Join(target, "owned")
	if err := os.WriteFile(owned, []byte("owned"), 0o644); err != nil {
		t.Fatalf("create owned path: %v", err)
	}
	if removed, err := RemoveOwnedPath(target, "owned", true); err != nil || !removed {
		t.Fatalf("dry-run owned path: removed=%t err=%v", removed, err)
	}
	if _, err := os.Stat(owned); err != nil {
		t.Fatalf("dry run removed owned path: %v", err)
	}
}

func TestRemoveOwnedPathSkipsMissingParent(t *testing.T) {
	removed, err := RemoveOwnedPath(t.TempDir(), ".claude/agents", false)
	if err != nil || removed {
		t.Fatalf("missing path: removed=%t err=%v", removed, err)
	}
}

func TestRemoveOwnedPathRejectsSymlinkedParentEscape(t *testing.T) {
	target, external := t.TempDir(), t.TempDir()
	victim := filepath.Join(external, "victim")
	if err := os.WriteFile(victim, []byte("keep"), 0o644); err != nil {
		t.Fatalf("create external file: %v", err)
	}
	if err := os.Symlink(external, filepath.Join(target, "link")); err != nil {
		t.Fatalf("create escaping symlink: %v", err)
	}
	if _, err := RemoveOwnedPath(target, "link/victim", false); err == nil {
		t.Fatal("expected symlinked parent escape to fail")
	}
	if _, err := os.Stat(victim); err != nil {
		t.Fatalf("external file was removed: %v", err)
	}
}
