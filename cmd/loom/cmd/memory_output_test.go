package cmd

// In-process tests for the `loom memory` rendering. The binary-driven tests
// in run_memory_test.go prove the wiring; these are what the coverage tool
// can see, and they are where the careful wording is pinned.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"

	"github.com/orieken/loom/internal/memory"
	"github.com/spf13/cobra"
)

func renderMemory(t *testing.T, render func(*cobra.Command)) string {
	t.Helper()
	return captured(t, render)
}

func sampleRuns() []memory.RunSummary {
	return []memory.RunSummary{
		{Feature: "user-auth", StartedAt: "2026-09-02T10:00:00Z", Complete: true,
			InputTokens: 1200, OutputTokens: 300, CostUSD: 0.0425, Corrections: 2},
		{Feature: "billing", StartedAt: "2026-09-01T09:00:00Z", WaitingGate: "confirm-design"},
	}
}

func TestRunsTableShowsStateAndCost(t *testing.T) {
	output := renderMemory(t, func(command *cobra.Command) { printRuns(command, nil, sampleRuns()) })

	for _, want := range []string{"user-auth", "complete", "$0.0425", "billing", "at gate"} {
		if !strings.Contains(output, want) {
			t.Errorf("runs table missing %q:\n%s", want, output)
		}
	}
}

// A cost of zero renders as an em dash, not $0.0000: a provider that
// reported nothing is a different fact from a run that cost nothing, and
// the table must not assert the second.
func TestUnreportedCostIsNotRenderedAsZeroDollars(t *testing.T) {
	output := renderMemory(t, func(command *cobra.Command) { printRuns(command, nil, sampleRuns()) })

	if strings.Contains(output, "$0.0000") {
		t.Errorf("a run with no reported usage rendered as costing nothing:\n%s", output)
	}
	if !strings.Contains(output, "—") {
		t.Errorf("unreported cost is not marked as unmeasured:\n%s", output)
	}
}

func TestEmptyRunsListSaysSo(t *testing.T) {
	output := renderMemory(t, func(command *cobra.Command) { printRuns(command, nil, nil) })

	if !strings.Contains(output, "No runs recorded yet") {
		t.Errorf("empty run list produced %q", output)
	}
}

func TestRunStateNaming(t *testing.T) {
	cases := []struct {
		name string
		run  memory.RunSummary
		want string
	}{
		{"completed", memory.RunSummary{Complete: true}, "complete"},
		{"halted at a gate", memory.RunSummary{WaitingGate: "confirm-ship"}, "at gate"},
		{"neither", memory.RunSummary{}, "unfinished"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := runState(testCase.run); got != testCase.want {
				t.Errorf("runState = %q, want %q", got, testCase.want)
			}
		})
	}
}

// "More than twice" has two defensible readings, so the header states which
// one produced the rows rather than leaving a reader to infer it.
func TestRetriesHeaderStatesTheThreshold(t *testing.T) {
	memoryArgs.moreThan = 2
	rows := []memory.RetryRow{
		{Feature: "user-auth", StartedAt: "2026-09-02T10:00:00Z", Stage: "code-reviewer", Iterations: 4},
	}

	output := renderMemory(t, func(command *cobra.Command) { printRetries(command, rows) })

	for _, want := range []string{"more than 2 iterations", "3 attempts", "code-reviewer", "4"} {
		if !strings.Contains(output, want) {
			t.Errorf("retries output missing %q:\n%s", want, output)
		}
	}
}

func TestRetriesWithNoRowsStillStatesTheThreshold(t *testing.T) {
	memoryArgs.moreThan = 3

	output := renderMemory(t, func(command *cobra.Command) { printRetries(command, nil) })

	if !strings.Contains(output, "more than 3 iterations") || !strings.Contains(output, "None recorded") {
		t.Errorf("empty retries output is not self-explaining:\n%s", output)
	}
}

// A correction was recorded, not adopted (roadmap L4.5). A table that let a
// reader assume otherwise would misrepresent the whole signal.
func TestCorrectionsTableSaysTheEditWasNotAdopted(t *testing.T) {
	rows := []memory.AgentCorrections{{Agent: "analyst", Corrections: 3, Runs: 2, LinesAdded: 9, LinesCut: 1}}

	output := renderMemory(t, func(command *cobra.Command) { printAgentCorrections(command, rows) })

	if !strings.Contains(output, "analyst") || !strings.Contains(output, "3") {
		t.Errorf("corrections table missing its data:\n%s", output)
	}
	if !strings.Contains(output, "recorded, not adopted") {
		t.Errorf("corrections table does not say the edit was not adopted:\n%s", output)
	}
}

func TestEmptyCorrectionsSaysSo(t *testing.T) {
	output := renderMemory(t, func(command *cobra.Command) { printAgentCorrections(command, nil) })

	if !strings.Contains(output, "No human corrections recorded") {
		t.Errorf("empty corrections produced %q", output)
	}
}

func TestShortTimeTrimsToMinutes(t *testing.T) {
	if got := shortTime("2026-09-02T10:00:00.123456Z"); got != "2026-09-02T10:00" {
		t.Errorf("shortTime = %q", got)
	}
	// A value too short to trim is returned unchanged rather than panicking.
	if got := shortTime("2026"); got != "2026" {
		t.Errorf("shortTime on a short value = %q", got)
	}
}

func TestCompleteSuffixMarksHaltedRuns(t *testing.T) {
	if got := completeSuffix(memory.Run{Complete: true}); got != "" {
		t.Errorf("completed run suffixed with %q", got)
	}
	if got := completeSuffix(memory.Run{}); !strings.Contains(got, "halted") {
		t.Errorf("halted run suffix = %q", got)
	}
}

// projectWithStore builds a temp project holding one ingested run and makes
// it the working directory, so the command bodies can be exercised
// in-process rather than only through the built binary.
func projectWithStore(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	archive := filepath.Join("docs", "features", "user-auth")
	if err := os.MkdirAll(archive, 0o755); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	started := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	state := orchestrator.NewRunState("deliver-feature", orchestrator.CreatedByExecutor)
	state.FeatureName = "user-auth"
	state.StartedAt = started
	state.Stages["code-reviewer"] = orchestrator.StageRecord{
		Status: orchestrator.StageStatusCompleted, Sequence: 1, Agent: "code-reviewer", Iteration: 4,
	}
	statePath := filepath.Join(archive, orchestrator.RunStateFileName)
	if err := orchestrator.NewStateStore(statePath).Save(state); err != nil {
		t.Fatalf("save state: %v", err)
	}
	if err := orchestrator.NewTimeline(statePath).Append(
		orchestrator.Event{At: started, Kind: orchestrator.EventRunCompleted}); err != nil {
		t.Fatalf("append event: %v", err)
	}
}

func runMemoryCommand(t *testing.T, body func(*cobra.Command, []string) error) string {
	t.Helper()
	return captured(t, func(command *cobra.Command) {
		if err := body(command, nil); err != nil {
			t.Fatalf("command failed: %v", err)
		}
	})
}

// Backfill from the archive, in-process: the path that makes a deleted
// store recoverable.
func TestIngestBackfillsFromTheArchiveInProcess(t *testing.T) {
	projectWithStore(t)
	memoryArgs.dir = ""

	output := runMemoryCommand(t, runMemoryIngest)

	if !strings.Contains(output, "1 of 1 archived runs ingested") {
		t.Fatalf("ingest did not walk the archive:\n%s", output)
	}
	listed := runMemoryCommand(t, runMemoryRuns)
	if !strings.Contains(listed, "user-auth") {
		t.Errorf("runs does not list the ingested run:\n%s", listed)
	}
}

// The done-when, end to end and in-process.
func TestRetriesFindsTheIngestedRun(t *testing.T) {
	projectWithStore(t)
	memoryArgs.dir = ""
	memoryArgs.agent = "code-reviewer"
	memoryArgs.moreThan = 2
	if output := runMemoryCommand(t, runMemoryIngest); !strings.Contains(output, "ingested") {
		t.Fatalf("ingest failed:\n%s", output)
	}

	output := runMemoryCommand(t, runMemoryRetries)

	if !strings.Contains(output, "code-reviewer") || !strings.Contains(output, "4") {
		t.Errorf("retries did not find the four-attempt stage:\n%s", output)
	}
}

func TestMemoryCommandsEmitJSON(t *testing.T) {
	projectWithStore(t)
	memoryArgs.dir = ""
	_ = runMemoryCommand(t, runMemoryIngest)
	memoryArgs.asJSON = true
	t.Cleanup(func() { memoryArgs.asJSON = false })

	for name, body := range map[string]func(*cobra.Command, []string) error{
		"runs": runMemoryRuns, "retries": runMemoryRetries, "corrections": runMemoryCorrections,
	} {
		t.Run(name, func(t *testing.T) {
			var decoded []map[string]any
			if err := json.Unmarshal([]byte(runMemoryCommand(t, body)), &decoded); err != nil {
				t.Errorf("--json output does not parse: %v", err)
			}
		})
	}
}

// A project with no archive at all must say so rather than erroring.
func TestIngestWithNothingArchived(t *testing.T) {
	t.Chdir(t.TempDir())
	memoryArgs.dir = ""

	output := runMemoryCommand(t, runMemoryIngest)

	if !strings.Contains(output, "nothing to ingest") {
		t.Errorf("empty archive produced %q", output)
	}
}

// A directory that holds no run state is skipped with a reason, and the
// others still import — one unreadable archive must not stop twenty.
func TestIngestSkipsADirectoryWithoutRunState(t *testing.T) {
	projectWithStore(t)
	if err := os.MkdirAll(filepath.Join("docs", "features", "legacy"), 0o755); err != nil {
		t.Fatalf("create legacy dir: %v", err)
	}
	memoryArgs.dir = ""

	output := runMemoryCommand(t, runMemoryIngest)

	if !strings.Contains(output, "1 of 1 archived runs ingested") {
		t.Errorf("a markdown-only feature directory was not skipped quietly:\n%s", output)
	}
}
