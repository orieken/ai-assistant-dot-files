package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/spf13/cobra"
)

// approvedRunState builds a run whose analyst stage completed with a real
// artifact and whose confirm-design gate is approved against it.
func approvedRunState(t *testing.T) (*orchestrator.Executor, string) {
	t.Helper()
	dir := t.TempDir()
	artifact := filepath.Join(dir, "analysis.md")
	if err := os.WriteFile(artifact, []byte("# analysis\n"), 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	digest, err := orchestrator.ArtifactSHA256(artifact)
	if err != nil {
		t.Fatalf("hash artifact: %v", err)
	}
	state := orchestrator.NewRunState("deliver-feature", orchestrator.CreatedByExecutor)
	state.Stages["analyst"] = orchestrator.StageRecord{
		Status: orchestrator.StageStatusCompleted, Sequence: 1, StartedAt: time.Now().UTC(),
		ArtifactPath: artifact, ArtifactSHA256: digest,
	}
	state.RecordApproval("confirm-design", orchestrator.ApprovalMethodTTY)

	store := orchestrator.NewStateStore(filepath.Join(dir, orchestrator.RunStateFileName))
	if err := store.Save(state); err != nil {
		t.Fatalf("save state: %v", err)
	}
	return orchestrator.NewExecutor(nil, store), artifact
}

func TestRefuseStaleApprovalPassesWhenNothingChanged(t *testing.T) {
	executor, _ := approvedRunState(t)

	if err := refuseStaleApproval(executor); err != nil {
		t.Errorf("refused an approval on an unchanged run: %v", err)
	}
}

func TestRefuseStaleApprovalNamesTheChangedArtifact(t *testing.T) {
	executor, artifact := approvedRunState(t)
	if err := os.WriteFile(artifact, []byte("# analysis, edited\n"), 0o644); err != nil {
		t.Fatalf("edit artifact: %v", err)
	}

	err := refuseStaleApproval(executor)

	if err == nil {
		t.Fatal("approving into a stale run was allowed")
	}
	for _, want := range []string{"cannot approve", "confirm-design", "analyst", "resume without --approve first"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("refusal missing %q: %v", want, err)
		}
	}
}

func TestApplyApproveFlagRequiresResume(t *testing.T) {
	executor, _ := approvedRunState(t)

	err := applyApproveFlag(executor, "confirm-design", false)

	if err == nil || !strings.Contains(err.Error(), "--resume") {
		t.Errorf("error = %v, want a refusal pointing at --resume", err)
	}
	if err := applyApproveFlag(executor, "", false); err != nil {
		t.Errorf("an absent --approve should be a no-op, got %v", err)
	}
}

func TestReportApprovalResetExplainsTheHalt(t *testing.T) {
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetErr(&output)

	reportApprovalReset(command, &orchestrator.StaleApprovalError{Gate: "confirm-ship", ChangedStage: "tech-writer"})

	for _, want := range []string{"confirm-ship", "tech-writer", "re-approve to continue"} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("reset message missing %q: %q", want, output.String())
		}
	}
}

func TestReportStaleStagesListsEveryDemotion(t *testing.T) {
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetErr(&output)

	reportStaleStages(command, []orchestrator.StaleStage{
		{StageID: "analyst", Reason: orchestrator.StaleReasonEdited},
		{StageID: "architect", Reason: orchestrator.StaleReasonUpstream},
	})

	if !strings.Contains(output.String(), "analyst") || !strings.Contains(output.String(), "architect") {
		t.Errorf("stale report = %q", output.String())
	}
}

func TestReportRouteSummarisesAndPointsAtTheDocument(t *testing.T) {
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetOut(&output)

	reportRoute(command, orchestrator.RouteSummary{
		Included: 9, Total: 14, Skipped: []string{"architect", "devops-engineer"},
	}, "/workspace")

	for _, want := range []string{"routed 9 of 14", "architect, devops-engineer", "/workspace/route.md"} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("route summary missing %q: %q", want, output.String())
		}
	}
}

func TestReportRouteSaysWhenNothingWasSkipped(t *testing.T) {
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetOut(&output)

	reportRoute(command, orchestrator.RouteSummary{Included: 14, Total: 14}, "/workspace")

	if !strings.Contains(output.String(), "nothing skipped") {
		t.Errorf("route summary = %q, want it to confirm the full plan runs", output.String())
	}
}

func TestReportLoopRoundSaysWhichRoundAndWhy(t *testing.T) {
	command := &cobra.Command{}
	var output bytes.Buffer
	command.SetOut(&output)

	reportLoopRound(command, orchestrator.LoopRound{
		Loop: "review", From: "developer", Iteration: 2, Max: 3,
	})

	for _, want := range []string{"review loop", "changes requested", "round 2 of 3", "developer"} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("loop message missing %q: %q", want, output.String())
		}
	}
}
