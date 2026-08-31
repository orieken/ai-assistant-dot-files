package orchestrator_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/mock"
)

// approvedAtGate runs the gated plan to its gate and approves it, returning
// the harness positioned exactly where a human has said yes.
func approvedAtGate(t *testing.T, scripts map[string]mock.Script) (*orchestrator.Executor, *mock.Provider, *orchestrator.StateStore, orchestrator.StageInput) {
	t.Helper()
	executor, provider, store, input := newHarness(t, scripts)
	_ = runUntilGate(t, executor, gatedPlan(), input)
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodTTY); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	return executor, provider, store, input
}

func resumeGated(t *testing.T, executor *orchestrator.Executor, input orchestrator.StageInput) error {
	t.Helper()
	return executor.Run(context.Background(), gatedPlan(), input)
}

// TestEditingAnArtifactAfterApprovalRefusesExecution is the L2.14
// done-when: an edit between approval and execution stops the gated stage
// from running, and the halt names the artifact that changed.
func TestEditingAnArtifactAfterApprovalRefusesExecution(t *testing.T) {
	executor, provider, store, input := approvedAtGate(t, completedScripts())

	editArtifact(t, input, "analyst", "# analysis, edited after the human approved")
	err := resumeGated(t, executor, input)

	var waiting *orchestrator.WaitingApprovalError
	if !errors.As(err, &waiting) || waiting.Gate != "confirm-design" {
		t.Fatalf("Run error = %v, want a halt at confirm-design", err)
	}
	state := mustLoad(t, store)
	if state.IsGateApproved("confirm-design") {
		t.Error("the approval survived an edit to a bound artifact")
	}
	approval := state.Approvals["confirm-design"]
	if approval.InvalidatedAt == nil || approval.InvalidatedBy != "analyst" {
		t.Errorf("invalidation record = %+v, want it to name the analyst", approval)
	}
	// The gated stage never ran: only the analyst was invoked, twice (once
	// originally, once re-run by the L2.12 cascade).
	assertInvocations(t, provider, []string{"analyst", "analyst"})
}

func TestInvalidationIsRecordedNotDeleted(t *testing.T) {
	executor, _, store, input := approvedAtGate(t, completedScripts())
	editArtifact(t, input, "analyst", "# analysis, edited")
	_ = resumeGated(t, executor, input)

	approval := mustLoad(t, store).Approvals["confirm-design"]
	if approval.ApprovedAt.IsZero() || approval.Approver == "" {
		t.Error("invalidation destroyed the record of who approved and when")
	}
	if len(approval.ArtifactDigests) == 0 {
		t.Error("invalidation destroyed the digests the human approved")
	}
}

// TestIdenticalRerunKeepsTheApproval pins the distinction the rule actually
// makes: any *edit* resets the gate, not any re-run. Forcing re-approval
// for provably identical content trains people to click through gates.
func TestIdenticalRerunKeepsTheApproval(t *testing.T) {
	executor, provider, store, input := approvedAtGate(t, completedScripts())

	// Touch the artifact with byte-identical content, as a re-run would.
	editArtifact(t, input, "analyst", "# analysis")
	if err := resumeGated(t, executor, input); err != nil {
		t.Fatalf("identical content forced re-approval: %v", err)
	}

	if !mustLoad(t, store).IsGateApproved("confirm-design") {
		t.Error("approval was invalidated by a re-run that changed nothing")
	}
	assertInvocations(t, provider, []string{"analyst", "developer", "qa-engineer"})
}

func TestReapprovingAfterInvalidationProceeds(t *testing.T) {
	executor, _, store, input := approvedAtGate(t, completedScripts())
	editArtifact(t, input, "analyst", "# analysis, edited")
	_ = resumeGated(t, executor, input)

	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("re-approve: %v", err)
	}
	if err := resumeGated(t, executor, input); err != nil {
		t.Fatalf("run after re-approval: %v", err)
	}

	approval := mustLoad(t, store).Approvals["confirm-design"]
	if !approval.IsValid() || approval.Method != orchestrator.ApprovalMethodFlag {
		t.Errorf("re-approval did not replace the invalidated record: %+v", approval)
	}
}

func TestDeletingABoundArtifactInvalidatesTheApproval(t *testing.T) {
	executor, _, store, input := approvedAtGate(t, completedScripts())

	if err := os.Remove(filepath.Join(input.WorkspaceDir, "analyst.md")); err != nil {
		t.Fatalf("remove artifact: %v", err)
	}
	if err := resumeGated(t, executor, input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v, want a halt", err)
	}
	if mustLoad(t, store).IsGateApproved("confirm-design") {
		t.Error("approval survived the deletion of a bound artifact")
	}
}

// TestApprovalBindsOnlyWhatWasCompleteWhenApproved keeps the binding
// coherent: a human cannot have approved work that had not happened yet.
func TestApprovalBindsOnlyWhatWasCompleteWhenApproved(t *testing.T) {
	executor, _, store, input := approvedAtGate(t, completedScripts())

	bound := mustLoad(t, store).Approvals["confirm-design"].ArtifactDigests
	if _, ok := bound["analyst"]; !ok {
		t.Error("the completed analyst stage was not bound")
	}
	if _, ok := bound["developer"]; ok {
		t.Error("a stage that had not run when the human approved was bound to the approval")
	}

	// Running on and editing a later artifact must not disturb this gate.
	if err := resumeGated(t, executor, input); err != nil {
		t.Fatalf("run: %v", err)
	}
	editArtifact(t, input, "developer", "# implementation, edited later")
	if err := resumeGated(t, executor, input); err != nil {
		t.Fatalf("editing a stage completed after approval halted the run: %v", err)
	}
	if !mustLoad(t, store).IsGateApproved("confirm-design") {
		t.Error("an edit to a stage completed after approval invalidated it")
	}
}

// TestEditBothDemotesTheStageAndInvalidatesTheApproval makes the
// interaction with L2.12 explicit: they are two mechanisms, and an edit
// trips both.
func TestEditBothDemotesTheStageAndInvalidatesTheApproval(t *testing.T) {
	executor, _, store, input := approvedAtGate(t, completedScripts())
	var reported []orchestrator.StaleStage
	executor.OnStale(func(stale []orchestrator.StaleStage) { reported = stale })

	editArtifact(t, input, "analyst", "# analysis, edited")
	_ = resumeGated(t, executor, input)

	if len(reported) == 0 || reported[0].StageID != "analyst" {
		t.Errorf("L2.12 did not demote the edited stage: %+v", reported)
	}
	if mustLoad(t, store).IsGateApproved("confirm-design") {
		t.Error("L2.14 did not invalidate the approval")
	}
}

func TestInvalidationEmitsATimelineEvent(t *testing.T) {
	executor, _, store, input := approvedAtGate(t, completedScripts())
	editArtifact(t, input, "analyst", "# analysis, edited")
	_ = resumeGated(t, executor, input)

	found := firstOfEachKind(readTimeline(t, store))
	event, ok := found[orchestrator.EventGateInvalidated]
	if !ok {
		t.Fatal("no gate.invalidated event was recorded")
	}
	if event.Gate != "confirm-design" || event.Stage != "analyst" {
		t.Errorf("gate.invalidated event = %+v, want confirm-design / analyst", event)
	}
}

func TestLoadRefusesSchemaVersionFourStateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), orchestrator.RunStateFileName)
	if err := os.WriteFile(path, []byte(`{"schemaVersion": 4, "planName": "p", "stages": {}}`), 0o644); err != nil {
		t.Fatalf("seed v4 state: %v", err)
	}
	if _, err := orchestrator.NewStateStore(path).Load(); err == nil {
		t.Fatal("Load accepted a v4 run state; schema 5 binds approvals to digests and must refuse it")
	}
}
