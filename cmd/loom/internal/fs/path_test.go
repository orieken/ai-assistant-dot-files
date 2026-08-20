package frameworkfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTargetRejectsFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "target.txt")
	if err := os.WriteFile(file, []byte("target"), 0o644); err != nil {
		t.Fatalf("create target file: %v", err)
	}
	if _, err := ResolveTarget(file); err == nil {
		t.Fatal("expected a non-directory target error")
	}
}

func TestSafeDestinationRejectsTraversal(t *testing.T) {
	if _, err := safeDestination(t.TempDir(), "../escape"); err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
}
