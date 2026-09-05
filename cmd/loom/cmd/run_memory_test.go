package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/memory"
	"github.com/orieken/loom/internal/orchestrator"
)

// A halted run must be recorded. Most runs anyone investigates are halted
// ones, so a store that only knew about clean completions would miss the
// interesting half.
func TestAHaltedRunIsArchivedAndIngested(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	run.Dir = projectDir
	if output, err := run.CombinedOutput(); err == nil {
		t.Fatalf("run should halt at confirm-design\n%s", output)
	}

	// Archived beside the artifacts, in git rather than in a temporary
	// workspace — that copy is the durable record.
	for _, name := range []string{orchestrator.RunStateFileName, orchestrator.RunEventsFileName} {
		path := filepath.Join(projectDir, "docs", "features", "user-auth", name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s was not archived: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(projectDir, memory.DefaultDir, memory.DefaultFileName)); err != nil {
		t.Fatalf("no episodic store was written: %v", err)
	}
}

// Recording is best-effort and must never fail a run. An unwritable store
// path is the cheapest way to prove it.
func TestAnUnwritableStoreDoesNotChangeHowARunEnds(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)
	// Occupy the memory directory's path with a file.
	if err := os.MkdirAll(filepath.Join(projectDir, ".claude"), 0o755); err != nil {
		t.Fatalf("create .claude: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, memory.DefaultDir), []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("create blocker: %v", err)
	}

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	run.Dir = projectDir
	output, err := run.CombinedOutput()
	if err == nil {
		t.Fatalf("run should still halt at confirm-design\n%s", output)
	}
	assertStagesCompleted(t, projectDir, stagesBeforeFirstGate, true)
	if !strings.Contains(string(output), "memory: could not record this run") {
		t.Errorf("the failure was not reported:\n%s", output)
	}
}
