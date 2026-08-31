package state

// Types shared by more than one stage's state document. They live here so
// the analyst and architect mean the same thing by the same word — the
// ubiquitous-language rule applied to the transport itself.

// FitnessFunction is a measurable architectural property and how CI checks
// it (architecture-guardrails.md #7). The analyst proposes them; the
// architect commits to them.
type FitnessFunction struct {
	Property     string `json:"property" jsonschema:"required,description=The property being measured"`
	Verification string `json:"verification" jsonschema:"required,description=How CI verifies it: tool, threshold, command"`
	Owner        string `json:"owner,omitempty" jsonschema:"description=Architect or DevOps"`
	JudgmentOnly bool   `json:"judgmentOnly,omitempty" jsonschema:"description=True when the property cannot be machine-verified and is documented as judgment-only"`
}

// BoundedContext is the DDD ownership of a feature and any boundary it
// crosses. A crossing is an architectural concern, so the architect reads
// this field directly.
type BoundedContext struct {
	Owning    string   `json:"owning" jsonschema:"required,description=The domain context that owns this feature"`
	Crossings []string `json:"crossings,omitempty" jsonschema:"description=Boundaries this feature crosses, empty when it stays inside its own context"`
}

// DomainEvents is the Event Storming Lite output.
type DomainEvents struct {
	Commands         []string `json:"commands,omitempty"`
	EventsProduced   []string `json:"eventsProduced,omitempty"`
	OwningAggregates []string `json:"owningAggregates,omitempty"`
	ReadModels       []string `json:"readModels,omitempty"`
}
