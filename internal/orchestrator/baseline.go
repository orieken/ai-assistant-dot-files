package orchestrator

// The presented baseline: the artifact state a human was last shown at a
// gate (roadmap L4.5).
//
// L2.14 records the digests an approval binds, which answers *that* an
// artifact changed and never *what* changed. Answering the second question
// needs the earlier content, and a digest is deliberately not reversible.
// So the content is retained, following the pattern L2.17 established for
// loop iterations: copy it aside, digest the copy, and record it in state,
// so an earlier state is verifiable rather than merely remembered.
//
// "Presented" rather than "approved" is the load-bearing word. The
// telemetry schema defines the correction signal as a checksum changing
// between gate-presented and gate-approved — a human halting at a gate,
// reading the artifact, fixing it, and then approving. No approval exists
// at the moment of that edit, so L2.14's invalidation cannot see it at all.
// Anchoring the baseline to what was *shown* covers that sequence and the
// approve-then-edit-then-reapprove one with the same comparison.

import (
	"fmt"
	"path/filepath"
	"time"
)

// ApprovedDirName holds the retained artifact state per gate, beside the
// loop's .iterations directory.
const ApprovedDirName = ".approved"

// BaselineReason says why a baseline was captured.
type BaselineReason string

// The two capture points. Both mean "this is what the human has seen":
// halting shows them the state, and approving accepts the state as it
// stands right then.
const (
	BaselinePresented BaselineReason = "PRESENTED"
	BaselineApproved  BaselineReason = "APPROVED"
)

// PresentedArtifact is one retained copy and its own digest.
type PresentedArtifact struct {
	Stage  string    `json:"stage"`
	Path   string    `json:"path"`
	SHA256 string    `json:"sha256"`
	At     time.Time `json:"at"`
}

// GateBaseline is everything a gate bound at the moment it was last shown
// to or accepted by a human.
type GateBaseline struct {
	At     time.Time      `json:"at"`
	Reason BaselineReason `json:"reason"`
	// Artifacts is keyed by producing stage ID, matching the keys of the
	// approval's ArtifactDigests so the two can be compared directly.
	Artifacts map[string]PresentedArtifact `json:"artifacts"`
}

// DigestFor returns the retained digest for a stage, and whether the
// baseline holds one at all.
func (b GateBaseline) DigestFor(stageID string) (string, bool) {
	artifact, ok := b.Artifacts[stageID]
	return artifact.SHA256, ok
}

// captureBaseline retains the artifacts a gate binds and records them
// against that gate.
//
// It is best-effort by design: a capture failure is reported and never
// fails the stage or the approval. This observes a human's action; it is
// not a control. A learning signal that can halt a delivery is a signal
// people turn off.
func (e *Executor) captureBaseline(state *RunState, gate string, reason BaselineReason) {
	if !shouldCapture(state, gate, reason) {
		return
	}
	baseline, err := e.retainGateArtifacts(state, gate, reason)
	if err != nil && e.onBaselineError != nil {
		e.onBaselineError(fmt.Errorf("gate %q: retaining what the human was shown: %w", gate, err))
	}
	if len(baseline.Artifacts) == 0 {
		return
	}
	state.Baselines[gate] = baseline
}

// shouldCapture decides whether this capture may replace what is already
// retained.
//
// A second PRESENTED capture must not overwrite the first. A run halted at
// a gate can be re-run and halt again — after the human has already edited
// the artifact — and re-capturing then would quietly replace the evidence
// with the edited content, so the correction the human made would compare
// clean and never be recorded. The first presentation is the one the human
// reacted to.
//
// An APPROVED capture always replaces: accepting the current state makes
// that state the new baseline, which is precisely what a *later* edit
// should be measured against.
func shouldCapture(state *RunState, gate string, reason BaselineReason) bool {
	if reason == BaselineApproved {
		return true
	}
	existing, ok := state.Baselines[gate]
	return !ok || existing.Reason != BaselinePresented
}

// retainGateArtifacts copies every artifact the gate binds. A failure on
// one artifact does not discard the ones that succeeded — a partial
// baseline still answers the question for the stages it covers.
func (e *Executor) retainGateArtifacts(state *RunState, gate string, reason BaselineReason) (GateBaseline, error) {
	baseline := GateBaseline{At: time.Now().UTC(), Reason: reason, Artifacts: map[string]PresentedArtifact{}}
	var failure error
	for stageID := range state.completedDigests() {
		retained, err := e.retainOne(state, gate, stageID)
		if err != nil {
			failure = err
			continue
		}
		if retained.Path == "" {
			continue
		}
		baseline.Artifacts[stageID] = retained
	}
	return baseline, failure
}

// correctionTarget is the file a human would edit to correct this stage.
//
// For a typed stage that is the rendered view, not the tracked state
// document: the view is what a person opens, and editing it cannot corrupt
// the run because the executor never reads it back. For an untyped stage
// the markdown IS the tracked artifact, so they are the same file — and
// editing it does have consequences (L2.12 re-runs the stage).
func correctionTarget(record StageRecord) string {
	if record.ViewPath != "" {
		return record.ViewPath
	}
	return record.ArtifactPath
}

func (e *Executor) retainOne(state *RunState, gate, stageID string) (PresentedArtifact, error) {
	source := correctionTarget(state.Stages[stageID])
	if source == "" {
		return PresentedArtifact{}, nil
	}
	target := e.baselinePath(gate, stageID, source)
	if err := copyFile(source, target); err != nil {
		return PresentedArtifact{}, fmt.Errorf("stage %q: %w", stageID, err)
	}
	sum, err := ArtifactSHA256(target)
	if err != nil {
		return PresentedArtifact{}, fmt.Errorf("stage %q: digest retained artifact: %w", stageID, err)
	}
	return PresentedArtifact{Stage: stageID, Path: target, SHA256: sum, At: time.Now().UTC()}, nil
}

// baselinePath puts each gate's retained state in its own directory, so two
// gates binding the same stage do not overwrite each other.
func (e *Executor) baselinePath(gate, stageID, artifactPath string) string {
	return filepath.Join(e.workspaceDir(), ApprovedDirName, gate, stageID+filepath.Ext(artifactPath))
}

// workspaceDir is the directory the run's state file lives in — the same
// derivation the timeline uses, so retained artifacts land beside the state
// they describe without the executor needing StageInput.
func (e *Executor) workspaceDir() string {
	return filepath.Dir(e.store.Path())
}
