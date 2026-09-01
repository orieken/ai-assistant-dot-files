package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/spf13/cobra"
)

// stateProject sets up a project directory as the working directory, with a
// spec, so the state subcommands resolve the same workspace `loom run` does.
func stateProject(t *testing.T) string {
	t.Helper()
	projectDir := t.TempDir()
	t.Chdir(projectDir)
	return writeSpec(t, projectDir)
}

// callState invokes one subcommand in-process with scripted flags and
// captures everything it printed.
func callState(t *testing.T, run func(*cobra.Command, []string) error, args stateFlags) (string, error) {
	t.Helper()
	stateArgs = args
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	err := run(command, nil)
	return output.String(), err
}

func writeArtifact(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	return path
}

func TestRunStateRecordHashesTheArtifactItself(t *testing.T) {
	spec := stateProject(t)
	artifact := writeArtifact(t, "analysis.md", "# analysis\n")

	output, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: artifact})
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if !strings.Contains(output, "recorded analyst #1") || !strings.Contains(output, "sha256:") {
		t.Errorf("record output = %q", output)
	}

	assertRecordMatchesFile(t, spec, artifact)
}

func assertRecordMatchesFile(t *testing.T, spec, artifact string) {
	t.Helper()
	want, err := orchestrator.ArtifactSHA256(artifact)
	if err != nil {
		t.Fatalf("hash artifact: %v", err)
	}
	record := loadStateFor(t, spec).Stages["analyst"]
	if record.ArtifactSHA256 != want {
		t.Errorf("recorded digest = %q, want the file's actual hash", record.ArtifactSHA256)
	}
	if record.Status != orchestrator.StageStatusCompleted || record.FinishedAt == nil {
		t.Errorf("record = %+v, want a finished COMPLETED stage", record)
	}
}

func loadStateFor(t *testing.T, spec string) *orchestrator.RunState {
	t.Helper()
	workspace := filepath.Join(".claude", "feature-workspace", orchestrator.FeatureNameFromSpec(spec))
	state, err := orchestrator.NewStateStore(filepath.Join(workspace, orchestrator.RunStateFileName)).Load()
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if state == nil {
		t.Fatal("no state recorded")
	}
	return state
}

func TestRunStateRecordAcceptsAStageWithNoArtifact(t *testing.T) {
	spec := stateProject(t)

	output, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "architect"})
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if !strings.Contains(output, "no artifact") {
		t.Errorf("record output should say no artifact was hashed: %q", output)
	}
	if !loadStateFor(t, spec).IsStageCompleted("architect") {
		t.Error("artifact-less stage was not recorded as completed")
	}
}

func TestRunStateRecordRejectsAMissingArtifact(t *testing.T) {
	spec := stateProject(t)

	_, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: "nope.md"})
	if err == nil {
		t.Fatal("record accepted an artifact that does not exist")
	}
}

func TestRunStateVerifyReportsPerStageAndFailsOnEdit(t *testing.T) {
	spec := stateProject(t)
	artifact := writeArtifact(t, "analysis.md", "# analysis\n")
	if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: artifact}); err != nil {
		t.Fatalf("record: %v", err)
	}

	output, err := callState(t, runStateVerify, stateFlags{spec: spec})
	if err != nil || !strings.Contains(output, "OK") {
		t.Fatalf("clean verify: err %v, output %q", err, output)
	}

	if err := os.WriteFile(artifact, []byte("# analysis, edited\n"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}
	assertVerifyFailsOnEdit(t, spec)
}

func assertVerifyFailsOnEdit(t *testing.T, spec string) {
	t.Helper()
	output, err := callState(t, runStateVerify, stateFlags{spec: spec})
	if !errors.Is(err, errStateVerifyFailed) {
		t.Fatalf("verify error = %v, want errStateVerifyFailed", err)
	}
	if !strings.Contains(output, string(orchestrator.StaleReasonEdited)) {
		t.Errorf("verify output = %q, want the stale reason", output)
	}
	if loadStateFor(t, spec).IsStageCompleted("analyst") {
		t.Error("verify left an edited stage marked completed")
	}
}

func TestRunStateApproveRecordsAndSaysItDoesNotEnforce(t *testing.T) {
	spec := stateProject(t)

	output, err := callState(t, runStateApprove, stateFlags{spec: spec, gate: "confirm-ship"})
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if !strings.Contains(output, "audit record") || !strings.Contains(output, "loom run") {
		t.Errorf("approve output must not overclaim enforcement: %q", output)
	}
	if !loadStateFor(t, spec).IsGateApproved("confirm-ship") {
		t.Error("approval was not recorded")
	}
}

func TestRunStateShowRendersSequenceAndApprovals(t *testing.T) {
	spec := stateProject(t)
	recordStages(t, spec, "analyst", "developer")

	output, err := callState(t, runStateShow, stateFlags{spec: spec})
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	for _, want := range []string{"user-auth", "markdown", "1. analyst", "2. developer", "approvals: none recorded"} {
		if !strings.Contains(output, want) {
			t.Errorf("show output missing %q:\n%s", want, output)
		}
	}

	if _, err := callState(t, runStateApprove, stateFlags{spec: spec, gate: "confirm-design"}); err != nil {
		t.Fatalf("approve: %v", err)
	}
	output, err = callState(t, runStateShow, stateFlags{spec: spec})
	if err != nil || !strings.Contains(output, "confirm-design") {
		t.Errorf("show after approval: err %v, output %q", err, output)
	}
}

// recordStages records one completed stage per name, each with its own
// artifact, in order.
func recordStages(t *testing.T, spec string, stageIDs ...string) {
	t.Helper()
	for _, stageID := range stageIDs {
		artifact := writeArtifact(t, stageID+".md", "# "+stageID+"\n")
		if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: stageID, artifact: artifact}); err != nil {
			t.Fatalf("record %s: %v", stageID, err)
		}
	}
}

func TestRunStateShowJSONIsParseable(t *testing.T) {
	spec := stateProject(t)
	artifact := writeArtifact(t, "analysis.md", "# analysis\n")
	if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: artifact}); err != nil {
		t.Fatalf("record: %v", err)
	}

	output, err := callState(t, runStateShow, stateFlags{spec: spec, asJSON: true})
	if err != nil {
		t.Fatalf("show --json: %v", err)
	}
	if !strings.Contains(output, `"createdBy": "markdown"`) || !strings.Contains(output, `"sequence": 1`) {
		t.Errorf("--json output = %s", output)
	}
}

func TestRunStateReadCommandsRequireCheckpoints(t *testing.T) {
	spec := stateProject(t)

	for name, run := range map[string]func(*cobra.Command, []string) error{"verify": runStateVerify, "show": runStateShow} {
		if _, err := callState(t, run, stateFlags{spec: spec}); err == nil {
			t.Errorf("%s succeeded with no checkpoints recorded", name)
		}
	}
}

func TestRunStateRefusesExecutorOwnedState(t *testing.T) {
	spec := stateProject(t)
	workspace := filepath.Join(".claude", "feature-workspace", "user-auth")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	store := orchestrator.NewStateStore(filepath.Join(workspace, orchestrator.RunStateFileName))
	if err := store.Save(orchestrator.NewRunState("deliver-feature", orchestrator.CreatedByExecutor)); err != nil {
		t.Fatalf("seed executor state: %v", err)
	}

	_, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst"})
	if err == nil || !strings.Contains(err.Error(), "executor") {
		t.Errorf("loom state accepted executor-owned state: %v", err)
	}
}

func TestRunStateWritesTheSameTimelineTheExecutorDoes(t *testing.T) {
	spec := stateProject(t)
	recordStages(t, spec, "analyst", "developer")
	if _, err := callState(t, runStateApprove, stateFlags{spec: spec, gate: "confirm-design"}); err != nil {
		t.Fatalf("approve: %v", err)
	}

	output, err := callState(t, runStateTimeline, stateFlags{spec: spec})
	if err != nil {
		t.Fatalf("timeline: %v", err)
	}
	for _, want := range []string{"stage.completed", "analyst", "developer", "gate.approved", "confirm-design"} {
		if !strings.Contains(output, want) {
			t.Errorf("timeline output missing %q:\n%s", want, output)
		}
	}
}

func TestRunStateTimelineRecordsVerificationFailures(t *testing.T) {
	spec := stateProject(t)
	artifact := writeArtifact(t, "analysis.md", "# analysis\n")
	if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: artifact}); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := os.WriteFile(artifact, []byte("# analysis, edited\n"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}
	if _, err := callState(t, runStateVerify, stateFlags{spec: spec}); !errors.Is(err, errStateVerifyFailed) {
		t.Fatalf("verify error = %v", err)
	}

	last := lastTimelineEvent(t, spec)
	if last.Kind != orchestrator.EventStageStale || last.StaleReason != orchestrator.StaleReasonEdited {
		t.Errorf("last event = %+v, want a stale analyst", last)
	}
}

func lastTimelineEvent(t *testing.T, spec string) orchestrator.Event {
	t.Helper()
	output, err := callState(t, runStateTimeline, stateFlags{spec: spec, asJSON: true})
	if err != nil {
		t.Fatalf("timeline --json: %v", err)
	}
	var events []orchestrator.Event
	if err := json.Unmarshal([]byte(output), &events); err != nil {
		t.Fatalf("--json output does not parse: %v\n%s", err, output)
	}
	return events[len(events)-1]
}

func TestRunStateTimelineWithoutEventsFails(t *testing.T) {
	spec := stateProject(t)

	if _, err := callState(t, runStateTimeline, stateFlags{spec: spec}); err == nil {
		t.Error("timeline succeeded for a run with no events")
	}
}

// TestRunStateVerifyReportsAnInvalidatedApproval is the markdown pipeline's
// half of L2.14: detection, not enforcement. `loom state` can say the
// approval no longer holds; it cannot stop a prose pipeline from running.
func TestRunStateVerifyReportsAnInvalidatedApproval(t *testing.T) {
	spec := stateProject(t)
	artifact := writeArtifact(t, "analysis.md", "# analysis\n")
	if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: artifact}); err != nil {
		t.Fatalf("record: %v", err)
	}
	if _, err := callState(t, runStateApprove, stateFlags{spec: spec, gate: "confirm-design"}); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if err := os.WriteFile(artifact, []byte("# analysis, edited after approval\n"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}

	output, err := callState(t, runStateVerify, stateFlags{spec: spec})
	if !errors.Is(err, errStateVerifyFailed) {
		t.Fatalf("verify error = %v, want failure", err)
	}
	assertOutputContains(t, output, "confirm-design", "INVALIDATED", "cannot prevent")
	assertInvalidatedByAnalyst(t, spec)
}

func assertInvalidatedByAnalyst(t *testing.T, spec string) {
	t.Helper()
	approval := loadStateFor(t, spec).Approvals["confirm-design"]
	if approval.IsValid() || approval.InvalidatedBy != "analyst" {
		t.Errorf("approval record = %+v, want invalidated by analyst", approval)
	}
}

func TestRunStateApprovalSurvivesAnUnrelatedEdit(t *testing.T) {
	spec := stateProject(t)
	recordStages(t, spec, "analyst")
	if _, err := callState(t, runStateApprove, stateFlags{spec: spec, gate: "confirm-design"}); err != nil {
		t.Fatalf("approve: %v", err)
	}
	// Recorded after the approval, so it is not part of what was approved.
	later := writeArtifact(t, "implementation.md", "# implementation\n")
	if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "developer", artifact: later}); err != nil {
		t.Fatalf("record developer: %v", err)
	}
	if err := os.WriteFile(later, []byte("# implementation, edited\n"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}

	if _, err := callState(t, runStateVerify, stateFlags{spec: spec}); !errors.Is(err, errStateVerifyFailed) {
		t.Fatalf("verify error = %v, want the stage itself to fail verification", err)
	}
	if !loadStateFor(t, spec).IsGateApproved("confirm-design") {
		t.Error("an edit to work recorded after the approval reset the gate")
	}
}

func TestRunStateShowMarksAnInvalidatedApproval(t *testing.T) {
	spec := stateProject(t)
	artifact := writeArtifact(t, "analysis.md", "# analysis\n")
	if _, err := callState(t, runStateRecord, stateFlags{spec: spec, stage: "analyst", artifact: artifact}); err != nil {
		t.Fatalf("record: %v", err)
	}
	if _, err := callState(t, runStateApprove, stateFlags{spec: spec, gate: "confirm-design"}); err != nil {
		t.Fatalf("approve: %v", err)
	}
	if err := os.WriteFile(artifact, []byte("# analysis, edited\n"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}
	if _, err := callState(t, runStateVerify, stateFlags{spec: spec}); !errors.Is(err, errStateVerifyFailed) {
		t.Fatalf("verify should have failed: %v", err)
	}

	output, err := callState(t, runStateShow, stateFlags{spec: spec})
	if err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(output, "INVALIDATED") {
		t.Errorf("show did not mark the invalidated approval:\n%s", output)
	}
}

func TestRunStateShowListsWhatWasRoutedOut(t *testing.T) {
	spec := stateProject(t)
	recordStages(t, spec, "analyst")
	state := loadStateFor(t, spec)
	state.Stages["devops-engineer"] = orchestrator.StageRecord{
		Status: orchestrator.StageStatusSkipped, Sequence: 2,
		SkipReason: "the analysis lists no DevOps tasks",
	}
	workspace := filepath.Join(".claude", "feature-workspace", "user-auth")
	if err := orchestrator.NewStateStore(filepath.Join(workspace, orchestrator.RunStateFileName)).Save(state); err != nil {
		t.Fatalf("save state: %v", err)
	}

	output, err := callState(t, runStateShow, stateFlags{spec: spec})
	if err != nil {
		t.Fatalf("show: %v", err)
	}

	for _, want := range []string{"routed out:", "devops-engineer", "no DevOps tasks"} {
		if !strings.Contains(output, want) {
			t.Errorf("show output missing %q:\n%s", want, output)
		}
	}
}
