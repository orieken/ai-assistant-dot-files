package state_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/state"
)

func analysisJSON(t *testing.T) []byte {
	t.Helper()
	raw, err := json.Marshal(validAnalysis())
	if err != nil {
		t.Fatalf("marshal analysis: %v", err)
	}
	return raw
}

// TestProjectionGivesTheArchitectOnlyWhatItReads is the field-level access
// L2.9 exists for: the architect gets constraints and structural facts, not
// the QA task list or the definition of done.
func TestProjectionGivesTheArchitectOnlyWhatItReads(t *testing.T) {
	projected, err := state.ProjectFor(state.KindArchitecture, state.KindAnalysis, analysisJSON(t))
	if err != nil {
		t.Fatalf("ProjectFor: %v", err)
	}

	var received state.ArchitectInput
	if err := json.Unmarshal(projected, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	source := validAnalysis()
	if received.Feature != source.Feature || received.BoundedContext.Owning != source.BoundedContext.Owning {
		t.Errorf("projection dropped fields the architect needs: %+v", received)
	}
	if len(received.DeveloperTasks) != len(source.Tasks.Developer) {
		t.Errorf("developer tasks = %v, want the analyst's", received.DeveloperTasks)
	}
	assertProjectionOmits(t, projected, "definitionOfDone", "edgeCases", `"qa"`, "techWriter")
}

func assertProjectionOmits(t *testing.T, projected []byte, absent ...string) {
	t.Helper()
	for _, field := range absent {
		if strings.Contains(string(projected), field) {
			t.Errorf("projection leaked %q:\n%s", field, projected)
		}
	}
}

func TestProjectionIsFieldSelectionNotSummarization(t *testing.T) {
	source := validAnalysis()
	projected, err := state.ProjectFor(state.KindArchitecture, state.KindAnalysis, analysisJSON(t))
	if err != nil {
		t.Fatalf("ProjectFor: %v", err)
	}

	var received state.ArchitectInput
	if err := json.Unmarshal(projected, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	// Every retained field is byte-identical to its source: nothing on this
	// path rewrites, shortens, or paraphrases (the L2.10 premise).
	if received.Summary != source.Summary {
		t.Errorf("summary was altered in transit:\n got %q\nwant %q", received.Summary, source.Summary)
	}
	if received.AcceptanceCriteria[0].Statement != source.AcceptanceCriteria[0].Statement {
		t.Error("acceptance criteria were altered in transit")
	}
}

func TestProjectionRefusesUndeclaredHops(t *testing.T) {
	if _, err := state.ProjectFor(state.KindAnalysis, state.KindArchitecture, analysisJSON(t)); err == nil {
		t.Error("ProjectFor invented a projection for an undeclared hop")
	}
}

func TestProjectionRefusesInvalidUpstreamState(t *testing.T) {
	_, err := state.ProjectFor(state.KindArchitecture, state.KindAnalysis, []byte(`{"schemaVersion":1}`))
	if err == nil || !strings.Contains(err.Error(), "upstream") {
		t.Errorf("error = %v, want a refusal naming the upstream state", err)
	}
}

func TestDecodeRejectsUnknownKindAndInventedFields(t *testing.T) {
	if _, err := state.Decode(state.Kind("implementation"), analysisJSON(t)); err == nil {
		t.Error("Decode accepted a kind that is not typed yet")
	}
	invented := strings.Replace(string(analysisJSON(t)), `{"schemaVersion"`, `{"vibes":"immaculate","schemaVersion"`, 1)
	_, err := state.Decode(state.KindAnalysis, []byte(invented))
	if err == nil || !strings.Contains(err.Error(), "vibes") {
		t.Errorf("error = %v, want rejection naming the invented field", err)
	}
}

func TestDecodeAcceptsValidStateOfEachKind(t *testing.T) {
	if _, err := state.Decode(state.KindAnalysis, analysisJSON(t)); err != nil {
		t.Errorf("valid analysis rejected: %v", err)
	}
	raw, err := json.Marshal(validArchitecture())
	if err != nil {
		t.Fatalf("marshal architecture: %v", err)
	}
	if _, err := state.Decode(state.KindArchitecture, raw); err != nil {
		t.Errorf("valid architecture rejected: %v", err)
	}
}
