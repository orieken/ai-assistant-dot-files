package state_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/state"
)

func renderedAnalysis(t *testing.T) (string, string) {
	t.Helper()
	raw, err := json.Marshal(validAnalysis())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	name, body, err := state.RenderView(state.KindAnalysis, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}
	return name, body
}

// TestRenderedAnalysisUsesTheContractFilenameAndHeadings matters because
// the untyped downstream stages were told to read analysis.md and to grep
// for these exact headings.
func TestRenderedAnalysisUsesTheContractFilenameAndHeadings(t *testing.T) {
	name, body := renderedAnalysis(t)

	if name != "analysis.md" {
		t.Errorf("rendered to %q, want the contract's analysis.md", name)
	}
	required := []string{
		"## Summary", "### Acceptance Criteria", "### Non-Functional Requirements",
		"## Proposed Fitness Functions", "## Out of Scope", "## Technical Breakdown",
		"### Bounded Context", "### Domain Events (Event Storming Lite)", "### Affected Components",
		"### Data Model Changes", "### API Changes", "### New Dependencies",
		"## Task List", "### Developer Tasks", "### QA Tasks", "### Tech Writer Tasks",
		"### DevOps Tasks", "## Edge Cases and Risks", "## Definition of Done",
	}
	for _, heading := range required {
		if !strings.Contains(body, "\n"+heading+"\n") {
			t.Errorf("rendered analysis is missing the contract heading %q", heading)
		}
	}
}

func TestRenderedAnalysisCarriesTheContent(t *testing.T) {
	_, body := renderedAnalysis(t)
	source := validAnalysis()

	for _, want := range []string{
		source.Summary,
		source.AcceptanceCriteria[0].Statement,
		source.AffectedComponents[0].Path,
		source.Tasks.Developer[0],
		"expand", // the migration phase survives as a visible fact
	} {
		if !strings.Contains(body, want) {
			t.Errorf("rendered analysis does not carry %q", want)
		}
	}
}

// TestEmptySectionsRenderAsNone keeps the distinction the contracts ask
// for: a section with nothing in it says so, rather than being blank.
func TestEmptySectionsRenderAsNone(t *testing.T) {
	analysis := validAnalysis()
	analysis.OutOfScope = nil
	analysis.APIChanges = nil
	raw, err := json.Marshal(analysis)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	_, body, err := state.RenderView(state.KindAnalysis, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}
	if !strings.Contains(body, "## Out of Scope\n\nNone") {
		t.Errorf("empty section did not render as None:\n%s", body)
	}
}

func TestRenderedArchitectureUsesTheContractFilenameAndHeadings(t *testing.T) {
	raw, err := json.Marshal(validArchitecture())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	name, body, err := state.RenderView(state.KindArchitecture, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}

	if name != "architecture-notes.md" {
		t.Errorf("rendered to %q, want the contract's architecture-notes.md", name)
	}
	for _, heading := range []string{
		"## Structural Decisions", "## Component Placement", "## Bounded Context",
		"## Stability Design", "## Observability Design", "## Layer Boundary Checks",
		"## Anti-Pattern Check", "## Fitness Functions", "## Refactoring Opportunities",
		"## Developer Handoff Notes", "## Open Architectural Questions",
	} {
		if !strings.Contains(body, "\n"+heading+"\n") {
			t.Errorf("rendered architecture is missing the contract heading %q", heading)
		}
	}
	// architecture-contract.md's validation rule looks for this line under
	// each decision.
	if !strings.Contains(body, "**Fitness Function**") {
		t.Error("a structural decision rendered without its Fitness Function line")
	}
}

func TestJudgmentOnlyDecisionRendersItsException(t *testing.T) {
	architecture := validArchitecture()
	architecture.StructuralDecisions[0].Fitness = nil
	architecture.StructuralDecisions[0].JudgmentOnly = true
	raw, err := json.Marshal(architecture)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	_, body, err := state.RenderView(state.KindArchitecture, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}
	if !strings.Contains(body, "judgment-only") {
		t.Errorf("judgment-only decision did not render its documented exception:\n%s", body)
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	first, second := "", ""
	for i := range 2 {
		_, body := renderedAnalysis(t)
		if i == 0 {
			first = body
			continue
		}
		second = body
	}
	if first != second {
		t.Error("rendering the same state twice produced different markdown")
	}
}

func TestRenderRefusesUnknownKindAndInvalidState(t *testing.T) {
	if _, _, err := state.RenderView(state.Kind("implementation"), []byte(`{}`)); err == nil {
		t.Error("rendered a kind that has no view")
	}
	if _, _, err := state.RenderView(state.KindAnalysis, []byte(`{"schemaVersion":1}`)); err == nil {
		t.Error("rendered a state document that fails validation")
	}
}

// TestRenderedViewCarriesRetrievalFrontmatter covers what validate-artifact
// checks (SKILL.md step 5) and what the retrieval corpus indexes on: all
// seven fields, always present.
func TestRenderedViewCarriesRetrievalFrontmatter(t *testing.T) {
	analysis := validAnalysis()
	analysis.Retrieval = state.Retrieval{
		DomainTerms: []string{"Session", "Credential"},
		IssueRefs:   []string{"PROJ-123"},
		LinkedADRs:  []string{"docs/adrs/ADR-006-loom-executes-pipelines.md"},
	}
	raw, err := json.Marshal(analysis)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	_, body, err := state.RenderView(state.KindAnalysis, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}

	if !strings.HasPrefix(body, "---\n") {
		t.Fatalf("rendered view does not start with a frontmatter block:\n%s", truncateForTest(body))
	}
	for _, field := range []string{"feature:", "bounded_context:", "domain_terms:", "files_touched:", "issue_refs:", "linked_adrs:", "linked_kis:"} {
		if !strings.Contains(body, "\n"+field) {
			t.Errorf("frontmatter is missing %q", field)
		}
	}
}

func TestFrontmatterDerivesWhatTheStateAlreadyKnows(t *testing.T) {
	raw, err := json.Marshal(validAnalysis())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	_, body, err := state.RenderView(state.KindAnalysis, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}

	source := validAnalysis()
	for _, want := range []string{
		`feature: "` + source.Feature + `"`,
		`bounded_context: "` + source.BoundedContext.Owning + `"`,
		`files_touched: ["` + source.AffectedComponents[0].Path + `"]`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("frontmatter missing derived value %q", want)
		}
	}
	// An agent that supplied nothing still gets well-formed empty lists —
	// a missing field is a WARN, an empty list is a fact.
	if !strings.Contains(body, "issue_refs: []") {
		t.Error("unsupplied retrieval fields should render as empty lists")
	}
}

func TestArchitectureFrontmatterUsesComponentPlacement(t *testing.T) {
	architecture := validArchitecture()
	architecture.ComponentPlacement = []state.ComponentPlacement{
		{Component: "SessionIssuer", Layer: "UseCases", Package: "internal/identity"},
	}
	raw, err := json.Marshal(architecture)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	_, body, err := state.RenderView(state.KindArchitecture, raw)
	if err != nil {
		t.Fatalf("RenderView: %v", err)
	}
	if !strings.Contains(body, `files_touched: ["internal/identity"]`) {
		t.Errorf("architecture frontmatter did not derive files_touched from placements:\n%s", truncateForTest(body))
	}
}

func truncateForTest(body string) string {
	if len(body) <= 400 {
		return body
	}
	return body[:400] + "…"
}
