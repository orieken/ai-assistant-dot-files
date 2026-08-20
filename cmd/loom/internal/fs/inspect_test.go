package frameworkfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectOwnedPathReportsBrokenSymlink(t *testing.T) {
	target := t.TempDir()
	link := filepath.Join(target, "agents")
	if err := os.Symlink(filepath.Join(target, "missing"), link); err != nil {
		t.Fatalf("create broken symlink: %v", err)
	}
	status, err := InspectOwnedPath(target, "agents")
	if err != nil {
		t.Fatalf("inspect owned path: %v", err)
	}
	if !status.Exists || !status.IsSymlink || !status.IsBroken {
		t.Fatalf("status = %+v", status)
	}
}

func TestInspectOwnedPathReportsMissingPath(t *testing.T) {
	status, err := InspectOwnedPath(t.TempDir(), ".claude/agents")
	if err != nil || status.Exists {
		t.Fatalf("status = %+v, err = %v", status, err)
	}
}
