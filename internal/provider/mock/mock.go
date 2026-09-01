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
	// Payload is the typed state document a typed stage returns (roadmap
	// L2.9). The executor validates it and writes it as the artifact, so a
	// script setting Payload leaves ArtifactContent empty.
	Payload []byte
	// Usage is what this invocation reports consuming (roadmap L3.8). Nil
	// is the default and the honest one for a mock: it consumed nothing
	// from a model, and reporting zeros would assert a measurement that
	// was never taken.
	Usage *orchestrator.Usage
}

// Provider is a scripted, deterministic orchestrator.Provider.
type Provider struct {
	mu          sync.Mutex
	scripts     map[string]Script
	invocations []string
	inputs      map[string]orchestrator.StageInput
	hook        func(stageID string, invocation int) *Script
	visits      map[string]int
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
	p.recordInput(stage.ID, input)
	script := p.scriptFor(stage.ID)
	if script.Hang {
		<-ctx.Done()
		return orchestrator.StageOutput{}, ctx.Err()
	}
	if script.Err != nil {
		return orchestrator.StageOutput{Usage: script.Usage}, script.Err
	}
	if len(script.Payload) > 0 {
		return orchestrator.StageOutput{Payload: script.Payload, Usage: script.Usage}, nil
	}
	return p.writeArtifact(stage, input, script)
}

func (p *Provider) writeArtifact(stage orchestrator.Stage, input orchestrator.StageInput, script Script) (orchestrator.StageOutput, error) {
	if script.ArtifactContent == "" {
		return orchestrator.StageOutput{Usage: script.Usage}, nil
	}
	path := filepath.Join(input.WorkspaceDir, stage.ID+".md")
	if err := os.WriteFile(path, []byte(script.ArtifactContent), 0o644); err != nil {
		return orchestrator.StageOutput{}, fmt.Errorf("mock artifact write: %w", err)
	}
	return orchestrator.StageOutput{ArtifactPath: path, Usage: script.Usage}, nil
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

// SetHook lets a test vary a stage's script by how many times it has been
// invoked — which is how a loop test scripts "reject twice, then approve".
// Returning nil falls back to the stage's standing script.
func (p *Provider) SetHook(hook func(stageID string, invocation int) *Script) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.hook = hook
}

// scriptFor consults the hook first, then the standing script.
func (p *Provider) scriptFor(stageID string) Script {
	p.mu.Lock()
	hook := p.hook
	if p.visits == nil {
		p.visits = map[string]int{}
	}
	p.visits[stageID]++
	visit := p.visits[stageID]
	p.mu.Unlock()
	if hook != nil {
		if scripted := hook(stageID, visit); scripted != nil {
			return *scripted
		}
	}
	return p.script(stageID)
}

// SetScript replaces one stage's script — used by resume tests to turn a
// hanging stage into a succeeding one between runs.
func (p *Provider) SetScript(stageID string, script Script) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.scripts[stageID] = script
}

// InputFor returns the input a stage was invoked with, so tests can assert
// what a stage actually received — the projected upstream state, in
// particular.
func (p *Provider) InputFor(stageID string) orchestrator.StageInput {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.inputs[stageID]
}

func (p *Provider) recordInput(stageID string, input orchestrator.StageInput) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.inputs == nil {
		p.inputs = map[string]orchestrator.StageInput{}
	}
	p.inputs[stageID] = input
}

// Invocations returns the stage IDs invoked so far, in order.
func (p *Provider) Invocations() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]string, len(p.invocations))
	copy(out, p.invocations)
	return out
}
