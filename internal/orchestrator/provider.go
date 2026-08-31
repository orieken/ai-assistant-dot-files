package orchestrator

import "context"

// StageInput is what every stage invocation receives: the feature spec being
// delivered and the workspace directory artifacts are written into.
type StageInput struct {
	SpecPath     string
	WorkspaceDir string
	// UpstreamState is the projected slice of the previous stage's typed
	// state this stage is allowed to read (roadmap L2.9) — the fields its
	// contract declares, not the whole document. Empty for stages that
	// still exchange markdown.
	UpstreamState []byte
	// UpstreamStage names where UpstreamState came from, so a provider can
	// tell the agent what it is reading.
	UpstreamStage string
}

// StageOutput is what a stage invocation produces. ArtifactPath may be empty
// when a stage produced no file artifact; when set, the executor records the
// artifact's SHA-256 in run state.
type StageOutput struct {
	ArtifactPath string
	// Payload is the JSON state document a typed stage produced. The
	// executor validates it against the stage's schema and writes it as the
	// stage's artifact; stages that still write markdown leave it empty.
	Payload []byte
}

// Provider invokes one stage's agent and blocks until it finishes or ctx is
// done. Defined here — at the consumer — per go-conventions.md; implemented
// in internal/provider/* (mock for tests, claude for real runs). Kept to one
// method deliberately: retries, routing, and telemetry are later roadmap
// items and must not grow this interface.
type Provider interface {
	Invoke(ctx context.Context, stage Stage, input StageInput) (StageOutput, error)
}
