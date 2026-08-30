package orchestrator

import "context"

// StageInput is what every stage invocation receives: the feature spec being
// delivered and the workspace directory artifacts are written into.
type StageInput struct {
	SpecPath     string
	WorkspaceDir string
}

// StageOutput is what a stage invocation produces. ArtifactPath may be empty
// when a stage produced no file artifact; when set, the executor records the
// artifact's SHA-256 in run state.
type StageOutput struct {
	ArtifactPath string
}

// Provider invokes one stage's agent and blocks until it finishes or ctx is
// done. Defined here — at the consumer — per go-conventions.md; implemented
// in internal/provider/* (mock for tests, claude for real runs). Kept to one
// method deliberately: retries, routing, and telemetry are later roadmap
// items and must not grow this interface.
type Provider interface {
	Invoke(ctx context.Context, stage Stage, input StageInput) (StageOutput, error)
}
