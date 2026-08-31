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
	onStale  func([]StaleStage)
}

// NewExecutor wires a provider and a state store into an executor.
func NewExecutor(provider Provider, store *StateStore) *Executor {
	return &Executor{provider: provider, store: store}
}

// OnStale registers a callback invoked once per Run, before any stage
// executes, with the stages that digest verification demoted. It is how the
// CLI tells the human which completed work is being redone and why.
func (e *Executor) OnStale(report func([]StaleStage)) {
	e.onStale = report
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
	for _, stage := range plan.Stages {
		if state.IsStageCompleted(stage.ID) {
			continue
		}
		if err := e.checkGate(stage, state); err != nil {
			return err
		}
		if err := e.runStage(ctx, stage, input, state); err != nil {
			return err
		}
	}
	return nil
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
	if stage.Gate == "" || state.IsGateApproved(stage.Gate) {
		return nil
	}
	record := StageRecord{Status: StageStatusWaitingApproval, StartedAt: time.Now().UTC(), Gate: stage.Gate}
	if err := e.persistStatus(state, stage.ID, record); err != nil {
		return err
	}
	return &WaitingApprovalError{Gate: stage.Gate, Stage: stage.ID}
}

func (e *Executor) runStage(ctx context.Context, stage Stage, input StageInput, state *RunState) error {
	if err := e.persistStatus(state, stage.ID, StageRecord{Status: StageStatusRunning, StartedAt: time.Now().UTC()}); err != nil {
		return err
	}
	output, invokeErr := e.invoke(ctx, stage, input)
	if invokeErr != nil {
		return e.persistFailure(ctx, state, stage, invokeErr)
	}
	return e.persistCompletion(state, stage, output)
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
	return fmt.Errorf("stage %q: %w", stage.ID, invokeErr)
}

func (e *Executor) persistCompletion(state *RunState, stage Stage, output StageOutput) error {
	record := state.Stages[stage.ID]
	now := time.Now().UTC()
	record.FinishedAt = &now
	record.Status = StageStatusCompleted
	record.ArtifactPath = output.ArtifactPath
	if output.ArtifactPath != "" {
		sum, err := ArtifactSHA256(output.ArtifactPath)
		if err != nil {
			return fmt.Errorf("stage %q completed but its artifact is unreadable: %w", stage.ID, err)
		}
		record.ArtifactSHA256 = sum
	}
	return e.persistStatus(state, stage.ID, record)
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
