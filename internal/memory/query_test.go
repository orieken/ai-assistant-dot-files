package memory_test

import (
	"testing"
	"time"

	"github.com/orieken/loom/internal/memory"
	"github.com/orieken/loom/internal/orchestrator"
)

// runWithRetries builds a run whose code-reviewer loop took a known number
// of rounds, so the done-when query is tested against counts the test
// chose rather than against whatever a fixture happened to contain.
func runWithRetries(feature string, started time.Time, iterations int) (*orchestrator.RunState, []orchestrator.Event) {
	state := orchestrator.NewRunState("deliver-feature", orchestrator.CreatedByExecutor)
	state.FeatureName = feature
	state.StartedAt = started
	state.Stages["code-reviewer"] = orchestrator.StageRecord{
		Status: orchestrator.StageStatusCompleted, Sequence: 1,
		Agent: "code-reviewer", Iteration: iterations,
	}
	state.Stages["analyst"] = orchestrator.StageRecord{
		Status: orchestrator.StageStatusCompleted, Sequence: 2, Agent: "analyst", Iteration: 1,
	}
	return state, []orchestrator.Event{{At: started, Kind: orchestrator.EventRunCompleted}}
}

func storeWithRuns(t *testing.T, iterations ...int) *memory.Store {
	t.Helper()
	store := openStore(t)
	base := time.Date(2026, 9, 2, 9, 0, 0, 0, time.UTC)
	for index, count := range iterations {
		state, events := runWithRetries(
			"feature-"+string(rune('a'+index)), base.Add(time.Duration(index)*time.Hour), count)
		if _, err := store.Ingest(state, events); err != nil {
			t.Fatalf("ingest run %d: %v", index, err)
		}
	}
	return store
}

// L3.5's done-when, verbatim: "show me every run where code-reviewer
// retried more than twice". Read as iteration count > 2 — three attempts.
func TestTheDoneWhenQuery(t *testing.T) {
	store := storeWithRuns(t, 1, 3, 2, 5)

	rows, err := store.Retries("code-reviewer", 2)
	if err != nil {
		t.Fatalf("Retries: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("got %d rows, want the two runs with more than 2 iterations: %+v", len(rows), rows)
	}
	// Ordered by iteration count descending, so the worst offender is first.
	if rows[0].Iterations != 5 || rows[1].Iterations != 3 {
		t.Errorf("iterations = %d, %d; want 5 then 3", rows[0].Iterations, rows[1].Iterations)
	}
	for _, row := range rows {
		if row.Agent != "code-reviewer" {
			t.Errorf("row names agent %q, want code-reviewer only", row.Agent)
		}
	}
}

// The boundary is where an off-by-one would hide: 2 iterations must not
// appear in "more than 2".
func TestRetryThresholdIsExclusive(t *testing.T) {
	store := storeWithRuns(t, 2)

	rows, err := store.Retries("code-reviewer", 2)
	if err != nil {
		t.Fatalf("Retries: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("a stage with exactly 2 iterations matched \"more than 2\": %+v", rows)
	}
}

// An empty agent means every agent, so the same query answers "which stages
// anywhere took more than N attempts".
func TestRetriesAcrossEveryAgent(t *testing.T) {
	store := storeWithRuns(t, 4)

	rows, err := store.Retries("", 0)
	if err != nil {
		t.Fatalf("Retries: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want both stages of the run: %+v", len(rows), rows)
	}
}

func TestRunsListsNewestFirst(t *testing.T) {
	store := storeWithRuns(t, 1, 1, 1)

	runs, err := store.Runs(10)
	if err != nil {
		t.Fatalf("Runs: %v", err)
	}
	if len(runs) != 3 {
		t.Fatalf("got %d runs, want 3", len(runs))
	}
	for index := 1; index < len(runs); index++ {
		if runs[index-1].StartedAt < runs[index].StartedAt {
			t.Errorf("runs are not newest-first: %q before %q", runs[index-1].StartedAt, runs[index].StartedAt)
		}
	}
}

func TestRunsRespectsTheLimit(t *testing.T) {
	store := storeWithRuns(t, 1, 1, 1)

	runs, err := store.Runs(2)
	if err != nil {
		t.Fatalf("Runs: %v", err)
	}
	if len(runs) != 2 {
		t.Errorf("got %d runs, want the limit of 2", len(runs))
	}
}

func TestCorrectionsRankAgents(t *testing.T) {
	store := openStore(t)
	state, events := fixtureRun()
	state.Corrections = append(state.Corrections, orchestrator.Correction{
		Stage: "analyst", Agent: "analyst", Gate: "confirm-design",
		Stat: orchestrator.DiffStat{Added: 2, Removed: 0},
	})
	if _, err := store.Ingest(state, events); err != nil {
		t.Fatalf("Ingest: %v", err)
	}

	rows, err := store.Corrections()
	if err != nil {
		t.Fatalf("Corrections: %v", err)
	}
	if len(rows) != 1 || rows[0].Agent != "analyst" {
		t.Fatalf("rows = %+v, want the analyst alone", rows)
	}
	want := memory.AgentCorrections{Agent: "analyst", Corrections: 2, LinesAdded: 6, LinesCut: 1, Runs: 1}
	if rows[0] != want {
		t.Errorf("row = %+v, want %+v", rows[0], want)
	}
}

// A store full of zeroes means the provider reported nothing, not that the
// runs were free. A cost report that cannot tell those apart is exactly the
// confident wrong number this framework keeps removing.
func TestCostReportedDistinguishesUnmeasuredFromFree(t *testing.T) {
	store := storeWithRuns(t, 1)

	reported, err := store.CostReported()
	if err != nil {
		t.Fatalf("CostReported: %v", err)
	}
	if reported {
		t.Error("CostReported is true for runs that reported no usage at all")
	}

	state, events := fixtureRun()
	if _, err := store.Ingest(state, events); err != nil {
		t.Fatalf("Ingest: %v", err)
	}
	if reported, err = store.CostReported(); err != nil || !reported {
		t.Errorf("CostReported = %v (err %v) after ingesting a run with usage", reported, err)
	}
}

// Every query must return nothing rather than erroring on a store nobody
// has written to.
func TestQueriesOnAnEmptyStoreReturnNothing(t *testing.T) {
	store := openStore(t)

	runs, err := store.Runs(10)
	if err != nil || len(runs) != 0 {
		t.Errorf("Runs on an empty store = (%v, %v)", runs, err)
	}
	retries, err := store.Retries("", 0)
	if err != nil || len(retries) != 0 {
		t.Errorf("Retries on an empty store = (%v, %v)", retries, err)
	}
	corrections, err := store.Corrections()
	if err != nil || len(corrections) != 0 {
		t.Errorf("Corrections on an empty store = (%v, %v)", corrections, err)
	}
}
