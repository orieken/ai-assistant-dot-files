package memory

// Reading a run's records from a workspace or a feature archive.
//
// Both locations are supported because the archived copy is the durable
// record: `.claude/feature-workspace/<feature>/` is temporary and gets
// cleaned, while `docs/features/<name>/` is committed. A store rebuilt from
// the archive after someone deletes episodes.db is the reason ingest is
// idempotent.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orieken/loom/internal/orchestrator"
)

// Records is one run's raw material.
type Records struct {
	State  *orchestrator.RunState
	Events []orchestrator.Event
	Source string
}

// ReadRecords loads a run's state and timeline from a directory holding
// them. A missing timeline is not an error — a run that halted before
// writing one still has state worth keeping, and refusing it would drop the
// runs a human most often looks at.
func ReadRecords(dir string) (Records, error) {
	statePath := filepath.Join(dir, orchestrator.RunStateFileName)
	state, err := orchestrator.NewStateStore(statePath).Load()
	if err != nil {
		return Records{}, fmt.Errorf("read run state: %w", err)
	}
	if state == nil {
		return Records{}, fmt.Errorf("no run state in %s — nothing to ingest", dir)
	}
	events, err := orchestrator.NewTimeline(statePath).Read()
	if err != nil {
		return Records{}, fmt.Errorf("read run events: %w", err)
	}
	return Records{State: state, Events: events, Source: dir}, nil
}

// ArchiveRecords copies a run's state and timeline into the feature
// archive, so the history is version-controlled beside the artifacts it
// describes and survives the workspace being cleaned.
//
// It is best-effort at its call sites: losing the archive copy costs
// durability, not the run.
func ArchiveRecords(workspaceDir, archiveDir string) error {
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return fmt.Errorf("create archive directory: %w", err)
	}
	for _, name := range []string{orchestrator.RunStateFileName, orchestrator.RunEventsFileName} {
		if err := copyIfPresent(filepath.Join(workspaceDir, name), filepath.Join(archiveDir, name)); err != nil {
			return err
		}
	}
	return nil
}

// copyIfPresent treats a missing source as nothing to do: a run with no
// timeline yet is a normal state, not a failure.
func copyIfPresent(source, target string) error {
	content, err := os.ReadFile(source)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", source, err)
	}
	if err := os.WriteFile(target, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", target, err)
	}
	return nil
}
