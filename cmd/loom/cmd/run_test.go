package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
)

func writeSpec(t *testing.T, dir string) string {
	t.Helper()
	spec := filepath.Join(dir, "user-auth.md")
	if err := os.WriteFile(spec, []byte("# Feature: user auth\n"), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	return spec
}

func loadRunState(t *testing.T, projectDir string) *orchestrator.RunState {
	t.Helper()
	path := filepath.Join(projectDir, ".claude", "feature-workspace", "user-auth", orchestrator.RunStateFileName)
	state, err := orchestrator.NewStateStore(path).Load()
	if err != nil {
		t.Fatalf("load run state: %v", err)
	}
	return state
}

// stagesBeforeFirstGate are the built-in plan's stages ahead of the
// confirm-design gate on the developer stage.
var stagesBeforeFirstGate = []string{
	"context-engineer", "analyst", "architect", "performance-engineer", "data-engineer",
}

// TestRunMockProviderHaltsAtFirstGate runs the real built binary against the
// mock provider: the run executes every ungated stage, halts at the
// confirm-design gate, and refuses to restart over the resulting state.
// Same build-the-binary pattern as mcp_serve_test.go.
func TestRunMockProviderHaltsAtFirstGate(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	run.Dir = projectDir
	output, err := run.CombinedOutput()
	if err == nil {
		t.Fatalf("run completed without approval; it must halt at confirm-design\n%s", output)
	}

	assertStagesCompleted(t, projectDir, stagesBeforeFirstGate, true)
	assertWaitingOnGate(t, projectDir, "developer", "confirm-design")
	assertFreshRunOverStateRefused(t, binary, projectDir, spec)
}

func assertStagesCompleted(t *testing.T, projectDir string, stageIDs []string, wantArtifacts bool) {
	t.Helper()
	state := loadRunState(t, projectDir)
	if state == nil {
		t.Fatal("no run state persisted")
	}
	for _, stageID := range stageIDs {
		if !state.IsStageCompleted(stageID) {
			t.Errorf("stage %q not COMPLETED", stageID)
		}
		if wantArtifacts {
			assertArtifactExists(t, projectDir, stageID)
		}
	}
}

func assertWaitingOnGate(t *testing.T, projectDir, stageID, gate string) {
	t.Helper()
	state := loadRunState(t, projectDir)
	if got := state.Stages[stageID].Status; got != orchestrator.StageStatusWaitingApproval {
		t.Fatalf("stage %q status = %q, want WAITING_APPROVAL", stageID, got)
	}
	if got := state.WaitingGate(); got != gate {
		t.Errorf("run waiting on gate %q, want %q", got, gate)
	}
	if len(state.Approvals) != 0 {
		t.Errorf("run recorded approvals %v without any human approving", state.Approvals)
	}
}

func assertArtifactExists(t *testing.T, projectDir, stageID string) {
	t.Helper()
	artifact := filepath.Join(projectDir, ".claude", "feature-workspace", "user-auth", stageID+".md")
	if _, err := os.Stat(artifact); err != nil {
		t.Errorf("stage %q artifact missing: %v", stageID, err)
	}
}

func assertFreshRunOverStateRefused(t *testing.T, binary, projectDir, spec string) {
	t.Helper()
	rerun := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	rerun.Dir = projectDir
	output, err := rerun.CombinedOutput()
	if err == nil {
		t.Fatal("second run without --resume succeeded; must refuse to start over existing state")
	}
	if !strings.Contains(string(output), "--resume") {
		t.Errorf("refusal message should point at --resume, got:\n%s", output)
	}
}

// TestRunInterruptCheckpointAndResume kills a run mid-stage with SIGINT
// through the real CLI, asserts a resumable checkpoint was persisted, then
// resumes to completion — the M0.4 done-when criterion.
func TestRunInterruptCheckpointAndResume(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("SIGINT delivery is POSIX-only in this test")
	}
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)

	interruptRunMidStage(t, binary, projectDir, spec)
	assertInterruptCheckpoint(t, projectDir)

	resume := exec.Command(binary, "run", "--spec", spec, "--provider", "mock", "--resume")
	resume.Dir = projectDir
	if output, err := resume.CombinedOutput(); err == nil {
		t.Fatalf("resumed run completed without approval; it must halt at confirm-design\n%s", output)
	}
	assertStagesCompleted(t, projectDir, stagesBeforeFirstGate, false)
	assertWaitingOnGate(t, projectDir, "developer", "confirm-design")
}

func assertInterruptCheckpoint(t *testing.T, projectDir string) {
	t.Helper()
	state := loadRunState(t, projectDir)
	if got := state.Stages["architect"].Status; got != orchestrator.StageStatusInterrupted {
		t.Fatalf("hung stage status = %q, want INTERRUPTED", got)
	}
	for _, done := range []string{"context-engineer", "analyst"} {
		if !state.IsStageCompleted(done) {
			t.Errorf("stage %q should be COMPLETED before the interrupt", done)
		}
	}
}

func interruptRunMidStage(t *testing.T, binary, projectDir, spec string) {
	t.Helper()
	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock", "--mock-hang-stage", "architect")
	run.Dir = projectDir
	if err := run.Start(); err != nil {
		t.Fatalf("start hung run: %v", err)
	}
	waitForStageStatus(t, projectDir, "architect", orchestrator.StageStatusRunning)
	if err := run.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("send SIGINT: %v", err)
	}
	if err := run.Wait(); err == nil {
		t.Fatal("interrupted run exited 0; a canceled run must report the interruption")
	}
}

func waitForStageStatus(t *testing.T, projectDir, stageID string, want orchestrator.StageStatus) {
	t.Helper()
	path := filepath.Join(projectDir, ".claude", "feature-workspace", "user-auth", orchestrator.RunStateFileName)
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if stageStatusIs(path, stageID, want) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("stage %q never reached %q", stageID, want)
}

func stageStatusIs(path, stageID string, want orchestrator.StageStatus) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var state orchestrator.RunState
	if json.Unmarshal(raw, &state) != nil {
		return false
	}
	return state.Stages[stageID].Status == want
}

func TestRunResumeWithoutStateFails(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock", "--resume")
	run.Dir = projectDir
	output, err := run.CombinedOutput()
	if err == nil {
		t.Fatal("--resume with no existing state succeeded; it must fail")
	}
	if !strings.Contains(string(output), "no run state") {
		t.Errorf("error should explain state is missing, got:\n%s", output)
	}
}
