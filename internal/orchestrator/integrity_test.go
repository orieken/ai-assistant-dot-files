package orchestrator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/mock"
)

func completedScripts() map[string]mock.Script {
	return map[string]mock.Script{
		"analyst":     {ArtifactContent: "# analysis"},
		"developer":   {ArtifactContent: "# implementation"},
		"qa-engineer": {ArtifactContent: "# qa report"},
	}
}

// completedRun runs the three-stage plan to completion and returns the
// harness so a test can tamper with the workspace and resume.
func completedRun(t *testing.T) (*orchestrator.Executor, *mock.Provider, *orchestrator.StateStore, orchestrator.StageInput) {
	t.Helper()
	executor, provider, store, input := newHarness(t, completedScripts())
	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("first run: %v", err)
	}
	return executor, provider, store, input
}

func editArtifact(t *testing.T, input orchestrator.StageInput, stageID, content string) {
	t.Helper()
	path := filepath.Join(input.WorkspaceDir, stageID+".md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("edit %s artifact: %v", stageID, err)
	}
}

func resume(t *testing.T, executor *orchestrator.Executor, input orchestrator.StageInput) []orchestrator.StaleStage {
	t.Helper()
	var reported []orchestrator.StaleStage
	executor.OnStale(func(stale []orchestrator.StaleStage) { reported = stale })
	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("resume run: %v", err)
	}
	return reported
}

// TestHandEditedArtifactIsNotTreatedAsComplete is the L2.12 done-when test:
// a stage whose artifact changed on disk re-runs, decided by a digest
// computed in Go rather than by a model reporting its own hash.
func TestHandEditedArtifactIsNotTreatedAsComplete(t *testing.T) {
	executor, provider, store, input := completedRun(t)
	editArtifact(t, input, "qa-engineer", "# qa report, hand-edited")

	stale := resume(t, executor, input)

	if len(stale) != 1 || stale[0].StageID != "qa-engineer" || stale[0].Reason != orchestrator.StaleReasonEdited {
		t.Fatalf("stale report = %+v, want qa-engineer/ARTIFACT_EDITED", stale)
	}
	if stale[0].Description() == "" {
		t.Error("stale stage has no human-readable description")
	}
	assertStatus(t, mustLoad(t, store), "qa-engineer", orchestrator.StageStatusCompleted)
	assertInvocations(t, provider, []string{"analyst", "developer", "qa-engineer", "qa-engineer"})
}

func TestEditingAnEarlyArtifactCascadesToLaterStages(t *testing.T) {
	executor, provider, store, input := completedRun(t)
	editArtifact(t, input, "analyst", "# analysis, rewritten by hand")

	stale := resume(t, executor, input)

	wantReasons := map[string]orchestrator.StaleReason{
		"analyst":     orchestrator.StaleReasonEdited,
		"developer":   orchestrator.StaleReasonUpstream,
		"qa-engineer": orchestrator.StaleReasonUpstream,
	}
	if len(stale) != len(wantReasons) {
		t.Fatalf("stale report = %+v, want all three stages", stale)
	}
	for _, got := range stale {
		if got.Reason != wantReasons[got.StageID] {
			t.Errorf("stage %q stale reason = %q, want %q", got.StageID, got.Reason, wantReasons[got.StageID])
		}
	}
	// Every stage re-runs: the two downstream artifacts were derived from
	// content that no longer exists.
	assertInvocations(t, provider, []string{
		"analyst", "developer", "qa-engineer",
		"analyst", "developer", "qa-engineer",
	})
	state := mustLoad(t, store)
	for stageID := range wantReasons {
		assertStatus(t, state, stageID, orchestrator.StageStatusCompleted)
	}
}

func TestDeletedArtifactCountsAsAMismatch(t *testing.T) {
	executor, _, _, input := completedRun(t)
	if err := os.Remove(filepath.Join(input.WorkspaceDir, "developer.md")); err != nil {
		t.Fatalf("remove artifact: %v", err)
	}

	stale := resume(t, executor, input)

	if len(stale) != 2 || stale[0].StageID != "developer" || stale[0].Reason != orchestrator.StaleReasonMissing {
		t.Fatalf("stale report = %+v, want developer/ARTIFACT_MISSING then the cascade", stale)
	}
}

func TestUnchangedArtifactsDemoteNothing(t *testing.T) {
	executor, provider, _, input := completedRun(t)

	stale := resume(t, executor, input)

	if len(stale) != 0 {
		t.Fatalf("clean resume demoted %+v; verification must not produce false positives", stale)
	}
	assertInvocations(t, provider, []string{"analyst", "developer", "qa-engineer"})
}

func TestStageWithoutAnArtifactStaysCompleted(t *testing.T) {
	executor, provider, store, input := newHarness(t, map[string]mock.Script{
		"analyst":     {ArtifactContent: "# analysis"},
		"developer":   {},
		"qa-engineer": {ArtifactContent: "# qa report"},
	})
	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("first run: %v", err)
	}

	if stale := resume(t, executor, input); len(stale) != 0 {
		t.Fatalf("artifact-less stage was demoted: %+v", stale)
	}
	assertStatus(t, mustLoad(t, store), "developer", orchestrator.StageStatusCompleted)
	assertInvocations(t, provider, []string{"analyst", "developer", "qa-engineer"})
}

// TestStaleStageDoesNotRevokeApproval pins the L2.14 boundary: editing an
// artifact invalidates the stage, but binding approvals to digests is the
// next epic's job. Asserted so that work starts from a known state.
func TestStaleStageDoesNotRevokeApproval(t *testing.T) {
	executor, _, store, input := newHarness(t, completedScripts())
	plan := gatedPlan()
	_ = runUntilGate(t, executor, plan, input)
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodTTY); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if err := executor.Run(context.Background(), plan, input); err != nil {
		t.Fatalf("run after approval: %v", err)
	}

	editArtifact(t, input, "analyst", "# analysis, hand-edited after approval")
	executor.OnStale(func([]orchestrator.StaleStage) {})
	if err := executor.Run(context.Background(), plan, input); err != nil {
		t.Fatalf("resume after edit: %v", err)
	}

	if !mustLoad(t, store).IsGateApproved("confirm-design") {
		t.Error("editing an artifact revoked an approval; reset-on-edit is L2.14, not this epic")
	}
}

func TestRunRefusesStateFromADifferentSpec(t *testing.T) {
	executor, _, _, input := completedRun(t)
	other := input
	other.SpecPath = filepath.Join(input.WorkspaceDir, "some-other-feature.md")

	if err := executor.Run(context.Background(), threeStagePlan(), other); err == nil {
		t.Fatal("Run accepted state belonging to a different spec")
	}
}

func TestFreshRunRecordsItsIdentity(t *testing.T) {
	_, _, store, input := completedRun(t)
	state := mustLoad(t, store)

	if state.SpecPath != input.SpecPath {
		t.Errorf("state spec path = %q, want %q", state.SpecPath, input.SpecPath)
	}
	if state.FeatureName != "spec" {
		t.Errorf("feature name = %q, want the spec file's base name", state.FeatureName)
	}
	if state.StartedAt.IsZero() {
		t.Error("fresh run did not record when it started")
	}
}

func TestLoadRefusesSchemaVersionTwoStateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), orchestrator.RunStateFileName)
	if err := os.WriteFile(path, []byte(`{"schemaVersion": 2, "planName": "test-plan", "stages": {}}`), 0o644); err != nil {
		t.Fatalf("seed v2 state: %v", err)
	}
	if _, err := orchestrator.NewStateStore(path).Load(); err == nil {
		t.Fatal("Load accepted a v2 run state; schema 3 adds stale tracking and must refuse it")
	}
}
