package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Executor runs a Plan's stages in order against a Provider, persisting
// state after every transition (stage started, completed, failed,
// interrupted) so a run can always resume from its last checkpoint.
type Executor struct {
	provider Provider
	store    *StateStore
	timeline *Timeline
	onStale  func([]StaleStage)
	onReset  func(*StaleApprovalError)
}

// NewExecutor wires a provider and a state store into an executor.
func NewExecutor(provider Provider, store *StateStore) *Executor {
	return &Executor{provider: provider, store: store, timeline: NewTimeline(store.Path())}
}

// emit appends one event to the run's timeline. Timestamps come from the
// clock here, not from anything a provider reported.
func (e *Executor) emit(event Event) error {
	return e.timeline.Append(event)
}

// OnStale registers a callback invoked once per Run, before any stage
// executes, with the stages that digest verification demoted. It is how the
// CLI tells the human which completed work is being redone and why.
func (e *Executor) OnStale(report func([]StaleStage)) {
	e.onStale = report
}

// OnApprovalReset registers a callback invoked when an edit resets a gate's
// approval. It is how the CLI explains a halt that a human's own edit
// caused, rather than leaving them to work it out from the state file.
func (e *Executor) OnApprovalReset(report func(*StaleApprovalError)) {
	e.onReset = report
}

// Run executes every stage of the plan that is not already COMPLETED in
// persisted state, stopping on the first failure. Cancelling ctx (SIGINT)
// persists a clean INTERRUPTED checkpoint for the in-flight stage before
// returning, so a later Run resumes by re-running that stage.
func (e *Executor) Run(ctx context.Context, plan Plan, input StageInput) error {
	state, err := e.prepareState(plan, input)
	if err != nil {
		return err
	}
	if err := e.emit(Event{Kind: EventRunStarted}); err != nil {
		return err
	}
	for _, stage := range plan.Stages {
		if state.IsStageCompleted(stage.ID) {
			continue
		}
		if err := e.checkGate(stage, state); err != nil {
			return err
		}
		if err := e.runStage(ctx, stage, plan, input, state); err != nil {
			return err
		}
	}
	return e.emit(Event{Kind: EventRunCompleted})
}

func (e *Executor) prepareState(plan Plan, input StageInput) (*RunState, error) {
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	state, err := e.store.Load()
	if err != nil {
		return nil, err
	}
	if state == nil {
		return newRunFor(plan, input), nil
	}
	if err := state.CheckCreatedBy(CreatedByExecutor); err != nil {
		return nil, err
	}
	if err := checkStateBelongsToRun(state, plan, input, e.store.Path()); err != nil {
		return nil, err
	}
	return e.verifyResumedState(state)
}

func (e *Executor) emitStale(state *RunState, stale []StaleStage) error {
	for _, item := range stale {
		event := Event{Kind: EventStageStale, Stage: item.StageID, StaleReason: item.Reason,
			Sequence: state.Stages[item.StageID].Sequence}
		if err := e.emit(event); err != nil {
			return err
		}
	}
	return nil
}

func newRunFor(plan Plan, input StageInput) *RunState {
	state := NewRunState(plan.Name, CreatedByExecutor)
	state.SpecPath = input.SpecPath
	state.FeatureName = FeatureNameFromSpec(input.SpecPath)
	return state
}

// checkStateBelongsToRun refuses state from a different plan or a different
// spec: resuming one feature's run against another's spec would replay the
// wrong work against the right state.
func checkStateBelongsToRun(state *RunState, plan Plan, input StageInput, path string) error {
	if state.PlanName != plan.Name {
		return fmt.Errorf("run state %s belongs to plan %q, not %q — refusing to mix runs",
			path, state.PlanName, plan.Name)
	}
	if state.SpecPath != "" && input.SpecPath != "" && state.SpecPath != input.SpecPath {
		return fmt.Errorf("run state %s belongs to spec %q, not %q — refusing to mix runs",
			path, state.SpecPath, input.SpecPath)
	}
	return nil
}

// verifyResumedState is the L2.12 integrity check: every completed stage's
// artifact is re-hashed in Go before the run loop trusts it.
func (e *Executor) verifyResumedState(state *RunState) (*RunState, error) {
	stale := VerifyCompletedStages(state)
	if len(stale) == 0 {
		return state, nil
	}
	if err := e.store.Save(state); err != nil {
		return nil, fmt.Errorf("persist verification result: %w", err)
	}
	if err := e.emitStale(state, stale); err != nil {
		return nil, err
	}
	if err := e.invalidateApprovalsFor(state, stale); err != nil {
		return nil, err
	}
	if err := e.store.Save(state); err != nil {
		return nil, fmt.Errorf("persist approval invalidation: %w", err)
	}
	if e.onStale != nil {
		e.onStale(stale)
	}
	return state, nil
}

// sequenceFor keeps a stage's original position when it is re-recorded: a
// re-run after an interrupt or a review loop is the same step of the run,
// not a new one.
func sequenceFor(state *RunState, stageID string) int {
	if existing := state.Stages[stageID].Sequence; existing > 0 {
		return existing
	}
	return state.NextSequence()
}

// checkGate is the process interrupt: a gated stage cannot start until run
// state records an approval for its gate. Nothing a provider returns can
// create that approval — only the CLI approval channels write it (roadmap
// L2.13). The halt is persisted so the run resumes exactly here.
func (e *Executor) checkGate(stage Stage, state *RunState) error {
	if stage.Gate == "" {
		return nil
	}
	if state.IsGateApproved(stage.Gate) {
		invalidated, err := e.enforceApprovalBinding(state, stage)
		if err != nil {
			return err
		}
		if !invalidated {
			return nil
		}
	}
	record := StageRecord{Status: StageStatusWaitingApproval, StartedAt: time.Now().UTC(), Gate: stage.Gate}
	if err := e.persistStatus(state, stage.ID, record); err != nil {
		return err
	}
	if err := e.emit(Event{Kind: EventGateWaiting, Stage: stage.ID, Gate: stage.Gate}); err != nil {
		return err
	}
	return &WaitingApprovalError{Gate: stage.Gate, Stage: stage.ID}
}

func (e *Executor) runStage(ctx context.Context, stage Stage, plan Plan, input StageInput, state *RunState) error {
	if err := e.persistStatus(state, stage.ID, StageRecord{Status: StageStatusRunning, StartedAt: time.Now().UTC()}); err != nil {
		return err
	}
	if err := e.emit(Event{Kind: EventStageStarted, Stage: stage.ID, Sequence: state.Stages[stage.ID].Sequence}); err != nil {
		return err
	}
	input, projectErr := projectUpstream(stage, plan, input)
	if projectErr != nil {
		return e.persistFailure(ctx, state, stage, projectErr)
	}
	output, invokeErr := e.invoke(ctx, stage, input)
	if invokeErr != nil {
		return e.persistFailure(ctx, state, stage, invokeErr)
	}
	return e.persistCompletion(state, stage, input, output)
}

// invoke calls the provider under the stage's timeout. A zero timeout means
// no per-stage deadline beyond the parent context.
func (e *Executor) invoke(ctx context.Context, stage Stage, input StageInput) (StageOutput, error) {
	if stage.Timeout <= 0 {
		return e.provider.Invoke(ctx, stage, input)
	}
	stageCtx, cancel := context.WithTimeout(ctx, stage.Timeout)
	defer cancel()
	return e.provider.Invoke(stageCtx, stage, input)
}

// artifactFor resolves what this stage's artifact is: a typed stage's
// validated state document, written here, or the markdown file a provider
// wrote itself.
func artifactFor(stage Stage, input StageInput, output StageOutput) (string, error) {
	if stage.StateKind == "" {
		return output.ArtifactPath, nil
	}
	return persistTypedOutput(stage, input, output)
}

// persistFailure distinguishes parent cancellation (SIGINT — checkpoint as
// INTERRUPTED, resumable) from a stage failure or timeout (FAILED, stops
// the run).
func (e *Executor) persistFailure(ctx context.Context, state *RunState, stage Stage, invokeErr error) error {
	record := state.Stages[stage.ID]
	now := time.Now().UTC()
	record.FinishedAt = &now
	record.Error = invokeErr.Error()
	if errors.Is(ctx.Err(), context.Canceled) {
		record.Status = StageStatusInterrupted
	} else {
		record.Status = StageStatusFailed
	}
	if err := e.persistStatus(state, stage.ID, record); err != nil {
		return errors.Join(invokeErr, err)
	}
	if err := e.emitStageEnd(state, stage.ID, record); err != nil {
		return errors.Join(invokeErr, err)
	}
	return fmt.Errorf("stage %q: %w", stage.ID, invokeErr)
}

func (e *Executor) persistCompletion(state *RunState, stage Stage, input StageInput, output StageOutput) error {
	artifactPath, err := artifactFor(stage, input, output)
	if err != nil {
		return e.persistFailure(context.Background(), state, stage, err)
	}
	record := state.Stages[stage.ID]
	now := time.Now().UTC()
	record.FinishedAt = &now
	record.Status = StageStatusCompleted
	record.ArtifactPath = artifactPath
	if artifactPath != "" {
		sum, err := ArtifactSHA256(artifactPath)
		if err != nil {
			return fmt.Errorf("stage %q completed but its artifact is unreadable: %w", stage.ID, err)
		}
		record.ArtifactSHA256 = sum
	}
	if err := e.persistStatus(state, stage.ID, record); err != nil {
		return err
	}
	return e.emitStageEnd(state, stage.ID, record)
}

// emitStageEnd records how a stage finished, using the same record the
// state file holds so the two never disagree.
func (e *Executor) emitStageEnd(state *RunState, stageID string, record StageRecord) error {
	kinds := map[StageStatus]EventKind{
		StageStatusCompleted:   EventStageCompleted,
		StageStatusFailed:      EventStageFailed,
		StageStatusInterrupted: EventStageInterrupted,
	}
	kind, ok := kinds[record.Status]
	if !ok {
		return nil
	}
	return e.emit(Event{Kind: kind, Stage: stageID, Sequence: state.Stages[stageID].Sequence, Error: record.Error})
}

func (e *Executor) persistStatus(state *RunState, stageID string, record StageRecord) error {
	if record.Sequence == 0 {
		record.Sequence = sequenceFor(state, stageID)
	}
	state.Stages[stageID] = record
	if err := e.store.Save(state); err != nil {
		return fmt.Errorf("persist state for stage %q: %w", stageID, err)
	}
	return nil
}
