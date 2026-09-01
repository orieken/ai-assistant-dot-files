package state

// QAState is the qa-engineer's output, modelling
// shared/contracts/qa-contract.md — including the rule that contract
// already asserts: `## Test Results` must show `Failed: 0`, "per the
// agent's own rule, tests must be green before the pipeline proceeds. A
// non-zero failure count is a FAIL, not a warning."
//
// As a grep over prose that rule could be satisfied by any line containing
// the right string. As a field it is a number.

import "fmt"

// TestResults is the run's outcome. Failed is the field the contract's rule
// turns on.
type TestResults struct {
	Passed  int `json:"passed" jsonschema:"required,minimum=0"`
	Failed  int `json:"failed" jsonschema:"minimum=0"`
	Skipped int `json:"skipped,omitempty" jsonschema:"minimum=0"`
	// SkipReasons explains any skipped tests; the template asks for a
	// reason and a count without one is not reviewable.
	SkipReasons []string `json:"skipReasons,omitempty"`
}

// CoverageSummary is what the qa-engineer measured. AcceptanceCriteriaCovered
// and AcceptanceCriteriaTotal are separate numbers rather than an "X/Y"
// string so a ratio can actually be computed from them.
type CoverageSummary struct {
	AcceptanceCriteriaCovered int     `json:"acceptanceCriteriaCovered" jsonschema:"minimum=0"`
	AcceptanceCriteriaTotal   int     `json:"acceptanceCriteriaTotal" jsonschema:"minimum=0"`
	NewTests                  int     `json:"newTests,omitempty" jsonschema:"minimum=0"`
	StatementCoveragePercent  float64 `json:"statementCoveragePercent,omitempty"`
}

// Bug is something QA found and what happened to it.
type Bug struct {
	Description string `json:"description" jsonschema:"required"`
	Resolution  string `json:"resolution" jsonschema:"required,description=How it was fixed, or why it was left"`
}

// KnownGap is an acceptance criterion that could not be tested, and why.
// This is deliberately structured: a gap without a reason is the kind of
// thing that disappears into prose.
type KnownGap struct {
	Criterion string `json:"criterion" jsonschema:"required"`
	Reason    string `json:"reason" jsonschema:"required"`
}

// QAState is the typed form of qa-report.md.
type QAState struct {
	SchemaVersion int    `json:"schemaVersion" jsonschema:"required"`
	Feature       string `json:"feature" jsonschema:"required"`

	TestFilesCreated  []string `json:"testFilesCreated,omitempty"`
	TestFilesModified []string `json:"testFilesModified,omitempty"`

	Coverage    CoverageSummary `json:"coverage" jsonschema:"required"`
	TestResults TestResults     `json:"testResults" jsonschema:"required"`

	AccessibilityCheck []string   `json:"accessibilityCheck,omitempty"`
	BugsFound          []Bug      `json:"bugsFound,omitempty"`
	KnownGaps          []KnownGap `json:"knownGaps,omitempty"`
	NotesForTechWriter []string   `json:"notesForTechWriter,omitempty"`

	Retrieval Retrieval `json:"retrieval,omitempty"`
}

// IsGreen reports whether the suite passed. The pipeline's "tests must be
// green" rule reads this rather than a rendered line.
func (q QAState) IsGreen() bool { return q.TestResults.Failed == 0 }

// Validate enforces the contract's rule in code: a non-zero failure count
// is a FAIL, not a warning.
func (q QAState) Validate() error {
	return firstError(
		requireSchemaVersion(q.SchemaVersion),
		requireText("feature", q.Feature),
		requireGreenSuite(q.TestResults),
		requireCoherentCoverage(q.Coverage),
	)
}

func requireGreenSuite(results TestResults) error {
	if results.Failed == 0 {
		return nil
	}
	return &ValidationError{Field: "testResults.failed",
		Reason: fmt.Sprintf("is %d — the suite must be green before the pipeline proceeds, and a failure count is not a warning", results.Failed)}
}

// requireCoherentCoverage catches a report claiming to have covered more
// acceptance criteria than the feature has.
func requireCoherentCoverage(coverage CoverageSummary) error {
	if coverage.AcceptanceCriteriaCovered <= coverage.AcceptanceCriteriaTotal {
		return nil
	}
	return &ValidationError{Field: "coverage.acceptanceCriteriaCovered",
		Reason: fmt.Sprintf("is %d of %d — more criteria covered than exist",
			coverage.AcceptanceCriteriaCovered, coverage.AcceptanceCriteriaTotal)}
}
