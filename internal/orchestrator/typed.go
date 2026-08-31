package orchestrator

// Typed stages (roadmap L2.9) exchange validated JSON instead of markdown.
// The executor validates a stage's payload, writes it as that stage's
// artifact, and hands the next stage only the fields its projection
// declares — no document, no summarization, no model on the data path.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orieken/loom/internal/state"
)

// typedStatePath is where a typed stage's document lives. It is recorded as
// the stage's artifact, so digest verification and the staleness cascade
// (L2.12) apply to it unchanged.
func typedStatePath(workspaceDir, stageID string) string {
	return filepath.Join(workspaceDir, state.TypedStateDir, stageID+".json")
}

// projectUpstream fills in the projected upstream state a typed stage
// reads. A stage with no declared upstream, or whose upstream has not run,
// receives nothing rather than an error — the run loop already refuses to
// reach a stage whose predecessor did not complete.
func projectUpstream(stage Stage, plan Plan, input StageInput) (StageInput, error) {
	if stage.Consumes == "" {
		return input, nil
	}
	payload, err := os.ReadFile(typedStatePath(input.WorkspaceDir, stage.Consumes))
	if os.IsNotExist(err) {
		return input, nil
	}
	if err != nil {
		return input, fmt.Errorf("read upstream state for stage %q: %w", stage.ID, err)
	}
	return withProjection(stage, plan, input, payload)
}

func withProjection(stage Stage, plan Plan, input StageInput, payload []byte) (StageInput, error) {
	upstreamKind, found := plan.stateKindOf(stage.Consumes)
	if !found {
		return input, fmt.Errorf("stage %q consumes %q, which produces no typed state", stage.ID, stage.Consumes)
	}
	projected, err := state.ProjectFor(state.Kind(stage.StateKind), state.Kind(upstreamKind), payload)
	if err != nil {
		return input, fmt.Errorf("project state for stage %q: %w", stage.ID, err)
	}
	input.UpstreamState = projected
	input.UpstreamStage = stage.Consumes
	return input, nil
}

// persistTypedOutput validates a typed stage's payload and writes it as the
// stage's artifact. An invalid payload fails the stage loudly: no repair
// prompt, no retry — those are L3.x, and a silent repair would hide the
// modelling failures this epic exists to surface.
func persistTypedOutput(stage Stage, input StageInput, output StageOutput) (string, error) {
	if len(output.Payload) == 0 {
		return "", fmt.Errorf("stage %q is typed but returned no state payload", stage.ID)
	}
	if _, err := state.Decode(state.Kind(stage.StateKind), output.Payload); err != nil {
		return "", fmt.Errorf("stage %q returned invalid state: %w", stage.ID, err)
	}
	return writeTypedState(stage, input, output.Payload)
}

func writeTypedState(stage Stage, input StageInput, payload []byte) (string, error) {
	path := typedStatePath(input.WorkspaceDir, stage.ID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create state directory: %w", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return "", fmt.Errorf("write state for stage %q: %w", stage.ID, err)
	}
	return path, nil
}

// stateKindOf returns the typed state kind a stage produces.
func (p Plan) stateKindOf(stageID string) (string, bool) {
	for _, stage := range p.Stages {
		if stage.ID == stageID {
			return stage.StateKind, stage.StateKind != ""
		}
	}
	return "", false
}
