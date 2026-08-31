package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/spf13/cobra"
)

// errStateVerifyFailed makes `loom state verify` exit non-zero without
// cobra printing a redundant error line — the per-stage report is the
// message.
var errStateVerifyFailed = errors.New("recorded artifacts failed verification")

func runStateRecord(cmd *cobra.Command, _ []string) error {
	store, state, err := openMarkdownState(stateArgs.spec)
	if err != nil {
		return err
	}
	record, err := completedRecord(state, stateArgs.stage, stateArgs.artifact)
	if err != nil {
		return err
	}
	state.Stages[stateArgs.stage] = record
	if err := store.Save(state); err != nil {
		return err
	}
	cmd.Printf("recorded %s #%d%s\n", stateArgs.stage, record.Sequence, digestSuffix(record))
	return nil
}

// completedRecord hashes the artifact here, in Go. There is deliberately no
// flag for supplying a digest: a caller that could hand us a hash could
// hand us a wrong one, which is the whole failure mode being closed.
func completedRecord(state *orchestrator.RunState, stageID, artifact string) (orchestrator.StageRecord, error) {
	record := state.Stages[stageID]
	record.Status = orchestrator.StageStatusCompleted
	record.StaleReason = ""
	record.FoundSHA256 = ""
	if record.Sequence == 0 {
		record.Sequence = state.NextSequence()
	}
	if record.StartedAt.IsZero() {
		record.StartedAt = time.Now().UTC()
	}
	now := time.Now().UTC()
	record.FinishedAt = &now
	return withArtifactDigest(record, artifact)
}

func withArtifactDigest(record orchestrator.StageRecord, artifact string) (orchestrator.StageRecord, error) {
	if artifact == "" {
		return record, nil
	}
	absArtifact, err := filepath.Abs(artifact)
	if err != nil {
		return record, fmt.Errorf("resolve artifact path: %w", err)
	}
	sum, err := orchestrator.ArtifactSHA256(absArtifact)
	if err != nil {
		return record, fmt.Errorf("hash artifact: %w", err)
	}
	record.ArtifactPath = absArtifact
	record.ArtifactSHA256 = sum
	return record, nil
}

func digestSuffix(record orchestrator.StageRecord) string {
	if record.ArtifactSHA256 == "" {
		return " (no artifact)"
	}
	return " sha256:" + record.ArtifactSHA256[:12]
}

func runStateVerify(cmd *cobra.Command, _ []string) error {
	store, state, err := requireExistingState(stateArgs.spec)
	if err != nil {
		return err
	}
	stale := orchestrator.VerifyCompletedStages(state)
	if err := store.Save(state); err != nil {
		return err
	}
	printVerifyReport(cmd, state, stale)
	if len(stale) > 0 {
		return errStateVerifyFailed
	}
	return nil
}

func printVerifyReport(cmd *cobra.Command, state *orchestrator.RunState, stale []orchestrator.StaleStage) {
	demoted := make(map[string]orchestrator.StaleReason, len(stale))
	for _, item := range stale {
		demoted[item.StageID] = item.Reason
	}
	for _, stageID := range state.StagesInSequence() {
		if reason, ok := demoted[stageID]; ok {
			cmd.Printf("%-28s %s\n", stageID, reason)
			continue
		}
		cmd.Printf("%-28s OK\n", stageID)
	}
}

func runStateApprove(cmd *cobra.Command, _ []string) error {
	store, state, err := openMarkdownState(stateArgs.spec)
	if err != nil {
		return err
	}
	state.RecordApproval(stateArgs.gate, orchestrator.ApprovalMethodCLI)
	if err := store.Save(state); err != nil {
		return err
	}
	cmd.Printf("recorded approval for gate %q\n", stateArgs.gate)
	cmd.PrintErrln("note: this is an audit record — gate enforcement applies to `loom run` stages only")
	return nil
}

func runStateShow(cmd *cobra.Command, _ []string) error {
	_, state, err := requireExistingState(stateArgs.spec)
	if err != nil {
		return err
	}
	if stateArgs.asJSON {
		return printStateJSON(cmd, state)
	}
	printStateReport(cmd, state)
	return nil
}

func printStateJSON(cmd *cobra.Command, state *orchestrator.RunState) error {
	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	cmd.Println(string(raw))
	return nil
}

func printStateReport(cmd *cobra.Command, state *orchestrator.RunState) {
	cmd.Printf("%s (%s pipeline, started %s)\n", state.FeatureName, state.CreatedBy, state.StartedAt.Format(time.RFC3339))
	for _, stageID := range state.StagesInSequence() {
		record := state.Stages[stageID]
		cmd.Printf("  %2d. %-28s %s\n", record.Sequence, stageID, record.Status)
	}
	printApprovals(cmd, state)
}

func printApprovals(cmd *cobra.Command, state *orchestrator.RunState) {
	if len(state.Approvals) == 0 {
		cmd.Println("  approvals: none recorded")
		return
	}
	cmd.Println("  approvals:")
	for gate, approval := range state.Approvals {
		cmd.Printf("    %-20s %s by %s (%s)\n", gate, approval.ApprovedAt.Format(time.RFC3339), approval.Approver, approval.Method)
	}
}
