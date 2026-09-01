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
	for _, upstream := range stage.Consumes {
		projected, err := projectOne(stage, plan, input, upstream)
		if err != nil {
			return input, err
		}
		if projected == nil {
			continue
		}
		input.UpstreamState = withUpstream(input.UpstreamState, upstream, projected)
	}
	return input, nil
}

// projectOne returns the projected fields of one upstream stage, or nil
// when that stage has not produced state yet — a stage the route skipped,
// or one the loop has not reached. The run loop already refuses to reach a
// stage whose predecessor did not complete, so absence here is legitimate.
func projectOne(stage Stage, plan Plan, input StageInput, upstream string) ([]byte, error) {
	payload, err := os.ReadFile(typedStatePath(input.WorkspaceDir, upstream))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %q state for stage %q: %w", upstream, stage.ID, err)
	}
	upstreamKind, found := plan.stateKindOf(upstream)
	if !found {
		return nil, fmt.Errorf("stage %q consumes %q, which produces no typed state", stage.ID, upstream)
	}
	projected, err := state.ProjectionFor(stage.ID, state.Kind(upstreamKind), payload)
	if err != nil {
		return nil, fmt.Errorf("project %q state for stage %q: %w", upstream, stage.ID, err)
	}
	return projected, nil
}

// withUpstream adds one projection to the input's map without the caller
// having to worry about whether it exists yet.
func withUpstream(existing map[string][]byte, upstream string, projected []byte) map[string][]byte {
	if existing == nil {
		existing = map[string][]byte{}
	}
	existing[upstream] = projected
	return existing
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
	path, err := writeTypedState(stage, input, output.Payload)
	if err != nil {
		return "", err
	}
	return path, renderView(stage, input, output.Payload)
}

// renderView writes the human-readable markdown for a typed stage under the
// contract's filename — `analysis.md`, not `analyst.md` — because the
// stages that are still untyped were told to read the contract's name. The
// view is derived, so it is deliberately not digest-tracked: editing it
// must not be able to corrupt a run.
func renderView(stage Stage, input StageInput, payload []byte) error {
	name, body, err := state.RenderView(state.Kind(stage.StateKind), payload)
	if err != nil {
		return fmt.Errorf("render view for stage %q: %w", stage.ID, err)
	}
	if err := os.WriteFile(filepath.Join(input.WorkspaceDir, name), []byte(body), 0o644); err != nil {
		return fmt.Errorf("write view for stage %q: %w", stage.ID, err)
	}
	return nil
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
