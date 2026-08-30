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
}

// NewExecutor wires a provider and a state store into an executor.
func NewExecutor(provider Provider, store *StateStore) *Executor {
	return &Executor{provider: provider, store: store}
}

// Run executes every stage of the plan that is not already COMPLETED in
// persisted state, stopping on the first failure. Cancelling ctx (SIGINT)
// persists a clean INTERRUPTED checkpoint for the in-flight stage before
// returning, so a later Run resumes by re-running that stage.
func (e *Executor) Run(ctx context.Context, plan Plan, input StageInput) error {
	state, err := e.prepareState(plan)
	if err != nil {
		return err
	}
	for _, stage := range plan.Stages {
		if state.IsStageCompleted(stage.ID) {
			continue
		}
		if err := e.runStage(ctx, stage, input, state); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) prepareState(plan Plan) (*RunState, error) {
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	state, err := e.store.Load()
	if err != nil {
		return nil, err
	}
	if state == nil {
		return NewRunState(plan.Name), nil
	}
	if state.PlanName != plan.Name {
		return nil, fmt.Errorf("run state %s belongs to plan %q, not %q — refusing to mix runs",
			e.store.Path(), state.PlanName, plan.Name)
	}
	return state, nil
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
	state.Stages[stageID] = record
	if err := e.store.Save(state); err != nil {
		return fmt.Errorf("persist state for stage %q: %w", stageID, err)
	}
	return nil
}
