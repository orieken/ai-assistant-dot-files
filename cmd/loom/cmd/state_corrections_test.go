package cmd

// In-process rendering tests. The binary-driven tests elsewhere in this
// package are invisible to the coverage tool, so the reporting logic is
// exercised directly here.

import (
	"bytes"
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
