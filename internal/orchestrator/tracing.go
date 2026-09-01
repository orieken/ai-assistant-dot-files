package orchestrator

// Tracing is the executor's telemetry seam (roadmap L3.8). The interface
// lives here, with its consumer, and `internal/telemetry` implements it —
// so no OpenTelemetry type reaches the use-case layer at all. That is
// architecture-guardrails.md #8 held structurally rather than by review:
// spans are emitted from the adapter layer because the inner layers cannot
// name a span if they try.
//
// It is deliberately semantic rather than generic. There is no
// SetAttribute(string, any) here, because every fact the executor knows
// about a run is already a typed field on RunSpan or StageSpan. A generic
// attribute bag would push the vocabulary decision out to call sites and
// let two of them disagree about what "stage" means.
//
// This does NOT replace run-events.jsonl. The timeline is the audit record
// — gates, digests, staleness — and stays readable with no collector and no
// exporter configured. Spans are the timing and cost record. They overlap
// in subject and answer different questions.

import "context"

// RunSpan describes the run a root span covers.
type RunSpan struct {
	Plan      string
	Feature   string
	SpecPath  string
	StateFile string
}

// StageSpan describes one stage execution. A loop that sends a stage round
// again produces another StageSpan with a higher Iteration, not a nested
// one: iterations are retries of one stage, and a trace viewer should show
// them side by side rather than three levels deep.
type StageSpan struct {
	ID        string
	Agent     string
	Sequence  int
	Iteration int
	Gate      string
	Internal  bool
}

// SpanOutcome is how a span ended. Status carries the executor's own
// vocabulary (COMPLETED, FAILED, SKIPPED, INTERRUPTED…) so a trace and the
// run state cannot describe the same stage differently.
type SpanOutcome struct {
	Status StageStatus
	Reason string
	Err    error
}

// Span is one open unit of work. End is safe to call on the zero value of
// any implementation, so callers never guard it.
type Span interface {
	End(outcome SpanOutcome)
}

// Tracer opens spans. A nil Tracer disables tracing entirely, which is the
// default: nothing is exported and nothing is measured beyond what the
// timeline already records.
type Tracer interface {
	StartRun(ctx context.Context, run RunSpan) (context.Context, Span)
	StartStage(ctx context.Context, stage StageSpan) (context.Context, Span)
}

// noopTracer is what an executor without a Tracer uses, so the run loop
// never branches on whether tracing is on.
type noopTracer struct{}

func (noopTracer) StartRun(ctx context.Context, _ RunSpan) (context.Context, Span) {
	return ctx, noopSpan{}
}

func (noopTracer) StartStage(ctx context.Context, _ StageSpan) (context.Context, Span) {
	return ctx, noopSpan{}
}

type noopSpan struct{}

func (noopSpan) End(SpanOutcome) {}

// WithTracer registers the tracer the run loop opens spans on. Passing nil
// restores the no-op tracer.
func (e *Executor) WithTracer(tracer Tracer) {
	if tracer == nil {
		e.tracer = noopTracer{}
		return
	}
	e.tracer = tracer
}
