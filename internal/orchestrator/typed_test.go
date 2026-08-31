package orchestrator_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/provider/mock"
	"github.com/orieken/loom/internal/state"
)

// typedPlan is the L2.9 slice: a typed analyst whose state the typed
// architect consumes, followed by an untyped stage still on markdown.
func typedPlan() orchestrator.Plan {
	return orchestrator.Plan{
		Name: "typed-plan",
		Stages: []orchestrator.Stage{
			{ID: "analyst", Agent: "analyst", StateKind: string(state.KindAnalysis), Timeout: 5 * time.Second},
			{ID: "architect", Agent: "architect", StateKind: string(state.KindArchitecture),
				Consumes: "analyst", Timeout: 5 * time.Second},
			{ID: "developer", Agent: "developer", Timeout: 5 * time.Second},
		},
	}
}

func typedScripts(t *testing.T) map[string]mock.Script {
	t.Helper()
	analysis, ok := mock.TypedScript(string(state.KindAnalysis))
	if !ok {
		t.Fatal("no scripted analysis payload")
	}
	architecture, ok := mock.TypedScript(string(state.KindArchitecture))
	if !ok {
		t.Fatal("no scripted architecture payload")
	}
	return map[string]mock.Script{
		"analyst":   analysis,
		"architect": architecture,
		"developer": {ArtifactContent: "# implementation"},
	}
}

// TestTypedStagesExchangeDataWithNoMarkdownOnThePath is the L2.9 done-when:
// the architect receives the analyst's fields as validated data, and no
// markdown file exists for either stage.
func TestTypedStagesExchangeDataWithNoMarkdownOnThePath(t *testing.T) {
	scripts := typedScripts(t)
	executor, provider, store, input := newHarness(t, scripts)

	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertNoMarkdownFor(t, input, "analyst", "architect")
	received := provider.InputFor("architect")
	if received.UpstreamStage != "analyst" || len(received.UpstreamState) == 0 {
		t.Fatalf("architect received no upstream state: %+v", received)
	}
	assertProjectionContent(t, received.UpstreamState)
	assertTypedArtifactRecorded(t, store, input, "analyst")
}

func assertNoMarkdownFor(t *testing.T, input orchestrator.StageInput, stageIDs ...string) {
	t.Helper()
	for _, stageID := range stageIDs {
		if _, err := os.Stat(filepath.Join(input.WorkspaceDir, stageID+".md")); err == nil {
			t.Errorf("stage %q wrote a markdown file; typed stages exchange state, not documents", stageID)
		}
	}
}

// assertProjectionContent checks the architect got the analysis fields it
// declares — and none of the ones it does not. Passing the whole document
// would defeat the point of L2.9.
func assertProjectionContent(t *testing.T, projected []byte) {
	t.Helper()
	var received state.ArchitectInput
	if err := json.Unmarshal(projected, &received); err != nil {
		t.Fatalf("projection does not parse: %v", err)
	}
	if received.Feature == "" || len(received.AcceptanceCriteria) == 0 || received.BoundedContext.Owning == "" {
		t.Errorf("projection is missing fields the architect needs: %+v", received)
	}
	if strings.Contains(string(projected), "definitionOfDone") || strings.Contains(string(projected), "\"qa\"") {
		t.Errorf("projection leaked fields the architect does not read:\n%s", projected)
	}
}

func assertTypedArtifactRecorded(t *testing.T, store *orchestrator.StateStore, input orchestrator.StageInput, stageID string) {
	t.Helper()
	record := mustLoad(t, store).Stages[stageID]
	want := filepath.Join(input.WorkspaceDir, state.TypedStateDir, stageID+".json")
	if record.ArtifactPath != want {
		t.Errorf("stage %q artifact = %q, want its typed state document", stageID, record.ArtifactPath)
	}
	if record.ArtifactSHA256 == "" {
		t.Error("typed state was not digest-recorded; L2.12 integrity must cover it")
	}
}

func TestTypedStageRejectsInvalidPayload(t *testing.T) {
	scripts := typedScripts(t)
	scripts["analyst"] = mock.Script{Payload: []byte(`{"schemaVersion":1,"feature":"x"}`)}
	executor, provider, store, input := newHarness(t, scripts)

	err := executor.Run(context.Background(), typedPlan(), input)

	if err == nil || !strings.Contains(err.Error(), "summary") {
		t.Fatalf("Run error = %v, want a validation failure naming the missing field", err)
	}
	assertStatus(t, mustLoad(t, store), "analyst", orchestrator.StageStatusFailed)
	assertInvocations(t, provider, []string{"analyst"})
}

func TestTypedStageRejectsInventedFields(t *testing.T) {
	scripts := typedScripts(t)
	valid, _ := mock.TypedScript(string(state.KindAnalysis))
	withExtra := strings.Replace(string(valid.Payload), `{"schemaVersion"`, `{"vibes":"immaculate","schemaVersion"`, 1)
	scripts["analyst"] = mock.Script{Payload: []byte(withExtra)}
	executor, _, _, input := newHarness(t, scripts)

	err := executor.Run(context.Background(), typedPlan(), input)

	if err == nil || !strings.Contains(err.Error(), "vibes") {
		t.Fatalf("Run error = %v, want rejection naming the invented field", err)
	}
}

func TestTypedStageWithNoPayloadFails(t *testing.T) {
	scripts := typedScripts(t)
	scripts["analyst"] = mock.Script{ArtifactContent: "# analysis in markdown"}
	executor, _, _, input := newHarness(t, scripts)

	err := executor.Run(context.Background(), typedPlan(), input)

	if err == nil || !strings.Contains(err.Error(), "no state payload") {
		t.Fatalf("Run error = %v, want a failure about the missing payload", err)
	}
}

func TestEditingTypedStateDemotesTheStage(t *testing.T) {
	executor, provider, store, input := newHarness(t, typedScripts(t))
	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("first run: %v", err)
	}

	statePath := filepath.Join(input.WorkspaceDir, state.TypedStateDir, "analyst.json")
	if err := os.WriteFile(statePath, []byte(`{"hand":"edited"}`), 0o644); err != nil {
		t.Fatalf("edit typed state: %v", err)
	}
	var reported []orchestrator.StaleStage
	executor.OnStale(func(stale []orchestrator.StaleStage) { reported = stale })
	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("resume: %v", err)
	}

	if len(reported) == 0 || reported[0].StageID != "analyst" {
		t.Fatalf("stale report = %+v, want the edited analyst", reported)
	}
	assertStatus(t, mustLoad(t, store), "analyst", orchestrator.StageStatusCompleted)
	assertInvocations(t, provider, []string{"analyst", "architect", "developer", "analyst", "architect", "developer"})
}

func TestUntypedStagesAreUnaffected(t *testing.T) {
	executor, _, store, input := newHarness(t, typedScripts(t))
	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}

	record := mustLoad(t, store).Stages["developer"]
	if !strings.HasSuffix(record.ArtifactPath, "developer.md") {
		t.Errorf("untyped stage artifact = %q, want the markdown file it wrote", record.ArtifactPath)
	}
}

func TestDefaultPlanTypesTheAnalystArchitectHop(t *testing.T) {
	kinds := typedStagesOf(orchestrator.DefaultDeliverFeaturePlan())
	want := map[string]string{"analyst": string(state.KindAnalysis), "architect": string(state.KindArchitecture)}

	if len(kinds) != len(want) {
		t.Fatalf("typed stages = %v, want exactly %v in this cut", kinds, want)
	}
	for stageID, kind := range want {
		if kinds[stageID] != kind {
			t.Errorf("stage %q kind = %q, want %q", stageID, kinds[stageID], kind)
		}
	}
	assertArchitectConsumesAnalyst(t)
}

func typedStagesOf(plan orchestrator.Plan) map[string]string {
	kinds := map[string]string{}
	for _, stage := range plan.Stages {
		if stage.StateKind != "" {
			kinds[stage.ID] = stage.StateKind
		}
	}
	return kinds
}

func assertArchitectConsumesAnalyst(t *testing.T) {
	t.Helper()
	for _, stage := range orchestrator.DefaultDeliverFeaturePlan().Stages {
		if stage.ID == "architect" && stage.Consumes != "analyst" {
			t.Errorf("architect consumes %q, want analyst", stage.Consumes)
		}
	}
}

// TestTypedStageRendersTheContractsMarkdownView checks the half of L2.9
// that keeps the still-untyped stages working: they were told to read
// analysis.md, so that is where the view lands.
func TestTypedStageRendersTheContractsMarkdownView(t *testing.T) {
	executor, _, store, input := newHarness(t, typedScripts(t))
	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("Run: %v", err)
	}

	for _, view := range []string{"analysis.md", "architecture-notes.md"} {
		body, err := os.ReadFile(filepath.Join(input.WorkspaceDir, view))
		if err != nil {
			t.Fatalf("rendered view %s missing: %v", view, err)
		}
		if !strings.Contains(string(body), "## ") {
			t.Errorf("%s does not look like a rendered document:\n%s", view, body)
		}
	}
	// The view is derived, so it is not what integrity tracks — editing it
	// must not be able to corrupt a run.
	record := mustLoad(t, store).Stages["analyst"]
	if strings.HasSuffix(record.ArtifactPath, ".md") {
		t.Errorf("the rendered view is being tracked as the artifact: %q", record.ArtifactPath)
	}
}

func TestEditingTheRenderedViewDoesNotDemoteTheStage(t *testing.T) {
	executor, provider, _, input := newHarness(t, typedScripts(t))
	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("first run: %v", err)
	}

	if err := os.WriteFile(filepath.Join(input.WorkspaceDir, "analysis.md"), []byte("# scribbled over\n"), 0o644); err != nil {
		t.Fatalf("edit view: %v", err)
	}
	var reported []orchestrator.StaleStage
	executor.OnStale(func(stale []orchestrator.StaleStage) { reported = stale })
	if err := executor.Run(context.Background(), typedPlan(), input); err != nil {
		t.Fatalf("resume: %v", err)
	}

	if len(reported) != 0 {
		t.Errorf("editing a derived view demoted %+v; only state is tracked", reported)
	}
	assertInvocations(t, provider, []string{"analyst", "architect", "developer"})
}
