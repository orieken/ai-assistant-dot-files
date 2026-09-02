package orchestrator

// The human-correction signal (roadmap L4.5).
//
// The telemetry schema has always specified `edited_then_approved` — "the
// human edited the artifact (checksum changed between gate-presented and
// gate-approved)" — and called it the richest corrective signal the system
// receives. Nothing has ever emitted it. This does.
//
// What a human corrects is the *rendered view*, not the tracked artifact.
// That is not a compromise: the view is the file a person opens, and
// because the executor never reads it back, an edit there cannot corrupt a
// run. The same property that made views safe (L2.9) makes them the right
// place to say "this should have said X". An edit to a tracked artifact is
// a different act with a different consequence — L2.12 re-runs the stage —
// so for the artifacts that are still markdown, the correction is captured
// before that re-run replaces the human's text with a second attempt.
//
// A view edit is ADVISORY. Nothing downstream adopts it and the next render
// overwrites the file. What survives is this record: the diff, and which
// agent's output needed correcting.

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CorrectionsDirName holds the retained diffs, beside the baselines.
const CorrectionsDirName = "corrections"

// Correction is one human edit to one stage's output.
type Correction struct {
	// Stage and Agent attribute the correction to whoever PRODUCED the
	// artifact, not to the gate that caught it. A gate name says where a
	// problem surfaced; the value of this signal is knowing whose output
	// needed fixing.
	Stage string    `json:"stage"`
	Agent string    `json:"agent"`
	Gate  string    `json:"gate"`
	At    time.Time `json:"at"`
	// DiffPath is the retained unified diff. Stat is the summary, kept on
	// the event so the timeline stays one greppable line per entry.
	DiffPath string   `json:"diffPath,omitempty"`
	Stat     DiffStat `json:"stat"`
}

// detectCorrections compares every retained baseline against what is on
// disk now, records a Correction for each divergence, and refreshes that
// baseline entry.
//
// Refreshing is what makes this safe to call from more than one place: a
// correction is reported exactly once, because after reporting it the
// baseline holds the corrected content. Without that, the next call would
// re-report the same edit — and worse, after a stage re-runs, would report
// the agent's own second attempt as a human correction.
func (e *Executor) detectCorrections(state *RunState) []Correction {
	corrections := make([]Correction, 0)
	for gate, baseline := range state.Baselines {
		corrections = append(corrections, e.correctionsAt(state, gate, baseline)...)
	}
	return corrections
}

func (e *Executor) correctionsAt(state *RunState, gate string, baseline GateBaseline) []Correction {
	found := make([]Correction, 0)
	for stageID, retained := range baseline.Artifacts {
		correction, ok := e.compareOne(state, gate, stageID, retained)
		if !ok {
			continue
		}
		found = append(found, correction)
		e.refreshBaselineEntry(state, gate, stageID)
	}
	return found
}

// compareOne diffs one stage's correction target against what was retained.
// A missing or unreadable file is not a correction — it is an absence, and
// reporting it as a human edit would be a fabrication.
func (e *Executor) compareOne(state *RunState, gate, stageID string, retained PresentedArtifact) (Correction, bool) {
	current, err := os.ReadFile(correctionTarget(state.Stages[stageID]))
	if err != nil {
		return Correction{}, false
	}
	before, err := os.ReadFile(retained.Path)
	if err != nil {
		return Correction{}, false
	}
	diff, stat := UnifiedDiff(string(before), string(current))
	if !stat.Changed() {
		return Correction{}, false
	}
	return Correction{
		Stage: stageID, Agent: state.Stages[stageID].Agent, Gate: gate, At: time.Now().UTC(),
		DiffPath: e.writeDiff(gate, stageID, diff), Stat: stat,
	}, true
}

// writeDiff retains the diff and returns its path, or "" when it could not
// be written. A lost diff must never fail an approval: this observes a
// human's action, it does not control the run.
func (e *Executor) writeDiff(gate, stageID, diff string) string {
	target := filepath.Join(e.workspaceDir(), ApprovedDirName, gate, CorrectionsDirName, stageID+".diff")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		e.reportBaselineError(fmt.Errorf("gate %q: create corrections directory: %w", gate, err))
		return ""
	}
	if err := os.WriteFile(target, []byte(diff), 0o644); err != nil {
		e.reportBaselineError(fmt.Errorf("gate %q stage %q: write correction diff: %w", gate, stageID, err))
		return ""
	}
	return target
}

// refreshBaselinesFor re-retains a stage in every baseline that holds it,
// called whenever the EXECUTOR writes that stage's output.
//
// Without this, a stage that re-runs after a human edit — which is exactly
// what L2.12 does to an edited markdown artifact — would leave the baseline
// holding the human's text while the file holds the agent's second attempt.
// The next comparison would then report the agent's own output as a human
// correction: a fabricated signal, and precisely the kind of confident
// wrong number this epic exists to avoid producing.
func (e *Executor) refreshBaselinesFor(state *RunState, stageID string) {
	for gate, baseline := range state.Baselines {
		if _, held := baseline.Artifacts[stageID]; !held {
			continue
		}
		e.refreshBaselineEntry(state, gate, stageID)
	}
}

// refreshBaselineEntry re-retains one stage so the reported correction is
// not reported again.
func (e *Executor) refreshBaselineEntry(state *RunState, gate, stageID string) {
	retained, err := e.retainOne(state, gate, stageID)
	if err != nil || retained.Path == "" {
		return
	}
	state.Baselines[gate].Artifacts[stageID] = retained
}

func (e *Executor) reportBaselineError(err error) {
	if e.onBaselineError != nil {
		e.onBaselineError(err)
	}
}

// recordCorrections persists the corrections and puts each on the timeline.
func (e *Executor) recordCorrections(state *RunState, corrections []Correction) error {
	for _, correction := range corrections {
		state.Corrections = append(state.Corrections, correction)
		if err := e.emit(correctionEvent(correction)); err != nil {
			return err
		}
	}
	return nil
}

func correctionEvent(correction Correction) Event {
	return Event{
		Kind: EventArtifactCorrected, Stage: correction.Stage, Gate: correction.Gate,
		Agent: correction.Agent, Correction: correction.Stat.String(), Reason: correction.DiffPath,
	}
}

// collectCorrections detects, records, and persists in one step, for
// callers that are not already saving state themselves.
func (e *Executor) collectCorrections(state *RunState) error {
	corrections := e.detectCorrections(state)
	if len(corrections) == 0 {
		return nil
	}
	if err := e.recordCorrections(state, corrections); err != nil {
		return err
	}
	return e.store.Save(state)
}
