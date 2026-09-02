package policy

// The policy document, and how a YAML file becomes one.

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// FileSuffix is the extension every policy file must carry.
const FileSuffix = ".policy.yaml"

// DefaultDir is where a project keeps its policies, relative to the project
// root.
const DefaultDir = ".claude/policies"

// ActionType is what a policy asks for when its condition holds.
type ActionType string

// The action types policy-schema.md declares.
const (
	ActionAutoApprove  ActionType = "auto-approve"
	ActionAutoReject   ActionType = "auto-reject"
	ActionRequireHuman ActionType = "require-human"
	ActionEscalate     ActionType = "escalate"
)

// Action is what happens when a policy matches.
type Action struct {
	Type       ActionType `yaml:"type"`
	Reason     string     `yaml:"reason"`
	EscalateTo string     `yaml:"escalateTo,omitempty"`
}

// Policy is one loaded, validated policy file.
type Policy struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description,omitempty"`
	Enabled     *bool     `yaml:"enabled,omitempty"`
	Matcher     Matcher   `yaml:"matcher"`
	Condition   Condition `yaml:"condition"`
	Action      Action    `yaml:"action"`
	// Source is the file it came from, so an error at evaluation time can
	// name something a human can open.
	Source string `yaml:"-"`
}

// IsEnabled defaults to true, matching the schema: `enabled` absent means
// the policy is live.
func (p Policy) IsEnabled() bool { return p.Enabled == nil || *p.Enabled }

// Matcher names the gates a policy watches.
type Matcher struct {
	Gate  GateID   `yaml:"gate,omitempty"`
	Gates []GateID `yaml:"gates,omitempty"`
}

// Watched returns every gate this matcher covers.
func (m Matcher) Watched() []GateID {
	if m.Gate != "" {
		return append([]GateID{m.Gate}, m.Gates...)
	}
	return m.Gates
}

// Parse decodes and validates one policy document.
//
// Decoding is strict: unknown fields are rejected, and duplicate mapping
// keys are an error rather than a last-one-wins silent overwrite. Both
// matter more here than in most parsers — a policy is an authorization
// document, and the shipped auto-approve-refactor example had a duplicate
// `filePaths` key that no parser had ever seen, so half of what it claimed
// to check was never going to be checked.
func Parse(source string, raw []byte) (Policy, error) {
	var document policyDocument
	decoder := yaml.NewDecoder(newReader(raw))
	decoder.KnownFields(true)
	if err := decoder.Decode(&document); err != nil {
		return Policy{}, fmt.Errorf("%s: %w", source, err)
	}
	policy, err := document.toPolicy(source)
	if err != nil {
		return Policy{}, fmt.Errorf("%s: %w", source, err)
	}
	return policy, nil
}

// policyDocument is the wire shape: the condition arrives as a raw node so
// it can be walked into typed checks rather than into a map of any.
type policyDocument struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description,omitempty"`
	Enabled     *bool     `yaml:"enabled,omitempty"`
	Matcher     Matcher   `yaml:"matcher"`
	Condition   yaml.Node `yaml:"condition"`
	Action      Action    `yaml:"action"`
}

func (d policyDocument) toPolicy(source string) (Policy, error) {
	condition, err := decodeCondition(&d.Condition)
	if err != nil {
		return Policy{}, err
	}
	policy := Policy{
		Name: d.Name, Version: d.Version, Description: d.Description, Enabled: d.Enabled,
		Matcher: d.Matcher, Condition: condition, Action: d.Action, Source: source,
	}
	return policy, policy.validate()
}

func (p Policy) validate() error {
	if p.Name == "" {
		return fmt.Errorf("policy has no name — every policy needs a unique kebab-case name")
	}
	if err := p.validateGates(); err != nil {
		return err
	}
	if p.Condition.IsEmpty() {
		return fmt.Errorf("policy %q has an empty condition — it would fire on every gate it watches", p.Name)
	}
	return p.validateAction()
}

func (p Policy) validateGates() error {
	watched := p.Matcher.Watched()
	if len(watched) == 0 {
		return fmt.Errorf("policy %q watches no gate — set matcher.gate or matcher.gates", p.Name)
	}
	for _, gate := range watched {
		if err := CheckGate(gate); err != nil {
			return fmt.Errorf("policy %q: %w", p.Name, err)
		}
	}
	return nil
}

func (p Policy) validateAction() error {
	switch p.Action.Type {
	case ActionAutoApprove, ActionAutoReject, ActionRequireHuman, ActionEscalate:
	default:
		return fmt.Errorf("policy %q has unknown action type %q", p.Name, p.Action.Type)
	}
	if p.Action.Type == ActionEscalate && p.Action.EscalateTo == "" {
		return fmt.Errorf("policy %q escalates but names no escalateTo target", p.Name)
	}
	if p.Action.Reason == "" {
		return fmt.Errorf("policy %q has no reason — a decision nobody can explain is not auditable", p.Name)
	}
	return nil
}
