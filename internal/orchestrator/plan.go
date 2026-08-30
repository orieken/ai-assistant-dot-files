// Package orchestrator is the minimal pipeline executor decided by ADR-006
// (loom executes pipelines) and roadmap item M0.4. It owns the run loop:
// load a plan, execute stages in order via a Provider, persist durable state,
// and stop. Routing, parallelism, gates, and policy are later roadmap items
// (L2.13, L3.x) that plug into this skeleton — they do not live here.
package orchestrator

import (
	"fmt"
	"time"
)

// Stage is one step of a Plan. Stages are identified by stable string IDs
// (agent names), never ordinals — hand-numbered step lists are a defect the
// roadmap (L2.15) explicitly retires.
type Stage struct {
	ID      string
	Agent   string
	Timeout time.Duration
}

// Plan is an ordered list of stages the executor runs sequentially.
type Plan struct {
	Name   string
	Stages []Stage
}

// Validate rejects plans the executor cannot run safely: empty names,
// stages without IDs, or duplicate stage IDs (state is keyed by stage ID,
// so duplicates would silently overwrite each other's records).
func (p Plan) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("plan has no name")
	}
	if len(p.Stages) == 0 {
		return fmt.Errorf("plan %q has no stages", p.Name)
	}
	return p.validateStageIDs()
}

func (p Plan) validateStageIDs() error {
	seen := make(map[string]bool, len(p.Stages))
	for _, stage := range p.Stages {
		if stage.ID == "" {
			return fmt.Errorf("plan %q has a stage with an empty ID", p.Name)
		}
		if seen[stage.ID] {
			return fmt.Errorf("plan %q has duplicate stage ID %q", p.Name, stage.ID)
		}
		seen[stage.ID] = true
	}
	return nil
}

// DefaultDeliverFeaturePlanName names the built-in plan.
const DefaultDeliverFeaturePlanName = "deliver-feature"

// defaultStageTimeout bounds each stage invocation. Every network or
// subprocess call must have an explicit timeout (architecture-guardrails.md
// #5); per-stage overrides belong in the plan, not in provider internals.
const defaultStageTimeout = 30 * time.Minute

// DefaultDeliverFeaturePlan returns the built-in plan encoding the existing
// linear agent sequence from shared/skills/deliver-feature/SKILL.md, in
// invocation order. Conditional-skip semantics (e.g. architect only when
// analysis.md flags architecture work) are routing, which this skeleton does
// not do — the linear order is preserved so behavior matches the markdown
// pipeline while the substrate changes underneath (roadmap M0.4).
func DefaultDeliverFeaturePlan() Plan {
	agents := []string{
		"context-engineer",
		"analyst",
		"architect",
		"performance-engineer",
		"data-engineer",
		"developer",
		"code-reviewer",
		"accessibility-engineer",
		"security-reviewer",
		"qa-engineer",
		"visual-qa-engineer",
		"sre-engineer",
		"tech-writer",
		"devops-engineer",
	}
	stages := make([]Stage, 0, len(agents))
	for _, agent := range agents {
		stages = append(stages, Stage{ID: agent, Agent: agent, Timeout: defaultStageTimeout})
	}
	return Plan{Name: DefaultDeliverFeaturePlanName, Stages: stages}
}
