package cmd

// In-process rendering tests. The binary-driven tests elsewhere in this
// package are invisible to the coverage tool, so the reporting logic is
// exercised directly here.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/spf13/cobra"
)

func captured(t *testing.T, render func(*cobra.Command)) string {
	t.Helper()
	output := &bytes.Buffer{}
	command := &cobra.Command{}
	command.SetOut(output)
	command.SetErr(output)
	render(command)
	return output.String()
}

func stateWithCorrection(correction orchestrator.Correction) *orchestrator.RunState {
	state := orchestrator.NewRunState("test-plan", orchestrator.CreatedByExecutor)
	state.FeatureName = "user-auth"
	state.Stages["analyst"] = orchestrator.StageRecord{
		Status: orchestrator.StageStatusCompleted, Sequence: 1, Agent: "analyst",
	}
	state.Corrections = []orchestrator.Correction{correction}
	return state
}

// "Whose output did a human have to fix?" must be answerable from
// `loom state show` without opening a JSON file.
func TestStateShowNamesWhoWasCorrected(t *testing.T) {
	state := stateWithCorrection(orchestrator.Correction{
		Stage: "analyst", Agent: "analyst", Gate: "confirm-design", At: time.Now().UTC(),
		DiffPath: "/tmp/ws/.approved/confirm-design/corrections/analyst.diff",
		Stat:     orchestrator.DiffStat{Added: 4, Removed: 1},
	})

	output := captured(t, func(command *cobra.Command) { printStateReport(command, state) })

	for _, want := range []string{"corrected by a human", "analyst", "+4/-1", "confirm-design", "analyst.diff"} {
		if !strings.Contains(output, want) {
			t.Errorf("state show output missing %q:\n%s", want, output)
		}
	}
}

// A run nobody corrected must say nothing, rather than print an empty
// heading that reads like a missing section.
func TestStateShowSaysNothingWhenNothingWasCorrected(t *testing.T) {
	state := stateWithCorrection(orchestrator.Correction{})
	state.Corrections = nil

	output := captured(t, func(command *cobra.Command) { printStateReport(command, state) })

	if strings.Contains(output, "corrected by a human") {
		t.Errorf("clean run printed a corrections heading:\n%s", output)
	}
}

// A correction whose diff could not be written is still reported — losing
// the diff loses the detail, not the fact.
func TestStateShowReportsACorrectionWithNoDiff(t *testing.T) {
	state := stateWithCorrection(orchestrator.Correction{
		Stage: "analyst", Agent: "analyst", Gate: "confirm-design",
		Stat: orchestrator.DiffStat{Added: 1},
	})

	output := captured(t, func(command *cobra.Command) { printStateReport(command, state) })

	if !strings.Contains(output, "analyst") || !strings.Contains(output, "+1/-0") {
		t.Errorf("correction with no diff was not reported:\n%s", output)
	}
	if strings.Contains(output, " — ") {
		t.Errorf("output includes an empty diff separator:\n%s", output)
	}
}

// The timeline line reads as "who was corrected, by how much" — the diff's
// absolute path would swamp it, and lives in --json instead.
func TestTimelineRendersACorrectionWithoutTheDiffPath(t *testing.T) {
	start := time.Now().UTC()
	events := []orchestrator.Event{
		{At: start, Kind: orchestrator.EventRunStarted},
		{
			At: start.Add(time.Second), Kind: orchestrator.EventArtifactCorrected,
			Stage: "analyst", Agent: "analyst", Gate: "confirm-design", Correction: "+4/-1",
			DiffPath: "/tmp/ws/.approved/confirm-design/corrections/analyst.diff",
		},
	}

	output := captured(t, func(command *cobra.Command) { printEventReport(command, events) })

	if !strings.Contains(output, "artifact.corrected") || !strings.Contains(output, "analyst (analyst) +4/-1") {
		t.Errorf("timeline does not render the correction legibly:\n%s", output)
	}
	if strings.Contains(output, "analyst.diff") {
		t.Errorf("timeline line carries the absolute diff path:\n%s", output)
	}
}

// --dry-run-policies evaluates a run that already happened, so it must not
// go through the resume check — that check refuses to start over existing
// state, which is exactly the state dry-run reads. Found end-to-end: the
// flag failed on every run it was meant to inspect.
func TestDryRunPoliciesDoesNotGoThroughTheResumeCheck(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)
	writeExamplePolicy(t, projectDir)

	// Produce a halted run, so state exists.
	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	run.Dir = projectDir
	if output, err := run.CombinedOutput(); err == nil {
		t.Fatalf("run should halt at confirm-design\n%s", output)
	}

	output, code := runLoom(t, binary, projectDir, "run", "--spec", spec, "--dry-run-policies")
	if code != 0 {
		t.Fatalf("dry-run exited %d over an existing run:\n%s", code, output)
	}
	for _, want := range []string{"Dry-run against", "confirm-design"} {
		if !strings.Contains(output, want) {
			t.Errorf("dry-run output missing %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "already exists") {
		t.Errorf("dry-run hit the resume check:\n%s", output)
	}
}

// The kill-switch has to work through the CLI, not only in the package.
func TestPoliciesEnabledFalseSkipsEvaluation(t *testing.T) {
	binary := buildLoomBinary(t)
	projectDir := t.TempDir()
	spec := writeSpec(t, projectDir)
	writeExamplePolicy(t, projectDir)
	if err := os.WriteFile(filepath.Join(projectDir, ".claude", "delivery-policy.yaml"),
		[]byte("policiesEnabled: false\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	run := exec.Command(binary, "run", "--spec", spec, "--provider", "mock")
	run.Dir = projectDir
	output, err := run.CombinedOutput()
	if err == nil {
		t.Fatalf("run should halt at confirm-design\n%s", output)
	}
	if !strings.Contains(string(output), "policy evaluation is off") {
		t.Errorf("output does not report the kill-switch:\n%s", output)
	}
	if strings.Contains(string(output), "Policy at gate") {
		t.Errorf("policies were evaluated despite policiesEnabled: false:\n%s", output)
	}
}

func writeExamplePolicy(t *testing.T, projectDir string) {
	t.Helper()
	dir := filepath.Join(projectDir, ".claude", "policies")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create policy dir: %v", err)
	}
	body := "name: watch-design\nversion: \"1.0\"\nmatcher:\n  gate: confirm-design\n" +
		"condition:\n  testsPass: true\naction:\n  type: require-human\n  reason: \"always ask\"\n"
	if err := os.WriteFile(filepath.Join(dir, "watch.policy.yaml"), []byte(body), 0o644); err != nil {
		t.Fatalf("write policy: %v", err)
	}
}
