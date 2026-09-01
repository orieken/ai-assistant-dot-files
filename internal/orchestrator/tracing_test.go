package orchestrator_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/mock"
)

// recordingTracer captures the span tree without any OpenTelemetry
// involvement, which is the point of the seam: the executor's tracing
// behaviour is testable without an SDK, an exporter, or a collector.
type recordingTracer struct {
	mutex  sync.Mutex
	runs   []orchestrator.RunSpan
	stages []orchestrator.StageSpan
	ended  []orchestrator.SpanOutcome
	open   int
}

func (r *recordingTracer) StartRun(ctx context.Context, run orchestrator.RunSpan) (context.Context, orchestrator.Span) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.runs = append(r.runs, run)
	r.open++
	return ctx, &recordingSpan{tracer: r}
}

func (r *recordingTracer) StartStage(ctx context.Context, stage orchestrator.StageSpan) (context.Context, orchestrator.Span) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.stages = append(r.stages, stage)
	r.open++
	return ctx, &recordingSpan{tracer: r}
}

type recordingSpan struct {
	tracer *recordingTracer
}

func (s *recordingSpan) End(outcome orchestrator.SpanOutcome) {
	s.tracer.mutex.Lock()
	defer s.tracer.mutex.Unlock()
	s.tracer.ended = append(s.tracer.ended, outcome)
	s.tracer.open--
}

func (r *recordingTracer) stageIDs() []string {
	ids := make([]string, 0, len(r.stages))
	for _, stage := range r.stages {
		ids = append(ids, stage.ID)
	}
	return ids
}

func (r *recordingTracer) statuses() []orchestrator.StageStatus {
	statuses := make([]orchestrator.StageStatus, 0, len(r.ended))
	for _, outcome := range r.ended {
		statuses = append(statuses, outcome.Status)
	}
	return statuses
}

func assertStrings(t *testing.T, got, want []string, label string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s = %v, want %v", label, got, want)
		}
	}
}

func TestRunTracesOneSpanPerStageUnderOneRunSpan(t *testing.T) {
	scripts := map[string]mock.Script{
		"analyst":     {ArtifactContent: "# analysis"},
		"developer":   {ArtifactContent: "# implementation"},
		"qa-engineer": {ArtifactContent: "# qa report"},
	}
	executor, _, _, input := newHarness(t, scripts)
	tracer := &recordingTracer{}
	executor.WithTracer(tracer)

	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(tracer.runs) != 1 {
		t.Fatalf("run spans = %d, want exactly 1", len(tracer.runs))
	}
	if tracer.runs[0].Plan != "test-plan" {
		t.Errorf("run span plan = %q, want %q", tracer.runs[0].Plan, "test-plan")
	}
	assertStrings(t, tracer.stageIDs(), []string{"analyst", "developer", "qa-engineer"}, "traced stages")
	if tracer.open != 0 {
		t.Errorf("%d spans left open — every span must be closed on every path", tracer.open)
	}
}

// A span must report the status run state records, not a status inferred
// from whether an error came back. The two describing the same stage
// differently is the defect the shared vocabulary exists to prevent.
func TestStageSpanReportsThePersistedStatus(t *testing.T) {
	boom := errors.New("agent exploded")
	scripts := map[string]mock.Script{
		"analyst":   {ArtifactContent: "# analysis"},
		"developer": {Err: boom},
	}
	executor, _, store, input := newHarness(t, scripts)
	tracer := &recordingTracer{}
	executor.WithTracer(tracer)

	if err := executor.Run(context.Background(), threeStagePlan(), input); !errors.Is(err, boom) {
		t.Fatalf("Run error = %v, want %v", err, boom)
	}

	state := mustLoad(t, store)
	want := []orchestrator.StageStatus{
		orchestrator.StageStatusCompleted, // analyst stage span
		orchestrator.StageStatusFailed,    // developer stage span
		orchestrator.StageStatusFailed,    // run span
	}
	got := tracer.statuses()
	if len(got) != len(want) {
		t.Fatalf("span statuses = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("span statuses = %v, want %v", got, want)
		}
	}
	if state.Stages["developer"].Status != orchestrator.StageStatusFailed {
		t.Error("precondition: run state should record developer as FAILED")
	}
}

// A run that halts at a gate is waiting on a human, not failing. It must
// still close its root span, or everything it recorded before the barrier
// is lost when the process exits.
func TestGateHaltClosesTheRunSpanAsWaiting(t *testing.T) {
	plan := orchestrator.Plan{
		Name: "gated-plan",
		Stages: []orchestrator.Stage{
			{ID: "analyst", Agent: "analyst", Timeout: 5 * time.Second},
			{ID: "developer", Agent: "developer", Gate: "confirm-design", Timeout: 5 * time.Second},
		},
	}
	executor, _, _, input := newHarness(t, map[string]mock.Script{"analyst": {ArtifactContent: "# analysis"}})
	tracer := &recordingTracer{}
	executor.WithTracer(tracer)

	err := executor.Run(context.Background(), plan, input)
	if !errors.Is(err, orchestrator.ErrWaitingApproval) {
		t.Fatalf("Run error = %v, want ErrWaitingApproval", err)
	}
	if tracer.open != 0 {
		t.Errorf("%d spans left open after a gate halt", tracer.open)
	}
	last := tracer.ended[len(tracer.ended)-1]
	if last.Status != orchestrator.StageStatusWaitingApproval {
		t.Errorf("run span status = %q, want %q", last.Status, orchestrator.StageStatusWaitingApproval)
	}
}

// With no tracer registered the run loop must behave identically, and cost
// nothing. Tracing is off by default and must never be load-bearing.
func TestRunWithoutATracerIsUnchanged(t *testing.T) {
	executor, provider, store, input := newHarness(t, map[string]mock.Script{
		"analyst":     {ArtifactContent: "# analysis"},
		"developer":   {ArtifactContent: "# implementation"},
		"qa-engineer": {ArtifactContent: "# qa report"},
	})

	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}
	assertInvocations(t, provider, []string{"analyst", "developer", "qa-engineer"})
	assertStatus(t, mustLoad(t, store), "qa-engineer", orchestrator.StageStatusCompleted)
}

// WithTracer(nil) must restore the no-op rather than panic on the next span.
func TestWithNilTracerDisablesTracing(t *testing.T) {
	executor, _, _, input := newHarness(t, map[string]mock.Script{
		"analyst":     {ArtifactContent: "# analysis"},
		"developer":   {ArtifactContent: "# implementation"},
		"qa-engineer": {ArtifactContent: "# qa report"},
	})
	executor.WithTracer(&recordingTracer{})
	executor.WithTracer(nil)

	if err := executor.Run(context.Background(), threeStagePlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}
}
