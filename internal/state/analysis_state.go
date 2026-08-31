package state

// AnalysisState is the analyst's output, modelling
// shared/contracts/analysis-contract.md as typed fields rather than twenty
// markdown headings. Only what a downstream stage actually reads is
// modelled — inventing fields nobody consumes is how a schema rots.

// AcceptanceCriterion is one testable statement of done. Examples carry
// Specification by Example data the analyst was asked to make concrete.
type AcceptanceCriterion struct {
	Statement string   `json:"statement" jsonschema:"required,description=What must be true, in business language, never implementation"`
	Examples  []string `json:"examples,omitempty" jsonschema:"description=Concrete examples or data rows for complex rules"`
}

// NonFunctionalRequirement is a performance, security, or scaling
// constraint. Threshold is separate from the prose so a fitness function
// can be checked against it.
type NonFunctionalRequirement struct {
	Category    string `json:"category" jsonschema:"required,enum=performance,enum=security,enum=scaling,enum=accessibility,enum=other"`
	Requirement string `json:"requirement" jsonschema:"required"`
	Threshold   string `json:"threshold,omitempty" jsonschema:"description=The measurable limit, e.g. p99 < 200ms"`
}

// AffectedComponent is one file or module the feature touches.
type AffectedComponent struct {
	Path   string `json:"path" jsonschema:"required,description=Repo-relative path"`
	Reason string `json:"reason" jsonschema:"required"`
}

// MigrationPhase is the Expand/Contract phase of a data model change.
// Destructive work cannot ship in the same release as the code it supports
// (architecture-guardrails.md #2), so the phase is a field, not prose.
type MigrationPhase string

// The phases a data model change can be in.
const (
	MigrationPhaseNone     MigrationPhase = "none"
	MigrationPhaseExpand   MigrationPhase = "expand"
	MigrationPhaseContract MigrationPhase = "contract"
)

// DataModelChange is one schema change and when it may run.
type DataModelChange struct {
	Description string         `json:"description" jsonschema:"required"`
	Phase       MigrationPhase `json:"phase" jsonschema:"required,enum=none,enum=expand,enum=contract"`
}

// APIChange is one endpoint or signature change.
type APIChange struct {
	Endpoint string `json:"endpoint" jsonschema:"required"`
	Change   string `json:"change" jsonschema:"required"`
}

// TaskList is the per-role work the analyst broke the feature into. The
// roles are separate fields because each downstream agent reads only its
// own list — the projection L2.9 exists to make possible.
type TaskList struct {
	Developer  []string `json:"developer,omitempty"`
	QA         []string `json:"qa,omitempty"`
	TechWriter []string `json:"techWriter,omitempty"`
	DevOps     []string `json:"devops,omitempty"`
}

// EdgeCase is a boundary condition and how the feature handles it.
type EdgeCase struct {
	Case     string `json:"case" jsonschema:"required"`
	Handling string `json:"handling" jsonschema:"required"`
}

// AnalysisState is the typed form of analysis.md.
type AnalysisState struct {
	SchemaVersion int    `json:"schemaVersion" jsonschema:"required"`
	Feature       string `json:"feature" jsonschema:"required,description=kebab-case feature slug"`
	Summary       string `json:"summary" jsonschema:"required,description=One paragraph of what this does and why"`

	AcceptanceCriteria        []AcceptanceCriterion      `json:"acceptanceCriteria" jsonschema:"required"`
	NonFunctionalRequirements []NonFunctionalRequirement `json:"nonFunctionalRequirements,omitempty"`
	ProposedFitnessFunctions  []FitnessFunction          `json:"proposedFitnessFunctions,omitempty"`
	OutOfScope                []string                   `json:"outOfScope,omitempty"`

	BoundedContext     BoundedContext      `json:"boundedContext" jsonschema:"required"`
	DomainEvents       DomainEvents        `json:"domainEvents,omitempty"`
	AffectedComponents []AffectedComponent `json:"affectedComponents" jsonschema:"required"`
	DataModelChanges   []DataModelChange   `json:"dataModelChanges,omitempty"`
	APIChanges         []APIChange         `json:"apiChanges,omitempty"`
	NewDependencies    []string            `json:"newDependencies,omitempty"`

	// ArchitecturalFlags is what deliver-feature/SKILL.md step 12 routes on
	// ("invoke architect if analysis.md has Architectural Flags != None").
	// The markdown template never had a heading for it, so the prose
	// pipeline reads a field its own template does not produce. Typing it
	// makes the routing input explicit; acting on it is L3.1.
	ArchitecturalFlags []string `json:"architecturalFlags,omitempty"`

	Tasks            TaskList   `json:"tasks" jsonschema:"required"`
	EdgeCases        []EdgeCase `json:"edgeCases,omitempty"`
	DefinitionOfDone []string   `json:"definitionOfDone" jsonschema:"required"`
	OpenQuestions    []string   `json:"openQuestions,omitempty"`
}

// Validate enforces the fields a downstream stage cannot work without.
// Optional lists stay optional: an analyst legitimately has nothing to say
// about API changes on a docs-only feature.
func (a AnalysisState) Validate() error {
	return firstError(
		requireSchemaVersion(a.SchemaVersion),
		requireText("feature", a.Feature),
		requireText("summary", a.Summary),
		requireItems("acceptanceCriteria", len(a.AcceptanceCriteria)),
		requireText("boundedContext.owning", a.BoundedContext.Owning),
		requireItems("affectedComponents", len(a.AffectedComponents)),
		requireItems("definitionOfDone", len(a.DefinitionOfDone)),
	)
}
