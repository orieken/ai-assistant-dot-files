// Package mock is a deterministic Provider implementation for executor
// tests: every stage's behavior is scripted — a canned artifact, a scripted
// failure, or a hang that blocks until the context is done (exercising
// timeouts and SIGINT checkpointing). Never used outside tests; the real
// provider spawning `claude -p` lands in internal/provider/claude (M0.4
// part 2).
package mock

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/orieken/loom/internal/orchestrator"
)

// Script describes what the mock does when a stage is invoked. Exactly one
// behavior applies, checked in order: Hang, Err, then artifact success.
type Script struct {
	// Hang blocks until ctx is done and returns ctx.Err() — used to
	// exercise stage timeouts and cancellation checkpoints.
	Hang bool
	// Err fails the invocation with this error.
	Err error
	// ArtifactContent, when non-empty, is written to <WorkspaceDir>/<stage ID>.md
	// and that path is returned as the stage's artifact.
	ArtifactContent string
}

// Provider is a scripted, deterministic orchestrator.Provider.
type Provider struct {
	mu          sync.Mutex
	scripts     map[string]Script
	invocations []string
}

// New returns a mock provider with per-stage-ID scripts. Stages without a
// script succeed with no artifact.
func New(scripts map[string]Script) *Provider {
	if scripts == nil {
		scripts = map[string]Script{}
	}
	return &Provider{scripts: scripts}
}

// Invoke runs the stage's script. It records the invocation order so tests
// can assert which stages actually ran (e.g. resume skips completed ones).
func (p *Provider) Invoke(ctx context.Context, stage orchestrator.Stage, input orchestrator.StageInput) (orchestrator.StageOutput, error) {
	p.record(stage.ID)
	script := p.script(stage.ID)
	if script.Hang {
		<-ctx.Done()
		return orchestrator.StageOutput{}, ctx.Err()
	}
	if script.Err != nil {
		return orchestrator.StageOutput{}, script.Err
	}
	return p.writeArtifact(stage, input, script)
}

func (p *Provider) writeArtifact(stage orchestrator.Stage, input orchestrator.StageInput, script Script) (orchestrator.StageOutput, error) {
	if script.ArtifactContent == "" {
		return orchestrator.StageOutput{}, nil
	}
	path := filepath.Join(input.WorkspaceDir, stage.ID+".md")
	if err := os.WriteFile(path, []byte(script.ArtifactContent), 0o644); err != nil {
		return orchestrator.StageOutput{}, fmt.Errorf("mock artifact write: %w", err)
	}
	return orchestrator.StageOutput{ArtifactPath: path}, nil
}

func (p *Provider) record(stageID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.invocations = append(p.invocations, stageID)
}

func (p *Provider) script(stageID string) Script {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.scripts[stageID]
}

// SetScript replaces one stage's script — used by resume tests to turn a
// hanging stage into a succeeding one between runs.
func (p *Provider) SetScript(stageID string, script Script) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.scripts[stageID] = script
}

// Invocations returns the stage IDs invoked so far, in order.
func (p *Provider) Invocations() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]string, len(p.invocations))
	copy(out, p.invocations)
	return out
}
