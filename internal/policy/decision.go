package policy

// What a gate's policies collectively decided, and why.
//
// This epic records a decision and does not act on it. `loom run` has never
// auto-approved a gate, and the first run that skips a barrier should not
// also be the first evidence the evaluator is correct. So Decision is an
// account of what WOULD have happened — which is exactly what the audit
// trail `approval-gates.md` promised and never had.

import (
	"fmt"
	"sort"
	"strings"
)

// Result is one policy's outcome at one gate. Named Result rather than
// PolicyOutcome because Outcome is already the tri-state a check resolves to,
// and two near-identically-named types in one package is how a caller reaches
// for the wrong one.
type Result struct {
	Name    string
	Action  ActionType
	Outcome Outcome
	Reason  string
	// Missing names the fields the policy tested that could not be
	// answered. A decision that says which fact it lacked is auditable;
	// one that only says "no match" is not.
	Missing []Field
	Source  string
}

// Matched reports whether this policy's condition held.
func (p Result) Matched() bool { return p.Outcome == OutcomeTrue }

// Decision is the resolved position for one gate.
type Decision struct {
	Gate GateID
	// Effect is what the policies collectively asked for. RequireHuman when
	// they conflict, or when nothing matched — the safe default is the one
	// the framework already has.
	Effect ActionType
	// Conflict names the policies that disagreed, empty when none did.
	Conflict []string
	Policies []Result
}

// WouldAutoApprove reports whether, had this epic honoured decisions, the
// gate would have been approved without a human. Nothing calls this to act
// yet — it exists so the recorded decision states its own consequence
// rather than leaving a reader to work it out.
func (d Decision) WouldAutoApprove() bool { return d.Effect == ActionAutoApprove }

// Summary is one line for a human at a halted gate.
func (d Decision) Summary() string {
	if len(d.Policies) == 0 {
		return "no policy watches this gate"
	}
	matched := d.matchedNames()
	if len(matched) == 0 {
		return fmt.Sprintf("%d policies evaluated, none matched", len(d.Policies))
	}
	return fmt.Sprintf("%s → %s (matched: %s)", d.Effect, d.effectNote(), strings.Join(matched, ", "))
}

func (d Decision) effectNote() string {
	if d.Effect == ActionAutoApprove {
		return "would auto-approve; this build still stops for a human"
	}
	return "stops for a human"
}

func (d Decision) matchedNames() []string {
	names := make([]string, 0, len(d.Policies))
	for _, outcome := range d.Policies {
		if outcome.Matched() {
			names = append(names, outcome.Name)
		}
	}
	sort.Strings(names)
	return names
}

// Decide evaluates every policy watching a gate and resolves the result.
//
// Conflict resolution is the rule the schema always stated and nothing ever
// applied: require-human beats auto-approve. So does auto-reject and
// escalate — anything that stops beats anything that proceeds, which is the
// only safe reading when two authorities disagree.
func Decide(policies []Policy, context GateContext) Decision {
	decision := Decision{Gate: context.Gate, Effect: ActionRequireHuman}
	for _, policy := range For(policies, context.Gate) {
		decision.Policies = append(decision.Policies, outcomeOf(policy, context))
	}
	decision.Effect, decision.Conflict = resolve(decision.Policies)
	return decision
}

func outcomeOf(policy Policy, context GateContext) Result {
	return Result{
		Name: policy.Name, Action: policy.Action.Type, Reason: policy.Action.Reason,
		Outcome: Evaluate(policy.Condition, context),
		Missing: MissingFacts(policy.Condition, context),
		Source:  policy.Source,
	}
}

// resolve returns the effect and any conflicting policy names. Nothing
// matching means require-human: the absence of a policy is not permission.
func resolve(outcomes []Result) (ActionType, []string) {
	approving, stopping := partitionMatched(outcomes)
	if len(approving) > 0 && len(stopping) > 0 {
		return ActionRequireHuman, append(approving, stopping...)
	}
	if len(approving) > 0 {
		return ActionAutoApprove, nil
	}
	return ActionRequireHuman, nil
}

func partitionMatched(outcomes []Result) (approving, stopping []string) {
	for _, outcome := range outcomes {
		if !outcome.Matched() {
			continue
		}
		if outcome.Action == ActionAutoApprove {
			approving = append(approving, outcome.Name)
			continue
		}
		stopping = append(stopping, outcome.Name)
	}
	sort.Strings(approving)
	sort.Strings(stopping)
	return approving, stopping
}
