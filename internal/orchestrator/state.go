package orchestrator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// StateSchemaVersion identifies the run-state JSON shape. Bump on any
// incompatible change so a future reader can refuse or migrate old files.
const StateSchemaVersion = 4

// RunStateFileName is the executor-owned state file inside the feature
// workspace. NOTE: this lives beside the markdown pipeline's
// pipeline-state.json; roadmap L2.12 migrates pipeline-state.json semantics
// (checksum-verified resume, tamper refusal) into this file and retires the
// prompt-owned one.
const RunStateFileName = "run-state.json"

// StageStatus is the lifecycle state of one stage in a run.
type StageStatus string

// Stage lifecycle values. A stage is only skipped on resume when COMPLETED;
// RUNNING, INTERRUPTED, and FAILED stages are re-run.
const (
	StageStatusRunning     StageStatus = "RUNNING"
	StageStatusCompleted   StageStatus = "COMPLETED"
	StageStatusFailed      StageStatus = "FAILED"
	StageStatusInterrupted StageStatus = "INTERRUPTED"
	// StageStatusWaitingApproval marks a gated stage the executor refused to
	// start because its gate has no approval yet (roadmap L2.13).
	StageStatusWaitingApproval StageStatus = "WAITING_APPROVAL"
	// StageStatusStale marks a stage that completed once but whose artifact
	// no longer matches its recorded digest, or whose input went stale
	// (roadmap L2.12). A stale stage re-runs.
	StageStatusStale StageStatus = "STALE"
)

// StageRecord is the persisted result of one stage invocation.
type StageRecord struct {
	Status         StageStatus `json:"status"`
	StartedAt      time.Time   `json:"startedAt"`
	FinishedAt     *time.Time  `json:"finishedAt,omitempty"`
	ArtifactPath   string      `json:"artifactPath,omitempty"`
	ArtifactSHA256 string      `json:"artifactSha256,omitempty"`
	Gate           string      `json:"gate,omitempty"`
	Sequence       int         `json:"sequence"`
	PreviousStatus StageStatus `json:"previousStatus,omitempty"`
	StaleReason    StaleReason `json:"staleReason,omitempty"`
	FoundSHA256    string      `json:"foundSha256,omitempty"`
	Error          string      `json:"error,omitempty"`
}

// RunState is the executor-owned durable state for one feature delivery run.
type RunState struct {
	SchemaVersion int                    `json:"schemaVersion"`
	PlanName      string                 `json:"planName"`
	CreatedBy     Creator                `json:"createdBy"`
	FeatureName   string                 `json:"featureName,omitempty"`
	SpecPath      string                 `json:"specPath,omitempty"`
	StartedAt     time.Time              `json:"startedAt"`
	Stages        map[string]StageRecord `json:"stages"`
	Approvals     map[string]Approval    `json:"approvals"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

// Creator identifies which pipeline owns a state file. The two pipelines
// route differently — the markdown one skips conditional agents and loops
// developer/code-reviewer, the executor's plan is linear — so sharing a
// file is fine but resuming each other's runs is not.
type Creator string

// The two pipelines that write run state.
const (
	CreatedByExecutor Creator = "executor"
	CreatedByMarkdown Creator = "markdown"
)

// NewRunState returns an empty state for a fresh run of the named plan.
func NewRunState(planName string, createdBy Creator) *RunState {
	return &RunState{
		SchemaVersion: StateSchemaVersion,
		PlanName:      planName,
		CreatedBy:     createdBy,
		StartedAt:     time.Now().UTC(),
		Stages:        map[string]StageRecord{},
		Approvals:     map[string]Approval{},
	}
}

// CheckCreatedBy refuses a state file written by the other pipeline.
func (s *RunState) CheckCreatedBy(want Creator) error {
	if s.CreatedBy == want {
		return nil
	}
	return fmt.Errorf("this run state was written by the %s pipeline, not the %s one — the two route differently and cannot resume each other's runs",
		s.CreatedBy, want)
}

// NextSequence returns the sequence number for a stage record being created
// for the first time. Sequences are monotonic in recording order and are
// preserved when a stage is re-recorded (a re-run after a review loop is
// still the same position in the run), so they order stages for a pipeline
// that has no fixed plan to order by.
func (s *RunState) NextSequence() int {
	highest := 0
	for _, record := range s.Stages {
		if record.Sequence > highest {
			highest = record.Sequence
		}
	}
	return highest + 1
}

// StagesInSequence returns the recorded stage IDs in recording order.
func (s *RunState) StagesInSequence() []string {
	ids := make([]string, 0, len(s.Stages))
	for id := range s.Stages {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return s.Stages[ids[i]].Sequence < s.Stages[ids[j]].Sequence })
	return ids
}

// IsStageCompleted reports whether a stage already finished successfully and
// may be skipped on resume.
func (s *RunState) IsStageCompleted(stageID string) bool {
	return s.Stages[stageID].Status == StageStatusCompleted
}

// StateStore persists RunState to a single JSON file with atomic writes.
type StateStore struct {
	path string
}

// NewStateStore returns a store writing to the given file path.
func NewStateStore(path string) *StateStore {
	return &StateStore{path: path}
}

// Path returns the file the store reads and writes.
func (st *StateStore) Path() string { return st.path }

// Load reads persisted state. It returns (nil, nil) when no state file
// exists — a fresh run, not an error.
func (st *StateStore) Load() (*RunState, error) {
	raw, err := os.ReadFile(st.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read run state: %w", err)
	}
	return decodeRunState(raw, st.path)
}

func decodeRunState(raw []byte, path string) (*RunState, error) {
	var state RunState
	if err := json.Unmarshal(raw, &state); err != nil {
		return nil, fmt.Errorf("parse run state %s: %w", path, err)
	}
	if state.SchemaVersion != StateSchemaVersion {
		return nil, fmt.Errorf("run state %s has schema version %d, this executor supports %d — finish or delete the old run",
			path, state.SchemaVersion, StateSchemaVersion)
	}
	if state.Stages == nil {
		state.Stages = map[string]StageRecord{}
	}
	if state.Approvals == nil {
		state.Approvals = map[string]Approval{}
	}
	return &state, nil
}

// Save persists state atomically: it writes a temp file in the same
// directory, then os.Rename over the target. A crash between the temp write
// and the rename leaves the previous state file intact — readers never see
// a partial JSON document.
func (st *StateStore) Save(state *RunState) error {
	state.UpdatedAt = time.Now().UTC()
	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode run state: %w", err)
	}
	return st.writeAtomic(raw)
}

func (st *StateStore) writeAtomic(raw []byte) error {
	dir := filepath.Dir(st.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}
	tmp, err := os.CreateTemp(dir, RunStateFileName+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp state file: %w", err)
	}
	return st.commitTemp(tmp, raw)
}

func (st *StateStore) commitTemp(tmp *os.File, raw []byte) error {
	tmpName := tmp.Name()
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write temp state file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp state file: %w", err)
	}
	if err := os.Rename(tmpName, st.path); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename temp state file into place: %w", err)
	}
	return nil
}

// FeatureNameFromSpec derives the feature name from its spec file, matching
// the workspace convention (features/user-auth.md -> user-auth).
func FeatureNameFromSpec(specPath string) string {
	base := filepath.Base(specPath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// ArtifactSHA256 computes the hex SHA-256 of an artifact file, in Go —
// never delegated to a prompt (roadmap L2.12 alignment).
func ArtifactSHA256(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read artifact for checksum: %w", err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
