package orchestrator_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/mock"
	"github.com/orieken/loom/internal/state"
)

// typedGatedPlan types the analyst stage, so its artifact is a state
// document and `analysis.md` is a rendered view.
func typedGatedPlan() orchestrator.Plan {
	return orchestrator.Plan{
		Name: "typed-gated-plan",
		Stages: []orchestrator.Stage{
			{ID: "analyst", Agent: "analyst", StateKind: string(state.KindAnalysis), Timeout: 5 * time.Second},
			{ID: "developer", Agent: "developer", Gate: "confirm-design", Timeout: 5 * time.Second},
		},
	}
}

func typedHarness(t *testing.T) (*orchestrator.Executor, *orchestrator.StateStore, orchestrator.StageInput) {
	t.Helper()
	analysis, ok := mock.TypedScript(string(state.KindAnalysis))
	if !ok {
		t.Fatal("mock has no scripted analysis payload")
	}
	executor, _, store, input := newHarness(t, map[string]mock.Script{"analyst": analysis})
	err := executor.Run(context.Background(), typedGatedPlan(), input)
	if !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v, want ErrWaitingApproval", err)
	}
	return executor, store, input
}

// The design's central claim: a human corrects the RENDERED VIEW, that
// correction is recorded and attributed, and it disturbs nothing — because
// the executor never reads the view back.
func TestEditingARenderedViewIsRecordedAndDisturbsNothing(t *testing.T) {
	executor, store, input := typedHarness(t)
	view := filepath.Join(input.WorkspaceDir, "analysis.md")
	original, err := os.ReadFile(view)
	if err != nil {
		t.Fatalf("read view: %v", err)
	}
	corrected := string(original) + "\n## Added by a human\nThe analyst missed the rate-limit case.\n"
	if err := os.WriteFile(view, []byte(corrected), 0o644); err != nil {
		t.Fatalf("edit view: %v", err)
	}

	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve: %v", err)
	}

	recorded := mustLoad(t, store)
	correction := onlyCorrection(t, recorded)
	assertAttributedToAnalyst(t, correction)
	assertDiffContains(t, correction.DiffPath, "+The analyst missed the rate-limit case.")
	assertRunUndisturbed(t, recorded)
}

func onlyCorrection(t *testing.T, recorded *orchestrator.RunState) orchestrator.Correction {
	t.Helper()
	if len(recorded.Corrections) != 1 {
		t.Fatalf("corrections = %+v, want exactly one", recorded.Corrections)
	}
	return recorded.Corrections[0]
}

// Attribution is to whoever produced the artifact, not to the gate that
// caught it: whose output needed fixing is the part worth learning from.
func assertAttributedToAnalyst(t *testing.T, correction orchestrator.Correction) {
	t.Helper()
	if correction.Stage != "analyst" || correction.Agent != "analyst" {
		t.Errorf("correction attributed to stage %q / agent %q, want the producing analyst",
			correction.Stage, correction.Agent)
	}
	if correction.Stat.Added == 0 {
		t.Errorf("stat = %v, want added lines", correction.Stat)
	}
}

func assertDiffContains(t *testing.T, path, want string) {
	t.Helper()
	diff, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read retained diff: %v", err)
	}
	if !strings.Contains(string(diff), want) {
		t.Errorf("diff does not carry %q:\n%s", want, diff)
	}
}

// The analyst's tracked state document is untouched, so nothing went stale
// and nothing re-runs — which is the whole reason the view is the right
// place for a human to write a correction.
func assertRunUndisturbed(t *testing.T, recorded *orchestrator.RunState) {
	t.Helper()
	if got := recorded.Stages["analyst"].Status; got != orchestrator.StageStatusCompleted {
		t.Errorf("analyst status = %q — editing a view must not disturb the run", got)
	}
	if !recorded.IsGateApproved("confirm-design") {
		t.Error("approval did not stand; a view edit must not invalidate it")
	}
}

// An approval with nothing changed is not a correction.
func TestUneditedApprovalRecordsNoCorrection(t *testing.T) {
	executor, store, _ := typedHarness(t)

	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if corrections := mustLoad(t, store).Corrections; len(corrections) != 0 {
		t.Errorf("corrections = %+v, want none for an untouched artifact", corrections)
	}
}

// The correction reaches the timeline as one greppable line, with the diff
// itself on disk rather than inline.
func TestCorrectionReachesTheTimeline(t *testing.T) {
	executor, store, input := typedHarness(t)
	view := filepath.Join(input.WorkspaceDir, "analysis.md")
	if err := os.WriteFile(view, []byte("# replaced wholesale\n"), 0o644); err != nil {
		t.Fatalf("edit view: %v", err)
	}
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve: %v", err)
	}

	found := correctionEventIn(t, readTimeline(t, store))
	if found.Stage != "analyst" || found.Agent != "analyst" {
		t.Errorf("event names stage %q agent %q, want the analyst", found.Stage, found.Agent)
	}
	if !strings.HasPrefix(found.Correction, "+") {
		t.Errorf("event summary = %q, want a +added/-removed stat", found.Correction)
	}
}

func correctionEventIn(t *testing.T, events []orchestrator.Event) orchestrator.Event {
	t.Helper()
	for _, event := range events {
		if event.Kind == orchestrator.EventArtifactCorrected {
			return event
		}
	}
	t.Fatalf("no %q event on the timeline", orchestrator.EventArtifactCorrected)
	return orchestrator.Event{}
}

// Reporting a correction consumes it: approving twice must not report the
// same edit again.
func TestACorrectionIsReportedOnlyOnce(t *testing.T) {
	executor, store, input := typedHarness(t)
	view := filepath.Join(input.WorkspaceDir, "analysis.md")
	if err := os.WriteFile(view, []byte("# corrected once\n"), 0o644); err != nil {
		t.Fatalf("edit view: %v", err)
	}
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("first Approve: %v", err)
	}
	if err := executor.Run(context.Background(), typedGatedPlan(), input); err != nil {
		t.Fatalf("Run after approval: %v", err)
	}

	if corrections := mustLoad(t, store).Corrections; len(corrections) != 1 {
		t.Errorf("corrections = %d (%+v), want the one edit reported once", len(corrections), corrections)
	}
}

// A diff that cannot be written must not fail the approval — this observes
// a human's action, it does not control the run.
func TestUnwritableDiffLeavesTheApprovalIntact(t *testing.T) {
	executor, store, input := typedHarness(t)
	var reported []error
	executor.OnBaselineError(func(err error) { reported = append(reported, err) })
	view := filepath.Join(input.WorkspaceDir, "analysis.md")
	if err := os.WriteFile(view, []byte("# corrected\n"), 0o644); err != nil {
		t.Fatalf("edit view: %v", err)
	}
	// Occupy the corrections directory path with a file.
	blocker := filepath.Join(input.WorkspaceDir, orchestrator.ApprovedDirName, "confirm-design",
		orchestrator.CorrectionsDirName)
	if err := os.WriteFile(blocker, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("create blocker: %v", err)
	}

	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve failed because a diff could not be written: %v", err)
	}
	recorded := mustLoad(t, store)
	if !recorded.IsGateApproved("confirm-design") {
		t.Error("approval did not stand")
	}
	// The correction is still recorded — losing the diff loses the detail,
	// not the fact that a human corrected this agent.
	if correction := onlyCorrection(t, recorded); correction.DiffPath != "" {
		t.Errorf("correction kept a diff path %q that could not be written", correction.DiffPath)
	}
	if len(reported) == 0 {
		t.Error("the unwritable diff was not reported")
	}
}
