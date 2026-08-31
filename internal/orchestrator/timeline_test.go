package orchestrator_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/mock"
)

func readTimeline(t *testing.T, store *orchestrator.StateStore) []orchestrator.Event {
	t.Helper()
	events, err := orchestrator.NewTimeline(store.Path()).Read()
	if err != nil {
		t.Fatalf("read timeline: %v", err)
	}
	return events
}

func eventKinds(events []orchestrator.Event) []orchestrator.EventKind {
	kinds := make([]orchestrator.EventKind, 0, len(events))
	for _, event := range events {
		kinds = append(kinds, event.Kind)
	}
	return kinds
}

func assertKinds(t *testing.T, got []orchestrator.EventKind, want []orchestrator.EventKind) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("event kinds = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("event kinds = %v, want %v", got, want)
		}
	}
}

func TestTimelineRecordsEveryTransitionOfAHappyRun(t *testing.T) {
	executor, _, store, input := newHarness(t, completedScripts())

	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}

	events := readTimeline(t, store)
	assertKinds(t, eventKinds(events), []orchestrator.EventKind{
		orchestrator.EventRunStarted,
		orchestrator.EventStageStarted, orchestrator.EventStageCompleted,
		orchestrator.EventStageStarted, orchestrator.EventStageCompleted,
		orchestrator.EventStageStarted, orchestrator.EventStageCompleted,
		orchestrator.EventRunCompleted,
	})
	assertTimelineIsOrdered(t, events)
	if events[1].Stage != "analyst" || events[1].Sequence != 1 {
		t.Errorf("first stage event = %+v, want analyst #1", events[1])
	}
}

func assertTimelineIsOrdered(t *testing.T, events []orchestrator.Event) {
	t.Helper()
	for i := 1; i < len(events); i++ {
		if events[i].At.Before(events[i-1].At) {
			t.Fatalf("event %d is timestamped before its predecessor", i)
		}
		if events[i].At.IsZero() {
			t.Fatalf("event %d has no timestamp", i)
		}
	}
}

func TestTimelineRecordsGateHaltApprovalAndStaleness(t *testing.T) {
	executor, _, store, input := newHarness(t, completedScripts())
	plan := gatedPlan()
	_ = runUntilGate(t, executor, plan, input)
	if err := executor.Approve("confirm-design", orchestrator.ApprovalMethodTTY); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if err := executor.Run(context.Background(), plan, input); err != nil {
		t.Fatalf("run after approval: %v", err)
	}
	// Since L2.14 the edit also resets the gate, so this resume halts —
	// which is the point: the timeline records the halt, the approval, and
	// the staleness that caused it.
	editArtifact(t, input, "analyst", "# analysis, hand-edited")
	if err := executor.Run(context.Background(), plan, input); !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("resume after edit = %v, want a halt at the reset gate", err)
	}

	assertGateAndStaleEvents(t, firstOfEachKind(readTimeline(t, store)))
}

func assertGateAndStaleEvents(t *testing.T, found map[orchestrator.EventKind]orchestrator.Event) {
	t.Helper()
	if got := found[orchestrator.EventGateWaiting]; got.Gate != "confirm-design" || got.Stage != "developer" {
		t.Errorf("gate.waiting event = %+v", got)
	}
	if got := found[orchestrator.EventGateApproved]; got.ApprovalMethod != orchestrator.ApprovalMethodTTY {
		t.Errorf("gate.approved event = %+v", got)
	}
	if got := found[orchestrator.EventStageStale]; got.Stage != "analyst" || got.StaleReason != orchestrator.StaleReasonEdited {
		t.Errorf("stage.stale event = %+v", got)
	}
}

// firstOfEachKind keeps the first event of every kind: the cascade emits
// one stale event per demoted stage, and the first names the actual edit.
func firstOfEachKind(events []orchestrator.Event) map[orchestrator.EventKind]orchestrator.Event {
	found := map[orchestrator.EventKind]orchestrator.Event{}
	for _, event := range events {
		if _, seen := found[event.Kind]; !seen {
			found[event.Kind] = event
		}
	}
	return found
}

func TestTimelineRecordsFailureAndInterruption(t *testing.T) {
	boom := errors.New("agent exploded")
	executor, _, store, input := newHarness(t, map[string]mock.Script{
		"analyst":   {ArtifactContent: "# analysis"},
		"developer": {Err: boom},
	})
	if err := executor.Run(context.Background(), threeStagePlan(), input); !errors.Is(err, boom) {
		t.Fatalf("Run error = %v", err)
	}

	events := readTimeline(t, store)
	last := events[len(events)-1]
	if last.Kind != orchestrator.EventStageFailed || last.Stage != "developer" || last.Error == "" {
		t.Errorf("last event = %+v, want a stage.failed carrying the error", last)
	}
	for _, event := range events {
		if event.Kind == orchestrator.EventRunCompleted {
			t.Error("a failed run recorded run.completed")
		}
	}
}

func TestTimelineOfAnInterruptedRunIsReadable(t *testing.T) {
	executor, _, store, input := newHarness(t, map[string]mock.Script{
		"analyst":   {ArtifactContent: "# analysis"},
		"developer": {Hang: true},
	})
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(20*time.Millisecond, cancel)
	defer timer.Stop()
	if err := executor.Run(ctx, threeStagePlan(), input); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run error = %v, want context.Canceled", err)
	}
	cancel()

	events := readTimeline(t, store)
	if last := events[len(events)-1]; last.Kind != orchestrator.EventStageInterrupted || last.Stage != "developer" {
		t.Errorf("last event = %+v, want stage.interrupted for developer", last)
	}
	assertTimelineIsOrdered(t, events)
}

// TestTimelineSkipsATornTrailingLine covers the crash case the append-only
// design accepts: a process killed mid-write can leave half a line, which
// must not hide the history before it.
func TestTimelineSkipsATornTrailingLine(t *testing.T) {
	executor, _, store, input := newHarness(t, completedScripts())
	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}
	timeline := orchestrator.NewTimeline(store.Path())
	before := len(readTimeline(t, store))

	file, err := os.OpenFile(timeline.Path(), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open timeline: %v", err)
	}
	if _, err := file.WriteString(`{"at":"2026-08-30T00:00:00Z","kind":"stage.st`); err != nil {
		t.Fatalf("write torn line: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close timeline: %v", err)
	}

	if got := len(readTimeline(t, store)); got != before {
		t.Errorf("read %d events after a torn line, want the %d written before it", got, before)
	}
}

func TestTimelineReadOfAMissingFileIsEmpty(t *testing.T) {
	timeline := orchestrator.NewTimeline(filepath.Join(t.TempDir(), orchestrator.RunStateFileName))
	events, err := timeline.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("read %d events from a run that never started", len(events))
	}
}
