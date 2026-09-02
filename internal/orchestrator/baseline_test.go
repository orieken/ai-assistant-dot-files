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

// twiceGatedPlan gates both the developer and qa-engineer stages, so a run
// reaches a second barrier after the first is approved. gatedPlan (one
// gate) lives in gate_test.go.
func twiceGatedPlan() orchestrator.Plan {
	plan := gatedPlan()
	plan.Stages[2].Gate = "confirm-ship"
	return plan
}

// haltAtGate runs the gated plan until it stops at confirm-design.
func haltAtGate(t *testing.T, scripts map[string]mock.Script) (*orchestrator.Executor, *orchestrator.StateStore, orchestrator.StageInput) {
	t.Helper()
	executor, _, store, input := newHarness(t, scripts)
	err := executor.Run(context.Background(), twiceGatedPlan(), input)
	if !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v, want ErrWaitingApproval", err)
	}
	return executor, store, input
}

func baselineFor(t *testing.T, store *orchestrator.StateStore, gate string) orchestrator.GateBaseline {
	t.Helper()
	baseline, ok := mustLoad(t, store).Baselines[gate]
	if !ok {
		t.Fatalf("no baseline retained for gate %q", gate)
	}
	return baseline
}

// The point of retaining content rather than only a digest: the artifact
// the human was shown must be recoverable from disk.
func TestGateHaltRetainsWhatTheHumanWasShown(t *testing.T) {
	_, store, input := haltAtGate(t, map[string]mock.Script{"analyst": {ArtifactContent: "# original analysis"}})

	baseline := baselineFor(t, store, "confirm-design")
	if baseline.Reason != orchestrator.BaselinePresented {
		t.Errorf("reason = %q, want %q", baseline.Reason, orchestrator.BaselinePresented)
	}
	retained, ok := baseline.Artifacts["analyst"]
	if !ok {
		t.Fatalf("analyst artifact not retained; baseline holds %v", baseline.Artifacts)
	}
	content, err := os.ReadFile(retained.Path)
	if err != nil {
		t.Fatalf("read retained artifact: %v", err)
	}
	if string(content) != "# original analysis" {
		t.Errorf("retained content = %q, want the artifact as presented", content)
	}
	// The copy carries its own digest, so it is verifiable rather than
	// merely remembered — same property the loop's retained iterations have.
	if retained.SHA256 != mustLoad(t, store).Stages["analyst"].ArtifactSHA256 {
		t.Error("retained copy's digest does not match the artifact it copied")
	}
	if !filepath.IsAbs(retained.Path) {
		t.Errorf("retained path %q is not absolute", retained.Path)
	}
	_ = input
}

// A stage that re-runs must never have its own output reported as a human
// correction.
//
// This is the sharp edge of the whole design. Editing a markdown artifact
// makes its stage STALE (L2.12), so the executor re-runs it and the agent's
// fresh output replaces the human's text. If the baseline still held what
// the human wrote, the next comparison would diff it against the agent's
// second attempt and attribute that to a person — a fabricated signal, and
// exactly the confident wrong number this epic exists not to produce. The
// executor therefore refreshes the baseline whenever it writes a stage's
// output itself.
func TestAReRunIsNotReportedAsAHumanCorrection(t *testing.T) {
	executor, provider, store, input := newHarness(t, map[string]mock.Script{
		"analyst": {ArtifactContent: "# original analysis"},
	})
	provider.SetHook(func(stageID string, invocation int) *mock.Script {
		if stageID == "analyst" && invocation > 1 {
			return &mock.Script{ArtifactContent: "# a different second attempt"}
		}
		return nil
	})
	runUntilHalt(t, executor, input)

	artifact := mustLoad(t, store).Stages["analyst"].ArtifactPath
	if err := os.WriteFile(artifact, []byte("# edited by a human"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}
	// The human resumes without approving: the edit is detected, then the
	// staleness cascade re-runs the stage over the top of it.
	runUntilHalt(t, executor, input)
	assertFileContent(t, artifact, "# a different second attempt",
		"precondition: the stage should have been re-run")

	corrections := mustLoad(t, store).Corrections
	if len(corrections) != 1 {
		t.Fatalf("corrections = %d (%+v), want exactly the human's one edit", len(corrections), corrections)
	}
	assertFileContent(t, corrections[0].DiffPath,
		"-# original analysis\n+# edited by a human\n",
		"the diff must describe the human's edit, not the re-run")

	// Approving now must find nothing further: the agent's second attempt
	// is not a correction.
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if after := mustLoad(t, store).Corrections; len(after) != 1 {
		t.Errorf("corrections = %d (%+v) after approval — the agent's own re-run was attributed to a human",
			len(after), after)
	}
}

// Approving accepts the state as it stands, so that state becomes the
// baseline a LATER edit is measured against.
func TestApprovalRefreshesTheBaseline(t *testing.T) {
	executor, store, _ := haltAtGate(t, map[string]mock.Script{"analyst": {ArtifactContent: "# original analysis"}})
	artifact := mustLoad(t, store).Stages["analyst"].ArtifactPath
	if err := os.WriteFile(artifact, []byte("# edited by a human"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}

	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve: %v", err)
	}

	baseline := baselineFor(t, store, "confirm-design")
	if baseline.Reason != orchestrator.BaselineApproved {
		t.Errorf("reason = %q, want %q", baseline.Reason, orchestrator.BaselineApproved)
	}
	content, err := os.ReadFile(baseline.Artifacts["analyst"].Path)
	if err != nil {
		t.Fatalf("read retained artifact: %v", err)
	}
	if string(content) != "# edited by a human" {
		t.Errorf("retained content = %q, want the state the human accepted", content)
	}
}

// Two gates binding the same stage must not overwrite each other.
func TestEachGateRetainsItsOwnCopy(t *testing.T) {
	executor, store, input := haltAtGate(t, map[string]mock.Script{
		"analyst":   {ArtifactContent: "# analysis"},
		"developer": {ArtifactContent: "# implementation"},
	})
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodFlag); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if err := executor.Run(context.Background(), twiceGatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v, want a halt at the second gate", err)
	}

	design := baselineFor(t, store, "confirm-design").Artifacts["analyst"]
	ship := baselineFor(t, store, "confirm-ship").Artifacts["analyst"]
	if design.Path == ship.Path {
		t.Errorf("both gates retained the analyst artifact at %q — one would overwrite the other", design.Path)
	}
}

// A re-run producing byte-identical output changes no digest, so the
// baseline it recorded still describes the truth.
func TestIdenticalReRunLeavesTheBaselineDigestUnchanged(t *testing.T) {
	executor, store, input := haltAtGate(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})
	before := baselineFor(t, store, "confirm-design").Artifacts["analyst"].SHA256

	if err := executor.Run(context.Background(), twiceGatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("second Run error = %v, want ErrWaitingApproval", err)
	}
	if after := baselineFor(t, store, "confirm-design").Artifacts["analyst"].SHA256; after != before {
		t.Errorf("baseline digest changed from %q to %q on an identical re-run", before, after)
	}
}

// A run with no gates retains nothing: there is no human to have been shown
// anything.
func TestUngatedRunRetainsNothing(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{
		"analyst":     {ArtifactContent: "# analysis"},
		"developer":   {ArtifactContent: "# implementation"},
		"qa-engineer": {ArtifactContent: "# qa report"},
	})
	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if baselines := mustLoad(t, store).Baselines; len(baselines) != 0 {
		t.Errorf("ungated run retained %v", baselines)
	}
	if _, err := os.Stat(filepath.Join(input.WorkspaceDir, orchestrator.ApprovedDirName)); !os.IsNotExist(err) {
		t.Errorf("ungated run created %s", orchestrator.ApprovedDirName)
	}
}

// Retention is an observation, not a control: a capture failure is reported
// and never fails the halt or the approval.
func TestBaselineFailureIsReportedButDoesNotFailTheRun(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})
	var reported []error
	executor.OnBaselineError(func(err error) { reported = append(reported, err) })

	// Make the retention target unwritable by occupying its directory path
	// with a file.
	blocker := filepath.Join(input.WorkspaceDir, orchestrator.ApprovedDirName)
	if err := os.WriteFile(blocker, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("create blocker: %v", err)
	}

	err := executor.Run(context.Background(), twiceGatedPlan(), input)
	if !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v — a retention failure must not change how the run ends", err)
	}
	if len(reported) == 0 {
		t.Error("retention failed silently; it must be reported")
	}
	if mustLoad(t, store).Stages["developer"].Status != orchestrator.StageStatusWaitingApproval {
		t.Error("the gate halt itself was disturbed by a retention failure")
	}
}

func runUntilHalt(t *testing.T, executor *orchestrator.Executor, input orchestrator.StageInput) {
	t.Helper()
	if err := executor.Run(context.Background(), twiceGatedPlan(), input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v, want ErrWaitingApproval", err)
	}
}

func assertFileContent(t *testing.T, path, want, why string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(content) != want {
		t.Errorf("%s: content = %q, want %q", why, content, want)
	}
}
