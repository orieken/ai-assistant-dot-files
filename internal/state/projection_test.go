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
	projected, err := state.ProjectionFor("architect", state.KindAnalysis, analysisJSON(t))
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
	projected, err := state.ProjectionFor("architect", state.KindAnalysis, analysisJSON(t))
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
	if _, err := state.ProjectionFor("architect", state.KindArchitecture, analysisJSON(t)); err == nil {
		t.Error("a projection was invented for an undeclared hop")
	}
	if _, err := state.ProjectionFor("sre-engineer", state.KindAnalysis, analysisJSON(t)); err == nil {
		t.Error("a projection was invented for a stage that declares none")
	}
}

func TestProjectionRefusesInvalidUpstreamState(t *testing.T) {
	_, err := state.ProjectionFor("architect", state.KindAnalysis, []byte(`{"schemaVersion":1}`))
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

// TestQAReadsDifferentSlicesOfEachUpstream is why projections are keyed by
// consumer as well as upstream: three sources, three different slices.
func TestQAReadsDifferentSlicesOfEachUpstream(t *testing.T) {
	analysis := analysisJSON(t)

	forQA, err := state.ProjectionFor("qa-engineer", state.KindAnalysis, analysis)
	if err != nil {
		t.Fatalf("qa-engineer from analysis: %v", err)
	}
	forArchitect, err := state.ProjectionFor("architect", state.KindAnalysis, analysis)
	if err != nil {
		t.Fatalf("architect from analysis: %v", err)
	}

	if string(forQA) == string(forArchitect) {
		t.Error("qa-engineer and architect received the same slice of the analysis")
	}
	assertProjectionOmits(t, forQA, "boundedContext", "dataModelChanges", "architecturalFlags")
	var received state.QAAcceptanceInput
	if err := json.Unmarshal(forQA, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	if len(received.AcceptanceCriteria) == 0 {
		t.Error("qa-engineer received no acceptance criteria — the thing it tests against")
	}
}

// TestTechWriterReadsScopeNotTheWholeAnalysis pins the second
// summarize-artifact replacement: selected fields, not a paraphrase.
func TestTechWriterReadsScopeNotTheWholeAnalysis(t *testing.T) {
	projected, err := state.ProjectionFor("tech-writer", state.KindAnalysis, analysisJSON(t))
	if err != nil {
		t.Fatalf("ProjectionFor: %v", err)
	}

	var received state.TechWriterScopeInput
	if err := json.Unmarshal(projected, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	if received.Summary != validAnalysis().Summary {
		t.Error("the summary was altered in transit; a projection selects fields, it does not summarise")
	}
	assertProjectionOmits(t, projected, "acceptanceCriteria", "affectedComponents", "definitionOfDone")
}

func TestSecurityAndQAReadDifferentSlicesOfTheImplementation(t *testing.T) {
	raw, err := json.Marshal(validImplementation())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	forSecurity, err := state.ProjectionFor("security-reviewer", state.KindImplementation, raw)
	if err != nil {
		t.Fatalf("security-reviewer: %v", err)
	}
	forQA, err := state.ProjectionFor("qa-engineer", state.KindImplementation, raw)
	if err != nil {
		t.Fatalf("qa-engineer: %v", err)
	}

	// Security wants the surface that came in; QA wants what to test and
	// where the build departed from the plan. Neither gets the developer's
	// self-review checklist, which belongs to the code-reviewer.
	assertProjectionOmits(t, forSecurity, "selfReview", "refactoringLog", "notesForQa")
	assertProjectionOmits(t, forQA, "selfReview", "interfaceDesign")
}

// TestEveryProjectionRefusesAMalformedUpstream keeps a bad upstream
// document from reaching a consumer as half a projection.
func TestEveryProjectionRefusesAMalformedUpstream(t *testing.T) {
	cases := []struct {
		consumer string
		upstream state.Kind
	}{
		{"architect", state.KindAnalysis},
		{"developer", state.KindReview},
		{"security-reviewer", state.KindImplementation},
		{"qa-engineer", state.KindImplementation},
		{"qa-engineer", state.KindSecurity},
		{"qa-engineer", state.KindAnalysis},
		{"tech-writer", state.KindQA},
		{"tech-writer", state.KindAnalysis},
	}
	for _, tc := range cases {
		t.Run(tc.consumer+"<-"+string(tc.upstream), func(t *testing.T) {
			if _, err := state.ProjectionFor(tc.consumer, tc.upstream, []byte(`{"schemaVersion":1}`)); err == nil {
				t.Error("a projection accepted an upstream document that fails validation")
			}
			if _, err := state.ProjectionFor(tc.consumer, tc.upstream, []byte(`not json`)); err == nil {
				t.Error("a projection accepted an upstream document that is not JSON")
			}
		})
	}
}

func TestSecurityFindingsReachQAWithTheirVerification(t *testing.T) {
	security := validSecurity()
	security.Findings = []state.SecurityFinding{{
		Severity: state.SeverityHigh, Title: "token in logs", Location: "x.go:1",
		Description: "logged at info", FixApplied: "redacted", Verification: "grep the log fixture",
	}}
	raw, err := json.Marshal(security)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	projected, err := state.ProjectionFor("qa-engineer", state.KindSecurity, raw)
	if err != nil {
		t.Fatalf("ProjectionFor: %v", err)
	}

	var received state.QASecurityInput
	if err := json.Unmarshal(projected, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	if len(received.Findings) != 1 || received.Findings[0].Verification == "" {
		t.Errorf("QA received %+v, want the finding with how to verify its fix", received.Findings)
	}
}

func TestTechWriterReadsWhatQAFoundSurprising(t *testing.T) {
	qa := validQA()
	qa.NotesForTechWriter = []string{"lockout resets after 15 minutes, not on next success"}
	qa.KnownGaps = []state.KnownGap{{Criterion: "concurrent sign-in", Reason: "no harness for it yet"}}
	raw, err := json.Marshal(qa)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	projected, err := state.ProjectionFor("tech-writer", state.KindQA, raw)
	if err != nil {
		t.Fatalf("ProjectionFor: %v", err)
	}

	var received state.TechWriterQAInput
	if err := json.Unmarshal(projected, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	if len(received.NotesForTechWriter) != 1 || len(received.KnownGaps) != 1 {
		t.Errorf("tech-writer received %+v, want the notes and the gaps", received)
	}
	assertProjectionOmits(t, projected, "testResults", "coverage", "testFilesCreated")
}
