package state_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/state"
)

// routableStages mirrors the built-in plan's routable set: the review
// stages are present but never skippable.
func routableStages() []state.RoutableStage {
	return []state.RoutableStage{
		{ID: "context-engineer"},
		{ID: "analyst"},
		{ID: "architect", Skippable: true},
		{ID: "performance-engineer", Skippable: true},
		{ID: "data-engineer", Skippable: true},
		{ID: "developer"},
		{ID: "code-reviewer"},
		{ID: "accessibility-engineer", Skippable: true},
		{ID: "security-reviewer"},
		{ID: "qa-engineer"},
		{ID: "visual-qa-engineer"},
		{ID: "sre-engineer"},
		{ID: "tech-writer"},
		{ID: "devops-engineer", Skippable: true},
	}
}

// minimalAnalysis is a feature that needs none of the optional stages: no
// crossing, no migration, no dependency, no threshold, no a11y requirement,
// no devops tasks.
func minimalAnalysis() state.AnalysisState {
	analysis := validAnalysis()
	analysis.BoundedContext.Crossings = nil
	analysis.DataModelChanges = nil
	analysis.NonFunctionalRequirements = nil
	analysis.NewDependencies = nil
	analysis.ArchitecturalFlags = nil
	analysis.Tasks = state.TaskList{Developer: []string{"do the thing"}}
	return analysis
}

func TestRouteSkipsWhatTheAnalysisDoesNotAskFor(t *testing.T) {
	route := state.RouteFor(minimalAnalysis(), routableStages())

	for _, stage := range []string{"architect", "performance-engineer", "data-engineer", "accessibility-engineer", "devops-engineer"} {
		if route.Includes(stage) {
			t.Errorf("stage %q was included for a feature that asks for none of it", stage)
		}
		if route.ReasonFor(stage) == "" {
			t.Errorf("stage %q was skipped with no recorded reason", stage)
		}
	}
	if err := route.Validate(); err != nil {
		t.Errorf("route invalid: %v", err)
	}
}

// TestReviewStagesAreNeverSkipped is the safety property: automation gets
// the cheap half of the asymmetry, never the expensive half.
func TestReviewStagesAreNeverSkipped(t *testing.T) {
	route := state.RouteFor(minimalAnalysis(), routableStages())

	for _, stage := range []string{"code-reviewer", "security-reviewer", "qa-engineer", "developer", "tech-writer", "sre-engineer"} {
		if !route.Includes(stage) {
			t.Errorf("stage %q was skipped; it is not skippable by routing", stage)
		}
	}
}

// TestVisualQARunsBecauseItsConditionIsEnvironmental records the ADR-007
// boundary: whether heatmap data exists is not a fact about the feature, so
// no predicate can decide it and the stage stays in.
func TestVisualQARunsBecauseItsConditionIsEnvironmental(t *testing.T) {
	route := state.RouteFor(minimalAnalysis(), routableStages())

	if !route.Includes("visual-qa-engineer") {
		t.Error("visual-qa-engineer was routed out; its condition is environmental (ADR-007), not analytical")
	}
}

func TestRouteIncludesWhatTheAnalysisAsksFor(t *testing.T) {
	cases := map[string]struct {
		stage string
		shape func(*state.AnalysisState)
	}{
		"crossing summons the architect": {"architect", func(a *state.AnalysisState) {
			a.BoundedContext.Crossings = []string{"billing"}
		}},
		"threshold summons performance": {"performance-engineer", func(a *state.AnalysisState) {
			a.NonFunctionalRequirements = []state.NonFunctionalRequirement{
				{Category: "performance", Requirement: "fast", Threshold: "p99 < 200ms"},
			}
		}},
		"migration summons data": {"data-engineer", func(a *state.AnalysisState) {
			a.DataModelChanges = []state.DataModelChange{{Description: "add table", Phase: state.MigrationPhaseExpand}}
		}},
		"a11y requirement summons accessibility": {"accessibility-engineer", func(a *state.AnalysisState) {
			a.NonFunctionalRequirements = []state.NonFunctionalRequirement{
				{Category: "accessibility", Requirement: "keyboard navigable"},
			}
		}},
		"devops tasks summon devops": {"devops-engineer", func(a *state.AnalysisState) {
			a.Tasks.DevOps = []string{"add a CI job"}
		}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			analysis := minimalAnalysis()
			tc.shape(&analysis)
			route := state.RouteFor(analysis, routableStages())
			if !route.Includes(tc.stage) {
				t.Errorf("stage %q was skipped: %s", tc.stage, route.ReasonFor(tc.stage))
			}
		})
	}
}

func TestRouteIsDeterministic(t *testing.T) {
	first := state.RouteFor(minimalAnalysis(), routableStages())
	second := state.RouteFor(minimalAnalysis(), routableStages())

	if len(first.Decisions) != len(second.Decisions) {
		t.Fatalf("route length differs between runs")
	}
	for i := range first.Decisions {
		if first.Decisions[i] != second.Decisions[i] {
			t.Fatalf("decision %d differs between runs: %+v vs %+v", i, first.Decisions[i], second.Decisions[i])
		}
	}
}

func TestUnroutedStageRuns(t *testing.T) {
	route := state.RouteFor(minimalAnalysis(), routableStages())

	if !route.Includes("some-stage-nobody-routed") {
		t.Error("a stage the route does not mention should run; an unrouted stage is not a skipped one")
	}
}

func TestRouteValidationRequiresAReasonForEveryDecision(t *testing.T) {
	route := state.Route{
		SchemaVersion: state.SchemaVersion, Feature: "user-auth",
		Decisions: []state.RouteDecision{{Stage: "devops-engineer", Included: false}},
	}

	err := route.Validate()

	if err == nil || !strings.Contains(err.Error(), "reason") {
		t.Errorf("error = %v, want a refusal naming the missing reason", err)
	}
}

func TestRouteRendersTheDecisionAndItsConsequence(t *testing.T) {
	route := state.RouteFor(minimalAnalysis(), routableStages())
	raw, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	name, body, err := state.RenderView(state.KindRoute, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}

	if name != "route.md" {
		t.Errorf("rendered to %q, want route.md", name)
	}
	for _, want := range []string{"devops-engineer", "no DevOps tasks", "resets the design gate's approval"} {
		if !strings.Contains(body, want) {
			t.Errorf("rendered route missing %q:\n%s", want, body)
		}
	}
}

func TestRouteRoundTripsThroughItsSchema(t *testing.T) {
	raw, err := json.Marshal(state.RouteFor(minimalAnalysis(), routableStages()))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	decoded, err := state.Decode(state.KindRoute, raw)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	route, ok := decoded.(*state.Route)
	if !ok || len(route.Decisions) != len(routableStages()) {
		t.Errorf("decoded route = %+v", decoded)
	}
}
