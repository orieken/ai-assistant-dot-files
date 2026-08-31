package mock

// Typed stages (roadmap L2.9) return validated state documents rather than
// markdown, so the mock needs a valid document per typed stage to script a
// deterministic run.

import (
	"encoding/json"

	"github.com/orieken/loom/internal/state"
)

// TypedScript returns a scripted payload for a state kind, and false for a
// stage that still exchanges markdown.
func TypedScript(kind string) (Script, bool) {
	payload, typed := typedPayload(state.Kind(kind))
	if !typed {
		return Script{}, false
	}
	return Script{Payload: payload}, true
}

func typedPayload(kind state.Kind) ([]byte, bool) {
	switch kind {
	case state.KindAnalysis:
		return mustEncode(sampleAnalysis()), true
	case state.KindArchitecture:
		return mustEncode(sampleArchitecture()), true
	default:
		return nil, false
	}
}

// mustEncode panics only on a programming error: these values are compiled
// in, so a failure means the structs and this file disagree, which a test
// catches immediately.
func mustEncode(value interface{}) []byte {
	raw, err := json.Marshal(value)
	if err != nil {
		panic("mock: cannot encode scripted state: " + err.Error())
	}
	return raw
}

func sampleAnalysis() state.AnalysisState {
	return state.AnalysisState{
		SchemaVersion:      state.SchemaVersion,
		Feature:            "mock-feature",
		Summary:            "Scripted analysis produced by the mock provider.",
		AcceptanceCriteria: []state.AcceptanceCriterion{{Statement: "the mocked behaviour happens"}},
		BoundedContext:     state.BoundedContext{Owning: "mock"},
		AffectedComponents: []state.AffectedComponent{{Path: "internal/mock/thing.go", Reason: "scripted"}},
		Tasks:              state.TaskList{Developer: []string{"implement the mocked behaviour"}},
		DefinitionOfDone:   []string{"tests green"},
	}
}

func sampleArchitecture() state.ArchitectureState {
	fitness := state.FitnessFunction{Property: "mock property", Verification: "go test ./..."}
	return state.ArchitectureState{
		SchemaVersion: state.SchemaVersion,
		Feature:       "mock-feature",
		StructuralDecisions: []state.StructuralDecision{
			{Decision: "keep the mock in the adapter layer", Rationale: "scripted", Fitness: &fitness},
		},
		BoundedContext:   state.BoundedContext{Owning: "mock"},
		FitnessFunctions: []state.FitnessFunction{fitness},
	}
}
